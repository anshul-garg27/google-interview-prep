from airflow import DAG
from airflow.hooks.base import BaseHook
from airflow.operators.python import PythonOperator
from airflow.utils.dates import days_ago
from airflow.models import Variable
# from slack_connection import SlackNotifier
import pika

from slack_connection import SlackNotifier

def create_post_colleciton_sentiment_report():
    import os
    import json
    import requests
    import clickhouse_connect

    os.environ["no_proxy"] = "*"
    max_collections = 250

    connection = BaseHook.get_connection("clickhouse_gcc")
    extra_params = json.loads(connection.extra)
    client = clickhouse_connect.get_client(host=connection.host, password=connection.password, username=connection.login, send_receive_timeout=extra_params['send_receive_timeout'])
    post_collection_ids = client.query('''                                    
                                        with
                                            collection_ids as (
                                                select scpc.id id, count(*) count
                                                from dbt.stg_coffee_post_collection scpc
                                                inner join dbt.stg_coffee_post_collection_item scpci on scpci.post_collection_id = scpc.id
                                                inner join _e.post_activity_log_events pale on pale.shortcode = scpci.short_code
                                                where scpc.partner_id = 11743 and pale.publish_time>= now() - INTERVAL 2 WEEK
                                                group by scpc.id
                                                having count>0
                                            )
                                        select id from collection_ids
                                        '''
                                 f"LIMIT {max_collections}")

    rabbitmq_conn = BaseHook.get_connection("rabbitmq_connection")
    rabbitmq_host = rabbitmq_conn.host
    rabbitmq_username = 'admin'
    rabbitmq_password = 'admin'
    rabbitmq_port = rabbitmq_conn.port
    rabbitmq_queue = "sentiment_collection_report_in_q"
    connection_params = pika.ConnectionParameters(
        host=rabbitmq_host,
        credentials=pika.PlainCredentials(rabbitmq_username, rabbitmq_password),
        port=rabbitmq_port
    )

    connection = pika.BlockingConnection(connection_params)
    channel = connection.channel()
    channel.queue_declare(queue=rabbitmq_queue, durable=True)
    for row in post_collection_ids.result_rows:
        id = row[0]
        
        message = {
            "collectionId": id,
            "collectionType": "POST"
        }
        
        json_message = json.dumps(message)
        channel.basic_publish(exchange='', routing_key=rabbitmq_queue, body=json_message)
    connection.close()

with DAG(
        dag_id='sync_post_collection_sentiment_report',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval="0 */3 * * *",
) as dag:
    task = PythonOperator(
        task_id='create_post_collection_sentiment_report',
        python_callable=create_post_colleciton_sentiment_report
    )
    task
