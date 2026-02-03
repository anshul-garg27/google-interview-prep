from datetime import datetime
from loguru import logger
from sqlalchemy import update

from core.amqp.amqp import publish
from utils.extracted_keyword import get_extracted_keywords
from utils.getter import dimension
from core.entities.entities import ProfileLog, PostLog, PostActivityLog, ProfileActivityLog, YtProfileRelationshipLog
from core.enums import enums
from core.helpers.session import sessionize
from core.models.models import Context

from utils.db import get_or_create
from youtube.entities.entities import YoutubeAccount, YoutubePost, YoutubeAccountTimeSeries, \
    YoutubeProfileInsights
from youtube.metric_dim_store import KEYWORDS, PUBLISH_TIME, TITLE, TEXT
from youtube.models.models import YoutubeProfileLog, YoutubePostLog, YoutubePostActivityLog, \
    YoutubeActivityLog, YoutubeProfileRelationshipLog
from youtube.tasks.transformer import profile_map, post_map, profile_map_ts, profile_insights_map

from utils.request import make_scrape_log_event

"""
    Process Data and Update DB
"""


@sessionize
async def upsert_profiles(profile_logs: list[YoutubeProfileLog], session=None) -> None:
    for log in profile_logs:
        await upsert_profile(log, session=session)


@sessionize
async def upsert_posts(post_logs: list[YoutubePostLog], session=None) -> None:
    for log in post_logs:
        await upsert_post(log, session=session)

@sessionize
async def upsert_activities(post_logs: list[YoutubePostLog], session=None) -> None:
    for log in post_logs:
        await upsert_post(log, session=session)


@sessionize
async def upsert_profile(profile_log: YoutubeProfileLog, session=None) -> None:
    now = datetime.now()
    context = Context(now)
    profile = ProfileLog(
        platform=enums.Platform.YOUTUBE.name,
        profile_id=profile_log.channel_id,
        metrics=[m.__dict__ for m in profile_log.metrics],
        dimensions=[d.__dict__ for d in profile_log.dimensions],
        source=profile_log.source,
        timestamp=now
    )
    # session.add(profile)
    await upsert_yt_account(context, profile_log, profile_log.channel_id,session=session)
    # await upsert_yt_account_ts(context, profile_log, profile_log.channel_id,session=session)
    await make_scrape_log_event("profile_log", profile)


@sessionize
async def upsert_post(post_log: YoutubePostLog, session=None) -> None:
    now = datetime.now()
    context = Context(now)
    logger.debug(post_log)
    if post_log.dimensions:
        title = next((item.value for item in post_log.dimensions if item.key == TITLE), "")
        post_log.dimensions.append(dimension(get_extracted_keywords(title) if title else '', KEYWORDS))
    post = PostLog(
        platform=enums.Platform.YOUTUBE.name,
        profile_id=post_log.channel_id,
        platform_post_id=post_log.shortcode,
        metrics=[m.__dict__ for m in post_log.metrics],
        dimensions=[d.__dict__ for d in post_log.dimensions],
        source=post_log.source,
        timestamp=now
    )
    #session.add(post)
    await upsert_yt_post(context, post_log,session=session)
    #await upsert_yt_post_ts(context, post_log)
    await make_scrape_log_event("post_log", post)


@sessionize
async def upsert_yt_account(context: Context, profile_log: YoutubeProfileLog, channel_id: str,
                            session=None) -> YoutubeAccount:
    await upsert_account(context, profile_log, channel_id, session=session)


