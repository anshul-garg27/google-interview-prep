import os

from datetime import datetime

import clickhouse_connect
import loguru

from core.amqp.amqp import publish


async def generate_post_collection_sentiment_report(collection_id: str, collection_type: str):
    client = clickhouse_connect.get_client(host=os.environ["CLICKHOUSE_URL"],
                                           username=os.environ["CLICKHOUSE_USERNAME"],
                                           password=os.environ["CLICKHOUSE_PASSWORD"])
    environment = os.environ["ENVIRONMENT"]
    current_date = datetime.today()
    parquet_file_path = f"sentiment_report/{environment}/{current_date.year}/{current_date.month}/{current_date.day}/{collection_type}_{collection_id}.parquet"

    query = """
                    INSERT INTO FUNCTION s3('https://%s.s3.ap-south-1.amazonaws.com/%s', '%s', '%s', 'Parquet')
                        with collection as (select *
                                            from dbt.stg_coffee_post_collection_item scpci
                                                     inner join dbt.stg_coffee_post_collection scpc
                                                                on scpc.id = scpci.post_collection_id
                                            where post_collection_id = '%s'),
                             shortcodes as (select short_code
                                            from collection),
                             sentiment as (select
                                            argMax(sle.shortcode, sle.event_timestamp) shortcode,
                                            sle.comment_id comment_id,
                                            argMax(sle.source, sle.event_timestamp) source,
                                            argMax(pale.actor_profile_id, pale.event_timestamp) actor_profile_id,
                                            argMax(sle.comment, sle.event_timestamp) comment,
                                            argMax(JSONExtractString(pale.dimensions, 'keywords') , pale.event_timestamp) keywords,
                                            argMax(sle.sentiment, sle.event_timestamp) sentiment,
                                            argMax(sle.score, sle.event_timestamp) score,
                                            argMax(sle.platform, sle.event_timestamp) platform
                                           from %s.sentiment_log_events sle
                                           left join %s.post_activity_log_events pale on sle.comment_id = JSONExtractString(pale.dimensions, 'comment_id')
                                           where sle.shortcode in shortcodes
                                           group by sle.comment_id)

                        select post_collection_id collection_id,
                               '%s'             collection_type,
                               source,
                               platform,
                               shortcode,
                               comment_id,
                               comment,
                               keywords,
                               sentiment,
                               round(score, 2) score,
                               post_type,
                               actor_profile_id
                        from sentiment
                                 left join collection on collection.short_code = sentiment.shortcode
                            SETTINGS s3_truncate_on_insert = 1;
                                """ % ( os.environ["AWS_BUCKET"],
        parquet_file_path, os.environ["AWS_ACCESS_KEY"], os.environ["AWS_SECRET_ACCESS_KEY"], collection_id,
        os.environ["EVENT_DB_NAME"], os.environ["EVENT_DB_NAME"], collection_type)

    client.query(query)
    payload = {
        "collectionType": collection_type,
        "collectionId": collection_id,
        "sentimentReportPath": parquet_file_path,
        "sentimentReportBucket": os.environ['AWS_BUCKET']
    }
    publish(payload, "coffee.dx", "sentiment_collection_report_out_rk")


async def generate_sentiment_report(payload, session=None):
    loguru.logger.debug(payload)

    collection_type = payload['collectionType']
    collection_id = payload['collectionId']
    if collection_type == "POST":
        await generate_post_collection_sentiment_report(collection_id, collection_type)
