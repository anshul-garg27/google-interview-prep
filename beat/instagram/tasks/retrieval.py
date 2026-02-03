from typing import Tuple
from core.helpers.session import sessionize
from core.models.models import PaginationContext
from instagram.functions.retriever.crawler import InstagramCrawler
from loguru import logger

@sessionize
async def retrieve_profile_data_by_handle(handle: str, session=None) -> Tuple[dict, str]:
    crawler = InstagramCrawler()
    (data, profile_id), source = await crawler.fetch_profile_by_handle(handle, session=session)
    return data, source


@sessionize
async def retrieve_profile_data_by_id(profile_id: str, session=None) -> Tuple[dict, str]:
    crawler = InstagramCrawler()
    data, source = await crawler.fetch_profile_by_id(profile_id, session=session)
    return data, source


@sessionize
async def retrieve_post_data_by_shortcode(shortcode: str, session=None) -> tuple[dict, str]:
    crawler = InstagramCrawler()
    data, source = await crawler.fetch_post_by_shortcode(shortcode, session=session)
    return data, source


@sessionize
async def retrieve_profile_posts_by_id(profile_id: str, cursor: str = None, session=None) -> Tuple[dict, bool, str, str]:
    crawler = InstagramCrawler()
    (data, has_next_page, cursor), source = await crawler.fetch_profile_posts_by_id(profile_id,
                                                                                    cursor,
                                                                                    session=session)
    return data, has_next_page, cursor, source


@sessionize
async def retrieve_reels_posts_by_id(profile_id: str, cursor: str = None, session=None) -> Tuple[dict, bool, str, str]:
    crawler = InstagramCrawler()
    (data, has_next_page, cursor), source = await crawler.fetch_reels_posts(profile_id, cursor, session=session)
    return data, has_next_page, cursor, source

@sessionize
async def retrieve_profile_data_from_graphapi(handle: str, session=None) -> dict:
    crawler = InstagramCrawler()
    data, _ = await crawler.fetch_profile_by_handle_from_graphapi(handle, session=session)
    return data


@sessionize
async def retrieve_post_insights(handle: str, post_id: str, post_type: str, session=None) -> Tuple[dict, str]:
    crawler = InstagramCrawler()
    data, source = await crawler.fetch_post_insights(handle, post_id, post_type, session=session)
    return data, source

@sessionize
async def retrieve_profile_insights(handle: str, session=None) -> Tuple[dict, str]:
    crawler = InstagramCrawler()
    data, source = await crawler.fetch_profile_insights(handle, session=session)
    return data, source

@sessionize
async def retrieve_profile_insights_from_graphapi(handle: str, session=None) -> dict:
    crawler = InstagramCrawler()
    data, _ = await crawler.fetch_profile_insights_from_graphapi(handle, session=session)
    return data
    
@sessionize
async def retrieve_story_insights(story_id: str, handle: str, session=None) -> Tuple[dict, str]:
    crawler = InstagramCrawler()
    data, source = await crawler.fetch_story_insights(story_id, handle, session=session)
    return data, source


@sessionize
async def retrieve_profile_followers(profile_id: str, pagination: PaginationContext = None, session=None) -> Tuple[dict, str]:
    crawler = InstagramCrawler()
    data, source = await crawler.fetch_followers(profile_id, pagination=pagination, session=session)
    return data, source


@sessionize
async def retrieve_profile_following(profile_id: str, cursor: str = None, session=None) -> Tuple[dict, str]:
    crawler = InstagramCrawler()
    data, source = await crawler.fetch_following(profile_id, cursor=cursor, session=session)
    return data, source


@sessionize
async def retrieve_post_likes(shortcode: str, cursor: str = None, session=None) -> Tuple[dict, str]:
    crawler = InstagramCrawler()
    data, source = await crawler.fetch_post_likes(shortcode, cursor=cursor, session=session)
    return data, source


@sessionize
async def retrieve_post_comments(shortcode: str, cursor=None, session=None) -> Tuple[dict, str]:
    crawler = InstagramCrawler()
    data, source = await crawler.fetch_post_comments(shortcode, cursor=cursor, session=session)
    return data, source


@sessionize
async def retrieve_stories_posts_data(handle: str, cursor: str = None, session=None) -> Tuple[dict, str]:
    crawler = InstagramCrawler()
    (data, has_next_page, cursor), source = await crawler.fetch_stories_posts(handle, cursor, session=session)
    return data, has_next_page, cursor, source


@sessionize
async def retrieve_hashtag_posts(hashtag: str, cursor=None, session=None) -> Tuple[dict, str]:
    crawler = InstagramCrawler()
    data, source = await crawler.fetch_hashtag_posts(hashtag, cursor=cursor, session=session)
    return data, source

@sessionize
async def retrieve_tagged_posts_by_profile_id(profile_id: str, cursor: str = None, session=None) -> Tuple[dict, bool, str, str]:
    crawler = InstagramCrawler()
    (data, has_next_page, cursor), source = await crawler.fetch_tagged_posts_by_profile_id(profile_id,
                                                                                    cursor,
                                                                                    session=session)
    return data, has_next_page, cursor, source

@sessionize
async def retrieve_story_posts_by_profile_id(profile_id: str, cursor: str = None, session=None) -> Tuple[dict]:
    crawler = InstagramCrawler()
    data, source = await crawler.fetch_story_posts_by_profile_id(profile_id,cursor,session=session)
    return data, source