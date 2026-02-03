import os
from gpt.metric_dim_store import ROW_DATA, STATE, COUNTRY, GENDER, LANGUAGES, TOPICS, CATEGORIES, CITY, HANDLE, \
    AUDIENCE_CITY, AUDIENCE_GENDER_AGE
from instagram.models.models import InstagramProfileLog
from utils.getter import safe_dimension


async def transform_instagram_gpt_data_base_gender(data) -> InstagramProfileLog:
    handle = data['handle']
    prompt_version = os.environ["PROMPT_VERSION"]
    source = f"openapi_{prompt_version}"
    dimensions = [
        safe_dimension(data, 'handle', HANDLE),
        safe_dimension(data, 'gender', GENDER),
    ]

    log = InstagramProfileLog(
        profile_id=handle,
        handle=handle,
        source=source,
        dimensions=dimensions,
        metrics=[]
    )
    return log


async def transform_instagram_gpt_data_base_location(data) -> InstagramProfileLog:
    handle = data['handle']
    prompt_version = os.environ["PROMPT_VERSION"]
    source = f"openapi_{prompt_version}"
    dimensions = [
        safe_dimension(data, 'handle', HANDLE),
        safe_dimension(data, 'city', CITY),
        safe_dimension(data, 'state', STATE),
        safe_dimension(data, 'country', COUNTRY)
    ]

    log = InstagramProfileLog(
        profile_id=handle,
        handle=handle,
        source=source,
        dimensions=dimensions,
        metrics=[]
    )
    return log


async def transform_instagram_gpt_data_base_categ_lang_topics(data) -> InstagramProfileLog:
    handle = data['handle']
    prompt_version = os.environ["PROMPT_VERSION"]
    source = f"openapi_{prompt_version}"
    dimensions = [
        safe_dimension(data, 'handle', HANDLE),
        safe_dimension(data, 'categories', CATEGORIES, type=str),
        safe_dimension(data, 'languages', LANGUAGES),
        safe_dimension(data, 'topics', TOPICS)
    ]

    log = InstagramProfileLog(
        profile_id=handle,
        handle=handle,
        source=source,
        dimensions=dimensions,
        metrics=[]
    )
    return log


async def transform_instagram_gpt_data_audience_age_gender(data) -> InstagramProfileLog:
    handle = data['handle']
    prompt_version = os.environ["PROMPT_VERSION"]
    source = f"openapi_{prompt_version}"
    dimensions = [
        safe_dimension(data, 'handle', HANDLE),
        safe_dimension(data, 'audience.age_gender', AUDIENCE_GENDER_AGE),
        safe_dimension(data, 'row_data', ROW_DATA)
    ]

    log = InstagramProfileLog(
        profile_id=handle,
        handle=handle,
        source=source,
        dimensions=dimensions,
        metrics=[]
    )
    return log


async def transform_instagram_gpt_data_audience_cities(data) -> InstagramProfileLog:
    handle = data['handle']
    prompt_version = os.environ["PROMPT_VERSION"]
    source = f"openapi_{prompt_version}"
    dimensions = [
        safe_dimension(data, 'handle', HANDLE),
        safe_dimension(data, 'audience.cities', AUDIENCE_CITY),
        safe_dimension(data, 'row_data', ROW_DATA)
    ]

    log = InstagramProfileLog(
        profile_id=handle,
        handle=handle,
        source=source,
        dimensions=dimensions,
        metrics=[]
    )
    return log


async def transform_instagram_gpt_data_gender_location_lang(data) -> InstagramProfileLog:
    handle = data['handle']
    prompt_version = os.environ["PROMPT_VERSION"]
    source = f"openapi_{prompt_version}"
    dimensions = [
        safe_dimension(data, 'handle', HANDLE),
        safe_dimension(data, 'city', CITY),
        safe_dimension(data, 'state', STATE),
        safe_dimension(data, 'country', COUNTRY),
        safe_dimension(data, 'languages', LANGUAGES),
        safe_dimension(data, 'gender', GENDER),
    ]

    log = InstagramProfileLog(
        profile_id=handle,
        handle=handle,
        source=source,
        dimensions=dimensions,
        metrics=[]
    )
    return log
