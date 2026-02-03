import math
import os
import traceback
from datetime import datetime, timedelta
from functools import lru_cache
from typing import Optional, List, Tuple

import uvicorn
from asyncio_redis_rate_limit import RateSpec, RateLimiter, RateLimitError
from dotenv import load_dotenv
from fastapi import FastAPI, Header, Body
from fastapi import Response, Depends, Request
from loguru import logger
from redis.asyncio import Redis as AsyncRedis
from sqlalchemy import select, and_
from sqlalchemy.ext.asyncio import create_async_engine
from starlette import status
from sqlalchemy import update, select, and_

import config
from core.entities.entities import ScrapeRequestLog, AssetLog, Credential
from core.enums import enums
from core.helpers.session import sessionize_api, sessionize_scl_api
from core.models.models import RequestContext, Context, ScrapeRequest
from credentials.models import CredentialModel
from credentials.validator import CredentialValidator
from credentials.manager import CredentialManager
from instagram.entities.entities import InstagramAccount, InstagramPost
from instagram.flows.refresh_profile import refresh_profile, refresh_post_insights, \
    refresh_profile_by_handle_from_graphapi, refresh_profile_basic, refresh_story_posts_by_profile_id,refresh_post_by_shortcode, \
    refresh_profile_by_profile_id, refresh_profile_by_handle, fetch_profile_insights_from_graphapi, fetch_profile_from_graphapi, process_profile_data
from instagram.functions.retriever.lama.lama import Lama
from instagram.helper import get_reach, get_engagement_rate
from instagram.tasks.processing import upsert_insta_post
from utils.db import get
from utils.exceptions import BeatException, TokenValidationFalied, DataAccessExpired
from utils.request import make_scrape_request_log_event
from youtube.entities.entities import YoutubeAccount
from youtube.flows.refresh_profile import refresh_yt_profiles
from youtube.utils import fetch_yt_channel_id_from_handle
from instagram.models.models import InstagramProfileLog, InstagramPostLog

app = FastAPI()
load_dotenv()

csv_flows = ["fetch_keyword_videos_csv", "fetch_playlist_videos_csv", "fetch_channels_csv", "fetch_channel_videos_csv", "fetch_video_details_csv", 
             "fetch_channel_demographic_csv", "fetch_channel_daily_stats_csv", "fetch_channel_monthly_stats_csv","fetch_channel_audience_geography_csv",
              "fetch_channel_engagement_csv"]

@app.on_event("startup")
async def startup_event():
    logger.info("Server Starting Up")
    session_engine = create_async_engine(os.environ["PG_URL"], isolation_level="AUTOCOMMIT", echo=False, pool_size=50,
                                         max_overflow=50)
    session_engine_scl = create_async_engine(os.environ["PG_URL"], isolation_level="AUTOCOMMIT", echo=False,
                                             pool_size=10, max_overflow=5)
    await session_engine.connect()
    app.state.db = session_engine
    app.state.scl_db = session_engine_scl
    logger.info("Server Startup Done")


@app.on_event("shutdown")
async def shutdown_event():
    if not app.state.db:
        await app.state.db.close()
    logger.info("Server Shutdown")


@lru_cache()
def get_settings():
    return config.Settings()


