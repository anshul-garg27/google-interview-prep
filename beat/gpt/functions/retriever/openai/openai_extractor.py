import json
import os.path

import openai
import yaml
from loguru import logger
from gpt.functions.retriever.interface import GptCrawlerInterface
from gpt.functions.retriever.openai.openai_parser import transform_instagram_gpt_data_audience_age_gender, \
    transform_instagram_gpt_data_audience_cities, transform_instagram_gpt_data_base_categ_lang_topics, \
    transform_instagram_gpt_data_base_gender, transform_instagram_gpt_data_base_location, \
    transform_instagram_gpt_data_gender_location_lang
from instagram.models.models import InstagramProfileLog


def fetch_prompt(prompt_type: str):
    dir_name = os.path.dirname
    logger.debug(dir_name(dir_name(dir_name(dir_name(__file__)))))
    directory_path = f"{dir_name(dir_name(dir_name(dir_name(__file__))))}/prompts"
    prompt_log_path = f"{directory_path}/{os.environ['PROMPT_VERSION']}.yaml"

    with open(prompt_log_path, 'r') as file:
        data = yaml.safe_load(file)
    return data[prompt_type]['system'], data[prompt_type]['query'], data[prompt_type]['model']


class OpenAi(GptCrawlerInterface):

    @staticmethod
    async def fetch_instagram_gpt_data_base_gender(data: dict) -> dict:
        openai.api_key = os.environ["OPENAI_KEY"]
        openai.api_type = os.environ["OPENAI_API_TYPE"]
        openai.api_base = os.environ["OPENAI_API_BASE"]
        openai.api_version = os.environ["OPENAI_API_VERSION"]

        system, query, model = fetch_prompt("base_gender")
        query = query.format(name=data['name'], handle=data['handle'], bio=data['bio'])

        response = openai.ChatCompletion.create(
            engine=model,
            messages=[
                {"role": "system", "content": system},
                {"role": "user", "content": query}
            ],
            temperature=0
        )
        reply = response['choices'][0]['message']['content']
        return json.loads(reply)

    @staticmethod
    async def fetch_instagram_gpt_data_base_location(data: dict) -> dict:
        openai.api_key = os.environ["OPENAI_KEY"]
        openai.api_type = os.environ["OPENAI_API_TYPE"]
        openai.api_base = os.environ["OPENAI_API_BASE"]
        openai.api_version = os.environ["OPENAI_API_VERSION"]

        system, query, model = fetch_prompt("base_localtion")
        query = query.format(name=data['name'], handle=data['handle'], bio=data['bio'])

        response = openai.ChatCompletion.create(
            engine=model,
            messages=[
                {"role": "system", "content": system},
                {"role": "user", "content": query}
            ],
            temperature=0
        )
        reply = response['choices'][0]['message']['content']
        return json.loads(reply)

    @staticmethod
    async def fetch_instagram_gpt_data_base_categ_lang_topics(data: dict) -> dict:
        openai.api_key = os.environ["OPENAI_KEY"]
        openai.api_type = os.environ["OPENAI_API_TYPE"]
        openai.api_base = os.environ["OPENAI_API_BASE"]
        openai.api_version = os.environ["OPENAI_API_VERSION"]

        system, query, model = fetch_prompt("base_categ_lang_topics")
        query = query.format(name=data['name'], handle=data['handle'], bio=data['bio'])

        response = openai.ChatCompletion.create(
            engine=model,
            messages=[
                {"role": "system", "content": system},
                {"role": "user", "content": query}
            ],
            temperature=0
        )
        reply = response['choices'][0]['message']['content']
        return json.loads(reply)

    @staticmethod
    async def fetch_instagram_gpt_data_audience_age_gender(data: dict) -> dict:
        openai.api_key = os.environ["OPENAI_KEY"]
        openai.api_type = os.environ["OPENAI_API_TYPE"]
        openai.api_base = os.environ["OPENAI_API_BASE"]
        openai.api_version = os.environ["OPENAI_API_VERSION"]

        system, query, model = fetch_prompt("audience_age_gender")
        query = query.format(name=data['name'], handle=data['handle'], bio=data['bio'])

        response = openai.ChatCompletion.create(
            engine=model,
            messages=[
                {"role": "system", "content": system},
                {"role": "user", "content": query}
            ],
            temperature=0
        )
        reply = response['choices'][0]['message']['content']
        return json.loads(reply)

    @staticmethod
    async def fetch_instagram_gpt_data_audience_cities(data: dict) -> dict:
        openai.api_key = os.environ["OPENAI_KEY"]
        openai.api_type = os.environ["OPENAI_API_TYPE"]
        openai.api_base = os.environ["OPENAI_API_BASE"]
        openai.api_version = os.environ["OPENAI_API_VERSION"]

        system, query, model = fetch_prompt("audience_cities")
        query = query.format(name=data['name'], handle=data['handle'], bio=data['bio'])

        response = openai.ChatCompletion.create(
            engine=model,
            messages=[
                {"role": "system", "content": system},
                {"role": "user", "content": query}
            ],
            temperature=0
        )
        reply = response['choices'][0]['message']['content']
        return json.loads(reply)

    @staticmethod
    async def fetch_instagram_gpt_data_gender_location_lang(data: dict) -> dict:
        openai.api_key = os.environ["OPENAI_KEY"]
        openai.api_type = os.environ["OPENAI_API_TYPE"]
        openai.api_base = os.environ["OPENAI_API_BASE"]
        openai.api_version = os.environ["OPENAI_API_VERSION"]

        system, query, model = fetch_prompt("base_gender_location_lang")
        query = query.format(name=data['name'], handle=data['handle'], bio=data['bio'])

        response = openai.ChatCompletion.create(
            engine=model,
            messages=[
                {"role": "system", "content": system},
                {"role": "user", "content": query}
            ],
            temperature=0
        )
        reply = response['choices'][0]['message']['content']
        result = None
        try:
            result = json.loads(reply)
        except json.JSONDecodeError as e:
            if "Extra data" == e.msg:
                json_objects = reply.split("\n\n")
                for json_object in json_objects:
                    temp = json.loads(json_object)
                    if temp['handle'] == data['handle']:
                        return temp
        return result

    @staticmethod
    async def parse_instagram_gpt_data_base_gender(resp_dict: dict) -> InstagramProfileLog:
        return await transform_instagram_gpt_data_base_gender(resp_dict)

    @staticmethod
    async def parse_instagram_gpt_data_base_location(resp_dict: dict) -> InstagramProfileLog:
        return await transform_instagram_gpt_data_base_location(resp_dict)

    @staticmethod
    async def parse_instagram_gpt_data_base_categ_lang_topics(resp_dict: dict) -> InstagramProfileLog:
        return await transform_instagram_gpt_data_base_categ_lang_topics(resp_dict)

    @staticmethod
    async def parse_instagram_gpt_data_audience_age_gender(resp_dict: dict) -> InstagramProfileLog:
        return await transform_instagram_gpt_data_audience_age_gender(resp_dict)

    @staticmethod
    async def parse_instagram_gpt_data_audience_cities(resp_dict: dict) -> InstagramProfileLog:
        return await transform_instagram_gpt_data_audience_cities(resp_dict)

    @staticmethod
    async def parse_instagram_gpt_data_gender_location_lang(resp_dict: dict) -> InstagramProfileLog:
        return await transform_instagram_gpt_data_gender_location_lang(resp_dict)
