from typing import List

import urllib3
from loguru import logger

from core.models.models import RequestContext
from instagram.functions.retriever.interface import InstagramCrawlerInterface
from instagram.functions.retriever.rapidapi.igapi.igapi_parser import transform_follower
from instagram.models.models import InstagramRelationshipLog
from utils.exceptions import ProfileScrapingFailed
from utils.request import make_request


class IGApi(InstagramCrawlerInterface):

    @staticmethod
    async def fetch_followers(ctx: RequestContext, profile_id: str, cursor: str = None) -> (dict, bool, str):
        headers = ctx.credential.credentials
        url = "https://igdata.p.rapidapi.com/web_followers"
        querystring = {"user_id": profile_id}
        if cursor:
            querystring["cursor"] = cursor
        headers['Accept-Encoding'] = ""
        response = await make_request("GET", url, headers=headers, params=querystring)
        logger.error(response.json())
        if response.status_code == 200:
            data = response.json()
            if 'edges' in data and 'page_info' in data and 'end_cursor' in data['page_info'] and 'has_next_page' in \
                    data['page_info']:
                page_info = data['page_info']
                cursor = page_info['end_cursor']
                has_next_page = page_info['has_next_page']
                # API return 200 with empty list but sometimes on retry returns data with the same cursor
                if data['count'] >= 10000 and (not has_next_page or not cursor or cursor == ""):
                    logger.error(f"Error getting followers for profile - {profile_id} - {response}")
                    raise ProfileScrapingFailed(response.status_code, "Something went wrong")
                return data, has_next_page, cursor
        logger.error(f"Error getting followers for profile - {profile_id} - {response}")
        raise ProfileScrapingFailed(response.status_code, "Something went wrong")

    @staticmethod
    def parse_followers(profile_id: str, resp_dict: dict) -> List[InstagramRelationshipLog]:
        return [transform_follower(profile_id, follower['node']) for follower in resp_dict['edges']]
