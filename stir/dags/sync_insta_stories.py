from airflow import DAG
from airflow.hooks.base import BaseHook
from airflow.operators.python import PythonOperator
from airflow.utils.dates import days_ago
from airflow.models import Variable
from slack_connection import SlackNotifier


def create_scrape_request_log():
    import os
    import json
    import requests
    import clickhouse_connect

    os.environ["no_proxy"] = "*"
    max_handles = 1000

    connection = BaseHook.get_connection("clickhouse_gcc")
    extra_params = json.loads(connection.extra)
    client = clickhouse_connect.get_client(host=connection.host, password=connection.password, username=connection.login, send_receive_timeout=extra_params['send_receive_timeout'])
    insta_handles = client.query('''with
                                        cred as (
                                            select
                                                    handle
                                                from dbt.stg_beat_credential as credential
                                                where credential.source = 'graphapi'
                                                and credential.enabled = TRUE
                                        ),
                                        attempted as (
                                            select
                                                JSONExtractString(params, 'handle') as handle
                                            from _e.scrape_request_log_events as scrape
                                            where scrape.flow = 'refresh_stories_posts' 
                                            and scrape.event_timestamp >= now() - INTERVAL 1 DAY
                                            group by handle
                                        )
                                        select cred.handle
                                        from cred
                                        where handle not in attempted
                                '''
                                 f"LIMIT {max_handles}")

    handles = [sublist for sublist in insta_handles.result_rows if all(sublist)]

    flow = 'refresh_stories_posts'
    beat_connection = BaseHook.get_connection("beat")
    base_url = beat_connection.host
    url = base_url + "/scrape_request_log/flow/%s" % flow
    for handle in handles:
        data = {"flow": flow, "platform": "INSTAGRAM", "params": {"handle": handle[0]}}
        requests.post(url=url, data=json.dumps(data))


posts_stories_schedule_interval = Variable.get('posts_stories_schedule_interval')

with DAG(
        dag_id='sync_insta_posts_stories',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval=posts_stories_schedule_interval,
) as dag:
    task = PythonOperator(
        task_id='create_insta_posts_stories_scrape_request_log',
        python_callable=create_scrape_request_log
    )
    task
