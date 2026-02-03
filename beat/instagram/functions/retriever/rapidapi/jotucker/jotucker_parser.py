from core.models.models import Dimension
from instagram.metric_dim_store import MEDIA_COUNT, FOLLOWERS, FOLLOWING, PROFILE_ID, PROFILE_PIC, BIO, IS_VERIFIED, \
    IS_PRIVATE, HASHTAGS, CATEGORY, CATEGORY_ID, FBID, FULL_NAME, COMMENTS, LIKES, DISPLAY_URL, THUMBNAIL, POST_TYPE, \
    PUBLISH_TIME, CAPTION, VIEW_COUNT, HANDLE, PLAY_COUNT, PROFILE_TYPE, TEXT, COMMENTS_THREAD, CITY, LOCATION_LONG, \
    LOCATION_LAT, PUBLIC_EMAIL, CONTACT_PHONE, EXTERNAL_URL, USER_TAGS, LOCATION, COMMENT_ID, USER_SPONSORS, IS_SPONSOR
from instagram.models.models import InstagramProfileLog, InstagramPostLog, \
    InstagramPostActivityLog, InstagramRelationshipLog
from utils.getter import safe_metric, safe_dimension, safe_get, dimension

SOURCE = "rapidapi-jotucker"


def _convert_profile_type(value: int) -> str:
    if value == 1: return 'personal'
    if value == 2: return 'business'
    if value == 3: return 'creator'


def _convert_post_type(value: str | int) -> str:
    if value == 'GraphImage' or value == 1:
        return 'image'
    if value == 'GraphVideo' or value == 2:
        return 'reels'
    if value == 'GraphSidecar' or value == 8:
        return 'carousel'

def _get_hashtags(caption: str) -> str:
    if not caption:
        return ""
    words = caption.split()
    hashtags = [word for word in words if word[0] == '#']
    return ','.join(hashtags)
    

def transform_profile(s: dict) -> InstagramProfileLog:
    handle = s['user']['username']
    profile_id = str(s['user']['pk'])
    metrics = [safe_metric(s, 'user.media_count', MEDIA_COUNT),
               safe_metric(s, 'user.follower_count', FOLLOWERS),
               safe_metric(s, 'user.following_count', FOLLOWING),
               ]

    dimensions = [
        safe_dimension(s, 'user.pk', PROFILE_ID, type=str),
        safe_dimension(s, 'user.username', HANDLE),
        safe_dimension(s, 'user.hd_profile_pic_url_info.url', PROFILE_PIC),
        safe_dimension(s, 'user.biography', BIO),
        safe_dimension(s, 'user.full_name', FULL_NAME),
        safe_dimension(s, 'user.public_email', PUBLIC_EMAIL),
        safe_dimension(s, 'user.external_url', EXTERNAL_URL),
        safe_dimension(s, 'user.contact_phone_number', CONTACT_PHONE),
        dimension(_convert_profile_type(safe_get(s, 'user.account_type')), PROFILE_TYPE),
        safe_dimension(s, 'user.interop_messaging_user_fbid', FBID, type=str),
        safe_dimension(s, 'user.category_id', CATEGORY_ID, type=str),
        safe_dimension(s, 'user.category', CATEGORY),
        safe_dimension(s, 'user.city_name', CITY),
        safe_dimension(s, 'user.latitude', LOCATION_LAT, type=str),
        safe_dimension(s, 'user.longitude', LOCATION_LONG, type=str),
        safe_dimension(s, 'user.is_private', IS_PRIVATE, type=bool),
        safe_dimension(s, 'user.is_verified', IS_VERIFIED, type=bool),
    ]
    log = InstagramProfileLog(
        handle=handle,
        profile_id=profile_id,
        metrics=metrics,
        dimensions=dimensions,
        source=SOURCE)
    return log


