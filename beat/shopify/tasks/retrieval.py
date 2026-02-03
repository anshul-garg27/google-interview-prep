from core.helpers.session import sessionize
from shopify.functions.retriever.crawler import ShoppifyCrawler


@sessionize
async def retrieve_orders_data_by_store(store: str, limit: int, updated_at_min=None, 
                                        created_at_min=None, created_at_max=None, session=None) -> (dict, str):
    crawler = ShoppifyCrawler()
    data, source = await crawler.fetch_orders_by_store(store, limit, updated_at_min=updated_at_min, 
                                                       created_at_min=created_at_min, created_at_max=created_at_max, session=session)
    return data, source
