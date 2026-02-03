import asyncio
from datetime import datetime
import time

from loguru import logger
import traceback
from keyword_collection.categorization import youtube_categorization
from keyword_collection.keyword_collection_logs import keyword_collection_logs
from utils.scraping_windows import fetch_scraping_windows
from core.entities.entities import PostLog
from core.enums import enums
from core.helpers.session import sessionize
from utils.exceptions import QuotaExceeded

from utils.request import make_scrape_log_event
from utils.getter import dimension
from youtube.entities.entities import YoutubeAccount
from youtube.metric_dim_store import PLAYLIST_ID
from youtube.models.models import YoutubeProfileLog, YoutubePostLog, YoutubeActivityLog, YoutubeProfileRelationshipLog
from youtube.tasks.ingestion import parse_profiles_data, parse_posts_data, parse_profile_insights_data, \
    parse_post_ids_by_playlist_id, parse_post_ids_by_genre, parse_post_ids_by_channel_id, parse_post_ids_by_search, \
    parse_activities_by_channel_id, parse_yt_profile_relationship_by_channel_id
from youtube.tasks.processing import upsert_profile, upsert_post, upsert_profiles, upsert_posts, \
    upsert_profile_insights, \
    upsert_post_type, upsert_yt_activities, upsert_yt_profile_relationship
from youtube.tasks.retrieval import retrieve_post_type, retrieve_profiles_data_by_channel_ids, \
    retrieve_posts_data_by_post_ids, \
    retrieve_profile_insights, retrieve_post_ids_by_playlist_id, retrieve_post_ids_by_genre, \
    retrieve_post_ids_by_channel_id, \
    retrieve_posts_by_search, retrieve_activities_by_channel_id, retrieve_yt_profile_relationship_by_channel_id
from utils.db import get


def _batchify(lst, max_len=50):
    result = []
    for i in range(0, len(lst), max_len):
        result.append(lst[i:i + max_len])
    return result


@sessionize
async def fetch_profiles(channel_ids: list[str], session=None) -> list[YoutubeProfileLog]:
    data, source = await retrieve_profiles_data_by_channel_ids(channel_ids, session=session)
    logs = await parse_profiles_data(data, source)
    return logs


@sessionize
async def fetch_posts(post_ids: list[str], category=None, topic=None, session=None) -> list[YoutubePostLog]:
    data, source = await retrieve_posts_data_by_post_ids(post_ids, session=session)
    logs = await parse_posts_data(data, source, category=category, topic=topic)
    return logs


@sessionize
async def fetch_profile_insights(channel_id: str, session=None) -> YoutubeProfileLog:
    data, source = await retrieve_profile_insights(channel_id, session=session)
    if data:
        insights_logs = await parse_profile_insights_data(data, source)
        insights_logs.channel_id = channel_id
        return insights_logs
    else:
        raise Exception("Unable to fetch data for profile")


@sessionize
async def fetch_post_ids_by_playlist_id(playlist_id: str, max_posts=None, session=None) -> list[str]:
    has_next_page = True
    result_post_ids = []
    cursor = None
    logger.debug(max_posts)
    if max_posts is None:
        condition = lambda: has_next_page
    else:
        max_post_ids = 50
        if max_posts:
            max_post_ids = max_posts
        condition = lambda: has_next_page and len(result_post_ids) < max_post_ids
    
    while condition():
        data, has_next_page, next_cursor, source = await retrieve_post_ids_by_playlist_id(playlist_id, cursor=cursor, session=session)
        post_ids = await parse_post_ids_by_playlist_id(data, source)
        result_post_ids = result_post_ids + post_ids
        cursor = next_cursor
    return result_post_ids


@sessionize
async def fetch_post_ids_by_genre(category=None, language=None, session=None) -> list[str]:
    cursor = None
    has_next_page = True

    result_post_ids = []
    max_post_ids = 100

    while has_next_page and len(result_post_ids) < max_post_ids:
        data, has_next_page, next_cursor, source = await retrieve_post_ids_by_genre(category, language, cursor=cursor,
                                                                                    session=session)
        post_ids = await parse_post_ids_by_genre(data, source)
        result_post_ids = result_post_ids + post_ids
        cursor = next_cursor

    return result_post_ids


