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
    max_profiles = 25000

    connection = BaseHook.get_connection("clickhouse_gcc")
    extra_params = json.loads(connection.extra)
    client = clickhouse_connect.get_client(host=connection.host, password='CuB31Z@3995ddOpT', username="cube", send_receive_timeout=extra_params['send_receive_timeout'])

    profiles = client.query(''' with abc as (
                                            SELECT profile_id, count(*) AS `count`
                                            FROM _e.profile_log_events
                                            WHERE (source = 'openapi_profile_info_v0.10'
                                                AND JSONHas(dimensions, 'gender')
                                                )
                                            GROUP BY profile_id
                                            ORDER BY count DESC
                                )
                                select name, handle, bio
                                from dbt.mart_instagram_account
                                where
                                    handle not in (select abc.profile_id from abc) and
                                    country = 'IN'
                                and followers >= 100000
                                '''
                            f"LIMIT {max_profiles}")

    flow = 'refresh_instagram_gpt_data_base_gender'
    beat_connection = BaseHook.get_connection("beat")
    base_url = beat_connection.host
    url = base_url + "/scrape_request_log/flow/%s" % flow
    for profile in profiles.result_rows:
        data = {"flow": flow, "platform": "INSTAGRAM", "params": {"data": {"handle": profile[1], "name": profile[0], "bio": profile[2]}}}
        requests.post(url=url, data=json.dumps(data))


with DAG(
        dag_id='sync_instagram_gpt_data_base_gender',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval="0 */8 * * *",
) as dag:
    task = PythonOperator(
        task_id='create_instagram_gpt_data_base_gender_scl',
        python_callable=create_scrape_request_log
    )
    task
