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
    insta_posts = client.query('''with  
                                        recent_requests as (
                                            SELECT
                                                JSONExtractString(params, 'shortcode') as shortcode 
                                            from _e.scrape_request_log_events as scrape
                                            where scrape.flow = 'refresh_post_insights' and scrape.event_timestamp >= now() - INTERVAL 24 HOUR
                                        ),
                                        handles as (
                                            SELECT cred.handle FROM dbt.stg_beat_credential as cred
                                            where cred.source = 'graphapi' and cred.enabled = TRUE
                                        ),
                                        candidate_set as (
                                            SELECT
                                                handle,
                                                post.post_id,
                                                post.shortcode,
                                                post.post_type
                                            FROM 
                                            dbt.stg_beat_instagram_post AS post 
                                            WHERE post.handle IN handles
                                            AND post.post_id is not NULL
                                            AND post.post_type in ('reels', 'image', 'carousel') 
                                            AND post.publish_time >= now() - INTERVAL 1 MONTH
                                            group by
                                                handle,
                                                post.post_id,
                                                post.shortcode,
                                                post.post_type,
                                                post.publish_time
                                            order by post.publish_time ASC
                                        )
                                    SELECT candidate_set.handle, candidate_set.post_id, candidate_set.shortcode, candidate_set.post_type
                                    FROM candidate_set
                                    where candidate_set.shortcode not in recent_requests
                                '''
                               f"LIMIT {max_handles}")

    posts = [sublist for sublist in insta_posts.result_rows if all(sublist)]

    flow = 'refresh_post_insights'
    beat_connection = BaseHook.get_connection("beat")
    base_url = beat_connection.host
    url = base_url + "/scrape_request_log/flow/%s" % flow

    for (handle, post_id, shortcode, post_type) in posts:
        data = {"flow": flow, "platform": "INSTAGRAM",
                "params": {"handle": handle, "post_id": post_id, "shortcode": shortcode, "post_type": post_type}}
        requests.post(url=url, data=json.dumps(data))

posts_insights_schedule_interval = Variable.get('posts_insights_schedule_interval')


with DAG(
        dag_id='sync_insta_post_insights',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval=posts_insights_schedule_interval
) as dag:
    task = PythonOperator(
        task_id='create_insta_post_insights_scrape_request_log',
        python_callable=create_scrape_request_log
    )
    task
