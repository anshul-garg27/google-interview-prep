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
    insta_handles = client.query('''select profile_id, handle from dbt.mart_explicitly_instagram_stories
                                '''
                                 f"LIMIT {max_handles}")

    stories = [sublist for sublist in insta_handles.result_rows if all(sublist)]

    flow = 'refresh_story_posts_by_profile_id'
    beat_connection = BaseHook.get_connection("beat")
    base_url = beat_connection.host
    url = base_url + "/scrape_request_log/flow/%s" % flow
    for (profile_id, handle) in stories:
        data = {"flow": flow, "platform": "INSTAGRAM", "params": {"profile_id": profile_id, "handle": handle}}
        requests.post(url=url, data=json.dumps(data))


posts_stories_schedule_interval = Variable.get('posts_stories_schedule_interval')

with DAG(
        dag_id='sync_insta_stories_explititly',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval="0 */5 * * *",
) as dag:
    task = PythonOperator(
        task_id='create_insta_stories_explicitly_scrape_request_log',
        python_callable=create_scrape_request_log
    )
    task
