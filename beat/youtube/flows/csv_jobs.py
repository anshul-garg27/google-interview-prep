from datetime import datetime
import os
from utils.request import make_request
from youtube.flows.refresh_profile import _batchify, refresh_yt_posts_by_channel_id, refresh_yt_profiles
from youtube.tasks.csv_report_generation import fetch_channel_audience_geography_report_generation, fetch_channel_daily_stats_csv_report_generation, \
    fetch_channel_demographic_csv_report_generation, fetch_channel_details_csv, fetch_channel_engagement_report_generation, \
    fetch_channel_monthly_stats_csv_report_generation, fetch_channel_videos_csv_report, \
    get_channel_details_from_mart, get_channel_ids, get_channel_videos_merged_data
from youtube.flows.refresh_profile import fetch_post_ids_by_playlist_id, refresh_yt_posts
from youtube.tasks.csv_report_generation import fetch_playlist_videos_csv_report_generation, \
    fetch_videos_details_csv_report_generation, get_ids_from_csv_file
import csv
import boto3
import asyncio
from loguru import logger
from youtube.flows.refresh_profile import refresh_yt_posts_by_search
from youtube.tasks.csv_report_generation import fetch_keyword_channel_csv_report_generation, \
    fetch_keyword_video_mapping_report_generation, fetch_keyword_videos_csv_report_generation
import re

def diff(li1, li2):
    li_dif = [i for i in li1 + li2 if i not in li1 or i not in li2]
    return li_dif


async def fetch_channel_videos_csv(scrape_id: str, start_date: str, end_date: str, channel_ids: list[str] = None,
                                   channel_ids_path: str = None, session=None) -> dict:
    if channel_ids_path is not None:
        logger.debug(f"channel_ids_path = {channel_ids_path}")
        channel_ids = await get_ids_from_csv_file(channel_ids_path)

    max_posts = 500
    filters = {'start_date': start_date, 'end_date': end_date}
    csv_path = None
    if channel_ids is not None:
        for channel_id in channel_ids:
            try:
                await refresh_yt_posts_by_channel_id(channel_id, filters, max_posts, session=session)
            except Exception as e:
                logger.error(e)
        job_id = f"beat_scrape_{scrape_id}"
        await asyncio.sleep(60)
        csv_path = await fetch_channel_videos_csv_report(job_id, filters['start_date'], filters['end_date'],
                                                         channel_ids, channel_ids_path)

    return {"channel_video_csv": csv_path}


async def fetch_channels_csv(scrape_id: str, channel_ids: list[str] = None, channel_ids_path: str = None,
                             session=None) -> dict:
    if channel_ids_path is not None:
        channel_ids = await get_ids_from_csv_file(channel_ids_path)

    csv_path = None
    if channel_ids is not None:
        mart_channel_ids = await get_channel_details_from_mart(channel_ids, channel_ids_path)
        event_channel_ids = diff(channel_ids, mart_channel_ids)
        batches = _batchify(event_channel_ids, max_len=40)
        for batch in batches:
            try:
                await refresh_yt_profiles(batch, session=session)
            except Exception as e:
                logger.error(e)
                continue
        job_id = f"beat_scrape_{scrape_id}"
        await asyncio.sleep(30)
        csv_path = await fetch_channel_details_csv(job_id, channel_ids, channel_ids_path)

    return {"channels_csv": csv_path}


async def fetch_playlist_videos_csv(scrape_id: str, playlist_ids: list[str] = None, playlist_ids_path: str = None,
                                    session=None) -> dict:
    job_id = f"beat_scrape_{scrape_id}"
    if playlist_ids_path is not None:
        playlist_ids = await get_ids_from_csv_file(playlist_ids_path)

    csv_path = None
    if playlist_ids is not None:
        for playlist_id in playlist_ids:
            post_ids = await fetch_post_ids_by_playlist_id(playlist_id, session=session)

            batches = _batchify(post_ids, max_len=40)
            for batch in batches:
                try:
                    await refresh_yt_posts(batch, playlist_id=playlist_id, session=session)
                except Exception as e:
                    logger.error(e)
                    continue
        await asyncio.sleep(60)
        csv_path = await fetch_playlist_videos_csv_report_generation(job_id, playlist_ids=playlist_ids,
                                                                     playlist_ids_path=playlist_ids_path)
    return {"playlist_videos_csv": csv_path}


