import time
from datetime import datetime
from typing import List, Optional, TypeVar, Tuple

from loguru import logger

from core.amqp.amqp import publish
from utils.getter import dimension, metric
from utils.extracted_keyword import get_extracted_keywords

from core.entities.entities import ProfileLog, PostLog, ProfileRelationshipLog, PostActivityLog
from core.enums import enums
from core.helpers.session import sessionize
from core.models.models import Dimension, Context
from instagram.entities.entities import InstagramAccount, InstagramPost, InstagramPostTimeSeries, \
    InstagramAccountTimeSeries, InstagramProfileInsights
from instagram.metric_dim_store import CAPTION, COLLAB_PROFILE_ID, KEYWORDS, PROFILE_ID, PUBLISH_TIME, TEXT
from instagram.models.models import InstagramProfileLog, InstagramPostLog, \
    InstagramRelationshipLog, InstagramPostActivityLog
from instagram.tasks.transformer import profile_map, profile_map_ts, post_map, profile_insights_map
from utils.db import get_or_create
from utils.request import make_scrape_log_event

"""
    Process Data and Update DB
"""
T = TypeVar("T")


def set_attr(model: T, key: str, value):
    if value is not None and value != "":
        setattr(model, key, value)


@sessionize
async def upsert_post(shortcode: str, post_log: InstagramPostLog, profile_id: str = None, session=None) -> None:
    now = datetime.now()
    context = Context(now)
    if post_log.dimensions:
        profile_id = next((item.value for item in post_log.dimensions if item.key==PROFILE_ID),None)
        caption = next((item.value for item in post_log.dimensions if item.key==CAPTION),"")
        post_log.dimensions.append(dimension(get_extracted_keywords(caption) if caption else '', KEYWORDS))
    post = PostLog(
        platform=enums.Platform.INSTAGRAM.name,
        metrics=[m.__dict__ for m in post_log.metrics],
        dimensions=[d.__dict__ for d in post_log.dimensions],
        platform_post_id=post_log.shortcode,
        profile_id=profile_id if profile_id else "TODO",
        source=post_log.source,
        timestamp=now
    )
    #session.add(post)
    await upsert_insta_post(context, post_log, session=session)
    #await upsert_insta_post_ts(context, post_log)
    await make_scrape_log_event("post_log", post)


@sessionize
async def upsert_stories_posts(profile_id: str, posts_stories_log: List[InstagramPostLog], session=None) -> None:
    for _post_story in posts_stories_log:
        _post_story.dimensions.append(Dimension(PROFILE_ID, profile_id))
        await upsert_post(_post_story.shortcode, _post_story, profile_id, session=session)


@sessionize
async def upsert_post_insights(shortcode: str, post_log: InstagramPostLog, session=None) -> None:
    now = datetime.now()
    context = Context(now)
    post = PostLog(
        platform=enums.Platform.INSTAGRAM.name,
        metrics=[m.__dict__ for m in post_log.metrics],
        dimensions=[d.__dict__ for d in post_log.dimensions],
        platform_post_id=post_log.shortcode,
        profile_id="TODO",
        source=post_log.source,
        timestamp=now
    )
    #session.add(post)
    await upsert_insta_post_insights(context, post_log, shortcode, session=session)
    #await upsert_insta_post_insights_ts(context, post_log, shortcode, session=session)
    await make_scrape_log_event("post_log", post)


@sessionize
async def upsert_profile_insights(profile_log: InstagramProfileLog, session=None) -> None:
    now = datetime.now()
    context = Context(now)
    profile_insights = ProfileLog(
        platform=enums.Platform.INSTAGRAM.name,
        profile_id=profile_log.profile_id,
        source=profile_log.source,
        metrics=[m.__dict__ for m in profile_log.metrics],
        dimensions=[m.__dict__ for m in profile_log.dimensions],
        timestamp=now
    )
    # session.add(profile_insights)
    await upsert_insta_profile_insights(context, profile_log, session=session)
    await make_scrape_log_event("profile_log", profile_insights)


@sessionize
async def upsert_story_insights(story_log: InstagramPostLog, profile_id: str, session=None) -> None:
    now = datetime.now()
    context = Context(now)
    story = PostLog(
        platform=enums.Platform.INSTAGRAM.name,
        metrics=[m.__dict__ for m in story_log.metrics],
        dimensions=[d.__dict__ for d in story_log.dimensions],
        platform_post_id=story_log.shortcode,
        profile_id=profile_id,
        source=story_log.source,
        timestamp=now
    )
    #session.add(story)
    await upsert_insta_story_insights(context, story_log, session=session)
    #await upsert_insta_story_insights_ts(context, story_log, session=session)
    await make_scrape_log_event("post_log", story)


