from airflow import DAG
from airflow.providers.postgres.operators.postgres import PostgresOperator
from airflow.providers.ssh.operators.ssh import SSHOperator
from airflow.utils.dates import days_ago
from airflow_clickhouse_plugin.operators.clickhouse import ClickHouseOperator
from slack_connection import SlackNotifier

with DAG(
    dag_id='sync_genre_overview_prod',
    start_date=days_ago(1),
    catchup=False,
    max_active_runs=1,
    concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
    schedule_interval="0 0 * * *",
) as dag:
    fetch_genre_overview_from_ch = ClickHouseOperator(
        task_id='fetch_genre_overview_from_ch',
        database='vidooly',
        sql=(
            '''
                INSERT INTO FUNCTION s3('https://gcc-social-data.s3.ap-south-1.amazonaws.com/data-pipeline/tmp/genre_overview.json', 'AKIAXGXUCIER7YGHF4X3', 'LNm/3YjC6L2dIZWSNGTahS9HtQvdJjlcme59ZWF1', 'JSONEachRow')
                SELECT * from dbt.mart_genre_overview
                SETTINGS s3_truncate_on_insert=1, output_format_json_escape_forward_slashes=0;
            '''
        ),
        clickhouse_conn_id='clickhouse_gcc',
    )
    download_cmd = "aws s3 cp s3://gcc-social-data/data-pipeline/tmp/genre_overview.json /tmp/genre_overview.json"
    import_to_pg_box = SSHOperator(
        ssh_conn_id='ssh_prod_pg',
        cmd_timeout=1000,
        task_id='download_to_pg_local',
        command=download_cmd,
        dag=dag)
    import_to_pg_db = PostgresOperator(
        task_id="import_to_pg_db",
        runtime_parameters={'statement_timeout': '3000s'},
        postgres_conn_id="prod_pg",
        sql="""
            ALTER SEQUENCE genre_overview_id_seq OWNED BY NONE;
            DROP TABLE IF EXISTS genre_overview_tmp;
            DROP TABLE IF EXISTS mart_genre_overview;
            DROP TABLE IF EXISTS genre_overview_old_bkp;
            
            CREATE TEMP TABLE mart_genre_overview(data jsonb);
                
            COPY mart_genre_overview from '/tmp/genre_overview.json';

            create table genre_overview_tmp (
                like genre_overview
                including defaults
                including constraints
                including indexes
            );
            
            INSERT INTO genre_overview_tmp
            (
                category,
                month,
                platform,
                language,
                profile_type,
                creators,
                uploads,
                views,
                engagement,
                engagement_rate,
                followers,
                likes,
                comments,
                audience_age_gender,
                audience_age,
                audience_gender,
                enabled
            )
            SELECT
                (data->>'category') AS category,
                (data->>'month') AS month,
                (data->>'platform') AS platform,
                (data->>'language') AS language,
                (data->>'profile_type') AS profile_type,
                (data->>'creators')::int8 AS creators,
                (data->>'uploads')::int8 AS uploads,
                (data->>'views')::int8 AS views,
                (data->>'engagement')::int8 AS engagement,
                (data->>'engagement_rate')::float8 AS engagement_rate,
                (data->>'followers')::int8 AS followers,
                (data->>'likes')::int8 AS likes,
                (data->>'comments')::int8 AS comments,
                (data->>'audience_age_gender')::jsonb AS audience_age_gender,
                (data->>'audience_age')::jsonb AS audience_age,
                (data->>'audience_gender')::jsonb AS audience_gender,
                TRUE enabled
            from mart_genre_overview
            where (data->>'language') != 'MISSING' AND (data->>'category') != 'MISSING';
                        
            BEGIN;
            ALTER TABLE "genre_overview" RENAME TO "genre_overview_old_bkp";
            ALTER TABLE "genre_overview_tmp" RENAME TO "genre_overview";
            COMMIT;
        """,
    )
    fetch_genre_overview_from_ch >> import_to_pg_box >> import_to_pg_db
        
        