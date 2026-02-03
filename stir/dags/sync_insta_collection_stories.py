from airflow import DAG
from airflow.hooks.base import BaseHook
from airflow.operators.python import PythonOperator
from airflow.utils.dates import days_ago
from airflow.models import Variable
from slack_connection import SlackNotifier


def create_scrape_request_log():
    import os
    import json
    import requests
    import clickhouse_connect

    os.environ["no_proxy"] = "*"
    profile_ids = []
    handles = []
    profile_to_story_ids_map = {}
    max_stories = 250

    connection = BaseHook.get_connection("clickhouse_gcc")
    extra_params = json.loads(connection.extra)
    client = clickhouse_connect.get_client(host=connection.host, password=connection.password, username=connection.login, send_receive_timeout=extra_params['send_receive_timeout'])

    posts = client.query('''
                            with stories_to_scrape as (select short_code shortcode, argMax(posted_by_handle, updated_at) as handle
                                                    from dbt.stg_coffee_post_collection_item
                                                    where platform = 'INSTAGRAM'
                                                        and post_type IN ('story')
                                                        and (match(short_code, '^[0-9]{5,}$') = 1 or short_code like 'story_dum_%')
                                                        and created_at >= now() - INTERVAL 1 DAY
                                                    group by short_code),
                                    stories_shortcode as (select shortcode
                                                                                            from stories_to_scrape),
                                                                        stories_existed as (select shortcode
                                                            from _e.post_log_events
                                                            where shortcode in stories_shortcode
                                                            and event_timestamp >= now() - INTERVAL 2 DAY
                                                            and JSONHas(metrics, 'profile_follower')
                                                            and JSONHas(metrics, 'profile_er')
                                                            group by shortcode
                                                            having ((JSONExtractRaw(argMax(metrics, insert_timestamp), 'profile_follower') != '0' or
                                                                    JSONExtractRaw(argMax(metrics, insert_timestamp), 'profile_follower') == 'null')
                                                                or (JSONExtractRaw(argMax(metrics, insert_timestamp), 'profile_er') != '0' or
                                                                    JSONExtractRaw(argMax(metrics, insert_timestamp), 'profile_er') == 'null'))),
                                    candiate_set as (select mia.ig_id as profile_id,
                                                            mia.handle as handle,
                                                            sts.shortcode as shortcode
                                                    from stories_to_scrape sts
                                                            inner join dbt.mart_instagram_account mia on mia.handle = sts.handle
                                                    where sts.shortcode not in stories_existed)
                                select profile_id, handle, groupArray(shortcode)
                                from candiate_set
                                group by profile_id, handle
                            '''
                            f"LIMIT {max_stories}")

    for row in posts.result_rows:
        profile_ids.append(row[0])
        handles.append(row[1])
        profile_to_story_ids_map[row[0]] = row[2]
    
    beat_connection = BaseHook.get_connection("beat")
    base_url = beat_connection.host
    for handle in handles:
        url = base_url + "/profiles/instagram/byhandle/%s" % handle
        params = {"full_refresh": "true", "force_refresh": "true"}
        try:
            response = requests.get(url=url, params=params)
            print(f"url = {url}, params = {params}")
            if response.status_code == 200:
                resp_dict = response.json()
                followers = None 
                er = None
                static_er = None
                reels_er = None
                profile_id = None
                if 'status' in resp_dict and 'type' in resp_dict['status'] and resp_dict['status']['type']=='SUCCESS':
                    if 'profile' in resp_dict and 'followers' in resp_dict['profile']:
                        followers = resp_dict['profile']['followers']
                        if 'post_metrics' in resp_dict['profile'] and 'all' in resp_dict['profile']['post_metrics'] and 'avg_engagement' in resp_dict['profile']['post_metrics']['all']:
                            er = resp_dict['profile']['post_metrics']['all']['avg_engagement']
                        if 'post_metrics' in resp_dict['profile'] and 'static' in resp_dict['profile']['post_metrics'] and 'avg_engagement' in resp_dict['profile']['post_metrics']['static']:
                            static_er = resp_dict['profile']['post_metrics']['static']['avg_engagement']
                        if 'post_metrics' in resp_dict['profile'] and 'reels' in resp_dict['profile']['post_metrics'] and 'avg_engagement' in resp_dict['profile']['post_metrics']['reels']:
                            reels_er = resp_dict['profile']['post_metrics']['reels']['avg_engagement']
                    profile_id = resp_dict['profile']['profile_id']
                    
                    flow = 'refresh_story_posts_by_profile_id'
                    beat_connection = BaseHook.get_connection("beat")
                    base_url = beat_connection.host
                    url = base_url + "/scrape_request_log/flow/%s" % flow
                    data = {"flow": flow, "platform": "INSTAGRAM", "params": {"profile_id": profile_id, "handle": handle, "story_ids": profile_to_story_ids_map[profile_id], "profile_follower": followers, "profile_er" : er, "profile_static_er": static_er, "profile_reels_er": reels_er}}
                    requests.post(url=url, data=json.dumps(data))
        except Exception as e:
            print(f"error occured for handle = {handle}. error: {e}")

with DAG(
        dag_id='sync_insta_collection_stories',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval="0 */3 * * *",
) as dag:
    task = PythonOperator(
        task_id='create_insta_story_scl',
        python_callable=create_scrape_request_log
    )
    task
