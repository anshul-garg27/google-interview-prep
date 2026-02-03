import pika
import json
import clickhouse_connect
from airflow import DAG
from airflow.hooks.base import BaseHook
from airflow.operators.python_operator import PythonOperator
from airflow.utils.dates import days_ago
from slack_connection import SlackNotifier

def create_payout_for_active_su_if_not_present():
    connection = BaseHook.get_connection("clickhouse_gcc")
    extra_params = json.loads(connection.extra)
    client = clickhouse_connect.get_client(host=connection.host, password=connection.password, username=connection.login, send_receive_timeout=extra_params['send_receive_timeout'])

    suIds = client.query("""select su.id
                            from dbt.stg_campaign_shortlisted_user su
                            where su.active = 1
                            and su.accepted = 1
                            and su.id not in (select payout.shortlisted_user_id from dbt.stg_campaign_campaign_payout payout)
                            and su.created_on > '2023-08-30 00:00:00'
                            order by su.id desc
                            limit 1000
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

    for row in suIds.result_rows:
        json_message = json.dumps(row[0])
        channel.basic_publish(exchange='campaign.dx',
                              routing_key='campaign.CREATE_PAYOUT_FOR_SU_IF_ALREADY_NOT_EXIST_RK', body=json_message)
    connection.close()


with DAG(
        dag_id='create_payout_for_active_su_if_not_present',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval="0 * * * *",
) as dag:
    task = PythonOperator(
        task_id='create_payout_for_active_su_if_not_present',
        python_callable=create_payout_for_active_su_if_not_present
    )
    task
