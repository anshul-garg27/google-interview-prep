from core.models.models import Dimension
from instagram.functions.retriever.rapidapi.jotucker.jotucker_parser import _get_hashtags
from instagram.metric_dim_store import PROFILE_ID, COMMENTS, LIKES, DISPLAY_URL, THUMBNAIL, POST_TYPE, \
    PUBLISH_TIME, CAPTION, VIEW_COUNT, HANDLE, PLAY_COUNT, BIO, EXTERNAL_URL, FOLLOWERS, FBID, FOLLOWING, \
    FULL_NAME, IS_BUSINESS_OR_CREATOR, CATEGORY, IS_BUSINESS, IS_PROFESSIONAL, BUSINESS_CATEGORY, PROFILE_PIC, \
    RECENT_POSTS, RECENT_VIDEOS, SEO_CATEGORY_INFO, MEDIA_COUNT, REELS_MEDIA_COUNT, SHORTCODE, DIMENSIONS, \
    TAGGED_USER_HANDLES, TAGGED_USER_IDS, PRODUCT_TYPE, SONG_ARTIST, SONG_NAME, POST_TYPE_ORIGINAL, VIDEO_DURATION, \
    USER_TAGS, LOCATION, LOCATION_LONG, LOCATION_LAT, IS_VERIFIED, IS_PRIVATE, CITY, PUBLIC_EMAIL, CONTACT_PHONE, \
    POST_ID, \
    CONTENT_TYPE, STORY_LINKS, STORY_HASHTAGS, HASHTAGS, IS_PAID_PARTNERSHIP, LOCATION_CITY
from instagram.models.models import InstagramPostLog, InstagramProfileLog, InstagramRelationshipLog
from utils.getter import safe_metric, safe_dimension, safe_get, dimension
from loguru import logger
SOURCE = "lama"


def _convert_post_type(value: str | int) -> str:
    if value == 'GraphImage' or value == 1:
        return 'image'
    if value == 'GraphVideo' or value == 2:
        return 'reels'
    if value == 'GraphSidecar' or value == 8:
        return 'carousel'


