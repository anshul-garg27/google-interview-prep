from typing import List

from loguru import logger

from core.models.models import PaginationContext, RequestContext
from instagram.functions.retriever.interface import InstagramCrawlerInterface
from instagram.functions.retriever.rapidapi.jotucker.jotucker_parser import transform_profile, transform_post, \
    transform_post_by_shortcode, transform_follower, transform_following, transform_post_likes, transform_post_comments, \
    transform_hashtag_post
from instagram.models.models import InstagramPostLog, InstagramProfileLog, \
    InstagramPostActivityLog, InstagramRelationshipLog
from utils.exceptions import ProfileScrapingFailed, ProfileMissingError, NotSupportedError
from utils.getter import safe_get
from utils.request import make_request


class JoTucker(InstagramCrawlerInterface):

    @staticmethod
    async def fetch_profile_by_id(ctx: RequestContext, profile_id: str) -> dict:
        headers = ctx.credential.credentials
        url = "https://instagram-scraper2.p.rapidapi.com/user_info_by_id"
        querystring = {"user_id": profile_id}
        response = await make_request("GET", url, headers=headers, params=querystring)
        if response.status_code == 200:
            resp_dict = response.json()
            if isinstance(resp_dict, dict) and 'user' in resp_dict and resp_dict['user']['pk'] == profile_id:
                return resp_dict
            logger.error(response)
            raise ProfileMissingError(response.status_code, "Invalid user id")
        if response.status_code == 403:
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "Forbidden exception")
        logger.error(response)
        raise ProfileScrapingFailed(response.status_code, "Something went wrong")

    @staticmethod
    async def fetch_profile_by_handle(ctx: RequestContext, handle: str) -> (dict, str):
        headers = ctx.credential.credentials
        url = "https://instagram-scraper2.p.rapidapi.com/user_info_by_id"
        querystring = {"user_name": handle}
        response = await make_request("GET", url, headers=headers, params=querystring)
        if response.status_code == 200:
            resp_dict = response.json()
            if isinstance(resp_dict, dict) and 'user' in resp_dict and resp_dict['user']['username'] == handle:
                profile_id = str(resp_dict['user']['pk'])
                return resp_dict, profile_id
            logger.error(response)
            raise ProfileMissingError(response.status_code, "Invalid username")
        if response.status_code == 403:
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "Forbidden exception")
        logger.error(response)
        raise ProfileScrapingFailed(response.status_code, "Something went wrong")

    @staticmethod
    async def fetch_profile_posts_by_id(ctx: RequestContext, profile_id: str, cursor: str = None) -> (dict, bool, str):
        headers = ctx.credential.credentials
        url = "https://instagram-scraper2.p.rapidapi.com/medias_v2"
        querystring = {"user_id": profile_id, "batch_size": "50"}
        next_cursor = None
        has_next_page = False
        if cursor:
            querystring["max_id"] = cursor
        response = await make_request("GET", url, headers=headers, params=querystring)
        if response.status_code == 200:
            resp_dict = response.json()
            if 'items' in resp_dict and resp_dict['items']:
                if 'next_max_id' in resp_dict:
                    next_cursor = resp_dict['next_max_id']
                if 'more_available' in resp_dict:
                    has_next_page = resp_dict['more_available']
                return resp_dict, has_next_page, next_cursor
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "API Failure")
        if response.status_code == 403:
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "Forbidden exception")
        logger.error(response)
        raise ProfileScrapingFailed(response.status_code, "Something went wrong")

    @staticmethod
    async def fetch_post_by_shortcode(ctx: RequestContext, shortcode: str) -> dict:
        headers = ctx.credential.credentials
        url = "https://instagram-scraper2.p.rapidapi.com/media_info_v2"
        querystring = {"short_code": shortcode}
        response = await make_request("GET", url, headers=headers, params=querystring)
        if response.status_code == 200:
            resp_dict = response.json()
            if 'items' in resp_dict and len(resp_dict['items']) > 0:
                return resp_dict
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "API Failure")
        if response.status_code == 403:
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "Forbidden exception")
        logger.error(response)
        raise ProfileScrapingFailed(response.status_code, "Something went wrong")

    @staticmethod
    async def fetch_following(ctx: RequestContext, profile_id: str, cursor: str = None) -> (dict, bool, str):
        headers = ctx.credential.credentials
        url = "https://instagram-scraper2.p.rapidapi.com/following"
        batch_size = "50"
        querystring = {"user_id": profile_id, "batch_size": batch_size}
        if cursor:
            querystring["next_cursor"] = cursor
        response = await make_request("GET", url, headers=headers, params=querystring)
        if response.status_code == 200:
            data = response.json()
            has_next_page = safe_get(data, 'data.user.edge_follow.page_info.has_next_page',
                                     type=bool)
            cursor = safe_get(data, 'data.user.edge_follow.page_info.end_cursor', type=str)
            if not cursor or cursor == "":
                has_next_page = False
            return data, has_next_page, cursor
        logger.error(f"Error getting following for profile - {profile_id} - {response}")
        raise ProfileScrapingFailed(response.status_code, "Something went wrong")

    @staticmethod
    async def fetch_followers(ctx: RequestContext, profile_id: str, pagination: PaginationContext = None) -> (dict, bool, PaginationContext):
        headers = ctx.credential.credentials
        url = "https://instagram-scraper2.p.rapidapi.com/followers"
        batch_size = "50"
        querystring = {"user_id": profile_id, "batch_size": batch_size}
        if pagination.cursor:
            querystring["next_cursor"] = pagination.cursor
        response = await make_request("GET", url, headers=headers, params=querystring)
        if response.status_code == 200:
            data = response.json()
            has_next_page = safe_get(data, 'data.user.edge_followed_by.page_info.has_next_page',
                                     type=bool)
            pagination.cursor = safe_get(data, 'data.user.edge_followed_by.page_info.end_cursor', type=str)
            return data, has_next_page, pagination
        logger.error(f"Error getting followers for profile - {profile_id} - {response}")
        raise ProfileScrapingFailed(response.status_code, "Something went wrong")

    @staticmethod
    async def fetch_post_likes(ctx: RequestContext, shortcode: str, cursor: str = None) -> (dict, bool, str):
        headers = ctx.credential.credentials
        url = "https://instagram-scraper2.p.rapidapi.com/media_likers"
        batch_size = "50"
        querystring = {"short_code": shortcode, "batch_size": batch_size}
        if cursor:
            querystring["next_cursor"] = cursor
        response = await make_request("GET", url, headers=headers, params=querystring)
        if response.status_code == 200:
            data = response.json()
            has_next_page = safe_get(data, 'data.shortcode_media.edge_liked_by.page_info.has_next_page', type=bool)
            cursor = safe_get(data, 'data.shortcode_media.edge_liked_by.page_info.end_cursor', type=str)
            return data, has_next_page, cursor
        logger.error(f"Error getting likes for post - {shortcode} - {response}")
        raise ProfileScrapingFailed(response.status_code, "Something went wrong")

    @staticmethod
    async def fetch_post_comments(ctx: RequestContext, shortcode: str, cursor: str = None) -> (dict, bool, str):
        headers = ctx.credential.credentials
        url = "https://instagram-scraper2.p.rapidapi.com/media_comments"
        batch_size = "50"
        querystring = {"short_code": shortcode, "batch_size": batch_size}
        if cursor:
            querystring["next_cursor"] = cursor
        response = await make_request("GET", url, headers=headers, params=querystring)
        if response.status_code == 200:
            data = response.json()
            has_next_page = safe_get(data, 'data.shortcode_media.edge_media_to_parent_comment.page_info.has_next_page', type=bool)
            cursor = safe_get(data, 'data.shortcode_media.edge_media_to_parent_comment.page_info.end_cursor', type=str)
            if not data['data']['shortcode_media']:
                raise ProfileScrapingFailed(response.status_code, "Post is missing")
            return data, has_next_page, cursor
        logger.error(f"Error getting comments for post - {shortcode} - {response}")
        raise ProfileScrapingFailed(response.status_code, "Something went wrong")

    @staticmethod
    async def fetch_hashtag_posts(ctx: RequestContext, hashtag: str, cursor: str = None) -> (dict, bool, str):
        headers = ctx.credential.credentials
        url = "https://instagram-scraper2.p.rapidapi.com/hash_tag_medias_v2"
        batch_size = "50"
        querystring = {"hash_tag": hashtag, "batch_size": batch_size}
        if cursor:
            querystring["next_cursor"] = cursor
        response = await make_request("GET", url, headers=headers, params=querystring)
        if response.status_code == 200:
            data = response.json()
            has_next_page = safe_get(data, 'data.hashtag.edge_hashtag_to_ranked_media.page_info.has_next_page',
                                     type=bool)
            cursor = safe_get(data, 'data.hashtag.edge_hashtag_to_ranked_media.page_info.end_cursor', type=str)
            return data, has_next_page, cursor
        logger.error(f"Error getting posts for hashtag - {hashtag} - {response}")
        raise ProfileScrapingFailed(response.status_code, "Something went wrong")
    
    @staticmethod
    async def fetch_reels_posts(ctx: RequestContext, profile_id: str) -> (dict, bool, str):
        raise NotSupportedError("Not Supported")

    @staticmethod
    async def fetch_profile_posts_by_handle(handle: str, cursor: str = None) -> (dict, bool, str):
        raise NotSupportedError("Not Supported")

    @staticmethod
    def parse_profile_data(resp_dict: dict) -> InstagramProfileLog:
        return transform_profile(resp_dict)

    @staticmethod
    def parse_profile_post_data(resp_dict: dict) -> List[InstagramPostLog]:
        posts = [transform_post(p) for p in resp_dict['items']]
        return posts

    @staticmethod
    def parse_post_by_shortcode(resp_dict: dict) -> InstagramPostLog:
        return transform_post_by_shortcode(resp_dict['items'][0])

    @staticmethod
    def parse_followers(profile_id: str, resp_dict: dict) -> List[InstagramRelationshipLog]:
        return [transform_follower(profile_id, follower['node']) for follower in resp_dict['data']['user']['edge_followed_by']['edges']]

    @staticmethod
    def parse_following(profile_id: str, resp_dict: dict) -> List[InstagramRelationshipLog]:
        return [transform_following(profile_id, follower['node']) for follower in resp_dict['data']['user']['edge_follow']['edges']]

    @staticmethod
    def parse_post_likes(shortcode: str, resp_dict: dict) -> List[InstagramPostActivityLog]:
        return [transform_post_likes(shortcode, post['node']) for post in resp_dict['data']['shortcode_media']['edge_liked_by']['edges']]

    @staticmethod
    def parse_post_comments(shortcode: str, resp_dict: dict) -> List[InstagramPostActivityLog]:
        return [transform_post_comments(shortcode, post['node']) for post in resp_dict['data']['shortcode_media']['edge_media_to_parent_comment']['edges']]

    @staticmethod
    def parse_hashtag_posts(resp_dict: dict) -> List[InstagramPostLog]:
        return [transform_hashtag_post(post['node']) for post in resp_dict['data']['hashtag']['edge_hashtag_to_ranked_media']['edges']]
