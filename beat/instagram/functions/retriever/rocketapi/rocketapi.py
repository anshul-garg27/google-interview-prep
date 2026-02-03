import json
from typing import List

from loguru import logger
from typing import Tuple
from core.models.models import PaginationContext, RequestContext
from instagram.functions.retriever.interface import InstagramCrawlerInterface
from instagram.functions.retriever.rocketapi.rocketapi_parser import transform_post_by_shortcode, transform_reels_post, \
    transform_story_post_data, transform_follower
from instagram.models.models import InstagramPostLog, InstagramProfileLog, InstagramRelationshipLog
from utils.exceptions import ProfileScrapingFailed, ProfileMissingError, NotSupportedError
from utils.getter import safe_get
from utils.request import make_request


class RocketApi(InstagramCrawlerInterface):

    @staticmethod
    async def fetch_post_by_shortcode(ctx: RequestContext, shortcode: str) -> dict:
        headers = ctx.credential.credentials
        url = "https://v1.rocketapi.io/instagram/media/get_info_by_shortcode"
        querystring = {"shortcode": shortcode}
        response = await make_request("GET", url, headers=headers, data=json.dumps(querystring))
        if response.status_code == 200:
            resp_dict = response.json()
            if 'body' in resp_dict['response'] and 'items' in resp_dict['response']['body'] and \
                    resp_dict['response']['body']['items']:
                return resp_dict
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "API Failure")
        if response.status_code == 403:
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "Forbidden exception")
        if response.status_code == 403:
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "Forbidden exception")
        logger.error(response)
        raise ProfileScrapingFailed(response.status_code, "Something went wrong")

    @staticmethod
    async def fetch_reels_posts(ctx: RequestContext, profile_id: str, cursor: str = None) -> Tuple[
        List[InstagramPostLog], bool, str]:
        headers = ctx.credential.credentials
        url = "https://v1.rocketapi.io/instagram/user/get_clips"
        querystring = {"id": profile_id}

        if cursor and cursor != "":
            querystring["max_id"] = cursor
        response = await make_request("GET", url, headers=headers, data=json.dumps(querystring))

        if response.status_code == 200:
            resp_dict = response.json()
            if 'body' in resp_dict['response'] and 'items' in resp_dict['response']['body'] and \
                    resp_dict['response']['body']['items']:
                max_id = None
                page_info = None
                if 'paging_info' in resp_dict['response']['body']:
                    page_info = resp_dict['response']['body']['paging_info']
                if page_info and 'max_id' in page_info:
                    max_id = page_info['max_id']
                more_available = None
                if page_info and 'more_available' in page_info:
                    more_available = page_info['more_available']
                return resp_dict, more_available, max_id
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "API Failure")
        if response.status_code == 403:
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "Forbidden exception")
        logger.error(response)
        raise ProfileScrapingFailed(response.status_code, "Something went wrong")

    @staticmethod
    async def fetch_story_posts_by_profile_id(ctx: RequestContext, profile_id: str, cursor: str = None) -> Tuple[dict]:
        headers = ctx.credential.credentials
        url = "https://v1.rocketapi.io/instagram/user/get_stories"
        datastring = {
            "ids": [int(profile_id)]
        }
        if cursor:
            datastring["end_cursor"] = cursor
        response = await make_request("POST", url, headers=headers, data=json.dumps(datastring))
        if response.status_code == 200:
            resp_dict = response.json()
            if 'body' in resp_dict['response'] and 'reels' in resp_dict['response']['body'] and \
                    resp_dict['response']['body']['reels'] and resp_dict['response']['body']['reels'][
                profile_id] is not None and 'user' in resp_dict['response']['body']['reels'][profile_id] and 'items' in \
                    resp_dict['response']['body']['reels'][profile_id]:
                return resp_dict
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "No Stories for this user")
        if response.status_code == 403:
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "Forbidden exception")
        if response.status_code == 403:
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "Forbidden exception")
        logger.error(response)
        raise ProfileScrapingFailed(response.status_code, "Something went wrong")


    @staticmethod
    async def fetch_followers(ctx: RequestContext, profile_id: str, pagination: PaginationContext = None) -> (dict, bool, PaginationContext):
        headers = ctx.credential.credentials
        url = "https://v1.rocketapi.io/instagram/user/get_followers"
        batch_size = 100
        querystring = {"id": profile_id, "count": batch_size}
        if pagination.cursor:
            querystring["max_id"] = pagination.cursor
        response = await make_request("GET", url, headers=headers, data=json.dumps(querystring))
        if response.status_code == 200:
            data = response.json()
            has_next_page = safe_get(data, 'response.body.big_list', type=bool)
            pagination.cursor = safe_get(data, 'response.body.next_max_id', type=str)
            return data, has_next_page, pagination
        logger.error(f"Error getting followers for profile - {profile_id} - {response}")
        raise ProfileScrapingFailed(response.status_code, "Something went wrong")

    @staticmethod
    def parse_story_posts_data(profile_id: str, resp_dict: dict) -> List[InstagramPostLog]:
        handle = resp_dict['response']['body']['reels'][profile_id]['user']['username']
        posts = [transform_story_post_data(p, handle) for p in
                 resp_dict['response']['body']['reels'][profile_id]['items']]
        return posts

    @staticmethod
    def parse_post_by_shortcode(resp_dict: dict) -> InstagramPostLog:
        return transform_post_by_shortcode(resp_dict)

    @staticmethod
    def parse_reels_post_data(resp_dict: dict) -> List[InstagramPostLog]:
        posts = [transform_reels_post(p) for p in resp_dict['response']['body']['items']]
        return posts

    @staticmethod
    def parse_followers(profile_id: str, resp_dict: dict) -> List[InstagramRelationshipLog]:
        return [transform_follower(profile_id, follower) for follower in resp_dict['response']['body']['users']]
