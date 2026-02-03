import asyncio
import os
from datetime import datetime

import clickhouse_connect
from loguru import logger

from core.entities.entities import PostLog
from core.enums import enums
from core.models.models import Context
from instagram.metric_dim_store import CATEGORIZATION
from keyword_collection.categorization import instagram_categorization
from keyword_collection.keyword_collection_logs import keyword_collection_logs
from utils.getter import dimension
from utils.request import make_scrape_log_event


async def intermediate_report_generation(payload):
    job_id = payload['job_id']
    start_date = datetime.strptime(payload['start_date'], "%Y-%m-%d").strftime("%Y-%m-%d %H:%M:%S")
    keyword_collection_logs(start_date, job_id)
    end_date = datetime.strptime(payload['end_date'], "%Y-%m-%d").replace(hour=23, minute=59, second=59).strftime(
        "%Y-%m-%d %H:%M:%S")
    keyword_collection_logs(end_date, job_id)
    job_id = payload['job_id']
    client = clickhouse_connect.get_client(host=os.environ["CLICKHOUSE_URL"],
                                           username=os.environ["CLICKHOUSE_USERNAME"],
                                           password=os.environ["CLICKHOUSE_PASSWORD"])
    current_date = datetime.today()
    parquet_file_path = f"keyword_collections/{current_date.year}/{current_date.month}/{current_date.day}/{job_id}.parquet"
    keyword_string = [f"(?i).*\\b(\\S*?)({keyword})\\b.*" for keyword in payload['keywords']]
    keyword_collection_logs(keyword_string, job_id)
    final_keyword_string = f"({keyword_string})"
    query = """
                    INSERT INTO FUNCTION s3('https://%s.s3.ap-south-1.amazonaws.com/%s', '%s', '%s', 'Parquet')
                    with posts as (select 
                                        shortcode,
                                        profile_id,
                                        argMax(likes, updated_at)    post_likes,
                                        argMax(plays, updated_at)    post_plays,
                                        argMax(comments, updated_at) post_comments,
                                        argMax(caption, updated_at)  post_caption,
                                        max(publish_time)            post_publish_time,
                                        max(post_type)               post_type_,
                                        ''                           extracted_keywords
                                from dbt.stg_beat_instagram_post post
                                where post_type != 'story'
                                    and multiMatchAny(caption, %s)
                                    and publish_time >= '%s'
                                    and publish_time <= '%s'
                                group by shortcode, profile_id),
                        profileIds as (select profile_id
                                        from posts),
                        postIds as (select shortcode
                                    from posts),
                        profile_pics as (select entity_id                                                   profile_id,
                                                concat('https://d24w28i6lzk071.cloudfront.net/', asset_url) profile_pic_url
                                        from dbt.stg_beat_asset_log
                                        where entity_type = 'PROFILE'
                                            and entity_id in profileIds),
                        post_pics as (select entity_id                                                   shortcode,
                                            concat('https://d24w28i6lzk071.cloudfront.net/', asset_url) post_pic_url
                                    from dbt.stg_beat_asset_log
                                    where entity_type = 'POST'
                                        and entity_id in postIds),
                        accounts as (select profile_id,
                                            argMax(followers, updated_at)       followers,
                                            argMax(following, updated_at)       following,
                                            argMax(media_count, updated_at)     profile_upload,
                                            argMax(full_name, updated_at)       profile_description,
                                            argMax(profile_pic_url, updated_at) profile_thumbnail,
                                            argMax(handle, updated_at)          handle
                                    from dbt.stg_beat_instagram_account
                                    where profile_id in profileIds
                                    group by profile_id),
                        insta_accounts as (select ig_id,
                                                primary_category                category,
                                                if(category == '', NULL, 'mia') category_source,
                                                primary_language                language,
                                                country
                                            from dbt.mart_instagram_account mia
                                            where ig_id in profileIds)
                    select posts.shortcode                                                            shortcode,
                        posts.profile_id                                                           profile_id,
                        posts.post_likes                                                           post_likes,
                        posts.post_plays                                                           post_plays,
                        posts.post_comments                                                        post_comments,
                        posts.post_caption                                                         post_caption,
                        posts.post_publish_time                                                    post_publish_time,
                        posts.post_type_                                                           post_type,
                        posts.extracted_keywords                                                   post_keywords,
                        accounts.followers                                                         followers,
                        accounts.following                                                         following,
                        accounts.profile_upload                                                    profile_upload,
                        accounts.handle                                                            handle,
                        accounts.profile_description                                               profile_description,
                        pp.post_pic_url                                                            post_pic_url,
                        pop.profile_pic_url                                                        profile_pic_url,
                        ia.category                                                                category,
                        ia.category_source                                                         category_source,
                        ia.language                                                                post_language,
                        if(followers!=0, round(((post_likes + post_comments) / followers), 2), 0)  engagement_rate,
                        if(post_type = 'reels', 'reels', 'static')                                 post_class,
                        post_plays * (0.94 - (log2(followers) * 0.001))                            _reach_reels,
                        (7.6 - (log10(post_likes) * 0.7)) * 0.85 * post_likes                      _reach_non_reels,
                        if(post_class = 'reels', max2(_reach_reels, 0), max2(_reach_non_reels, 0)) reach
                    from posts
                            inner join accounts on posts.profile_id = accounts.profile_id
                            left join post_pics pp on posts.shortcode = pp.shortcode
                            left join profile_pics pop on posts.profile_id = pop.profile_id
                            left join insta_accounts ia on posts.profile_id = ia.ig_id
                    where ia.country = 'IN'
                    SETTINGS s3_truncate_on_insert=1;
                    """ % (os.environ["AWS_BUCKET"], parquet_file_path, os.environ["AWS_ACCESS_KEY"], os.environ["AWS_SECRET_ACCESS_KEY"],
                           final_keyword_string, start_date, end_date)
    client.query(query)
    return parquet_file_path


