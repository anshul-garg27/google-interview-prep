import csv
import os
from datetime import datetime, timedelta
import secrets
import string
import clickhouse_connect
from loguru import logger
import boto3
from dateutil.relativedelta import relativedelta


def get_sub_query(ids, ids_path) -> str:
    sub_query = ""
    if ids_path:
        sub_query = f"ids as ( select * from url('{ids_path}', 'CSVWithNames'))"
    else:
        sub_query = f"ids as (SELECT arrayJoin({ids}) AS id)"
    return sub_query


async def fetch_keyword_videos_csv_report_generation(keywords, job_id, start_date, end_date) -> str:
    client = get_clickhouse_connection()

    current_date = datetime.today()
    csv_file_path = f"Data_request/fetch_keyword_videos_csv/{current_date.year}/{current_date.month}/{current_date.day}/{job_id}/video_list.csv"
    csv_url = f'https://{os.environ["AWS_BUCKET"]}.s3.ap-south-1.amazonaws.com/{csv_file_path}'
    start_date = start_date + ' 00:00:00'
    end_date = end_date + ' 23:59:59'
    query = f"""
                    
                    INSERT INTO FUNCTION s3('{csv_url}',
                                            '{os.environ["AWS_ACCESS_KEY"]}', '{os.environ["AWS_SECRET_ACCESS_KEY"]}', 'CSVWithNames')
                    with 
                    post as (
                        select * from {os.environ["EVENT_DB_NAME"]}.post_log_events where event_timestamp >= now() - interval 1 day and source = 'ytapi'
                    ),
                    ple as (
                        select * from {os.environ["EVENT_DB_NAME"]}.profile_log_events where event_timestamp >= now() - interval 1 day  and source = 'ytapi'
                    ),
                    final as (select post.shortcode                                                                   id,
                                          argMax(post.publish_time, post.insert_timestamp)                                 publishedAt,
                                          max(post.profile_id)                                                             channelId,
                                          JSONExtractString(argMax(post.dimensions, post.insert_timestamp), 'title')       title,
                                          JSONExtractString(argMax(post.dimensions, post.insert_timestamp), 'description') description,
                                          JSONExtractString(argMax(ple.dimensions, ple.insert_timestamp), 'title')         channelTitle,
                                          JSONExtractString(argMax(post.dimensions, post.insert_timestamp), 'tags')        tags,
                                          JSONExtractString(argMax(post.dimensions, post.insert_timestamp), 'category_id') categoryId,
                                          JSONExtractString(argMax(post.dimensions, post.insert_timestamp),
                                                            'audio_language')                                              defaultAudioLanguage,
                                          JSONExtractString(argMax(post.dimensions, post.insert_timestamp), 'duration')    duration,
                                          extractAll(duration, '([0-9]+)H')[1]                                       AS    hours,
                                          extractAll(duration, '([0-9]+)M')[1]                                       AS    minutes,
                                          extractAll(duration, '([0-9]+)S')[1]                                       AS    seconds,
                                          if(hours == '', 0, toInt64OrZero(hours)) * 60 * 60 +
                                          if(minutes == '', 0, toInt64OrZero(minutes)) * 60 + toInt64OrZero(seconds) AS    total_seconds,
                                          JSONExtractInt(argMax(post.metrics, post.insert_timestamp), 'view_count')        viewCount,
                                          JSONExtractInt(argMax(post.metrics, post.insert_timestamp), 'likes')             likeCount,
                                          JSONExtractInt(argMax(post.metrics, post.insert_timestamp), 'comments')          commentCount,
                                          JSONExtractString(argMax(post.dimensions, post.insert_timestamp), 'keywords')    post_keywords,
                                          JSONExtractString(argMax(post.dimensions, post.insert_timestamp), 'thumbnail')   post_thumbnail,
                                          JSONExtractString(argMax(post.dimensions, post.insert_timestamp), 'post_type')   post_type
                                   from post
                                        left join ple
                                   on post.profile_id = ple.profile_id
                                   where toStartOfDay(post.event_timestamp) >= today() - INTERVAL 1 DAY
                                     and JSONExtractString(post.dimensions
                                       , 'topic') in {keywords}
                                     and (publish_time >= '{start_date}'
                                      or publish_time <= '{end_date}')
                                   group by shortcode
                        )
                    select id,
                           publishedAt,
                           channelId,
                           title,
                           description,
                           channelTitle,
                           tags,
                           categoryId,
                           defaultAudioLanguage,
                           total_seconds duration,
                           viewCount,
                           likeCount,
                           commentCount,
                           post_keywords,
                           post_thumbnail,
                           post_type
                    from final
                        SETTINGS s3_truncate_on_insert = 1;
                    """

    client.query("SET format_csv_null_representation = ''")
    client.query(query)
    return csv_url


