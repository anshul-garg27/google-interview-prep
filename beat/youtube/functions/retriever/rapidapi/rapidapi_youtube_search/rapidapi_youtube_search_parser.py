import json

import dateutil
import loguru

from utils.getter import safe_dimension, dimension, safe_get
from youtube.metric_dim_store import TITLE, CHANNEL_TITLE, DESCRIPTION, THUMBNAILS, THUMBNAIL, PUBLISH_TIME, TEXT, \
    CHANNEL_ID, LIKES, COMMENT_ID
from youtube.models.models import YoutubePostLog, YoutubePostActivityLog

SOURCE = "rapidapi-yt-search"


def transform_post_ids_by_search(s: dict) -> str:
    if 'videoId' in s['id']:
        return s['id']['videoId']


def transform_channel_ids_by_search(s: dict) -> str:
    if 'channelId' in s['snippet']:
        return s['snippet']['channelId']


def transform_post_log_by_search(s: dict) -> YoutubePostLog:
    snippet = s["snippet"]
    shortcode = s['id']['videoId']
    channel_id = snippet["channelId"]
    dimensions = [
        safe_dimension(snippet, 'title', TITLE),
        safe_dimension(snippet, 'channelTitle', CHANNEL_TITLE),
        safe_dimension(snippet, 'description', DESCRIPTION),
        dimension(int(dateutil.parser.isoparse(safe_get(snippet, 'publishedAt')).strftime('%s')), PUBLISH_TIME),
        dimension(json.dumps(safe_get(snippet, 'thumbnails')), THUMBNAILS),
        safe_dimension(snippet, 'thumbnails.default.url', THUMBNAIL)

    ]
    metrics = []
    return YoutubePostLog(channel_id, shortcode, SOURCE, dimensions=dimensions, metrics=metrics)


def transform_post_comments(shortcode: str, s: dict) -> YoutubePostActivityLog:
    snippet = s["snippet"]
    shortcode = snippet["videoId"]
    topLevelComment = snippet['topLevelComment']
    dimensions = [
        # safe_dimension(snippet, 'title', TITLE),
        safe_dimension(topLevelComment, 'snippet.textDisplay', TEXT),
        # safe_dimension(snippet, 'snippet.textOriginal', ),
        safe_dimension(topLevelComment, 'snippet.authorDisplayName', CHANNEL_TITLE),
        safe_dimension(topLevelComment, 'snippet.authorProfileImageUrl', THUMBNAIL),
        dimension(int(dateutil.parser.isoparse(safe_get(topLevelComment, 'snippet.publishedAt')).strftime('%s')), PUBLISH_TIME),
        safe_dimension(topLevelComment, 'id', COMMENT_ID)
    ]
    metrics = [
        safe_dimension(topLevelComment, 'snippet.likeCount', LIKES),
    ]

    activity_type = "COMMENT"
    actor_profile_id = safe_get(topLevelComment, 'snippet.authorChannelId.value')
    log = YoutubePostActivityLog(
        activity_type=activity_type,
        shortcode=shortcode,
        actor_profile_id=actor_profile_id,
        dimensions=dimensions,
        metrics=metrics,
        source=SOURCE,
    )
    return log