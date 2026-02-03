from typing import List

from loguru import logger

from core.models.models import RequestContext
from instagram.functions.retriever.interface import InstagramCrawlerInterface
from instagram.functions.retriever.lama.lama_parser import transform_post_by_shortcode, transform_profile, \
    transform_profile_post, transform_profiles_post, transform_follower, transform_following, transform_story_post_data, \
    transform_hashtag_post
from instagram.models.models import InstagramPostLog, InstagramProfileLog, InstagramRelationshipLog
from utils.exceptions import ProfileScrapingFailed, ProfileMissingError
from utils.getter import safe_get
from utils.request import make_request


class Lama(InstagramCrawlerInterface):

    @staticmethod
    async def fetch_profile_by_handle(ctx: RequestContext, handle: str) -> (dict, str):
        headers = ctx.credential.credentials
        url = "https://api.lamadava.com/a1/user"
        querystring = {"username": handle}
        response = await make_request("GET", url, headers=headers, params=querystring)
        if response.status_code == 200:
            resp_dict = response.json()
            if isinstance(resp_dict, dict) and 'graphql' in resp_dict and resp_dict['graphql']['user']['username'] == handle:
                profile_id = str(resp_dict['graphql']['user']['id'])
                return resp_dict, profile_id
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "Something went wrong")
        if response.status_code == 400:
            logger.error(response)
            raise ProfileMissingError(404, "Profile not found")
        if response.status_code == 403:
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "Forbidden exception")
        logger.error(response)
        raise ProfileScrapingFailed(response.status_code, "Something went wrong")

    @staticmethod
    async def fetch_profile_by_id(ctx: RequestContext, profile_id: str) -> (dict, str):
        headers = ctx.credential.credentials
        url = "https://api.lamadava.com/v1/user/by/id"
        querystring = {"id": profile_id}
        response = await make_request("GET", url, headers=headers, params=querystring)
        if response.status_code == 200:
            resp_dict = response.json()
            if isinstance(resp_dict, dict) and 'pk' in resp_dict and str(resp_dict['pk']) == profile_id:
                return resp_dict
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "Something went wrong")
        if response.status_code == 400:
            logger.error(response)
            raise ProfileMissingError(404, "Profile not found")
        if response.status_code == 403:
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "Forbidden exception")
        logger.error(response)
        raise ProfileScrapingFailed(response.status_code, "Something went wrong")

    @staticmethod
    async def fetch_profile_posts_by_id(ctx: RequestContext, profile_id: str, cursor: str = None) -> (dict, bool, str):
        headers = ctx.credential.credentials
        url = "https://api.lamadava.com/v1/user/medias/chunk"
        querystring = {"user_id": profile_id, "max_amount": 200}
        if cursor:
            querystring["end_cursor"] = cursor
        response = await make_request("GET", url, headers=headers, params=querystring)
        if response.status_code == 200:
            resp = response.json()
            next_cursor = resp[1]
            has_next_page = False
            if len(resp) != 2:
                logger.error(response)
                raise ProfileScrapingFailed(response.status_code, "Format Changed")
            if next_cursor and next_cursor != "":
                has_next_page = True
            return resp, has_next_page, next_cursor
        if response.status_code == 400:
            logger.error(response)
            raise ProfileMissingError(404, "Profile not found")
        if response.status_code == 403:
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "Forbidden exception")
        logger.error(response)
        raise ProfileScrapingFailed(response.status_code, "Something went wrong")

    @staticmethod
    async def fetch_post_by_shortcode(ctx: RequestContext, shortcode: str) -> dict:
        headers = ctx.credential.credentials
        url = "https://api.lamadava.com/a1/media"
        querystring = {"code": shortcode}
        response = await make_request("GET", url, headers=headers, params=querystring)
        if response.status_code == 200:
            resp_dict = response.json()
            if 'items' in resp_dict and resp_dict['items'] and resp_dict['items'][0]['code'] == shortcode:
                return resp_dict
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "API Failure")
        if response.status_code == 403:
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "Forbidden exception")
        if response.status_code == 404:
            logger.error(response)
            raise ProfileScrapingFailed(400, "Failed to fetch post")
        logger.error(response)
        raise ProfileScrapingFailed(response.status_code, "Something went wrong")


    @staticmethod
    async def fetch_followers(ctx: RequestContext, profile_id: str, cursor: str = None) -> (dict, bool, str):
        headers = ctx.credential.credentials
        url = "https://api.lamadava.com/gql/user/followers/chunk"
        batch_size = "200"
        querystring = {"user_id": profile_id, "amount": batch_size}
        if cursor:
            querystring["end_cursor"] = cursor
        response = await make_request("GET", url, headers=headers, params=querystring)
        if response.status_code == 200:
            data = response.json()
            has_next_page = False
            cursor = None
            if len(data) > 1 and data[1] != "":
                cursor = data[1]
                has_next_page = True
            return data[0], has_next_page, cursor
        logger.error(f"Error getting followers for profile - {profile_id} - {response}")
        raise ProfileScrapingFailed(response.status_code, "Something went wrong")

    @staticmethod
    async def fetch_following(ctx: RequestContext, profile_id: str, cursor: str = None) -> (dict, bool, str):
        headers = ctx.credential.credentials
        url = "https://api.lamadava.com/gql/user/following/chunk"
        batch_size = "200"
        querystring = {"user_id": profile_id, "amount": batch_size}
        if cursor:
            querystring["end_cursor"] = cursor
        response = await make_request("GET", url, headers=headers, params=querystring)
        if response.status_code == 200:
            data = response.json()
            has_next_page = False
            cursor = None
            if len(data) > 1 and data[1] != "":
                cursor = data[1]
                has_next_page = True
            return data[0], has_next_page, cursor
        logger.error(f"Error getting followers for profile - {profile_id} - {response}")
        raise ProfileScrapingFailed(response.status_code, "Something went wrong")
    
    @staticmethod
    async def fetch_story_posts_by_profile_id(ctx: RequestContext, profile_id: str, cursor: str = None) -> dict:
        headers = ctx.credential.credentials
        url = "https://api.lamadava.com/v1/user/stories"
        querystring = {"user_id": profile_id}
        response = await make_request("GET", url, headers=headers, params=querystring)
        if response.status_code == 200:
            data = response.json()
            return data
        logger.error(f"Error getting stories for profile - {profile_id} - {response}")
        raise ProfileScrapingFailed(response.status_code, "Something went wrong")

    @staticmethod
    def parse_followers(profile_id: str, resp_array: dict) -> List[InstagramRelationshipLog]:
        return [transform_follower(profile_id, follower) for follower in resp_array]

    @staticmethod
    def parse_following(profile_id: str, resp_array: dict) -> List[InstagramRelationshipLog]:
        return [transform_following(profile_id, following) for following in resp_array]

    @staticmethod
    def parse_post_by_shortcode(resp_dict: dict) -> InstagramPostLog:
        return transform_post_by_shortcode(resp_dict['items'][0])

    @staticmethod
    def parse_profile_data(resp_dict: dict) -> InstagramProfileLog:
        return transform_profile(resp_dict)

    @staticmethod
    def parse_profile_post_data(resp_dict: dict) -> List[InstagramPostLog]:
        posts = [transform_profile_post(p) for p in resp_dict[0]]
        return posts

    @staticmethod
    def parse_profiles_post_data(resp_dict: dict) -> List[InstagramPostLog]:
        posts = [transform_profiles_post(p) for p in resp_dict['edges']]
        return posts

    @staticmethod
    def parse_story_posts_data(profile_id: str, resp_dict: dict) -> List[InstagramPostLog]:
        posts = [transform_story_post_data(p) for p in resp_dict]
        return posts

    @staticmethod
    async def fetch_hashtag_posts(ctx: RequestContext, hashtag: str, cursor: str = None) -> (dict, bool, str):
        headers = ctx.credential.credentials
        url = "https://api.lamadava.com/v1/hashtag/medias/recent/chunk"
        batch_size = "24"
        querystring = {"name": hashtag, "max_amount": batch_size}
        if cursor:
            querystring["max_id"] = cursor
        response = await make_request("GET", url, headers=headers, params=querystring)
        if response.status_code == 200:
            data = response.json()
            has_next_page = False
            cursor = None
            if len(data) > 1 and data[1] != "":
                cursor = data[1]
                has_next_page = True
            return data[0], has_next_page, cursor
        logger.error(f"Error getting posts for hashtag - {hashtag} - {response}")
        raise ProfileScrapingFailed(response.status_code, "Something went wrong")

    @staticmethod
    def parse_hashtag_posts(resp_list: dict) -> List[InstagramPostLog]:
        return [transform_hashtag_post(post) for post in resp_list]