async def fetch_keyword_channel_csv_report_generation(keywords, end_date, start_date, job_id) -> str:
    client = get_clickhouse_connection()

    current_date = datetime.today()
    csv_file_path = f"Data_request/fetch_keyword_videos_csv/{current_date.year}/{current_date.month}/{current_date.day}/{job_id}/channel_list.csv"
    csv_url = f'https://{os.environ["AWS_BUCKET"]}.s3.ap-south-1.amazonaws.com/{csv_file_path}'
    query = f"""    
                    
                    INSERT INTO FUNCTION s3('{csv_url}', '{os.environ["AWS_ACCESS_KEY"]}', '{os.environ["AWS_SECRET_ACCESS_KEY"]}', 'CSVWithNames')
                    with profileIds as (select profile_id
                                        from {os.environ["EVENT_DB_NAME"]}.post_log_events post
                                        where toStartOfDay(post.event_timestamp) >= today() - INTERVAL 1 DAY and source = 'ytapi'
                                        and JSONExtractString(post.dimensions
                                            , 'topic') in {keywords}
                                        and (publish_time >= {start_date}
                                        or publish_time <= {end_date})
                                        group by profile_id),
                        profiles as (select profile.profile_id                                                                 channelId,
                                            JSONExtractString(argMax(profile.dimensions, profile.insert_timestamp), 'title')   title,
                                            JSONExtractString(argMax(profile.dimensions, profile.insert_timestamp),
                                                            'description')                                                   description,
                                            JSONExtractString(argMax(profile.dimensions, profile.insert_timestamp),
                                                            'custom_url')                                                    custom_url,
                                            JSONExtractString(argMax(profile.dimensions, profile.insert_timestamp),
                                                            'taken_at_timestamp')                                            publishedAt,
                                            JSONExtractString(argMax(profile.dimensions, profile.insert_timestamp), 'country') country,
                                            JSONExtractInt(argMax(profile.metrics, profile.insert_timestamp), 'view_count')    viewCount,

                                            JSONExtractInt(argMax(profile.metrics, profile.insert_timestamp),
                                                            'subscriber_count')                                                 subscriberCount,
                                            JSONExtractInt(argMax(profile.metrics, profile.insert_timestamp), 'video_count')   videoCount,
                                            JSONExtractString(argMax(profile.dimensions, profile.insert_timestamp),
                                                            'HIDDEN_SUBSCRIBER_COUNT')                                       hiddenSubscriberCount,
                                            JSONExtractString(argMax(profile.dimensions, profile.insert_timestamp),
                                                            'thumbnail')                                                     profile_pic_url

                                    from {os.environ["EVENT_DB_NAME"]}.profile_log_events profile
                                    where 
                                    toStartOfDay(profile.event_timestamp) >= today() - INTERVAL 1 DAY and source = 'ytapi'
                                    and profile_id in profileIds
                                    group by profile_id)
                    select channelId,
                        title,
                        description,
                        custom_url,
                        publishedAt,
                        country,
                        viewCount,
                        subscriberCount,
                        videoCount,
                        hiddenSubscriberCount,
                        COALESCE(
                                NULLIF(REPLACE(profiles.profile_pic_url, '/default.jpg', '/hqdefault.jpg'), profiles.profile_pic_url),
                                profiles.profile_pic_url
                            ) as profile_pic_url
                    from profiles
                        SETTINGS s3_truncate_on_insert = 1
                    """

    client.query("SET format_csv_null_representation = ''")
    client.query(query)
    return csv_url


async def fetch_keyword_video_mapping_report_generation(keywords, job_id, start_date, end_date) -> str:
    client = get_clickhouse_connection()

    current_date = datetime.today()
    csv_file_path = f"Data_request/fetch_keyword_videos_csv/{current_date.year}/{current_date.month}/{current_date.day}/{job_id}/keyword_video_mapping.csv"
    csv_url = f'https://{os.environ["AWS_BUCKET"]}.s3.ap-south-1.amazonaws.com/{csv_file_path}'
    
    query = f"""        
                        INSERT INTO FUNCTION s3('{csv_url}', '{os.environ["AWS_ACCESS_KEY"]}', '{os.environ["AWS_SECRET_ACCESS_KEY"]}', 'CSVWithNames')
                        select JSONExtractString(argMax(post.dimensions, post.insert_timestamp), 'topic') topic,
                            post.shortcode                                                             video_id
                        from {os.environ["EVENT_DB_NAME"]}.post_log_events post
                        where toStartOfDay(post.event_timestamp) >= today() - INTERVAL 1 DAY
                        and platform = 'YOUTUBE'
                        and source = 'ytapi'
                        and JSONExtractString(post.dimensions
                                , 'topic') in {keywords}
                        and (publish_time >= {start_date}
                            or publish_time <= {end_date})
                        group by shortcode
                        SETTINGS s3_truncate_on_insert = 1
                    """

    client.query("SET format_csv_null_representation = ''")
    client.query(query)
    return csv_url


async def fetch_videos_details_csv_report_generation(scrape_id: str, video_ids: list[str] = None,
                                                     video_ids_path: str = None) -> str:
    client = get_clickhouse_connection()

    current_date = datetime.today()
    csv_path = f"Data_request/fetch_video_details_csv/{current_date.year}/{current_date.month}/{current_date.day}/{scrape_id}_video_details.csv"

    sub_query = get_sub_query(video_ids, video_ids_path)
    csv_url = f'https://{os.environ["AWS_BUCKET"]}.s3.ap-south-1.amazonaws.com/{csv_path}'
    query = f"""
                
                INSERT INTO FUNCTION s3('{csv_url}', '{os.environ["AWS_ACCESS_KEY"]}', '{os.environ["AWS_SECRET_ACCESS_KEY"]}', 'CSVWithNames')
                with 
                {sub_query},
                final as (select post.shortcode id,
                    max(post.publish_time) publishedAt,
                    max(post.profile_id) channelId,
                    JSONExtractString(argMax(post.dimensions, post.insert_timestamp), 'title') title,
                    JSONExtractString(argMax(post.dimensions, post.insert_timestamp), 'description') description,
                    JSONExtractString(argMax(post.dimensions, post.insert_timestamp), 'channel_title') channelTitle,
                    JSONExtractString(argMax(post.dimensions, post.insert_timestamp), 'tags') tags,
                    JSONExtractString(argMax(post.dimensions, post.insert_timestamp), 'category_id') categoryId,
                    JSONExtractString(argMax(post.dimensions, post.insert_timestamp), 'audio_language') defaultAudioLanguage,
                    JSONExtractString(argMax(post.dimensions, post.insert_timestamp), 'duration') duration,
                    extractAll(duration, '([0-9]+)H')[1]                                       AS    hours,
                    extractAll(duration, '([0-9]+)M')[1]                                       AS    minutes,
                    extractAll(duration, '([0-9]+)S')[1]                                       AS    seconds,
                    if(hours == '', 0, toInt64OrZero(hours)) * 60 * 60 +
                    if(minutes == '', 0, toInt64OrZero(minutes)) * 60 + toInt64OrZero(seconds) AS    total_seconds,
                    JSONExtractString(argMax(post.metrics, post.insert_timestamp), 'view_count') viewCount,
                    JSONExtractString(argMax(post.metrics, post.insert_timestamp), 'likes') likeCount,
                    JSONExtractString(argMax(post.metrics, post.insert_timestamp), 'comments') commentCount        
                from {os.environ["EVENT_DB_NAME"]}.post_log_events post
                where 
                insert_timestamp >= now()-INTERVAL 1 DAY
                and platform='YOUTUBE'
                and source = 'ytapi'
                and post.shortcode in ids
                GROUP BY
                    post.shortcode)
                select
                    id,
                    publishedAt,
                    channelId,
                    title,
                    description,
                    channelTitle,
                    tags,
                    categoryId,
                    defaultAudioLanguage,
                    total_seconds duration,
                    viewCount,
                    likeCount,
                    commentCount
                from final
                SETTINGS s3_truncate_on_insert=1
    """

    client.query("SET format_csv_null_representation = ''")
    client.query(query)
    return csv_url


