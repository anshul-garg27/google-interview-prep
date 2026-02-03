from airflow import DAG
from airflow.providers.postgres.operators.postgres import PostgresOperator
from airflow.providers.ssh.operators.ssh import SSHOperator
from airflow.utils.dates import days_ago
from airflow_clickhouse_plugin.operators.clickhouse import ClickHouseOperator
from slack_connection import SlackNotifier

with DAG(
        dag_id='sync_group_metrics_prod',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval="0 0 * * *",
) as dag:
    fetch_from_ch = ClickHouseOperator(
        task_id='fetch_group_metrics_from_ch',
        database='vidooly',
        sql=(
            '''
                INSERT INTO FUNCTION s3('https://gcc-social-data.s3.ap-south-1.amazonaws.com/data-pipeline/tmp/group_metrics.json', 'AKIAXGXUCIER7YGHF4X3', 'LNm/3YjC6L2dIZWSNGTahS9HtQvdJjlcme59ZWF1', 'JSONEachRow')
                select
                    group_key,
                    group_avg_avg_likes group_avg_likes,
                    group_avg_avg_comments group_avg_comments,
                    group_avg_comments_rate,
                    group_avg_followers,
                    group_avg_engagement_rate,
                    group_avg_reels_reach,
                    group_avg_image_reach,
                    group_avg_story_reach,
                    group_avg_video_reach,
                    group_avg_shorts_reach,
                    group_avg_likes_to_comment_ratio,
                    group_avg_followers_growth7d,
                    group_avg_followers_growth30d,
                    group_avg_followers_growth90d,
                    group_avg_audience_reachability,
                    group_avg_audience_authencity,
                    group_avg_audience_quality,
                    group_avg_post_count,
                    group_avg_reaction_rate,
                    group_avg_likes_spread,
                    group_avg_followers_growth1y group_avg_followers_growth1y,
                    group_avg_avg_posts_per_week group_avg_posts_per_week,
                    profiles,
                    bin_start,
                    bin_end,
                    bin_height
                from dbt.mart_primary_group_metrics
                SETTINGS s3_truncate_on_insert=1;
            '''
        ),
        clickhouse_conn_id='clickhouse_gcc',
    )
    download_cmd = "aws s3 cp s3://gcc-social-data/data-pipeline/tmp/group_metrics.json /tmp/group_metrics.json"
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
                DROP TABLE IF EXISTS group_metrics_tmp;
                DROP TABLE IF EXISTS mart_group_metrics;
                DROP TABLE IF EXISTS group_metrics_old_bkp;
                
                CREATE TEMP TABLE mart_group_metrics(data jsonb);
                
                COPY mart_group_metrics from '/tmp/group_metrics.json';

                create table group_metrics_tmp (
                    like group_metrics
                    including defaults
                    including constraints
                    including indexes
                );
                
                INSERT INTO group_metrics_tmp
                (
                    group_key,
                    group_avg_likes,
                    group_avg_comments,
                    group_avg_comments_rate,
                    group_avg_followers,
                    group_avg_engagement_rate,
                    group_avg_reels_reach,
                    group_avg_image_reach,
                    group_avg_story_reach,
                    group_avg_shorts_reach,
                    group_avg_video_reach,
                    group_avg_likes_to_comment_ratio,
                    group_avg_followers_growth7d,
                    group_avg_followers_growth30d,
                    group_avg_followers_growth90d,
                    group_avg_audience_reachability,
                    group_avg_audience_authencity,
                    group_avg_audience_quality,
                    group_avg_post_count,
                    group_avg_reaction_rate,
                    group_avg_likes_spread,
                    group_avg_followers_growth1y,
                    group_avg_posts_per_week,
                    profiles,
                    bin_start,
                    bin_end,
                    bin_height
                )
                SELECT
                    data->>'group_key' group_key,
                    (data->>'group_avg_likes')::float 	group_avg_likes,
                    (data->>'group_avg_comments')::float 	group_avg_comments,
                    (data->>'group_avg_comments_rate')::float 	group_avg_comments_rate,
                    (data->>'group_avg_followers')::float 	group_avg_followers,
                    (data->>'group_avg_engagement_rate')::float 	group_avg_engagement_rate,
                    (data->>'group_avg_reels_reach')::float 	group_avg_reels_reach,
                    (data->>'group_avg_image_reach')::float 	group_avg_image_reach,
                    (data->>'group_avg_story_reach')::float 	group_avg_story_reach,
                    (data->>'group_avg_shorts_reach')::float 	group_avg_shorts_reach,
                    (data->>'group_avg_video_reach')::float 	group_avg_video_reach,
                    (data->>'group_avg_likes_to_comment_ratio')::float 	group_avg_likes_to_comment_ratio,
                    (data->>'group_avg_followers_growth7d')::float 	group_avg_followers_growth7d,
                    (data->>'group_avg_followers_growth30d')::float 	group_avg_followers_growth30d,
                    (data->>'group_avg_followers_growth90d')::float 	group_avg_followers_growth90d,
                    (data->>'group_avg_audience_reachability')::float 	group_avg_audience_reachability,
                    (data->>'group_avg_audience_authencity')::float 	group_avg_audience_authencity,
                    (data->>'group_avg_audience_quality')::float 	group_avg_audience_quality,
                    (data->>'group_avg_post_count')::float 	group_avg_post_count,
                    (data->>'group_avg_reaction_rate')::float 	group_avg_reaction_rate,
                    (data->>'group_avg_likes_spread')::float 	group_avg_likes_spread,
                    (data->>'group_avg_followers_growth1y')::float 	group_avg_followers_growth1y,
                    (data->>'group_avg_posts_per_week')::float 	group_avg_posts_per_week,
                    (data->>'profiles')::int 	profiles,
                    jsonb_array_to_text_array((data->>'bin_start')::jsonb) 	bin_start,
                    jsonb_array_to_text_array((data->>'bin_end')::jsonb) 	bin_end,
                    jsonb_array_to_text_array((data->>'bin_height')::jsonb) 	bin_height
                from mart_group_metrics;
                
                BEGIN;
                ALTER TABLE "group_metrics" RENAME TO "group_metrics_old_bkp";
                ALTER TABLE "group_metrics_tmp" RENAME TO "group_metrics";
                COMMIT;
              """,
    )
    fetch_from_ch >> import_to_pg_box >> import_to_pg_db

