import pika
from airflow import DAG
from airflow.models import BaseOperator
from airflow.operators.python_operator import PythonOperator
from airflow.utils.dates import days_ago
from datetime import datetime, timedelta
from sqlalchemy import create_engine, text
from airflow.hooks.base import BaseHook
from airflow.hooks.postgres_hook import PostgresHook
import os
import json
import requests
import clickhouse_connect

from slack_connection import SlackNotifier

def upload_to_content_verification():
    connection = BaseHook.get_connection("clickhouse_gcc")
    extra_params = json.loads(connection.extra)
    client = clickhouse_connect.get_client(host=connection.host, password=connection.password, username=connection.login, send_receive_timeout=extra_params['send_receive_timeout'])
    max_posts = 250
    posts = client.query('''
                        with 
                        content_data as (
                            select sub.id                              as submissionId,
                                sub.campaign_profile_id                as cpId,
                                CASE
                                    WHEN sub.campaign_post_format_id = 2 THEN 'story'
                                    WHEN sub.campaign_post_format_id = 11 THEN 'reels'
                                    WHEN sub.campaign_post_format_id = 42 THEN 'reels'
                                    END                             as postype,
                                sub.short_code                      as shortCode,
                                deliverable.reference_url           as referenceUrl,
                                brief.auto_verification             as autoVerificationEnabled,
                                brief.content_auto_verification     as contentAutoVerificationEnabled,
                                brief.unique_link_auto_verification as uniqueAutoVerificationEnabled,
                                sub.created_on                      as submissionCreatedOn,
                                sub.campaign_id                     as campaingId
                            from dbt.stg_campaign_deliverable_submission sub
                                    inner join dbt.stg_campaign_campaign campaign on campaign.id = sub.campaign_id
                                    inner join dbt.stg_campaign_brief brief on brief.id = sub.brief_id
                                    inner join dbt.stg_campaign_campaign_deliverable deliverable on deliverable.id = sub.campaign_deliverable_id
                            where sub.status = 'UNDER_REVIEW'
                            and sub.level = 'FINAL_POST'
                            and sub.campaign_post_format_id in (2, 11, 42) --2 story 11 normal reel 42 repost reel
                            and campaign.campaign_type in ('REWARD_LADDER', 'REPOST_CHALLENGE')
                            AND campaign.platform_id = 1
                            and brief.auto_verification = 1
                            and (brief.content_auto_verification = 1 or brief.unique_link_auto_verification = 1)
                            ORDER BY sub.id DESC
                        ),
                        cpids as (
                            select toString(cpId) from content_data group by cpId
                        ),
                        profile_data as (
                            select  
                                sccp.id cpId,
                                sccp.platform as platform,
                                sccp.platform_account_id as accountId,
                                social_data.ig_id as profile_id
                            from dbt.stg_coffee_campaign_profiles sccp
                            left join dbt.stg_coffee_view_instagram_account_lite social_data on social_data.id = sccp.platform_account_id
                            where sccp.platform = 'INSTAGRAM' and toString(sccp.id) IN cpids
                        ),
                        final as (
                            select 
                                profile_data.cpId cpId,
                                profile_data.platform platform,
                                profile_data.profile_id profile_id,
                                content_data.submissionId submissionId,
                                content_data.postype postype,
                                content_data.shortCode shortCode,
                                content_data.referenceUrl referenceUrl
                            from content_data left join profile_data on profile_data.cpId = content_data.cpId
                            where profile_id is not null and profile_id != ''
                        )
                        select * from final
                        
                        '''
                            f"LIMIT {max_posts}")

    rabbitmq_conn = BaseHook.get_connection("rabbitmq_connection")
    rabbitmq_host = rabbitmq_conn.host
    rabbitmq_username = 'admin'
    rabbitmq_password = 'admin'
    rabbitmq_port = rabbitmq_conn.port
    rabbitmq_queue = "content_match_q"
    connection_params = pika.ConnectionParameters(
        host=rabbitmq_host,
        credentials=pika.PlainCredentials(rabbitmq_username, rabbitmq_password),
        port=rabbitmq_port
    )
    connection = pika.BlockingConnection(connection_params)
    channel = connection.channel()
    channel.queue_declare(queue=rabbitmq_queue, durable=True)
    for row in posts.result_rows:
        platform = row[1]
        profile_id = row[2]
        submission_id = row[3]
        post_type = row[4] 
        short_code = row[5] 
        reference_url = row[6]
        
        post_obj = {
            "type": post_type,
            "shortcode": short_code,
            "profile_id": profile_id,
            "submission_id": submission_id,
            "platform": platform
        }
        
        message = {
            "expected_content_url": reference_url,
            "post": post_obj
        }
        
        json_message = json.dumps(message)
        channel.basic_publish(exchange='', routing_key=rabbitmq_queue, body=json_message)
    connection.close()

with DAG(
        dag_id='upload_content_verification_data',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval="0 */3 * * *",
) as dag:
    task = PythonOperator(
        task_id='upload_content_verification',
        python_callable=upload_to_content_verification
    )
    task 
