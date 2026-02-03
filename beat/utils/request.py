import asyncio
import datetime
import os
import socket
import time
import urllib.parse
import uuid

import asks
import asyncio_redis_rate_limit
from asyncio_redis_rate_limit import RateLimiter, RateSpec
from loguru import logger
from redis.asyncio import Redis as AsyncRedis

from core.enums import enums
from core.amqp.amqp import publish
from core.entities.entities import OrderLogs, PostLog, ProfileLog, ProfileRelationshipLog, SentimentLog, \
    PostActivityLog
from core.models.models import ScrapeRequest
from youtube.models.models import YoutubeActivityLog, YoutubeProfileRelationshipLog

youtube138 = 'youtube138'
youtubev31 = 'youtubev31'
insta_best_perf = 'insta-best-performance'
insta_jt = 'instagram-scraper2'
arraybobo = 'instagram-scraper-2022'
youtubev311 = 'youtubev311'
rocketapi = 'rocketapi'

source_specs = {
    youtube138: RateSpec(requests=850, seconds=60),
    insta_best_perf: RateSpec(requests=2, seconds=1),
    insta_jt: RateSpec(requests=5, seconds=1),
    arraybobo: RateSpec(requests=100, seconds=30),
    youtubev31: RateSpec(requests=500, seconds=60),
    youtubev311: RateSpec(requests=3, seconds=1),
    rocketapi: RateSpec(requests=100, seconds=30),
}


async def emit_trace_log_event(method, address, resp, time_taken, kwargs):
    headers = {}
    if 'headers' in kwargs:
        headers = kwargs['headers']
    url = address
    if 'params' in kwargs:
        params = kwargs['params']
        url = f"{url}?{urllib.parse.urlencode(params)}"
    status_code = 0
    if resp.status_code:
        status_code = int(resp.status_code)
    payload = {
        "id": str(uuid.uuid4()),
        "hostname": socket.gethostname(),
        "serviceName": "beat",
        "timestamp": datetime.datetime.now().isoformat(),
        "timeTaken": int(time_taken),
        "method": method,
        "uri": url,
        "headers": headers,
        "responseStatus": status_code,
        "remoteAddress": address,
        "onlyRequest": False,
        "requestBody": "",
        "responseBody": "",
    }
    publish(payload, "identity.dx", "trace_log")


async def make_request(method: str, url: str, **kwargs: any) -> any:
    if 'instagram-api-cheap-best-performance.p.rapidapi.com' in url:
        return await make_request_limited(insta_best_perf, method, url, **kwargs)
    elif 'youtube138.p.rapidapi.com' in url:
        return await make_request_limited(youtube138, method, url, **kwargs)
    elif 'instagram-scraper22' in url:
        return await make_request_limited(insta_jt, method, url, **kwargs)
    elif 'instagram-scraper-2022' in url:
        return await make_request_limited(arraybobo, method, url, **kwargs)
    elif 'youtube-v31.p.rapidapi.com' in url:
        return await make_request_limited(youtubev31, method, url, **kwargs)
    elif 'youtube-v311.p.rapidapi.com' in url:
        return await make_request_limited(youtubev311, method, url, **kwargs)
    elif 'v1.rocketapi.io' in url:
        return await make_request_limited(rocketapi, method, url, **kwargs)
    logger.debug("Making Request :: %s %s %s" % (method, url, kwargs))
    start_time = time.perf_counter()
    resp = await asks.request(method, url, **kwargs)
    end_time = time.perf_counter()
    logger.debug("Completed Request :: %s %s %s %s" % (method, url, kwargs, resp.status_code))
    try:
        await emit_trace_log_event(method, url, resp, end_time - start_time, kwargs)
    except Exception as e:
        logger.error(f"Failed to emit trace log event - {e}")
    return resp


async def make_request_limited(source, method: str, url: str, **kwargs: any) -> any:
    redis = AsyncRedis.from_url(os.environ["REDIS_URL"])
    while True:
        try:
            async with RateLimiter(
                    unique_key=source,
                    backend=redis,
                    cache_prefix="beat_",
                    rate_spec=source_specs[source]):
                logger.debug("Making Request :: %s %s %s" % (method, url, kwargs))
                start_time = time.perf_counter()
                resp = await asks.request(method, url, **kwargs)
                end_time = time.perf_counter()
                logger.debug("Completed Request :: %s %s %s %s" % (method, url, kwargs, resp.status_code))
                try:
                    await emit_trace_log_event(method, url, resp, end_time - start_time, kwargs)
                except Exception as e:
                    logger.error(f"Failed to emit trace log event - {e}")
                return resp
        except asyncio_redis_rate_limit.RateLimitError as e:
            await asyncio.sleep(1)