def get_recent_post_metrics(account, recent_posts):
    if not recent_posts:
        return {}
    static_posts = []
    reels_posts = []
    static_posts_count = 0
    reels_posts_count = 0
    all_posts = sorted(recent_posts, key=lambda x: x['publish_time'], reverse=True)[0:12]
    all_posts_count = len(all_posts)
    for post in recent_posts:
        if 'metrics' not in post:
            post['metrics'] = {'likes': 0, 'comments': 0, 'play_count': 0}
        if 'likes' not in post['metrics'] or not post['metrics']['likes']:
            post['metrics']['likes'] = 0
        if 'play_count' not in post['metrics'] or not post['metrics']['play_count']:
            post['metrics']['play_count'] = 0
        if 'comments' not in post['metrics'] or not post['metrics']['comments']:
            post['metrics']['comments'] = 0
        if post['post_type'] == 'image' or post['post_type'] == 'carousel':
            static_posts.append(post)
            static_posts_count += 1
        if post['post_type'] == 'reels':
            reels_posts.append(post)
            reels_posts_count += 1

    static_posts.sort(key=lambda x: x['metrics']['likes'], reverse=True)
    if static_posts_count > 4:
        static_posts = static_posts[2:-2] # for outlier removal.
        static_posts_count -= 4

    reels_posts.sort(key=lambda x: x['metrics']['play_count'], reverse=True)
    if reels_posts_count > 4:
        reels_posts = reels_posts[2:-2]
        reels_posts_count -= 4

    all_posts.sort(key=lambda x: x['metrics']['likes'], reverse=True)
    if all_posts_count > 4:
        all_posts = all_posts[2:-2]
        all_posts_count -= 4

    likes = 0
    comments = 0
    for post in all_posts:
        if 'likes' in post['metrics'] and post['metrics']['likes']:
            likes += post['metrics']['likes']
        if 'comments' in post['metrics'] and post['metrics']['comments']:
            comments += post['metrics']['comments']

    reels_plays = 0
    reels_reach = 0
    reels_likes = 0
    reels_comments = 0
    for post in reels_posts:
        if 'play_count' in post['metrics'] and post['metrics']['play_count']:
            reels_plays += post['metrics']['play_count']
        if 'reach' in post['metrics'] and post['metrics']['reach']:
            reels_reach += post['metrics']['reach']
        if 'likes' in post['metrics'] and post['metrics']['likes']:
            reels_likes += post['metrics']['likes']
        if 'comments' in post['metrics'] and post['metrics']['comments']:
            reels_comments += post['metrics']['comments']
    static_reach = 0
    static_likes = 0
    static_comments = 0
    for post in static_posts:
        if 'reach' in post['metrics'] and post['metrics']['reach']:
            static_reach += post['metrics']['reach']
        if 'likes' in post['metrics'] and post['metrics']['likes']:
            static_likes += post['metrics']['likes']
        if 'comments' in post['metrics'] and post['metrics']['comments']:
            static_comments += post['metrics']['comments']
    static_metrics = {}
    reels_metrics = {}
    all_metrics = {}
    story_metrics = {}
    avg_engagement = 0
    if all_posts_count > 0 and account.followers > 0:
        all_metrics = {
            "avg_engagement": 100 * ((likes + comments) / (account.followers * all_posts_count)),
            "avg_likes": likes / all_posts_count,
            "avg_comments": comments / all_posts_count,
        }
        avg_engagement = all_metrics['avg_engagement']
    if static_posts_count > 0 and account.followers > 0:
        static_metrics = {
            "avg_engagement": 100 * ((static_likes + static_comments) / (account.followers * static_posts_count)),
            "avg_likes": static_likes / static_posts_count,
            "avg_comments": static_comments / static_posts_count,
            "avg_reach": static_reach / static_posts_count
        }
    if reels_posts_count > 0 and account.followers > 0:
        reels_metrics = {
            "avg_engagement": 100 * ((reels_likes + reels_comments) / (account.followers * reels_posts_count)),
            "avg_plays": reels_plays / reels_posts_count,
            "avg_likes": reels_likes / reels_posts_count,
            "avg_comments": reels_comments / reels_posts_count,
            "avg_reach": reels_reach / reels_posts_count,
        }
        avg_engagement = reels_metrics['avg_engagement']

    story_metrics = {
        "avg_reach": 0.0
    }
    if account.followers > 0:
        if avg_engagement > 0:
            story_metrics = {
                "avg_reach": round((-0.000025017) * account.followers + (
                        1.11 * (account.followers * abs(math.log((avg_engagement+2), 2)) * 2 / 100)))
            }
    return {
        "static": static_metrics,
        "reels": reels_metrics,
        "story": story_metrics,
        "all": all_metrics
    }


def transform_insta_profile(entity: Optional[InstagramAccount], recent_posts: Optional[InstagramPost]) -> dict:
    if entity is None:
        return {}
    post_metrics = {}
    profile_type = "personal"
    if entity.profile_type and entity.profile_type != "":
        profile_type = entity.profile_type
    else:
        if entity.is_business_or_creator:
            profile_type = "business_or_creator"
    if recent_posts:
        post_metrics = get_recent_post_metrics(entity, recent_posts)
    return {
        "handle": entity.handle,
        "fbid": entity.fbid,
        "profile_pic_url": entity.profile_pic_url,
        "profile_id": entity.profile_id,
        "biography": entity.biography,
        "following": entity.following,
        "followers": entity.followers,
        "full_name": entity.full_name,
        "is_private": entity.is_private,
        "updated_at": entity.updated_at,
        "post_metrics": post_metrics,
        "recent_posts": recent_posts,
        "profile_type": profile_type,
        "media_count": entity.media_count
    }


def transform_yt_profile(entity: Optional[YoutubeAccount]) -> dict:
    if entity is None:
        return {}
    return {
        "channel_id": entity.channel_id,
        "title": entity.title,
        "subscribers": entity.subscribers,
        "uploads": entity.uploads,
        "views": entity.views,
        "thumbnail": entity.thumbnail
    }


