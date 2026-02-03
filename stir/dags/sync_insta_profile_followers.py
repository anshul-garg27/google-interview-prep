from airflow import DAG
from airflow.hooks.base import BaseHook
from airflow.operators.python import PythonOperator
from airflow.utils.dates import days_ago
from slack_connection import SlackNotifier


def create_scrape_request_logs():
    import os
    import json
    import requests
    import clickhouse_connect
    os.environ["no_proxy"] = "*"
    profile_ids = []
    connection = BaseHook.get_connection("clickhouse_gcc")
    extra_params = json.loads(connection.extra)
    client = clickhouse_connect.get_client(host=connection.host, password=connection.password, username=connection.login, send_receive_timeout=extra_params['send_receive_timeout'])
    # result = client.query('''with base as (
    #                                         select
    #                                             ig_id, max(handle) handle, max(followers) f, uniqExact(f.source_profile_id) followers_scraped, uniqExactIf(scl.id, scl.id > 0) requests
    #                                         from dbt.mart_instagram_account ia
    #                                         left join  beat_replica.profile_relationship_log f on f.target_profile_id = ia.ig_id
    #                                         left join beat_replica.scrape_request_log scl on scl.flow = 'fetch_profile_followers' and JSONExtractString(params, 'profile_id') = ia.ig_id and scl.created_at >= now() - INTERVAL 1 DAY
    #                                         where country = 'IN' and followers > 50000
    #                                         group by ig_id
    #                                         having followers_scraped = 0 and requests = 0
    #                                         order by f desc
    #                                     )
    #                                     select ig_id from base LIMIT 10
    #                             ''')

    ## Temporary change to prioritize handles for HUL
    # result = client.query('''with base as (
    #                                 select 
    #                                     platform_id ig_id, 
    #                                     max(handle) handle, 
    #                                     max(followers) f, 
    #                                     uniqExact(f.source_profile_id) followers_scraped, 
    #                                     uniqExactIf(scl.event_id, scl.event_id != '' and scl.event_id is not null) requests
    #                                 from dbt.mart_instagram_leaderboard ia
    #                                 left join  dbt.stg_beat_profile_relationship_log f on f.target_profile_id = ia.platform_id
    #                                 left join _e.scrape_request_log_events as scl on 
    #                                     scl.flow = 'fetch_profile_followers' and 
    #                                     JSONExtractString(params, 'profile_id') = ia.platform_id and 
    #                                     scl.event_timestamp >= now() - INTERVAL 1 DAY
    #                                 where 
    #                                     (month = toStartOfMonth(now())  or month = toStartOfMonth(date('2023-02-01'))) and 
    #                                     (   
    #                                         (followers_rank_by_cat_lang < 200 or views_rank_by_cat_lang < 200 or plays_rank_by_cat_lang < 200)
    #                                     )
    #                                 group by ig_id
    #                                 having followers_scraped = 0 and requests = 0
    #                                 order by f desc
    #                             )
    #                             select ig_id, handle from base LIMIT 10
    #                                 ''')
    
    result = client.query('''with creator_handles as (select *
                                                    from dbt.creator_handles_fake_follower_poc),
                                profiles as (select ig_id, handle
                                            from dbt.mart_instagram_account
                                            where handle in creator_handles
                                                and followers > 0),
                                scrapped_profile_with_followers as (select target_profile_id
                                                                    from dbt.mart_instagram_profiles_followers_count
                                                                    where followers_in_db > 500),
                                recent_profile_with_follower_event_data as (select log.target_profile_id as                                  target_profile_id,
                                                                                    uniqExact(JSONExtractString(source_dimensions, 'handle')) followers_in_db
                                                                            from _e.profile_relationship_log_events log
                                                                            where event_timestamp >= now() - INTERVAL 1 DAY
                                                                            and JSONExtractString(source_dimensions, 'handle') is not null
                                                                            and JSONExtractString(source_dimensions, 'full_name') is not null
                                                                            and JSONExtractString(source_dimensions, 'handle') != ''
                                                                            and JSONExtractString(source_dimensions, 'full_name') != ''
                                                                            group by target_profile_id),
                                recent_profile_with_followers as (select target_profile_id
                                                                from recent_profile_with_follower_event_data
                                                                where followers_in_db > 500),
                                recent_tried as (select JSONExtractString(params, 'profile_id') profile_id
                                                from _e.scrape_request_log_events
                                                where flow = 'fetch_profile_followers'
                                                    and event_timestamp >= now() - INTERVAL 1 day
                                                group by profile_id
                                                having uniqExact(scl_id) > 3)
                            select ig_id
                            from profiles
                            where ig_id not in scrapped_profile_with_followers
                            and ig_id not in recent_profile_with_followers
                            and ig_id not in recent_tried
                            limit 500
                          ''')

    for profile_id in result.result_rows:
        profile_ids.append(profile_id[0])

    flow = 'fetch_profile_followers'
    beat_connection = BaseHook.get_connection("beat")
    base_url = beat_connection.host
    url = base_url + "/scrape_request_log/flow/%s" % flow
    for profile_id in profile_ids:
        data = {"flow": flow, "platform": "INSTAGRAM", "params": {"profile_id": profile_id}}
        requests.post(url=url, data=json.dumps(data))


with DAG(
        dag_id='sync_insta_profile_followers',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval="0 */1 * * *",
) as dag:
    task = PythonOperator(
        task_id='create_insta_followers_scrape_request_log',
        python_callable=create_scrape_request_logs
    )
    task

