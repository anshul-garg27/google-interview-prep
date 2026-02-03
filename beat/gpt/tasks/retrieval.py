from core.helpers.session import sessionize
from gpt.functions.retriever.crawler import GptCrawler


@sessionize
async def retrieve_instagram_gpt_data_base_gender(data: dict, session=None) -> dict:
    crawler = GptCrawler()
    data = await crawler.fetch_instagram_gpt_data_base_gender(data, session=session)
    return data


@sessionize
async def retrieve_instagram_gpt_data_base_location(data: dict, session=None) -> dict:
    crawler = GptCrawler()
    data = await crawler.fetch_instagram_gpt_data_base_location(data, session=session)
    return data


@sessionize
async def retrieve_instagram_gpt_data_base_categ_lang_topics(data: dict, session=None) -> dict:
    crawler = GptCrawler()
    data = await crawler.fetch_instagram_gpt_data_base_categ_lang_topics(data, session=session)
    return data


@sessionize
async def retrieve_instagram_gpt_data_audience_age_gender(data: dict, session=None) -> dict:
    crawler = GptCrawler()
    data = await crawler.fetch_instagram_gpt_data_audience_age_gender(data, session=session)
    return data


@sessionize
async def retrieve_instagram_gpt_data_audience_cities(data: dict, session=None) -> dict:
    crawler = GptCrawler()
    data = await crawler.fetch_instagram_gpt_data_audience_cities(data, session=session)
    return data


@sessionize
async def retrieve_instagram_gpt_data_gender_location_lang(data: dict, session=None) -> dict:
    crawler = GptCrawler()
    data = await crawler.fetch_instagram_gpt_data_gender_location_lang(data, session=session)
    return data
