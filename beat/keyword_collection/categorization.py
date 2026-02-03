import os

from loguru import logger

from utils.getter import dimension
from utils.request import make_request
from youtube.metric_dim_store import CATEGORIZATION


async def get_categorization(title):
    headers = {
        'Content-Type': 'application/json',
    }

    json_data = {
        'model': 'CATEGORIZER',
        'input': {
            'text': title,
        },
    }
    url = os.environ["RAY_URL"]
    response = await make_request("POST", url=url, headers=headers, json=json_data)

    return response.json()


async def youtube_categorization(post_log):
    try:
        post_category_log = await get_categorization(post_log.dimensions[0].value)
        post_log.dimensions.append(dimension(post_category_log, CATEGORIZATION))
        logger.debug(post_log)
    except Exception as e:
        logger.debug(f"getting error while categorization :- {e}")
    return post_log


async def instagram_categorization(post_log):
    try:
        post_category_log = await get_categorization(post_log[2])
        post_log = post_log + (post_category_log,)
        logger.debug(post_log)
    except Exception as e:
        logger.debug(f"getting error while categorization :- {e}")
        post_log = post_log + ({},)
    return post_log
