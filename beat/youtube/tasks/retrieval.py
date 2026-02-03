from core.helpers.session import sessionize
from youtube.functions.retriever.crawler import YoutubeCrawler


@sessionize
async def retrieve_profiles_data_by_channel_ids(channel_ids: list[str], session=None) -> (dict, str):
    crawler = YoutubeCrawler()
    data, source = await crawler.fetch_profiles_by_channel_ids(channel_ids, session=session)
    return data, source


@sessionize
async def retrieve_posts_data_by_post_ids(post_ids: list[str], session=None) -> (dict, str):
    crawler = YoutubeCrawler()
    data, source = await crawler.fetch_posts_by_post_ids(post_ids, session=session)
    return data, source


@sessionize
async def retrieve_profile_insights(channel_id: str, session=None) -> (dict, str):
    crawler = YoutubeCrawler()
    data, source = await crawler.fetch_profile_insights(channel_id, session=session)
    return data, source


@sessionize
async def retrieve_post_ids_by_playlist_id(playlist_id: str, cursor=None, session=None) -> (dict, bool, str, str):
    crawler = YoutubeCrawler()
    (data, has_next_page, cursor), source = await crawler.fetch_post_ids_by_playlist_id(playlist_id, cursor=cursor, session=session)
    return data, has_next_page, cursor, source


@sessionize
async def retrieve_post_ids_by_genre(category=None, language=None, cursor=None, session=None) -> (dict, bool, str, str):
    crawler = YoutubeCrawler()
    (data, has_next_page, cursor), source = await crawler.fetch_post_ids_by_genre(category, language, cursor=cursor,
                                                                                  session=session)
    return data, has_next_page, cursor, source


@sessionize
async def retrieve_post_ids_by_channel_id(channel_id: str, filter=None, cursor=None, session=None) -> (
dict, bool, str, str):
    crawler = YoutubeCrawler()
    (data, has_next_page, cursor), source = await crawler.fetch_post_ids_by_channel_id(channel_id, filter=filter,
                                                                                       cursor=cursor, session=session)
    return data, has_next_page, cursor, source


@sessionize
async def retrieve_posts_by_search(keyword: str, start_date: str, end_date: str, cursor=None, session=None) -> (
dict, bool, str, str):
    crawler = YoutubeCrawler()
    (data, has_next_page, cursor), source = await crawler.fetch_posts_by_search(keyword, start_date, end_date,
                                                                                cursor=cursor, session=session)
    return data, has_next_page, cursor, source


@sessionize
async def retrieve_post_type(post_id: str, session=None) -> str:
    crawler = YoutubeCrawler()
    post_type = await crawler.fetch_post_type(post_id, session=session)
    return post_type


@sessionize
async def retrieve_post_comments(shortcode: str, cursor=None, session=None) -> (dict, str):
    crawler = YoutubeCrawler()
    data, source = await crawler.fetch_post_comments(shortcode, cursor=cursor, session=session)
    return data, source


@sessionize
async def retrieve_activities_by_channel_id(channel_id: str, cursor=None, session=None) -> (dict, bool, str, str):
    crawler = YoutubeCrawler()
    (data, has_next_page, cursor), source = await crawler.fetch_activities_by_channel_id(channel_id, cursor=cursor,
                                                                                         session=session)
    return data, has_next_page, cursor, source


@sessionize
async def retrieve_yt_profile_relationship_by_channel_id(channel_id: str, cursor=None, session=None) -> (dict, bool, str, str):
    crawler = YoutubeCrawler()
    (data, has_next_page, cursor), source = await crawler.fetch_yt_profile_relationship_by_channel_id(channel_id, cursor=cursor,
                                                                                         session=session)
    return data, has_next_page, cursor, source
