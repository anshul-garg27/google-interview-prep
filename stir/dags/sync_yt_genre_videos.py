from airflow import DAG
from airflow.hooks.base import BaseHook
from airflow.operators.python import PythonOperator
from airflow.utils.dates import days_ago
from slack_connection import SlackNotifier

categories = [
    "Entertainment",
    "Music",
    "Comedy",
    "Food & Drinks",
    "Beauty",
    "Kids & Animation",
    "News & Politics",
    "Fashion & Style",
    "Science & Technology",
    "Autos & Vehicles",
    "Education",
    "Sports",
    "Travel & Leisure",
    "Pets & Animals",
    "Gaming",
    "Films",
    "Animation",
    "Non profits",
    "Vloging",
    "DIY",
    "Reviews",
    "Arts & Craft",
    "Events",
    "Family & Parenting",
    "Health & Fitness",
    "Motivational",
    "Religious Content",
    "Finance",
    "Real Estate",
    "Book",
    "Electronics",
    "Adult",
    "Legal",
    "Government",
    "Movie & Shows",
    "Devotional",
    "Blogs and Travel",
    "Astrology",
    "Supernatural",
    "Agriculture & Allied Sectors",
    "Infotainment",
    "People & Culture",
    "Photography & Editing",
    "IT & ITES",
    "Defence",
    "Miscellaneous"
]

languages = [
    "",
    "en",
    "hi",
    "kn",
    "ta",
    "gu",
    "te",
    "bn",
    "ml",
    "mr",
    "pa",
    "or"
]

def create_scrape_request_log():
    import os
    import json
    import requests

    os.environ["no_proxy"] = "*"

    flow = 'refresh_yt_posts_by_search'
    beat_connection = BaseHook.get_connection("beat")
    base_url = beat_connection.host
    url = base_url + "/scrape_request_log/flow/%s" % flow

    for category in categories:
        for language in languages:
            data = {"flow": flow, "platform": "YOUTUBE", "params": {"category": category, "language": language}}
            requests.post(url=url, data=json.dumps(data))

with DAG(
        dag_id='sync_yt_genre_videos',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval="0 0 * * 0",
) as dag:
    task = PythonOperator(
        task_id='create_yt_posts_by_search_scrape_request_log',
        python_callable=create_scrape_request_log
    )
    task