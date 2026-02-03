"""
    Map scrapper keys to entity keys
"""
from shopify.entities.entities import Order
from shopify.metric_dim_store import PLATFORM_ORDER_ID, STATUS, NAME, LANDING_SITE, \
    ORDER_DATA, UTM_PARAMS, STORE, PLATFORM, STORE_NAME

order_map = {
    PLATFORM_ORDER_ID: Order.platform_order_id.name,
    STATUS: Order.status.name,
    NAME: Order.name.name,
    LANDING_SITE: Order.landing_site.name,
    ORDER_DATA: Order.order_data.name,
    UTM_PARAMS: Order.utm_params.name,
    STORE: Order.store.name,
    STORE_NAME: Order.store_name.name,
    PLATFORM: Order.platform.name
}
