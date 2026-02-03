import time
from asyncio import sleep
from typing import List, Optional
from datetime import datetime, timedelta
from instagram.entities.entities import InstagramAccount
from loguru import logger
from typing import List, Tuple

from core.helpers.session import sessionize
from instagram.metric_dim_store import MEDIA_COUNT, POST_ID, POST_TYPE, PROFILE_ER, PROFILE_FOLLOWER, PROFILE_ID, \
    PROFILE_REELS_ER, PROFILE_STATIC_ER, HANDLE
from instagram.entities.entities import InstagramAccount, InstagramPost
from instagram.models.models import InstagramPostLog, InstagramProfileLog
from instagram.tasks.ingestion import parse_profile_data, parse_profile_insights, parse_profile_post_data, \
    parse_post_data, parse_post_insights, parse_story_insights, parse_reels_post_data, parse_tagged_posts_data, \
    parse_story_posts_data
from instagram.tasks.processing import upsert_profile, upsert_post, upsert_post_insights, upsert_profile_insights, \
    upsert_stories_posts, upsert_story_insights
from instagram.tasks.retrieval import retrieve_profile_data_by_handle, retrieve_profile_data_by_id, \
    retrieve_profile_insights, \
    retrieve_profile_posts_by_id, retrieve_post_data_by_shortcode, retrieve_post_insights, \
    retrieve_profile_data_from_graphapi, retrieve_story_insights, retrieve_stories_posts_data, \
    retrieve_reels_posts_by_id, retrieve_tagged_posts_by_profile_id, retrieve_story_posts_by_profile_id, retrieve_profile_insights_from_graphapi
from utils.db import get, get_or_create
from utils.getter import dimension, metric

@sessionize
async def fetch_recent_posts(profile_id: str, days: int = 30, session=None) -> List[InstagramPostLog]:
    all_posts_data = []
    has_next_page = True
    cursor = None
    max_posts_to_fetch = 50
    while has_next_page:
        posts_data, has_next_page, cursor, source = await retrieve_profile_posts_by_id(profile_id,
                                                                                       cursor,
                                                                                       session=session)
        posts_data_cleaned = await parse_profile_post_data(profile_id, posts_data, source)
        all_posts_data += posts_data_cleaned
        min_post_epoch = min([[int(dim.value) for dim in post.dimensions if dim.key == 'taken_at_timestamp'][0]
                              for post in posts_data_cleaned])
        min_req_epoch = int((datetime.now() - timedelta(days=days)).strftime('%s'))
        if min_post_epoch < min_req_epoch or len(all_posts_data) > max_posts_to_fetch:
            break
    return all_posts_data

@sessionize
async def fetch_recent_posts_custom(profile_id: str, max_posts_to_fetch: int = 30, session=None) -> List[InstagramPostLog]:
    all_posts_data = []
    has_next_page = True
    cursor = None
    retry_counter = 0
    while has_next_page:
        try:
            posts_data, has_next_page, cursor, source = await retrieve_profile_posts_by_id(profile_id,
                                                                                           cursor,
                                                                                           session=session)
            posts_data_cleaned = await parse_profile_post_data(profile_id, posts_data, source)
            all_posts_data += posts_data_cleaned
            if len(all_posts_data) > max_posts_to_fetch:
                break
        except Exception as e:
            if retry_counter > 5:
                logger.error("Max Retries Reached - ")
                return all_posts_data
            retry_counter += 1
            logger.error("Error - Retrying in 10s")
            await sleep(10)
    return all_posts_data

