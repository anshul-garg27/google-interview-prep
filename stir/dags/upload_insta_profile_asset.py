import struct

from airflow import DAG
from airflow.hooks.base import BaseHook
from airflow.operators.python import PythonOperator
from airflow.utils.dates import days_ago
import hashlib
import logging

from slack_connection import SlackNotifier



def create_scrape_request_log():
    import os
    import json
    import requests
    import clickhouse_connect

    os.environ["no_proxy"] = "*"
    max_handles = 10000

    connection = BaseHook.get_connection("clickhouse_gcc")
    extra_params = json.loads(connection.extra)
    client = clickhouse_connect.get_client(host=connection.host, password=connection.password, username=connection.login, send_receive_timeout=extra_params['send_receive_timeout'])

    result = client.query('''with 
                                        recent_requests as (
                                            select 
                                                JSONExtractUInt(params, 'hash') hash
                                            FROM _e.scrape_request_log_events as scrape
                                            where flow = 'asset_upload_flow' and event_timestamp >= now() - INTERVAL 24 HOUR
                                        ),
                                        profiles as (
                                            select 
                                                profile_id,
                                                argMax(JSONExtractString(dimensions, 'profile_pic_url'), event_timestamp) profile_pic_url 
                                            FROM _e.profile_log_events
                                            where
                                                date(event_timestamp) >= date(now()) - INTERVAL 1 DAY
                                                and event_timestamp >= now() - INTERVAL 12 HOUR
                                                and platform = 'INSTAGRAM'
                                                and JSONExtractString(dimensions, 'profile_pic_url') is not null and JSONExtractString(dimensions, 'profile_pic_url') != ''
                                            group by profile_id
                                        ),
                                        profile_ids as (
                                            select profile_id
                                            FROM profiles
                                        )
                                        SELECT profile_id, profile_pic_url, halfMD5(profile_pic_url) hash
                                        FROM profiles
                                        WHERE halfMD5(profiles.profile_pic_url) not in recent_requests'''
                                    f" LIMIT {max_handles}")

    flow = 'asset_upload_flow'
    entity_type = "PROFILE"
    asset_type = "IMAGE"
    platform = "INSTAGRAM"

    beat_connection = BaseHook.get_connection("beat")
    base_url = beat_connection.host
    url = base_url + "/scrape_request_log/flow/%s" % flow

    for profile in result.result_rows:
        profile_id = profile[0]
        profile_pic_url = profile[1]
        hash = profile[2]
        try:
            data = {"flow": flow, "platform": platform,
                    "params": {"entity_type": entity_type, "asset_type": asset_type, "platform": platform,
                               "entity_id": profile_id, "original_url": profile_pic_url,
                               "hash": hash}}
            requests.post(url=url, data=json.dumps(data))
        except Exception as e:
            logging.error("Error occurred : {e}".format(e=str(e)))


with DAG(
        dag_id='sync_insta_profile_asset',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval="*/10 * * * *",
) as dag:
    task = PythonOperator(
        task_id='create_asset_insta_profile_scrape_request_log',
        python_callable=create_scrape_request_log
    )
    task