# def transform_insta_post(entity: Optional[InstagramPost]) -> dict:
#     if entity is None:
#         return {}
#     return {
#         "shortcode": entity.shortcode,
#         "post_id": entity.post_id,
#         "post_type": entity.post_type,
#         "publish_time": entity.publish_time,
#         "updated_at": entity.updated_at,
#         "metrics": {
#             "reach": entity.reach,
#             "views": entity.views,
#             "likes": entity.likes,
#             "play_count": entity.plays,
#             "comments": entity.comments
#         }
#     }


def transform_insta_recent_posts(entity: InstagramPost, account: InstagramAccount) -> dict:
    if entity is None:
        return {}

    insta_url = " "
    if entity.post_type == 'carousel' or entity.post_type == 'image':
        insta_url = "https://www.instagram.com/p/"
    elif entity.post_type == 'reels':
        insta_url = "https://www.instagram.com/reel/"

    reach = get_reach(entity, account) if account else 0
    engagement_rate = get_engagement_rate(entity, account) if account else 0
    thumbnail_url = entity.thumbnail_url
    if not thumbnail_url or thumbnail_url == "None" or thumbnail_url == "":
        thumbnail_url = None

    display_url = entity.display_url
    if not display_url or display_url == "None" or display_url == "":
        display_url = None

    return {
        "shortcode": entity.shortcode,
        "post_id": entity.post_id,
        "post_type": entity.post_type,
        "post_url": insta_url + entity.shortcode,
        "publish_time": entity.publish_time,
        "updated_at": entity.updated_at,
        "metrics": {
            "reach": reach,
            "views": entity.views,
            "likes": entity.likes,
            "comments": entity.comments,
            "play_count": entity.plays,
            "engagement_rate": engagement_rate
        },
        "dimensions": {
            "caption": entity.caption,
            "comments": entity.comments,
            "thumbnail_url": thumbnail_url,
            "display_url": display_url
        }
    }


@app.post("/tokens")
async def insert_token(token: dict) -> dict:
    return success("Inserted Successfully")


@app.get("/profiles/{platform}/byhandle/{handle}")
@app.get("/social-profile-service/{platform}/byhandle/{handle}")
@sessionize_api
async def get_profile_details_for_handle(request: Request, platform: enums.Platform, handle: str, full_refresh=False, force_refresh=False,
                                         session=None) -> dict:
    if platform == enums.Platform.INSTAGRAM:
        profile = await get(session, InstagramAccount, handle=handle)
        if force_refresh or not profile or profile.updated_at is None or profile.updated_at < datetime.now() - timedelta(days=1):
            try:
                redis = AsyncRedis.from_url(os.environ["REDIS_URL"])
                global_limit_day = RateSpec(requests=20000, seconds=86400)
                global_limit_minute = RateSpec(requests=60, seconds=60)
                handle_limit = RateSpec(requests=1, seconds=1)
                if full_refresh:
                    # TODO: Figure out a better way to stack rate limits
                    async with RateLimiter(
                            unique_key="refresh_profile_insta_daily",
                            backend=redis,
                            cache_prefix="beat_server_",
                            rate_spec=global_limit_day):
                        async with RateLimiter(
                                unique_key="refresh_profile_insta_per_minute",
                                backend=redis,
                                cache_prefix="beat_server_",
                                rate_spec=global_limit_minute):
                            async with RateLimiter(
                                    unique_key="refresh_profile_insta_per_handle_" + handle,
                                    backend=redis,
                                    cache_prefix="beat_server_",
                                    rate_spec=handle_limit):
                                await refresh_profile(None, handle, session=session)
                else:
                    await refresh_profile_basic(None, handle, session=session)
                profile = await get(session, InstagramAccount, handle=handle)
            except RateLimitError as e:
                logger.error(e)
                return error("Concurrency Limit Exceeded", 2)
            except Exception as e:
                logger.error(e)
                return error("Profile not found - %s" % handle, 404)
        recent_posts = None
        if full_refresh:
            profile_id = profile.profile_id
            recent_reels = []
            recent_statics = []
            recent_reels_data = await fetch_recent_posts_by_profile_id(request,
                                                                       platform,
                                                                       profile_id,
                                                                       post_types=['reels'],
                                                                       session=session)
            recent_statics_data = await fetch_recent_posts_by_profile_id(request,
                                                                         platform,
                                                                         profile_id,
                                                                         post_types=['image', 'carousel'],
                                                                         session=session)
            if 'data' in recent_reels_data:
                recent_reels = recent_reels_data['data']['recent_posts']
            if 'data' in recent_statics_data:
                recent_statics = recent_statics_data['data']['recent_posts']
            recent_posts = recent_reels + recent_statics
        response = {
            "platform": platform,
            **transform_insta_profile(profile, recent_posts=recent_posts)
        }
        return success("Profile Retrieved", profile=response)
    elif platform == enums.Platform.YOUTUBE:
        profile = await get(session, YoutubeAccount, channel_id=handle)
        if not profile or profile.updated_at is None or profile.updated_at < datetime.now() - timedelta(days=1):
            try:
                await refresh_yt_profiles([handle], session=session)
                profile = await get(session, YoutubeAccount, channel_id=handle)
            except Exception as e:
                logger.error(traceback.format_exc())
                logger.error(e)
                return error("Profile not found - %s" % handle, 404)
        response = {
            "platform": platform,
            **transform_yt_profile(profile)
        }
        return success("Profile Retrieved", profile=response)
    return error("Platform not supported - %s" % platform, 400)


