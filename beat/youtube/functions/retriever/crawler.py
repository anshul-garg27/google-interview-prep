from core.crawler.crawler import Crawler
from core.helpers.session import sessionize
from core.models.models import RequestContext
from credentials.manager import CredentialManager
from utils.exceptions import NoAvailableSources
from youtube.functions.retriever.ytapi.ytapi import YoutubeDataApi
from youtube.functions.retriever.rapidapi.rapidapi_youtube.rapidapi_youtube import RapidApiYoutube
from youtube.functions.retriever.rapidapi.yt_v31.yt_v31 import YoutubeV31Api
from youtube.functions.retriever.rapidapi.rapidapi_youtube_search.rapidapi_youtube_search import RapidApiYoutubeSearch
from youtube.models.models import YoutubeProfileLog, YoutubePostLog, YoutubePostActivityLog, YoutubeActivityLog, \
    YoutubeProfileRelationshipLog

cred = CredentialManager()


class YoutubeCrawler(Crawler):

    def __init__(self) -> None:
        self.available_sources = {
            'fetch_profiles_by_ids': ['ytapitoken'],
            'fetch_posts_by_ids': ['ytapitoken'],
            'fetch_post_ids_by_playlist_id': ['ytapitoken'],
            'fetch_post_ids_by_genre': ['ytapitoken'],
            'fetch_post_ids_by_channel_id': ['ytapitoken', 'yt-v31', 'rapidapi-yt'],
            'fetch_post_ids_by_search': ['ytapitoken', 'rapidapi-yt-search'],
            'fetch_post_comments': ['ytapitoken', 'rapidapi-yt-search'],
            'fetch_activities_by_channel_id': ['ytapitoken'],
            'fetch_yt_profile_relationship_by_channel_id': ['ytapitoken'],
        }
        self.providers = {
            'ytapitoken': YoutubeDataApi,
            'ytapi': YoutubeDataApi,
            'rapidapi-yt': RapidApiYoutube,
            'rapidapi-yt-search': RapidApiYoutubeSearch,
            'yt-v31': YoutubeV31Api
        }

    @sessionize
    async def fetch_profiles_by_channel_ids(self, channel_ids: list[str], session=None) -> (dict, str):
        creds = await self.fetch_enabled_sources('fetch_profiles_by_ids', session=session)
        if not creds:
            raise NoAvailableSources("No source is available at the moment")
        provider = self.providers[creds.source]
        return await provider.fetch_profiles_by_channel_ids(RequestContext(creds), channel_ids), creds.source

    @sessionize
    async def fetch_posts_by_post_ids(self, post_ids: list[str], session=None) -> (dict, str):
        creds = await self.fetch_enabled_sources('fetch_posts_by_ids', session=session)
        if not creds:
            raise NoAvailableSources("No source is available at the moment")
        provider = self.providers[creds.source]
        return await provider.fetch_posts_by_post_ids(RequestContext(creds), post_ids), creds.source

    @sessionize
    async def fetch_profile_insights(self, channel_id: str, session=None) -> (dict, str):
        creds = await CredentialManager.get_enabled_cred_for_handle('ytapi', channel_id, session=session)
        if not creds:
            raise NoAvailableSources("No source is available at the moment")
        provider = self.providers[creds.source]
        return await provider.fetch_profile_insights(RequestContext(creds)), creds.source

    @sessionize
    async def fetch_post_ids_by_playlist_id(self, playlist_id: str, cursor=None, session=None) -> (dict, bool, str, str):
        creds = await self.fetch_enabled_sources('fetch_post_ids_by_playlist_id', session=session)
        if not creds:
            raise NoAvailableSources("No source is available at the moment")
        provider = self.providers[creds.source]
        return await provider.fetch_post_ids_by_playlist_id(RequestContext(creds), playlist_id, cursor), creds.source

    @sessionize
    async def fetch_post_ids_by_channel_id(self, channel_id: str, filter=None, cursor=None, session=None) -> (
    dict, bool, str, str):
        creds = await self.fetch_enabled_sources('fetch_post_ids_by_channel_id', session=session)
        if not creds:
            raise NoAvailableSources("No source is available at the moment")
        provider = self.providers[creds.source]
        return await provider.fetch_post_ids_by_channel_id(RequestContext(creds), channel_id, filter=filter,
                                                           cursor=cursor), creds.source

    @sessionize
    async def fetch_post_ids_by_genre(self, category=None, language=None, cursor=None, session=None) -> (
    dict, bool, str, str):
        creds = await self.fetch_enabled_sources('fetch_post_ids_by_genre', session=session)
        if not creds:
            raise NoAvailableSources("No source is available at the moment")
        provider = self.providers[creds.source]
        return await provider.fetch_post_ids_by_genre(RequestContext(creds), category, language, cursor), creds.source

    @sessionize
    async def fetch_posts_by_search(self, keyword: str, start_date: str, end_date: str, cursor=None, session=None) -> (
            dict, bool, str, str):
        creds = await self.fetch_enabled_sources('fetch_post_ids_by_search', session=session)
        if not creds:
            raise NoAvailableSources("No source is available at the moment")
        provider = self.providers[creds.source]
        return await provider.fetch_posts_by_search(RequestContext(creds), keyword, start_date, end_date,
                                                    cursor), creds.source

    @sessionize
    async def fetch_post_type(self, post_id: str, session=None) -> str:
        return await YoutubeDataApi.fetch_post_type(post_id)


    @sessionize
    async def fetch_post_comments(self, shortcode: str, cursor: str = None, session=None) -> (list[dict], bool, str):
        creds = await self.fetch_enabled_sources('fetch_post_comments', session=session)
        if not creds:
            raise NoAvailableSources("No source is available at the moment")
        provider = self.providers[creds.source]
        return await provider.fetch_post_comments(RequestContext(creds), shortcode, cursor), creds.source

    @sessionize
    async def fetch_activities_by_channel_id(self, channel_id: str, cursor: str = None, session=None) -> (dict, bool, str, str):
        creds = await self.fetch_enabled_sources('fetch_activities_by_channel_id', session=session)
        if not creds:
            raise NoAvailableSources("No source is available at the moment")
        provider = self.providers[creds.source]
        return await provider.fetch_activities_by_channel_id(RequestContext(creds), channel_id, cursor=cursor), creds.source

    @sessionize
    async def fetch_yt_profile_relationship_by_channel_id(self, channel_id: str, cursor: str = None, session=None) -> (dict, bool, str, str):
        creds = await self.fetch_enabled_sources('fetch_yt_profile_relationship_by_channel_id', session=session)
        if not creds:
            raise NoAvailableSources("No source is available at the moment")
        provider = self.providers[creds.source]
        return await provider.fetch_yt_profile_relationship_by_channel_id(RequestContext(creds), channel_id, cursor=cursor), creds.source

    def parse_profiles_data(self, source: str, resp_dict: dict) -> list[YoutubeProfileLog]:
        return self.providers[source].parse_profiles_data(resp_dict)

    def parse_posts_data(self, source: str, resp_dict: dict, category=None, topic=None) -> list[YoutubePostLog]:
        return self.providers[source].parse_posts_data(resp_dict, category=category, topic=topic)

    def parse_profile_insights_data(self, source: str, resp_dict: dict) -> YoutubeProfileLog:
        return self.providers[source].parse_profile_insights_data(resp_dict)

    def parse_post_ids_by_playlist_id(self, source: str, resp_dict: dict) -> list[str]:
        return self.providers[source].parse_post_ids_by_playlist_id(resp_dict)

    def parse_post_ids_by_genre(self, source: str, resp_dict: dict) -> list[str]:
        return self.providers[source].parse_post_ids_by_genre(resp_dict)

    def parse_post_ids_by_channel_id(self, source: str, resp_dict: dict, filter=None) -> list[str]:
        return self.providers[source].parse_post_ids_by_channel_id(resp_dict, filter=filter)

    def parse_post_ids_by_search(self, source: str, resp_dict: dict) -> tuple[
        list[str], list[str], list[YoutubePostLog]]:
        return self.providers[source].parse_post_ids_by_search(resp_dict)

    def parse_post_comments(self, source: str, shortcode: str, resp_dict: dict) -> list[YoutubePostActivityLog]:
        return self.providers[source].parse_post_comments(shortcode, resp_dict)

    def parse_activities_by_channel_id(self, source: str, resp_dict: dict) -> list[YoutubeActivityLog]:
        return self.providers[source].parse_activities_by_channel_id(resp_dict)

    def parse_yt_profile_relationship_by_channel_id(self, source: str, resp_dict: dict) -> list[YoutubeProfileRelationshipLog]:
        return self.providers[source].parse_yt_profile_relationship_by_channel_id(resp_dict)