async def fetch_playlist_videos_csv_report_generation(job_id: str, playlist_ids: list[str] = None,
                                                      playlist_ids_path: str = None) -> str:
    client = get_clickhouse_connection()

    current_date = datetime.today()
    csv_path = f"Data_request/fetch_playlist_videos_csv/{current_date.year}/{current_date.month}/{current_date.day}/{job_id}_playlist_video_details.csv"
    csv_url = f'https://{os.environ["AWS_BUCKET"]}.s3.ap-south-1.amazonaws.com/{csv_path}'
    
    sub_query = get_sub_query(playlist_ids, playlist_ids_path)

    query = f"""
                
                INSERT INTO FUNCTION s3('{csv_url}', '{os.environ["AWS_ACCESS_KEY"]}', '{os.environ["AWS_SECRET_ACCESS_KEY"]}', 'CSVWithNames')
                with
                {sub_query},
                final as (
                    select post.shortcode id,
                        max(post.publish_time) publishedAt,
                        max(post.profile_id) channelId,
                        JSONExtractString(argMax(post.dimensions, post.insert_timestamp), 'playlist_id') playlistId,
                        JSONExtractString(argMax(post.dimensions, post.insert_timestamp), 'title') title,
                        JSONExtractString(argMax(post.dimensions, post.insert_timestamp), 'description') description,
                        JSONExtractString(argMax(post.dimensions, post.insert_timestamp), 'channel_title') channelTitle,
                        JSONExtractString(argMax(post.dimensions, post.insert_timestamp), 'tags') tags,
                        JSONExtractString(argMax(post.dimensions, post.insert_timestamp), 'category_id') categoryId,
                        JSONExtractString(argMax(post.dimensions, post.insert_timestamp), 'audio_language') defaultAudioLanguage,
                        JSONExtractString(argMax(post.dimensions, post.insert_timestamp), 'duration') duration,
                        extractAll(duration, '([0-9]+)H')[1]                                       AS    hours,
                      extractAll(duration, '([0-9]+)M')[1]                                       AS    minutes,
                      extractAll(duration, '([0-9]+)S')[1]                                       AS    seconds,
                      if(hours == '', 0, toInt64OrZero(hours)) * 60 * 60 +
                      if(minutes == '', 0, toInt64OrZero(minutes)) * 60 + toInt64OrZero(seconds) AS    total_seconds,
                        JSONExtractInt(argMax(post.metrics, post.insert_timestamp), 'view_count') viewCount,
                        JSONExtractInt(argMax(post.metrics, post.insert_timestamp), 'likes') likeCount,
                        JSONExtractInt(argMax(post.metrics, post.insert_timestamp), 'comments') commentCount        
                    from {os.environ["EVENT_DB_NAME"]}.post_log_events post
                    where insert_timestamp >= now()-INTERVAL 1 DAY and platform='YOUTUBE'
                    and source = 'ytapi'
                    and JSONExtractString(post.dimensions, 'playlist_id') in ids
                    GROUP BY
                        post.shortcode
                    )
                select
                id,
               publishedAt,
               channelId,
               playlistId,
               title,
               description,
               channelTitle,
               tags,
               categoryId,
               defaultAudioLanguage,
               total_seconds seconds,
               viewCount,
               likeCount,
               commentCount
                from final
                SETTINGS s3_truncate_on_insert=1
    """

    client.query("SET format_csv_null_representation = ''")
    client.query(query)
    return csv_url


async def get_ids_from_csv_file(ids_path: str) -> list:
    client = get_clickhouse_connection()

    query = f"""
                SELECT
                    *
                FROM url('{ids_path}', 'CSVWithNames')
            """
    result = client.query(query)
    ids = [row[0] for row in result.result_rows]
    return ids


async def get_channel_details_from_mart(channel_ids: list[str] = None, channel_ids_path: str = None) -> list:
    client = get_clickhouse_connection()

    sub_query = get_sub_query(channel_ids, channel_ids_path)
    query = f""" 
                with
                {sub_query}
                SELECT
                    channel_id
                FROM dbt.mart_youtube_account where channel_id IN ids
            """
    result = client.query(query)
    output_chanel_ids = [row[0] for row in result.result_rows]
    return output_chanel_ids


