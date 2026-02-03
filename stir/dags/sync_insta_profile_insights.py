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
    handles = []
    max_handles = 1000

    connection = BaseHook.get_connection("clickhouse_gcc")
    extra_params = json.loads(connection.extra)
    client = clickhouse_connect.get_client(host=connection.host, password=connection.password, username=connection.login, send_receive_timeout=extra_params['send_receive_timeout'])
    insta_handles = client.query('''with
                                            candidate_set as (
                                                SELECT
                                                    cred.handle as handle
                                                FROM dbt.stg_beat_credential as cred
                                                LEFT JOIN dbt.stg_beat_instagram_profile_insights AS profile_insights ON cred.handle = profile_insights.handle 
                                                AND profile_insights.updated_at<= now() - INTERVAL 1 WEEK 
                                                WHERE cred.enabled = TRUE 
                                                AND cred.source = 'graphapi'
                                                group by handle
                                            ),
                                            attempts as (
                                                SELECT
                                                    JSONExtractString(params, 'handle') as handle from dbt.stg_beat_scrape_request_log as scrape
                                            where scrape.flow = 'refresh_profile_insights' and scrape.created_at >= now() - INTERVAL 1 DAY
                                            group by handle
                                            )
                                        SELECT candidate_set.handle
                                        FROM candidate_set
                                        where candidate_set.handle not in attempts
                                '''
                                f"LIMIT {max_handles}")

    for handle in insta_handles.result_rows:
        handles.append(handle[0])
    handles = list(set(handles))
    flow = 'refresh_profile_insights'
    beat_connection = BaseHook.get_connection("beat")
    base_url = beat_connection.host
    url = base_url + "/scrape_request_log/flow/%s" % flow
    for handle in handles:
        data = {"flow": flow, "platform": "INSTAGRAM", "params": {"handle": handle}}
        requests.post(url=url, data=json.dumps(data))


insta_profile_insights_schedule_interval = Variable.get('insta_profile_insights_schedule_interval')

with DAG(
        dag_id='sync_insta_profile_insights',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval=insta_profile_insights_schedule_interval,
) as dag:
    task = PythonOperator(
        task_id='create_insta_insights_scrape_request_log',
        python_callable=create_scrape_request_log
    )
    task