async def make_scrape_request_log_event(scl_id: str, log: ScrapeRequest):
    payload = {
        "event_id": str(uuid.uuid4()),
        "scl_id": scl_id,
        "flow": log.flow,
        "platform": log.platform,
        "params": log.params,
        "status": log.status,
        "priority": log.priority,
        "reason": log.reason,
        "event_timestamp": log.event_timestamp,
        "picked_at": log.picked_at,
        "expires_at": log.expires_at,
        "retry_count": log.retry_count
    }
    publish(payload, "beat.dx", "scrape_request_log_events")


async def emit_post_log_event(log: PostLog):
    now = datetime.datetime.now()
    publish_time = None
    handle = None

    for dim in log.dimensions:
        if dim["key"] == "taken_at_timestamp":
            publish_time = dim["value"]

    if publish_time:
        publish_time = datetime.datetime.fromtimestamp(publish_time).strftime("%Y-%m-%d %H:%M:%S")

    if log.platform == 'YOUTUBE':
        handle = log.profile_id
    else:
        for dim in log.dimensions:
            if dim["key"] == "handle":
                handle = dim["value"]

    payload = {
        "event_id": str(uuid.uuid4()),
        "source": log.source,
        "platform": log.platform,
        "profile_id": log.profile_id,
        "handle": handle,
        "shortcode": log.platform_post_id,
        "publish_time": publish_time,
        "event_timestamp": now.strftime("%Y-%m-%d %H:%M:%S"),
        "metrics": {metrics["key"]: metrics["value"] for metrics in log.metrics},
        "dimensions": {dimension["key"]: dimension["value"] for dimension in log.dimensions}
    }
    publish(payload, "beat.dx", "post_log_events")


async def emit_post_activity_log_event(log: PostActivityLog):
    now = datetime.datetime.now()
    publish_time = None

    for dim in log.dimensions:
        if dim["key"] == "taken_at_timestamp":
            publish_time = dim["value"]

    if publish_time:
        publish_time = datetime.datetime.fromtimestamp(publish_time).strftime("%Y-%m-%d %H:%M:%S")

    payload = {
        "event_id": str(uuid.uuid4()),
        "source": log.source,
        "activity_type": log.activity_type,
        "platform": log.platform,
        "actor_profile_id": log.actor_profile_id,
        "shortcode": log.platform_post_id,
        "publish_time": publish_time,
        "event_timestamp": now.strftime("%Y-%m-%d %H:%M:%S"),
        "metrics": {metrics["key"]: metrics["value"] for metrics in log.metrics},
        "dimensions": {dimension["key"]: dimension["value"] for dimension in log.dimensions}
    }
    publish(payload, "beat.dx", "post_activity_log_events")


async def emit_sentiment_log_event(log: SentimentLog):
    now = datetime.datetime.now()

    payload = {
        "event_id": str(uuid.uuid4()),
        "source": log.source,
        "platform": log.platform,
        "shortcode": log.platform_post_id,
        "comment_id": log.comment_id,
        "event_timestamp": now.strftime("%Y-%m-%d %H:%M:%S"),
        "comment": log.comment,
        "sentiment": log.sentiment,
        "score": log.score,
        "metrics": {metrics["key"]: metrics["value"] for metrics in log.metrics},
        "dimensions": {dimension["key"]: dimension["value"] for dimension in log.dimensions}
    }
    publish(payload, "beat.dx", "sentiment_log_events")


async def emit_profile_log_event(log: ProfileLog):
    now = datetime.datetime.now()
    handle = None

    if log.platform == 'YOUTUBE':
        handle = log.profile_id
    else:
        for dim in log.dimensions:
            if dim["key"] == "handle":
                handle = dim["value"]

    payload = {
        "event_id": str(uuid.uuid4()),
        "source": log.source,
        "platform": log.platform,
        "profile_id": log.profile_id,
        "handle": handle,
        "event_timestamp": now.strftime("%Y-%m-%d %H:%M:%S"),
        "metrics": {metrics["key"]: metrics["value"] for metrics in log.metrics},
        "dimensions": {dimension["key"]: dimension["value"] for dimension in log.dimensions}
    }
    publish(payload, "beat.dx", "profile_log_events")


