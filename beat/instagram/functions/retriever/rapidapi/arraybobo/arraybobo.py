from typing import List

from loguru import logger
from typing import Tuple
from core.models.models import RequestContext
from instagram.functions.retriever.interface import InstagramCrawlerInterface
from instagram.functions.retriever.rapidapi.arraybobo.arraybobo_parser import transform_profile, transform_post, \
    transform_post_by_shortcode, transform_reels_post, transform_tagged_post, transform_story_post_data
from instagram.models.models import InstagramPostLog, InstagramProfileLog
from utils.exceptions import ProfileScrapingFailed, ProfileMissingError, NotSupportedError
from utils.request import make_request


class ArrayBobo(InstagramCrawlerInterface):

    @staticmethod
    async def fetch_profile_by_id(ctx: RequestContext, profile_id: str) -> dict:
        headers = ctx.credential.credentials
        url = "https://instagram-scraper-2022.p.rapidapi.com/ig/info/"
        querystring = {"id_user": profile_id}
        response = await make_request("GET", url, headers=headers, params=querystring)
        if response.status_code == 200:
            resp_dict = response.json()
            if 'user' in resp_dict and resp_dict['user']['pk'] == int(profile_id):
                return resp_dict
            logger.error(response)
            raise ProfileMissingError(response.status_code, "Invalid user id")
        if response.status_code == 403:
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "Forbidden exception")
        logger.error(response)
        raise ProfileScrapingFailed(response.status_code, "Something went wrong")

    @staticmethod
    async def fetch_profile_by_handle(ctx: RequestContext, handle: str) -> Tuple[dict, str]:
        headers = ctx.credential.credentials
        url = "https://instagram-scraper-2022.p.rapidapi.com/ig/info_username/"
        querystring = {"user": handle}
        response = await make_request("GET", url, headers=headers, params=querystring)
        if response.status_code == 200:
            resp_dict = response.json()
            if 'user' in resp_dict and resp_dict['user']['username'] == handle:
                profile_id = str(resp_dict['user']['pk'])
                return resp_dict, profile_id
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "API Failure")
        if response.status_code == 403:
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "Forbidden exception")
        logger.error(response)
        raise ProfileScrapingFailed(response.status_code, "Something went wrong")

    @staticmethod
    async def fetch_profile_posts_by_id(ctx: RequestContext, profile_id: str, cursor: str = None) -> Tuple[List[InstagramPostLog], bool, str]:
        headers = ctx.credential.credentials
        url = "https://instagram-scraper-2022.p.rapidapi.com/ig/posts/"
        querystring = {"id_user": profile_id}
        if cursor:
            querystring["end_cursor"] = cursor
        response = await make_request("GET", url, headers=headers, params=querystring)
        if response.status_code == 200:
            resp_dict = response.json()
            if 'data' in resp_dict and 'user' in resp_dict['data'] and resp_dict['data']['user']:
                page_info = resp_dict['data']['user']['edge_owner_to_timeline_media']['page_info']
                return resp_dict, page_info['has_next_page'], page_info['end_cursor']
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
        url = "https://instagram-scraper-2022.p.rapidapi.com/ig/post_info/"
        querystring = {"shortcode": shortcode}
        response = await make_request("GET", url, headers=headers, params=querystring)
        if response.status_code == 200:
            resp_dict = response.json()
            if 'shortcode' in resp_dict and resp_dict['shortcode'] == shortcode:
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
    async def fetch_reels_posts(ctx: RequestContext, profile_id: str, cursor: str = None) -> Tuple[List[InstagramPostLog], bool, str]:
        headers = ctx.credential.credentials
        url = "https://instagram-scraper-2022.p.rapidapi.com/ig/reels_posts/"
        querystring = {"id_user": profile_id}
        if cursor and cursor != "":
            querystring["max_id"] = cursor
        response = await make_request("GET", url, headers=headers, params=querystring)
        if response.status_code == 200:
            resp_dict = response.json()
            if 'items' in resp_dict and resp_dict['items']:
                max_id = None
                page_info = None
                if 'paging_info' in resp_dict:
                    page_info = resp_dict['paging_info']
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
    async def fetch_tagged_posts_by_profile_id(ctx: RequestContext, profile_id: str, cursor: str = None) -> Tuple[dict, bool, str]:
        headers = ctx.credential.credentials
        url = "https://instagram-scraper-2022.p.rapidapi.com/ig/tagged/"
        querystring = {"id_user": profile_id}
        if cursor:
            querystring["end_cursor"] = cursor
        response = await make_request("GET", url, headers=headers, params=querystring)
        if response.status_code == 200:
            resp_dict = response.json()
            if 'data' in resp_dict and 'user' in resp_dict['data'] and resp_dict['data']['user']:
                page_info = resp_dict['data']['user']['edge_user_to_photos_of_you']['page_info']
                return resp_dict, page_info['has_next_page'], page_info['end_cursor']
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
    async def fetch_profile_posts_by_handle(handle: str, cursor: str = None) -> Tuple[dict, bool, str]:
        raise NotSupportedError("Not Supported")

    @staticmethod
    def parse_profile_data(resp_dict: dict) -> InstagramProfileLog:
        return transform_profile(resp_dict)

    @staticmethod
    def parse_profile_post_data(resp_dict: dict) -> List[InstagramPostLog]:
        posts = [transform_post(p) for p in resp_dict['data']['user']['edge_owner_to_timeline_media']['edges']]
        return posts

    @staticmethod
    def parse_post_by_shortcode(resp_dict: dict) -> InstagramPostLog:
        return transform_post_by_shortcode(resp_dict)

    @staticmethod
    def parse_reels_post_data(resp_dict: dict) -> List[InstagramPostLog]:
        posts = [transform_reels_post(p) for p in resp_dict['items']]
        return posts
    
    @staticmethod
    def parse_tagged_posts_data(resp_dict: dict, profile_id: str) -> List[InstagramPostLog]:
        posts = [transform_tagged_post(p, profile_id) for p in resp_dict['data']['user']['edge_user_to_photos_of_you']['edges']]
        return posts

    @staticmethod
    async def fetch_story_posts_by_profile_id(ctx: RequestContext, profile_id: str, cursor: str = None) -> Tuple[dict]:
        headers = ctx.credential.credentials
        url = "https://instagram-scraper-2022.p.rapidapi.com/ig/stories/"
        querystring = {"id_user": profile_id}
        if cursor:
            querystring["end_cursor"] = cursor
        response = await make_request("GET", url, headers=headers, params=querystring)
        if response.status_code == 200:
            resp_dict = response.json()
            if 'reels' in resp_dict and resp_dict['reels'] and resp_dict['reels'][profile_id] is not None and 'user' in resp_dict['reels'][profile_id] and 'items' in resp_dict['reels'][profile_id]:
                return resp_dict
            logger.error(response)
            return []
        if response.status_code == 403:
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "Forbidden exception")
        if response.status_code == 403:
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "Forbidden exception")
        logger.error(response)
        raise ProfileScrapingFailed(response.status_code, "Something went wrong")
    
    @staticmethod
    def parse_story_posts_data(profile_id: str, resp_dict: dict) -> List[InstagramPostLog]:
        handle = resp_dict['reels'][profile_id]['user']['username']
        posts = [transform_story_post_data(p, handle) for p in resp_dict['reels'][profile_id]['items']]
        return posts