@sessionize
async def fetch_recent_reels(profile_id: str, days: int = 30, session=None) -> List[InstagramPostLog]:
    all_posts_data = []
    has_next_page = True
    cursor = None
    max_posts_to_fetch = 10
    while has_next_page:
        posts_data, has_next_page, cursor, source = await retrieve_reels_posts_by_id(profile_id,
                                                                                       cursor,
                                                                                       session=session)
        posts_data_cleaned = await parse_reels_post_data(profile_id, posts_data, source)
        all_posts_data += posts_data_cleaned
        min_post_epoch = min([[int(dim.value) for dim in post.dimensions if dim.key == 'taken_at_timestamp'][0]
                              for post in posts_data_cleaned])
        min_req_epoch = int((datetime.now() - timedelta(days=days)).strftime('%s'))
        if min_post_epoch < min_req_epoch or len(all_posts_data) > max_posts_to_fetch:
            break
    return all_posts_data

@sessionize
async def fetch_recent_reels_custom(profile_id: str, max_posts_to_fetch: int = 30, session=None) -> List[InstagramPostLog]:
    all_posts_data = []
    has_next_page = True
    cursor = None
    while has_next_page:
        posts_data, has_next_page, cursor, source = await retrieve_reels_posts_by_id(profile_id,
                                                                                       cursor,
                                                                                       session=session)
        posts_data_cleaned = await parse_reels_post_data(profile_id, posts_data, source)
        all_posts_data += posts_data_cleaned
        if len(all_posts_data) > max_posts_to_fetch:
            break
    return all_posts_data


@sessionize
async def fetch_profile(profile_id: str, handle: str, session=None) -> Tuple[InstagramProfileLog, List[InstagramPostLog]]:
    data, source = None, None
    if profile_id and profile_id != "":
        data, source = await retrieve_profile_data_by_id(profile_id, session=session)
    elif handle and handle != "":
        data, source = await retrieve_profile_data_by_handle(handle, session=session)
    if data:
        profile_data, recent_posts_data = await parse_profile_data(data, source)
        return profile_data, recent_posts_data
    else:
        raise Exception("Unable to fetch data for profile")

@sessionize
async def fetch_post(shortcode: str, session=None) -> InstagramPostLog:
    data, source = await retrieve_post_data_by_shortcode(shortcode, session=session)
    if data:
        post_data = await parse_post_data(data, source)
        return post_data
    else:
        raise Exception("Unable to fetch data for profile")


@sessionize
async def fetch_post_insights(handle: str, shortcode: str, post_id: str, post_type: str, session=None) -> InstagramPostLog:
    data, source = await retrieve_post_insights(handle, post_id, post_type, session=session)
    if data:
        insights_data = await parse_post_insights(data, shortcode, source)
        return insights_data
    else:
        raise Exception("Unable to fetch data for profile")

@sessionize
async def fetch_profile_insights(handle: str, session=None) -> InstagramProfileLog:
    data, source = await retrieve_profile_insights(handle, session=session)
    if data:
        insights_data = await parse_profile_insights(data, source)
        profile = await get(session, InstagramAccount, handle=handle)
        insights_data.profile_id = profile.profile_id
        insights_data.handle = handle
        return insights_data
    else:
        raise Exception("Unable to fetch data for profile")

@sessionize
async def fetch_profile_insights_from_graphapi(handle: str, session=None) -> InstagramProfileLog:
    data = await retrieve_profile_insights_from_graphapi(handle, session=session)
    if data:
        insights_data = await parse_profile_insights(data, 'graphapi')
        profile = await get(session, InstagramAccount, handle=handle)
        insights_data.profile_id = profile.profile_id
        insights_data.handle = handle
        return insights_data
    else:
        raise Exception("Unable to fetch data for profile")

@sessionize
async def fetch_profile_from_graphapi(handle: str, session=None) -> Tuple[InstagramProfileLog, List[InstagramPostLog]]:
    data = await retrieve_profile_data_from_graphapi(handle, session=session)
    if data:
        profile_data, recent_posts_data = await parse_profile_data(data, 'graphapi')
        return profile_data, recent_posts_data
    else:
        raise Exception("Unable to fetch data for profile")


