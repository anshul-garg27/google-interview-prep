import pika
import json
import clickhouse_connect
from airflow import DAG
from airflow.hooks.base import BaseHook
from airflow.operators.python_operator import PythonOperator
from airflow.utils.dates import days_ago
from slack_connection import SlackNotifier

#syncing non synced uca as saas item with saas for active uca where campaign already synced with saas and uca created on >30 aug 2023


def uca_su_saas_item_sync():
    connection = BaseHook.get_connection("clickhouse_gcc")
    extra_params = json.loads(connection.extra)
    client = clickhouse_connect.get_client(host=connection.host, password=connection.password, username=connection.login, send_receive_timeout=extra_params['send_receive_timeout'])

    ucaIds = client.query("""select uca.id
                                from dbt.stg_campaign_user_campaign_application uca
                                where uca.active = 1
                                and uca.campaign_id in (select campaign.id from dbt.stg_campaign_campaign campaign where saas_campaign_collection_id is not null)
                                and uca.saas_campaign_collection_item_id is  null
                                and uca.created_on > '2023-08-30 00:00:00'
                                order by uca.id asc limit 5000
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

    for row in ucaIds.result_rows:
        json_message = json.dumps(row[0])
        channel.basic_publish(exchange='campaign.dx',
                              routing_key='campaign.SYNC_UCA_WITH_SAAS_AS_PROFILE_ITEM_RK', body=json_message)
    connection.close()


with DAG(
        dag_id='uca_su_saas_item_sync',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval="0 * * * *",
) as dag:
    task = PythonOperator(
        task_id='uca_su_saas_item_sync',
        python_callable=uca_su_saas_item_sync
    )
    task


