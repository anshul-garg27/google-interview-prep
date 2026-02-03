from typing import List, Tuple

from loguru import logger

from instagram.functions.retriever.crawler import InstagramCrawler
from instagram.metric_dim_store import RECENT_POSTS, RECENT_STORIES
from instagram.models.models import InstagramPostLog, InstagramProfileLog, \
    InstagramPostActivityLog, InstagramRelationshipLog


async def parse_profile_data(data: dict, source: str) -> Tuple[InstagramProfileLog, List[InstagramPostLog]]:
    crawler = InstagramCrawler()
    data = crawler.parse_profile_data(source, data)
    # logger.debug(data)
    recent_posts_data = None
    posts_dim = [d for d in data.dimensions if d.key == RECENT_POSTS]
    stories_dim = [d for d in data.dimensions if d.key == RECENT_STORIES]
    if posts_dim and len(posts_dim) > 0:
        try:
            recent_posts_data = crawler.parse_profile_post_data(source, posts_dim[0].value)
            if stories_dim[0].value and len(stories_dim) > 0:
                recent_posts_data += crawler.parse_profile_post_data(source, stories_dim[0].value)
        except Exception as e:
            logger.error("Unable to fetch posts from profile API", e)
    return data, recent_posts_data


async def parse_profile_post_data(profile_id: str, data: dict, source: str) -> List[InstagramPostLog]:
    crawler = InstagramCrawler()
    posts = crawler.parse_profile_post_data(source, data)
    return posts


async def parse_reels_post_data(profile_id: str, data: dict, source: str) -> List[InstagramPostLog]:
    crawler = InstagramCrawler()
    posts = crawler.parse_reels_post_data(source, data)
    return posts


async def parse_post_data(data: dict, source: str) -> InstagramPostLog:
    crawler = InstagramCrawler()
    post = crawler.parse_post_by_shortcode(source, data)
    return post


async def parse_post_insights(data: dict, shortcode: str, source: str) -> InstagramPostLog:
    crawler = InstagramCrawler()
    insights = crawler.parse_post_insights_data(source, shortcode, data)
    return insights


async def parse_profile_insights(data: dict, source: str) -> InstagramProfileLog:
    crawler = InstagramCrawler()
    insights = crawler.parse_profile_insights_data(source, data)
    return insights


async def parse_story_insights(story_id: str, data: dict, source: str) -> InstagramPostLog:
    crawler = InstagramCrawler()
    insights = crawler.parse_story_insights_data(story_id, source, data)
    return insights


async def parse_followers(data: dict, profile_id: str, source: str) -> List[InstagramRelationshipLog]:
    crawler = InstagramCrawler()
    insights = crawler.parse_followers(source, profile_id, data)
    return insights


async def parse_following(data: dict, profile_id: str, source: str) -> List[InstagramRelationshipLog]:
    crawler = InstagramCrawler()
    insights = crawler.parse_following(source, profile_id, data)
    return insights


async def parse_post_likes(data: dict, profile_id: str, source: str) -> List[InstagramPostActivityLog]:
    crawler = InstagramCrawler()
    insights = crawler.parse_post_likes(source, profile_id, data)
    return insights


async def parse_post_comments(data: dict, profile_id: str, source: str) -> List[InstagramPostActivityLog]:
    crawler = InstagramCrawler()
    insights = crawler.parse_post_comments(source, profile_id, data)
    return insights


async def parse_hashtag_posts(data: dict, source: str) -> List[InstagramPostLog]:
    crawler = InstagramCrawler()
    hashtag_posts = crawler.parse_hashtag_posts(source, data)
    return hashtag_posts


async def parse_tagged_posts_data(data: dict, profile_id: str, source: str) -> List[InstagramPostLog]:
    crawler = InstagramCrawler()
    posts = crawler.parse_tagged_posts_data(source, profile_id, data)
    return posts


async def parse_story_posts_data(profile_id: str, data: dict, source: str) -> List[InstagramPostLog]:
    crawler = InstagramCrawler()
    posts = crawler.parse_story_posts_data(profile_id, source, data)
    return posts
