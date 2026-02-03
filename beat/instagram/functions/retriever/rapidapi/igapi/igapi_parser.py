from core.models.models import Dimension
from instagram.metric_dim_store import PROFILE_PIC, IS_VERIFIED, \
    IS_PRIVATE, FULL_NAME, HANDLE
from instagram.models.models import InstagramRelationshipLog
from utils.getter import safe_get

SOURCE = "rapidapi-igapi"


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
