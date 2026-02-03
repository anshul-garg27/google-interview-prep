import hashlib
import logging
import struct

from airflow import DAG
from airflow.hooks.base import BaseHook
from airflow.operators.python import PythonOperator
from airflow.utils.dates import days_ago
from slack_connection import SlackNotifier


def half_md5(string):
    """
        # Convert string to bytes
        # Compute the MD5 hash
        # Convert the hash to an integer
        # Take the upper 64 bits of the integer
        # Convert the upper bits to a UInt64 in little-endian byte order
    """
    bytes_string = string.encode('utf-8')
    md5_hash = hashlib.md5(bytes_string).hexdigest()
    hash_int = int(md5_hash, 16)
    upper_bits = hash_int >> 64
    uint64_bytes = struct.pack('<Q', upper_bits)
    return uint64_bytes


def create_scrape_request_log():
    import os
    import json
    import requests
    import clickhouse_connect

    os.environ["no_proxy"] = "*"
    max_posts = 25000

    connection = BaseHook.get_connection("clickhouse_gcc")
    extra_params = json.loads(connection.extra)
    client = clickhouse_connect.get_client(host=connection.host, password=connection.password, username=connection.login, send_receive_timeout=extra_params['send_receive_timeout'])

    # Instagram Posts
    result = client.query('''
                            with 
                                recent_requests as (
                                    select 
                                        JSONExtractString(params, 'entity_id') entity_id
                                    FROM _e.scrape_request_log_events
                                    where flow = 'asset_upload_flow' and event_timestamp >= now() - INTERVAL 24 HOUR
                                ),
                                posts as (
                                    select 
                                        shortcode, 
                                        argMax(JSONExtractString(dimensions, 'thumbnail_src'), event_timestamp) thumbnail_url 
                                    from _e.post_log_events
                                    where 
                                    date(event_timestamp) >= date(now()) - INTERVAL 1 DAY 
                                    and JSONExtractString(dimensions, 'post_type') == 'story'
                                    and event_timestamp >= now() - INTERVAL 12 HOUR
                                    and platform = 'INSTAGRAM'
                                    and JSONExtractString(dimensions, 'thumbnail_src') is not null and JSONExtractString(dimensions, 'thumbnail_src') != ''
                                    group by shortcode
                                ),
                                shortcodes as (
                                    select shortcode
                                    FROM posts
                                ),
                                with_assets as (
                                    select entity_id shortcode
                                    from dbt.stg_beat_asset_log
                                    where entity_id in shortcodes
                                )
                                SELECT shortcode, thumbnail_url
                                FROM posts
                                WHERE shortcode not in with_assets and shortcode not in recent_requests

                                '''
                            f" LIMIT {max_posts}")

    flow = 'asset_upload_flow_stories'
    entity_type = "POST"
    platform = "INSTAGRAM"
    beat_connection = BaseHook.get_connection("beat")
    base_url = beat_connection.host
    url = base_url + "/scrape_request_log/flow/%s" % flow

    for post in result.result_rows:
        shortcode = post[0]
        thumbnail_url = post[1]
        try:
            data = {"flow": flow, "platform": platform,
                    "params": {"entity_type": entity_type, "asset_type": "IMAGE", "platform": platform,
                               "entity_id": shortcode, "original_url": thumbnail_url,
                               "hash": int.from_bytes(half_md5(thumbnail_url), byteorder='little')}}
            requests.post(url=url, data=json.dumps(data))
        except Exception as e:
            logging.error("Error occurred : {e}".format(e=str(e)))


with DAG(
        dag_id='sync_post_stories_asset_log',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval="*/10 * * * *",
) as dag:
    task = PythonOperator(
        task_id='create_post_stories_asset_scrape_request_log',
        python_callable=create_scrape_request_log
    )
    task
