import json
import datetime
from core.models.models import Dimension

from utils.getter import safe_metric, safe_dimension, safe_get, dimension
from youtube.models.models import YoutubeProfileLog, YoutubePostLog
from youtube.metric_dim_store import TOPIC, VIDEO_COUNT, VIEW_COUNT, SUBSCRIBER_COUNT, TITLE, EXTERNAL_URLS, \
    DESCRIPTION, COUNTRY, THUMBNAILS, THUMBNAIL, PUBLISH_TIME, CUSTOM_URL, LIKES, CONTENT_RATING, CATEGORY,\
    COMMENTS, TAGS, CATEGORY_ID, AUDIO_LANGUAGE, IS_LICENSED, CHANNEL_TITLE, DURATION, IS_VERIFIED, PLAYLIST_ID, POST_TYPE

SOURCE = "rapidapi-yt"

def _convert_profile_time_to_epoch(value: str) -> int:
    if not value:
        return None
    datetime_object = datetime.datetime.strptime(value, "%Y-%m-%d")
    epoch = int(datetime_object.timestamp())
    return epoch

def _convert_post_time_to_epoch(value: str) -> int:
    if not value:
        return None
    datetime_object = datetime.datetime.strptime(value, "%Y-%m-%dT%H:%M:%S%z")
    epoch = int(datetime_object.timestamp())
    return epoch

def _convert_seconds_to_iso(time_in_seconds):
    duration = datetime.timedelta(seconds=time_in_seconds)
    hours = duration // datetime.timedelta(hours=1)
    minutes = (duration // datetime.timedelta(minutes=1)) % 60
    seconds = (duration // datetime.timedelta(seconds=1)) % 60
    iso_duration = f"PT{hours}H{minutes:02}M{seconds:02}S"

    return iso_duration

def transform_profile(s: dict) -> YoutubeProfileLog:
    channel_id = s["channelId"]
    statistics = s["stats"]
    metrics = [
        safe_metric(statistics, 'videos', VIDEO_COUNT, type=int),
        safe_metric(statistics, 'subscribers', SUBSCRIBER_COUNT, type=int),
        safe_metric(statistics, 'views', VIEW_COUNT, type=int),
    ]

    dimensions = [
        safe_dimension(s, 'title', TITLE),
        safe_dimension(s, 'description', DESCRIPTION),
        safe_dimension(s, 'username', CUSTOM_URL),
        dimension(json.dumps(safe_get(s, 'links')), EXTERNAL_URLS),
        dimension(json.dumps(safe_get(s, 'banner.desktop')), THUMBNAILS),
        safe_dimension(s, 'banner.desktop.5.url', THUMBNAIL),
        safe_dimension(s, 'country', COUNTRY),
        dimension(_convert_profile_time_to_epoch(safe_get(s, 'joinedDate')), PUBLISH_TIME)
    ]

    if "badges" in s and s["badges"] and s["badges"][0]['text'] == "Verified":
        dimensions.append(Dimension(IS_VERIFIED, True))
    else:
        dimensions.append(Dimension(IS_VERIFIED, False))

    return YoutubeProfileLog(channel_id, SOURCE, dimensions=dimensions, metrics=metrics)


def transform_post(s: dict, category=None, topic=None) -> YoutubePostLog:
    channel_id = None
    if "author" in s:
        channel_id = s["author"]["channelId"]
    shortcode = s["videoId"]
    metrics = [
        safe_metric(s, 'stats.comments', COMMENTS, type=int),
        safe_metric(s, 'stats.likes', LIKES, type=int),
        safe_metric(s, 'stats.views', VIEW_COUNT, type=int),
    ]

    dimensions = [
        safe_dimension(s, 'title', TITLE),
        safe_dimension(s, 'author.title', CHANNEL_TITLE),
        safe_dimension(s, 'description', DESCRIPTION),
        dimension(_convert_seconds_to_iso(safe_get(s, 'lengthSeconds')), DURATION),
        dimension(_convert_post_time_to_epoch(safe_get(s, 'date.published')), PUBLISH_TIME),
        dimension(json.dumps(safe_get(s, 'thumbnails')), THUMBNAILS),
        safe_dimension(s, 'thumbnails.high.url', THUMBNAIL),
        dimension(json.dumps(safe_get(s, 'keywords')), TAGS),
        safe_dimension(s, 'category.id', CATEGORY_ID, type=str),
        safe_dimension(s, 'language.defaultAudioLanguage', AUDIO_LANGUAGE),
        safe_dimension(s, 'content.licensed', IS_LICENSED, type=bool),
        safe_dimension(s, 'endScreen.items.1.playlist.playlistId', PLAYLIST_ID),
        dimension(json.dumps(safe_get(s, 'content.rating')), CONTENT_RATING),
    ]

    if category:
        dimensions.append(Dimension(CATEGORY, category))

    if topic:
        dimensions.append(Dimension(TOPIC, topic))

    duration = safe_get(s, 'lengthSeconds', type=float)
    if duration and duration > 60:
        dimensions.append(Dimension(POST_TYPE, 'video'))
    
    return YoutubePostLog(channel_id, shortcode, SOURCE, dimensions=dimensions, metrics=metrics)


def transform_post_ids_by_channel_id(s: dict, post_type) -> str:
    return s[post_type]['videoId']