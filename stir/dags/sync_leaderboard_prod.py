import random
from airflow import DAG
from airflow.providers.postgres.operators.postgres import PostgresOperator
from airflow.providers.ssh.operators.ssh import SSHOperator
from airflow.utils.dates import days_ago
from airflow_clickhouse_plugin.operators.clickhouse import ClickHouseOperator
import string

from slack_connection import SlackNotifier

with DAG(
        dag_id='sync_leaderboard_prod',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval="15 20 * * *",
) as dag:
    fetch_leaderboard_from_ch = ClickHouseOperator(
        task_id='fetch_leaderboard_from_ch',
        database='vidooly',
        sql=(
            '''
                INSERT INTO FUNCTION s3('https://gcc-social-data.s3.ap-south-1.amazonaws.com/data-pipeline/tmp/leaderboard.json', 'AKIAXGXUCIER7YGHF4X3', 'LNm/3YjC6L2dIZWSNGTahS9HtQvdJjlcme59ZWF1', 'JSONEachRow')
                SELECT * from dbt.mart_leaderboard
                SETTINGS s3_truncate_on_insert=1, output_format_json_quote_64bit_integers=0;
            '''
        ),
        clickhouse_conn_id='clickhouse_gcc',
    )
    download_cmd = "aws s3 cp s3://gcc-social-data/data-pipeline/tmp/leaderboard.json /tmp/leaderboard.json"
    import_to_pg_box = SSHOperator(
        ssh_conn_id='ssh_prod_pg',
        cmd_timeout=1000,
        task_id='download_to_pg_local',
        command=download_cmd,
        dag=dag)
    uniq = ''.join(random.choices(string.ascii_letters + string.digits, k=10))
    import_to_pg_db = PostgresOperator(
        task_id="import_to_pg_db",
        runtime_parameters={'statement_timeout': '14400s'},
        postgres_conn_id="prod_pg",
        sql=f"""
                DROP TABLE IF EXISTS leaderboard_tmp;
                DROP TABLE IF EXISTS mart_leaderboard;
                DROP TABLE IF EXISTS leaderboard_old_bkp;

                CREATE TEMP TABLE mart_leaderboard(data jsonb);

                COPY mart_leaderboard from '/tmp/leaderboard.json';

                create UNLOGGED table leaderboard_tmp (
                    like leaderboard
                    including defaults
                );

                INSERT INTO leaderboard_tmp
                (
                    month,
                    language,
                    category,
                    platform,
                    profile_type,
                    handle,
                    followers,
                    engagement_rate,
                    avg_likes,
                    avg_comments,
                    avg_views,
                    avg_plays,
                    followers_rank,
                    followers_change_rank,
                    views_rank,
                    platform_id,
                    country,
                    views,
                    yt_views,
                    ia_views,
                    uploads,
                    prev_followers,
                    prev_views,
                    prev_plays,
                    followers_change,
                    likes,
                    comments,
                    plays,
                    engagement,
                    plays_rank,
                    profile_platform,
                    enabled,
                    last_month_ranks,
                    current_month_ranks
                )
                SELECT
                    (data->>'month')::timestamp AS month,
                    (data->>'language') AS language,
                    (data->>'category') AS category,
                    (data->>'platform') AS platform,
                    (data->>'profile_type') AS profile_type,
                    (data->>'handle') AS handle,
                    (data->>'followers')::int8 AS followers,
                    (data->>'engagement_rate')::float AS engagement_rate,
                    (data->>'avg_likes')::float AS avg_likes,
                    (data->>'avg_comments')::float AS avg_comments,
                    (data->>'avg_views')::float AS avg_views,
                    (data->>'avg_plays')::float AS avg_plays,
                    (data->>'followers_rank')::int8 AS followers_rank,
                    (data->>'followers_change_rank')::int8 AS followers_change_rank,
                    (data->>'views_rank')::int8 AS views_rank,
                    (data->>'platform_id') AS platform_id,
                    (data->>'country') AS country,
                    (data->>'views')::int8 AS views,
                    (data->>'yt_views')::int8 AS yt_views,
                    (data->>'ia_views')::int8 AS ia_views,
                    (data->>'uploads')::int8 AS uploads,
                    (data->>'prev_followers')::int8 AS prev_followers,
                    (data->>'prev_views')::int8 AS prev_views,
                    (data->>'prev_plays')::int8 AS prev_plays,
                    (data->>'followers_change')::int8 AS followers_change,
                    (data->>'likes')::int8 AS likes,
                    (data->>'comments')::int8 AS comments,
                    (data->>'plays')::int8 AS plays,
                    (data->>'engagement')::int8 AS engagement,
                    (data->>'plays_rank')::int8 AS plays_rank,
                    (data->>'profile_platform') AS profile_platform,
                    TRUE enabled,
                    (data->>'last_month_ranks')::jsonb AS last_month_ranks,
                    (data->>'current_month_ranks')::jsonb AS current_month_ranks
                from mart_leaderboard;

                update leaderboard_tmp lv
                set platform_profile_id  = ia.id, audience_gender = ia.audience_gender, audience_age = ia.audience_age, search_phrase = ia.search_phrase
                from instagram_account ia
                where ia.ig_id  = lv.platform_id and lv.platform_id is not null;

                update leaderboard_tmp lv
                set platform_profile_id  = ya.id, audience_gender = ya.audience_gender, audience_age = ya.audience_age, search_phrase = ya.search_phrase
                from youtube_account ya
                where ya.channel_id  = lv.platform_id and lv.platform_id is not null;
                ALTER TABLE leaderboard_tmp SET LOGGED;

                CREATE UNIQUE INDEX leaderboard_{uniq}_pkey1 ON public.leaderboard_tmp USING btree (id);
                CREATE INDEX leaderboard_{uniq}_month_platform_platform_profile_id_idx1 ON public.leaderboard_tmp USING btree (month, platform, platform_profile_id);
                CREATE INDEX leaderboard_{uniq}_month_platform_followers_change_rank_id_idx1 ON public.leaderboard_tmp USING btree (month, platform, followers_change_rank, id);
                CREATE INDEX leaderboard_{uniq}_month_platform_plays_rank_id_idx1 ON public.leaderboard_tmp USING btree (month, platform, plays_rank, id);
                CREATE INDEX leaderboard_{uniq}_month_platform_views_rank_id_idx1 ON public.leaderboard_tmp USING btree (month, platform, views_rank, id);
                CREATE INDEX leaderboard_{uniq}_month_platform_idx1 ON public.leaderboard_tmp USING btree (month, platform);
                
                BEGIN;
                ALTER TABLE "leaderboard" RENAME TO "leaderboard_old_bkp";
                ALTER TABLE "leaderboard_tmp" RENAME TO "leaderboard";
                COMMIT;
                """,
    )
    vacuum_to_pg_db = PostgresOperator(
        task_id="vacuum_to_pg_db",
        runtime_parameters={'statement_timeout': '3000s'},
        postgres_conn_id="prod_pg",
        sql=f"""
            vacuum analyze leaderboard;
        """,
        autocommit=True,
    )
    fetch_leaderboard_from_ch >> import_to_pg_box >> import_to_pg_db >> vacuum_to_pg_db
