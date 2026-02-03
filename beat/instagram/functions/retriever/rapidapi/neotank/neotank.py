from typing import List

from loguru import logger

from core.models.models import RequestContext
from credentials.manager import CredentialManager
from instagram.functions.retriever.interface import InstagramCrawlerInterface
from instagram.functions.retriever.rapidapi.neotank.neotank_parser import transform_profile, transform_post
from instagram.models.models import InstagramPostLog, InstagramProfileLog
from utils.exceptions import ProfileScrapingFailed, ProfileMissingError, NotSupportedError
from utils.request import make_request

cred_provider = CredentialManager()

POST_FETCH_SIZE = 80
FAILED_STATUS_CODES = [401, 422, 503]


class NeoTank(InstagramCrawlerInterface):

    @staticmethod
    async def fetch_profile_by_id(ctx: RequestContext, profile_id: str) -> dict:
        headers = ctx.credential.credentials
        url = "https://instagram130.p.rapidapi.com/username-by-id"
        querystring = {"userid": profile_id}
        response = await make_request("GET", url, headers=headers, params=querystring)
        if response.status_code == 200:
            resp_dict = response.json()
            if 'data' in resp_dict and resp_dict['data']:
                return resp_dict
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "API Failure")
        if response.status_code in FAILED_STATUS_CODES:
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "API Failure")
        if response.status_code == 403:
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "Forbidden exception")
        logger.error(response)
        raise ProfileScrapingFailed(response.status_code, "Something went wrong")

    @staticmethod
    async def fetch_profile_by_handle(ctx: RequestContext, handle: str) -> (dict, str):
        headers = ctx.credential.credentials
        url = "https://instagram130.p.rapidapi.com/account-info"
        querystring = {"username": handle}
        response = await make_request("GET", url, headers=headers, params=querystring)
        if response.status_code == 200:
            resp_dict = response.json()
            if 'username' in resp_dict and resp_dict['username'] == handle:
                profile_id = str(resp_dict['id'])
                return resp_dict, profile_id
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "API Failure")
        if response.status_code == 404:
            logger.error(response)
            raise ProfileMissingError(response.status_code, "Invalid user id")
        if response.status_code in FAILED_STATUS_CODES:
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "API Failure")
        if response.status_code == 403:
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "Forbidden exception")
        logger.error(response)
        raise ProfileScrapingFailed(response.status_code, "Something went wrong")

    @staticmethod
    async def fetch_profile_posts_by_id(ctx: RequestContext, profile_id: str, cursor: str = None) -> (List[InstagramPostLog], bool, str):
        headers = ctx.credential.credentials
        url = "https://instagram130.p.rapidapi.com/account-medias"
        querystring = {"userid": profile_id, "first": POST_FETCH_SIZE}
        if cursor:
            querystring["after"] = cursor
        response = await make_request("GET", url, headers=headers, params=querystring)
        if response.status_code == 200:
            resp_dict = response.json()
            if 'owner' in resp_dict and resp_dict['edges'][0]['node']['owner']['id'] == profile_id:
                page_info = resp_dict['page_info']
                return resp_dict, page_info['has_next_page'], page_info['end_cursor']
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "API Failure")
        if response.status_code in FAILED_STATUS_CODES:
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "API Failure")
        if response.status_code == 403:
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "Forbidden exception")
        logger.error(response)
        raise ProfileScrapingFailed(response.status_code, "Something went wrong")

    @staticmethod
    async def fetch_profile_posts_by_handle(ctx: RequestContext, handle: str, cursor: str = None) -> (dict, bool, str):
        raise NotSupportedError("Not Supported")

    @staticmethod
    async def fetch_post_by_shortcode(ctx: RequestContext, shortcode: str) -> dict:
        raise NotSupportedError("Not Supported")

    @staticmethod
    async def fetch_reels_posts(ctx: RequestContext, profile_id: str) -> (dict, bool, str):
        raise NotSupportedError("Not Supported")

    @staticmethod
    def parse_profile_data(resp_dict: dict) -> InstagramProfileLog:
        return transform_profile(resp_dict)

    @staticmethod
    def parse_profile_post_data(resp_dict: dict) -> list[InstagramPostLog]:
        posts = [transform_post(p) for p in resp_dict['edges']]
        return posts