async def fetch_channel_details_csv(job_id: str, channel_ids: list[str] = None, channel_ids_path: str = None) -> str:
    client = get_clickhouse_connection()

    current_date = datetime.today()
    file_path = f"Data_request/fetch_channel_details_csv/{current_date.year}/{current_date.month}/{current_date.day}/{job_id}_channel_details.csv"

    csv_url = f'https://{os.environ["AWS_BUCKET"]}.s3.ap-south-1.amazonaws.com/{file_path}'
    sub_query = get_sub_query(channel_ids, channel_ids_path)
    query = f"""
                
                INSERT INTO FUNCTION s3('{csv_url}',
                                        '{os.environ["AWS_ACCESS_KEY"]}', '{os.environ["AWS_SECRET_ACCESS_KEY"]}', 'CSVWithNames')
                with
                    {sub_query}
                        ,
                    martChannelDetails as (select channel_id  as id,
                                                  title,
                                                  description,
                                                  username    as customUrl,
                                                  country,
                                                  views_count as viewCount,
                                                  followers   as subscriberCount,
                                                  uploads     as videoCount
                                           from dbt.mart_youtube_account
                                           where channel_id in ids),
                    martChannelIds as (
                        select id from martChannelDetails
                    ),
                    remainingChannelIds as (
                        select * from ids where id not in martChannelIds
                    ),
                    eventChannelDetails as (select profile.profile_id                 id,
                                                   JSONExtractString(argMax(profile.dimensions, profile.insert_timestamp),
                                                                     'title')         title,
                                                   JSONExtractString(argMax(profile.dimensions, profile.insert_timestamp),
                                                                     'description')   description,
                                                   JSONExtractString(argMax(profile.dimensions, profile.insert_timestamp),
                                                                     'custom_url')    customUrl,
                                                   JSONExtractString(argMax(profile.dimensions, profile.insert_timestamp),
                                                                     'country')       country,
                                                   JSONExtractInt(argMax(profile.metrics, profile.insert_timestamp),
                                                                  'view_count')       viewCount,

                                                   JSONExtractInt(argMax(profile.metrics, profile.insert_timestamp),
                                                                  'subscriber_count') subscriberCount,
                                                   JSONExtractInt(argMax(profile.metrics, profile.insert_timestamp),
                                                                  'video_count')      videoCount
                                            from {os.environ["EVENT_DB_NAME"]}.profile_log_events profile
                                            where toStartOfDay(profile.event_timestamp) >= today() - INTERVAL 1 DAY and
                                            source = 'ytapi' and profile_id in remainingChannelIds
                                            group by profile_id),

                    final as (select *
                              from martChannelDetails
                              union all
                              select *
                              from eventChannelDetails)
                select *
                from final
                    SETTINGS s3_truncate_on_insert = 1
    """
    client.query("SET format_csv_null_representation = ''")
    client.query(query)
    return csv_url


async def fetch_channel_videos_csv_report(job_id: str, start_date: str, end_date: str, channel_ids: list[str] = None,
                                          channel_ids_path: str = None):    
    client = get_clickhouse_connection()

    current_date = datetime.today()
    file_path = f"Data_request/fetch_channel_videos_csv/{current_date.year}/{current_date.month}/{current_date.day}/{job_id}_channel_videos.csv"
    csv_url = f'https://{os.environ["AWS_BUCKET"]}.s3.ap-south-1.amazonaws.com/{file_path}'
    sub_query = get_sub_query(channel_ids, channel_ids_path)
    start_date = start_date + ' 00:00:00'
    end_date = end_date + ' 23:59:59'
    query = f"""
                
                INSERT INTO FUNCTION s3('{csv_url}', '{os.environ["AWS_ACCESS_KEY"]}', '{os.environ["AWS_SECRET_ACCESS_KEY"]}', 'CSVWithNames')
                with 
                {sub_query},
                final as (
                    select post.shortcode id,
                                        argMax(publish_time, insert_timestamp) publishedAt,
                                        argMax(profile_id, insert_timestamp) channel_id,
                                        JSONExtractString(argMax(post.dimensions, post.insert_timestamp),'title') title,
                                        JSONExtractString(argMax(post.dimensions, post.insert_timestamp),'description') description,
                                        JSONExtractString(argMax(post.dimensions, post.insert_timestamp),'channel_title') channelTitle,
                                        JSONExtractString(argMax(post.dimensions, post.insert_timestamp),'keywords') tags,
                                        JSONExtractString(argMax(post.dimensions, post.insert_timestamp),'category_id') categoryId,
                                        JSONExtractString(argMax(post.dimensions, post.insert_timestamp),'audio_language') defaultAudioLanguage,         
                                        JSONExtractString(argMax(post.dimensions, post.insert_timestamp),'duration') duration,
                                        extractAll(duration, '([0-9]+)H')[1]                                       AS    hours,
                                      extractAll(duration, '([0-9]+)M')[1]                                       AS    minutes,
                                      extractAll(duration, '([0-9]+)S')[1]                                       AS    seconds,
                                      if(hours == '', 0, toInt64OrZero(hours)) * 60 * 60 +
                                      if(minutes == '', 0, toInt64OrZero(minutes)) * 60 + toInt64OrZero(seconds) AS    total_seconds,
                                        JSONExtractString(argMax(post.metrics, post.insert_timestamp),'view_count') viewCount,
                                        JSONExtractString(argMax(post.metrics, post.insert_timestamp),'likes') likeCount,
                                        JSONExtractString(argMax(post.metrics, post.insert_timestamp),'comments') commentCount
                                        from {os.environ["EVENT_DB_NAME"]}.post_log_events post
                                        where toStartOfDay(post.event_timestamp) >= today() - INTERVAL 1 DAY and 
                                        platform = 'YOUTUBE' and source = 'ytapi' and publish_time >= '{start_date}' and publish_time <= '{end_date}' and profile_id in ids
                                        group by shortcode
                                        )
                                        select 
                                        id,
                                       publishedAt,
                                       channel_id channelId,
                                       title,
                                       description,
                                       channelTitle,
                                       tags,
                                       categoryId,
                                       defaultAudioLanguage,
                                       total_seconds seconds,
                                       viewCount,
                                       likeCount,
                                       commentCount
                                        from final
                SETTINGS s3_truncate_on_insert = 1
    """
    client.query("SET format_csv_null_representation = ''")
    client.query(query)
    return csv_url


