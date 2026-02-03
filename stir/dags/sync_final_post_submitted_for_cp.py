import pika
import json
import clickhouse_connect
from airflow import DAG
from airflow.hooks.base import BaseHook
from airflow.operators.python_operator import PythonOperator
from airflow.utils.dates import days_ago
from slack_connection import SlackNotifier

#creating payout for active sus where not exist already and su created_on> 30 aug 2023



def sync_final_post_submitted_for_cp():
    connection = BaseHook.get_connection("clickhouse_gcc")
    extra_params = json.loads(connection.extra)
    client = clickhouse_connect.get_client(host=connection.host, password=connection.password, username=connection.login, send_receive_timeout=extra_params['send_receive_timeout'])

    finalPostData = client.query("""select sub.id, sub.campaign_profile_id
                                from dbt.stg_identity_social_account sa,
                                dbt.stg_campaign_deliverable_submission sub
                                where sub.level = 'FINAL_POST'
                                and sub.platform_id in (1, 2)
                                and sa.campaign_profile_id = sub.campaign_profile_id
                                and (sa.final_post_link_submitted_once is null or sa.final_post_link_submitted_once = 0)
                                order by sub.id desc limit 1000
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

    for row in finalPostData.result_rows:
        json_message = json.dumps(row[1])
        channel.basic_publish(exchange='campaign.dx',
                              routing_key='campaign.SYNC_POST_CREATION_FOR_CP_WITH_IDENTITY_RK', body=json_message)
    connection.close()


with DAG(
        dag_id='sync_final_post_submitted_for_cp',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval="0 * * * *",
) as dag:
    task = PythonOperator(
        task_id='sync_final_post_submitted_for_cp',
        python_callable=sync_final_post_submitted_for_cp
    )
    task
