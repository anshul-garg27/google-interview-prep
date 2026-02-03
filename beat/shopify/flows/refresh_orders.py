from core.helpers.session import sessionize
from core.models.models import Dimension
from shopify.metric_dim_store import STORE
from shopify.models.models import OrderLog
from shopify.tasks.ingestion import parse_orders_data
from shopify.tasks.processing import upsert_orders
from shopify.tasks.retrieval import retrieve_orders_data_by_store

@sessionize
async def fetch_orders_data_by_store(store: str, limit: int, updated_at_min=None, 
                                     created_at_min=None, created_at_max=None, session=None) -> list[OrderLog]:
    data, source = await retrieve_orders_data_by_store(store, limit, updated_at_min=updated_at_min, created_at_min=created_at_min, created_at_max=created_at_max, session=session)
    logs = await parse_orders_data(data, source)
    for log in logs:
        log.dimensions.append(Dimension(STORE, store))
    return logs


@sessionize
async def refresh_orders_by_store(store: str, limit: int, updated_at_min=None, 
                                  created_at_min=None, created_at_max=None, session=None) -> None:
    orders_data = await fetch_orders_data_by_store(store, limit, updated_at_min=updated_at_min, created_at_min=created_at_min, created_at_max=created_at_max, session=session)
    await process_orders_data(orders_data, session=session)


@sessionize
async def process_orders_data(profile_logs: list[OrderLog], session=None) -> None:
    await upsert_orders(profile_logs, session=session)