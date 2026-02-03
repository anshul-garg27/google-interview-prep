from loguru import logger
from youtube.tasks.ingestion import parse_post_comments
from youtube.tasks.processing import upsert_activities
from youtube.tasks.retrieval import retrieve_post_comments


async def fetch_yt_post_comments(shortcode: str, max_comments: int, cursor=None, session=None) -> None:
    has_next_page = True
    total = 0
    max_count = 1000
    if max_comments:
        max_count = max_comments
    while has_next_page and total < max_count:
        logger.debug(shortcode)
        (data, has_next_page, cursor), source = await retrieve_post_comments(shortcode, cursor=cursor, session=session)
        logger.debug(f"Fetching post_comments - {has_next_page} {cursor}")
        logger.debug(data)
        logs = await parse_post_comments(data, shortcode, source)
        total += len(logs)
        await upsert_activities(logs, session=session)