async def emit_profile_relationship_log_event(log: ProfileRelationshipLog):
    now = datetime.datetime.now()

    payload = {
        "event_id": str(uuid.uuid4()),
        "source": log.source,
        "relationship_type": log.relationship_type,
        "platform": log.platform,
        "source_profile_id": log.source_profile_id,
        "target_profile_id": log.target_profile_id,
        "event_timestamp": now.strftime("%Y-%m-%d %H:%M:%S"),
        "source_metrics": log.source_metrics,
        "target_metrics": log.target_metrics,
        "source_dimensions": log.source_dimensions,
        "target_dimensions": log.target_dimensions
    }
    publish(payload, "beat.dx", "profile_relationship_log_events")


async def emit_order_log_event(log: OrderLogs):
    now = datetime.datetime.now()
    store = None

    for dim in log.dimensions:
        if dim["key"] == "store":
            store = dim["value"]

    payload = {
        "event_id": str(uuid.uuid4()),
        "source": log.source,
        "platform": log.platform,
        "platform_order_id": log.platform_order_id,
        "store": store,
        "event_timestamp": now.strftime("%Y-%m-%d %H:%M:%S"),
        "metrics": {metrics["key"]: metrics["value"] for metrics in log.metrics},
        "dimensions": {dimension["key"]: dimension["value"] for dimension in log.dimensions}
    }
    publish(payload, "beat.dx", "order_log_events")


async def emit_youtube_activity_log_event(log: YoutubeActivityLog):
    now = datetime.datetime.now()

    payload = {
        "event_id": str(uuid.uuid4()),
        "activity_id": log.activity_id,
        "activity_type": log.activity_type,
        "source": log.source,
        "platform": 'YOUTUBE',
        "actor_channel_id": log.actor_channel_id,
        "metrics": {metrics["key"]: metrics["value"] for metrics in log.metrics},
        "dimensions": {dimension["key"]: dimension["value"] for dimension in log.dimensions},
        "event_timestamp": now.strftime("%Y-%m-%d %H:%M:%S"),
    }
    publish(payload, "beat.dx", "yt_activity_log_events")


async def emit_youtube_profile_relationship_log_event(log: YoutubeProfileRelationshipLog):
    now = datetime.datetime.now()

    payload = {
        "event_id": str(uuid.uuid4()),
        "relationship_id": log.relationship_id,
        "relationship_type": log.relationship_type,
        "source": log.source,
        "platform": 'YOUTUBE',
        "source_channel_id": log.source_channel_id,
        "target_channel_id": log.target_channel_id,
        "source_metrics": {metrics["key"]: metrics["value"] for metrics in log.source_metrics},
        "target_metrics": {metrics["key"]: metrics["value"] for metrics in log.target_metrics},
        "source_dimensions": {dimension["key"]: dimension["value"] for dimension in log.source_dimensions},
        "target_dimensions": {dimension["key"]: dimension["value"] for dimension in log.target_dimensions},
        "subscribed_on": log.subscribed_on,
        "event_timestamp": now.strftime("%Y-%m-%d %H:%M:%S"),
    }
    publish(payload, "beat.dx", "yt_profile_relationship_log_events")


async def make_scrape_log_event(log_type: str, scrape_log: any):
    if log_type == "post_log":
        try:
            await emit_post_log_event(scrape_log)
        except Exception as e:
            logger.error(f"Failed to emit post log event - {e}")
    elif log_type == "profile_log":
        try:
            await emit_profile_log_event(scrape_log)
        except Exception as e:
            logger.error(f"Failed to emit profile log event - {e}")
    elif log_type == "profile_relationship_log":
        try:
            await emit_profile_relationship_log_event(scrape_log)
        except Exception as e:
            logger.error(f"Failed to emit profile relationship log event - {e}")
    elif log_type == "order_log":
        try:
            await emit_order_log_event(scrape_log)
        except Exception as e:
            logger.error(f"Failed to emit order log event - {e}")
    elif log_type == "sentiment_log":
        try:
            await emit_sentiment_log_event(scrape_log)
        except Exception as e:
            logger.error(f"Failed to emit sentiment log event - {e}")
    elif log_type == "post_activity_log":
        try:
            await emit_post_activity_log_event(scrape_log)
        except Exception as e:
            logger.error(f"Failed to emit post activity log event - {e}")
    elif log_type == "yt_activity_log":
        try:
            await emit_youtube_activity_log_event(scrape_log)
        except Exception as e:
            logger.error(f"Failed to emit activity log event - {e}")
    elif log_type == "yt_profile_relationship_logs":
        try:
            await emit_youtube_profile_relationship_log_event(scrape_log)
        except Exception as e:
            logger.error(f"Failed to emit youtube profile relationship log event - {e}")

