from core.crawler.crawler import Crawler
from core.helpers.session import sessionize
from core.models.models import RequestContext
from credentials.manager import CredentialManager
from shopify.functions.retriever.shopify.shopify import ShoppifyApi
from shopify.models.models import OrderLog
from utils.exceptions import NoAvailableSources

cred = CredentialManager()


class ShoppifyCrawler(Crawler):

    def __init__(self) -> None:
        self.available_sources = {
            'fetch_orders_by_store': ['shopify']
        }
        self.providers = {
            'shopify': ShoppifyApi
        }

    @sessionize
    async def fetch_orders_by_store(self, store: str, limit: int, updated_at_min=None, 
                                    created_at_min=None, created_at_max=None, session=None) -> (dict, str):
        cred = await CredentialManager.get_enabled_cred_for_handle('shopify', store, session=session)
        if not cred:
            raise NoAvailableSources("No source is available at the moment")
        return await ShoppifyApi.fetch_orders(RequestContext(cred), store, limit, updated_at_min=updated_at_min, 
                                              created_at_min=created_at_min, created_at_max=created_at_max), 'shopify'
    

    def parse_orders_data(self, source: str, resp_dict: dict) -> list[OrderLog]:
        return self.providers[source].parse_orders_data(resp_dict)
