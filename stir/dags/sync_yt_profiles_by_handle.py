from airflow import DAG
from airflow.hooks.base import BaseHook
from airflow.operators.python import PythonOperator
from airflow.utils.dates import days_ago
from slack_connection import SlackNotifier


def batchify(lst, max_len=50):
    result = []
    for i in range(0, len(lst), max_len):
        result.append(lst[i:i + max_len])
    return result


def create_scrape_request_log():
    import os
    import json
    import requests
    import clickhouse_connect

    os.environ["no_proxy"] = "*"
    handles = []
    max_handles = 100000

    connection = BaseHook.get_connection("clickhouse_gcc")
    extra_params = json.loads(connection.extra)
    client = clickhouse_connect.get_client(host=connection.host, password=connection.password, username=connection.login,
                                           send_receive_timeout=3600)
    priority_wise_queries = []

    priority_wise_queries.append('''
                                    with 
                                        handles_in_progress_raw as (
                                            select 
                                                JSONExtractArrayRaw(ifNull(params, ''), 'channel_ids') as handles
                                            from _e.scrape_request_log_events
                                            where flow = 'refresh_yt_profiles'
                                            and event_timestamp >= now() - INTERVAL 12 HOUR
                                            group by handles
                                        ),
                                        handles_in_progress as (
                                            select trim(BOTH '"' FROM arrayJoin(handles)) handle
                                            from handles_in_progress_raw
                                            group by handle
                                        )
                                        select 
                                            handle
                                        from 
                                            dbt.mart_youtube_tracked_profiles profiles
                                        where 
                                            handle not in handles_in_progress
                                            and next_attempt <= now()
                                        order by priority asc, next_attempt asc
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

    batches = batchify(handles, max_len=40)

    refresh_profiles_flow = 'refresh_yt_profiles'

    beat_connection = BaseHook.get_connection("beat")
    base_url = beat_connection.host

    refresh_profiles_url = base_url + "/scrape_request_log/flow/%s" % refresh_profiles_flow

    for batch in batches:
        refresh_profile_scl = {"flow": refresh_profiles_flow, "platform": "YOUTUBE", "params": {"channel_ids": batch}}
        requests.post(url=refresh_profiles_url, data=json.dumps(refresh_profile_scl))


with DAG(
        dag_id='sync_yt_profiles_by_handle',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval="*/30 * * * *",
) as dag:
    task = PythonOperator(
        task_id='create_yt_scrape_request_log',
        python_callable=create_scrape_request_log
    )
    task
