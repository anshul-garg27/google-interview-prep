import os
from datetime import datetime

import clickhouse_connect
from loguru import logger

from keyword_collection.keyword_collection_logs import keyword_collection_logs
from youtube.flows.refresh_profile import refresh_yt_posts_by_search


async def fetch_data(payload, session=None):
    await refresh_yt_posts_by_search(payload['job_id'], payload['keywords'], payload['start_date'],
                                     payload['end_date'], session=session)


async def final_report_generation(payload) -> str:
    client = clickhouse_connect.get_client(host=os.environ["CLICKHOUSE_URL"],
                                           username=os.environ["CLICKHOUSE_USERNAME"],
                                           password=os.environ["CLICKHOUSE_PASSWORD"])
    job_id = payload['job_id']
    current_date = datetime.today()
    parquet_file_path = f"keyword_collections/{current_date.year}/{current_date.month}/{current_date.day}/{job_id}.parquet"

    query = """
                    INSERT INTO FUNCTION s3('https://%s.s3.ap-south-1.amazonaws.com/%s', '%s', '%s', 'Parquet')
                    with posts as (select max(post.profile_id)                                      profile_id,
                                        post.shortcode                                       shortcode,
                                        argMax(post.publish_time, post.insert_timestamp)                                    post_publish_time,
                                        JSONExtractInt(argMax(post.metrics, post.insert_timestamp), 'comments')             post_comments,
                                        JSONExtractInt(argMax(post.metrics, post.insert_timestamp), 'likes')                post_likes,
                                        JSONExtractInt(argMax(post.metrics, post.insert_timestamp), 'view_count')           post_views,
                                        JSONExtractString(argMax(post.dimensions, post.insert_timestamp), 'audio_language') post_audio_language,
                                        JSONExtractString(argMax(post.dimensions, post.insert_timestamp), 'keywords')       post_keywords,
                                        JSONExtractString(argMax(post.dimensions, post.insert_timestamp), 'thumbnail')      post_thumbnail,
                                        JSONExtractString(argMax(post.dimensions, post.insert_timestamp), 'title')          post_title,
                                        JSONExtractString(argMax(post.dimensions, post.insert_timestamp), 'post_type')      post_type
                                from %s.post_log_events post
                                where 
                                    toStartOfDay(post.event_timestamp) >= today() - INTERVAL 1 DAY
                                    and JSONExtractString(post.dimensions, 'topic') in %s
                                group by shortcode),
                        profileIds as (select profile_id
                                        from posts),
                        shortcodes as (select shortcode from posts),
                        profiles as (select profile.profile_id                                   profile_id,
                                            JSONExtractString(argMax(profile.dimensions, profile.insert_timestamp), 'custom_url')     custom_url,
                                            JSONExtractInt(argMax(profile.metrics, profile.insert_timestamp), 'subscriber_count')  profile_subscriber,
                                            JSONExtractInt(argMax(profile.metrics, profile.insert_timestamp), 'video_count')       profile_upload,
                                            JSONExtractString(argMax(profile.dimensions, profile.insert_timestamp), 'description') profile_description,
                                            JSONExtractString(argMax(profile.dimensions, profile.insert_timestamp), 'thumbnail')   profile_pic_url,
                                            JSONExtractString(argMax(profile.dimensions, profile.insert_timestamp), 'title')       profile_title
                                    from %s.profile_log_events profile
                                    where profile_id in profileIds
                                        and toStartOfDay(profile.event_timestamp) >= today() - INTERVAL 1 DAY
                                        group by profile_id),
                        youtube_account as (select channel_id                        profile_id,
                                                    category,
                                                    if(category is NULL, NULL, 'mya') category_source
                                            from dbt.mart_youtube_account
                                            where channel_id in profileIds),
                        categories as (select shortcode,
                                            JSONExtractString(JSONExtractRaw(argMax(dimensions, insert_timestamp), 'categorization'), 'label') label
                                        from %s.post_log_events
                                        where source = 'categorization'
                                        and toStartOfDay(event_timestamp) >= today() - INTERVAL 1 DAY
                                        and shortcode in shortcodes
                                        group by shortcode)
                    select posts.profile_id                                                            profile_id,
                        posts.shortcode                                                             shortcode,
                        posts.post_title                                                            post_title,
                        COALESCE(
                            NULLIF(REPLACE(posts.post_thumbnail, '/default.jpg', '/hqdefault.jpg'), posts.post_thumbnail),
                            posts.post_thumbnail
                        )  post_thumbnail_url,
                        posts.post_keywords                                                         post_keywords,
                        toDateTime(posts.post_publish_time)                                                     post_publish_time,
                        posts.post_comments                                                         post_comments,
                        posts.post_likes                                                            post_likes,
                        posts.post_views                                                            post_views,
                        posts.post_audio_language                                                   post_language,
                        posts.post_type                                                             post_type,
                        posts.post_views                                                            post_reach,
                        profiles.profile_subscriber                                                 profile_followers,
                        0                                                                           profile_following,
                        profiles.profile_upload                                                     profile_upload,
                        profiles.profile_title                                                      profile_title,
                        COALESCE(
                            NULLIF(REPLACE(profiles.profile_pic_url, '/default.jpg', '/hqdefault.jpg'), profiles.profile_pic_url),
                            profiles.profile_pic_url
                        )  profile_pic_url,
                        profiles.custom_url                                                         handle,
                        if(ya.category is NULL, categories.label, ya.category)                      category,
                        if(ya.category_source is NULL and category != '', 'ml', ya.category_source) category_source,
                        if(profiles.profile_subscriber != 0, round(((post_likes + post_comments) / profiles.profile_subscriber), 2),
                            0)                                                                       engagement_rate
                    from posts
                            left join profiles on posts.profile_id = profiles.profile_id
                            left join categories on posts.shortcode = categories.shortcode
                            left join youtube_account ya on posts.profile_id = ya.profile_id
                    SETTINGS s3_truncate_on_insert=1;
                    """ % (os.environ["AWS_BUCKET"], parquet_file_path, os.environ["AWS_ACCESS_KEY"], os.environ["AWS_SECRET_ACCESS_KEY"], os.environ["EVENT_DB_NAME"], payload['keywords']
                           , os.environ["EVENT_DB_NAME"], os.environ["EVENT_DB_NAME"])

    client.query(query)
    return parquet_file_path


