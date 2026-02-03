from asyncio import sleep

from loguru import logger
from core.models.models import PaginationContext

from instagram.tasks.ingestion import parse_followers, parse_following, parse_post_likes, parse_post_comments, \
    parse_hashtag_posts
from instagram.tasks.processing import upsert_relationships, upsert_activities, upsert_post
from instagram.tasks.retrieval import retrieve_profile_followers, retrieve_profile_following, retrieve_post_likes, \
    retrieve_post_comments, retrieve_hashtag_posts

async def fetch_profile_followers(profile_id: str, session=None) -> None:
    has_next_page = True
    total = 0
    retry_counter = 0
    pagination = PaginationContext(
        cursor=None,
        source=None
    )
    while has_next_page and total < 2000:
        try:
            (data, has_more_pages, next_pagination), source = await retrieve_profile_followers(profile_id, pagination=pagination,
                                                                                           session=session)
            logger.debug(f"Fetching profile_followers - {has_next_page} {next_pagination}")
            logs = await parse_followers(data, profile_id, source)
            total += len(logs)
            await upsert_relationships(logs, session=session)
            await sleep(5)
            pagination = next_pagination
            has_next_page = has_more_pages
            retry_counter = 0
        except Exception as e:
            if retry_counter > 5:
                logger.error("Max Retries Reached - ")
                raise e
            retry_counter += 1
            logger.error("Error - Retrying in 10s")
            await sleep(10)


async def fetch_profile_following(profile_id: str, cursor=None, session=None) -> None:
    has_next_page = True
    total = 0
    retry_counter = 0
    while has_next_page and total < 200:
        try:
            (data, has_more_pages, next_cursor), source = await retrieve_profile_following(profile_id, cursor=cursor,
                                                                                           session=session)
            logger.debug(f"Fetching profile_following - {has_next_page} {cursor}")
            logs = await parse_following(data, profile_id, source)
            total += len(logs)
            await upsert_relationships(logs, session=session)
            await sleep(10)
            cursor = next_cursor
            has_next_page = has_more_pages
            retry_counter = 0
        except Exception as e:
            logger.error(f"Error: {e}")
            if retry_counter > 1:
                logger.error("Max Retries Reached - ")
                raise e
            retry_counter += 1
            logger.error("Error - Retrying in 10s")
            await sleep(30)


async def fetch_post_likes(shortcode: str, max_likes=None, cursor=None, session=None) -> None:
    has_next_page = True
    total = 0
    max_count = 2000
    if max_likes:
        max_count = max_likes

    while has_next_page and total < max_count:
        (data, has_next_page, cursor), source = await retrieve_post_likes(shortcode, cursor=cursor, session=session)
        logger.debug(f"Fetching post_likes - {has_next_page} {cursor}")
        logs = await parse_post_likes(data, shortcode, source)
        total += len(logs)
        await upsert_activities(logs, session=session)


async def fetch_post_comments(shortcode: str, max_comments: int, cursor=None, session=None) -> None:
    has_next_page = True
    total = 0
    max_count = 1000
    if max_comments:
        max_count = max_comments
    while has_next_page and total < max_count:
        (data, has_next_page, cursor), source = await retrieve_post_comments(shortcode, cursor=cursor, session=session)
        logger.debug(f"Fetching post_comments - {has_next_page} {cursor}")
        logs = await parse_post_comments(data, shortcode, source)
        total += len(logs)
        await upsert_activities(logs, session=session)


async def fetch_hashtag_posts(hashtag: str, max_posts: int, cursor=None, session=None) -> None:
    has_next_page = True
    total = 0
    retry_counter = 0
    while has_next_page and total < max_posts:
        try:
            (data, has_more_pages, next_cursor), source = await retrieve_hashtag_posts(hashtag, cursor=cursor,
                                                                                       session=session)
            logger.debug(f"Fetching hashtag posts - {has_next_page} {cursor}")
            logs = await parse_hashtag_posts(data, source)
            total += len(logs)

            if logs:
                for log in logs:
                    await upsert_post(log.shortcode, log, session=session)

            await sleep(5)
            cursor = next_cursor
            has_next_page = has_more_pages
            retry_counter = 0
        except Exception as e:
            if retry_counter > 5:
                logger.error("Max Retries Reached - ")
                raise e
            retry_counter += 1
            logger.error("Error - Retrying in 10s")
            await sleep(10)