@app.get("/profiles/{platform}/byprofileid/{profile_id}")
@sessionize_api
async def get_profile_details_for_profile_id(request: Request, platform: enums.Platform, profile_id: str,
                                             session=None) -> dict:
    if platform == enums.Platform.INSTAGRAM:
        profile = await get(session, InstagramAccount, profile_id=profile_id)
        if not profile or profile.updated_at is None or profile.updated_at < datetime.now() - timedelta(days=7):
            try:
                await refresh_profile(profile_id, None, session=session)
                profile = await get(session, InstagramAccount, profile_id=profile_id)
            except Exception as e:
                logger.error(e)
                return error("Profile not found - %s" % profile_id, 404)
        response = {
            "platform": platform,
            "handle": profile.handle,
            "profile": transform_insta_profile(profile, recent_posts=None)
        }
        return success("Profile Retrieved", profile=response)
    return error("Platform not supported - %s" % platform, 400)


@app.get("/posts/{platform}/byshortcode/{shortcode}")
@sessionize_api
async def get_post_details_for_shortcode(request: Request, platform: enums.Platform, shortcode: str, handle: str,
                                         third_party_fallback: bool = False, session=None) -> dict:
    if platform == enums.Platform.INSTAGRAM:
        post = await get(session, InstagramPost, shortcode=shortcode)
        if not post or post.post_id is None or "_" in post.post_id:
            try:
                await refresh_profile_by_handle_from_graphapi(handle, session=session)
                post = await get(session, InstagramPost, shortcode=shortcode)
                if not post:
                    return error("Post not found - %s" % shortcode, 404)
                await session.refresh(post)
            except BeatException as e:
                if third_party_fallback:
                    logger.error(
                        f"Failed to refresh profile - falling back to lama for post {shortcode} by {handle} - {e}")
                    return await lama_fallback(handle, platform, shortcode, session=session)
                return error(e.message, e.status_code)
            except Exception as e:
                if third_party_fallback:
                    logger.error(
                        f"Failed to refresh profile - falling back to lama for post {shortcode} by {handle} - {e}")
                    return await lama_fallback(handle, platform, shortcode, session=session)
                return error(str(e), 500)

        try:
            await refresh_post_insights(post.shortcode, handle, post.post_id, post.post_type, session=session)
        except BeatException as e:
            if third_party_fallback:
                logger.error(f"Failed to refresh profile - falling back to lama for post {shortcode} by {handle} - {e}")
                return await lama_fallback(handle, platform, shortcode, session=session)
            return error(e.message, e.status_code)
        except Exception as e:
            if third_party_fallback:
                logger.error(f"Failed to refresh profile - falling back to lama for post {shortcode} by {handle} - {e}")
                return await lama_fallback(handle, platform, shortcode, session=session)
            return error(str(e), 500)

        response = {
            "platform": platform,
            "handle": handle,
            "data_provider": "GRAPHAPI",
            **transform_insta_post(post),
        }
        return success("Post Retrieved", post=response)
    return error("Platform not supported - %s" % platform, 400)


@app.get("/token/validate")
async def token_validate(access_token: str, social_user_id: str, handle: str, platform: str, type: str) -> dict:
    validator = CredentialValidator()
    cred = CredentialModel(
        id=0,
        source=type,
        handle=handle,
        credentials={"user_id": social_user_id, "token": access_token})
    try:
        await validator.validate(cred)
    except TokenValidationFalied:
        return error("Token is invalid", 190)
    except DataAccessExpired:
        return error("Data access expired", 500)
    else:
        return success("Token is valid", code=200)


