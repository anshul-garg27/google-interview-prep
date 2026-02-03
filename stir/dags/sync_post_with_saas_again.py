import pika
import json
import clickhouse_connect
from airflow import DAG
from airflow.hooks.base import BaseHook
from airflow.operators.python_operator import PythonOperator
from airflow.utils.dates import days_ago
from slack_connection import SlackNotifier

def sync_post_with_saas_again():
    connection = BaseHook.get_connection("clickhouse_gcc")
    extra_params = json.loads(connection.extra)
    client = clickhouse_connect.get_client(host=connection.host, password=connection.password, username=connection.login, send_receive_timeout=extra_params['send_receive_timeout'])

    submissionIds = client.query("""select ds.id
                            from dbt.stg_campaign_deliverable_submission ds
                            where ds.saas_post_item_status is not null
                            and ((ds.status = 'UNDER_REVIEW' AND ds.saas_post_item_status != 'UNDER_REVIEW_ADDED') OR
                            (ds.status = 'APPROVED' AND ds.saas_post_item_status != 'ADDED') OR
                            (ds.status = 'REJECTED' AND ds.saas_post_item_status != 'REMOVED'))
                            """)

    rabbitmq_conn = BaseHook.get_connection("rabbitmq_connection")
    rabbitmq_host = rabbitmq_conn.host
    rabbitmq_username = 'admin'
    rabbitmq_password = 'admin'
    rabbitmq_port = rabbitmq_conn.port
    connection_params = pika.ConnectionParameters(
        host=rabbitmq_host,
        credentials=pika.PlainCredentials(rabbitmq_username, rabbitmq_password),
        port=rabbitmq_port
    )
    connection = pika.BlockingConnection(connection_params)
    channel = connection.channel()

    for row in submissionIds.result_rows:
        json_message = json.dumps(row[0])
        channel.basic_publish(exchange='campaign.dx',
                              routing_key='campaign.SYNC_POST_DATA_WITH_SAAS_RK', body=json_message)
    connection.close()


with DAG(
        dag_id='sync_post_with_saas_again',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval="0 * * * *",
) as dag:
    task = PythonOperator(
        task_id='sync_post_with_saas_again',
        python_callable=sync_post_with_saas_again
    )
    task
