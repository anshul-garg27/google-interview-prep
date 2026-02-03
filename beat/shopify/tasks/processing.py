from datetime import datetime

from loguru import logger
from core.enums import enums
from core.entities.entities import OrderLogs
from core.helpers.session import sessionize
from core.models.models import Context
from shopify.entities.entities import Order
from shopify.models.models import OrderLog
from shopify.tasks.transformer import order_map
from utils.db import get_or_create
from utils.request import make_scrape_log_event

"""
    Process Data and Update DB
"""


@sessionize
async def upsert_orders(order_logs: list[OrderLog], session=None) -> None:
    for log in order_logs:
        await upsert_order(log, session=session)


@sessionize
async def upsert_order(order_log: OrderLog, session=None) -> None:
    now = datetime.now()
    context = Context(now)
    order = OrderLogs(
        platform=enums.Platform.SHOPIFY.name,
        platform_order_id=order_log.platform_order_id,
        metrics=[m.__dict__ for m in order_log.metrics],
        dimensions=[d.__dict__ for d in order_log.dimensions],
        source=order_log.source,
        timestamp=now
    )
    await upsert_shopify_order(context, order_log, order_log.platform_order_id, session=session)
    await make_scrape_log_event("order_log", order)


@sessionize
async def upsert_shopify_order(context: Context, order_log: OrderLog, platform_order_id: str, session=None) -> Order:
    order: Order = (await get_or_create(session, Order, platform_order_id=platform_order_id))[0]
    order.updated_at = context.now
    order.platform_order_id = order_log.platform_order_id

    for metric in order_log.metrics:
        if metric.key in order_map:
            setattr(order, order_map[metric.key], metric.value)
        else:
            logger.debug("Key is missing - %s" % metric.key)

    for dimension in order_log.dimensions:
        if dimension.key in order_map:
            setattr(order, order_map[dimension.key], dimension.value)
        else:
            logger.debug("Key is missing - %s" % dimension.key)
    return order