from instagram.metric_dim_store import MEDIA_COUNT, FOLLOWERS, FOLLOWING, PROFILE_ID, PROFILE_PIC, BIO, IS_VERIFIED, \
    IS_PRIVATE, \
    CATEGORY, CATEGORY_ID, FBID, FULL_NAME, COMMENTS, LIKES, DISPLAY_URL, THUMBNAIL, POST_TYPE, \
    PUBLISH_TIME, CAPTION, HANDLE, VIEW_COUNT, PLAY_COUNT, PROFILE_TYPE, EXTERNAL_URL, CITY, CONTACT_PHONE, \
    PUBLIC_EMAIL, LOCATION_LAT, LOCATION_LONG, USER_TAGS, LOCATION, VIDEO_DURATION, CONTENT_TYPE, STORY_LINKS, \
    STORY_HASHTAGS, \
    POST_ID, TAGGED_USER_IDS, USER_SPONSORS, IS_SPONSOR
from instagram.models.models import InstagramProfileLog, InstagramPostLog
from utils.getter import safe_metric, safe_dimension, safe_get, dimension
from loguru import logger
SOURCE = "rapidapi-arraybobo"


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
        dimension(_convert_profile_type(safe_get(s, 'user.account_type')), PROFILE_TYPE),
        safe_dimension(s, 'user.interop_messaging_user_fbid', FBID, type=str),
        safe_dimension(s, 'user.category_id', CATEGORY_ID, type=str),
        safe_dimension(s, 'user.category', CATEGORY),
        safe_dimension(s, 'user.external_url', EXTERNAL_URL),
        safe_dimension(s, 'user.city_name', CITY),
        safe_dimension(s, 'user.latitude', LOCATION_LAT),
        safe_dimension(s, 'user.longitude', LOCATION_LONG),
        safe_dimension(s, 'user.public_email', PUBLIC_EMAIL),
        safe_dimension(s, 'user.contact_phone_number', CONTACT_PHONE),
        safe_dimension(s, 'user.is_private', IS_PRIVATE, type=bool),
        safe_dimension(s, 'user.is_verified', IS_VERIFIED, type=bool),
    ]
    log = InstagramProfileLog(handle=handle,
                              profile_id=profile_id,
                              metrics=metrics,
                              dimensions=dimensions,
                              source=SOURCE)
    return log


