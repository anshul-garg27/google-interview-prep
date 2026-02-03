import copy
from loguru import logger

from core.helpers.session import sessionize
from gpt.helper import is_data_consumable, normalize_audience_age_gender, normalize_audience_cities
from gpt.tasks.ingestion import parse_instagram_gpt_data_audience_age_gender, parse_instagram_gpt_data_audience_cities, \
    parse_instagram_gpt_data_base_categ_lang_topics, parse_instagram_gpt_data_base_gender, \
    parse_instagram_gpt_data_base_location, parse_instagram_gpt_data_gender_location_lang
from gpt.tasks.processing import upsert_instagram_gpt_data
from gpt.tasks.retrieval import retrieve_instagram_gpt_data_audience_age_gender, \
    retrieve_instagram_gpt_data_audience_cities, retrieve_instagram_gpt_data_base_categ_lang_topics, \
    retrieve_instagram_gpt_data_base_gender, retrieve_instagram_gpt_data_base_location, \
    retrieve_instagram_gpt_data_gender_location_lang
from instagram.models.models import InstagramProfileLog


@sessionize
async def fetch_instagram_gpt_data_base_gender(data: dict, session=None):
    max_retries = 2
    base_data = await retrieve_instagram_gpt_data_base_gender(data, session=session)

    for _ in range(max_retries):
        keys_to_update = []
        for key, value in base_data.items():
            if not value or (key == 'gender' and value in ["", "UNKNOWN", "unknown"]):
                keys_to_update.append(key)
        logger.debug(keys_to_update)
        if not is_data_consumable(base_data, "base_gender"):
            updated_data = await retrieve_instagram_gpt_data_base_gender(data, session=session)

            for key in keys_to_update:
                base_data[key] = updated_data.get(key, base_data[key])
        else:
            break

    profile_log = await parse_instagram_gpt_data_base_gender(base_data)
    profile_log.handle = data['handle']
    profile_log.profile_id = data['handle']
    return profile_log


@sessionize
async def fetch_instagram_gpt_data_base_location(data: dict, session=None):
    max_retries = 2
    base_data = await retrieve_instagram_gpt_data_base_location(data, session=session)

    for _ in range(max_retries):
        keys_to_update = []
        for key, value in base_data.items():
            if not value or (key in ["city", "state", "country"] and value in ["", "UNKNOWN", "unknown"]):
                keys_to_update.append(key)
        logger.debug(keys_to_update)
        if not is_data_consumable(base_data, "base_localtion"):
            updated_data = await retrieve_instagram_gpt_data_base_location(data, session=session)

            for key in keys_to_update:
                base_data[key] = updated_data.get(key, base_data[key])
        else:
            break

    profile_log = await parse_instagram_gpt_data_base_location(base_data)
    profile_log.handle = data['handle']
    profile_log.profile_id = data['handle']
    return profile_log


@sessionize
async def fetch_instagram_gpt_data_base_categ_lang_topics(data: dict, session=None):
    max_retries = 2
    base_data = await retrieve_instagram_gpt_data_base_categ_lang_topics(data, session=session)

    for _ in range(max_retries):
        keys_to_update = []
        for key, value in base_data.items():
            if not value:
                keys_to_update.append(key)
        logger.debug(keys_to_update)
        if not is_data_consumable(base_data, "base_categ_lang_topics"):
            updated_data = await retrieve_instagram_gpt_data_base_categ_lang_topics(data, session=session)

            for key in keys_to_update:
                base_data[key] = updated_data.get(key, base_data[key])
        else:
            break

    profile_log = await parse_instagram_gpt_data_base_categ_lang_topics(base_data)
    profile_log.handle = data['handle']
    profile_log.profile_id = data['handle']
    return profile_log