@sessionize
async def upsert_account(context: Context, profile_log: YoutubeProfileLog, channel_id: str,
                         session=None) -> YoutubeAccount:
    account: YoutubeAccount = (await get_or_create(session, YoutubeAccount, channel_id=channel_id))[0]
    account.updated_at = context.now
    account.channel_id = profile_log.channel_id

    for metric in profile_log.metrics:
        if metric.key in profile_map:
            setattr(account, profile_map[metric.key], metric.value)
        else:
            logger.debug("Key is missing - %s" % metric.key)

    for dimension in profile_log.dimensions:
        if dimension.key in profile_map:
            if dimension.key == PUBLISH_TIME and isinstance(dimension.value, int):
                setattr(account, profile_map[dimension.key], datetime.fromtimestamp(dimension.value))
            else:
                setattr(account, profile_map[dimension.key], dimension.value)
        else:
            logger.debug("Key is missing - %s" % dimension.key)
    return account


@sessionize
async def upsert_yt_post(context: Context, post_log: YoutubePostLog, session=None) -> YoutubePost:
    post: YoutubePost = (await get_or_create(session, YoutubePost, shortcode=post_log.shortcode))[0]
    post.updated_at = context.now
    post.channel_id = post_log.channel_id

    for metric in post_log.metrics:
        if metric.key in post_map:
            setattr(post, post_map[metric.key], metric.value)
        else:
            logger.debug("Key is missing - %s" % metric.key)

    for dimension in post_log.dimensions:
        if dimension.key in post_map:
            if dimension.key == PUBLISH_TIME and isinstance(dimension.value, int):
                setattr(post, post_map[dimension.key], datetime.fromtimestamp(dimension.value))
            else:
                setattr(post, post_map[dimension.key], dimension.value)
        else:
            logger.debug("Key is missing - %s" % dimension.key)
    return post


# @sessionize
# async def upsert_yt_account_ts(context: Context, profile_log: YoutubeProfileLog,
#                                channel_id: str, session=None) -> YoutubeAccountTimeSeries:
#     account: YoutubeAccountTimeSeries = YoutubeAccountTimeSeries(
#         channel_id=channel_id,
#         created_at=datetime.now()
#     )

#     account.channel_id = profile_log.channel_id
#     account.timestamp = context.now
#     for metric in profile_log.metrics:
#         if metric.key in profile_map_ts:
#             setattr(account, profile_map_ts[metric.key], metric.value)
#         else:
#             logger.debug("Key is missing - %s" % metric.key)

#     for dimension in profile_log.dimensions:
#         if dimension.key in profile_map_ts:
#             if dimension.key == PUBLISH_TIME and isinstance(dimension.value, int):
#                 setattr(account, profile_map_ts[dimension.key], datetime.fromtimestamp(dimension.value))
#             else:
#                 setattr(account, profile_map_ts[dimension.key], dimension.value)
#         else:
#             logger.debug("Key is missing - %s" % dimension.key)

#     session.add(account)
#     return account


# @sessionize
# async def upsert_yt_post_ts(context: Context, post_log: YoutubePostLog, session=None) -> YoutubePostTimeSeries:
#     post: YoutubePostTimeSeries = YoutubePostTimeSeries(
#         shortcode=post_log.shortcode,
#         created_at=datetime.now()
#     )

#     post.channel_id = post_log.channel_id
#     post.timestamp = context.now
#     for metric in post_log.metrics:
#         if metric.key in post_map_ts:
#             setattr(post, post_map_ts[metric.key], metric.value)
#         else:
#             logger.debug("Key is missing - %s" % metric.key)

#     for dimension in post_log.dimensions:
#         if dimension.key in post_map_ts:
#             if dimension.key == PUBLISH_TIME and isinstance(dimension.value, int):
#                 setattr(post, post_map[dimension.key], datetime.fromtimestamp(dimension.value))
#             else:
#                 setattr(post, post_map_ts[dimension.key], dimension.value)
#         else:
#             logger.debug("Key is missing - %s" % dimension.key)

#     session.add(post)
#     return post


