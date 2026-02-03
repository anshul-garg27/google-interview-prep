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
    # result = client.query('''with base as
    #                             (
    #                                 select ia.handle handle, ia.followers followers, sa.account_id account_id, sa.graph_access_token graph_access_token, sa.access_token_expired access_token_expired, up.phone phone, ua.created_on created_on
    #                                 from dbt.stg_winkl_instagram_accounts ia
    #                                 left join dbt.stg_identity_social_account sa on ia.handle = sa.handle 
    #                                 left join dbt.stg_identity_user_account ua on sa.account_id = ua.id
    #                                 left join dbt.stg_identity_user_phone up on ua.user_id = up.user_id
    #                                 where 
    #                                 ia.followers >= 500
    #                                 and sa.account_id is not null
    #                                 and up.is_primary = 1
    #                             ),
    #                             handles_raw as (select handle from base group by handle),
    #                             pids as (select profile_id from dbt.stg_beat_instagram_account where handle in handles_raw group by profile_id),
    #                             recent_attempts as (
    #                                 select JSONExtractString(params, 'profile_id') profile_id from 
    #                                 _e.scrape_request_log_events
    #                                 where flow = 'fetch_profile_following' and event_timestamp >= now() - INTERVAL 1 WEEK
    #                             )
    #                             select * from pids where profile_id not in recent_attempts LIMIT 300
    #                         ''')
    result = client.query('''with pids as 
                            (
                                select profile_id
                                from dbt.mart_trinity_graph 
                                where length(profiles_i_follow) = 0 and following > 0
                            ),
                            recent_attempts as (
                                select JSONExtractString(params, 'profile_id') profile_id from 
                                _e.scrape_request_log_events
                                where flow = 'fetch_profile_following' and event_timestamp >= now() - INTERVAL 3 DAY
                            )
                            select * from pids where profile_id not in recent_attempts LIMIT 300
                            ''')

    for profile_id in result.result_rows:
        profile_ids.append(profile_id[0])

    flow = 'fetch_profile_following'
    beat_connection = BaseHook.get_connection("beat")
    base_url = beat_connection.host
    url = base_url + "/scrape_request_log/flow/%s" % flow
    for profile_id in profile_ids:
        data = {"flow": flow, "platform": "INSTAGRAM", "params": {"profile_id": profile_id}}
        requests.post(url=url, data=json.dumps(data))


with DAG(
        dag_id='sync_insta_profile_following',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval="*/10 * * * *",
) as dag:
    task = PythonOperator(
        task_id='create_insta_followers_scrape_request_log',
        python_callable=create_scrape_request_logs
    )
    task

