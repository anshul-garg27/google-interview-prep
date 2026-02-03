from airflow import DAG
from airflow.providers.postgres.operators.postgres import PostgresOperator
from airflow.providers.ssh.operators.ssh import SSHOperator
from airflow.utils.dates import days_ago
from airflow_clickhouse_plugin.operators.clickhouse import ClickHouseOperator
from slack_connection import SlackNotifier

with DAG(
        dag_id='sync_hashtags_prod',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval="0 0 * * *",
) as dag:
    fetch_from_ch = ClickHouseOperator(
        task_id='fetch_hashtags_from_ch',
        database='vidooly',
        sql=(
            '''
                INSERT INTO FUNCTION s3('https://gcc-social-data.s3.ap-south-1.amazonaws.com/data-pipeline/tmp/hashtags.json', 'AKIAXGXUCIER7YGHF4X3', 'LNm/3YjC6L2dIZWSNGTahS9HtQvdJjlcme59ZWF1', 'JSONEachRow')
                select 
                    platform,
                    platform_id,
                    hashtags,
                    hashtags_counts
                from dbt.mart_instagram_hashtags
                SETTINGS s3_truncate_on_insert=1;
            '''
        ),
        clickhouse_conn_id='clickhouse_gcc',
    )
    download_cmd = "aws s3 cp s3://gcc-social-data/data-pipeline/tmp/hashtags.json /tmp/hashtags.json"
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
                DROP TABLE IF EXISTS social_profile_hashtags_tmp;
                DROP TABLE IF EXISTS mart_hashtags;
                DROP TABLE IF EXISTS social_profile_hashtags_old_bkp;
                
                CREATE TEMP TABLE mart_hashtags(data jsonb);
                
                COPY mart_hashtags from '/tmp/hashtags.json';

                create table social_profile_hashtags_tmp (
                    like social_profile_hashtags
                    including defaults
                    including constraints
                    including indexes
                );
                
                INSERT INTO social_profile_hashtags_tmp
                (
                    platform_id,
                    platform,
                    hashtags,
                    hashtags_counts
                )
                SELECT
                    data->>'platform_id' platform_id,
                    data->>'platform' platform,
                    jsonb_array_to_text_array((data->>'hashtags')::jsonb) 	hashtags,
                    jsonb_array_to_text_array((data->>'hashtags_counts')::jsonb) 	hashtags_counts
                from mart_hashtags;
                
                update social_profile_hashtags_tmp lv 
                set platform_profile_id  = ia.id
                from instagram_account ia
                where ia.ig_id  = lv.platform_id;
                
                BEGIN;
                ALTER TABLE "social_profile_hashtags" RENAME TO "social_profile_hashtags_old_bkp";
                ALTER TABLE "social_profile_hashtags_tmp" RENAME TO "social_profile_hashtags";
                COMMIT;
              """,
    )
    fetch_from_ch >> import_to_pg_box >> import_to_pg_db

