from typing import Tuple, List

import loguru

from core.models.models import RequestContext
from credentials.manager import CredentialManager
from utils.exceptions import ProfileScrapingFailed
from utils.request import make_request
from utils.getter import safe_get
from youtube.functions.retriever.interface import YoutubeCrawlerInterface
from youtube.functions.retriever.rapidapi.rapidapi_youtube_search.rapidapi_youtube_search_parser import \
    transform_post_ids_by_search, transform_channel_ids_by_search, \
    transform_post_log_by_search, transform_post_comments
from youtube.models.models import YoutubePostLog, YoutubePostActivityLog

cred_provider = CredentialManager()


class RapidApiYoutubeSearch(YoutubeCrawlerInterface):

    @staticmethod
    async def fetch_post_ids_by_search(ctx: RequestContext, category: str, language=None, cursor=None) -> (
    dict, bool, str):
        headers = ctx.credential.credentials
        url = "https://yt-api.p.rapidapi.com/search"
        querystring = {
            "query": category,
            "geo": "IN",
            "type": "video",
            "sort_by": "views",
            "upload_date": "month"
        }
        if language:
            querystring["lang"] = language
        if cursor:
            querystring["token"] = cursor
        resp = await make_request("GET", url, headers=headers, params=querystring)

        if resp.status_code == 200:
            resp_dict = resp.json()
            next_cursor = safe_get(resp_dict, 'continuation', type=str)
            if next_cursor:
                has_next_page = True
            else:
                has_next_page = False
            return resp_dict, has_next_page, next_cursor
        else:
            raise ProfileScrapingFailed(resp.status_code, "Something went wrong")

    @staticmethod
    async def fetch_posts_by_search(ctx: RequestContext, keyword: str, start_date: str, end_date: str, cursor=None) -> (
            dict, bool, str):

        headers = ctx.credential.credentials
        url = "https://youtube-v311.p.rapidapi.com/search/"
        querystring = {
            "part": "snippet",
            "maxResults": "50",
            "order": "viewCount",
            "publishedAfter": start_date,
            "publishedBefore": end_date,
            "q": keyword,
            "regionCode": "IN",
            "type": "video"
        }
        if cursor:
            querystring["pageToken"] = cursor
        resp = await make_request("GET", url, headers=headers, params=querystring)
        if not resp.json():
            loguru.logger.debug("getting none in response")
            resp = await make_request("GET", url, headers=headers, params=querystring)
        if resp.status_code == 200:
            resp_dict = resp.json()
            next_cursor = safe_get(resp_dict, 'nextPageToken', type=str)
            loguru.logger.debug(next_cursor)
            if next_cursor:
                has_next_page = True
            else:
                has_next_page = False
            return resp_dict, has_next_page, next_cursor
        else:
            loguru.logger.debug(resp)
            raise ProfileScrapingFailed(resp.status_code, "Something went wrong")

    @staticmethod
    async def fetch_post_comments(ctx: RequestContext, video_id: str, cursor: str = None) -> (dict, bool, str):
        headers = ctx.credential.credentials
        url = "https://youtube-v311.p.rapidapi.com/commentThreads/"
        querystring = {
            "part": "snippet",
            "videoId": video_id,
            "maxResults": "50",
            "order": "time",
            "textFormat": "html"
        }

        if cursor:
            querystring["pageToken"] = cursor
        response = await make_request("GET", url, headers=headers, params=querystring)
        if response.status_code == 200:
            data = response.json()
            next_cursor = safe_get(data, 'nextPageToken', type=str)
            if next_cursor:
                has_next_page = True
            else:
                has_next_page = False
            return data, has_next_page, next_cursor
        loguru.logger.error(f"Error getting comments for post - {video_id} - {response.json()}")
        raise ProfileScrapingFailed(response.status_code, "Something went wrong")


    @staticmethod
    def parse_post_ids_by_search(resp_dict: dict) -> tuple[list[str], list[str], list[YoutubePostLog]]:
        resp_data = resp_dict["items"]
        return [transform_post_ids_by_search(data) for data in resp_data], [transform_channel_ids_by_search(data) for
                                                                            data in resp_data], [
            transform_post_log_by_search(data) for data in resp_data]

    @staticmethod
    def parse_post_comments(shortcode: str, resp_dict: dict) -> List[YoutubePostActivityLog]:
        resp_data = resp_dict["items"]
        return [transform_post_comments(shortcode, data) for data in resp_data]