@sessionize
async def fetch_story_insights(story_id: str, handle: str, session=None) -> InstagramPostLog:
    data, source = await retrieve_story_insights(story_id, handle, session=session)
    if data:
        story_insights_data = await parse_story_insights(story_id, data, source)
        return story_insights_data
    else:
        raise Exception("Unable to fetch data for story")

@sessionize
async def fetch_tagged_posts_by_profile_id(profile_id: str, max_posts: int, session=None) -> list[InstagramPostLog]:
    all_posts_data = []
    has_next_page = True
    cursor = None
    retry_counter = 0
    while has_next_page:
        try:
            posts_data, has_next_page, next_cursor, source = await retrieve_tagged_posts_by_profile_id(profile_id, cursor=cursor, session=session)
            posts_data_cleaned = await parse_tagged_posts_data(posts_data, profile_id, source)
            all_posts_data += posts_data_cleaned
            cursor = next_cursor
            if len(all_posts_data) > max_posts:
                    break
        except Exception as e:
            if retry_counter > 5:
                logger.error("Max Retries Reached - ")
                return all_posts_data
            retry_counter += 1
            logger.error("Error - Retrying in 10s")
            await sleep(10)
    return all_posts_data

@sessionize
async def fetch_stories_posts(handle: str, days=90, session=None) -> InstagramPostLog:
    all_posts_data = []
    max_posts_to_fetch = 150
    stories_posts_data, has_next_page, cursor, source = await retrieve_stories_posts_data(handle, session=session)
    profile_data, recent_post_story_data = await parse_profile_data(stories_posts_data, source)
    profile_id = profile_data.profile_id
    previous_cursor = cursor
    while has_next_page:
        posts_data, has_next_page, cursor, source = await retrieve_stories_posts_data(handle, cursor,
                                                                                              session=session)
        post_story_data = await parse_profile_post_data(profile_id, posts_data["business_discovery"]["media"]["data"], source)
        all_posts_data += post_story_data
        min_post_epoch = min([[int(dim.value) for dim in post.dimensions if dim.key == 'taken_at_timestamp'][0] for post in post_story_data])
        min_req_epoch = int((datetime.now() - timedelta(days=days)).strftime('%s'))
        if min_post_epoch < min_req_epoch or len(all_posts_data) > max_posts_to_fetch:
            break
    return profile_id, profile_data, recent_post_story_data+all_posts_data
@sessionize
async def process_profile_data(handle: str,
                               profile_data: InstagramProfileLog,
                               recent_posts: Optional[List[InstagramPostLog]],
                               session=None) -> Tuple[InstagramAccount, List[InstagramPost]]:
    return await upsert_profile(handle, profile_data, recent_posts, session=session)


@sessionize
async def process_post_data(shortcode: str, post: InstagramPostLog, session=None) -> None:
    await upsert_post(shortcode, post, session=session)


@sessionize
async def process_insights_data(shortcode: str, insights_data: InstagramPostLog, session=None) -> None:
    await upsert_post_insights(shortcode, insights_data, session=session)

@sessionize
async def process_profile_insights_data(insights_data: InstagramProfileLog, session=None) -> None:
    await upsert_profile_insights(insights_data, session=session)


@sessionize
async def process_story_insights_data(insights_data: InstagramPostLog, profile_id: str, session=None) -> None:
    await upsert_story_insights(insights_data, profile_id, session=session)


@sessionize
async def process_stories_posts_data(profile_id: str, post_stories_data: InstagramPostLog,
                                     session=None) -> None:
    await upsert_stories_posts(profile_id, post_stories_data, session=session)

"""
    Possible Flags
        - metrics_set: [profile, posts, reels]
        - lookback_period: int
"""