@app.post("/scrape_request_log/flow/{flow}")
@sessionize_scl_api
async def create_scrape_request_log(request: Request, body: ScrapeRequest, session=None) -> dict:

    account_id = None
    if 'x-bb-account-id' in request.headers:
        account_id = request.headers['x-bb-account-id']
    scrape_request = ScrapeRequestLog(
        params=body.params,
        platform=body.platform,
        account_id=account_id,
        status=enums.ScrapeLogRequestStatus.PENDING.name,
        flow=body.flow,
        priority=1,
        created_at=datetime.now(),
        expires_at=datetime.now() + timedelta(hours=1),
        picked_at=None,
        retry_count=1)
    try:
        session.add(scrape_request)
        await session.flush()
        
        now = datetime.now()
        body.status="PENDING"
        body.retry_count=1
        body.event_timestamp = now.strftime("%Y-%m-%d %H:%M:%S")
        body.picked_at = now.strftime("%Y-%m-%d %H:%M:%S")
        body.expires_at = now + timedelta(hours=1)
        body.priority=1
        await make_scrape_request_log_event(str(scrape_request.id), body)
        
    except Exception as e:
        return error("Something went wrong - %s" % e, 400)
    else:
        return success("Scrape request log added", code=200, scrape_id=scrape_request.id)

@app.post("/scrape_request_log/flow/update/{scrape_id}")
@sessionize_scl_api
async def update_scrape_request_log(request: Request, scrape_id: int, body: dict, session=None) -> dict:

    now = datetime.now()
    
    picked_at = None 
    if body["status"] != enums.ScrapeLogRequestStatus.PENDING.name:
        picked_at = now.strftime("%Y-%m-%d %H:%M:%S")
    
    event_body = ScrapeRequest(
            flow=body["flow"],
            platform=body["platform"],
            params=body["params"],
            retry_count=body["retry_count"],
            status=body["status"],
            event_timestamp=now.strftime("%Y-%m-%d %H:%M:%S"),
            picked_at=picked_at,
            expires_at= now + timedelta(hours=1),
            priority=1
        )
    
    update_query = update(ScrapeRequestLog).values({"retry_count": body["retry_count"], "status": body["status"]}).where(and_(ScrapeRequestLog.id == scrape_id))
    await session.execute(update_query)
    await make_scrape_request_log_event(str(scrape_id), event_body)
    

@app.get("/scrape_request_log/flow/{id}")
@sessionize_scl_api
async def get_scrape_data(request: Request, id: int, session=None) -> dict:
    try:
        scrape_log = await get(session, ScrapeRequestLog, id=id)

        if scrape_log is not None:
            response = {
                "scrape_id": scrape_log.id,
                "data": scrape_log.data,
                "status": scrape_log.status,
                "params": scrape_log.params,
                "created_at": scrape_log.created_at,
                "flow": scrape_log.flow,
                "scraped_at": scrape_log.scraped_at,
                "platform": scrape_log.platform,
            }
            logger.debug(response)
            return success("Scrape Log Data Retrieved", post=response)
        return error(f"scrape log not found for scrape_id - {id}", 404)
    except Exception as e:
        return error(f"scrape log not found - {e}", 404)


@app.get("/list_scrape_data/{account_id}")
@sessionize_scl_api
async def list_scrape_data(request: Request, account_id: str, session=None) -> dict:
    size = request.query_params.get('size')
    size = 10 if size is None else size
    try:
        scrape_logs = await session.execute(select(ScrapeRequestLog).filter(
            ScrapeRequestLog.account_id == account_id
        ).order_by(ScrapeRequestLog.created_at.desc()).limit(size))
        scrape_logs_list = []
        for scrape_temp in scrape_logs:
            for scrape_log in scrape_temp:
                scrape_logs_list.append({
                    "scrape_id": scrape_log.id,
                    "account_id": scrape_log.account_id,
                    "data": scrape_log.data,
                    "status": scrape_log.status,
                    "params": scrape_log.params,
                    "created_at": scrape_log.created_at,
                    "flow": scrape_log.flow,
                    "scraped_at": scrape_log.scraped_at,
                    "platform": scrape_log.platform,
                })
        response = {
            "account_id": account_id,
            "scrape_list": scrape_logs_list
        }
        return response
    except Exception as e:
        return error(f"scrape log not found for account id - {e}", 404)


@app.get("/youtube/channel/byhandle/{handle}")
async def fetch_channel_id_by_handle(handle: str):
    try:
        channel_id = await fetch_yt_channel_id_from_handle(handle)
    except Exception as e:
        return error("Something went wrong - %s" % e, 400)

    if channel_id:
        response = {
            "handle": handle,
            "channel_id": channel_id
        }
        return success("Channel ID Retrieved", response=response)
    else:
        return error(f"Failed to fetch Channel ID for - {handle}", 400)


