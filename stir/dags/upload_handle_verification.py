import pika
import json
import clickhouse_connect
from airflow import DAG
from airflow.hooks.base import BaseHook
from airflow.operators.python_operator import PythonOperator
from airflow.utils.dates import days_ago
from slack_connection import SlackNotifier

def handle_verification_for_posts():
    connection = BaseHook.get_connection("clickhouse_gcc")
    extra_params = json.loads(connection.extra)
    client = clickhouse_connect.get_client(host=connection.host, password=connection.password, username=connection.login, send_receive_timeout=extra_params['send_receive_timeout'])

    posts = client.query('''
        with submission_data  as (
            select sub.id as submissionId,
                CAST(sub.campaign_profile_id AS Int64) as cpId,
                sub.short_code as shortCode
            from dbt.stg_campaign_deliverable_submission sub
                inner join dbt.stg_campaign_campaign campaign on campaign.id = sub.campaign_id
                where sub.status = 'UNDER_REVIEW'
                   and sub.level = 'FINAL_POST'
                   and sub.campaign_post_format_id in (1, 2, 3, 11, 12, 42)
                   and campaign.platform_id = 1
                   and campaign.saas_collection_id is not null
                   and sub.handle_verified is null
                   and sub.created_on > '2023-10-04 00:00:00'
        ),
        post_submitted_profile_data as (
            select
                bip.shortcode as shortcode,
                argMin(bip.profile_id, bip.updated_at) as profile_id,
                argMin(bip.handle, bip.updated_at) as handle
            from
                dbt.stg_beat_instagram_post bip
                where bip.shortcode in (select shortCode from submission_data)
                group by shortcode
        ),
        user_profile_data as (
            select
                cp.id as cpId,
                ia.ig_id as profile_id,
                ia.handle as handle
            from
                dbt.stg_coffee_campaign_profiles cp
                left join dbt.stg_coffee_view_instagram_account_lite ia
                    on ia.id = cp.platform_account_id and cp.platform = 'INSTAGRAM'
                    where cp.id in (select cpId from submission_data)
        ),
        handle_verification_data as (
           select
                sd.submissionId as submissionId,
                upd.profile_id as submissionCPIGId,
                upd.handle as submissionCPHandle,
                pspd.profile_id as urlIGId,
                pspd.handle as urlCPHandle
           from submission_data sd
           inner join
               post_submitted_profile_data pspd
           on sd.shortCode = pspd.shortcode
           inner join
               user_profile_data upd
           on sd.cpId = upd.cpId
        )
        select * from handle_verification_data
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
        submissionId = row[0]
        submissionCPIGId = row[1]
        submissionCPHandle = row[2]
        urlIGId = row[3]
        urlCPHandle = row[4]

        post_obj = {
            "submissionId": submissionId,
            "submissionCPIGId": submissionCPIGId,
            "submissionCPHandle": submissionCPHandle,
            "urlIGId": urlIGId,
            "urlCPHandle": urlCPHandle,
        }

        json_message = json.dumps(post_obj)
        channel.basic_publish(exchange='campaign.dx', routing_key='campaign.handle_verification', body=json_message)
    connection.close()


with DAG(
        dag_id='upload_handle_verification_data',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval="0 * * * *",
) as dag:
    task = PythonOperator(
        task_id='upload_handle_verification',
        python_callable=handle_verification_for_posts
    )
    task