async def fetch_channel_monthly_stats_csv_report_generation(job_id: str, start_date: str, end_date: str,
                                                            channel_ids: list[str] = None,
                                                            channel_ids_path: str = None):
    previous_month_start_date = (datetime.strptime(start_date, '%Y-%m-%d') - relativedelta(months=1)).strftime('%Y-%m-%d')
    client = get_clickhouse_connection()

    current_date = datetime.today()
    csv_path = f"Data_request/fetch_channel_monthly_stats_csv/{current_date.year}/{current_date.month}/{current_date.day}/{job_id}_channel_monthly_stats.csv"
    csv_url = f'https://{os.environ["AWS_BUCKET"]}.s3.ap-south-1.amazonaws.com/{csv_path}'
    
    sub_query = get_sub_query(channel_ids, channel_ids_path)
    start_date = start_date + ' 00:00:00'
    end_date = end_date + ' 23:59:59'
    query = f"""
                
                INSERT INTO FUNCTION s3('{csv_url}', '{os.environ["AWS_ACCESS_KEY"]}', '{os.environ["AWS_SECRET_ACCESS_KEY"]}', 'CSVWithNames')
                                with 
                                {sub_query}
                                ,
                                old_stats as (select toStartOfMonth(parseDateTimeBestEffortOrZero(month)) as month,
                                        id                                                      channel_id,
                                        toInt64OrZero(lifetime_views)                        as lifetime_views,
                                        toInt64OrZero(lifetime_subscribers)                     lifetime_followers,
                                        toInt64OrZero(lifetime_uploads)                      as lifetime_uploads,
                                        1                                                       priority
                                from vidooly.channels_monthly_stats_20230219
                                where id in ids),
                    new_stats_vidooly as (select toStartOfMonth(parseDateTimeBestEffortOrZero(month)) as month,
                                                id                                                      channel_id,
                                                lifetime_views                                       as lifetime_views,
                                                lifetime_subscribers                                    lifetime_followers,
                                                lifetime_uploads                                     as lifetime_uploads,
                                                1                                                       priority
                                        from vidooly.channels_monthly_stats_20230709
                                        where id in ids),
                    new_stats_gcc as (select toStartOfMonth(event_timestamp)                                      month,
                                            handle                                                               channel_id,
                                            argMax(JSONExtractInt(metrics, 'view_count'), event_timestamp)       lifetime_views,
                                            argMax(JSONExtractInt(metrics, 'subscriber_count'), event_timestamp) lifetime_followers,
                                            argMax(JSONExtractInt(metrics, 'video_count'), event_timestamp)      lifetime_uploads,
                                            0                                                                    priority
                                    from {os.environ["EVENT_DB_NAME"]}.profile_log_events
                                    where platform = 'YOUTUBE'
                                        and handle in ids
                                        and event_timestamp >= date('2023-06-01')
                                        and JSONHas(metrics, 'view_count')
                                    group by month, channel_id
                                    order by month desc, channel_id desc),
                    monthly_stats as (select *
                                    from (select *
                                            from old_stats
                                            union all
                                            select *
                                            from new_stats_vidooly
                                            union all
                                            select *
                                            from new_stats_gcc)
                                    order by channel_id desc, month desc),
                    semi as (select month,
                                    channel_id,
                                    argMax(lifetime_views, priority)     as lifetime_views,
                                    argMax(lifetime_followers, priority) as lifetime_followers,
                                    argMax(lifetime_uploads, priority)   as lifetime_uploads
                            from monthly_stats
                            where month >= date('{previous_month_start_date}')
                                and month <=date('{end_date}')
                            group by month, channel_id),
                    final as (select month                                                                   date,
                                    channel_id,
                                    toUInt64(lifetime_uploads)                                              as lifetime_uploads,
                                    toUInt64(lifetime_views)                                                as lifetime_views,
                                    toUInt64(lifetime_followers)                                            as lifetime_followers,
                                    groupArray(date) OVER (PARTITION BY channel_id ORDER BY date ASC
                                        ROWS BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS                        date_values,
                                    date_values[1]                                                          p_date,
                                    abs(p_date - date) <= 31                                                last_month_available,
                                    groupArray(lifetime_followers) OVER (PARTITION BY channel_id ORDER BY date ASC
                                        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS                        followers_values,
                                    groupArray(lifetime_views) OVER (PARTITION BY channel_id ORDER BY date ASC
                                        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS                        views_values,
                                    groupArray(lifetime_uploads) OVER (PARTITION BY channel_id ORDER BY date ASC
                                        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS                        uploads_values,
                                    followers_values[1]                                                     followers_prev,
                                    views_values[1]                                                         views_prev,
                                    uploads_values[1]                                                       uploads_prev,
                                    if(last_month_available = 1, lifetime_followers - followers_prev, NULL) followers_change,
                                    if(last_month_available = 1, lifetime_views - views_prev, NULL)         views_change,
                                    if(last_month_available = 1, lifetime_uploads - uploads_prev, NULL)     uploads_change
                            from semi)
                select channel_id         id,
                    date               month,
                    lifetime_views,
                    lifetime_followers lifetime_subscribers,
                    lifetime_uploads,
                    views_change       netviewchange,
                    followers_change   subchange,
                    uploads_change     uploadchange
                from final
                where date >= date('2022-01-01') and date != date('{previous_month_start_date}')
                SETTINGS s3_truncate_on_insert=1
    """
    client.query("SET format_csv_null_representation = ''")
    client.query(query)
    return csv_url