@app.get("/recent/posts/{platform}/byprofileid/{profile_id}")
@sessionize_api
async def fetch_recent_posts_by_profile_id(request: Request,
                                           platform: enums.Platform,
                                           profile_id: str,
                                           post_types=None,
                                           session=None) -> dict:
    if platform == enums.Platform.INSTAGRAM:
        try:
            if post_types:
                if post_types == ['reels']:
                    result = await session.execute(select(InstagramPost).filter(
                        and_(
                            and_(InstagramPost.plays >= 0),
                            and_(InstagramPost.post_type.in_(post_types),
                                 InstagramPost.profile_id == profile_id)))
                                                   .order_by(InstagramPost.publish_time.desc()).limit(12))
                else:
                    result = await session.execute(select(InstagramPost).filter(and_(
                        and_(InstagramPost.post_type.in_(post_types),
                             InstagramPost.profile_id == profile_id))).order_by(
                        InstagramPost.publish_time.desc()).limit(12))
            else:
                result = await session.execute(select(InstagramPost).filter(and_(
                    and_(InstagramPost.thumbnail_url != None, InstagramPost.thumbnail_url != "None",
                         InstagramPost.thumbnail_url != ' '),
                    and_(InstagramPost.post_type != 'story', InstagramPost.profile_id == profile_id))).order_by(
                    InstagramPost.publish_time.desc()).limit(12))
        except Exception:
            return error("Failed to fetch recent posts for profile_id - %s" % profile_id, 500)

        recent_posts = []
        posts = [r[0] for r in result]

        if not posts:
            return error("Failed to fetch recent posts for profile_id - %s" % profile_id, 500)

        try:
            account = await session.execute(select(InstagramAccount).filter(InstagramAccount.profile_id == profile_id))
        except Exception:
            return error("Failed to fetch Instagram account for handle - %s" % profile_id, 500)

        account = [a[0] for a in account]
        account = account[0] if account else None
        try:
            asset_result = await session.execute(
                select(AssetLog).filter(AssetLog.entity_id.in_(p.shortcode for p in posts if p is not None)))
        except Exception as e:
            return error(f"Failed to fetch asset logs for posts shortcode {e}", 500)

        asset_dict = {asset[0].entity_id: asset[0].asset_url for asset in asset_result}
        for post in posts:
            if post.shortcode in asset_dict:
                recent_post = transform_insta_recent_posts(post, account)
                if recent_post:
                    recent_post['dimensions'][
                        'thumbnail_url'] = f"https://d24w28i6lzk071.cloudfront.net/{asset_dict[post.shortcode]}"
                recent_posts.append(recent_post)
            else:
                recent_posts.append(transform_insta_recent_posts(post, account))

        response = {
            "platform": platform,
            "profile_id": profile_id,
            "recent_posts": recent_posts
        }
        return success("Recent Posts Retrieved", data=response)
    return error("Platform not supported - %s" % platform, 400)


async def lama_fallback(handle: str, platform: str, shortcode: str, session=None) -> dict:
    # Call lama for fallback
    lama = Lama()
    ctx = RequestContext(CredentialModel(0, "", "", {"x-access-key": "AUIiOivV7oO9vqrv3nnDdL5KYlPzcgu8"}))
    try:
        resp_dict = await lama.fetch_post_by_shortcode(ctx, shortcode)
        log = lama.parse_post_by_shortcode(resp_dict)
        now = datetime.now()
        context = Context(now)
        post = await upsert_insta_post(context, log, session=session)
        response = {
            "platform": platform,
            "handle": post.handle,
            "data_provider": "LAMA",
            **transform_insta_post(post),
        }
        post_handle = post.handle
        if handle.lower() != post_handle.lower():
            logger.error(f"Handle mismatch {handle.lower()} {post_handle.lower()}")
            return error("Post not found - %s" % shortcode, 404)
        return success("Post Retrieved", post=response)
    except Exception as e:
        logger.error(e)
        return error("Post not found - %s" % shortcode, 404)


@app.get("/heartbeat")
async def getHealth(response: Response, settings: config.Settings = Depends(get_settings)):
    if settings.heartbeat==True:
        response.status_code = status.HTTP_200_OK
    else:
        response.status_code=status.HTTP_410_GONE
    return response


@app.put("/heartbeat")
async def setHealth(beat: bool, response: Response, settings: config.Settings = Depends(get_settings)):
    settings.heartbeat = beat
    if beat:
        response.status_code = status.HTTP_200_OK
    else:
        response.status_code = status.HTTP_410_GONE
    return response


