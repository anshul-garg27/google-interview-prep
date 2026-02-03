import datetime
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

    connection = BaseHook.get_connection("clickhouse_gcc")
    extra_params = json.loads(connection.extra)
    client = clickhouse_connect.get_client(host=connection.host, password=connection.password, username=connection.login, send_receive_timeout=extra_params['send_receive_timeout'])
    stores = client.query('''select
                                handle
                            from 
                                dbt.stg_beat_credential as credential
                            where  
                                credential.source = 'shopify'
                                and credential.enabled = TRUE
                            ''')

    flow = 'refresh_orders_by_store'
    beat_connection = BaseHook.get_connection("beat")
    base_url = beat_connection.host
    url = base_url + "/scrape_request_log/flow/%s" % flow

    current_time = datetime.datetime.now()
    last_twelve_hours = current_time - datetime.timedelta(hours=12)
    last_twelve_hours_format = last_twelve_hours.strftime('%Y-%m-%dT%H:%M:%S%z')

    total_limit = 1000
    limit_per_batch = 250
    num_batches = total_limit // limit_per_batch

    for store in stores.result_rows:
        for i in range(num_batches):
            data = {"flow": flow, "platform": "SHOPIFY", "params": {"limit": limit_per_batch, "store": store[0], "updated_at_min": last_twelve_hours_format}}
            requests.post(url=url, data=json.dumps(data))

with DAG(
        dag_id='sync_shopify_orders',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval='0 */4 * * *',
) as dag:
    task = PythonOperator(
        task_id='create_shopify_orders_scrape_request_log',
        python_callable=create_scrape_request_log
    )
    task
