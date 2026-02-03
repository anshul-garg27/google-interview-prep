from airflow import DAG
from airflow.providers.postgres.operators.postgres import PostgresOperator
from airflow.providers.ssh.operators.ssh import SSHOperator
from airflow.utils.dates import days_ago
from airflow_clickhouse_plugin.operators.clickhouse import ClickHouseOperator
from slack_connection import SlackNotifier

with DAG(
        dag_id='sync_cp_ts_prod',
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
                INSERT INTO FUNCTION s3('https://gcc-social-data.s3.ap-south-1.amazonaws.com/data-pipeline/tmp/cpts.json', 'AKIAXGXUCIER7YGHF4X3', 'LNm/3YjC6L2dIZWSNGTahS9HtQvdJjlcme59ZWF1', 'JSONEachRow')
                select 
                    platform,
                    post_short_code,
                    post_type,
                    '' post_title,
                    post_link,
                    post_thumbnail,
                    published_at,
                    profile_social_id,
                    profile_handle,
                    profile_name,
                    collection_type,
                    collection_id,
                    post_collection_item_id,
                    followers,
                    link_clicks,
                    orders,
                    stats_date,
                    views,
                    likes,
                    comments,
                    impressions,
                    saves,
                    plays,
                    reach,
                    swipe_ups,
                    mentions,
                    sticker_taps,
                    shares,
                    story_exits,
                    story_back_taps,
                    story_forward_taps,
                    total_engagement,
                    engagement_rate
                from dbt.mart_collection_post_ts
                SETTINGS s3_truncate_on_insert=1;
            '''
        ),
        clickhouse_conn_id='clickhouse_gcc',
    )
    download_cmd = "aws s3 cp s3://gcc-social-data/data-pipeline/tmp/cpts.json /tmp/cpts.json"
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
                DROP TABLE IF EXISTS collection_post_metrics_ts_tmp;
                DROP TABLE IF EXISTS mart_cp_ts;
                DROP TABLE IF EXISTS collection_post_metrics_ts_old_bkp;
                
                CREATE TEMP TABLE mart_cp_ts(data jsonb);
                
                COPY mart_cp_ts from '/tmp/cpts.json';
                
                create UNLOGGED table collection_post_metrics_ts_tmp (
                    like collection_post_metrics_ts
                    including defaults
                    including constraints
                    including indexes
                );
                
                alter table collection_post_metrics_ts_tmp alter column collection_share_id drop not null;
                alter table collection_post_metrics_ts_tmp alter column platform_account_code drop not null;
                
                INSERT INTO collection_post_metrics_ts_tmp
                    (
                        platform,
                        post_short_code,
                        post_type,
                        post_title,
                        post_link,
                        post_thumbnail,
                        published_at,
                        profile_social_id,
                        profile_handle,
                        profile_name,
                        collection_type,
                        collection_id,
                        post_collection_item_id,
                        link_clicks,
                        orders,
                        stats_date,
                        views,
                        likes,
                        comments,
                        impressions,
                        saves,
                        plays,
                        reach,
                        swipe_ups,
                        mentions,
                        sticker_taps,
                        shares,
                        story_exits,
                        story_back_taps,
                        story_forward_taps,
                        er
                    )
                    select
                       (data->>'platform') platform,
                       (data->>'post_short_code') post_short_code,
                       (data->>'post_type') post_type,
                       (data->>'post_title') post_title,
                       coalesce((data->>'post_link'),'') post_link,
                       (data->>'post_thumbnail') post_thumbnail,
                       coalesce((data->>'published_at')::timestamp, date('1970-01-01')) published_at,
                       (data->>'profile_social_id') profile_social_id,
                       (data->>'profile_handle') profile_handle,
                       (data->>'profile_name') profile_name,
                       (data->>'collection_type') collection_type,
                       (data->>'collection_id') collection_id,
                       (data->>'post_collection_item_id')::int8 post_collection_item_id,
                       (data->>'link_clicks')::int8 link_clicks,
                       (data->>'orders')::int8 orders,
                       (data->>'stats_date')::timestamp stats_date,
                       round(coalesce(((data->>'views')::float8),0.0)) views,
                       (data->>'likes')::int8 likes,
                       (data->>'comments')::int8 comments,
                       round((data->>'impressions')::float8) impressions,
                       (data->>'saves')::int8 saves,
                       (data->>'plays')::int8 plays,
                       round(coalesce(((data->>'reach')::float8), 0.0)) reach,
                       (data->>'swipe_ups')::int8 swipe_ups,
                       (data->>'mentions')::int8 mentions,
                       (data->>'sticker_taps')::int8 sticker_taps,
                       (data->>'shares')::int8 shares,
                       (data->>'story_exits')::int8 story_exits,
                       (data->>'story_back_taps')::int8 story_back_taps,
                       (data->>'story_forward_taps')::int8 story_forward_taps,
                       coalesce((data->>'engagement_rate')::float8, 0.0) er
                    from mart_cp_ts
                    where (data->>'profile_social_id') is not null;
                    
                    update collection_post_metrics_ts_tmp
                    set platform_account_code = ia.id, profile_pic = ia.thumbnail, followers = coalesce(ia.followers, 0)
                    from instagram_account ia
                    where ia.ig_id = profile_social_id;
                    
                    update collection_post_metrics_ts_tmp
                    set platform_account_code = ya.id, profile_pic = ya.thumbnail, followers = coalesce(ya.followers, 0)
                    from youtube_account ya
                    where ya.channel_id = profile_social_id;
                    
                    update collection_post_metrics_ts_tmp
                    set collection_share_id = pc.share_id
                    from profile_collection pc
                    where pc.id = collection_id::int8 and collection_type = 'PROFILE';
                    
                    update collection_post_metrics_ts_tmp
                    set collection_share_id = pc.share_id
                    from post_collection pc
                    where pc.id = collection_id and collection_type = 'POST';
                    
                    DELETE FROM collection_post_metrics_ts_tmp where platform_account_code is null;
                    alter table collection_post_metrics_ts_tmp alter column collection_share_id set not null;
                    alter table collection_post_metrics_ts_tmp alter column platform_account_code set not null;
                    
                    ALTER TABLE collection_post_metrics_ts_tmp SET LOGGED; 
                    
                    BEGIN;
                        ALTER TABLE "collection_post_metrics_ts" RENAME TO "collection_post_metrics_ts_old_bkp";
                        ALTER TABLE "collection_post_metrics_ts_tmp" RENAME TO "collection_post_metrics_ts";
                    COMMIT;
              """,
    )
    fetch_from_ch >> import_to_pg_box >> import_to_pg_db

