from core.crawler.crawler import Crawler
from core.helpers.session import sessionize
from gpt.functions.retriever.openai.openai_extractor import OpenAi
from instagram.models.models import InstagramProfileLog


class GptCrawler(Crawler):

    @sessionize
    async def fetch_instagram_gpt_data_base_gender(self, data: dict, session=None) -> dict:
        return await OpenAi.fetch_instagram_gpt_data_base_gender(data)

    @sessionize
    async def fetch_instagram_gpt_data_base_location(self, data: dict, session=None) -> dict:
        return await OpenAi.fetch_instagram_gpt_data_base_location(data)

    @sessionize
    async def fetch_instagram_gpt_data_base_categ_lang_topics(self, data: dict, session=None) -> dict:
        return await OpenAi.fetch_instagram_gpt_data_base_categ_lang_topics(data)

    @sessionize
    async def fetch_instagram_gpt_data_audience_age_gender(self, data: dict,
                                                           session=None) -> dict:
        return await OpenAi.fetch_instagram_gpt_data_audience_age_gender(data)

    @sessionize
    async def fetch_instagram_gpt_data_audience_cities(self, data: dict, session=None) -> dict:
        return await OpenAi.fetch_instagram_gpt_data_audience_cities(data)

    @sessionize
    async def fetch_instagram_gpt_data_gender_location_lang(self, data: dict, session=None) -> dict:
        return await OpenAi.fetch_instagram_gpt_data_gender_location_lang(data)

    async def parse_instagram_gpt_data_base_gender(self, resp_dict: dict) -> InstagramProfileLog:
        return await OpenAi.parse_instagram_gpt_data_base_gender(resp_dict)

    async def parse_instagram_gpt_data_base_location(self, resp_dict: dict) -> InstagramProfileLog:
        return await OpenAi.parse_instagram_gpt_data_base_location(resp_dict)

    async def parse_instagram_gpt_data_base_categ_lang_topics(self, resp_dict: dict) -> InstagramProfileLog:
        return await OpenAi.parse_instagram_gpt_data_base_categ_lang_topics(resp_dict)

    async def parse_instagram_gpt_data_audience_age_gender(self, resp_dict: dict) -> InstagramProfileLog:
        return await OpenAi.parse_instagram_gpt_data_audience_age_gender(resp_dict)

    async def parse_instagram_gpt_data_audience_cities(self, resp_dict: dict) -> InstagramProfileLog:
        return await OpenAi.parse_instagram_gpt_data_audience_cities(resp_dict)

    async def parse_instagram_gpt_data_gender_location_lang(self, resp_dict: dict) -> InstagramProfileLog:
        return await OpenAi.parse_instagram_gpt_data_gender_location_lang(resp_dict)