@sessionize
async def fetch_post_ids_by_search(job_id: str, keyword: str, start_date: str, end_date: str, frequency: int = None, session=None) -> \
        tuple[list[str], list[str], list[YoutubePostLog]]:
    result_post_ids = []
    result_channel_ids = set()
    result_post_logs = []

    cursor = None
    scraping_windows = fetch_scraping_windows(start_date, end_date, frequency)
    keyword_collection_logs(scraping_windows, job_id)
    has_next_page = True
    next_date = False
    for dates in scraping_windows:
        if next_date:
            has_next_page = True
        while has_next_page:
            try:
                data, has_next_page, next_cursor, source = await retrieve_posts_by_search(keyword=keyword,
                                                                                        start_date=dates[0],
                                                                                        end_date=dates[1],
                                                                                        cursor=cursor,
                                                                                        session=session)
                if data["items"]:
                    post_ids, channel_ids, post_logs = await parse_post_ids_by_search(data, source)
                    result_post_ids = result_post_ids + post_ids
                    result_post_logs += post_logs
                    result_channel_ids.update(channel_ids)
                    cursor = next_cursor
                    if ~has_next_page:
                        next_date = True
                else:
                    keyword_collection_logs("item have none values", job_id)
                    if ~has_next_page:
                        next_date = True
            except Exception as e:
                logger.error("Failed to fetch keywords")
                logger.error(traceback.format_exc())
                logger.error(e)
                has_next_page = False
    result_channel_ids = list(result_channel_ids)
    return result_post_ids, result_channel_ids, result_post_logs


@sessionize
async def fetch_post_ids_by_channel_id(channel_id: str, filter=None, max_posts=None, session=None) -> list[str]:
    post_ids = []
    result_post_ids = []
    condition = ''
    has_next_page = True
    if filter is not None and filter['start_date'] is not None and filter['end_date'] is not None:
        condition = lambda: has_next_page
    else:
        max_post_ids = 50
        if max_posts:
            max_post_ids = max_posts
        condition = lambda: has_next_page and len(result_post_ids) < max_post_ids

    cursor = None

    while condition():
        try:
            data, has_next_page, next_cursor, source = await retrieve_post_ids_by_channel_id(channel_id, filter=filter,
                                                                                             cursor=cursor,
                                                                                             session=session)
            post_ids = await parse_post_ids_by_channel_id(data, source, filter=filter)
            result_post_ids = result_post_ids + post_ids
            cursor = next_cursor
        except QuotaExceeded:
            return result_post_ids
    return result_post_ids


@sessionize
async def fetch_post_type(post_id: str, session=None):
    post_type = await retrieve_post_type(post_id, session=session)
    return post_type


@sessionize
async def fetch_activities_by_channel_id(channel_id: str, max_activities=None, session=None) -> list[
    YoutubeActivityLog]:
    activities = []
    has_next_page = True
    cursor = None

    if max_activities is None:
        condition = lambda: has_next_page
    else:
        condition = lambda: has_next_page and len(activities) <= max_activities

    while condition():
        data, has_next_page, next_cursor, source = await retrieve_activities_by_channel_id(channel_id, cursor=cursor,
                                                                                           session=session)
        logs = await parse_activities_by_channel_id(data, source)
        activities = activities + logs
        cursor = next_cursor
        subscription_activity_exist = False
        for log in logs:
            if log.activity_type == "subscription":
                subscription_activity_exist = True
                break
        if not subscription_activity_exist:
            break
    return activities