async def fetch_keyword_videos_csv(start_date: str, end_date: str, scrape_id: str, keywords: list[str]=None, keywords_list_path: str=None, frequency: int = None, session=None):
    job_id = f"beat_scrape_{scrape_id}"
    if keywords_list_path is not None:
        keywords = await get_ids_from_csv_file(keywords_list_path)
        
    keywords = [re.sub(r'[^\w\s]', '', keyword) for keyword in keywords]
    await refresh_yt_posts_by_search(job_id, keywords, start_date, end_date, False, frequency, session=session)
    await asyncio.sleep(60)
    videos_list_path = await fetch_keyword_videos_csv_report_generation(keywords, job_id, start_date, end_date)

    channel_list_path = await fetch_keyword_channel_csv_report_generation(keywords, start_date, end_date, job_id)

    video_keyword_mapping_path = await fetch_keyword_video_mapping_report_generation(keywords, job_id, start_date,
                                                                                     end_date)
    result = {
        "channel_list_path": channel_list_path,
        "video_list_path": videos_list_path,
        "keyword_video_mapping_path": video_keyword_mapping_path
    }
    return result


async def fetch_video_details_csv(scrape_id: str, video_ids: list[str] = None, video_ids_path: str = None,
                                  session=None) -> dict:
    if video_ids_path is not None:
        video_ids = await get_ids_from_csv_file(video_ids_path)

    csv_path = None
    if video_ids is not None:
        batches = _batchify(video_ids, max_len=40)
        for batch in batches:
            try:
                await refresh_yt_posts(batch, session=session)
            except Exception as e:
                logger.error(e)
                continue    
        await asyncio.sleep(60)
        job_id = f"beat_scrape_{scrape_id}"
        csv_path = await fetch_videos_details_csv_report_generation(job_id, video_ids, video_ids_path)

    return {"video_details_csv": csv_path}


async def fetch_channel_monthly_stats_csv(scrape_id: str, start_date: str, end_date: str, channel_ids: list[str] = None,
                                          channel_ids_path: str = None, session=None) -> dict:
    csv_path = None
    if channel_ids is not None or channel_ids_path is not None:
        job_id = f"beat_scrape_{scrape_id}"
        csv_path = await fetch_channel_monthly_stats_csv_report_generation(job_id, start_date, end_date, channel_ids,
                                                                           channel_ids_path)

    return {"channel_monthly_stats_csv": csv_path}


async def fetch_channel_daily_stats_csv(scrape_id: str, start_date: str, end_date: str, channel_ids: list[str] = None,
                                        channel_ids_path: str = None,
                                        session=None) -> dict:

    csv_path = None
    if channel_ids is not None  or channel_ids_path is not None:
        job_id = f"beat_scrape_{scrape_id}"
        csv_path = await fetch_channel_daily_stats_csv_report_generation(job_id, start_date, end_date, channel_ids,
                                                                         channel_ids_path)

    return {"channel_daily_stats_csv": csv_path}


async def fetch_channel_demographic_csv(scrape_id: str, channel_ids: list[str] = None, channel_ids_path: str = None,
                                        session=None) -> dict:
    csv_path, json_path = None, None
    if channel_ids is not None  or channel_ids_path is not None:
        job_id = f"beat_scrape_{scrape_id}"
        json_path, csv_path = await fetch_channel_demographic_csv_report_generation(job_id, channel_ids,
                                                                                    channel_ids_path)
    
    return {"channel_demographic_csv": csv_path, "channel_demographic_json": json_path}


