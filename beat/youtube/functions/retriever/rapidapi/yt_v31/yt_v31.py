from loguru import logger
from core.models.models import RequestContext
from credentials.manager import CredentialManager
from utils.getter import safe_get
from utils.exceptions import QuotaExceeded, ProfileScrapingFailed, PostScrapingFailed
from utils.request import make_request
from youtube.functions.retriever.interface import YoutubeCrawlerInterface
from youtube.functions.retriever.rapidapi.yt_v31.yt_v31_parser import transform_profile, transform_post, \
    transform_playlist_for_post_ids, transform_channel_for_post_ids
from youtube.models.models import YoutubeProfileLog, YoutubePostLog

cred_provider = CredentialManager()


class YoutubeV31Api(YoutubeCrawlerInterface):

    @staticmethod
    async def fetch_profiles_by_channel_ids(ctx: RequestContext, channel_ids: list[str]) -> dict:
        headers = ctx.credential.credentials
        parts = "snippet,statistics,contentDetails"
        base_uri = "https://youtube-v31.p.rapidapi.com/channels"
        ids = ",".join(channel_ids)
        url = f"{base_uri}"
        querystring = {
            "part": parts,
            "id": ids
        }
        resp = await make_request("GET", url, headers=headers, params=querystring)

        if resp.status_code == 200 and "items" in resp.json():
            return resp.json()
        elif resp.status_code == 403:
            raise QuotaExceeded(resp.status_code, "Quota exceeded")
        else:
            raise ProfileScrapingFailed(resp.status_code, "Something went wrong")

    @staticmethod
    async def fetch_posts_by_post_ids(ctx: RequestContext, post_ids: list[str]) -> dict:
        headers = ctx.credential.credentials
        parts = "snippet,statistics,contentDetails"
        base_uri = "https://youtube-v31.p.rapidapi.com/videos"
        ids = ",".join(post_ids)
        url = f"{base_uri}"
        querystring = {
            "part": parts,
            "id": ids
        }
        resp = await make_request("GET", url, headers=headers, params=querystring)
        if resp.status_code == 200 and "items" in resp.json():
            return resp.json()
        elif resp.status_code == 403:
            raise QuotaExceeded(resp.status_code, "Quota exceeded")
        else:
            raise ProfileScrapingFailed(resp.status_code, "Something went wrong")

    @staticmethod
    async def fetch_post_ids_by_playlist_id(ctx: RequestContext, playlist_id: str) -> dict:
        headers = ctx.credential.credentials
        parts = "snippet,contentDetails,status,id"
        max_results = 50
        base_uri = "https://youtube-v31.p.rapidapi.com/playlistItems"
        url = f"{base_uri}"
        querystring = {
            "part": parts,
            "playlistId": playlist_id,
            "maxResults": max_results

        }
        resp = await make_request("GET", url, headers=headers, params=querystring)

        if resp.status_code == 200:
            return resp.json()
        elif resp.status_code == 403:
            raise QuotaExceeded(resp.status_code, "Quota exceeded")
        else:
            raise ProfileScrapingFailed(resp.status_code, "Something went wrong")

    @staticmethod
    async def fetch_post_ids_by_channel_id(ctx: RequestContext, channel_id: str, filter=None, cursor=None) -> (
    dict, bool, str):
        headers = ctx.credential.credentials
        base_uri = "https://youtube-v31.p.rapidapi.com/search"
        url = f"{base_uri}"
        querystring = {
            "channelId": channel_id,
            "part": "id",
            "order": "date",
            "maxResults": 50
        }
        if cursor:
            querystring["pageToken"] = cursor

        resp = await make_request("GET", url, headers=headers, params=querystring)

        if resp.status_code == 200:
            resp_dict = resp.json()
            next_cursor = safe_get(resp_dict, 'nextPageToken', type=str)
            if next_cursor:
                has_next_page = True
            else:
                has_next_page = False
            return resp_dict, has_next_page, next_cursor
        elif resp.status_code == 403:
            raise QuotaExceeded(resp.status_code, "Quota exceeded")
        else:
            raise ProfileScrapingFailed(resp.status_code, "Something went wrong")

    @staticmethod
    def parse_profiles_data(resp_dict: dict) -> list[YoutubeProfileLog]:
        return [transform_profile(item) for item in resp_dict["items"]]

    @staticmethod
    def parse_posts_data(resp_dict: dict, category=None, topic=None) -> list[YoutubePostLog]:
        return [transform_post(item, category=category, topic=topic) for item in resp_dict["items"]]

    @staticmethod
    def parse_post_ids_by_playlist_id(resp_dict: dict) -> list[str]:
        return [transform_playlist_for_post_ids(item) for item in resp_dict["items"]]

    @staticmethod
    def parse_post_ids_by_channel_id(resp_dict: dict, filter=None) -> list[str]:
        post_ids = []
        for item in resp_dict["items"]:
            try:
                post_ids.append(transform_channel_for_post_ids(item))
            except PostScrapingFailed as e:
                logger.error(e)
        return post_ids