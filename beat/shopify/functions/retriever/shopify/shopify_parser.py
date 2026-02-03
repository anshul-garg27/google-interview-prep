import json
from core.models.models import Dimension
from utils.exceptions import OrderScrapingFailed
from utils.getter import safe_metric, safe_dimension, safe_get, dimension
from urllib.parse import parse_qs, urlparse
from shopify.models.models import OrderLog
from shopify.metric_dim_store import PLATFORM_ORDER_ID, STATUS, NAME, STORE, LANDING_SITE, \
        ORDER_DATA, UTM_PARAMS, PLATFORM, STORE_NAME

SOURCE = "shopify"

def _get_utm_params(url: str) -> dict:
    parsed_url = urlparse(url)
    query_params = parse_qs(parsed_url.query)
    params_dict = {}

    for key, value in query_params.items():
        if key.startswith("utm_"):
            param_key = key[4:]
            param_value = value[0]
            params_dict[param_key] = param_value
    return params_dict


def transform_order(s: dict) -> OrderLog:
    platform_order_id = str(s['id'])
    metrics = []
    dimensions =[
        safe_dimension(s, 'id', PLATFORM_ORDER_ID, type=str),
        safe_dimension(s, 'name', NAME),
        safe_dimension(s, 'fulfillments.0.status', STATUS),
        safe_dimension(s, 'landing_site', LANDING_SITE),
        safe_dimension(s, 'line_items.0.vendor', STORE_NAME),
        dimension(json.dumps(_get_utm_params(safe_get(s, 'landing_site'))), UTM_PARAMS),
    ]
    dimensions.append(Dimension(ORDER_DATA, s))
    dimensions.append(Dimension(PLATFORM, 'SHOPIFY'))

    return OrderLog(platform_order_id, SOURCE, dimensions, metrics)