from loguru import logger
from core.models.models import RequestContext
from credentials.manager import CredentialManager
from shopify.models.models import OrderLog
from utils.exceptions import OrderScrapingFailed
from utils.request import make_request
from shopify.functions.retriever.interface import ShoppifyCrawlerInterface
from shopify.functions.retriever.shopify.shopify_parser import transform_order

cred_provider = CredentialManager()

class ShoppifyApi(ShoppifyCrawlerInterface):
     

    @staticmethod
    async def fetch_orders(ctx: RequestContext, store, limit: int, updated_at_min=None, 
                           created_at_min=None, created_at_max=None) -> dict:
        headers = ctx.credential.credentials
        url = f"https://{store}/admin/api/2023-04/orders.json"
        querystring = {
            "status": "any",
            "limit": limit
        }
        if updated_at_min:
            querystring["updated_at_min"] = updated_at_min
        if created_at_min:
            querystring["created_at_min"] = created_at_min
        if created_at_max:
            querystring["created_at_max"] = created_at_max

        response = await make_request("GET", url, headers=headers, params=querystring)

        if response.status_code == 200:
            resp_dict = response.json()
            if "orders" in resp_dict and resp_dict["orders"]:
                return resp_dict
        else:
            raise OrderScrapingFailed(response.status_code, "Something went wrong")


    @staticmethod
    def parse_orders_data(resp: dict) -> list[OrderLog]:
        return [transform_order(item) for item in resp['orders']]