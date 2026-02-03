from gpt.functions.retriever.crawler import GptCrawler
from instagram.models.models import InstagramProfileLog


async def parse_instagram_gpt_data_base_gender(data: dict) -> InstagramProfileLog:
    crawler = GptCrawler()
    data = await crawler.parse_instagram_gpt_data_base_gender(data)
    return data


async def parse_instagram_gpt_data_base_location(data: dict) -> InstagramProfileLog:
    crawler = GptCrawler()
    data = await crawler.parse_instagram_gpt_data_base_location(data)
    return data


async def parse_instagram_gpt_data_base_categ_lang_topics(data: dict) -> InstagramProfileLog:
    crawler = GptCrawler()
    data = await crawler.parse_instagram_gpt_data_base_categ_lang_topics(data)
    return data


async def parse_instagram_gpt_data_audience_age_gender(data: dict) -> InstagramProfileLog:
    crawler = GptCrawler()
    data = await crawler.parse_instagram_gpt_data_audience_age_gender(data)
    return data


async def parse_instagram_gpt_data_audience_cities(data: dict) -> InstagramProfileLog:
    crawler = GptCrawler()
    data = await crawler.parse_instagram_gpt_data_audience_cities(data)
    return data


async def parse_instagram_gpt_data_gender_location_lang(data: dict) -> InstagramProfileLog:
    crawler = GptCrawler()
    data = await crawler.parse_instagram_gpt_data_gender_location_lang(data)
    return data
