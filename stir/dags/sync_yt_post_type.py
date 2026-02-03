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
    post_ids = []
    max_post_ids = 1000

    connection = BaseHook.get_connection("clickhouse_gcc")
    extra_params = json.loads(connection.extra)
    client = clickhouse_connect.get_client(host=connection.host, password=connection.password, username=connection.login, send_receive_timeout=extra_params['send_receive_timeout'])

    yt_post_ids = client.query('''with
                                candidate_set as (
                                    select 
                                        shortcode AS post_id
                                    from 
                                        dbt.stg_beat_youtube_post
                                    where 
                                        post_type is NULL or post_type == ''
                                ),
                                failures as (
                                        select 
                                            JSONExtractString(params, 'post_id') as post_id
                                        from dbt.stg_beat_scrape_request_log
                                        where 
                                            flow = 'refresh_yt_post_type'
                                            and status == 'FAILED'
                                        group by post_id
                                        having 
                                            uniqExactIf(id, status = 'FAILED') > 3
                                )
                                select 
                                    post_id 
                                from 
                                    candidate_set 
                                where 
                                    post_id not in failures
                                '''
                            f"LIMIT {max_post_ids}")

    for post_id in yt_post_ids.result_rows:
        post_ids.append(post_id[0])

    flow = 'refresh_yt_post_type'
    beat_connection = BaseHook.get_connection("beat")
    base_url = beat_connection.host
    url = base_url + "/scrape_request_log/flow/%s" % flow

    for post_id in post_ids:
        data = {"flow": flow, "platform": "YOUTUBE", "params": {"post_id": post_id}}
        requests.post(url=url, data=json.dumps(data))

with DAG(
        dag_id='sync_yt_post_type',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval="*/10 * * * *",
) as dag:
    task = PythonOperator(
        task_id='update_yt_post_type',
        python_callable=create_scrape_request_log
    )
    task