@sessionize
async def refresh_profile(profile_id: str, handle: str, session=None) -> None:

    start_time = time.perf_counter()
    profile_data, recent_posts = await fetch_profile(profile_id, handle, session=session)
    end_time = time.perf_counter()
    logger.debug(f"Time taken for fetch_profile - {end_time - start_time}")

    profile_id = profile_data.profile_id
    if recent_posts and len(recent_posts) > 0:
        start_time = time.perf_counter()
        await process_profile_data(profile_id, profile_data, recent_posts, session=session)
        end_time = time.perf_counter()
        logger.debug(f"Time taken for process_profile_data_1 - {end_time - start_time}")

    try:
        start_time = time.perf_counter()
        recent_posts = await fetch_recent_posts(profile_id=profile_id, session=session)  # Also gets reels data
        end_time = time.perf_counter()
        logger.debug(f"Time taken for fetch_recent_posts - {end_time - start_time}")
    except Exception as ex:
        logger.error(ex)
        logger.error("Unable to fetch profile posts for profile_id: %s" % profile_id)

    try:
        media_count_metric = [d for d in profile_data.metrics if d.key == MEDIA_COUNT]
        if media_count_metric and len(media_count_metric) > 0 and media_count_metric[0].value and int(media_count_metric[0].value) > 0:
            start_time = time.perf_counter()
            recent_reels = await fetch_recent_reels(profile_id=profile_id, session=session)  # Also gets reels data
            if not recent_posts:
                recent_posts = []
            recent_posts += recent_reels
            end_time = time.perf_counter()
            logger.debug(f"Time taken for fetch_recent_reels - {end_time - start_time}")
    except Exception as ex:
        logger.error(ex)
        logger.error("Unable to fetch profile reels for profile_id: %s" % profile_id)

    finally:
        start_time = time.perf_counter()
        await process_profile_data(profile_id, profile_data, recent_posts, session=session)
        end_time = time.perf_counter()
        logger.debug(f"Time taken for process_profile_data_2 - {end_time - start_time}")


@sessionize
async def refresh_profile_custom(profile_id: str, handle: str, max_posts: int, session=None) -> None:
    start_time = time.perf_counter()
    profile_data, recent_posts = await fetch_profile(profile_id, handle, session=session)
    profile_id = profile_data.profile_id
    end_time = time.perf_counter()
    logger.debug(f"Time taken for fetch_profile - {end_time - start_time}")
    start_time = time.perf_counter()
    recent_posts = await fetch_recent_posts_custom(profile_id=profile_id, max_posts_to_fetch=max_posts, session=session)
    end_time = time.perf_counter()
    logger.debug(f"Time taken for fetch_recent_posts - {end_time - start_time}")
    start_time = time.perf_counter()
    await process_profile_data(profile_id, profile_data, recent_posts, session=session)
    end_time = time.perf_counter()
    logger.debug(f"Time taken for process_profile_data_2 - {end_time - start_time}")


@sessionize
async def refresh_profile_basic(profile_id: str, handle: str, session=None) -> None:
    profile_data, recent_posts = await fetch_profile(profile_id, handle, session=session)
    if not recent_posts:
        recent_posts = []
    profile_id = profile_data.profile_id
    await process_profile_data(profile_id, profile_data, recent_posts, session=session)

@sessionize
async def refresh_profile_by_handle(handle: str, session=None) -> None:
    await refresh_profile(None, handle, session=session)


@sessionize
async def refresh_profile_by_profile_id(profile_id: str, session=None) -> None:
    await refresh_profile(profile_id, None, session=session)


@sessionize
async def refresh_post_by_shortcode(shortcode: str, session=None) -> None:
    post_data = await fetch_post(shortcode, session=session)
    await process_post_data(shortcode, post_data, session=session)


@sessionize
async def refresh_post_insights(shortcode: str, handle: str, post_id: str, post_type: str, session=None) -> None:
    insights_data = await fetch_post_insights(handle, shortcode, post_id, post_type, session=session)
    await process_insights_data(shortcode, insights_data, session=session)


@sessionize
async def refresh_profile_insights(handle: str, session=None) -> None:
    profile_insights_data = await fetch_profile_insights(handle, session=session)
    await process_profile_insights_data(profile_insights_data, session=session)