@sessionize
async def fetch_yt_profile_relationship_by_channel_id(channel_id: str, max_relationship=None, session=None) -> list[
    YoutubeProfileRelationshipLog]:
    total_profile_relationship = []
    has_next_page = True
    cursor = None

    if max_relationship is None:
        condition = lambda: has_next_page
    else:
        condition = lambda: has_next_page and len(total_profile_relationship) <= max_relationship

    while condition():
        data, has_next_page, next_cursor, source = await retrieve_yt_profile_relationship_by_channel_id(channel_id, cursor=cursor,
                                                                                              session=session)
        profile_relationship_logs = await parse_yt_profile_relationship_by_channel_id(data, source)
        total_profile_relationship = total_profile_relationship + profile_relationship_logs
        cursor = next_cursor
    return total_profile_relationship


@sessionize
async def refresh_yt_profiles(channel_ids: list[str], session=None) -> None:
    batches = _batchify(channel_ids, max_len=40)
    for batch in batches:
        profile_logs = await fetch_profiles(batch, session=session)
        await process_profiles_data(profile_logs, session=session)


@sessionize
async def refresh_yt_posts(post_ids: list[str], playlist_id=None, category=None, topic=None, session=None) -> None:
    batches = _batchify(post_ids, max_len=40)
    for batch in batches:
        post_logs = await fetch_posts(batch, category=category, topic=topic, session=session)
        for postLog in post_logs:
            postLog.dimensions.append(dimension(playlist_id, PLAYLIST_ID))
        await process_posts_data(post_logs, session=session)
        

@sessionize
async def refresh_yt_profile_insights(channel_id: str, session=None) -> None:
    profile_insights_logs = await fetch_profile_insights(channel_id, session=session)
    await process_profile_insights_data(profile_insights_logs, session=session)


@sessionize
async def refresh_yt_posts_by_playlist_id(channel_id: str, max_posts=None, session=None) -> None:
    channel = await get(session, YoutubeAccount, channel_id=channel_id)
    playlist_id = channel.uploads_playlist_id
    post_ids = await fetch_post_ids_by_playlist_id(playlist_id,max_posts, session=session)
    await refresh_yt_posts(post_ids, session=session)


@sessionize
async def refresh_yt_posts_by_channel_id(channel_id: str, filter=None, max_posts=None, session=None) -> None:
    post_ids = await fetch_post_ids_by_channel_id(channel_id, filter=filter, max_posts=max_posts, session=session)
    await refresh_yt_posts(post_ids, category=None, session=session)


@sessionize
async def refresh_yt_posts_by_genre(category: str, language: str, session=None):
    post_ids = await fetch_post_ids_by_genre(category, language, session=session)
    post_ids_batches = [post_ids[x:x + 50] for x in range(0, len(post_ids), 50)]
    for batch in post_ids_batches:
        await refresh_yt_posts(batch, category=category, session=session)


