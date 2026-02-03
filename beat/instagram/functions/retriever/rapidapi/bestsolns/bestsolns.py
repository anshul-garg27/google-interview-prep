from typing import List

from loguru import logger

from core.models.models import RequestContext
from instagram.functions.retriever.interface import InstagramCrawlerInterface
from instagram.functions.retriever.rapidapi.bestsolns.bestsolns_parser import transform_profile, transform_post
from instagram.models.models import InstagramPostLog, InstagramProfileLog
from utils.exceptions import ProfileScrapingFailed, ProfileMissingError
from utils.request import make_request


class BestSolutions(InstagramCrawlerInterface):

    @staticmethod
    async def fetch_profile_by_id(ctx: RequestContext, profile_id: str) -> dict:
        headers = ctx.credential.credentials
        url = "https://instagram-api-cheap-best-performance.p.rapidapi.com/profile"
        querystring = {"user_id": profile_id}
        response = await make_request("GET", url, headers=headers, params=querystring)
        if response.status_code == 200:
            resp_dict = response.json()
            if isinstance(resp_dict, dict) and 'pk' in resp_dict and str(resp_dict['pk']) == profile_id:
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
        url = "https://instagram-api-cheap-best-performance.p.rapidapi.com/profile"
        querystring = {"username": handle}
        response = await make_request("GET", url, headers=headers, params=querystring)
        if response.status_code == 200:
            resp_dict = response.json()
            if isinstance(resp_dict, dict) and 'username' in resp_dict and resp_dict['username'] == handle:
                profile_id = str(resp_dict['pk'])
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
        url = "https://instagram-api-cheap-best-performance.p.rapidapi.com/feed"
        querystring = {"user_id": profile_id}
        if cursor:
            querystring["next_max_id"] = cursor
        response = await make_request("GET", url, headers=headers, params=querystring)
        if response.status_code == 200:
            resp_dict = response.json()
            if 'items' in resp_dict and resp_dict['items']:
                more_available = False
                next_max_id = None
                if 'more_available' in resp_dict and 'next_max_id' in resp_dict:
                    if resp_dict['next_max_id'] and resp_dict['next_max_id'] != "" and resp_dict['more_available']:
                        next_max_id = resp_dict['next_max_id']
                        more_available = True
                return resp_dict, more_available, next_max_id
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "API Failure")
        if response.status_code == 403:
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "Forbidden exception")
        logger.error(response)
        raise ProfileScrapingFailed(response.status_code, "Something went wrong")

    @staticmethod
    def parse_profile_data(resp_dict: dict) -> InstagramProfileLog:
        return transform_profile(resp_dict)

    @staticmethod
    def parse_profile_post_data(resp_dict: dict) -> List[InstagramPostLog]:
        posts = [transform_post(p) for p in resp_dict['items']]
        return posts