@app.get("/posts/{platform}/{post_type}/{shortcode}")
@sessionize_api
async def get_post_details_using_shortcode(platform: enums.Platform, post_type: str, shortcode: str, request: Request,
                                           profile_id: str,
                                           session=None) -> dict:
    if platform == enums.Platform.INSTAGRAM:
        if post_type == "story":
            try:
                await refresh_story_posts_by_profile_id(profile_id, handle="", session=session)
                story = await get(session, InstagramPost, shortcode=shortcode)
                if story is not None:
                    response = {
                        "platform": platform,
                        "handle": story.handle,
                        "data_provider": "GRAPHAPI",
                        **transform_insta_post(story),
                    }
                    logger.debug(response)
                    return success("Stories Retrieved", post=response)
                return error(f"story not found for shortcode - {shortcode}", 404)
            except Exception as e:
                return error(f"Story not found - {e}", 404)

        else:
            try:
                await refresh_post_by_shortcode(shortcode=shortcode)
                post = await get(session, InstagramPost, shortcode=shortcode)
                if post is not None:
                    response = {
                        "platform": platform,
                        **transform_insta_post(post),
                    }
                    logger.debug(response)
                    message = post_type + " Retrieved"
                    return success(message, post=response)
                return error(f"{post_type} not found for shortcode - {shortcode}", 404)
            except Exception as e:
                return error(f"{post_type} not found - {e}", 404)
    return error(f"Platform not supported - {platform}", 400)


@app.get("/profiles/{platform}/byid/{id}")
@sessionize_api
async def get_profile_details_for_id(request: Request, platform: enums.Platform, id: str,
                                     third_party_fallback: bool = False, session=None) -> dict:
    if platform == enums.Platform.INSTAGRAM:
        profile = await get(session, InstagramAccount, profile_id=id)
        if not profile or profile.updated_at is None or profile.updated_at < datetime.now() - timedelta(days=1):
            try:
                redis = AsyncRedis.from_url(os.environ["REDIS_URL"])
                global_limit_day = RateSpec(requests=20000, seconds=86400)
                global_limit_minute = RateSpec(requests=60, seconds=60)
                id_limit = RateSpec(requests=1, seconds=1)
                async with RateLimiter(
                        unique_key="refresh_profile_insta_by_id_daily",
                        backend=redis,
                        cache_prefix="beat_server_",
                        rate_spec=global_limit_day):
                    async with RateLimiter(
                            unique_key="refresh_profile_insta_by_id_per_minute",
                            backend=redis,
                            cache_prefix="beat_server_",
                            rate_spec=global_limit_minute):
                        async with RateLimiter(
                                unique_key="refresh_profile_insta_per_id_" + id,
                                backend=redis,
                                cache_prefix="beat_server_",
                                rate_spec=id_limit):
                            await refresh_profile_by_profile_id(id, session=session)
                profile = await get(session, InstagramAccount, profile_id=id)
            except RateLimitError as e:
                logger.error(e)
                return error("Concurrency Limit Exceeded", 2)
            except Exception as e:
                logger.error(e)
                return error(f"Error - {e}", 404)
        if profile is not None:
            response = {
                "platform": platform,
                **transform_insta_profile(profile, recent_posts=None)
            }
            logger.debug(response)
            return success("Profile Retrieved", profile=response)
        return error(f"profile not found for id - {id}", 404)
    elif platform == enums.Platform.YOUTUBE:
        profile = await get(session, YoutubeAccount, channel_id=id)
        if not profile or profile.updated_at is None or profile.updated_at < datetime.now() - timedelta(days=1):
            try:
                redis = AsyncRedis.from_url(os.environ["REDIS_URL"])
                global_limit_day = RateSpec(requests=20000, seconds=86400)
                global_limit_minute = RateSpec(requests=60, seconds=60)
                id_limit = RateSpec(requests=1, seconds=1)
                async with RateLimiter(
                        unique_key="refresh_profile_yt_by_id_daily",
                        backend=redis,
                        cache_prefix="beat_server_",
                        rate_spec=global_limit_day):
                    async with RateLimiter(
                            unique_key="refresh_profile_yt_by_id_per_minute",
                            backend=redis,
                            cache_prefix="beat_server_",
                            rate_spec=global_limit_minute):
                        async with RateLimiter(
                                unique_key="refresh_profile_yt_per_id_" + id,
                                backend=redis,
                                cache_prefix="beat_server_",
                                rate_spec=id_limit):
                            await refresh_yt_profiles([id], session=session)
                profile = await get(session, YoutubeAccount, channel_id=id)
            except RateLimitError as e:
                logger.error(e)
                return error("Concurrency Limit Exceeded", 2)
            except Exception as e:
                logger.error(e)
                return error(f"Error - {e}", 404)
        if profile is not None:
            response = {
                "platform": platform,
                **transform_yt_profile(profile)
            }
            logger.debug(response)
            return success("Profile Retrieved", profile=response)
        return error(f"profile not found for id - {id}", 404)
    return error("Platform not supported - %s" % platform, 400)


