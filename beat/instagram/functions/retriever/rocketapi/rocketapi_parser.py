from loguru import logger

from core.models.models import Dimension
from instagram.metric_dim_store import PROFILE_ID, COMMENTS, LIKES, DISPLAY_URL, THUMBNAIL, POST_TYPE, \
    PUBLISH_TIME, CAPTION, HANDLE, VIEW_COUNT, PLAY_COUNT, USER_TAGS, LOCATION, VIDEO_DURATION, RESHARE_COUNT, POST_ID, \
    CONTENT_TYPE, STORY_LINKS, STORY_HASHTAGS, IS_PAID_PARTNERSHIP, FULL_NAME, PROFILE_PIC, IS_PRIVATE, IS_VERIFIED
from instagram.models.models import InstagramPostLog, InstagramRelationshipLog
from utils.getter import safe_metric, safe_dimension, safe_get, dimension

SOURCE = "rocketapi"


def _convert_profile_type(value: int) -> str:
    if value == 1: return 'personal'
    if value == 2: return 'business'
    if value == 3: return 'creator'


def _convert_post_type(value: str | int) -> str:
    if value == 'feed' or value == 1:
        return 'image'
    if value == 'clips' or value == 2:
        return 'reels'
    if value == 'carousel_container' or value == 8:
        return 'carousel'


def transform_post(s: dict) -> InstagramPostLog:
    metrics = [
        safe_metric(s, 'node.edge_media_to_comment.count', COMMENTS),
        safe_metric(s, 'node.edge_media_preview_like.count', LIKES),
        safe_metric(s, 'node.video_view_count', VIEW_COUNT)
    ]

    dimensions = [
        safe_dimension(s, 'node.owner.id', PROFILE_ID, type=str),
        safe_dimension(s, 'node.owner.username', HANDLE),
        dimension(_convert_post_type(safe_get(s, 'node.__typename')), POST_TYPE),
        safe_dimension(s, 'node.display_url', DISPLAY_URL),
        safe_dimension(s, 'node.thumbnail_src', THUMBNAIL),
        safe_dimension(s, 'node.taken_at_timestamp', PUBLISH_TIME),
        safe_dimension(s, 'node.edge_media_to_caption.edges.0.node.text', CAPTION),
    ]

    list(filter(None, metrics))
    list(filter(None, dimensions))
    shortcode = s['node']['shortcode']
    log = InstagramPostLog(
        shortcode=shortcode,
        source=SOURCE,
        metrics=metrics,
        dimensions=dimensions)
    return log


def transform_post_by_shortcode(s: dict) -> InstagramPostLog:
    items = s['response']['body']['items'][0]
    metrics = [
        safe_metric(items, 'comment_count', COMMENTS),
        safe_metric(items, 'like_count', LIKES),
        safe_metric(items, 'play_count', PLAY_COUNT),
        safe_metric(items, 'video_duration', VIDEO_DURATION),
        safe_dimension(items, 'reshare_count', RESHARE_COUNT)
    ]

    tagged_users_raw = safe_get(items, 'usertags.in')
    if not tagged_users_raw or tagged_users_raw == "":
        tagged_users_raw = []
    tagged_users = []
    for tagged_user in tagged_users_raw:
        tagged_users.append(tagged_user)  # tagged_users should be a list of objects with user.username
    dimensions = [
        safe_dimension(items, 'user.id', PROFILE_ID, type=str),
        safe_dimension(items, 'user.username', HANDLE),
        safe_dimension(items, 'image_versions2.candidates.0.url', DISPLAY_URL),
        safe_dimension(items, 'image_versions2.candidates.0.url', THUMBNAIL),
        dimension(_convert_post_type(safe_get(items, 'product_type')), POST_TYPE),
        safe_dimension(items, 'taken_at', PUBLISH_TIME),
        dimension(tagged_users, USER_TAGS),
        safe_dimension(items, 'caption.text', CAPTION),
        safe_dimension(items, 'location', LOCATION),
        safe_dimension(items, 'is_paid_partnership', IS_PAID_PARTNERSHIP)
        # dimension("video", CONTENT_TYPE),
    ]

    list(filter(None, metrics))
    list(filter(None, dimensions))

    shortcode = items['code']
    log = InstagramPostLog(
        shortcode=shortcode,
        source=SOURCE,
        metrics=metrics,
        dimensions=dimensions)
    return log


