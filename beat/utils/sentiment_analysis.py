import asyncio
import os
import loguru
from utils.request import make_scrape_log_event, make_request


async def get_sentiment(comments: list):
    headers = {
        'Content-Type': 'application/json',
    }

    json_data = {
        'model': 'SENTIMENT',
        'input': comments,
    }
    url = os.environ["RAY_URL"]
    response = await make_request("POST", url=url, headers=headers, json=json_data)
    return response.json()


async def get_sentiments(input_payload, session=None):
    try:
        response = await get_sentiment(input_payload)
    except Exception as e:
        loguru.logger.debug(f"getting error while sentiment :- {e}")
        loguru.logger.debug("retried again")
        await asyncio.sleep(2)
        response = await get_sentiment(input_payload)

    await asyncio.sleep(3)
    return response
