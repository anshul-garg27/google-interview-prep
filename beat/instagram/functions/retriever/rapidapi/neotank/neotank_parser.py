from instagram.metric_dim_store import MEDIA_COUNT, FOLLOWERS, FOLLOWING, PROFILE_ID, PROFILE_PIC, BIO, IS_VERIFIED, IS_PRIVATE, \
    CATEGORY, CATEGORY_ID, FBID, FULL_NAME, CAPTION, PUBLISH_TIME, THUMBNAIL, DISPLAY_URL, \
    LIKES, COMMENTS, HANDLE, POST_ID
from instagram.models.models import InstagramProfileLog, InstagramPostLog
from utils.getter import safe_metric, safe_dimension

SOURCE = "rapidapi-neotank"

# def convert_profile_type(value):
#     if value == 1: return 'personal'
#     if value == 2: return 'business'
#     if value == 3: return 'creator'

# def convert_post_type(value):
#     if value == 'GraphImage' or value == 1: return 'image'
#     if value == 'GraphVideo' or value == 2: return 'reels'
#     if value == 'GraphSidecar' or value == 8: return 'carousel'

def transform_profile(s: dict) -> InstagramProfileLog:
    handle = s['username']
    profile_id = str(s['id'])
    metrics = [safe_metric(s, 'edge_owner_to_timeline_media.count', MEDIA_COUNT),
               safe_metric(s, 'edge_followed_by.count', FOLLOWERS),
               safe_metric(s, 'edge_follow.count', FOLLOWING),
               ]

    dimensions = [
                  safe_dimension(s, 'id', PROFILE_ID, type=str),
                  safe_dimension(s, 'username', HANDLE),
                  safe_dimension(s, 'profile_pic_url', PROFILE_PIC),
                  safe_dimension(s, 'biography', BIO),
                  safe_dimension(s, 'full_name', FULL_NAME),
                  #dimension(convert_profile_type(safe_get(s, 'user.account_type')), PROFILE_TYPE),
                  safe_dimension(s, 'fbid', FBID),
                  safe_dimension(s, 'user.category_id', CATEGORY_ID),
                  safe_dimension(s, 'category_name', CATEGORY),
                  safe_dimension(s, 'is_private', IS_PRIVATE, type=bool),
                  safe_dimension(s, 'is_verified', IS_VERIFIED, type=bool),
                  ]
    log = InstagramProfileLog(source=SOURCE,
                              dimensions=dimensions,
                              metrics=metrics,
                              profile_id=profile_id,
                              handle=handle)
    return log


def transform_post(s: dict) -> InstagramPostLog:
    metrics = [
        safe_metric(s, 'node.edge_media_to_comment.count', COMMENTS),
        safe_metric(s, 'node.edge_media_preview_like.count', LIKES)
    ]

    dimensions = [
        safe_dimension(s, 'node.id', POST_ID),
        safe_dimension(s, 'node.display_url', DISPLAY_URL),
        safe_dimension(s, 'node.thumbnail_src', THUMBNAIL),
        safe_dimension(s, 'node.__typename', THUMBNAIL),
        #dimension(convert_post_type(safe_get(s, 'node.__typename')), POST_TYPE),
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
