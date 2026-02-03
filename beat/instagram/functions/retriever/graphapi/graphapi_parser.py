from collections import defaultdict
from datetime import datetime

import loguru

from instagram.metric_dim_store import AUDIENCE_GENDER_AGE, MEDIA_COUNT, FOLLOWERS, FOLLOWING, PROFILE_ID, PROFILE_PIC, \
    BIO, FBID, FULL_NAME, COMMENTS, LIKES, DISPLAY_URL, AUTOMATIC_FORWARD, SWIPE_BACK, \
    SWIPE_DOWN, SWIPE_FORWARD, SWIPE_UP, TAP_BACK, \
    TAP_EXIT, TAP_FORWARD, THUMBNAIL, POST_TYPE, SHARES, IMPRESSIONS, REACH, ENGAGEMENT, SAVED, \
    PUBLISH_TIME, CAPTION, HANDLE, RECENT_POSTS, RECENT_POSTS_CURSOR, POST_ID, IS_BUSINESS_OR_CREATOR, PLAY_COUNT, \
    INTERACTIONS, AUDIENCE_COUNTRY, AUDIENCE_CITY, AUDIENCE_LOCALE, EXTERNAL_URL, REPLIES, NAVIGATION, PROFILE_VISITS, \
    PROFILE_ACTIVITY, FOLLOWS, CONTENT_TYPE, VIDEO_VIEWS, RECENT_STORIES, REACHED_AUDIENCE_COUNTRY, \
    REACHED_AUDIENCE_CITY, REACHED_AUDIENCE_GENDER_AGE, ENGAGED_AUDIENCE_COUNTRY, ENGAGED_AUDIENCE_CITY, \
    ENGAGED_AUDIENCE_GENDER_AGE
from instagram.models.models import InstagramProfileLog, InstagramPostLog
from utils.getter import safe_metric, safe_dimension, safe_get, dimension

SOURCE = "graphapi"


def _convert_to_epoch(value: str) -> int:
    datetime_object = datetime.strptime(value, "%Y-%m-%dT%H:%M:%S%z")
    epoch = int(datetime_object.timestamp())
    return epoch


def _convert_post_type(value: str) -> str:
    if value == 'VIDEO': return 'reels'
    if value == 'IMAGE': return 'image'
    if value == 'CAROUSEL_ALBUM': return 'carousel'


def transform(s: dict) -> InstagramProfileLog:
    details = s['business_discovery']
    handle = details['username']
    profile_id = str(details['ig_id'])
    metrics = [safe_metric(details, 'media_count', MEDIA_COUNT),
               safe_metric(details, 'followers_count', FOLLOWERS),
               safe_metric(details, 'follows_count', FOLLOWING),
               ]

    dimensions = [
        safe_dimension(details, 'ig_id', PROFILE_ID, type=str),
        safe_dimension(details, 'username', HANDLE),
        safe_dimension(details, 'profile_picture_url', PROFILE_PIC),
        safe_dimension(details, 'biography', BIO),
        safe_dimension(details, 'name', FULL_NAME),
        safe_dimension(details, 'website', EXTERNAL_URL),
        safe_dimension(details, 'id', FBID, type=str),
        safe_dimension(details, 'media.data', RECENT_POSTS),
        safe_dimension(details, 'stories.data', RECENT_STORIES),
        safe_dimension(details, 'media.paging.cursors.after', RECENT_POSTS_CURSOR),
        dimension(True, IS_BUSINESS_OR_CREATOR)
    ]
    log = InstagramProfileLog(handle=handle,
                              profile_id=profile_id,
                              metrics=metrics,
                              dimensions=dimensions,
                              source=SOURCE)
    return log


