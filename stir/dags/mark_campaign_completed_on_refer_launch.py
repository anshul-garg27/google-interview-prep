import pika
import json
import clickhouse_connect
from airflow import DAG
from airflow.hooks.base import BaseHook
from airflow.operators.python_operator import PythonOperator
from airflow.utils.dates import days_ago
from slack_connection import SlackNotifier

def mark_campaign_completed_on_refer_launch():
    connection = BaseHook.get_connection("clickhouse_gcc")
    extra_params = json.loads(connection.extra)
    client = clickhouse_connect.get_client(host=connection.host, password=connection.password, username=connection.login, send_receive_timeout=extra_params['send_receive_timeout'])

    posts = client.query('''
        with dynamic_link as (
            select 
                distinct splitByChar('?', COALESCE(dynamic_link, ''))[1] as launch_link 
            from _e.launch_refer_events 
            where event_timestamp >= toDateTime(now()- INTERVAL 1 DAY )
        ),
        final as (
            select 
                distinct shortlisted_user_id 
            from dbt.stg_campaign_deliverable_shortlisted_user_link 
            where JSONExtractString(JSONExtractArrayRaw(COALESCE(unique_link_data, ''))[1], 'link') in (select launch_link from dynamic_link)
        )
        select * from final
        ''')

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

    for row in posts.result_rows:
        shortlisted_user_id = row[0]
        json_message = json.dumps(shortlisted_user_id)
        channel.basic_publish(exchange='campaign.dx', routing_key='campaign.HANDLE_SUBMISSION_FOR_REFERRAL_CAMPAIGN_SU_RK', body=json_message)
    connection.close()


with DAG(
        dag_id='mark_campaign_completed_on_refer_launch',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval="0 * * * *",
) as dag:
    task = PythonOperator(
        task_id='mark_campaign_completed_on_refer_launch',
        python_callable=mark_campaign_completed_on_refer_launch
    )
    task