@sessionize
async def upsert_profile_insights(profile_insights_log: YoutubeProfileLog, session=None) -> None:
    now = datetime.now()
    context = Context(now)
    profile_insight_log = ProfileLog(
        platform=enums.Platform.YOUTUBE.name,
        profile_id=profile_insights_log.channel_id,
        metrics=[m.__dict__ for m in profile_insights_log.metrics],
        dimensions=[d.__dict__ for d in profile_insights_log.dimensions],
        source=profile_insights_log.source,
        timestamp=now
    )
    # session.add(profile_insight_log)
    await upsert_yt_profile_insights(context, profile_insights_log, session=session)
    await make_scrape_log_event("profile_log", profile_insight_log)


@sessionize
async def upsert_yt_profile_insights(context: Context, profile_insights_log: YoutubeProfileLog,
                                     session=None) -> YoutubeProfileInsights:
    profile_insight: YoutubeProfileInsights = \
        (await get_or_create(session, YoutubeProfileInsights, channel_id=profile_insights_log.channel_id))[0]

    profile_insight.updated_at = context.now
    profile_insight.channel_id = profile_insights_log.channel_id
    for metric in profile_insights_log.metrics:
        if metric.key in profile_map:
            setattr(profile_insight, profile_map[metric.key], metric.value)
        else:
            logger.debug("Key is missing - %s" % metric.key)

    for dimension in profile_insights_log.dimensions:
        if dimension.key in profile_insights_map:
            setattr(profile_insight, profile_insights_map[dimension.key], dimension.value)
        else:
            logger.debug("Key is missing - %s" % dimension.key)
    return profile_insight


@sessionize
async def upsert_post_type(post_id: str, post_type: str, session=None):
    update_query = update(YoutubePost).values({"post_type": post_type}).where(YoutubePost.shortcode == post_id)
    await session.execute(update_query)


@sessionize
async def upsert_activities(logs: list[YoutubePostActivityLog], session=None) -> None:
    now = datetime.now()
    post_activity_logs = []
    for log in logs:
        if log.dimensions:
            comment = next((item.value for item in log.dimensions if item.key == TEXT), "")
            log.dimensions.append(dimension(get_extracted_keywords(comment) if comment else '', KEYWORDS))

        post = PostActivityLog(
            activity_type=log.activity_type,
            platform=enums.Platform.YOUTUBE.name,
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

    activity_type = post_activity_logs[0]['activity_type'] if len(post_activity_logs)>0 else None
    payload = {
        "platform": enums.Platform.YOUTUBE.name,
        "activity_type": activity_type,
        "post_activity_logs": post_activity_logs
    }
    publish(payload, "beat.dx", "post_activity_log_bulk")


@sessionize
async def upsert_yt_activities(logs: list[YoutubeActivityLog], session=None) -> None:
    for log in logs:
        activity_log = ProfileActivityLog(
            activity_id = log.activity_id,
            activity_type=log.activity_type,
            actor_channel_id=log.actor_channel_id,
            platform=enums.Platform.YOUTUBE.name,
            source=log.source,
            dimensions=[m.__dict__ for m in log.dimensions],
            metrics=[m.__dict__ for m in log.metrics],
        )

        await make_scrape_log_event("yt_activity_log", activity_log)


@sessionize
async def upsert_yt_profile_relationship(logs: list[YoutubeProfileRelationshipLog], session=None) -> None:
    for log in logs:
        profile_relationship_log = YtProfileRelationshipLog(
            relationship_id=log.relationship_id,
            relationship_type=log.relationship_type,
            source_channel_id=log.source_channel_id,
            target_channel_id=log.target_channel_id,
            platform=enums.Platform.YOUTUBE.name,
            source=log.source,
            source_dimensions=[m.__dict__ for m in log.source_dimensions],
            source_metrics=[m.__dict__ for m in log.source_metrics],
            target_dimensions=[m.__dict__ for m in log.target_dimensions],
            target_metrics=[m.__dict__ for m in log.target_metrics],
            subscribed_on=log.subscribed_on
        )
        await make_scrape_log_event("yt_profile_relationship_logs", profile_relationship_log)