def transform_post(s: dict) -> InstagramPostLog:
    metrics = [
        safe_metric(s, 'node.edge_media_to_comment.count', COMMENTS),
        safe_metric(s, 'node.edge_media_preview_like.count', LIKES),
        safe_metric(s, 'node.video_view_count', VIEW_COUNT)
    ]

    sponsor_users_raw = safe_get(s, 'node.edge_media_to_sponsor_user.edges')
    if not sponsor_users_raw or sponsor_users_raw == "":
        sponsor_users_raw = []
    sponsor_users = []
    is_sponsored = False
    for sponsor_user in sponsor_users_raw:
        is_sponsored = True
        if 'node' in sponsor_user:
            sponsor_users.append({"profile_id": sponsor_user['node']['sponsor']['id'],
                                  "handle": sponsor_user['node']['sponsor']['username']})

    tagged_users_raw = safe_get(s, 'node.edge_media_to_tagged_user.edges')
    if not tagged_users_raw or tagged_users_raw == "":
        tagged_users_raw = []
    tagged_users = []
    for tagged_user in tagged_users_raw:
        if 'node' in tagged_user:
            tagged_users.append(tagged_user['node'])  # tagged_users should be a list of objects with user.username

    dimensions = [
        safe_dimension(s, 'node.owner.id', PROFILE_ID, type=str),
        safe_dimension(s, 'node.owner.username', HANDLE),
        dimension(_convert_post_type(safe_get(s, 'node.__typename')), POST_TYPE),
        safe_dimension(s, 'node.display_url', DISPLAY_URL),
        safe_dimension(s, 'node.thumbnail_src', THUMBNAIL),
        safe_dimension(s, 'node.taken_at_timestamp', PUBLISH_TIME),
        dimension(tagged_users, USER_TAGS),
        dimension(sponsor_users, USER_SPONSORS),
        dimension(is_sponsored, IS_SPONSOR),
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
    metrics = [
        safe_metric(s, 'edge_media_to_comment.count', COMMENTS),
        safe_metric(s, 'edge_media_preview_like.count', LIKES),
        safe_metric(s, 'video_play_count', PLAY_COUNT),
        safe_metric(s, 'video_view_count', VIEW_COUNT),
        safe_metric(s, 'video_duration', VIDEO_DURATION),
    ]

    sponsor_users_raw = safe_get(s, 'edge_media_to_sponsor_user.edges')
    if not sponsor_users_raw or sponsor_users_raw == "":
        sponsor_users_raw = []
    sponsor_users = []
    is_sponsored = False
    for sponsor_user in sponsor_users_raw:
        is_sponsored = True
        if 'node' in sponsor_user:
            sponsor_users.append({"profile_id": sponsor_user['node']['sponsor']['id'],
                                  "handle": sponsor_user['node']['sponsor']['username']})

    tagged_users_raw = safe_get(s, 'edge_media_to_tagged_user.edges')
    if not tagged_users_raw or tagged_users_raw == "":
        tagged_users_raw = []
    tagged_users = []
    for tagged_user in tagged_users_raw:
        if 'node' in tagged_user:
            tagged_users.append(tagged_user['node'])  # tagged_users should be a list of objects with user.username

    dimensions = [
        safe_dimension(s, 'owner.id', PROFILE_ID, type=str),
        safe_dimension(s, 'owner.username', HANDLE),
        safe_dimension(s, 'video_url', DISPLAY_URL),
        safe_dimension(s, 'display_resources.0.src', THUMBNAIL),
        dimension(_convert_post_type(safe_get(s, '__typename')), POST_TYPE),
        safe_dimension(s, 'taken_at_timestamp', PUBLISH_TIME),
        dimension(tagged_users, USER_TAGS),
        dimension(sponsor_users, USER_SPONSORS),
        dimension(is_sponsored, IS_SPONSOR),
        safe_dimension(s, 'edge_media_to_caption.edges.0.node.text', CAPTION),
        safe_dimension(s, 'location', LOCATION),
        # dimension("video", CONTENT_TYPE),
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


def transform_reels_post(s: dict) -> InstagramPostLog:
    metrics = [
        safe_metric(s, 'media.comment_count', COMMENTS),
        safe_metric(s, 'media.like_count', LIKES),
        safe_metric(s, 'media.view_count', VIEW_COUNT),
        safe_metric(s, 'media.play_count', PLAY_COUNT)
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

    dimensions = [
        safe_dimension(s, 'media.user.pk', PROFILE_ID, type=str),
        safe_dimension(s, 'media.user.username', HANDLE),
        safe_dimension(s, 'media.video_versions.0.url', DISPLAY_URL),
        safe_dimension(s, 'media.image_versions2.candidates.0.url', THUMBNAIL),
        dimension(_convert_post_type(safe_get(s, 'media.media_type')), POST_TYPE),
        safe_dimension(s, 'media.taken_at', PUBLISH_TIME),
        safe_dimension(s, 'media.caption.text', CAPTION),
        safe_dimension(s, 'usertags', USER_TAGS),
        dimension(sponsor_users, USER_SPONSORS),
        dimension(is_sponsored, IS_SPONSOR),
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


def transform_tagged_post(s: dict, profile_id: str) -> InstagramPostLog:
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
        dimension(profile_id, TAGGED_USER_IDS)
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

def transform_story_post_data(s: dict, handle: str) -> InstagramPostLog:
    story_links=[]
    hashtag_names=[]
    metrics = []
    if 'story_link_stickers' in s: 
        story_links = [story_link_sticker['story_link']['display_url'] for story_link_sticker in s['story_link_stickers']]
    
    if 'story_hashtags' in s:
        hashtag_names = [story_hashtag['hashtag']['name'] for story_hashtag in s['story_hashtags']]
    
    display_url = s['image_versions2']['candidates'][0]['url']
    if 'video_versions' in s:
        display_url = s['video_versions'][0]['url']

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
        safe_dimension(s, 'taken_at', PUBLISH_TIME),
        dimension(sponsor_users, USER_SPONSORS),
        dimension(is_sponsored, IS_SPONSOR),
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