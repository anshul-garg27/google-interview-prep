from airflow import DAG
from airflow.providers.postgres.operators.postgres import PostgresOperator
from airflow.providers.ssh.operators.ssh import SSHOperator
from airflow.utils.dates import days_ago
from airflow_clickhouse_plugin.operators.clickhouse import ClickHouseOperator
from slack_connection import SlackNotifier

with DAG(
        dag_id='sync_time_series_staging',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval="0 4 * * *",
) as dag:
    fetch_from_ch = ClickHouseOperator(
        task_id='fetch_from_ch',
        database='vidooly',
        sql=(
            '''
                INSERT INTO FUNCTION s3('https://gcc-social-data.s3.ap-south-1.amazonaws.com/data-pipeline/tmp/mart_time_series.json', 'AKIAXGXUCIER7YGHF4X3', 'LNm/3YjC6L2dIZWSNGTahS9HtQvdJjlcme59ZWF1', 'JSONEachRow')
                SELECT * from dbt.mart_time_series
                SETTINGS s3_truncate_on_insert=1;
            '''
        ),
        clickhouse_conn_id='clickhouse_gcc',
    )
    download_cmd = "aws s3 cp s3://gcc-social-data/data-pipeline/tmp/mart_time_series.json /tmp/mart_time_series.json"
    import_to_pg_box = SSHOperator(
        ssh_conn_id='ssh_stage_pg',
        cmd_timeout=1000,
        task_id='download_to_pg_local',
        command=download_cmd,
        dag=dag)
    import_to_pg_db = PostgresOperator(
        task_id="import_to_pg_db",
        runtime_parameters={'statement_timeout': '20000s'},
        postgres_conn_id="stage_pg",
        sql="""
                DROP TABLE IF EXISTS time_series_tmp;
                DROP TABLE IF EXISTS mart_time_series;
                DROP TABLE IF EXISTS social_profile_time_series_bkp;
                DROP TABLE IF EXISTS leaderboard_bkp;
                
                CREATE TEMP TABLE mart_time_series(data jsonb);
                
                COPY mart_time_series from '/tmp/mart_time_series.json';

                create table time_series_tmp (
                    like social_profile_time_series
                    including defaults
                    including constraints
                    including indexes
                );
                
                INSERT INTO time_series_tmp
                (
                    platform_id,
                    platform,
                    followers,
                    following,
                    views,
                    views_total,
                    plays,
                    plays_total,
                    engagement_rate,
                    uploads,
                    followers_change,
                    date,
                    prev_date,
                    monthly_stats,
                    enabled
                )
                SELECT
                    data->>'platform_id' platform_id,
                    data->>'platform' platform,
                    (data->>'followers')::int8 followers,
                    (data->>'following')::int8 following,
                    (data->>'views')::int8 views,
                    (data->>'views_total')::int8 views_total,
                    (data->>'plays')::int8 plays,
                    (data->>'plays_total')::int8 plays_total,
                    (data->>'engagement_rate')::float engagement_rate,
                    (data->>'uploads')::int8 uploads,
                    (data->>'followers_change')::int8 followers_change,
                    (data->>'date')::timestamp date,
                    (data->>'p_date')::timestamp prev_date,
                    (data->>'monthly_stats')::bool as monthly_stats,
                    TRUE enabled
                from mart_time_series;
                
                update time_series_tmp lv 
                set platform_profile_id  = ya.id
                from youtube_account ya
                where ya.channel_id  = lv.platform_id and lv.platform_id is not null;
                                
                update time_series_tmp lv 
                set platform_profile_id  = ia.id
                from instagram_account ia
                where ia.ig_id  = lv.platform_id and lv.platform_id is not null;
                
                BEGIN;
                ALTER TABLE "social_profile_time_series" RENAME TO "social_profile_time_series_bkp";
                ALTER TABLE "time_series_tmp" RENAME TO "social_profile_time_series";
                DROP TABLE IF EXISTS social_profile_time_series_bkp;
                COMMIT;
              """,
    )
    fetch_from_ch >> import_to_pg_box >> import_to_pg_db

