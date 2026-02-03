import json
import isodate

import dateutil.parser
from core.models.models import Dimension
from utils.exceptions import PostScrapingFailed

from utils.getter import safe_metric, safe_dimension, safe_get, dimension
from youtube.models.models import YoutubeProfileLog, YoutubePostLog
from youtube.metric_dim_store import POST_TYPE, TOPIC, VIDEO_COUNT, VIEW_COUNT, HIDDEN_SUBSCRIBER_COUNT, \
    SUBSCRIBER_COUNT, TITLE, DESCRIPTION, UPLOADS_PLAYLIST, COUNTRY, THUMBNAILS, THUMBNAIL, PUBLISH_TIME, CUSTOM_URL, \
    LIKES, FAVOURITE_COUNT, COMMENTS, TAGS, CATEGORY_ID, AUDIO_LANGUAGE, IS_LICENSED, CHANNEL_TITLE, \
    DURATION, CONTENT_RATING, CATEGORY

SOURCE = "ytapi"


def transform_profile(s: dict) -> YoutubeProfileLog:
    channel_id = s["id"]
    snippet = s["snippet"]
    content_details = s["contentDetails"]
    statistics = s["statistics"]
    metrics = [
        safe_metric(statistics, 'videoCount', VIDEO_COUNT, type=int),
        safe_metric(statistics, 'subscriberCount', SUBSCRIBER_COUNT, type=int),
        safe_metric(statistics, 'hiddenSubscriberCount', HIDDEN_SUBSCRIBER_COUNT, type=int),
        safe_metric(statistics, 'viewCount', VIEW_COUNT, type=int),
    ]

    dimensions = [
        safe_dimension(snippet, 'title', TITLE),
        safe_dimension(snippet, 'description', DESCRIPTION),
        safe_dimension(snippet, 'customUrl', CUSTOM_URL),
        dimension(int(dateutil.parser.isoparse(safe_get(snippet, 'publishedAt')).strftime('%s')), PUBLISH_TIME),
        dimension(json.dumps(safe_get(snippet, 'thumbnails')), THUMBNAILS),
        safe_dimension(snippet, 'thumbnails.default.url', THUMBNAIL),
        safe_dimension(snippet, 'country', COUNTRY),
        safe_dimension(content_details, 'relatedPlaylists.uploads', UPLOADS_PLAYLIST),
    ]
    return YoutubeProfileLog(channel_id, SOURCE, dimensions=dimensions, metrics=metrics)


def transform_post(s: dict, category=None, topic=None) -> YoutubePostLog:
    snippet = s["snippet"]
    shortcode = s["id"]
    content_details = s["contentDetails"]
    statistics = s["statistics"]
    channel_id = snippet["channelId"]
    metrics = [
        safe_metric(statistics, 'commentCount', COMMENTS, type=int),
        safe_metric(statistics, 'favoriteCount', FAVOURITE_COUNT, type=int),
        safe_metric(statistics, 'likeCount', LIKES, type=int),
        safe_metric(statistics, 'viewCount', VIEW_COUNT, type=int),
    ]

    dimensions = [
        safe_dimension(snippet, 'title', TITLE),
        safe_dimension(snippet, 'channelTitle', CHANNEL_TITLE),
        safe_dimension(snippet, 'description', DESCRIPTION),
        safe_dimension(content_details, 'duration', DURATION),
        dimension(int(dateutil.parser.isoparse(safe_get(snippet, 'publishedAt')).strftime('%s')), PUBLISH_TIME),
        dimension(json.dumps(safe_get(snippet, 'thumbnails')), THUMBNAILS),
        safe_dimension(snippet, 'thumbnails.default.url', THUMBNAIL),
        dimension(json.dumps(safe_get(snippet, 'tags')), TAGS),
        safe_dimension(snippet, 'categoryId', CATEGORY_ID),
        safe_dimension(snippet, 'defaultAudioLanguage', AUDIO_LANGUAGE),
        safe_dimension(content_details, 'licensedContent', IS_LICENSED),
        dimension(json.dumps(safe_get(content_details, 'contentRating')), CONTENT_RATING),
    ]

    if category:
        dimensions.append(Dimension(CATEGORY, category))

    if topic:
        dimensions.append(Dimension(TOPIC, topic))

    duration_str = safe_get(content_details, 'duration', type=str)
    if duration_str:
        duration = isodate.parse_duration(duration_str)
        seconds = duration.total_seconds()
        if seconds > 60:
            dimensions.append(Dimension(POST_TYPE, 'video'))

    return YoutubePostLog(channel_id, shortcode, SOURCE, dimensions=dimensions, metrics=metrics)


def transform_playlist_for_post_ids(s: dict) -> str:
    content_details = s["contentDetails"]
    return content_details['videoId']


def transform_channel_for_post_ids(s: dict) -> str:
    if 'videoId' not in s['id']:
        raise PostScrapingFailed("Failed to fetch post_id")
    return s['id']['videoId']