@sessionize
async def fetch_instagram_gpt_data_audience_age_gender(data: dict, session=None):
    max_retries = 2
    audience_age_gender_data = await retrieve_instagram_gpt_data_audience_age_gender(data, session=session)
    audience_age_gender_data['row_data'] = copy.deepcopy(audience_age_gender_data['audience'])
    for _ in range(max_retries):
        if not is_data_consumable(audience_age_gender_data, "audience_age_gender"):
            audience_age_gender_data = await retrieve_instagram_gpt_data_audience_age_gender(data,
                                                                                             session=session)
        else:
            break
    category = data['category']
    logger.debug(audience_age_gender_data)
    normalize_audience_age_gender(audience_age_gender_data, category)
    logger.debug(audience_age_gender_data)
    profile_log = await parse_instagram_gpt_data_audience_age_gender(audience_age_gender_data)
    profile_log.handle = data['handle']
    profile_log.profile_id = data['handle']
    return profile_log


@sessionize
async def fetch_instagram_gpt_data_audience_cities(data: dict, session=None):
    max_retries = 2
    audience_cities_data = await retrieve_instagram_gpt_data_audience_cities(data, session=session)
    audience_cities_data['row_data'] = copy.deepcopy(audience_cities_data['audience'])
    audience_cities_data['audience']['cities'] = dict(
        sorted(audience_cities_data['audience']['cities'].items(), key=lambda item: item[1], reverse=True))
    for _ in range(max_retries):
        if not is_data_consumable(audience_cities_data, "audience_cities"):
            audience_cities_data = await retrieve_instagram_gpt_data_audience_cities(data, session=session)
        else:
            break
    normalize_audience_cities(audience_cities_data)
    logger.debug(audience_cities_data)
    profile_log = await parse_instagram_gpt_data_audience_cities(audience_cities_data)
    profile_log.handle = data['handle']
    profile_log.profile_id = data['handle']
    return profile_log


@sessionize
async def fetch_instagram_gpt_data_gender_location_lang(data: dict, session=None):
    max_retries = 1
    base_data = await retrieve_instagram_gpt_data_gender_location_lang(data, session=session)
    logger.debug(base_data)
    for _ in range(max_retries):
        keys_to_update = []
        for key, value in base_data.items():
            if not value or (key in ["city", "state", "country", "gender", "languages"] and value in ["", "UNKNOWN",
                                                                                                      "unknown"]):
                keys_to_update.append(key)
        logger.debug(f"keys_to_update :- {keys_to_update}")
        if not is_data_consumable(base_data, "base_localtion_gender_lang"):
            updated_data = await retrieve_instagram_gpt_data_gender_location_lang(data, session=session)

            for key in keys_to_update:
                base_data[key] = updated_data.get(key, base_data[key])
        else:
            break

    profile_log = await parse_instagram_gpt_data_gender_location_lang(base_data)
    profile_log.handle = data['handle']
    profile_log.profile_id = data['handle']
    return profile_log


@sessionize
async def process_instagram_gpt_data(profile: InstagramProfileLog, session=None) -> None:
    await upsert_instagram_gpt_data(profile, session=session)


@sessionize
async def refresh_instagram_gpt_data_base_gender(data: dict, session=None):
    log = await fetch_instagram_gpt_data_base_gender(data)
    await process_instagram_gpt_data(log, session=session)


@sessionize
async def refresh_instagram_gpt_data_base_location(data: dict, session=None):
    log = await fetch_instagram_gpt_data_base_location(data)
    await process_instagram_gpt_data(log, session=session)


@sessionize
async def refresh_instagram_gpt_data_base_categ_lang_topics(data: dict, session=None):
    log = await fetch_instagram_gpt_data_base_categ_lang_topics(data)
    await process_instagram_gpt_data(log, session=session)


@sessionize
async def refresh_instagram_gpt_data_audience_age_gender(data: dict, session=None):
    log = await fetch_instagram_gpt_data_audience_age_gender(data)
    await process_instagram_gpt_data(log, session=session)


@sessionize
async def refresh_instagram_gpt_data_audience_cities(data: dict, session=None):
    log = await fetch_instagram_gpt_data_audience_cities(data)
    await process_instagram_gpt_data(log, session=session)


@sessionize
async def refresh_instagram_gpt_data_gender_location_lang(data: dict, session=None):
    log = await fetch_instagram_gpt_data_gender_location_lang(data)
    await process_instagram_gpt_data(log, session=session)