@sessionize
async def upsert_profile(profile_id: str, profile_log: InstagramProfileLog,
                         recent_posts_log: Optional[List[InstagramPostLog]], session=None) -> Tuple[InstagramAccount, List[InstagramPost]]:
    now = datetime.now()
    context = Context(now)
    start_time = time.perf_counter()
    profile = ProfileLog(
        platform=enums.Platform.INSTAGRAM.name,
        profile_id=profile_id,
        metrics=[m.__dict__ for m in profile_log.metrics],
        dimensions=[d.__dict__ for d in profile_log.dimensions],
        source=profile_log.source,
        timestamp=now
    )
    # session.add(profile)
    await make_scrape_log_event("profile_log", profile)
    account = await upsert_insta_account(context, profile_log, profile_id, session=session)
    # await upsert_insta_account_ts(context, profile_log, profile_id, session=session)
    insta_posts = []
    end_time = time.perf_counter()
    logger.debug(f"Time taken for account insertion- {end_time - start_time}")
    start_time = time.perf_counter()

    if not recent_posts_log:
        recent_posts_log = []

    for _post in recent_posts_log:
        post_profile_id = None
        if _post.dimensions:
            post_profile_id = next((item.value for item in _post.dimensions if item.key==PROFILE_ID),None)

        if post_profile_id and post_profile_id!=profile_id:
            _post.dimensions.append(Dimension(COLLAB_PROFILE_ID, profile_id))

        if not post_profile_id:
            post_profile_id = profile_id
            
        if _post.dimensions:
            caption = next((item.value for item in _post.dimensions if item.key==CAPTION),"")
            _post.dimensions.append(dimension(get_extracted_keywords(caption) if caption else '', KEYWORDS))
        post = PostLog(
            platform=enums.Platform.INSTAGRAM.name,
            metrics=[m.__dict__ for m in _post.metrics],
            dimensions=[d.__dict__ for d in _post.dimensions],
            platform_post_id=_post.shortcode,
            profile_id=post_profile_id,
            source=_post.source,
            timestamp=now
        )
        #session.add(post)
        insta_post = await upsert_insta_post(context, _post, session=session)
        insta_posts.append(insta_post)
        #await upsert_insta_post_ts(context, _post, session=session)
        await make_scrape_log_event("post_log", post)
    end_time = time.perf_counter()
    logger.debug(f"Time taken for post insertion ({len(recent_posts_log)})- {end_time - start_time}")
    return account, insta_posts


@sessionize
async def upsert_insta_post(context: Context, post_log: InstagramPostLog, session=None) -> InstagramPost:
    post: InstagramPost = (await get_or_create(session, InstagramPost, shortcode=post_log.shortcode))[0]
    post.updated_at = context.now
    for metric in post_log.metrics:
        if metric.key in post_map:
            if metric.key in ['profile_follower', 'profile_er', 'profile_static_er', 'profile_reels_er']:
                if metric.key == 'profile_follower' and post.profile_follower is None:
                    set_attr(post, post_map[metric.key], metric.value)
                elif metric.key == 'profile_er' and post.profile_er is None:
                    set_attr(post, post_map[metric.key], metric.value)
                elif metric.key == 'profile_static_er' and post.profile_static_er is None:
                    set_attr(post, post_map[metric.key], metric.value)
                elif metric.key == 'profile_reels_er' and post.profile_reels_er is None:
                    set_attr(post, post_map[metric.key], metric.value)
            else:
                set_attr(post, post_map[metric.key], metric.value)
        else:
            logger.debug("Key is missing - %s" % metric.key)     
    for dimension in post_log.dimensions:
        if dimension.key in post_map:
            if dimension.key == PUBLISH_TIME and isinstance(dimension.value, int):
                set_attr(post, post_map[dimension.key], datetime.fromtimestamp(dimension.value))
            else:
                set_attr(post, post_map[dimension.key], dimension.value)
        else:
            logger.debug("Key is missing - %s" % dimension.key)
    return post


