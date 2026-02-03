from youtube.functions.retriever.crawler import YoutubeCrawler
from youtube.models.models import YoutubeProfileLog, YoutubePostLog, YoutubePostActivityLog, YoutubeActivityLog, \
    YoutubeProfileRelationshipLog


async def parse_profiles_data(data: dict, source: str) -> list[YoutubeProfileLog]:
    crawler = YoutubeCrawler()
    data = crawler.parse_profiles_data(source, data)
    return data


async def parse_posts_data(data: dict, source: str, category=None, topic=None) -> list[YoutubePostLog]:
    crawler = YoutubeCrawler()
    data = crawler.parse_posts_data(source, data, category=category, topic=topic)
    return data


async def parse_profile_insights_data(data: dict, source: str) -> YoutubeProfileLog:
    crawler = YoutubeCrawler()
    data = crawler.parse_profile_insights_data(source, data)
    return data


async def parse_post_ids_by_playlist_id(data: dict, source: str) -> list[str]:
    crawler = YoutubeCrawler()
    data = crawler.parse_post_ids_by_playlist_id(source, data)
    return data


async def parse_post_ids_by_genre(data: dict, source: str) -> list[str]:
    crawler = YoutubeCrawler()
    data = crawler.parse_post_ids_by_genre(source, data)
    return data


async def parse_post_ids_by_channel_id(data: dict, source: str, filter=None) -> list[str]:
    crawler = YoutubeCrawler()
    data = crawler.parse_post_ids_by_channel_id(source, data, filter=filter)
    return data


async def parse_post_ids_by_search(data: dict, source: str) -> tuple[list[str], list[str], list[YoutubePostLog]]:
    crawler = YoutubeCrawler()
    data = crawler.parse_post_ids_by_search(source, data)
    return data


async def parse_post_comments(data: dict, profile_id: str, source: str) -> list[YoutubePostActivityLog]:
    crawler = YoutubeCrawler()
    insights = crawler.parse_post_comments(source, profile_id, data)
    return insights


async def parse_activities_by_channel_id(data: dict, source: str) -> list[YoutubeActivityLog]:
    crawler = YoutubeCrawler()
    data = crawler.parse_activities_by_channel_id(source, data)
    return data


async def parse_yt_profile_relationship_by_channel_id(data: dict, source: str) -> list[YoutubeProfileRelationshipLog]:
    crawler = YoutubeCrawler()
    data = crawler.parse_yt_profile_relationship_by_channel_id(source, data)
    return data
