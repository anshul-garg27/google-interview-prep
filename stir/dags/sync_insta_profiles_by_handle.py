from airflow import DAG, AirflowException
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
    max_handles = 800

    connection = BaseHook.get_connection("clickhouse_gcc")
    extra_params = json.loads(connection.extra)
    client = clickhouse_connect.get_client(host=connection.host, password=connection.password, username=connection.login, send_receive_timeout=extra_params['send_receive_timeout'])
    handles_so_far = 0
    priority_wise_queries = []

    # App Handles (At least once)
    priority_wise_queries.append('''
        with handles_in_progress as (select JSONExtractString(params, 'handle') handle
                                     from _e.scrape_request_log_events
                                     where flow = 'refresh_profile_by_handle'
                                       and event_timestamp >= now() - INTERVAL 12 HOUR
                                     group by handle),
             failures as (select JSONExtractString(params, 'handle') handle
                          from dbt.stg_beat_scrape_request_log
                          where flow = 'refresh_profile_by_handle'
                            and data NOT LIKE '(204%'
                            and data NOT LIKE '(409%'
                            and data NOT LIKE '(429%'
                            and data NOT LIKE '%NoneType%'
                            and status = 'FAILED'
                            and created_at >= now() - INTERVAL 6 MONTH
                          group by handle
                          having uniqExactIf(id, status = 'FAILED') > 3)
        select handle, source
        from dbt.mart_instagram_tracked_profiles
        where handle not in handles_in_progress and handle not in failures 
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

    flow = 'refresh_profile_by_handle'
    beat_connection = BaseHook.get_connection("beat")
    base_url = beat_connection.host
    url = base_url + "/scrape_request_log/flow/%s" % flow
    handles = list(set(handles))
    errored = False
    for handle in handles:
        data = {"flow": flow, "platform": "INSTAGRAM", "params": {"handle": handle}}
        try:
            response = requests.post(url=url, data=json.dumps(data))
            response.raise_for_status()
        except requests.exceptions.HTTPError as e:  
    	    errored=True

    if errored:
        raise AirflowException()


with DAG(
        dag_id='sync_insta_profiles_by_handle',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval="*/10 * * * *",
) as dag:
    task = PythonOperator(
        task_id='create_insta_scrape_request_log',
        python_callable=create_scrape_request_log
    )
    task