async def demographic_report_generation(payload, parquet_file_path: str) -> str:
    client = clickhouse_connect.get_client(host=os.environ["CLICKHOUSE_URL"],
                                           username=os.environ["CLICKHOUSE_USERNAME"],
                                           password=os.environ["CLICKHOUSE_PASSWORD"])

    job_id = payload['job_id']
    current_date = datetime.today()
    demographic_parquet_file_path = f"keyword_collections/{current_date.year}/{current_date.month}/{current_date.day}/{job_id}_demographic.parquet"
    keyword_collection_logs(demographic_parquet_file_path, job_id)

    demographic_query = """
                                INSERT INTO FUNCTION s3('https://%s.s3.ap-south-1.amazonaws.com/%s', '%s', '%s', 'Parquet')
                                WITH
                                shortcode as (
                                    select shortcode
                                    from s3('https://%s.s3.ap-south-1.amazonaws.com/%s', '%s', '%s', 'Parquet')
                                    ),
                                categories as (
                                    select DISTINCT JSONExtractString(JSONExtractRaw(dimensions, 'categorization'), 'label') label
                                    from %s.post_log_events
                                    where 
                                    toStartOfDay(event_timestamp) >= today() - INTERVAL 1 DAY
                                    and source = 'categorization' and shortcode in shortcode 
                                    and toFloat64OrZero(JSONExtractString(JSONExtractRaw(dimensions, 'categorization'), 'score')) > 0.50
                                ),
                                categorization AS (SELECT category, male_audience_age_gender, female_audience_age_gender, audience_age, audience_gender
                                                    FROM dbt.mart_genre_overview
                                                    WHERE category IN categories
                                                    AND month = 'ALL'
                                                    AND platform = 'INSTAGRAM'
                                                    AND language = 'ALL'
                                                    AND profile_type = 'ALL'),
                                demographic as (SELECT  
                                                        SUM(arrayElement(male_audience_age_gender, '13-17')) / COUNT(*)         AS m_13_17,
                                                        SUM(arrayElement(male_audience_age_gender, '18-24')) / COUNT(*)         AS m_18_24,
                                                        SUM(arrayElement(male_audience_age_gender, '25-34')) / COUNT(*)         AS m_25_34,
                                                        SUM(arrayElement(male_audience_age_gender, '35-44')) / COUNT(*)         AS m_35_44,
                                                        SUM(arrayElement(male_audience_age_gender, '45-54')) / COUNT(*)         AS m_45_54,
                                                        SUM(arrayElement(male_audience_age_gender, '55-64')) / COUNT(*)         AS m_55_64,
                                                        SUM(arrayElement(male_audience_age_gender, '65+')) / COUNT(*)           AS m_65,
                                                        SUM(arrayElement(female_audience_age_gender, '13-17')) / COUNT(*)         AS f_13_17,
                                                        SUM(arrayElement(female_audience_age_gender, '18-24')) / COUNT(*)         AS f_18_24,
                                                        SUM(arrayElement(female_audience_age_gender, '25-34')) / COUNT(*)         AS f_25_34,
                                                        SUM(arrayElement(female_audience_age_gender, '35-44')) / COUNT(*)         AS f_35_44,
                                                        SUM(arrayElement(female_audience_age_gender, '45-54')) / COUNT(*)         AS f_45_54,
                                                        SUM(arrayElement(female_audience_age_gender, '55-64')) / COUNT(*)         AS f_55_64,
                                                        SUM(arrayElement(female_audience_age_gender, '65+')) / COUNT(*)           AS f_65,
                                                        SUM(arrayElement(audience_age, '13-17')) / COUNT(*)         AS g_13_17,
                                                        SUM(arrayElement(audience_age, '18-24')) / COUNT(*)         AS g_18_24,
                                                        SUM(arrayElement(audience_age, '25-34')) / COUNT(*)         AS g_25_34,
                                                        SUM(arrayElement(audience_age, '35-44')) / COUNT(*)         AS g_35_44,
                                                        SUM(arrayElement(audience_age, '45-54')) / COUNT(*)         AS g_45_54,
                                                        SUM(arrayElement(audience_age, '55-64')) / COUNT(*)         AS g_55_64,
                                                        SUM(arrayElement(audience_age, '65+')) / COUNT(*)           AS g_65,
                                                        SUM(arrayElement(audience_gender, 'male_per')) / COUNT(*)   AS male_per,
                                                        SUM(arrayElement(audience_gender, 'female_per')) / COUNT(*) AS female_per,
                                                        map(
                                                            '13-17', round(m_13_17, 2),
                                                            '18-24', round(m_18_24, 2),
                                                            '25-34', round(m_25_34, 2),
                                                            '35-44', round(m_35_44, 2),
                                                            '45-54', round(m_45_54, 2),
                                                            '55-64', round(m_55_64, 2),
                                                            '65+', round(m_65, 2)) male_audience_age_gender_,
                                                        map(
                                                            '13-17', round(f_13_17, 2),
                                                            '18-24', round(f_18_24, 2),
                                                            '25-34', round(f_25_34, 2),
                                                            '35-44', round(f_35_44, 2),
                                                            '45-54', round(f_45_54, 2),
                                                            '55-64', round(f_55_64, 2),
                                                            '65+', round(f_65, 2)) female_audience_age_gender_,
                                                        map(
                                                            'male', male_audience_age_gender_,
                                                            'female', female_audience_age_gender_
                                                        ) audience_age_gender,
                                                        map(
                                                                '13-17', round(g_13_17, 2),
                                                                '18-24', round(g_18_24, 2),
                                                                '25-34', round(g_25_34, 2),
                                                                '35-44', round(g_35_44, 2),
                                                                '45-54', round(g_45_54, 2),
                                                                '55-64', round(g_55_64, 2),
                                                                '65+', round(g_65, 2)
                                                            ) audience_age_,
                                                        map(
                                                                'male_per', round(male_per, 2),
                                                                'female_per', round(female_per, 2)
                                                            ) audience_gender_
                                                FROM categorization)
                            select 
                                audience_age_gender audience_age_gender,
                                audience_age_    audience_age,
                                audience_gender_ audience_gender
                            from demographic
                            SETTINGS s3_truncate_on_insert=1
            """ % (os.environ["AWS_BUCKET"], demographic_parquet_file_path, os.environ["AWS_ACCESS_KEY"], os.environ["AWS_SECRET_ACCESS_KEY"],
                   os.environ["AWS_BUCKET"], parquet_file_path, os.environ["AWS_ACCESS_KEY"], os.environ["AWS_SECRET_ACCESS_KEY"], os.environ["EVENT_DB_NAME"])

    client.query(demographic_query)
    return demographic_parquet_file_path