def transform_post_by_shortcode(s: dict) -> InstagramPostLog:
    metrics = [
        safe_metric(s, 'comment_count', COMMENTS),
        safe_metric(s, 'like_count', LIKES),
        safe_metric(s, 'view_count', VIEW_COUNT),
        safe_metric(s, 'play_count', PLAY_COUNT)
    ]

    dimensions = [
        safe_dimension(s, 'user.pk', PROFILE_ID, type=str),
        safe_dimension(s, 'user.username', HANDLE),
        safe_dimension(s, 'display_url', DISPLAY_URL),
        safe_dimension(s, 'image_versions2.candidates.0.url', THUMBNAIL),
        dimension(_convert_post_type(safe_get(s, 'media_type')), POST_TYPE),
        safe_dimension(s, 'taken_at', PUBLISH_TIME, type=int),
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


def transform_profile(s: dict) -> InstagramProfileLog:
    if 'graphql' in s:
        return transform_gql_profile(s)
    return transform_v1_profile(s)


def transform_v1_profile(s: dict) -> InstagramProfileLog:
    user = s
    profile_id = user['pk']
    handle = user['username']
    metrics = [
        safe_metric(user, 'follower_count', FOLLOWERS),
        safe_metric(user, 'following_count', FOLLOWING),
        safe_metric(user, 'media_count', MEDIA_COUNT),
    ]
    is_business_or_creator = user['is_business']
    dimensions = [
        safe_dimension(user, 'biography', BIO),
        safe_dimension(user, 'is_verified', IS_VERIFIED, type=bool),
        safe_dimension(user, 'is_private', IS_PRIVATE, type=bool),
        safe_dimension(user, 'external_url', EXTERNAL_URL),
        safe_dimension(user, 'city_name', CITY),
        safe_dimension(user, 'latitude', LOCATION_LAT),
        safe_dimension(user, 'longitude', LOCATION_LONG),
        safe_dimension(user, 'public_email', PUBLIC_EMAIL),
        safe_dimension(user, 'contact_phone_number', CONTACT_PHONE),
        safe_dimension(user, 'interop_messaging_user_fbid', FBID),
        safe_dimension(user, 'username', HANDLE),
        safe_dimension(user, 'full_name', FULL_NAME),
        safe_dimension(user, 'pk', PROFILE_ID),
        safe_dimension(user, 'is_business', IS_BUSINESS, type=bool),
        dimension(is_business_or_creator, IS_BUSINESS_OR_CREATOR),
        safe_dimension(user, 'category_name', CATEGORY),
        safe_dimension(user, 'business_category_name', BUSINESS_CATEGORY),
        safe_dimension(user, 'profile_pic_url_hd', PROFILE_PIC),
    ]

    list(filter(None, metrics))
    list(filter(None, dimensions))
    log = InstagramProfileLog(
        handle=handle,
        profile_id=profile_id,
        source=SOURCE,
        metrics=metrics,
        dimensions=dimensions)
    return log


def transform_gql_profile(s: dict) -> InstagramProfileLog:
    user = s['graphql']['user']
    profile_id = user['id']
    handle = user['username']
    metrics = [
        safe_metric(user, 'edge_followed_by.count', FOLLOWERS),
        safe_metric(user, 'edge_follow.count', FOLLOWING),
        safe_metric(user, 'edge_owner_to_timeline_media.count', MEDIA_COUNT),
        safe_metric(user, 'edge_felix_video_timeline.count', REELS_MEDIA_COUNT),
    ]
    is_business_or_creator = user['is_business_account'] or user['is_professional_account']
    dimensions = [
        safe_dimension(s, 'seo_category_infos', SEO_CATEGORY_INFO, type=str),
        safe_dimension(user, 'biography', BIO),
        safe_dimension(user, 'is_verified', IS_VERIFIED, type=bool),
        safe_dimension(user, 'is_private', IS_PRIVATE, type=bool),
        safe_dimension(user, 'external_url', EXTERNAL_URL),
        safe_dimension(user, 'fbid', FBID),
        safe_dimension(user, 'username', HANDLE),
        safe_dimension(user, 'full_name', FULL_NAME),
        safe_dimension(user, 'id', PROFILE_ID),
        safe_dimension(user, 'is_business_account', IS_BUSINESS, type=bool),
        safe_dimension(user, 'is_professional_account', IS_PROFESSIONAL, type=bool),
        dimension(is_business_or_creator, IS_BUSINESS_OR_CREATOR),
        safe_dimension(user, 'category_name', CATEGORY),
        safe_dimension(user, 'business_category_name', BUSINESS_CATEGORY),
        safe_dimension(user, 'profile_pic_url_hd', PROFILE_PIC),
        safe_dimension(user, 'edge_felix_video_timeline', RECENT_VIDEOS),
        safe_dimension(user, 'edge_owner_to_timeline_media', RECENT_POSTS),
    ]

    list(filter(None, metrics))
    list(filter(None, dimensions))
    log = InstagramProfileLog(
        handle=handle,
        profile_id=profile_id,
        source=SOURCE,
        metrics=metrics,
        dimensions=dimensions)
    return log


def transform_profile_post(s: dict) -> InstagramPostLog:
    post = s
    metrics = [
        safe_metric(post, 'view_count', VIEW_COUNT),
        safe_metric(post, 'comment_count', COMMENTS),
        safe_metric(post, 'like_count', LIKES),
        safe_metric(post, 'play_count', PLAY_COUNT),
        safe_metric(post, 'video_duration', VIDEO_DURATION),
    ]

    dimensions = [
        dimension(_convert_post_type(post['media_type']), POST_TYPE),
        safe_dimension(post, 'code', SHORTCODE, type=str),
        safe_dimension(post, 'dimensions', DIMENSIONS),
        safe_dimension(post, 'thumbnail_url', THUMBNAIL, type=str),
        safe_dimension(post, 'user.pk', PROFILE_ID, type=str),
        safe_dimension(post, 'user.username', HANDLE, type=str),
        safe_dimension(post, 'caption_text', CAPTION, type=str),
        safe_dimension(post, 'product_type', PRODUCT_TYPE, type=str),
        safe_dimension(post, 'taken_at_ts', PUBLISH_TIME, type=int),
        safe_dimension(post, 'usertags', USER_TAGS),
        safe_dimension(post, 'location.name', LOCATION),
        safe_dimension(post, 'location.lng', LOCATION_LONG),
        safe_dimension(post, 'location.lat', LOCATION_LAT),
    ]
    shortcode = post['code']
    list(filter(None, metrics))
    list(filter(None, dimensions))
    log = InstagramPostLog(
        shortcode=shortcode,
        source=SOURCE,
        metrics=metrics,
        dimensions=dimensions)
    return log


def transform_profiles_post(s: dict) -> InstagramPostLog:
    post = s['node']
    metrics = [
        safe_metric(post, 'video_view_count', VIEW_COUNT),
        safe_metric(post, 'edge_media_to_comment.count', COMMENTS),
        safe_metric(post, 'edge_liked_by.count', LIKES),
        safe_metric(post, 'edge_liked_by.count', LIKES),

    ]
    tagged_users = safe_get(post, 'edge_media_to_tagged_user.edges')
    tagged_users_handles = [node['user']['username'] for node in tagged_users]
    tagged_users_ids = [node['user']['id'] for node in tagged_users]

    dimensions = [
        dimension(_convert_post_type(post['__typename']), POST_TYPE),
        safe_dimension(post, '__typename', POST_TYPE_ORIGINAL, type=str),
        safe_dimension(post, 'shortcode', SHORTCODE, type=str),
        safe_dimension(post, 'dimensions', DIMENSIONS),
        safe_dimension(post, 'display_url', THUMBNAIL, type=str),
        dimension(tagged_users_handles, TAGGED_USER_HANDLES),
        dimension(tagged_users_ids, TAGGED_USER_IDS),
        safe_dimension(post, 'owner.id', PROFILE_ID, type=str),
        safe_dimension(post, 'owner.username', HANDLE, type=str),
        safe_dimension(post, 'edge_media_to_caption.edges.0.node.text', CAPTION, type=str),
        safe_dimension(post, 'thumbnail_src', THUMBNAIL, type=str),
        safe_dimension(post, 'product_type', PRODUCT_TYPE, type=str),
        safe_dimension(post, 'clips_music_attribution_info.artist_name', SONG_ARTIST, type=str),
        safe_dimension(post, 'clips_music_attribution_info.song_name', SONG_NAME, type=str),
    ]
    shortcode = post['shortcode']
    list(filter(None, metrics))
    list(filter(None, dimensions))
    log = InstagramPostLog(
        shortcode=shortcode,
        source=SOURCE,
        metrics=metrics,
        dimensions=dimensions)
    return log


def transform_follower(profile_id: str, s: dict) -> InstagramRelationshipLog:
    follower_id = s['pk']
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
    following_id = s['pk']
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

def transform_story_post_data(s: dict) -> InstagramPostLog:
    story_links=[]
    hashtag_names=[]
    metrics = []
    if len(s['stickers']) > 0: 
        story_links = [story_link_sticker['story_link']['display_url'] for story_link_sticker in s['stickers']]
    
    if len(s['hashtags']) > 0:
        hashtag_names = [story_hashtag['hashtag']['name'] for story_hashtag in s['hashtags']]
    
    display_url = s['thumbnail_url']
    if s['video_url'] is not None:
        display_url = s['video_url']
    dimensions = [
        safe_dimension(s, 'pk', POST_ID, type=str),
        safe_dimension(s, 'user.pk', PROFILE_ID, type=str),
        safe_dimension(s, 'user.username', HANDLE, type=str),
        dimension(_convert_story_type(safe_get(s, 'media_type')), CONTENT_TYPE),
        dimension("story", POST_TYPE),
        dimension(story_links, STORY_LINKS),
        dimension(hashtag_names, STORY_HASHTAGS),
        dimension(display_url, DISPLAY_URL),
        safe_dimension(s, 'thumbnail_url', THUMBNAIL),
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

def _convert_story_type(value: str | int) -> str:
    if value == 1:
        return 'image'
    if value == 2:
        return 'video'


def transform_hashtag_post(s: dict) -> InstagramPostLog:
    metrics = [
        safe_metric(s, 'comment_count', COMMENTS),
        safe_metric(s, 'like_count', LIKES),
        safe_metric(s, 'view_count', VIEW_COUNT),
        safe_metric(s, 'play_count', PLAY_COUNT),
        safe_metric(s, 'video_duration', VIDEO_DURATION),
    ]

    dimensions = [
        safe_dimension(s, 'user.pk', PROFILE_ID, type=str),
        safe_dimension(s, 'user.username', HANDLE, type=str),
        safe_dimension(s, 'user.full_name', FULL_NAME, type=str),
        safe_dimension(s, 'display_url', DISPLAY_URL),
        safe_dimension(s, 'thumbnail_url', THUMBNAIL),
        safe_dimension(s, 'is_paid_partnership', IS_PAID_PARTNERSHIP, type=bool),
        safe_dimension(s, 'video_url', DISPLAY_URL),
        dimension(_convert_post_type(safe_get(s, 'media_type')), POST_TYPE),
        dimension(_get_hashtags(safe_get(s, 'caption_text')), HASHTAGS),
        safe_dimension(s, 'taken_at_ts', PUBLISH_TIME, type=int),
        safe_dimension(s, 'caption_text', CAPTION),
        safe_dimension(s, 'location.name', LOCATION),
        safe_dimension(s, 'location.city', LOCATION_CITY),
        safe_dimension(s, 'location.lng', LOCATION_LONG),
        safe_dimension(s, 'location.lat', LOCATION_LAT),
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