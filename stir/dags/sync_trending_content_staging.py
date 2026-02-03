from airflow import DAG
from airflow.providers.postgres.operators.postgres import PostgresOperator
from airflow.providers.ssh.operators.ssh import SSHOperator
from airflow.utils.dates import days_ago
from airflow_clickhouse_plugin.operators.clickhouse import ClickHouseOperator
from slack_connection import SlackNotifier

with DAG(
    dag_id='sync_trending_content_staging',
    start_date=days_ago(1),
    catchup=False,
    max_active_runs=1,
    concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
    schedule_interval="0 4 * * *",
) as dag:
    fetch_trending_content_from_ch = ClickHouseOperator(
        task_id='fetch_trending_content_from_ch',
        database='vidooly',
        sql=(
            '''
                INSERT INTO FUNCTION s3('https://gcc-social-data.s3.ap-south-1.amazonaws.com/data-pipeline/tmp/trending_content.json', 'AKIAXGXUCIER7YGHF4X3', 'LNm/3YjC6L2dIZWSNGTahS9HtQvdJjlcme59ZWF1', 'JSONEachRow')
                SELECT * from dbt.mart_genre_trending_content
                SETTINGS s3_truncate_on_insert=1;
            '''
        ),
        clickhouse_conn_id='clickhouse_gcc',
    )
    download_cmd = "aws s3 cp s3://gcc-social-data/data-pipeline/tmp/trending_content.json /tmp/trending_content_raw.json && sed -e 's/\\\\/\\\\\\\\/g' /tmp/trending_content_raw.json > /tmp/trending_content.json"
    import_to_pg_box = SSHOperator(
        ssh_conn_id='ssh_stage_pg',
        cmd_timeout=1000,
        task_id='download_to_pg_local',
        command=download_cmd,
        dag=dag)
    import_to_pg_db = PostgresOperator(
        task_id="import_to_pg_db",
        runtime_parameters={'statement_timeout': '3000s'},
        postgres_conn_id="stage_pg",
        sql="""
            ALTER SEQUENCE trending_content_id_seq1 OWNED BY NONE;
            DROP TABLE IF EXISTS trending_content_tmp;
            DROP TABLE IF EXISTS mart_genre_trending_content;
            DROP TABLE IF EXISTS trending_content_old_bkp;
            
            CREATE TEMP TABLE mart_genre_trending_content(data jsonb);
                
            COPY mart_genre_trending_content from '/tmp/trending_content.json';

            create table trending_content_tmp (
                like trending_content
                including defaults
                including constraints
                including indexes
            );
            
            INSERT INTO trending_content_tmp
            (
                month,
                category,
                platform,
                language,
                profile_type,
                title,
                description,
                engagement,
                engagement_rate,
                views,
                plays,
                likes,
                comments,
                published_at,
                post_id,
                post_type,
                post_link,
                thumbnail,
                profile_social_id,
                plays_rank,
                views_rank,
                enabled
            )
            SELECT
                (data->>'month') AS month,
                (data->>'category') AS category,
                (data->>'platform') AS platform,
                (data->>'language') AS language,
                (data->>'profile_type') AS profile_type,
                (data->>'title') AS title,
                (data->>'description') AS description,
                (data->>'engagement')::int8 AS engagement,
                (data->>'engagement_rate')::float8 AS engagement_rate,
                (data->>'views')::int8 AS views,
                (data->>'plays')::int8 AS plays,
                (data->>'likes')::int8 AS likes,
                (data->>'comments')::int8 AS comments,
                (data->>'publish_time')::timestamp AS published_at,
                (data->>'short_code') AS post_id,
                (data->>'post_type') AS post_type,
                (data->>'post_url') AS post_link,
                (data->>'thumbnail_url') AS thumbnail,
                (data->>'profile_id') AS profile_social_id,
                (data->>'plays_rank')::int8 AS plays_rank,
                (data->>'views_rank')::int8 AS views_rank,
                TRUE enabled
            from mart_genre_trending_content
            where (data->>'short_code') != 'MISSING' AND (data->>'profile_id') != 'MISSING';
            
            update trending_content_tmp tcv
            set platform_account_id = ia.id, audience_gender = ia.audience_gender, audience_age = ia.audience_age
            from instagram_account ia
            where ia.ig_id = tcv.profile_social_id and tcv.profile_social_id is not null;
            
            update trending_content_tmp tcv
            set platform_account_id = ya.id, audience_gender = ya.audience_gender, audience_age = ya.audience_age
            from youtube_account ya
            where ya.channel_id = tcv.profile_social_id and tcv.profile_social_id is not null;
            
            BEGIN;
            ALTER TABLE "trending_content" RENAME TO "trending_content_old_bkp";
            ALTER TABLE "trending_content_tmp" RENAME TO "trending_content";
            DROP TABLE IF EXISTS trending_content_old_bkp;
            COMMIT;
        """,
    )
    fetch_trending_content_from_ch >> import_to_pg_box >> import_to_pg_db