@sessionize
async def refresh_profile_by_handle_from_graphapi(handle:str, session=None) -> None:
    profile_data, recent_posts_data = await fetch_profile_from_graphapi(handle, session=session)
    await process_profile_data(profile_data.profile_id, profile_data, recent_posts_data, session=session)

@sessionize
async def refresh_stories_posts(handle: str, session=None) -> None:
    profile_id, profile_data, post_stories_data = await fetch_stories_posts(handle, session=session)
    await process_stories_posts_data(profile_id, post_stories_data, session=session)


@sessionize
async def refresh_story_insights(story_id: str, shortcode: str, profile_id: str, handle: str, session=None) -> None:
    # Instagram shortcode and post_id are different now Using arraybobo or jotucker, we get the same value for both shortcode and post_id
    # But this value is not a valid post_id, so we cannot fetch story insights for it
    # If post_id and shortcode are the same, then we use graph API with handle to update story data Then we can fetch story id from pg using this shortcode
    # We can use this post_id to fetch story insights
    story = await get(session, InstagramPost, shortcode=shortcode)
    if story.post_id == story.shortcode:
        await refresh_stories_posts(handle, session=session)
        story = await get(session, InstagramPost, shortcode=shortcode)
        story_id = story.post_id
    
    insights_data = await fetch_story_insights(story_id, handle, session=session)
    insights_data.shortcode = shortcode
    await process_story_insights_data(insights_data, profile_id, session=session)

@sessionize
async def refresh_tagged_posts_by_profile_id(profile_id: str, max_posts: int, session=None) -> None:
    posts = await fetch_tagged_posts_by_profile_id(profile_id, max_posts, session=session)
    for post in posts:
        shortcode = post.shortcode
        await process_post_data(shortcode, post, session=session)

@sessionize
async def fetch_story_posts_by_profile_id(profile_id: str, session=None) -> list[InstagramPostLog]:
    stories_data, source = await retrieve_story_posts_by_profile_id(profile_id, cursor=None, session=session)
    stories_data_cleaned = []
    if stories_data:
        stories_data_cleaned = await parse_story_posts_data(profile_id,stories_data, source)            
    return stories_data_cleaned

@sessionize
async def refresh_story_posts_by_profile_id(profile_id: str, handle: str, story_ids: list=None, profile_follower: int=None, profile_er: float=None, profile_static_er: float=None, profile_reels_er:float = None, session=None) -> None:
    posts = await fetch_story_posts_by_profile_id(profile_id, session=session)
    story_ids_from_api = []
    for post in posts:
        shortcode = post.shortcode
        story_ids_from_api.append(shortcode)
        post.metrics.append(metric(profile_follower, PROFILE_FOLLOWER))
        post.metrics.append(metric(profile_er, PROFILE_ER))
        post.metrics.append(metric(profile_static_er, PROFILE_STATIC_ER))
        post.metrics.append(metric(profile_reels_er, PROFILE_REELS_ER))
        await process_post_data(shortcode, post, session=session)
    
    if story_ids is not None:
        not_present = list(set(story_ids) - set(story_ids_from_api))
        for id in not_present:
            shortcode = id
            post = InstagramPostLog(
                shortcode=id,
                source="",
                metrics=[],
                dimensions=[])

            post.dimensions.append(dimension(id, POST_ID))
            post.dimensions.append(dimension(handle, HANDLE))
            post.dimensions.append(dimension(profile_id, PROFILE_ID))
            post.dimensions.append(dimension("story", POST_TYPE))
            post.metrics.append(metric(profile_follower, PROFILE_FOLLOWER))
            post.metrics.append(metric(profile_er, PROFILE_ER))
            post.metrics.append(metric(profile_static_er, PROFILE_STATIC_ER))
            post.metrics.append(metric(profile_reels_er, PROFILE_REELS_ER))
            await process_post_data(shortcode, post, session=session)