# @sessionize
# async def upsert_insta_post_ts(context: Context, post_log: InstagramPostLog, session=None) -> InstagramPostTimeSeries:
#     post: InstagramPostTimeSeries = InstagramPostTimeSeries(
#         shortcode=post_log.shortcode,
#         created_at=datetime.now()
#     )
#     post.timestamp = context.now
#     for metric in post_log.metrics:
#         if metric.key in post_map_ts:
#             set_attr(post, post_map_ts[metric.key], metric.value)
#         else:
#             logger.debug("Key is missing - %s" % metric.key)

#     for dimension in post_log.dimensions:
#         if dimension.key in post_map_ts:
#             if dimension.key == PUBLISH_TIME and isinstance(dimension.value, int):
#                 set_attr(post, post_map[dimension.key], datetime.fromtimestamp(dimension.value))
#             else:
#                 set_attr(post, post_map_ts[dimension.key], dimension.value)
#         else:
#             logger.debug("Key is missing - %s" % dimension.key)
#     session.add(post)
#     return post


@sessionize
async def upsert_insta_account(context: Context, profile_log: InstagramProfileLog, profile_id: str,
                               session=None) -> InstagramAccount:
    insta_account: InstagramAccount = (await get_or_create(session, InstagramAccount, profile_id=profile_id))[0]
    insta_account.updated_at = context.now
    insta_account.handle = profile_log.handle

    for metric in profile_log.metrics:
        if metric.key in profile_map:
            set_attr(insta_account, profile_map[metric.key], metric.value)
        else:
            logger.debug("Key is missing - %s" % metric.key)

    for dimension in profile_log.dimensions:
        if dimension.key in profile_map:
            set_attr(insta_account, profile_map[dimension.key], dimension.value)
        else:
            logger.debug("Key is missing - %s" % dimension.key)
    return insta_account


# @sessionize
# async def upsert_insta_account_ts(context: Context,
#                                   profile_log: InstagramProfileLog,
#                                   profile_id: str, session=None) -> InstagramAccountTimeSeries:
#     insta_account: InstagramAccountTimeSeries = \
#         InstagramAccountTimeSeries(
#             profile_id=profile_id,
#             created_at=datetime.now()
#         )
#     insta_account.timestamp = context.now
#     for metric in profile_log.metrics:
#         if metric.key in profile_map_ts:
#             set_attr(insta_account, profile_map_ts[metric.key], metric.value)
#         else:
#             logger.debug("Key is missing - %s" % metric.key)
#     for dimension in profile_log.dimensions:
#         if dimension.key in profile_map_ts:
#             set_attr(insta_account, profile_map_ts[dimension.key], dimension.value)
#         else:
#             logger.debug("Key is missing - %s" % dimension.key)
#     session.add(insta_account)
#     return insta_account


@sessionize
async def upsert_insta_post_insights(context: Context, post_log: InstagramPostLog, shortcode: str,
                                     session=None) -> InstagramPost:
    post: InstagramPost = (await get_or_create(session, InstagramPost, shortcode=shortcode))[0]
    post.updated_at = context.now
    for metric in post_log.metrics:
        if metric.key in post_map:
            set_attr(post, post_map[metric.key], metric.value)
        else:
            logger.debug("Key is missing - %s" % metric.key)
    for dimension in post_log.dimensions:
        if dimension.key in post_map:
            set_attr(post, post_map[dimension.key], dimension.value)
        else:
            logger.debug("Key is missing - %s" % dimension.key)
    return post


# @sessionize
# async def upsert_insta_post_insights_ts(context: Context, post_log: InstagramPostLog, shortcode: str,
#                                         session=None) -> InstagramPost:
#     post: InstagramPostTimeSeries = InstagramPostTimeSeries(
#         shortcode=post_log.shortcode,
#         created_at=datetime.now()
#     )
#     post.publish_time = context.now

#     for metric in post_log.metrics:
#         if metric.key in post_map:
#             set_attr(post, post_map[metric.key], metric.value)
#         else:
#             logger.debug("Key is missing - %s" % metric.key)
#     for dimension in post_log.dimensions:
#         if dimension.key in post_map:
#             set_attr(post, post_map[dimension.key], dimension.value)
#         else:
#             logger.debug("Key is missing - %s" % dimension.key)

#     session.add(post)
#     return post


