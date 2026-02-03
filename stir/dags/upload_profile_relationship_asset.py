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
    max_handles = 1000

    connection = BaseHook.get_connection("clickhouse_gcc")
    extra_params = json.loads(connection.extra)
    client = clickhouse_connect.get_client(host=connection.host, password=connection.password, username=connection.login, send_receive_timeout=extra_params['send_receive_timeout'])

    profile_relationship = client.query('''with
                                            profile_relationship as
                                            (SELECT
                                                JSONExtractString(source_dimensions, 'profile_pic_url') source_pic,
                                                JSONExtractString(target_dimensions, 'profile_pic_url') target_pic,
                                                coalesce(
                                                        if(target_pic = '', NULL, target_pic),
                                                        if(source_pic = '', NULL, source_pic)
                                                    ) profile_pic_url,
                                                coalesce(
                                                        if(target_pic = '', NULL, target_profile_id),
                                                        if(source_pic = '', NULL, source_profile_id)
                                                    ) profile_id,
                                                     timestamp,
                                                     platform
                                                FROM beat_replica.profile_relationship_log
                                        )
                                        select pr.profile_id, pr.profile_pic_url, pr.platform
                                        from profile_relationship as pr
                                        LEFT JOIN dbt.stg_beat_asset_log AS asset ON pr.profile_id = asset.entity_id
                                        LEFT JOIN _e.scrape_request_log_events AS scrape ON JSONExtractUInt(scrape.params, 'hash') = halfMD5(profile_pic_url) AND scrape.flow = 'asset_upload_flow'
                                        WHERE asset.entity_id IS NULL and scrape.id <= 0 and pr.timestamp >= now() - INTERVAL 1 DAY'''
                                        f" LIMIT {max_handles}")

    profile_relationship = [sublist for sublist in profile_relationship.result_rows if all(sublist)]
    flow = 'asset_upload_flow'
    entity_type = "PROFILE"
    asset_type = "IMAGE"
    platform = "INSTAGRAM"
    beat_connection = BaseHook.get_connection("beat")
    base_url = beat_connection.host
    url = base_url + "/scrape_request_log/flow/%s" % flow

    for (entity_id, profile_pic_url, platform) in profile_relationship:
        try:
            data = {"flow": flow, "platform": platform,
                    "params": {"entity_type": entity_type, "asset_type": asset_type, "platform": platform,
                               "entity_id": entity_id, "original_url": profile_pic_url,
                               "hash": int.from_bytes(half_md5(profile_pic_url), byteorder='little')}}
            requests.post(url=url, data=json.dumps(data))
        except Exception as e:
            logging.error("Error occurred : {e}".format(e=str(e)))


with DAG(
        dag_id='sync_profile_relationship_asset',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval="*/10 * * * *",
) as dag:
    task = PythonOperator(
        task_id='create_asset_profile_relationship_scrape_request_log_relationship',
        python_callable=create_scrape_request_log
    )
    task
