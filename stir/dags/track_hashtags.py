from airflow import DAG
from airflow.hooks.base import BaseHook
from airflow.operators.python import PythonOperator
from airflow.utils.dates import days_ago
from airflow.models import Variable
from slack_connection import SlackNotifier


def create_insta_hashtag_scl():
    import os
    import json
    import requests

    os.environ["no_proxy"] = "*"
    hashtags = ["IndulekhaSvetakutaja", "TreatFromTheRoots", "ScalpGoals", "NoMoreDandruff", "HairCare", "IndulekhaPartner", "indulekha_care"]
    flow = 'fetch_hashtag_posts'
    batch_size = 2000
    beat_connection = BaseHook.get_connection("beat")
    base_url = beat_connection.host
    url = base_url + "/scrape_request_log/flow/%s" % flow
    for hashtag in hashtags:
        data = {"flow": flow, "platform": "INSTAGRAM", "params": {"hashtag": hashtag, "max_posts": batch_size}}
        requests.post(url=url, data=json.dumps(data))


with DAG(
        dag_id='track_hashtags',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval="0 */6 * * *",
) as dag:
    task = PythonOperator(
        task_id='create_insta_hashtag_scl',
        python_callable=create_insta_hashtag_scl
    )
    task