def transform_post(s: dict) -> InstagramPostLog:
    is_story_post = safe_get(s, 'media_product_type') == "STORY"
    is_video = safe_get(s, 'media_type') == "VIDEO"
    metrics = [
        safe_metric(s, 'comments_count', COMMENTS),
        safe_metric(s, 'like_count', LIKES)
    ]

    dimensions = [
        safe_dimension(s, 'username', HANDLE),
        safe_dimension(s, 'id', POST_ID, type=str),
        safe_dimension(s, 'media_url', DISPLAY_URL),
        safe_dimension(s, 'media_url', THUMBNAIL) if not is_video else dimension(None, THUMBNAIL),
        dimension(_convert_post_type(safe_get(s, 'media_type')), POST_TYPE) if not is_story_post else dimension('story', POST_TYPE),
        dimension(_convert_to_epoch(safe_get(s, 'timestamp')), PUBLISH_TIME),
        dimension(safe_get(s, 'media_type').lower(), CONTENT_TYPE) if is_story_post else dimension(None, CONTENT_TYPE),
        safe_dimension(s, 'caption', CAPTION),
    ]

    shortcode = str(s["permalink"].split('/')[-2]) if not is_story_post else str(s["permalink"].split('/')[-1])
    log = InstagramPostLog(
        shortcode=shortcode,
        source=SOURCE,
        metrics=metrics,
        dimensions=dimensions)
    return log


def transform_post_insights(shortcode: str, resp_dict: dict) -> InstagramPostLog:
    dimensions = []
    metrics = []
    for s in resp_dict['data']:
        metric = s["name"]

        if metric == "reach":
            metrics.append(safe_metric(s, 'values.0.value', REACH, type=int))
        elif metric == "impressions":
            metrics.append(safe_metric(s, 'values.0.value', IMPRESSIONS, type=int))
        elif metric == "engagement":
            metrics.append(safe_metric(s, 'values.0.value', ENGAGEMENT, type=float))
        elif metric == "saved":
            metrics.append(safe_metric(s, 'values.0.value', SAVED, type=int))
        elif metric == "likes":
            metrics.append(safe_metric(s, 'values.0.value', LIKES, type=int))
        elif metric == "comments":
            metrics.append(safe_metric(s, 'values.0.value', COMMENTS, type=int))
        elif metric == "shares":
            metrics.append(safe_metric(s, 'values.0.value', SHARES, type=int))
        elif metric == "plays":
            metrics.append(safe_metric(s, 'values.0.value', PLAY_COUNT, type=int))
        elif metric == "total_interactions":
            metrics.append(safe_metric(s, 'values.0.value', INTERACTIONS, type=int))
        elif metric == "video_views":
            metrics.append(safe_metric(s, 'values.0.value', VIDEO_VIEWS, type=int))
        elif metric == "follows":
            metrics.append(safe_metric(s, 'values.0.value', FOLLOWS, type=int))
        elif metric == "profile_visits":
            metrics.append(safe_metric(s, 'values.0.value', PROFILE_VISITS, type=int))
        elif metric == "profile_activity":
            metrics.append(safe_metric(s, 'values.0.value', PROFILE_ACTIVITY, type=int))

    log = InstagramPostLog(
        shortcode=shortcode,
        source=SOURCE,
        metrics=metrics,
        dimensions=dimensions)
    return log

