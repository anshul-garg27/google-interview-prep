from core.models.models import RequestContext
from shopify.models.models import OrderLog


class ShoppifyCrawlerInterface:

    @staticmethod
    async def fetch_orders(ctx: RequestContext, limit: int, created_at_min=None, created_at_max=None) -> dict:
        pass


    @staticmethod
    def parse_orders_data(resp_dict: dict) -> list[OrderLog]:
        pass