async def categorize_data(parquet_file_path, job_id):
    top_1000_post_query = """
                                select
                                shortcode,
                                profile_id,
                                post_caption
                                from s3('https://%s.s3.ap-south-1.amazonaws.com/%s', '%s', '%s', 'Parquet')
                                order by post_plays desc
                                limit 1000
            """ % (os.environ["AWS_BUCKET"], parquet_file_path, os.environ["AWS_ACCESS_KEY"], os.environ["AWS_SECRET_ACCESS_KEY"])

    client = clickhouse_connect.get_client(host=os.environ["CLICKHOUSE_URL"],
                                           username=os.environ["CLICKHOUSE_USERNAME"],
                                           password=os.environ["CLICKHOUSE_PASSWORD"])

    result = client.query(top_1000_post_query)

    limit = 100
    total_post = len(result.result_rows)
    keyword_collection_logs(total_post, job_id)
    total_post = total_post if len(result.result_rows) < 1000 else 1000
    total_iterations = (len(result.result_rows) + limit - 1) // limit

    for iteration in range(total_iterations):
        start_index = iteration * limit
        end_index = min(start_index + limit, total_post)
        tasks = []
        for i in range(start_index, end_index):
            task = asyncio.create_task(instagram_categorization(result.result_rows[i]))
            tasks.append(task)

        tasks, _ = await asyncio.wait(tasks)
        tasks = [task.result() for task in tasks]

        for task in tasks:
            profile_id = task[1]
            shortcode = task[0]
            caption = task[2]
            category_score = task[3]
            now = datetime.now()
            context = Context(now)
            dimensions = [dimension(category_score, CATEGORIZATION)]
            post = PostLog(
                platform=enums.Platform.INSTAGRAM.name,
                profile_id=profile_id,
                platform_post_id=shortcode,
                metrics=[],
                dimensions=[d.__dict__ for d in dimensions],
                source='categorization',
                timestamp=now
            )
            keyword_collection_logs(f"categorization --- {post.platform_post_id}", job_id)
            await make_scrape_log_event("post_log", post)


async def final_report_generation(parquet_file_path):
    client = clickhouse_connect.get_client(host=os.environ["CLICKHOUSE_URL"],
                                           username=os.environ["CLICKHOUSE_USERNAME"],
                                           password=os.environ["CLICKHOUSE_PASSWORD"])

    query = """
                    INSERT INTO FUNCTION s3('https://%s.s3.ap-south-1.amazonaws.com/%s', '%s', '%s', 'Parquet')
                    with posts as (select *
                                from s3('https://%s.s3.ap-south-1.amazonaws.com/%s',
                                        '%s', '%s', 'Parquet')),
                        shortcodes as (select shortcode
                                        from posts),
                        categorization as (select shortcode, JSONExtractString(JSONExtractRaw(argMax(dimensions, insert_timestamp), 'categorization'), 'label') label
                                            from %s.post_log_events
                                            where
                                            toStartOfDay(event_timestamp) >= today() - INTERVAL 1 DAY
                                            and shortcode in shortcodes
                                            and source = 'categorization' and platform = 'INSTAGRAM'
                                            group by shortcode
                                            )
                    select shortcode,
                        profile_id,
                        post_likes,
                        post_plays post_views,
                        post_comments,
                        reach post_reach,
                        post_caption post_title,
                        toDateTime(post_publish_time) post_publish_time,
                        post_type,
                        post_pic_url post_thumbnail_url,
                        post_keywords,
                        followers profile_followers,
                        following profile_following,
                        profile_upload,
                        profile_description profile_title,
                        profile_pic_url,
                        handle,
                        engagement_rate,
                        if(category=='', categorization.label, category) category,
                        if (category_source is NULL and category!='', 'ml', category_source) category_source,
                        post_language

                    from posts
                            left join categorization on posts.shortcode = categorization.shortcode
                    SETTINGS s3_truncate_on_insert=1
            """ % (os.environ["AWS_BUCKET"], parquet_file_path, os.environ["AWS_ACCESS_KEY"], os.environ["AWS_SECRET_ACCESS_KEY"],
                   os.environ["AWS_BUCKET"], parquet_file_path, os.environ["AWS_ACCESS_KEY"], os.environ["AWS_SECRET_ACCESS_KEY"],
                   os.environ["EVENT_DB_NAME"])
    client.query(query)


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
                                    and shortcode in shortcode and source = 'categorization' and platform = 'INSTAGRAM'
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
                   os.environ["AWS_BUCKET"], parquet_file_path, os.environ["AWS_ACCESS_KEY"], os.environ["AWS_SECRET_ACCESS_KEY"],
                   os.environ["EVENT_DB_NAME"])

    client.query(demographic_query)
    return demographic_parquet_file_path