async def fetch_channel_daily_stats_csv_report_generation(job_id: str, start_date: str, end_date: str,
                                                          channel_ids: list[str] = None, channel_ids_path: str = None):
    previous_day_start_date = (datetime.strptime(start_date, '%Y-%m-%d') - timedelta(days=1)).strftime('%Y-%m-%d')
    client = get_clickhouse_connection()
    current_date = datetime.today()
    csv_path = f"Data_request/fetch_channel_daily_stats_csv/{current_date.year}/{current_date.month}/{current_date.day}/{job_id}_channel_daily_stats.csv"

    csv_url = f'https://{os.environ["AWS_BUCKET"]}.s3.ap-south-1.amazonaws.com/{csv_path}'
    sub_query = get_sub_query(channel_ids, channel_ids_path)
    start_date = start_date + ' 00:00:00'
    end_date = end_date + ' 23:59:59'
    query = f"""
                INSERT INTO FUNCTION s3('{csv_url}', '{os.environ["AWS_ACCESS_KEY"]}', '{os.environ["AWS_SECRET_ACCESS_KEY"]}', 'CSVWithNames')
                                with 
                                {sub_query}
                                ,
                                old_stats as (select toStartOfDay(parseDateTimeBestEffortOrZero(month)) as created_at,
                                        id                                                    channel_id,
                                        toInt64OrZero(lifetime_views)                      as lifetime_views,
                                        toInt64OrZero(lifetime_subscribers)                   lifetime_followers,
                                        toInt64OrZero(lifetime_uploads)                    as lifetime_uploads,
                                        1                                                     priority
                                from vidooly.channels_monthly_stats_20230219
                                where id in ids),
                    new_stats_vidooly as (select toStartOfDay(parseDateTimeBestEffortOrZero(month)) as created_at,
                                                id                                                    channel_id,
                                                lifetime_views                                     as lifetime_views,
                                                lifetime_subscribers                                  lifetime_followers,
                                                lifetime_uploads                                   as lifetime_uploads,
                                                1                                                     priority
                                        from vidooly.channels_monthly_stats_20230709
                                        where id in ids),
                    new_stats_gcc as (select toStartOfDay(event_timestamp)                                        created_at,
                                            handle                                                               channel_id,
                                            argMax(JSONExtractInt(metrics, 'view_count'), event_timestamp)       lifetime_views,
                                            argMax(JSONExtractInt(metrics, 'subscriber_count'), event_timestamp) lifetime_followers,
                                            argMax(JSONExtractInt(metrics, 'video_count'), event_timestamp)      lifetime_uploads,
                                            0                                                                    priority
                                    from {os.environ["EVENT_DB_NAME"]}.profile_log_events
                                    where platform = 'YOUTUBE'
                                        and handle in ids
                                        and event_timestamp >= date('2023-06-01')
                                        and JSONHas(metrics, 'view_count')
                                    group by created_at, channel_id
                                    order by created_at desc, channel_id desc),
                    monthly_stats as (select *
                                    from (select *
                                            from old_stats
                                            union all
                                            select *
                                            from new_stats_vidooly
                                            union all
                                            select *
                                            from new_stats_gcc)
                                    order by channel_id desc, created_at desc),
                    final as (select created_at,
                                    channel_id,
                                    argMax(lifetime_views, priority)     as lifetime_views,
                                    argMax(lifetime_followers, priority) as lifetime_followers,
                                    argMax(lifetime_uploads, priority)   as lifetime_uploads
                            from monthly_stats
                            where created_at >= date('{previous_day_start_date}')
                                and created_at <= date('{end_date}')
                            group by created_at, channel_id),
                    final1 as (select created_at,
                                    channel_id,
                                    toUInt64(lifetime_uploads)                                              lifetime_uploads,
                                    toUInt64(lifetime_views)                                                lifetime_views,
                                    toUInt64(lifetime_followers)                                            lifetime_followers,
                                    groupArray(created_at) OVER (PARTITION BY channel_id ORDER BY created_at ASC
                                        ROWS BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS                        date_values,
                                    date_values[1]                                                          p_date,
                                    abs(date(p_date) - date(created_at)) <= 1                                     last_month_available,
                                    groupArray(lifetime_followers) OVER (PARTITION BY channel_id ORDER BY created_at ASC
                                        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS                        followers_values,
                                    groupArray(lifetime_views) OVER (PARTITION BY channel_id ORDER BY created_at ASC
                                        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS                        views_values,
                                    groupArray(lifetime_uploads) OVER (PARTITION BY channel_id ORDER BY created_at ASC
                                        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS                        uploads_values,
                                    followers_values[1]                                                     followers_prev,
                                    views_values[1]                                                         views_prev,
                                    uploads_values[1]                                                       uploads_prev,
                                    if(last_month_available = 1, lifetime_followers - followers_prev, NULL) followers_change,
                                    if(last_month_available = 1, lifetime_views - views_prev, NULL)         views_change,
                                    if(last_month_available = 1, lifetime_uploads - uploads_prev, NULL)     uploads_change
                                from final)
                select channel_id         id,
                    created_at,
                    lifetime_views,
                    lifetime_followers lifetime_subscribers,
                    lifetime_uploads,
                    views_change       netviewchange,
                    followers_change   subchange,
                    uploads_change     uploadchange
                from final1
                where created_at >= date('2022-01-01') and created_at != date('{previous_day_start_date}')
                SETTINGS s3_truncate_on_insert=1
    """
    client.query("SET format_csv_null_representation = ''")
    client.query(query)
    return csv_url

def get_clickhouse_connection():
    client = clickhouse_connect.get_client(host=os.environ["CLICKHOUSE_URL"],
                                           username=os.environ["CLICKHOUSE_USERNAME"],
                                           password=os.environ["CLICKHOUSE_PASSWORD"])
    return client

