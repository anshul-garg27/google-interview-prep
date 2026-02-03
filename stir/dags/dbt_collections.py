import datetime as dt
from datetime import datetime
from airflow import DAG
from airflow.providers.postgres.operators.postgres import PostgresOperator
from airflow.providers.ssh.operators.ssh import SSHOperator
from airflow.utils.dates import days_ago
from airflow_clickhouse_plugin.operators.clickhouse import ClickHouseOperator
from airflow_dbt_python.operators.dbt import DbtRunOperator
from slack_connection import SlackNotifier
from airflow.operators.python import ShortCircuitOperator

def check_hour():
    hour = datetime.now().hour
    if 19 <= hour < 22:
        return False
    return True

with DAG(
        dag_id="dbt_collections",
        schedule_interval="*/30 * * * *",
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        dagrun_timeout=dt.timedelta(minutes=180),
) as dag:
    check_hour_task = ShortCircuitOperator(
        task_id="check_hour",
        python_callable=check_hour,
    )
    
    truncate_fake_ts = ClickHouseOperator(
        task_id='truncate_fake_ts',
        database='dbt',
        sql=(
            '''
                truncate table dbt.mart_fake_events
            '''
        ),
        clickhouse_conn_id='clickhouse_gcc',
    )
    create_ts = DbtRunOperator(
        task_id="create_ts",
        project_dir="/gcc/airflow/stir/gcc_social/",
        profiles_dir="/gcc/airflow/.dbt/",
        select=["tag:collections"],
        exclude=["tag:deprecated", "tag:post_ranker", "tag:core"],
        target="production",
        profile="gcc_warehouse",
    )
    # fill_gaps = DbtRunOperator(
    #     task_id="fill_gaps",
    #     project_dir="/gcc/airflow/stir/gcc_social/",
    #     profiles_dir="/gcc/airflow/.dbt/",
    #     select=["tag:collections_fake_events"],
    #     exclude=["tag:deprecated", "tag:post_ranker", "tag:core"],
    #     target="production",
    #     profile="gcc_warehouse",
    # )
    # create_final_ts = DbtRunOperator(
    #     task_id="create_final_ts",
    #     project_dir="/gcc/airflow/stir/gcc_social/",
    #     profiles_dir="/gcc/airflow/.dbt/",
    #     select=["tag:collections"],
    #     exclude=["tag:deprecated", "tag:post_ranker", "tag:core"],
    #     target="production",
    #     profile="gcc_warehouse",
    # )

    fetch_from_ch = ClickHouseOperator(
        task_id='fetch_from_ch',
        database='vidooly',
        sql=(
            '''
                INSERT INTO FUNCTION s3('https://gcc-social-data.s3.ap-south-1.amazonaws.com/data-pipeline/tmp/cps.json', 'AKIAXGXUCIER7YGHF4X3', 'LNm/3YjC6L2dIZWSNGTahS9HtQvdJjlcme59ZWF1', 'JSONEachRow')
                select 
                    platform,
                    post_short_code,
                    post_type,
                    post_title post_title,
                    replaceRegexpAll(post_link, '[^a-zA-Z0-9_%\/\:\.\-\?]', '') post_link,
                    replaceRegexpAll(post_thumbnail, '[^a-zA-Z0-9_%\/\:\.\-]', '') post_thumbnail,
                    published_at,
                    profile_social_id,
                    coalesce(profile_handle, '') profile_handle,
                    profile_name profile_name,
                    followers,
                    collection_id,
                    collection_type,
                    post_collection_item_id,
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
                    link_clicks,
                    orders,
                    delivered_orders,
                    completed_orders,
                    leaderboard_overall_orders,
                    leaderboard_delivered_orders,
                    leaderboard_completed_orders,
                    total_engagement,
                    engagement_rate,
                    arrayMap(x -> lower(x), hashtags) hashtags
                from dbt.mart_collection_post
                SETTINGS s3_truncate_on_insert=1;
            '''
        ),
        clickhouse_conn_id='clickhouse_gcc',
    )
    download_cmd = "aws s3 cp s3://gcc-social-data/data-pipeline/tmp/cps.json /tmp/cps_raw.json && sed -e 's/\\\\/\\\\\\\\/g' /tmp/cps_raw.json > /tmp/cps.json"

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
                    DROP TABLE IF EXISTS collection_post_metrics_summary_tmp;
                    DROP TABLE IF EXISTS mart_cps;
                    DROP TABLE IF EXISTS collection_post_metrics_summary_old_bkp;

                    CREATE TEMP TABLE mart_cps(data jsonb);

                    COPY mart_cps from '/tmp/cps.json';

                    create table collection_post_metrics_summary_tmp (
                        like collection_post_metrics_summary
                        including defaults
                        including constraints
                        including indexes
                    );

                    alter table collection_post_metrics_summary_tmp alter column collection_share_id drop not null;

                    INSERT INTO collection_post_metrics_summary_tmp
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
                            followers,
                            collection_id,
                            post_collection_item_id,
                            collection_type,
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
                            link_clicks,
                            orders,
                            er,
                            hashtags,
                            delivered_orders,
                            completed_orders,
                            leaderboard_overall_orders,
                            leaderboard_delivered_orders,
                            leaderboard_completed_orders
                        )
                        select
                            (data->>'platform') platform,
                            (data->>'post_short_code') post_short_code,
                            coalesce((data->>'post_type'), '') post_type,
                            coalesce((data->>'post_title'), '') post_title,
                            coalesce((data->>'post_link'), '') post_link,
                            (data->>'post_thumbnail') post_thumbnail,
                            coalesce((data->>'published_at')::timestamp, date('1970-01-01')) published_at,
                            (data->>'profile_social_id') profile_social_id,
                            (data->>'profile_handle') profile_handle,
                            coalesce((data->>'profile_name'), '') profile_name,
                            coalesce((data->>'followers')::int8,0) followers,
                            (data->>'collection_id') collection_id,
                            (data->>'post_collection_item_id')::int8 post_collection_item_id,
                            (data->>'collection_type') collection_type,
                            round(coalesce((data->>'views')::float8, 0.0)) views,
                            (data->>'likes')::int8 likes,
                            (data->>'comments')::int8 comments,
                            round((data->>'impressions')::float8) impressions,
                            (data->>'saves')::int8 saves,
                            (data->>'plays')::int8 plays,
                            coalesce(round((data->>'reach')::float8), 0.0) reach,
                            (data->>'swipe_ups')::int8 swipe_ups,
                            (data->>'mentions')::int8 mentions,
                            (data->>'sticker_taps')::int8 sticker_taps,
                            (data->>'shares')::int8 shares,
                            (data->>'story_exits')::int8 story_exits,
                            (data->>'story_back_taps')::int8 story_back_taps,
                            (data->>'story_forward_taps')::int8 story_forward_taps,
                            coalesce((data->>'link_clicks')::int8, 0) link_clicks,
                            coalesce((data->>'orders')::int8, 0) orders,
                            coalesce((data->>'engagement_rate')::float8, 0.0) er,
                            (data->>'hashtags')::jsonb hashtags,
                            coalesce((data->>'delivered_orders')::int8, 0) delivered_orders,
                            coalesce((data->>'completed_orders')::int8, 0) completed_orders,
                            coalesce((data->>'leaderboard_overall_orders')::int8, 0) leaderboard_overall_orders,
                            coalesce((data->>'leaderboard_delivered_orders')::int8, 0) leaderboard_delivered_orders,
                            coalesce((data->>'leaderboard_completed_orders')::int8, 0) leaderboard_completed_orders
                        from mart_cps;

                        alter table collection_post_metrics_summary_tmp add column IF NOT EXISTS show_in_report boolean default false;
                        alter table collection_post_metrics_summary_tmp add column IF NOT EXISTS posted_by_cp_id int8 default null;

                        update collection_post_metrics_summary_tmp 
                            SET show_in_report = post_collection_item.show_in_report,
                            posted_by_cp_id = post_collection_item.posted_by_cp_id 
                        from post_collection_item where collection_post_metrics_summary_tmp.collection_id = post_collection_item.post_collection_id and 
                        collection_post_metrics_summary_tmp.post_short_code = post_collection_item.short_code and
                        collection_post_metrics_summary_tmp.collection_type = 'POST';

                        update collection_post_metrics_summary_tmp
                        set profile_pic = ia.thumbnail, 
                            followers = coalesce(ia.followers, 0),
                            profile_name = coalesce(ia.name, ia.handle, '')
                        from instagram_account ia
                        where ia.ig_id = profile_social_id and profile_social_id is not null and ia.ig_id is not null and ia.ig_id != '';

                        update collection_post_metrics_summary_tmp
                        set profile_pic = ya.thumbnail, 
                            followers = coalesce(ya.followers, 0),
                            profile_name = coalesce(ya.username, ya.title, '')
                        from youtube_account ya
                        where ya.channel_id = profile_social_id and profile_social_id is not null and ya.channel_id is not null and ya.channel_id != '';

                        update collection_post_metrics_summary_tmp
                        set collection_share_id = pc.share_id
                        from profile_collection pc
                        where pc.id::varchar = collection_id and collection_type = 'PROFILE';

                        update collection_post_metrics_summary_tmp
                        set collection_share_id = pc.share_id
                        from post_collection pc
                        where pc.id = collection_id and collection_type = 'POST';

                        update collection_post_metrics_summary_tmp
                        set cost = pci.cost
                        from post_collection_item pci
                        where pci.id = post_collection_item_id and post_collection_item_id > 0;
                        
                        delete from collection_post_metrics_summary_tmp where collection_type = 'POST' and show_in_report = false;
                        delete from collection_post_metrics_summary_tmp where collection_share_id is null;
                        alter table collection_post_metrics_summary_tmp alter column collection_share_id set not null;
                        ALTER TABLE collection_post_metrics_summary_tmp DROP COLUMN show_in_report;
                        ALTER TABLE collection_post_metrics_summary_tmp DROP COLUMN posted_by_cp_id;

                        BEGIN;
                            ALTER TABLE "collection_post_metrics_summary" RENAME TO "collection_post_metrics_summary_old_bkp";
                            ALTER TABLE "collection_post_metrics_summary_tmp" RENAME TO "collection_post_metrics_summary";
                        COMMIT;
                  """,
    )

    #truncate_fake_ts >> create_ts >> fetch_from_ch >> import_to_pg_box >> import_to_pg_db >> fill_gaps >> create_final_ts
    check_hour_task >> truncate_fake_ts >> create_ts >> fetch_from_ch >> import_to_pg_box >> import_to_pg_db
