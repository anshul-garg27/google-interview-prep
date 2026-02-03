from airflow import DAG
from airflow.hooks.base import BaseHook
from airflow.operators.python import PythonOperator
from airflow.utils.dates import days_ago
from airflow.models import Variable
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
    max_handles = 5000

    connection = BaseHook.get_connection("clickhouse_gcc")
    extra_params = json.loads(connection.extra)
    client = clickhouse_connect.get_client(host=connection.host, password=connection.password, username=connection.login, send_receive_timeout=extra_params['send_receive_timeout'])

    yt_handles = client.query('''with 
                                    candiate_set as (
                                        select 
                                            channel_id
                                        from dbt.stg_coffee_view_youtube_account_lite ya
                                        where ya.created_at >= now() - INTERVAL 6 MONTH 
                                        and followers = 0 and match(channel_id, '^UC[\w-]{21}[AQgw]$') = 1
                                    ),
                                    attempts_raw as (
                                        SELECT
                                            JSONExtractArrayRaw(ifNull(params, ''), 'channel_ids') as channel_ids
                                        from 
                                            _e.scrape_request_log_events as scrape
                                        where flow = 'refresh_yt_profiles' and event_timestamp >= now() - INTERVAL 1 DAY
                                    ),
                                    attempts as (
                                        select trim(BOTH '"' FROM arrayJoin(channel_ids)) channel_id
                                        from attempts_raw
                                        group by channel_id
                                    )
                                    select * from candiate_set where channel_id not in attempts
                                '''
                            f"LIMIT {max_handles}")

    for handle in yt_handles.result_rows:
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
        dag_id='sync_yt_channels',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval="0 0 * * *",
) as dag:
    task = PythonOperator(
        task_id='create_yt_channel_scl',
        python_callable=create_scrape_request_log
    )
    task
