from typing import List

from core.models.models import RequestContext
from instagram.models.models import InstagramProfileLog, InstagramPostLog, \
    InstagramPostActivityLog, InstagramRelationshipLog


class InstagramCrawlerInterface:

    @staticmethod
    async def fetch_profile_by_id(ctx: RequestContext, profile_id: str) -> dict:
        pass

    @staticmethod
    async def fetch_profile_by_handle(ctx: RequestContext, handle: str) -> (dict, str):
        pass

    @staticmethod
    async def fetch_profile_posts_by_id(ctx: RequestContext, profile_id: str, cursor: str = None) -> (dict, bool, str):
        pass

    @staticmethod
    async def fetch_profile_posts_by_handle(ctx: RequestContext, handle: str, cursor: str = None) -> (dict, bool, str):
        pass

    @staticmethod
    async def fetch_post_by_shortcode(ctx: RequestContext, shortcode: str) -> dict:
        pass

    @staticmethod
    async def fetch_reels_posts(ctx: RequestContext, profile_id: str) -> (list[dict], bool, str):
        pass

    @staticmethod
    async def fetch_followers(ctx: RequestContext, profile_id: str, cursor: str = None) -> (list[dict], bool, str):
        pass

    @staticmethod
    async def fetch_following(ctx: RequestContext, profile_id: str, cursor: str = None) -> (list[dict], bool, str):
        pass

    @staticmethod
    async def fetch_post_likes(ctx: RequestContext, shortcode: str, cursor: str = None) -> (list[dict], bool, str):
        pass

    @staticmethod
    async def fetch_post_comments(ctx: RequestContext, shortcode: str, cursor: str = None) -> (list[dict], bool, str):
        pass

    @staticmethod
    async def fetch_post_insights(ctx: RequestContext, post_id: str, post_type: str) -> dict:
        pass
    
    @staticmethod
    async def fetch_profile_insights(cx: RequestContext, handle: str) -> dict:
        pass

    @staticmethod
    async def fetch_stories_posts(cx: RequestContext, handle: str, cursor: str = None) -> (dict, str):
        pass

    @staticmethod
    async def fetch_story_insights(cx: RequestContext, handle: str, story_id: str) -> dict:
        pass

    @staticmethod
    async def fetch_hashtag_posts(ctx: RequestContext, hashtag: str, cursor: str = None) -> (dict, bool, str):
        pass

    @staticmethod
    def parse_profile_data(resp_dict: dict) -> InstagramProfileLog:
        pass

    @staticmethod
    def parse_profile_post_data(resp_dict: dict) -> List[InstagramPostLog]:
        pass

    @staticmethod
    def parse_post_by_shortcode(resp_dict: dict) -> InstagramPostLog:
        pass

    @staticmethod
    def parse_post_insights_data(shortcode: str, resp_dict: dict) -> List[InstagramPostLog]:
        pass
    
    @staticmethod
    def parse_profile_insights_data(resp_dict: dict) -> InstagramProfileLog:
        pass

    @staticmethod
    def parse_followers(profile_id: str, resp_dict: dict) -> List[InstagramRelationshipLog]:
        pass

    @staticmethod
    def parse_following(profile_id: str, resp_dict: dict) -> List[InstagramRelationshipLog]:
        pass

    @staticmethod
    def parse_post_likes(shortcode: str, resp_dict: dict) -> List[InstagramPostActivityLog]:
        pass

    @staticmethod
    def parse_post_comments(shortcode: str, resp_dict: dict) -> List[InstagramPostActivityLog]:
        pass

    @staticmethod
    def parse_story_insights_data(story_id: str, resp_dict: dict) -> InstagramPostLog:
        pass

    @staticmethod
    def parse_hashtag_posts(resp_dict: dict) -> List[InstagramPostLog]:
        pass