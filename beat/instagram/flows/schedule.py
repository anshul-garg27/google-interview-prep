from datetime import datetime
from datetime import timedelta

from loguru import logger
from sqlalchemy import and_, text, desc, null, asc, func
from sqlalchemy import select, or_

from core.entities.entities import ScrapeRequestLog
from core.enums import enums
from core.helpers.session import sessionize
from core.models.models import Context
from instagram.entities.entities import InstagramAccount, InstagramAccountRequest, InstagramPost
from instagram.flows.refresh_profile import refresh_profile_by_profile_id, refresh_profile_by_handle, \
    refresh_post_by_shortcode

"""
TODO:
    - Get List of Handles to scrape
      -> Last Scraped over a week ago
      -> In order of decreasing priority
"""


@sessionize
async def get_handles_to_refresh(size: int = 100, session=None) -> list:
    query = select([InstagramAccountRequest.handle]) \
        .where(InstagramAccountRequest.status == "PENDING") \
        .limit(size)
    result = await session.execute(query)
    return [handle for handle in result]


@sessionize
async def get_profile_ids(context: Context, is_creator: bool = True, size: int = 100, session=None) -> list:
    query = select([InstagramAccount.profile_id]) \
        .where(and_(
        or_(InstagramAccount.is_business_or_creator == is_creator, InstagramAccount.profile_type == 'personal'),
        or_(
            InstagramAccount.refreshed_at <= text("now() - interval '7 day'"),
            InstagramAccount.refreshed_at.is_(null())))) \
        .order_by(desc(InstagramAccount.id)) \
        .limit(size)
    result = await session.execute(query)
    return [r[0] for r in result]


@sessionize
async def get_handles(context: Context, is_business_or_creator: bool = True, size: int = 100, session=None) -> list:
    query = select([InstagramAccount.handle]) \
        .where(and_(
        or_(InstagramAccount.is_business_or_creator == is_business_or_creator,
            InstagramAccount.is_business_or_creator.is_(null())),
        or_(
            InstagramAccount.refreshed_at <= text("now() - interval '7 day'"),
            InstagramAccount.refreshed_at.is_(null())))) \
        .order_by(desc(InstagramAccount.id)) \
        .limit(size)
    result = await session.execute(query)
    return [r[0] for r in result]


@sessionize
async def get_profile_ids_for_view_count_refresh(size: int = 100, session=None) -> list:
    query = select([InstagramPost.profile_id, func.min(InstagramPost.created_at)]) \
        .where(and_(InstagramPost.views.is_(null()), InstagramPost.profile_id.is_not(null()),
                    InstagramPost.created_at >= text("now() - interval '7 day'"))) \
        .order_by(asc(func.min(InstagramPost.created_at))) \
        .group_by(InstagramPost.profile_id) \
        .limit(size)
    result = await session.execute(query)
    return [r[0] for r in result]


@sessionize
async def get_short_codes_to_refresh(size: int, session=None) -> list:
    return []  # TODO: Add Collection Post fetch


@sessionize
async def refresh_creator_profiles_insta(session=None) -> None:
    now = datetime.now()
    context = Context(now)
    await refresh_profiles(context, is_creator=True, session=session)


@sessionize
async def refresh_non_creator_profiles_insta(session=None) -> None:
    now = datetime.now()
    context = Context(now)
    await refresh_profiles(context, is_creator=False, session=session)


@sessionize
async def refresh_profiles(context: Context, is_creator: bool = True, session=None) -> None:
    if is_creator:
        handles = await get_handles(context, is_creator, size=1000)
        logger.error(handles)
        await create_scrape_request_logs_for_handles(context, handles, session=session)
    else:
        # profile_ids = await get_profile_ids(context, is_creator, size=1000)
        # await create_scrape_request_logs_for_profile_ids(context, profile_ids)
        handles = await get_handles(context, is_creator, size=1000)
        logger.error(handles)
        await create_scrape_request_logs_for_handles(context, handles)


@sessionize
async def refresh_missing_handles(session=None) -> None:
    handles = await get_handles_to_refresh(size=100, session=session)
    await create_scrape_request_logs_for_handles(handles, session=session)


@sessionize
async def refresh_collection_posts(session=None) -> None:
    shortcodes = await get_short_codes_to_refresh(size=100, session=session)
    await create_scrape_request_logs_for_shortcodes(shortcodes, session=session)


@sessionize
async def create_scrape_request_logs_for_handles(context, handles: list, session=None) -> None:
    for handle in handles:
        request = ScrapeRequestLog(params={"handle": handle},
                                   platform=enums.Platform.INSTAGRAM.name,
                                   status=enums.ScrapeLogRequestStatus.PENDING.name,
                                   flow=refresh_profile_by_handle.__name__,
                                   priority=1,
                                   created_at=datetime.now(),
                                   expires_at=datetime.now() + timedelta(hours=1))
        session.add(request)


@sessionize
async def create_scrape_request_logs_for_profile_ids(context: Context, profile_ids: list, session=None) -> None:
    for profile_id in profile_ids:
        request = ScrapeRequestLog(params={"profile_id": profile_id},
                                   platform=enums.Platform.INSTAGRAM.name,
                                   status=enums.ScrapeLogRequestStatus.PENDING.name,
                                   flow=refresh_profile_by_profile_id.__name__,
                                   priority=1,
                                   created_at=datetime.now(),
                                   expires_at=datetime.now() + timedelta(hours=1))
        session.add(request)


@sessionize
async def create_scrape_request_logs_for_shortcodes(shortcodes: list, session=None) -> None:
    for shortcode in shortcodes:
        request = ScrapeRequestLog(params={"short_code": shortcode},
                                   platform=enums.Platform.INSTAGRAM.name,
                                   status=enums.ScrapeLogRequestStatus.PENDING.name,
                                   flow=refresh_post_by_shortcode.name,
                                   priority=1,
                                   created_at=datetime.now(),
                                   expires_at=datetime.now() + timedelta(hours=1))
        session.add(request)
