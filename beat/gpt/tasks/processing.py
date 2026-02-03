"""
    Process Data and Update DB
"""
from datetime import datetime
from typing import TypeVar

from loguru import logger

from core.entities.entities import ProfileLog
from core.enums import enums
from core.helpers.session import sessionize
from instagram.models.models import InstagramProfileLog
from utils.request import make_scrape_log_event

T = TypeVar("T")


@sessionize
async def upsert_instagram_gpt_data(log: InstagramProfileLog, session=None):
    logger.debug(log)
    now = datetime.now()
    profile_insights = ProfileLog(
        platform=enums.Platform.INSTAGRAM.name,
        profile_id=log.profile_id,
        source=log.source,
        metrics=[m.__dict__ for m in log.metrics],
        dimensions=[m.__dict__ for m in log.dimensions],
        timestamp=now
    )
    # session.add(profile_insights)
    await make_scrape_log_event("profile_log", profile_insights)