def transform_reels_post(s: dict) -> InstagramPostLog:
    metrics = [
        safe_metric(s, 'media.comment_count', COMMENTS),
        safe_metric(s, 'media.like_count', LIKES),
        safe_metric(s, 'media.view_count', VIEW_COUNT),
        safe_metric(s, 'media.play_count', PLAY_COUNT),
        safe_metric(s, 'media.video_duration', VIDEO_DURATION),
        safe_metric(s, 'media.reshare_count', RESHARE_COUNT),
    ]
    tagged_users_raw = safe_get(s, 'usertags.in')
    if not tagged_users_raw or tagged_users_raw == "":
        tagged_users_raw = []
    tagged_users = []
    for tagged_user in tagged_users_raw:
        tagged_users.append(tagged_user) 
        
    dimensions = [
        safe_dimension(s, 'media.user.pk', PROFILE_ID, type=str),
        safe_dimension(s, 'media.user.username', HANDLE),
        safe_dimension(s, 'media.video_versions.0.url', DISPLAY_URL),
        safe_dimension(s, 'media.image_versions2.candidates.0.url', THUMBNAIL),
        dimension(_convert_post_type(safe_get(s, 'media.media_type')), POST_TYPE),
        safe_dimension(s, 'media.taken_at', PUBLISH_TIME),
        safe_dimension(s, 'media.caption.text', CAPTION),
        dimension(tagged_users, USER_TAGS),
        safe_dimension(s, 'location', LOCATION),
        safe_dimension(s, 'is_paid_partnership', IS_PAID_PARTNERSHIP),
    ]

    list(filter(None, metrics))
    list(filter(None, dimensions))
    shortcode = s['media']['code']
    log = InstagramPostLog(
        shortcode=shortcode,
        source=SOURCE,
        metrics=metrics,
        dimensions=dimensions)
    return log


def transform_story_post_data(s: dict, handle: str) -> InstagramPostLog:
    story_links = []
    hashtag_names = []
    metrics = []
    if 'story_link_stickers' in s:
        story_links = [story_link_sticker['story_link']['display_url'] for story_link_sticker in
                       s['story_link_stickers']]

    if 'story_hashtags' in s:
        hashtag_names = [story_hashtag['hashtag']['name'] for story_hashtag in s['story_hashtags']]

    display_url = s['image_versions2']['candidates'][0]['url']
    if 'video_versions' in s:
        display_url = s['video_versions'][0]['url']
    
    tagged_users_raw = safe_get(s, 'reel_mentions')
    if not tagged_users_raw or tagged_users_raw == "":
        tagged_users_raw = []
    tagged_users = []
    for tagged_user in tagged_users_raw:
        tagged_users.append(tagged_user['user']) 
    
    metrics = [
        safe_metric(s, 'video_duration', VIDEO_DURATION),
    ]
    
    dimensions = [
        safe_dimension(s, 'pk', POST_ID, type=str),
        safe_dimension(s, 'user.pk', PROFILE_ID, type=str),
        dimension(handle, HANDLE),
        dimension(_convert_story_type(safe_get(s, 'media_type')), CONTENT_TYPE),
        dimension("story", POST_TYPE),
        dimension(story_links, STORY_LINKS),
        dimension(hashtag_names, STORY_HASHTAGS),
        dimension(display_url, DISPLAY_URL),
        safe_dimension(s, 'image_versions2.candidates.0.url', THUMBNAIL),
        safe_dimension(s, 'location', LOCATION),
        dimension(tagged_users, USER_TAGS),
        safe_dimension(s, 'is_paid_partnership', IS_PAID_PARTNERSHIP)
    ]

    list(filter(None, metrics))
    list(filter(None, dimensions))
    shortcode = str(s['pk'])
    log = InstagramPostLog(
        shortcode=shortcode,
        source=SOURCE,
        metrics=metrics,
        dimensions=dimensions)
    return log


def transform_follower(profile_id: str, s: dict) -> InstagramRelationshipLog:
    follower_id = s['pk_id']
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

def _convert_story_type(value: str | int) -> str:
    if value == 1:
        return 'image'
    if value == 2:
        return 'video'
