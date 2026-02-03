from airflow import DAG
from airflow.providers.postgres.operators.postgres import PostgresOperator
from airflow.providers.ssh.operators.ssh import SSHOperator
from airflow.utils.dates import days_ago
from airflow_clickhouse_plugin.operators.clickhouse import ClickHouseOperator
from slack_connection import SlackNotifier

with DAG(
        dag_id='sync_cp_hashtags_prod',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval="0 0 * * *",
) as dag:
    fetch_from_ch = ClickHouseOperator(
        task_id='fetch_from_ch',
        database='vidooly',
        sql=(
            '''
                INSERT INTO FUNCTION s3('https://gcc-social-data.s3.ap-south-1.amazonaws.com/data-pipeline/tmp/cp_hashtags.json', 'AKIAXGXUCIER7YGHF4X3', 'LNm/3YjC6L2dIZWSNGTahS9HtQvdJjlcme59ZWF1', 'JSONEachRow')
                with hashtags as (
                    select
                        collection_id,
                        arrayJoin(hashtags) hashtag,
                        count(*) count,
                        row_number() over (partition by collection_id order by count desc) rank
                    from dbt.mart_collection_post
                    group by collection_id, hashtag
                    order by collection_id desc, rank asc
                )
                select * from hashtags where rank <= 100
                SETTINGS s3_truncate_on_insert=1;
            '''
        ),
        clickhouse_conn_id='clickhouse_gcc',
    )
    download_cmd = "aws s3 cp s3://gcc-social-data/data-pipeline/tmp/cp_hashtags.json /tmp/cp_hashtags.json"
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
                DROP TABLE IF EXISTS collection_hashtag_tmp;
                DROP TABLE IF EXISTS mart_collection_hashtags;
                DROP TABLE IF EXISTS collection_hashtag_old_bkp;

                CREATE TEMP TABLE mart_collection_hashtags(data jsonb);

                COPY mart_collection_hashtags from '/tmp/cp_hashtags.json';
                
                create table collection_hashtag_tmp (
                    like collection_hashtag
                    including defaults
                    including constraints
                    including indexes
                );
                
                alter table collection_hashtag_tmp alter column collection_share_id drop not null;
                alter table collection_hashtag_tmp alter column collection_type drop not null;

                INSERT INTO collection_hashtag_tmp
                    (
                        collection_id,
                        hashtag_name,
                        tagged_count,
                        ugc_tagged_count
                    )
                    select
                        (data->>'collection_id') collection_id,
                        (data->>'hashtag') hashtag_name,
                        (data->>'count')::int8 tagged_count,
                        (data->>'count')::int8 ugc_tagged_count
                    from mart_collection_hashtags;
                    
                    update collection_hashtag_tmp
                    set collection_share_id = pc.share_id, collection_type = 'PROFILE'
                    from profile_collection pc
                    where pc.id::varchar = collection_id;
                    
                    update collection_hashtag_tmp
                    set collection_share_id = pc.share_id, collection_type = 'POST'
                    from post_collection pc
                    where pc.id = collection_id;
                    
                    alter table collection_hashtag_tmp alter column collection_share_id set not null;
                    alter table collection_hashtag_tmp alter column collection_type set not null;
                    
                    BEGIN;
                        ALTER TABLE "collection_hashtag" RENAME TO "collection_hashtag_old_bkp";
                        ALTER TABLE "collection_hashtag_tmp" RENAME TO "collection_hashtag";
                    COMMIT;
              """,
    )
    fetch_from_ch >> import_to_pg_box >> import_to_pg_db

