from core.models.models import RequestContext
from youtube.models.models import YoutubeProfileLog, YoutubePostLog, YoutubePostActivityLog, YoutubeActivityLog, \
    YoutubeProfileRelationshipLog


class YoutubeCrawlerInterface:

    @staticmethod
    async def fetch_profiles_by_channel_ids(ctx: RequestContext, channel_ids: list[str]) -> dict:
        pass

    @staticmethod
    async def fetch_posts_by_post_ids(ctx: RequestContext, post_ids: list[str]) -> dict:
        pass

    @staticmethod
    async def fetch_profile_insights(ctx: RequestContext) -> dict:
        pass

    @staticmethod
    async def fetch_post_ids_by_playlist_id(ctx: RequestContext, playlist_id: str, cursor: str = None, session=None) -> (dict, str, str):
        pass

    @staticmethod
    async def fetch_post_ids_by_genre(ctx: RequestContext, category: str, language: str, cursor=None, session=None) -> (
            dict, bool, str):
        pass

    @staticmethod
    async def fetch_post_comments(ctx: RequestContext, video_id: str, cursor: str = None) -> (list[dict], bool, str):
        pass

    @staticmethod
    async def fetch_activities_by_channel_id(ctx: RequestContext, channel_id: str, filter=None, cursor: str = None) -> (
    dict, bool, str):
        pass

    @staticmethod
    async def fetch_profile_relationship_by_channel_id(ctx: RequestContext, channel_id: str, filter=None, cursor: str = None) -> (
    dict, bool, str):
        pass

    @staticmethod
    def parse_profiles_data(resp_dict: dict) -> list[YoutubeProfileLog]:
        pass

    @staticmethod
    def parse_posts_data(resp_dict: dict) -> list[YoutubePostLog]:
        pass

    @staticmethod
    def parse_profile_insights_data(resp_dict: dict) -> YoutubeProfileLog:
        pass

    @staticmethod
    def parse_post_ids_by_playlist_id(source: str, resp_dict: dict) -> list[str]:
        pass

    @staticmethod
    def parse_post_ids_by_genre(source: str, resp_dict: dict) -> list[str]:
        pass

    @staticmethod
    def parse_post_comments(shortcode: str, resp_dict: dict) -> list[YoutubePostActivityLog]:
        pass

    @staticmethod
    def parse_activities_by_channel_id(resp_dict: dict) -> list[YoutubeActivityLog]:
        pass

    @staticmethod
    def parse_profile_relationship_by_channel_id(resp_dict: dict) -> list[YoutubeProfileRelationshipLog]:
        pass