@sessionize
async def upsert_insta_story_insights(context: Context, story_log: InstagramPostLog, session=None) -> InstagramPost:
    story: InstagramPost = (await get_or_create(session, InstagramPost, shortcode=story_log.shortcode))[0]
    story.updated_at = context.now

    for metric in story_log.metrics:
        if metric.key in post_map:
            set_attr(story, post_map[metric.key], metric.value)
        else:
            logger.debug("Key is missing - %s" % metric.key)
    for dimension in story_log.dimensions:
        if dimension.key in post_map:
            set_attr(story, post_map[dimension.key], dimension.value)
        else:
            logger.debug("Key is missing - %s" % dimension.key)
    return story


# @sessionize
# async def upsert_insta_story_insights_ts(context: Context, story_log: InstagramPostLog, session=None) -> InstagramPost:
#     story: InstagramPostTimeSeries = InstagramPostTimeSeries(
#         shortcode=story_log.shortcode,
#         created_at=datetime.now()
#     )
#     story.publish_time = context.now

#     for metric in story_log.metrics:
#         if metric.key in post_map:
#             set_attr(story, post_map[metric.key], metric.value)
#         else:
#             logger.debug("Key is missing - %s" % metric.key)
#     for dimension in story_log.dimensions:
#         if dimension.key in post_map:
#             set_attr(story, post_map[dimension.key], dimension.value)
#         else:
#             logger.debug("Key is missing - %s" % dimension.key)

#     session.add(story)
#     return story


@sessionize
async def upsert_insta_profile_insights(context: Context, profile_log: InstagramProfileLog,
                                        session=None) -> InstagramProfileInsights:
    profile: InstagramProfileInsights = \
        (await get_or_create(session, InstagramProfileInsights, profile_id=profile_log.profile_id))[0]
    profile.updated_at = context.now
    profile.handle = profile_log.handle
    for metric in profile_log.metrics:
        if metric.key in profile_insights_map:
            set_attr(profile, str(profile_insights_map[metric.key]), metric.value)
        else:
            logger.debug("Key is missing - %s" % metric.key)
    for dimension in profile_log.dimensions:
        if dimension.key in profile_insights_map:
            set_attr(profile, str(profile_insights_map[dimension.key]), dimension.value)
        else:
            logger.debug("Key is missing - %s" % dimension.key)
    return profile


@sessionize
async def upsert_relationships(logs: list[InstagramRelationshipLog], session=None) -> None:
    now = datetime.now()
    for log in logs:
        post = ProfileRelationshipLog(
            relationship_type=log.relationship_type,
            platform=enums.Platform.INSTAGRAM.name,
            source_profile_id=log.source_profile_id,
            target_profile_id=log.target_profile_id,
            source=log.source,
            timestamp=now
        )
        post.source_dimensions = {dimension.key: dimension.value for dimension in log.source_dimensions}
        post.target_dimensions = {dimension.key: dimension.value for dimension in log.target_dimensions}
        post.source_metrics = {metric.key: metric.value for metric in log.source_metrics}
        post.target_metrics = {metric.key: metric.value for metric in log.target_metrics}
        session.add(post)
        await make_scrape_log_event("profile_relationship_log", post)


@sessionize
async def upsert_activities(logs: list[InstagramPostActivityLog], session=None) -> None:
    now = datetime.now()
    post_activity_logs = []
    for log in logs:
        if log.dimensions:
            comment = next((item.value for item in log.dimensions if item.key == TEXT), "")
            log.dimensions.append(dimension(get_extracted_keywords(comment) if comment else '', KEYWORDS))

        post = PostActivityLog(
            activity_type=log.activity_type,
            platform=enums.Platform.INSTAGRAM.name,
            platform_post_id=log.shortcode,
            actor_profile_id=log.actor_profile_id,
            dimensions=[m.__dict__ for m in log.dimensions],
            metrics=[m.__dict__ for m in log.metrics],
            source=log.source,
            timestamp=now
        )
        session.add(post)
        await make_scrape_log_event("post_activity_log", post)
        if isinstance(post, PostActivityLog):
            post_dict = post.__dict__
            del post_dict['_sa_instance_state']
            post_activity_logs.append(post_dict)

    activity_type = post_activity_logs[0]['activity_type'] if len(post_activity_logs) > 0 else None
    payload = {
        "platform": enums.Platform.INSTAGRAM.name,
        "activity_type": activity_type,
        "post_activity_logs": post_activity_logs
    }
    publish(payload, "beat.dx", "post_activity_log_bulk")
