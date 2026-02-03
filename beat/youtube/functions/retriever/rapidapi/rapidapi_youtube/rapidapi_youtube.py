from loguru import logger
from core.models.models import RequestContext
from credentials.manager import CredentialManager
from utils.exceptions import QuotaExceeded, ProfileScrapingFailed
from utils.request import make_request
from utils.getter import safe_get
from youtube.functions.retriever.interface import YoutubeCrawlerInterface
from youtube.functions.retriever.rapidapi.rapidapi_youtube.rapidapi_youtube_parser import transform_profile, \
    transform_post, transform_post_ids_by_channel_id
from youtube.models.models import YoutubeProfileLog, YoutubePostLog

cred_provider = CredentialManager()

class RapidApiYoutube(YoutubeCrawlerInterface):
     

    @staticmethod
    async def fetch_profiles_by_channel_ids(ctx: RequestContext, channel_ids: list[str]) -> list:
        profiles = []
        headers = ctx.credential.credentials
        url = "https://youtube138.p.rapidapi.com/v1/channel/details/"

        for channel_id in channel_ids:
            querystring = {"id": channel_id}
            response = await make_request("GET", url, headers=headers, params=querystring)
            if response.status_code == 200:
                profiles.append(response.json())
            elif response.status_code == 429:
                raise QuotaExceeded(response.status_code, "Quota exceeded")
            else:
                raise ProfileScrapingFailed(response.status_code, "Something went wrong")
        return profiles
    

    @staticmethod
    async def fetch_posts_by_post_ids(ctx: RequestContext, post_ids: list[str]) -> list:
        posts = []
        headers = ctx.credential.credentials
        url = "https://youtube138.p.rapidapi.com/v2/video/details/"

        for post_id in post_ids:
            querystring = {"id": post_id}
            response = await make_request("GET", url, headers=headers, params=querystring)

            if response.status_code == 200:
                posts.append(response.json())
            elif response.status_code == 429:
                raise QuotaExceeded(response.status_code, "Quota exceeded")
            else:
                raise ProfileScrapingFailed(response.status_code, "Something went wrong")        
        return posts


    @staticmethod
    async def fetch_post_ids_by_channel_id(ctx: RequestContext, channel_id: str, filter=None, cursor=None) -> (dict, bool, str):
        headers = ctx.credential.credentials
        url = "https://youtube138.p.rapidapi.com/v1/channel/videos/"

        querystring = {"id": channel_id}
        if cursor:
            querystring["cursorNext"] = cursor

        if filter:
            querystring["filter"] = filter
        else:
            querystring["filter"] = "videos_latest"

        response = await make_request("GET", url, headers=headers, params=querystring)

        if response.status_code == 200:
            resp_dict = response.json()
            next_cursor = safe_get(resp_dict, 'cursorNext', type=str)
            if next_cursor:
                has_next_page = True
            else:
                has_next_page = False
            return resp_dict, has_next_page, next_cursor
        elif response.status_code == 429:
            raise QuotaExceeded(response.status_code, "Quota exceeded")
        else:
            raise ProfileScrapingFailed(response.status_code, "Something went wrong")


    @staticmethod
    def parse_profiles_data(resp: list) -> list[YoutubeProfileLog]:
        return [transform_profile(item) for item in resp]


    @staticmethod
    def parse_posts_data(resp: list, category=None, topic=None) -> list[YoutubePostLog]:
        return [transform_post(item, category=category, topic=topic) for item in resp]
    

    @staticmethod
    def parse_post_ids_by_channel_id(resp_dict: dict, filter=None) -> list[str]:
        if filter == "shorts_latest":
            post_type = "short"
        elif filter == "live_now":
            post_type = "live"
        else:
            post_type = "video"
        return [transform_post_ids_by_channel_id(content, post_type) for content in resp_dict["contents"]]