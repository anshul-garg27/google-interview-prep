from airflow import DAG
from airflow.hooks.base import BaseHook
from airflow.operators.python import PythonOperator
from airflow.utils.dates import days_ago
from slack_connection import SlackNotifier


def retry_yt_scrape_events():
    import os
    import json
    import requests
    import clickhouse_connect

    os.environ["no_proxy"] = "*"
    max_handles = 1000

    connection = BaseHook.get_connection("clickhouse_gcc")
    extra_params = json.loads(connection.extra)
    client = clickhouse_connect.get_client(host=connection.host, password=connection.password, username=connection.login,
                                           send_receive_timeout=extra_params['send_receive_timeout'])
    print(client)
    query = '''
                select  scl_id scrape_id, 
                        groupArray(status) status_list,
                        argMax(flow, event_timestamp) flow,
                        argMax(params, event_timestamp) params,
                        argMax(retry_count, event_timestamp) retry_count,
                        argMax(platform, event_timestamp) platform,
                        has(status_list, 'COMPLETE') is_complete,
                        has(status_list, 'PROCESSING') is_processing,
                        has(status_list, 'FAILED') is_fail,
                        (is_processing and not is_complete and not is_fail) flag,
                        argMax(picked_at, event_timestamp) processing_start,
                        round((now() - processing_start)/60, 2) processing_since
                from _e.scrape_request_log_events
                group by scl_id
                having flag=1
            '''

    result = client.query(query + f" LIMIT {max_handles}")
    
    beat_connection = BaseHook.get_connection("beat_stage")
    base_url = beat_connection.host
    
    for row in result.result_rows:
        scrape_id = row[0]
        flow = row[2]
        params = json.loads(row[3])
        retry_count = row[4]
        platform = row[5]
        print(params)
        refresh_post_scl = {"flow": flow, "platform": platform, "params": params, "retry_count": retry_count+1}
        if retry_count == 3 or flow.endswith("csv")==False:
            refresh_post_scl["status"]="FAILED"
        else:
            refresh_post_scl["status"]="PENDING"
        print(refresh_post_scl)
        refresh_post_url = base_url + "/scrape_request_log/flow/update/%s" % scrape_id
        requests.post(url=refresh_post_url, data=json.dumps(refresh_post_scl))


with DAG(
        dag_id='retry_yt_scrape_events',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval="0 */3 * * *",
) as dag:
    task = PythonOperator(
        task_id='retry_yt_scrape_events',
        python_callable=retry_yt_scrape_events
    )
    task
