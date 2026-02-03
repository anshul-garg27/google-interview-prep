from airflow import DAG
from airflow.hooks.base import BaseHook
from airflow.operators.python import PythonOperator
from airflow.utils.dates import days_ago
from slack_connection import SlackNotifier


def sync_yt_profile_relationship_by_channel_id():
    import os
    import json
    import requests
    import clickhouse_connect

    os.environ["no_proxy"] = "*"
    handles = []
    max_handles = 1000

    connection = BaseHook.get_connection("clickhouse_gcc")
    extra_params = json.loads(connection.extra)
    client = clickhouse_connect.get_client(host=connection.host, password=connection.password, username=connection.login,
                                           send_receive_timeout=extra_params['send_receive_timeout'])
    print(client)
    query = '''
                with
                    vidooly_channel_id as (
                        select id
                        from vidooly.youtube_channels_india
                        group by id),
                    mart_profile_relationship_channel_id as (
                        select handle, last_activities_attempt
                        from dbt.mart_youtube_profile_relationship),
                    attempted_recently as (
                        select JSONExtractString(params, 'channel_id') channel_id 
                        from _e.scrape_request_log_events 
                        where flow = 'refresh_yt_profile_relationship' and event_timestamp>=now() - INTERVAL 1 DAY
                    )
                    select vci.id
                    from vidooly_channel_id vci left join mart_profile_relationship_channel_id mprci on vci.id = mprci.handle
                    where mprci.handle == '' and vci.id not in attempted_recently
                    order by mprci.last_activities_attempt asc
                '''

    
    result = client.query(query + f" LIMIT {max_handles}")
    for handle in result.result_rows:
        handles.append(handle[0])

    refresh_yt_profile_relationship = 'refresh_yt_profile_relationship'

    beat_connection = BaseHook.get_connection("beat")
    base_url = beat_connection.host

    refresh_post_url = base_url + "/scrape_request_log/flow/%s" % refresh_yt_profile_relationship

    for handle in handles:
        refresh_post_scl = {"flow": refresh_yt_profile_relationship, "platform": "YOUTUBE", "params": {"channel_id": handle}}
        requests.post(url=refresh_post_url, data=json.dumps(refresh_post_scl))


with DAG(
        dag_id='sync_yt_profile_relationship_by_channel_id',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval="*/10 * * * *",
) as dag:
    task = PythonOperator(
        task_id='sync_yt_profile_relationship_by_channel_id',
        python_callable=sync_yt_profile_relationship_by_channel_id
    )
    task