async def fetch_channel_audience_geography_csv(scrape_id: str, channel_ids: list[str] = None, channel_ids_path: str = None,
                                        session=None) -> dict:
    csv_path = None
    if channel_ids is not None  or channel_ids_path is not None:
        job_id = f"beat_scrape_{scrape_id}"
        csv_path = await fetch_channel_audience_geography_report_generation(job_id, channel_ids,
                                                                                    channel_ids_path)
    return {"channel_audience_geography_csv": csv_path}


async def fetch_channel_engagement_csv(scrape_id: str, start_date: str, end_date: str,filter_by :str, channel_ids: list[str] = None, channel_ids_path: str = None,session=None) -> dict:
    csv_path = None
    if channel_ids is not None  or channel_ids_path is not None:
        job_id = f"beat_scrape_{scrape_id}"
        csv_path = await fetch_channel_engagement_report_generation(job_id,start_date,end_date,filter_by, channel_ids, channel_ids_path)
    
    return {"channel_engagement_csv": csv_path}

async def fetch_ad_vs_organic_csv(scrape_id: str, video_ids: list[str] = None, video_ids_path: str = None, session=None) -> dict:
    
    if video_ids_path is not None:
        video_ids = await get_ids_from_csv_file(video_ids_path)

    if video_ids is not None:
        batches = _batchify(video_ids, max_len=40)
        for batch in batches:
            try:
                await refresh_yt_posts(batch, session=session)
            except Exception as e:
                logger.error(e)
                continue
        await asyncio.sleep(60)
        channel_ids = await get_channel_ids(video_ids)
        batches = _batchify(channel_ids, max_len=40)
        for batch in batches:
            try:
                await refresh_yt_profiles(batch, session=session)
            except Exception as e:
                logger.error(e)
                continue
        await asyncio.sleep(60)
        result = await get_channel_videos_merged_data(video_ids, channel_ids)
        filename = f"beat_scrape_{scrape_id}_ad_vs_organic.csv"
        current_dir = os.path.dirname(os.path.abspath(__file__))
        dir_file_path = os.path.join(current_dir, filename)
        
        with open(dir_file_path, 'w', newline="") as f:
            title = "videoId,ad_predicted,view,Channel_videoCount,duration,Channel_SubscriberCount,DislikeCount,Videoage,likecount,commentcount".split(",") 
            cw = csv.DictWriter(f,title,delimiter=',', quotechar='|', quoting=csv.QUOTE_MINIMAL)
            cw.writeheader()
            for index in range(0, len(result), 10):
                response = await get_ad_vs_organic(result[index:index+10])
                cw.writerows(response['data'])

    session = boto3.Session(
        aws_access_key_id=os.environ["AWS_ACCESS_KEY"],
        aws_secret_access_key=os.environ["AWS_SECRET_ACCESS_KEY"]
    )
    
    # Pushing File To s3
    s3 = session.client('s3')
    bucket_name = os.environ["AWS_BUCKET"]
    current_date = datetime.today()
    file_path = f"Data_request/fetch_ad_vs_organic_csv/{current_date.year}/{current_date.month}/{current_date.day}/{filename}"
    csv_url = f'https://{os.environ["AWS_BUCKET"]}.s3.ap-south-1.amazonaws.com/{file_path}'
    s3.upload_file(dir_file_path, bucket_name, file_path)

    # removing File From OS Directory
    if(os.path.exists(dir_file_path) and os.path.isfile(dir_file_path)): 
        os.remove(dir_file_path)
            
    return {"video_details_csv": csv_url}


async def get_ad_vs_organic(data):
    headers = {
        'Content-Type': 'application/json',
    }
    url = f'{os.environ["AD_VS_ORGANIC"]}ad_vs_organic'
    response = await make_request("POST", url=url, headers=headers, json=data)
    return response.json()