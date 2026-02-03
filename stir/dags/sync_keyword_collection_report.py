from airflow import DAG
from airflow.hooks.base import BaseHook
from airflow.operators.python import PythonOperator
from airflow.utils.dates import days_ago
from airflow.models import Variable
from slack_connection import SlackNotifier
import pika

def create_keyword_collection_report():
    import os
    import json
    import requests
    import clickhouse_connect

    os.environ["no_proxy"] = "*"
    max_handles = 250

    connection = BaseHook.get_connection("clickhouse_gcc")
    extra_params = json.loads(connection.extra)
    client = clickhouse_connect.get_client(host=connection.host, password=connection.password, username=connection.login, send_receive_timeout=extra_params['send_receive_timeout'])
    post_collection_ids = client.query('''
                                            select end_time, keywords, platform, start_time, id
                                            from dbt.stg_coffee_keyword_collection scpc
                                            where scpc.created_at <= now() - INTERVAL 2 HOUR and report_path is null and enabled=true
                                        '''
                                 f"LIMIT {max_handles}")

    rabbitmq_conn = BaseHook.get_connection("rabbitmq_connection")
    rabbitmq_host = rabbitmq_conn.host
    rabbitmq_username = 'admin'
    rabbitmq_password = 'admin'
    rabbitmq_port = rabbitmq_conn.port
    rabbitmq_queue = "keyword_collection_q"
    connection_params = pika.ConnectionParameters(
        host=rabbitmq_host,
        credentials=pika.PlainCredentials(rabbitmq_username, rabbitmq_password),
        port=rabbitmq_port
    )
    connection = pika.BlockingConnection(connection_params)
    channel = connection.channel()
    channel.queue_declare(queue=rabbitmq_queue, durable=True)
    for row in post_collection_ids.result_rows:
        end_date_only = row[0].date()
        end_date = end_date_only.strftime('%Y-%m-%d')
        keywords = row[1]
        platform = row[2]
        start_date_only = row[3].date()
        start_date = start_date_only.strftime('%Y-%m-%d')
        job_id = f"PROD_KC_{row[4]}"

        message = {
            "end_date": end_date,
            "job_id": job_id,
            "keywords": keywords,
            "platform": platform,
            "start_date": start_date
        }

        json_message = json.dumps(message)
        channel.basic_publish(exchange='', routing_key=rabbitmq_queue, body=json_message)
    connection.close()

with DAG(
        dag_id='sync_keyword_collection_report',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval="0 */3 * * *",
) as dag:
    task = PythonOperator(
        task_id='create_keyword_collection_report',
        python_callable=create_keyword_collection_report
    )
    task
