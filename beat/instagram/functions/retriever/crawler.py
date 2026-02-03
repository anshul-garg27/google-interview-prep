from typing import List, Tuple

from loguru import logger

from core.crawler.crawler import Crawler
from core.helpers.session import sessionize
from core.models.models import PaginationContext, RequestContext
from credentials.manager import CredentialManager
from instagram.functions.retriever.graphapi.graphapi import GraphApi
from instagram.functions.retriever.lama.lama import Lama
from instagram.functions.retriever.rapidapi.arraybobo.arraybobo import ArrayBobo
from instagram.functions.retriever.rapidapi.bestsolns.bestsolns import BestSolutions
from instagram.functions.retriever.rapidapi.igapi.igapi import IGApi
from instagram.functions.retriever.rapidapi.jotucker.jotucker import JoTucker
from instagram.functions.retriever.rapidapi.neotank.neotank import NeoTank
from instagram.functions.retriever.rocketapi.rocketapi import RocketApi
from instagram.models.models import InstagramProfileLog, InstagramPostLog, \
    InstagramPostActivityLog, InstagramRelationshipLog
from utils.exceptions import NoAvailableSources

cred = CredentialManager()


class InstagramCrawler(Crawler):

    def __init__(self) -> None:
        self.available_sources = {
            'fetch_profile_by_id': ['rapidapi-jotucker', 'rapidapi-arraybobo'],
            'fetch_profile_by_handle': ['graphapi',  'rapidapi-arraybobo', 'rapidapi-jotucker', 'lama', 'rapidapi-bestsolutions'],
            'fetch_profile_posts_by_id': ['rapidapi-jotucker', 'lama', 'rapidapi-arraybobo'],
            'fetch_profile_posts_by_handle': ['graphapi'],
            'fetch_reels_posts': ['rapidapi-arraybobo', 'rocketapi'],
            'fetch_post_by_shortcode': ['rapidapi-arraybobo', 'rapidapi-jotucker', 'rocketapi','lama'],
            'fetch_post_insights': ['graphapi'],
            'fetch_post_comments': ['rapidapi-jotucker'],
            'fetch_post_likes': ['rapidapi-jotucker'],
            'fetch_following': ['rapidapi-jotucker','lama'],
            'fetch_followers': ['rapidapi-jotucker', 'rocketapi'],
            'fetch_profile_insights': ['graphapi'],
            'fetch_hashtag_posts': ['rapidapi-jotucker', 'lama'],
            'fetch_tagged_posts_by_profile_id': ['rapidapi-arraybobo'],
            'fetch_story_posts_by_profile_id': ['rapidapi-arraybobo', 'rocketapi', 'lama']
        }
        self.providers = {
            'rapidapi-arraybobo': ArrayBobo,
            'rapidapi-neotank': NeoTank,
            'rapidapi-jotucker': JoTucker,
            'graphapi': GraphApi,
            'rapidapi-igapi': IGApi,
            'lama': Lama,
            'rapidapi-bestsolutions': BestSolutions,
            'rocketapi': RocketApi
        }
        self.fetch_follower_count = 0

    @sessionize
    async def fetch_profile_by_id(self, profile_id: str, session=None) -> Tuple[dict, str]:
        creds = await self.fetch_enabled_sources('fetch_profile_by_id', session=session)
        if not creds:
            raise NoAvailableSources("No source is available at the moment")
        provider = self.providers[creds.source]
        return await provider.fetch_profile_by_id(RequestContext(creds), profile_id), creds.source

    @sessionize
    async def fetch_profile_by_handle(self, handle: str, session=None) -> Tuple[dict, str]:
        creds = await self.fetch_enabled_sources('fetch_profile_by_handle', session=session)
        if not creds:
            raise NoAvailableSources("No source is available at the moment")
        provider = self.providers[creds.source]
        try:
            return await provider.fetch_profile_by_handle(RequestContext(creds), handle), creds.source
        except Exception as e:
            logger.error(f"Something went wrong - {e}")
            if creds.source == "graphapi":
                creds = await self.fetch_enabled_sources('fetch_profile_by_handle', exclusions=["graphapi"],
                                                         session=session)
                if not creds:
                    raise NoAvailableSources("No source is available at the moment")
                provider = self.providers[creds.source]
                return await provider.fetch_profile_by_handle(RequestContext(creds), handle), creds.source
            
    @sessionize
    async def fetch_profile_by_handle_from_graphapi(self, handle: str, session=None) -> dict:
        graphapi_cred = await CredentialManager.get_enabled_cred_for_handle('graphapi', handle, session=session)
        if not graphapi_cred:
            raise NoAvailableSources("No cred is available at the moment")
        return await GraphApi.fetch_profile_by_handle(RequestContext(graphapi_cred), handle)

    @sessionize
    async def fetch_profile_posts_by_id(self, profile_id: str, cursor: str = None, session=None) -> Tuple[dict, str]:
        creds = await self.fetch_enabled_sources('fetch_profile_posts_by_id', session=session)
        if not creds:
            raise NoAvailableSources("No source is available at the moment")
        provider = self.providers[creds.source]
        return await provider.fetch_profile_posts_by_id(RequestContext(creds), profile_id, cursor), creds.source

    @sessionize
    async def fetch_profile_posts_by_handle(self, handle: str, cursor: str = None, session=None) -> Tuple[dict, str]:
        creds = await self.fetch_enabled_sources('fetch_profile_posts_by_handle', session=session)
        if not creds:
            raise NoAvailableSources("No source is available at the moment")
        provider = self.providers[creds.source]
        return await provider.fetch_profile_posts_by_handle(RequestContext(creds), handle), creds.source

    @sessionize
    async def fetch_post_by_shortcode(self, shortcode: str, session=None) -> Tuple[dict, str]:
        creds = await self.fetch_enabled_sources('fetch_post_by_shortcode', session=session)
        if not creds:
            raise NoAvailableSources("No source is available at the moment")
        provider = self.providers[creds.source]
        return await provider.fetch_post_by_shortcode(RequestContext(creds), shortcode), creds.source
    

    @sessionize
    async def fetch_reels_posts(self, profile_id: str, cursor: str, session=None) -> Tuple[dict, bool, str]:
        creds = await self.fetch_enabled_sources('fetch_reels_posts', session=session)
        if not creds:
            raise NoAvailableSources("No source is available at the moment")
        provider = self.providers[creds.source]
        return await provider.fetch_reels_posts(RequestContext(creds), profile_id, cursor), creds.source

    @sessionize
    async def fetch_post_insights(self, handle: str, post_id: str, post_type: str, session=None) -> Tuple[dict, str]:
        creds = await CredentialManager.get_enabled_cred_for_handle('graphapi', handle, session=session)
        if not creds:
            raise NoAvailableSources("No cred is available at the moment")
        provider = self.providers[creds.source]
        return await provider.fetch_post_insights(RequestContext(creds), handle, post_id, post_type), creds.source

    @sessionize
    async def fetch_profile_insights(self, handle: str, session=None) -> Tuple[dict, str]:
        creds = await CredentialManager.get_enabled_cred_for_handle('graphapi', handle, session=session)
        if not creds:
            raise NoAvailableSources("No cred is available at the moment")
        provider = self.providers[creds.source]
        return await provider.fetch_profile_insights(RequestContext(creds), handle, False), creds.source
    
    @sessionize
    async def fetch_profile_insights_from_graphapi(self, handle: str, session=None) -> Tuple[dict, str]:
        graphApiCreds = await CredentialManager.get_enabled_cred_for_handle('graphapi', handle, session=session)
        if not graphApiCreds:
            raise NoAvailableSources("No cred is available at the moment")
        return await GraphApi.fetch_profile_insights(RequestContext(graphApiCreds), handle, True), graphApiCreds.source

    @sessionize
    async def fetch_story_insights(self, story_id: str, handle: str, session=None) -> Tuple[dict, str]:
        creds = await CredentialManager.get_enabled_cred_for_handle('graphapi', handle, session=session)
        if not creds:
            raise NoAvailableSources("No source is available at the moment")
        provider = self.providers[creds.source]
        return await provider.fetch_story_insights(RequestContext(creds), handle, story_id), creds.source
   
    @sessionize
    async def fetch_followers(self, profile_id: str, pagination: PaginationContext = None, session=None, fetch_follower_count=None) -> Tuple[list[dict], bool, PaginationContext]:
        creds = await self.fetch_enabled_sources('fetch_followers', session=session, pagination=pagination)
        if not creds:
            raise NoAvailableSources("No source is available at the moment")
        provider = self.providers[creds.source]
        pagination.source = creds.source
        return await provider.fetch_followers(RequestContext(creds), profile_id, pagination=pagination), creds.source

    @sessionize
    async def fetch_following(self, profile_id: str, cursor: str = None, session=None) -> Tuple[list[dict], bool, str]:
        creds = await self.fetch_enabled_sources('fetch_following', session=session)
        if not creds:
            raise NoAvailableSources("No source is available at the moment")
        provider = self.providers[creds.source]
        return await provider.fetch_following(RequestContext(creds), profile_id, cursor), creds.source

    @sessionize
    async def fetch_post_likes(self, shortcode: str, cursor: str = None, session=None) -> Tuple[list[dict], bool, str]:
        creds = await self.fetch_enabled_sources('fetch_post_likes', session=session)
        if not creds:
            raise NoAvailableSources("No source is available at the moment")
        provider = self.providers[creds.source]
        return await provider.fetch_post_likes(RequestContext(creds), shortcode, cursor), creds.source

    @sessionize
    async def fetch_post_comments(self, shortcode: str, cursor: str = None, session=None) -> Tuple[list[dict], bool, str]:
        creds = await self.fetch_enabled_sources('fetch_post_comments', session=session)
        if not creds:
            raise NoAvailableSources("No source is available at the moment")
        provider = self.providers[creds.source]
        return await provider.fetch_post_comments(RequestContext(creds), shortcode, cursor), creds.source

    @sessionize
    async def fetch_stories_posts(self, handle: str, cursor: str = None, session=None) -> Tuple[list[dict], bool, str]:
        creds = await CredentialManager.get_enabled_cred_for_handle('graphapi', handle, session=session)
        if not creds:
            raise NoAvailableSources("No source is available at the moment")
        provider = self.providers[creds.source]
        return await provider.fetch_stories_posts(RequestContext(creds), handle, cursor), creds.source

    @sessionize
    async def fetch_hashtag_posts(self, hashtag: str, cursor: str = None, session=None) -> Tuple[list[dict], bool, str]:
        creds = await self.fetch_enabled_sources('fetch_hashtag_posts', session=session)
        if not creds:
            raise NoAvailableSources("No source is available at the moment")
        provider = self.providers[creds.source]
        return await provider.fetch_hashtag_posts(RequestContext(creds), hashtag, cursor), creds.source
    
    @sessionize
    async def fetch_tagged_posts_by_profile_id(self, profile_id: str, cursor: str = None, session=None) -> Tuple[dict, bool, str, str]:
        creds = await self.fetch_enabled_sources('fetch_tagged_posts_by_profile_id', session=session)
        if not creds:
            raise NoAvailableSources("No source is available at the moment")
        provider = self.providers[creds.source]
        return await provider.fetch_tagged_posts_by_profile_id(RequestContext(creds), profile_id, cursor=cursor), creds.source

    def parse_profile_data(self, source: str, resp_dict: dict) -> InstagramProfileLog:
        return self.providers[source].parse_profile_data(resp_dict)

    def parse_profile_post_data(self, source: str, resp_dict: dict) -> List[InstagramPostLog]:
        return self.providers[source].parse_profile_post_data(resp_dict)

    def parse_profiles_post_data(self, source: str, resp_dict: dict) -> List[InstagramPostLog]:
        return self.providers[source].parse_profiles_post_data(resp_dict)

    def parse_post_by_shortcode(self, source: str, resp_dict: dict) -> InstagramPostLog:
        return self.providers[source].parse_post_by_shortcode(resp_dict)

    def parse_reels_post_data(self, source: str, resp_dict: dict) -> List[InstagramPostLog]:
        return self.providers[source].parse_reels_post_data(resp_dict)

    def parse_post_insights_data(self, source: str, shortcode: str, resp_dict: dict) -> InstagramPostLog:
        return self.providers[source].parse_post_insights_data(shortcode, resp_dict, )
    
    def parse_profile_insights_data(self, source: str, resp_dict: dict) -> InstagramProfileLog:
        return self.providers[source].parse_profile_insights_data(resp_dict)

    def parse_story_insights_data(self, story_id: str, source: str, resp_dict: dict) -> InstagramPostLog:
        return self.providers[source].parse_story_insights_data(story_id, resp_dict)

    def parse_followers(self, source: str, profile_id: str, resp_dict: dict) -> List[InstagramRelationshipLog]:
        return self.providers[source].parse_followers(profile_id, resp_dict)

    def parse_following(self, source: str, profile_id: str, resp_dict: dict) -> List[InstagramRelationshipLog]:
        return self.providers[source].parse_following(profile_id, resp_dict)

    def parse_post_likes(self, source: str, shortcode: str, resp_dict: dict) -> List[InstagramPostActivityLog]:
        return self.providers[source].parse_post_likes(shortcode, resp_dict)

    def parse_post_comments(self, source: str, shortcode: str, resp_dict: dict) -> List[InstagramPostActivityLog]:
        return self.providers[source].parse_post_comments(shortcode, resp_dict)

    def parse_hashtag_posts(self, source: str, resp_dict: dict) -> List[InstagramPostLog]:
        return self.providers[source].parse_hashtag_posts(resp_dict)
    
    def parse_tagged_posts_data(self, source: str, profile_id: str, resp_dict: dict) -> List[InstagramPostLog]:
        return self.providers[source].parse_tagged_posts_data(resp_dict, profile_id)
    
    @sessionize
    async def fetch_story_posts_by_profile_id(self, profile_id: str, cursor: str = None, session=None) -> Tuple[dict, bool, str, str]:
        creds = await self.fetch_enabled_sources('fetch_story_posts_by_profile_id', session=session)
        if not creds:
            raise NoAvailableSources("No source is available at the moment")
        provider = self.providers[creds.source]
        return await provider.fetch_story_posts_by_profile_id(RequestContext(creds), profile_id, cursor=cursor), creds.source
    
    def parse_story_posts_data(self, profile_id: str, source: str, resp_dict: dict) -> List[InstagramPostLog]:
        return self.providers[source].parse_story_posts_data(profile_id, resp_dict)
