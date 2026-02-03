from airflow import DAG
from airflow.hooks.base import BaseHook
from airflow.operators.python import PythonOperator
from airflow.utils.dates import days_ago
from airflow.models import Variable
from slack_connection import SlackNotifier


def batchify(lst, max_len=50):
    result = []
    for i in range(0, len(lst), max_len):
        result.append(lst[i:i + max_len])
    return result


def create_scrape_request_log():
    import os
    import json
    import requests
    import clickhouse_connect

    os.environ["no_proxy"] = "*"
    shortcodes = []
    max_posts = 250

    connection = BaseHook.get_connection("clickhouse_gcc")
    extra_params = json.loads(connection.extra)
    client = clickhouse_connect.get_client(host=connection.host, password=connection.password, username=connection.login, send_receive_timeout=extra_params['send_receive_timeout'])

    posts = client.query('''with 
                                posts_to_scrape as (
                                    select
                                        short_code shortcode
                                    from dbt.stg_coffee_post_collection_item scpci
                                    inner join dbt.stg_coffee_post_collection scpc on scpci.post_collection_id = scpc.id and scpc.partner_id = 11743
                                    where platform = 'YOUTUBE'
                                    and match(short_code, '^[0-9]*$') = 0
                                ),
                                candiate_set as (
                                    select shortcode 
                                    from posts_to_scrape 
                                    group by shortcode
                                ),
                                post_details as (
                                    select shortcode, max(published_at) published_at 
                                    from dbt.stg_beat_youtube_post p
                                    where shortcode in (select shortcode from candiate_set) and p.published_at <= today() - INTERVAL 4 DAY and p.comments>0
                                    group by shortcode
                                ),
                                attempts_raw as (
                                    SELECT
                                        JSONExtractArrayRaw(ifNull(params, ''), 'post_ids') as post_ids
                                    from 
                                        _e.scrape_request_log_events
                                    where flow = 'fetch_yt_post_comments'
                                    and event_timestamp >= now() - INTERVAL 12 HOUR
                                ),
                                attempts as (
                                    select trim(BOTH '"' FROM arrayJoin(post_ids)) post_id
                                    from attempts_raw
                                    group by post_id
                                ),
                                top_candidates as (
                                    select candiate_set.shortcode shortcode
                                    from candiate_set
                                    left join post_details on candiate_set.shortcode = post_details.shortcode
                                ),
                                failures_raw as (
                                    select 
                                        JSONExtractArrayRaw(ifNull(params, ''), 'post_ids') as post_ids
                                    from 
                                        dbt.stg_beat_scrape_request_log
                                    where 
                                        platform = 'YOUTUBE'
                                        and status = 'FAILED'
                                        and data not like '%429%'
                                        and flow = 'fetch_yt_post_comments'
                                        and scraped_at >= now() - INTERVAL 1 MONTH
                                    group by post_ids
                                    having count(*) > 3
                                ),
                                failures as (
                                    select trim(BOTH '"' FROM arrayJoin(post_ids)) post_id
                                    from failures_raw
                                    group by post_id
                                ),
                                comment_existed as (select shortcode, count(*) total_comment
                                                     from _e.post_activity_log_events
                                                     where shortcode in top_candidates
                                                     group by shortcode
                                                     having total_comment > 0)
                                select shortcode from top_candidates where shortcode not in failures and shortcode not in attempts
                                and shortcode not in (select shortcode from comment_existed)
                                '''
                            f"LIMIT {max_posts}")

    for shortcode in posts.result_rows:
        shortcodes.append(shortcode[0])
    flow = 'fetch_yt_post_comments'
    beat_connection = BaseHook.get_connection("beat")
    base_url = beat_connection.host
    max_comments = 1000
    url = base_url + "/scrape_request_log/flow/%s" % flow

    for shortcode in shortcodes:
        data = {"flow": flow, "platform": "YOUTUBE", "params": {"shortcode": shortcode, "max_comments": max_comments}}
        requests.post(url=url, data=json.dumps(data))


with DAG(
        dag_id='sync_yt_post_comments',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval="0 */3 * * *",
) as dag:
    task = PythonOperator(
        task_id='create_yt_post_comments_scl',
        python_callable=create_scrape_request_log
    )
    task
