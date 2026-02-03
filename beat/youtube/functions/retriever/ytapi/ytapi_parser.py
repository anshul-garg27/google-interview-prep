import json
from datetime import datetime

import dateutil.parser
import isodate
from loguru import logger

from core.models.models import Dimension
from utils.exceptions import PostScrapingFailed
from utils.getter import safe_metric, safe_dimension, safe_get, dimension
from youtube.models.models import YoutubeProfileLog, YoutubePostLog, YoutubePostActivityLog, \
    YoutubeActivityLog, YoutubeProfileRelationshipLog
from youtube.metric_dim_store import CHANNELITEM, DEFINITION, PLAYLISTITEM, POST_TYPE, RECOMMENDATION, SOCIAL, TOPIC, VIDEO_COUNT, \
    VIEW_COUNT, HIDDEN_SUBSCRIBER_COUNT, \
    SUBSCRIBER_COUNT, TITLE, \
    DESCRIPTION, UPLOADS_PLAYLIST, COUNTRY, THUMBNAILS, THUMBNAIL, PUBLISH_TIME, CUSTOM_URL, LIKES, FAVOURITE_COUNT, \
    COMMENTS, TAGS, CATEGORY_ID, AUDIO_LANGUAGE, IS_LICENSED, CHANNEL_TITLE, DURATION, CONTENT_RATING, CATEGORY, \
    CITY, GENDER_AGE, TEXT, COMMENT_ID, SUBSCRIBED_CHANNEL_ID, SUBSCRIBED_CHANNEL_TITLE, SUBSCRIBED_DESCRIPTION, SUBSCRIBED_THUMBNAILS, \
    SUBSCRIBED_THUMBNAIL, SHORTCODE, PLAYLIST_ID

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
        safe_dimension(content_details, 'definition', DEFINITION)
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


def transform_profile_insights(s: dict) -> YoutubeProfileLog:
    profile_city_insights = s['city']
    profile_country_insights = s['country']
    profile_gender_age_insights = s['gender_age']

    profile_city_insights['rows'] = convert_rows_to_key_value(profile_city_insights)
    profile_country_insights['rows'] = convert_rows_to_key_value(profile_country_insights)
    profile_gender_age_insights['rows'] = reformat_gender_age_data(profile_gender_age_insights)

    metrics = []
    dimensions = [
        safe_dimension(profile_city_insights, 'rows', CITY),
        safe_dimension(profile_country_insights, 'rows', COUNTRY),
        safe_dimension(profile_gender_age_insights, 'rows', GENDER_AGE)
    ]
    return YoutubeProfileLog(None, SOURCE, dimensions=dimensions, metrics=metrics)


def transform_playlist_for_post_ids(s: dict) -> str:
    content_details = s["contentDetails"]
    return content_details['videoId']


def transform_post_ids_by_genre(s: dict) -> str:
    id = s["id"]
    return id['videoId']


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
        dimension(int(dateutil.parser.isoparse(safe_get(topLevelComment, 'snippet.publishedAt')).strftime('%s')),
                  PUBLISH_TIME),
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


def convert_rows_to_key_value(logs: dict) -> dict:
    key_value = {row[0]: row[1] for row in logs['rows']}
    return key_value


def reformat_gender_age_data(profile_gender_age_insights: dict) -> dict:
    if not profile_gender_age_insights['rows']:
        return {}
    gender_age = {"female": {}, "male": {}, "uncategorized": {}}
    for row in profile_gender_age_insights['rows']:
        if row[0] == "female":
            gender_age['female'][row[1]] = row[2]
        elif row[0] == "male":
            gender_age['male'][row[1]] = row[2]
        else:
            gender_age['uncategorized'][row[1]] = row[2]
    return gender_age


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


def transform_channel_for_post_ids(s: dict) -> str:
    if 'videoId' not in s['id']:
        raise PostScrapingFailed("Failed to fetch post_id")
    return s['id']['videoId']


