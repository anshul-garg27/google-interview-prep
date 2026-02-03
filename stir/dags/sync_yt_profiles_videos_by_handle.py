from airflow import DAG
from airflow.hooks.base import BaseHook
from airflow.operators.python import PythonOperator
from airflow.utils.dates import days_ago
from slack_connection import SlackNotifier


def create_scrape_request_log():
    import os
    import json
    import requests
    import clickhouse_connect

    os.environ["no_proxy"] = "*"
    handles = []
    max_handles = 500

    connection = BaseHook.get_connection("clickhouse_gcc")
    extra_params = json.loads(connection.extra)
    client = clickhouse_connect.get_client(host=connection.host, password=connection.password, username=connection.login,
                                           send_receive_timeout=extra_params['send_receive_timeout'])
    priority_wise_queries = []

    priority_wise_queries.append('''
                                with 
                                    handles_in_progress as (
                                        select 
                                            JSONExtractString(params, 'channel_id') as handles
                                        from _e.scrape_request_log_events
                                        where flow = 'refresh_yt_posts_by_channel_id'
                                        and event_timestamp >= now() - INTERVAL 12 HOUR
                                        group by handles
                                    ),
                                    failures as (
                                        select 
                                            JSONExtractString(params, 'channel_id') as handle
                                        from dbt.stg_beat_scrape_request_log
                                        where 
                                            flow = 'refresh_yt_posts_by_channel_id'
                                            and data NOT LIKE '(204%'
                                            and data NOT LIKE '(409%'
                                            and data NOT LIKE '(429%'
                                            and status = 'FAILED'
                                            and created_at >= now() - INTERVAL 1 MONTH
                                        group by handle
                                        having 
                                            uniqExactIf(id, status = 'FAILED') > 3
                                    )
                                    select 
                                        handle
                                    from 
                                        dbt.mart_youtube_tracked_profiles profiles
                                    where 
                                        handle not in handles_in_progress
                                        and next_videos_attempt <= now()
                                        and handle not in failures
                                    order by priority asc, next_videos_attempt asc
                                ''')

    handles_so_far = 0
    for query in priority_wise_queries:
        handles_to_fetch = max_handles - handles_so_far
        if handles_to_fetch <= 0:
            break
        result = client.query(query + f" LIMIT {handles_to_fetch}")
        handles_so_far += len(result.result_rows)
        for handle in result.result_rows:
            handles.append(handle[0])

    refresh_post_flow = 'refresh_yt_posts_by_playlist_id'

    beat_connection = BaseHook.get_connection("beat")
    base_url = beat_connection.host

    refresh_post_url = base_url + "/scrape_request_log/flow/%s" % refresh_post_flow

    for handle in handles:
        refresh_post_scl = {"flow": refresh_post_flow, "platform": "YOUTUBE", "params": {"channel_id": handle, "max_posts": 50}}
        requests.post(url=refresh_post_url, data=json.dumps(refresh_post_scl))


with DAG(
        dag_id='sync_yt_profiles_videos_by_handle',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval="*/10 * * * *",
) as dag:
    task = PythonOperator(
        task_id='create_yt_scrape_request_log',
        python_callable=create_scrape_request_log
    )
    task
