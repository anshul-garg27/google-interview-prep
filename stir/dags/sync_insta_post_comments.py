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

    posts = client.query('''with posts_to_scrape as (select short_code shortcode
                                                     from dbt.stg_coffee_post_collection_item scpci
                                                              inner join dbt.stg_coffee_post_collection scpc
                                                                         on scpci.post_collection_id = scpc.id and scpc.partner_id = 11743
                                                     where platform = 'INSTAGRAM'
                                                       and post_type NOT IN ('story')
                                                       and match(short_code, '^[0-9]*$') = 0),
                                 recently_attempted as (SELECT JSONExtractString(params, 'shortcode') as shortcode
                                                        from _e.scrape_request_log_events as scrape
                                                        where flow = 'fetch_post_comments'
                                                          and scrape.event_timestamp >= now() - INTERVAL 12 HOUR),
                                 failures as (SELECT JSONExtractString(params, 'shortcode') as shortcode
                                              from dbt.stg_beat_scrape_request_log
                                              where flow = 'fetch_post_comments'
                                                and status = 'FAILED'
                                                and data NOT LIKE '(204%'
                                                and data NOT LIKE '(409%'
                                                and data NOT LIKE '(429%'
                                                and created_at >= now() - INTERVAL 1 MONTH
                                                and shortcode in posts_to_scrape
                                              group by shortcode
                                              having uniqExactIf(id, status = 'FAILED') > 1
                                                 and uniqExactIf(id, status = 'COMPLETE') = 0),
                                 post_details as (select p.shortcode, max(p.publish_time) publish_time
                                                  from dbt.stg_beat_instagram_post p
                                                  where p.shortcode in posts_to_scrape
                                                    and p.publish_time <= today() - INTERVAL 4 DAY
                                                    and p.comments > 0
                                                  group by shortcode),
                                 comment_existed as (select shortcode, count(*) total_comment
                                                     from _e.post_activity_log_events
                                                     where shortcode in posts_to_scrape
                                                     group by shortcode
                                                     having total_comment > 0),
                                 candiate_set as (select pts.shortcode   shortcode,
                                                         pd.publish_time publish_time
                                                  from posts_to_scrape pts
                                                           left join post_details pd on pd.shortcode = pts.shortcode
                                                  where pts.shortcode not in recently_attempted
                                                    and pts.shortcode not in failures
                                                    and pts.shortcode not in (select shortcode from comment_existed))
                            select *
                            from candiate_set
                            order by publish_time asc nulls first

                                '''
                            f"LIMIT {max_posts}")

    for shortcode in posts.result_rows:
        shortcodes.append(shortcode[0])
    flow = 'fetch_post_comments'
    beat_connection = BaseHook.get_connection("beat")
    base_url = beat_connection.host
    max_comments = 1000
    url = base_url + "/scrape_request_log/flow/%s" % flow
    for shortcode in shortcodes:
        data = {"flow": flow, "platform": "INSTAGRAM", "params": {"shortcode": shortcode, "max_comments": max_comments}}
        requests.post(url=url, data=json.dumps(data))


with DAG(
        dag_id='sync_insta_post_comments',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval="0 */3 * * *",
) as dag:
    task = PythonOperator(
        task_id='create_insta_post_comments_scl',
        python_callable=create_scrape_request_log
    )
    task
