from functools import reduce

from loguru import logger

from core.models.models import Dimension, Metric


def dimension(value: str | list | bool | dict, key: str) -> Dimension:
    dim = Dimension(key, value)
    return dim


def metric(value: float, key: str) -> Metric:
    metric = Metric(key, value)
    return metric


def safe_get(dictionary: dict, keys: str, default: any = None, type: any = None) -> any:
    value = reduce(
        lambda d,
               key: d.get(key, default) if isinstance(d, dict)
        else d[int(key)] if isinstance(d, list) and key.isnumeric() and len(d) > int(key)
        else default,
        keys.split("."),
        dictionary
    )
    if value is None:
        logger.debug("Key is missing - %s" % keys)
        return value
    if type is not None:
        return type(value)
    return value 


def safe_metric(dictionary: dict, keys: str, metric_key: str, type: any = None) -> any:
    return metric(safe_get(dictionary, keys, type=type), metric_key)


def safe_dimension(dictionary: dict, keys: str, dim_key: str, default: any = None, type: any = None) -> any:
    return dimension(safe_get(dictionary, keys, default=None, type=type), dim_key)