@sessionize
async def refresh_yt_posts_by_search(job_id: str, keywords: list[str], start_date: str, end_date: str,
                                     categorization: bool = True, frequency: int =None, session=None) -> tuple[list[str], list[str]]:
    total_post_ids = []
    total_channel_ids = []
    for keyword in keywords:
        try:
            keyword_collection_logs(keyword, job_id)
            search_start = time.perf_counter()
            post_ids, channel_ids, post_logs = await fetch_post_ids_by_search(job_id, keyword, start_date, end_date, frequency,
                                                                            session=session)

            await process_posts_data(post_logs, session=session)

            search_end = time.perf_counter()

            message = f"total time taken for search for keyword {keyword} is {search_end - search_start}"
            keyword_collection_logs(message, job_id)

            post_start = time.perf_counter()
            batches = _batchify(post_ids, max_len=40)
            for batch in batches:
                try:
                    await refresh_yt_posts(batch, category=None, topic=keyword, session=session)
                except Exception as e:
                    logger.error(e)
                    continue            
            post_end = time.perf_counter()

            message = f"total time taken to fetch post for keyword {keyword} is {post_end - post_start}"
            keyword_collection_logs(message, job_id)

            await asyncio.sleep(15)
            profile_start = time.perf_counter()
            batches = _batchify(channel_ids, max_len=40)
            for batch in batches:
                try:
                    await refresh_yt_profiles(batch, session=session)
                except Exception as e:
                    logger.error(e)
                    continue
            
            profile_end = time.perf_counter()

            message = f"total time taken to fetch profile for keyword {keyword} is {profile_end - profile_start}"
            keyword_collection_logs(message, job_id)

            await asyncio.sleep(10)

            if categorization:
                limit = 100
                total_post = len(post_logs)
                total_post = total_post if total_post < 1000 else 1000
                total_iterations = (total_post + limit - 1) // limit
                keyword_collection_logs(f"total post :- {total_post}", job_id)
                await asyncio.sleep(10)

                cate_start = time.perf_counter()
                for iteration in range(total_iterations):
                    start_index = iteration * limit
                    end_index = min(start_index + limit, total_post)
                    tasks = []
                    for i in range(start_index, end_index):
                        if i < len(post_logs):
                            task = asyncio.create_task(youtube_categorization(post_logs[i]))
                            tasks.append(task)
                    if tasks:
                        tasks, _ = await asyncio.wait(tasks)
                    posts_logs = [task.result() for task in tasks]

                    for post_log in posts_logs:
                        now = datetime.now()
                        post = PostLog(
                            platform=enums.Platform.YOUTUBE.name,
                            profile_id=post_log.channel_id,
                            platform_post_id=post_log.shortcode,
                            metrics=[m.__dict__ for m in post_log.metrics],
                            dimensions=[d.__dict__ for d in post_log.dimensions],
                            source='categorization',
                            timestamp=now
                        )
                        await make_scrape_log_event("post_log", post)
                cate_end = time.perf_counter()

                message = f"total time taken to perform categorization for keyword {keyword} is {cate_end - cate_start}"
                keyword_collection_logs(message, job_id)

            total_post_ids += post_ids
            total_channel_ids += channel_ids
            await asyncio.sleep(10)
        except Exception as e:
            logger.debug(e)
            continue
    return total_post_ids, total_channel_ids


@sessionize
async def refresh_yt_post_type(post_id: str, session=None):
    post_type = await fetch_post_type(post_id, session=session)
    await process_post_type(post_id, post_type, session=session)


@sessionize
async def refresh_yt_activities(channel_id=None, max_activities=None, session=None) -> None:
    yt_activity_logs = await fetch_activities_by_channel_id(channel_id=channel_id, max_activities=max_activities,
                                                            session=session)
    await process_activities_data(yt_activity_logs, session=session)


@sessionize
async def refresh_yt_profile_relationship(channel_id=None, max_relationship=None, session=None) -> None:
    yt_profile_relationship_logs = await fetch_yt_profile_relationship_by_channel_id(channel_id=channel_id, max_relationship=max_relationship, session=session)
    
    await process_yt_profile_relationship_data(yt_profile_relationship_logs, session=session)


@sessionize
async def process_profiles_data(profile_logs: list[YoutubeProfileLog], session=None) -> None:
    await upsert_profiles(profile_logs, session=session)


@sessionize
async def process_posts_data(post_logs: list[YoutubePostLog], session=None) -> None:
    await upsert_posts(post_logs, session=session)


@sessionize
async def process_activities_data(activity_logs: list[YoutubeActivityLog], session=None) -> None:
    await upsert_yt_activities(activity_logs, session=session)


@sessionize
async def process_yt_profile_relationship_data(profile_relationship_logs: list[YoutubeProfileRelationshipLog], session=None) -> None:
    await upsert_yt_profile_relationship(profile_relationship_logs, session=session)


@sessionize
async def process_profile_data(profile_log: YoutubeProfileLog, session=None) -> None:
    await upsert_profile(profile_log, session=session)


@sessionize
async def process_post_data(post_log: YoutubePostLog, session=None) -> None:
    await upsert_post(post_log, session=session)


@sessionize
async def process_profile_insights_data(profile_insights_log: YoutubeProfileLog, session=None) -> None:
    await upsert_profile_insights(profile_insights_log, session=session)


@sessionize
async def process_post_type(post_id: str, post_type: str, session=None):
    await upsert_post_type(post_id, post_type, session=session)
