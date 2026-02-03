from core.models.models import Dimension
from instagram.metric_dim_store import MEDIA_COUNT, FOLLOWERS, FOLLOWING, PROFILE_ID, PROFILE_PIC, BIO, IS_VERIFIED, \
    IS_PRIVATE, HASHTAGS, FULL_NAME, COMMENTS, LIKES, DISPLAY_URL, THUMBNAIL, POST_TYPE, \
    PUBLISH_TIME, CAPTION, VIEW_COUNT, HANDLE, PLAY_COUNT, PROFILE_TYPE, TEXT, COMMENTS_THREAD, EXTERNAL_URL, USER_TAGS
from instagram.models.models import InstagramProfileLog, InstagramPostLog, \
    InstagramPostActivityLog, InstagramRelationshipLog
from utils.getter import safe_metric, safe_dimension, safe_get, dimension

SOURCE = "rapidapi-bestsolutions"


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


def transform_profile(s: dict) -> InstagramProfileLog:
    handle = s['username']
    profile_id = str(s['pk'])
    metrics = [safe_metric(s, 'media_count', MEDIA_COUNT),
               safe_metric(s, 'follower_count', FOLLOWERS),
               safe_metric(s, 'following_count', FOLLOWING),
               ]

    dimensions = [
        safe_dimension(s, 'pk', PROFILE_ID, type=str),
        safe_dimension(s, 'username', HANDLE),
        safe_dimension(s, 'hd_profile_pic_url_info.url', PROFILE_PIC),
        safe_dimension(s, 'biography', BIO),
        safe_dimension(s, 'full_name', FULL_NAME),
        safe_dimension(s, 'external_url', EXTERNAL_URL),
        dimension(_convert_profile_type(safe_get(s, 'account_type')), PROFILE_TYPE),
        safe_dimension(s, 'is_private', IS_PRIVATE, type=bool),
        safe_dimension(s, 'is_verified', IS_VERIFIED, type=bool),
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

    dimensions = [
        safe_dimension(s, 'user.pk', PROFILE_ID, type=str),
        safe_dimension(s, 'user.username', HANDLE),
        safe_dimension(s, thumbnail, DISPLAY_URL),
        safe_dimension(s, thumbnail, THUMBNAIL),
        dimension(post_type, POST_TYPE),
        safe_dimension(s, 'taken_at', PUBLISH_TIME),
        safe_dimension(s, 'caption.text', CAPTION),
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
