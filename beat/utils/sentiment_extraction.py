from datetime import datetime

from core.entities.entities import SentimentLog
from utils.request import make_scrape_log_event
from utils.sentiment_analysis import get_sentiments


async def comment_sentiment_extraction(payload, session=None):
    platform = payload['platform']
    post_activity_logs = payload['post_activity_logs']
    now = datetime.now()
    input_payload = []
    comment_id = ''
    comment = ''
    for post_activity_log in post_activity_logs:
        for dimension in post_activity_log['dimensions']:
            if dimension['key'] == "comment_id":
                comment_id = dimension['value']
            if dimension['key'] == "text":
                comment = dimension['value']
        input_payload.append({"text": comment, "id": comment_id})

    response = await get_sentiments(input_payload, session=session)

    for post_activity_log in post_activity_logs:
        for dimension in post_activity_log['dimensions']:
            if dimension['key'] == "text":
                comment = dimension['value']
            if dimension['key'] == "comment_id":
                comment_id = dimension['value']

        sentiment_log = SentimentLog(
            platform=platform,
            platform_post_id=post_activity_log['platform_post_id'],
            comment_id=comment_id,
            comment=comment,
            sentiment=response[comment_id]['sentiment'],
            score=response[comment_id]['score'],
            dimensions=[],
            metrics=[],
            source="beat",
            timestamp=now
        )
        await make_scrape_log_event("sentiment_log", sentiment_log)
