from airflow import DAG
from airflow.providers.postgres.operators.postgres import PostgresOperator
from airflow.providers.ssh.operators.ssh import SSHOperator
from airflow.utils.dates import days_ago
from airflow_clickhouse_plugin.operators.clickhouse import ClickHouseOperator
from slack_connection import SlackNotifier

with DAG(
        dag_id='sync_cp_keywords_prod',
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
                INSERT INTO FUNCTION s3('https://gcc-social-data.s3.ap-south-1.amazonaws.com/data-pipeline/tmp/cp_keywords.json', 'AKIAXGXUCIER7YGHF4X3', 'LNm/3YjC6L2dIZWSNGTahS9HtQvdJjlcme59ZWF1', 'JSONEachRow')
                with keywords as (
                    select
                        collection_id,
                        arrayJoin(keywords) keyword,
                        count(*) count,
                        row_number() over (partition by collection_id order by count desc) rank
                    from dbt.mart_collection_post
                    where lower(keyword) not IN 
                    ('your', 'from', 'this', 'stop', 'with', 'that', 'have', 'https', 'using', 
                    'also', 'which', 'these', 'will', 'what', 'this', 'about',
                        'more', 'their', 'they', 'used', 'what', 'some', 'just', 'like')
                    group by collection_id, keyword
                    order by collection_id desc, rank asc
                )
                select * from keywords where rank <= 100
                SETTINGS s3_truncate_on_insert=1;
            '''
        ),
        clickhouse_conn_id='clickhouse_gcc',
    )
    download_cmd = "aws s3 cp s3://gcc-social-data/data-pipeline/tmp/cp_keywords.json /tmp/cp_keywords.json"
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
                DROP TABLE IF EXISTS collection_keyword_tmp;
                DROP TABLE IF EXISTS mart_collection_keywords;
                DROP TABLE IF EXISTS collection_keyword_old_bkp;

                CREATE TEMP TABLE mart_collection_keywords(data jsonb);

                COPY mart_collection_keywords from '/tmp/cp_keywords.json';
                
                create table collection_keyword_tmp (
                    like collection_keyword
                    including defaults
                    including constraints
                    including indexes
                );
                
                alter table collection_keyword_tmp alter column collection_share_id drop not null;
                alter table collection_keyword_tmp alter column collection_type drop not null;

                INSERT INTO collection_keyword_tmp
                    (
                        collection_id,
                        keyword_name,
                        tagged_count
                    )
                    select
                        (data->>'collection_id') collection_id,
                        (data->>'keyword') keyword_name,
                        (data->>'count')::int8 tagged_count
                    from mart_collection_keywords;
                    
                    update collection_keyword_tmp
                    set collection_share_id = pc.share_id, collection_type = 'PROFILE'
                    from profile_collection pc
                    where pc.id::varchar = collection_id;
                    
                    update collection_keyword_tmp
                    set collection_share_id = pc.share_id, collection_type = 'POST'
                    from post_collection pc
                    where pc.id = collection_id;
                    
                    alter table collection_keyword_tmp alter column collection_share_id set not null;
                    alter table collection_keyword_tmp alter column collection_type set not null;
                    
                    BEGIN;
                        ALTER TABLE "collection_keyword" RENAME TO "collection_keyword_old_bkp";
                        ALTER TABLE "collection_keyword_tmp" RENAME TO "collection_keyword";
                    COMMIT;
              """,
    )
    fetch_from_ch >> import_to_pg_box >> import_to_pg_db