async def fetch_channel_demographic_csv_report_generation(job_id: str, channel_ids: list[str] = None,
                                                          channel_ids_path: str = None):
    client = get_clickhouse_connection()
    current_date = datetime.today()
    json_path = f"Data_request/fetch_channel_demographic_csv/{current_date.year}/{current_date.month}/{current_date.day}/{job_id}/channel_demographic.json"
    csv_path = f"Data_request/fetch_channel_demographic_csv/{current_date.year}/{current_date.month}/{current_date.day}/{job_id}/channel_demographic.csv"
    csv_url = f'https://{os.environ["AWS_BUCKET"]}.s3.ap-south-1.amazonaws.com/{csv_path}'
    json_url = f'https://{os.environ["AWS_BUCKET"]}.s3.ap-south-1.amazonaws.com/{json_path}'
    sub_query = get_sub_query(channel_ids, channel_ids_path)

    query = f"""
            with
                                    {sub_query}
                                    ,
                                    followers_data as (
                                        select channel_id , followers
                                        from dbt.mart_youtube_account
                                        where channel_id in ids
                                    ),
                                    semi as (
                                select platform_id id, CAST((age_groups, age_group_percs), 'Map(String, Float64)') demography,
                                    demography['M_13-17'] m_13_17,
                                        demography['M_18-24'] m_18_24,
                                        demography['M_25-34'] m_25_34,
                                        demography['M_35-44'] m_35_44,
                                        demography['M_45-54'] m_45_54,
                                        demography['M_55-64'] m_55_64,
                                        demography['M_65+'] m_65,
                                        demography['F_13-17'] f_13_17,
                                        demography['F_18-24'] f_18_24,
                                        demography['F_25-34'] f_25_34,
                                        demography['F_35-44'] f_35_44,
                                        demography['F_45-54'] f_45_54,
                                        demography['F_55-64'] f_55_64,
                                        demography['F_65+'] f_65,
                                    male_followers_perc male, female_followers_perc female
                                from dbt.mart_audience_info
                                where platform = 'YOUTUBE' and platform_id in ids
                                and (male_followers_perc>0 or female_followers_perc>0))
                                select
                                    m_13_17,
                                    m_18_24,
                                    m_25_34,
                                    m_35_44,
                                    m_45_54,
                                    m_55_64,
                                    m_65,
                                    f_13_17,
                                    f_18_24,
                                    f_25_34,
                                    f_35_44,
                                    f_45_54,
                                    f_55_64,
                                    f_65,
                                    male,
                                    female,
                                    fd.followers total_audience,
                                    id
                                from semi
                                inner join followers_data fd on fd.channel_id = id
    """

    query_json = f"""
                INSERT INTO FUNCTION s3('{json_url}', '{os.environ["AWS_ACCESS_KEY"]}', '{os.environ["AWS_SECRET_ACCESS_KEY"]}', 'JSONEachRow')
                {query}
                SETTINGS s3_truncate_on_insert=1
    """

    client.query(query_json)

    query_csv = f"""
        
        INSERT INTO FUNCTION s3('{csv_url}', '{os.environ["AWS_ACCESS_KEY"]}', '{os.environ["AWS_SECRET_ACCESS_KEY"]}', 'CSVWithNames')
            {query}
                SETTINGS s3_truncate_on_insert=1
    """
    client.query("SET format_csv_null_representation = ''")
    client.query(query_csv)

    return json_url, csv_url




async def fetch_channel_audience_geography_report_generation(job_id: str, channel_ids: list[str] = None,
                                                          channel_ids_path: str = None):
    client = get_clickhouse_connection()

    current_date = datetime.today()
    file_path = f"Data_request/fetch_channel_audience_geography_report_generation/{current_date.year}/{current_date.month}/{current_date.day}/{job_id}_channel_geography.csv"
    csv_url = f'https://{os.environ["AWS_BUCKET"]}.s3.ap-south-1.amazonaws.com/{file_path}'
    sub_query = get_sub_query(channel_ids, channel_ids_path)
    query = f"""
                
                INSERT INTO FUNCTION s3('{csv_url}',
                                        '{os.environ["AWS_ACCESS_KEY"]}', '{os.environ["AWS_SECRET_ACCESS_KEY"]}', 'CSVWithNames')
                with
                    {sub_query},
                    channels as (
                        select
                            s.subscriber_id as subscriber_id,
                            s.channel_id as channel_id,
                            s.subscriber_country as subscriber_country
                            FROM dbt.mart_yt_subscriber AS s
                            where channel_id in ids and subscriber_country not in ('', 'MISSING')
                    )
                    SELECT s.channel_id                                                               as channel_id,
                        s.subscriber_country                                                       AS country,
                        round(uniqExact(s.subscriber_id) * 100.0 /
                                (sum(uniqExact(s.subscriber_id)) over (partition by channel_id)), 2) AS percentage
                    FROM channels AS s
                    GROUP BY s.channel_id, s.subscriber_country
                    ORDER BY s.channel_id
                    SETTINGS s3_truncate_on_insert = 1
    """
    client.query("SET format_csv_null_representation = ''")
    client.query(query)
    return csv_url