def transform_activities_for_channel_id(s: dict) -> YoutubeActivityLog:
    
    snippet = s["snippet"]
    content_details = s["contentDetails"]
    activity_id = s['id']
    activity_type = safe_get(snippet, 'type')
    dimensions = []
    if activity_type == "subscription":
        dimensions = [
            safe_dimension(content_details, 'subscription.resourceId.channelId', SUBSCRIBED_CHANNEL_ID),
            safe_dimension(snippet, 'channelTitle', SUBSCRIBED_CHANNEL_TITLE),
            safe_dimension(snippet, 'description', SUBSCRIBED_DESCRIPTION),
            dimension(json.dumps(safe_get(snippet, 'thumbnails')), SUBSCRIBED_THUMBNAILS),
            safe_dimension(snippet, 'thumbnails.default.url', SUBSCRIBED_THUMBNAIL),
            safe_dimension(snippet, 'publishedAt', PUBLISH_TIME)
        ]
    elif activity_type == "playlistItem":
        dimensions = [
            safe_dimension(content_details, 'playlistItem.resourceId.videoId', SHORTCODE),
            safe_dimension(content_details, 'playlistItem.playlistId', PLAYLIST_ID),
            safe_dimension(snippet, 'channelTitle', CHANNEL_TITLE),
            safe_dimension(snippet, 'title', TITLE),
            safe_dimension(snippet, 'description', DESCRIPTION),
            dimension(json.dumps(safe_get(snippet, 'thumbnails')), THUMBNAILS),
            safe_dimension(snippet, 'thumbnails.default.url', THUMBNAIL),
            safe_dimension(snippet, 'publishedAt', PUBLISH_TIME)
        ]
    elif activity_type == "upload":
        dimensions = [
            safe_dimension(content_details, 'upload.videoId', SHORTCODE),
            safe_dimension(snippet, 'channelTitle', CHANNEL_TITLE),
            safe_dimension(snippet, 'title', TITLE),
            safe_dimension(snippet, 'description', DESCRIPTION),
            dimension(json.dumps(safe_get(snippet, 'thumbnails')), THUMBNAILS),
            safe_dimension(snippet, 'thumbnails.default.url', THUMBNAIL),
            safe_dimension(snippet, 'publishedAt', PUBLISH_TIME)
        ]

    return YoutubeActivityLog(
        source=SOURCE,
        activity_id=activity_id,
        activity_type=activity_type,
        actor_channel_id=safe_get(snippet, 'channelId'),
        metrics=[],
        dimensions=dimensions)


def transform_yt_profile_relationship_for_channel_id(s: dict) -> YoutubeProfileRelationshipLog:
    snippet = s["snippet"]
    subscriber_snippet = s['subscriberSnippet']
    relationship_id = s['id']

    source_channel_id = safe_get(subscriber_snippet, 'channelId')
    target_channel_id = safe_get(snippet, 'resourceId.channelId')
    source_dimensions = [
        safe_dimension(subscriber_snippet, 'title', CHANNEL_TITLE),
        safe_dimension(subscriber_snippet, 'description', DESCRIPTION),
        dimension(json.dumps(safe_get(subscriber_snippet, 'thumbnails')), THUMBNAILS),
        safe_dimension(subscriber_snippet, 'thumbnails.default.url', THUMBNAIL),
    ]

    target_dimensions = [
        safe_dimension(snippet, 'title', CHANNEL_TITLE),
        safe_dimension(snippet, 'description', DESCRIPTION),
        dimension(json.dumps(safe_get(snippet, 'thumbnails')), THUMBNAILS),
        safe_dimension(snippet, 'thumbnails.default.url', THUMBNAIL),
    ]
    published_at = safe_get(snippet, 'publishedAt')

    subscribed_on = datetime.fromisoformat(published_at.replace("Z", "+00:00")).strftime("%Y-%m-%d %H:%M:%S")
    return YoutubeProfileRelationshipLog(
        source=SOURCE,
        relationship_id=relationship_id,
        relationship_type='subscription',
        source_channel_id=source_channel_id,
        target_channel_id=target_channel_id,
        source_dimensions=source_dimensions,
        source_metrics=[],
        target_dimensions=target_dimensions,
        target_metrics=[],
        subscribed_on=subscribed_on
        )