def transform_story_insights(story_id: str, resp_dict: dict) -> InstagramPostLog:
    dimensions = []
    metrics = []
    for s in resp_dict['data']:
        metric = s["name"]
        if metric == "impressions":
            metrics.append(safe_metric(s, 'values.0.value', IMPRESSIONS, type=int))
        elif metric == 'reach':
            metrics.append(safe_metric(s, 'values.0.value', REACH, type=int))
        elif metric == 'engagement':
            metrics.append(safe_metric(s, 'values.0.value', ENGAGEMENT, type=int))
        elif metric == 'saved':
            metrics.append(safe_metric(s, 'values.0.value', SAVED, type=int))
        elif metric == 'shares':
            metrics.append(safe_metric(s, 'values.0.value', SHARES, type=int))
        elif metric == "replies":
            metrics.append(safe_metric(s, 'values.0.value', REPLIES, type=int))
        elif metric == "total_interactions":
            metrics.append(safe_metric(s, 'values.0.value', INTERACTIONS, type=int))
        elif metric == "navigation":
            metrics.append(safe_metric(s, 'values.0.value', NAVIGATION, type=int))
            breakdowns = s['total_value']['breakdowns'][0]['results']
            for dim in breakdowns:
                metric = dim["dimension_values"][0]
                if metric == "automatic_forward":
                    metrics.append(safe_metric(dim, 'value', AUTOMATIC_FORWARD, type=int))
                elif metric == "tap_back":
                    metrics.append(safe_metric(dim, 'value', TAP_BACK, type=int))
                elif metric == "tap_exit":
                    metrics.append(safe_metric(dim, 'value', TAP_EXIT, type=int))
                elif metric == "tap_forward":
                    metrics.append(safe_metric(dim, 'value', TAP_FORWARD, type=int))
                elif metric == "swipe_back":
                    metrics.append(safe_metric(dim, 'value', SWIPE_BACK, type=int))
                elif metric == "swipe_down":
                    metrics.append(safe_metric(dim, 'value', SWIPE_DOWN, type=int))
                elif metric == "swipe_forward":
                    metrics.append(safe_metric(dim, 'value', SWIPE_FORWARD, type=int))
                elif metric == "swipe_up":
                    metrics.append(safe_metric(dim, 'value', SWIPE_UP, type=int))
        elif metric == "profile_activity":
            metrics.append(safe_metric(s, 'values.0.value', PROFILE_ACTIVITY, type=int))
        elif metric == "profile_visits":
            metrics.append(safe_metric(s, 'values.0.value', PROFILE_VISITS, type=int))
        elif metric == "follows":
            metrics.append(safe_metric(s, 'values.0.value', FOLLOWS, type=int))

    log = InstagramPostLog(
        shortcode=story_id,
        source=SOURCE,
        metrics=metrics,
        dimensions=dimensions)
    return log

def transform_profile_insights(resp_dict: dict) -> InstagramProfileLog:
    dimensions = []
    metrics = []
    for name, s in resp_dict.items():
        dimension_insights = name
        cities = {}
        countries = {}
        age_gender = defaultdict(dict)
        for item in s:
            breakdown = item['dimension_keys']
            if "results" in item:
                if breakdown[0] == "country":
                    countries = {result["dimension_values"][0]: result["value"] for result in item["results"]}
                elif breakdown[0] == "city":
                    cities = {result["dimension_values"][0]: result["value"] for result in item["results"]}
                elif breakdown[0] == "gender" or breakdown[1] == "age":
                    for result in item["results"]:
                        age_gender[result["dimension_values"][0]][result["dimension_values"][1]] = result["value"]
        mapping = {"M": "male", "F": "female", "U": "uncategorized"}
        age_gender_final = {mapping[key]: value for key, value in age_gender.items()}

        if dimension_insights == "follower_demographics":
            dimensions.append(dimension(countries, AUDIENCE_COUNTRY))
            dimensions.append(dimension(cities, AUDIENCE_CITY))
            dimensions.append(dimension(age_gender_final, AUDIENCE_GENDER_AGE))
        elif dimension_insights == "engaged_audience_demographics":
            dimensions.append(dimension(countries, REACHED_AUDIENCE_COUNTRY))
            dimensions.append(dimension(cities, REACHED_AUDIENCE_CITY))
            dimensions.append(dimension(age_gender_final, REACHED_AUDIENCE_GENDER_AGE))
        elif dimension_insights == "reached_audience_demographics":
            dimensions.append(dimension(countries, ENGAGED_AUDIENCE_COUNTRY))
            dimensions.append(dimension(cities, ENGAGED_AUDIENCE_CITY))
            dimensions.append(dimension(age_gender_final, ENGAGED_AUDIENCE_GENDER_AGE))

    log = InstagramProfileLog(
        profile_id=None,
        handle=None,
        source=SOURCE,
        dimensions=dimensions,
        metrics=metrics
    )
    return log


def reformat_gender_age_data(resp_dict: dict) -> dict:
    gender_age = {"female": {}, "male": {}, "uncategorized": {}}
    audience_gender_age = resp_dict['values'][0]['value']
    for key, value in audience_gender_age.items():
        if key[0] == 'F':
            gender_age['female'][key[2:]] = value
        elif key[0] == 'M':
            gender_age['male'][key[2:]] = value
        else:
            gender_age['uncategorized'][key[2:]] = value
    return gender_age
