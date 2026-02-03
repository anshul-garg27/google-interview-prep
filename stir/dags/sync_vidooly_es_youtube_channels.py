import datetime as dt
import json
from datetime import date

from airflow import DAG
from airflow.operators.bash import BashOperator
from airflow.utils.dates import days_ago
from airflow_clickhouse_plugin.operators.clickhouse import ClickHouseOperator
from slack_connection import SlackNotifier

with DAG(
        dag_id="sync_vidooly_es_youtube_channels",
        schedule_interval="0 0 * * *",
        start_date=days_ago(1),
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        catchup=False,
        dagrun_timeout=dt.timedelta(hours=6),
) as dag:
    today = date.today().strftime("%Y%m%d")
    es_search_query = {
                "query": {
                    "match_phrase": {
                        "snippet.country": "IN"
                        }
                    }
                }
    upload_data_to_s3 = BashOperator(
        task_id="fetch_data",
        bash_command=f"elasticdump \
            --s3AccessKeyId=AKIAXGXUCIERUOPVUUPW \
            --s3SecretAccessKey=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 \
            --input=http://elastic:accelastic123@52.6.85.30:8091/youtube_channels_api \
            --output=s3://gcc-social-data/temp/youtube_channels_api_{today}.json \
            --limit=10000 \
            --type=data \
            --searchBody='{json.dumps(es_search_query)}'",
    )
    upload_to_ch = ClickHouseOperator(
        task_id='upload_to_ch',
        database='vidooly',
        sql=(
            f'CREATE TABLE vidooly.youtube_channels_data_{today} '
            '''
                ENGINE = MergeTree
                ORDER BY id AS
                SELECT
                    _path,
                    _file,
                    JSONExtractString(toString(_source), 'id') AS id,
                    JSONExtractString(toString(_source), 'created_at') AS created,
                    JSONExtractString(toString(_source), 'snippet.country') AS country,
                    JSONExtractString(toString(_source), 'snippet.customUrl') AS custom_url,
                    JSONExtractString(toString(_source), 'snippet.description') AS description,
                    JSONExtractString(toString(_source), 'snippet.title') AS title,
                    JSONExtractString(toString(_source), 'snippet.publiashedAt') AS published_at,
                    JSONExtractString(toString(_source), 'snippet.thumbnails.default.url') AS thumbnail,
                    JSONExtractString(toString(_source), 'last_updated') AS last_updated_at,
                    JSONExtractString(toString(_source), 'statistics.videoCount') AS video_count,
                    JSONExtractString(toString(_source), 'statistics.viewCount') AS view_count,
                    JSONExtractBool(toString(_source), 'statistics.hiddenSubscriberCount') AS hidden_subscriber,
                    JSONExtractString(toString(_source), 'statistics.subscriberCount') AS subscriber_count'''
            f" FROM s3('https://gcc-social-data.s3.ap-south-1.amazonaws.com/temp/youtube_channels_api_{today}.json', 'AKIAXGXUCIER5DX7YI4F', 'WnVvyEs+3NFFcNWnEGJr3Ctt/KuMUmaRYG9xbPWj', 'JSONEachRow', '_source JSON')"
        ),
        clickhouse_conn_id='clickhouse_gcc',
    )
    upload_data_to_s3 >> upload_to_ch