def transform_post(s: dict) -> InstagramPostLog:
    metrics = [
        safe_metric(s, 'comment_count', COMMENTS),
        safe_metric(s, 'like_count', LIKES),
        safe_metric(s, 'view_count', VIEW_COUNT),
        safe_metric(s, 'play_count', PLAY_COUNT),
    ]

    post_type = _convert_post_type(safe_get(s, 'media_type'))
    thumbnail = safe_get(s, 'image_versions2.candidates.0.url')
    if post_type == "carousel":
        thumbnail = safe_get(s, 'carousel_media.0.image_versions2.candidates.0.url')

    sponsor_users_raw = safe_get(s, 'sponsor_tags')
    if not sponsor_users_raw or sponsor_users_raw == "":
        sponsor_users_raw = []
    sponsor_users = []
    is_sponsored = False
    for sponsor_user in sponsor_users_raw:
        is_sponsored = True
        if 'sponsor' in sponsor_user:
            sponsor_users.append({"profile_id": sponsor_user['sponsor']['pk'],
                                  "handle": sponsor_user['sponsor']['username']})

    tagged_users_raw = safe_get(s, 'usertags.in')
    if not tagged_users_raw or tagged_users_raw == "":
        tagged_users_raw = []
    tagged_users = []
    for tagged_user in tagged_users_raw:
        if 'user' in tagged_user:
            tagged_users.append(tagged_user['user'])  # tagged_users should be a list of objects with user.username

    dimensions = [
        safe_dimension(s, 'user.pk', PROFILE_ID, type=str),
        safe_dimension(s, 'user.username', HANDLE),
        dimension(thumbnail, DISPLAY_URL),
        dimension(thumbnail, THUMBNAIL),
        dimension(post_type, POST_TYPE),
        safe_dimension(s, 'taken_at', PUBLISH_TIME),
        safe_dimension(s, 'caption.text', CAPTION),
        safe_dimension(s, 'location.name', LOCATION),
        safe_dimension(s, 'location.lng', LOCATION_LONG),
        safe_dimension(s, 'location.lat', LOCATION_LAT),
        dimension(tagged_users, USER_TAGS),
        dimension(sponsor_users, USER_SPONSORS),
        dimension(is_sponsored, IS_SPONSOR),
    ]

    list(filter(None, metrics))
    list(filter(None, dimensions))
    shortcode = s['code']
    log = InstagramPostLog(
        shortcode=shortcode,
        source=SOURCE,
        metrics=metrics,
        dimensions=dimensions)
    return log


def transform_post_by_shortcode(s: dict) -> InstagramPostLog:
    metrics = [
        safe_metric(s, 'comment_count', COMMENTS),
        safe_metric(s, 'like_count', LIKES),
        safe_metric(s, 'play_count', PLAY_COUNT)
    ]

    sponsor_users_raw = safe_get(s, 'sponsor_tags')
    if not sponsor_users_raw or sponsor_users_raw == "":
        sponsor_users_raw = []
    sponsor_users = []
    is_sponsored = False
    for sponsor_user in sponsor_users_raw:
        is_sponsored = True
        if 'sponsor' in sponsor_user:
            sponsor_users.append({"profile_id": sponsor_user['sponsor']['pk'],
                                  "handle": sponsor_user['sponsor']['username']})

    tagged_users_raw = safe_get(s, 'usertags.in')
    if not tagged_users_raw or tagged_users_raw == "":
        tagged_users_raw = []
    tagged_users = []
    for tagged_user in tagged_users_raw:
        if 'user' in tagged_user:
            tagged_users.append(tagged_user['user'])  # tagged_users should be a list of objects with user.username

    dimensions = [
        safe_dimension(s, 'user.pk', PROFILE_ID, type=str),
        safe_dimension(s, 'user.username', HANDLE),
        safe_dimension(s, 'video_versions.0.url', DISPLAY_URL),
        safe_dimension(s, 'image_versions2.candidates.0.url', THUMBNAIL),
        dimension(_convert_post_type(safe_get(s, 'media_type')), POST_TYPE),
        safe_dimension(s, 'taken_at', PUBLISH_TIME),
        safe_dimension(s, 'caption.text', CAPTION),
        dimension(tagged_users, USER_TAGS),
        dimension(sponsor_users, USER_SPONSORS),
        dimension(is_sponsored, IS_SPONSOR),
    ]

    list(filter(None, metrics))
    list(filter(None, dimensions))
    shortcode = s['code']
    log = InstagramPostLog(
        shortcode=shortcode,
        source=SOURCE,
        metrics=metrics,
        dimensions=dimensions)
    return log


def transform_follower(profile_id: str, s: dict) -> InstagramRelationshipLog:
    follower_id = s['id']
    dimensions = [
        Dimension(key=HANDLE, value=safe_get(s, 'username')),
        Dimension(key=FULL_NAME, value=safe_get(s, 'full_name')),
        Dimension(key=PROFILE_PIC, value=safe_get(s, 'profile_pic_url')),
        Dimension(key=IS_PRIVATE, value=safe_get(s, 'is_private')),
        Dimension(key=IS_VERIFIED, value=safe_get(s, 'is_verified')),
    ]
    metrics = []
    log = InstagramRelationshipLog(
        source_profile_id=follower_id,
        relationship_type="FOLLOWS",
        target_profile_id=profile_id,
        source_dimensions=dimensions,
        source_metrics=metrics,
        target_dimensions=[],
        target_metrics=[],
        source=SOURCE,
    )
    return log


