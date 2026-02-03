from shopify.models.models import OrderLog
from shopify.functions.retriever.crawler import ShoppifyCrawler


async def parse_orders_data(data: dict, source: str) -> list[OrderLog]:
    crawler = ShoppifyCrawler()
    data = crawler.parse_orders_data(source, data)
    return data