async def fetch_channel_engagement_report_generation(job_id: str,start_date: str,end_date: str,filter_by: str, channel_ids: list[str] = None,
                                                          channel_ids_path: str = None):
    client = get_clickhouse_connection()
    current_date = datetime.today()
    file_path = f"Data_request/fetch_channel_engagement_report_generation/{current_date.year}/{current_date.month}/{current_date.day}/{job_id}_channel_engagement.csv"
    
    start_date = start_date + ' 00:00:00'
    end_date = end_date + ' 23:59:59'
    
    if channel_ids_path:
        sub_query = f" select * from url('{channel_ids_path}', 'CSVWithNames')"
    else:
        sub_query = f"SELECT arrayJoin({channel_ids}) AS id"
    channel_ids = []
    channel_ids_obj = client.query(sub_query)
    for row in channel_ids_obj.result_rows:
        channel_ids.append(row[0])

    final_data = []
    #  filter_by VALUES "channel_gcc_category" OR "channel_yt_category"
    for channel_id in channel_ids:
        query = f"""
                with subscriberList AS (SELECT channel_id, subscriber_id
                                        from dbt.mart_yt_subscriber mys
                                        where channel_id = '{channel_id}'
                                        group by channel_id, subscriber_id)
                , subcribers_channel_categories AS (SELECT subscriberList.channel_id as channel_id,
                                                            mys2.channel_id              subscriber_channel_id,
                                                            mys2.{filter_by},
                                                            uniqExact(channel_id)     as subscribers_count
                                                    from dbt.mart_yt_subscriber mys2
                                                    inner join subscriberList on subscriberList.subscriber_id = mys2.subscriber_id
                                                    WHERE mys2.timestamp >= '{start_date}' and mys2.timestamp <= '{end_date}' and channel_yt_category not in ('', 'MISSING')
                                                    GROUP BY channel_id, subscriber_channel_id, {filter_by})
                Select channel_id,
                    {filter_by} as category,
                    COUNT(*) as activityCount
                from subcribers_channel_categories
                group by channel_id, {filter_by}
                order by channel_id"""
        
        engagement_data_obj = client.query(query)
        for row in engagement_data_obj.result_rows:
            final_data.append(list(row))

    filename = job_id+"_channel_engagement.csv"
    current_dir = os.path.dirname(os.path.abspath(__file__))
    dir_file_path = os.path.join(current_dir, filename+".csv")
    
    with open(dir_file_path, "a", newline='') as f:
        writer = csv.writer(f)
        headers = ['Channel Id', filter_by, 'Subscriber Count']
        writer.writerow(headers)
        for row in final_data:
            writer.writerow(row)

    session = boto3.Session(
        aws_access_key_id=os.environ["AWS_ACCESS_KEY"],
        aws_secret_access_key=os.environ["AWS_SECRET_ACCESS_KEY"]
    )
    # Pushing File To s3
    s3 = session.client('s3')
    bucket_name = os.environ["AWS_BUCKET"]
    file_path = f"Data_request/fetch_channel_engagement_report_generation/{current_date.year}/{current_date.month}/{current_date.day}/{filename}"
    csv_url = f'https://{os.environ["AWS_BUCKET"]}.s3.ap-south-1.amazonaws.com/{file_path}'
    s3.upload_file(dir_file_path, bucket_name, file_path)

    # removing File From OS Directory
    if(os.path.exists(dir_file_path) and os.path.isfile(dir_file_path)): 
        os.remove(dir_file_path)
    return csv_url


async def get_channel_ids(video_ids: list[str]):
    client = get_clickhouse_connection()
    query = f"""
        select profile_id from {os.environ["EVENT_DB_NAME"]}.post_log_events where shortcode in {video_ids} group by profile_id
    """
    
    result = client.query(query)
    channel_ids = [row[0] for row in result.result_rows]
    return channel_ids

async def get_channel_videos_merged_data(video_ids: list[str], channel_ids: list[str]):
    client = get_clickhouse_connection()
    
    query = f"""
            with posts as (select *
                        from {os.environ["EVENT_DB_NAME"]}.post_log_events
                        where 
                            platform = 'YOUTUBE'
                            and source = 'ytapi' and
                            shortcode in
                                {video_ids}
                        ),
                final_posts as (select JSONExtractString(argMax(dimensions, event_timestamp), 'duration')    duration,
                                        JSONExtractString(argMax(dimensions, event_timestamp), 'is_licensed') licensedContent,
                                        JSONExtractString(argMax(dimensions, event_timestamp), 'definition')  definition,
                                        JSONExtractInt(argMax(metrics, event_timestamp), 'comments')          commentCount,
                                        JSONExtractInt(argMax(metrics, event_timestamp), 'view_count')        viewCount,
                                        JSONExtractInt(argMax(metrics, event_timestamp), 'likes')             likeCount,
                                        shortcode                                                             vid,
                                        JSONExtractString(argMax(dimensions, event_timestamp), 'category_id') categoryId,
                                        argMax(profile_id, event_timestamp)                                   channelId,
                                        argMax(toString(publish_time), event_timestamp)                publishedAt
                                from posts
                                group by shortcode),
                channels as (select JSONExtractString(argMax(dimensions, event_timestamp), 'country')           country,
                                    profile_id                                         id,
                                    JSONExtractInt(argMax(metrics, event_timestamp), 'hidden_subscriber_count') hidden_subscriber_count,
                                    JSONExtractInt(argMax(metrics, event_timestamp), 'subscriber_count')        subscriber_count,
                                    JSONExtractInt(argMax(metrics, event_timestamp), 'view_count')              view_count,
                                    JSONExtractInt(argMax(metrics, event_timestamp), 'video_count')              video_count,
                                    toString(toDateTimeOrZero(JSONExtractString(argMax(dimensions, event_timestamp), 'taken_at_timestamp')))          taken_at_timestamp
                            from {os.environ["EVENT_DB_NAME"]}.profile_log_events
                            where event_timestamp >= now() - interval 1 day and
                            platform = 'YOUTUBE'
                            and source = 'ytapi' and
                            profile_id in {channel_ids}
                            group by profile_id
                            )
            select duration,
                licensedContent,
                definition,
                commentCount,
                viewCount,
                likeCount,
                vid,
                categoryId,
                channelId,
                publishedAt,
                hidden_subscriber_count Channel_hiddenSubscriberCount,
                view_count              Channel_viewCount,
                subscriber_count        Channel_subscriberCount,
                video_count             Channel_videoCount,
                country                 Channel_country,
                taken_at_timestamp     ChannelPublish_date
            from final_posts
                    inner join channels on channels.id = final_posts.channelId
    """
    
    result = client.query(query)
    response = []
    columns = ["duration", "licensedContent", "definition", "commentCount", "viewCount", "likeCount", "vid", "categoryId", "channelId", "publishedAt",
                "Channel_hiddenSubscriberCount", "Channel_viewCount", "Channel_subscriberCount", "Channel_videoCount", "Channel_country", "ChannelPublish_date"]
    for row in result.result_rows:
        response.append(dict(zip(columns, row)))
    print(response)
    return response