def transform_insta_post(entity: Optional[InstagramPost]) -> dict:
    if entity is None:
        return {}
    return {
        "shortcode": entity.shortcode,
        "handle": entity.handle,
        "post_id": entity.post_id,
        "post_type": entity.post_type,
        "content_type": entity.content_type,
        "thumbnail": entity.thumbnail_url,
        "display_url": entity.display_url,
        "story_links": entity.story_links,
        "story_hashtags": entity.story_hashtags,
        "publish_time": entity.publish_time,
        "updated_at": entity.updated_at,
        "metrics": {
            "reach": entity.reach,
            "views": entity.views,
            "likes": entity.likes,
            "play_count": entity.plays,
            "comments": entity.comments
        }
    }


@app.get("/profiles/INSTAGRAM/byhandle/{handle}/insights")
@sessionize_api
async def get_insights_using_handle(handle: str, request: Request, token: str, user_id: str, 
                                         session=None) -> dict:
    try:
        await update_token(session, token, handle, user_id)
        profile_data, recent_posts = await fetch_profile_from_graphapi(handle, session=session)
        profile_id = profile_data.profile_id
        account, posts = await process_profile_data(profile_id, profile_data, recent_posts, session=session)
        transformed_posts = []
        for post in posts:
            recent_post1 = transform_insta_recent_posts(post, account)
            transformed_posts.append(recent_post1)

        if profile_data is not None:
            response = {
                "handle": handle,
                **transform_insights(account, transformed_posts)
            }
            logger.debug(response)
            return success("Profile Retrieved", profile=response)
    except Exception as e:
        return error(f"Error - {e}", 404)
    
@app.get("/profiles/INSTAGRAM/byhandle/{handle}/audienceinsights")
@sessionize_api
async def get_audience_insights_using_handle(handle: str, request: Request, token: str, user_id: str,
                                         session=None) -> dict:
    try:
        await update_token(session, token, handle, user_id)
        audience_data = await fetch_profile_insights_from_graphapi(handle, session=session)
        if audience_data is not None:
            response = {
                "handle": handle,
                **transform_audience_insights(audience_data)
            }
            logger.debug(response)
            return success("Profile Retrieved", profile=response)
    except Exception as e:
        return error(f"Error: {e}", 404)

async def update_token(session, token: str, handle: str, user_id: str) -> None:
    source = "graphapi"
    load_dotenv()
    _result = await session.execute(
        select(Credential).filter(
            and_(Credential.handle == handle, Credential.source == source)))
    result = _result.first()
    if result:
        entity: Credential = result[0]
        if entity.data_access_expired == True:
            credential = {"token": token, "user_id": user_id}
            update_query = update(Credential).values(
                {"credentials": credential, "enabled": True, "data_access_expired": False}).where(
                and_(Credential.handle == handle, Credential.source == source))
            await session.execute(update_query)
    else:
        credential = {"token": token, "user_id": user_id}
        await CredentialManager.insert_creds(source, credential, handle, session=session)

def transform_insights(account: InstagramAccount, posts: List[InstagramPost]) -> dict:
    post_metrics = get_recent_post_metrics(account, posts)
    return {
        "followers": account.followers,
        "following": account.following,
        "media_count": account.media_count,
        "post_metrics": post_metrics
    }

def transform_audience_insights(audience_data: InstagramProfileLog) -> dict:
    audience_city_value, audience_gender_age_value = None, None
    for dimension in audience_data.dimensions:
        if dimension.key == 'audience_city':
            audience_city_value = dimension.value
        if dimension.key == 'audience_gender_age':
            audience_gender_age_value = dimension.value
    return {
        "audience_data": {
            "city": audience_city_value,
            "gender_age": audience_gender_age_value,
        }
    }

def error(message: str, code: int) -> dict:
    return {
        "status": {
            "type": "ERROR",
            "message": message,
            "code": code
        }
    }


def success(message: str, **kwargs: any) -> dict:
    resp = {
        "status": {
            "type": "SUCCESS",
            "message": message
        }
    }
    for k, v in kwargs.items():
        resp[k] = v
    return resp


# def write_pid():
#     with open('beat_server.pid', 'w', encoding='utf-8') as f:
#         f.write(str(os.getpid()))

if __name__ == "__main__":
    # write_pid()
    uvicorn.run(app, host="0.0.0.0", port=8000)