def transform_following(profile_id: str, s: dict) -> InstagramRelationshipLog:
    following_id = s['id']
    dimensions = [
        Dimension(key=HANDLE, value=safe_get(s, 'username')),
        Dimension(key=FULL_NAME, value=safe_get(s, 'full_name')),
        Dimension(key=PROFILE_PIC, value=safe_get(s, 'profile_pic_url')),
        Dimension(key=IS_PRIVATE, value=safe_get(s, 'is_private')),
        Dimension(key=IS_VERIFIED, value=safe_get(s, 'is_verified')),
    ]
    metrics = []
    log = InstagramRelationshipLog(
        source_profile_id=profile_id,
        relationship_type="FOLLOWS",
        target_profile_id=following_id,
        source_metrics=[],
        source_dimensions=[],
        target_dimensions=dimensions,
        target_metrics=metrics,
        source=SOURCE,
    )
    return log


def transform_post_likes(shortcode: str, s: dict) -> InstagramPostActivityLog:
    activity_type = "LIKE"
    actor_profile_id = safe_get(s, 'id')
    dimensions = [
        safe_dimension(s, 'profile_pic_url', PROFILE_PIC, type=str),
        safe_dimension(s, 'username', HANDLE, type=str),
        safe_dimension(s, 'full_name', FULL_NAME, type=str),
        safe_dimension(s, 'is_private', IS_PRIVATE, type=bool),
        safe_dimension(s, 'is_verified', IS_VERIFIED, type=bool),
    ]
    metrics = [
    ]
    log = InstagramPostActivityLog(
        activity_type=activity_type,
        shortcode=shortcode,
        actor_profile_id=actor_profile_id,
        dimensions=dimensions,
        metrics=metrics,
        source=SOURCE,
    )
    return log


def transform_post_comments(shortcode: str, s: dict) -> InstagramPostActivityLog:
    activity_type = "COMMENT"
    actor_profile_id = safe_get(s, 'owner.id')
    dimensions = [
        safe_dimension(s, 'id', COMMENT_ID, type=str),
        safe_dimension(s, 'text', TEXT, type=str),
        safe_dimension(s, 'created_at', PUBLISH_TIME, type=int),
        safe_dimension(s, 'owner.profile_pic_url', PROFILE_PIC, type=str),
        safe_dimension(s, 'owner.username', HANDLE, type=str),
        safe_dimension(s, 'owner.is_verified', IS_VERIFIED, type=bool),
    ]
    metrics = [
        safe_metric(s, 'edge_liked_by.count', LIKES),
        safe_metric(s, 'edge_threaded_comments.count', COMMENTS),
        safe_metric(s, 'edge_threaded_comments.edges', COMMENTS_THREAD),
    ]
    log = InstagramPostActivityLog(
        activity_type=activity_type,
        shortcode=shortcode,
        actor_profile_id=actor_profile_id,
        dimensions=dimensions,
        metrics=metrics,
        source=SOURCE,
    )
    return log

def transform_hashtag_post(s: dict) -> InstagramPostLog:
    metrics = [
        safe_metric(s, 'edge_media_to_comment.count', COMMENTS),
        safe_metric(s, 'edge_liked_by.count', LIKES),
        safe_metric(s, 'video_view_count', VIEW_COUNT),
    ]

    dimensions = [
        safe_dimension(s, 'owner.id', PROFILE_ID, type=str),
        safe_dimension(s, 'display_url', DISPLAY_URL),
        safe_dimension(s, 'thumbnail_src', THUMBNAIL),
        dimension(_convert_post_type(safe_get(s, '__typename')), POST_TYPE),
        dimension(_get_hashtags(safe_get(s, 'edge_media_to_caption.edges.0.node.text')), HASHTAGS),
        safe_dimension(s, 'taken_at_timestamp', PUBLISH_TIME),
        safe_dimension(s, 'edge_media_to_caption.edges.0.node.text', CAPTION),
    ]

    list(filter(None, metrics))
    list(filter(None, dimensions))
    shortcode = s['shortcode']
    log = InstagramPostLog(
        shortcode=shortcode,
        source=SOURCE,
        metrics=metrics,
        dimensions=dimensions)
    return log