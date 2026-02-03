import asyncio
import datetime
import multiprocessing
import os
import sys
import time
import traceback
from asyncio import Semaphore
from random import randint

import uvloop
from dotenv import load_dotenv
from loguru import logger
from sqlalchemy import text
from sqlalchemy.ext.asyncio import create_async_engine, AsyncConnection

from core.amqp.amqp import listener
from core.amqp.models import AmqpListener
from core.flows.scraper import perform_scrape_task
from core.helpers.session import sessionize
from core.models.models import ScrapeRequest, ScrapeRequestLog
from core.models.models import ScrapeRequestLog as ScrapeRequestLogModel
from credentials.listener import credential_validate, upsert_credential_from_identity, fetch_keyword_collection, \
    fetch_sentiment_report, sentiment_extraction
from utils.request import make_scrape_request_log_event

_whitelist = {
              'refresh_profile_custom': {'no_of_workers': 1, 'no_of_concurrency': 2},
              'refresh_profile_by_handle': {'no_of_workers': 10, 'no_of_concurrency': 5},
              'refresh_profile_basic': {'no_of_workers': 5, 'no_of_concurrency': 5},
              'refresh_profile_by_profile_id': {'no_of_workers': 1, 'no_of_concurrency': 5},
              'refresh_post_by_shortcode': {'no_of_workers': 1, 'no_of_concurrency': 5},
              'refresh_tagged_posts_by_profile_id': {'no_of_workers': 1, 'no_of_concurrency': 5},
              'refresh_yt_profiles': {'no_of_workers': 10, 'no_of_concurrency': 5},
              'refresh_yt_posts': {'no_of_workers': 1, 'no_of_concurrency': 5},
              'refresh_yt_profile_insights': {'no_of_workers': 1, 'no_of_concurrency': 5},
              'asset_upload_flow': {'no_of_workers': 15, 'no_of_concurrency': 5},
              'asset_upload_flow_stories': {'no_of_workers': 3, 'no_of_concurrency': 5},
              'refresh_profile_insights': {'no_of_workers': 1, 'no_of_concurrency': 5},
              'fetch_profile_followers': {'no_of_workers': 3, 'no_of_concurrency': 5},
              'fetch_profile_following': {'no_of_workers': 1, 'no_of_concurrency': 15},
              'fetch_hashtag_posts': {'no_of_workers': 1, 'no_of_concurrency': 5},
              'refresh_stories_posts': {'no_of_workers': 1, 'no_of_concurrency': 5},
              'refresh_story_insights': {'no_of_workers': 1, 'no_of_concurrency': 5},
              'refresh_post_insights': {'no_of_workers': 1, 'no_of_concurrency': 5},
              'refresh_yt_posts_by_channel_id': {'no_of_workers': 5, 'no_of_concurrency': 5},
              'refresh_yt_posts_by_playlist_id': {'no_of_workers': 1, 'no_of_concurrency': 5},
              'refresh_yt_posts_by_genre': {'no_of_workers': 1, 'no_of_concurrency': 5},
              'refresh_yt_posts_by_search': {'no_of_workers': 1, 'no_of_concurrency': 5},
              'refresh_orders_by_store': {'no_of_workers': 1, 'no_of_concurrency': 5},
              'refresh_yt_post_type': {'no_of_workers': 1, 'no_of_concurrency': 5},
              'fetch_post_comments': {'no_of_workers': 1, 'no_of_concurrency': 5},
              'fetch_post_likes': {'no_of_workers': 1, 'no_of_concurrency': 5},
              'fetch_yt_post_comments': {'no_of_workers': 1, 'no_of_concurrency': 5},
              'refresh_instagram_gpt_data_base_gender': {'no_of_workers': 3, 'no_of_concurrency': 5},
              'refresh_instagram_gpt_data_base_location': {'no_of_workers': 3, 'no_of_concurrency': 5},
              'refresh_instagram_gpt_data_base_categ_lang_topics': {'no_of_workers': 3, 'no_of_concurrency': 5},
              'refresh_instagram_gpt_data_audience_age_gender': {'no_of_workers': 3, 'no_of_concurrency': 5},
              'refresh_instagram_gpt_data_audience_cities': {'no_of_workers': 3, 'no_of_concurrency': 5},
              'refresh_story_posts_by_profile_id': {'no_of_workers': 1, 'no_of_concurrency': 5},
              'fetch_channels_csv': {'no_of_workers': 1, 'no_of_concurrency': 5},
              'fetch_channel_videos_csv': {'no_of_workers': 1, 'no_of_concurrency': 5},
              'fetch_playlist_videos_csv': {'no_of_workers': 1, 'no_of_concurrency': 5},
              'fetch_keyword_videos_csv': {'no_of_workers': 1, 'no_of_concurrency': 5},
              'fetch_video_details_csv': {'no_of_workers': 1, 'no_of_concurrency': 5},
              'fetch_channel_monthly_stats_csv': {'no_of_workers': 1, 'no_of_concurrency': 5},
              'fetch_channel_daily_stats_csv': {'no_of_workers': 1, 'no_of_concurrency': 5},
              'fetch_channel_demographic_csv': {'no_of_workers': 1, 'no_of_concurrency': 5},
              'refresh_yt_activities': {'no_of_workers': 1, 'no_of_concurrency': 5},
              'fetch_channel_audience_geography_csv':{'no_of_workers': 1, 'no_of_concurrency': 5},
              'fetch_channel_engagement_csv':{'no_of_workers': 1, 'no_of_concurrency': 5},
              'refresh_yt_profile_relationship': {'no_of_workers': 1, 'no_of_concurrency': 5},
              'refresh_instagram_gpt_data_gender_location_lang': {'no_of_workers': 5, 'no_of_concurrency': 5},
              'refresh_story_posts_by_profile_id': {'no_of_workers': 1, 'no_of_concurrency': 5},
              'fetch_ad_vs_organic_csv': {'no_of_workers': 1, 'no_of_concurrency': 5}
            }

amqp_listeners = [
    AmqpListener("beat.dx",
                 "direct",
                 "credentials_validate",
                 "credentials_validate_q",
                 5,
                 10,
                 credential_validate),
    AmqpListener("identity.dx",
                 "direct",
                 "new_access_token_rk",
                 "identity_token_q",
                 5,
                 10,
                 upsert_credential_from_identity),
    AmqpListener("beat.dx",
                 "direct",
                 "keyword_collection_rk",
                 "keyword_collection_q",
                 5,
                 1,
                 fetch_keyword_collection),
    AmqpListener("beat.dx",
                     "direct",
                     "post_activity_log_bulk",
                     "sentiment_analysis_q",
                     5,
                     1,
                     sentiment_extraction),
    AmqpListener("beat.dx",
                 "direct",
                 "sentiment_collection_report_in_rk",
                 "sentiment_collection_report_in_q",
                 5,
                 1,
                 fetch_sentiment_report)
]


# def set_cpu_affinity(cpu_id: int):
#     proc = psutil.Process()
#     n_cpus = psutil.cpu_count()
#     if cpu_id >= n_cpus:
#         cpu_id = floor(random() * n_cpus)
#     proc.cpu_affinity([cpu_id])


def looper(id: int, limit: Semaphore, flows: list, max_conc=5) -> None:
    # TODO: Only do this on non-mac env / stage/prod
    # set_cpu_affinity(id - 1) # Worker ID starts at 1, CPU ID starts at 0
    while True:
        try:
            with asyncio.Runner(loop_factory=uvloop.new_event_loop) as runner:
                runner.run(poller(id, limit, flows, max_conc))
        except Exception as e:
            logger.error(f"Error: {e}")
            msg = traceback.format_exc()
            logger.error(f"Error Message - {msg}")

async def poller(id: int, limit: Semaphore, flows: list, max_conc=5) -> None:
    logger.debug("Starting Worker")
    engine = create_async_engine(os.environ["PG_URL"], isolation_level="AUTOCOMMIT", echo=False, pool_size=max_conc, max_overflow=5)
    conn = await engine.connect()
    while True:
        await poll(id, limit, flows, conn, engine)
        await asyncio.sleep(1)


async def check_and_cancel_long_running_tasks(background_tasks: set):
    logger.debug(f"Checking and cancelling {len(background_tasks)} long running tasks")
    current_time = time.time()
    task_timeout = 60 * 10 * 1000  # Timeout tasks after 10 minutes
    canceled_tasks = 0
    completed_tasks = 0
    for start_time, task in background_tasks:
        if task.done():
            background_tasks.remove((start_time, task))
            completed_tasks += 1
        elif current_time - start_time > task_timeout:
            task.cancel()
            background_tasks.remove((start_time, task))
            canceled_tasks += 1
    logger.debug(f"Completed {completed_tasks} tasks. Cancelled {canceled_tasks} tasks.")

async def poll(id: int, limit: Semaphore, flows: list, conn: AsyncConnection, engine) -> None:
    background_tasks = set()
    while limit.locked():
        logger.debug("Workers are busy, will try again in a bit.")
        await asyncio.sleep(30)
    for flow in flows:
        values = {'flow': flow}
        statement = text("""
            update scrape_request_log
                set status='PROCESSING'
                where id IN (
                    select id from scrape_request_log e
                    where status = 'PENDING' and
                    flow = :flow
                    for update skip locked
                    limit 1)
            RETURNING *
        """)
        rs = await conn.execute(statement, values)
        for row in rs:
            logger.debug(str(id) + ": Found Work - " + str(row.id))
            task = asyncio.create_task(perform_task(conn, row, limit, session_engine=engine))
            background_tasks.add((int(time.time()*1000), task))
        if randint(0, 1000) < 10:
            await check_and_cancel_long_running_tasks(background_tasks)

@sessionize
async def perform_task(con: AsyncConnection, row: any, limit: Semaphore, session=None, session_engine=None) -> None:
    logger.debug(f"TASK-{row.id}::Acquiring Locks")
    await limit.acquire()
    logger.debug(f"TASK-{row.id}::Starting Execution")
    body = ScrapeRequest(
        flow=row.flow,
        platform=row.platform,
        params=row.params,
        retry_count=row.retry_count,
        priority=1
    )
    
    now = datetime.datetime.now()
    body.status="PROCESSING"
    body.event_timestamp = now.strftime("%Y-%m-%d %H:%M:%S")
    body.picked_at = now.strftime("%Y-%m-%d %H:%M:%S")
    body.expires_at = now + datetime.timedelta(hours=1)
    await make_scrape_request_log_event(str(row.id), body)
    
    try:
        logger.debug(f"TASK-{row.id}::Creating Task")
        data = await execute(row, session=session)
        logger.debug(f"TASK-{row.id}::Completed Task")
        statement = text("""
                        update scrape_request_log
                            set status='COMPLETE', scraped_at=now(), data=:data
                            where id = %s
                    """ % row.id)
        await con.execute(statement, parameters={'data': str(data)})
        
        now = datetime.datetime.now()
        body.status="COMPLETE"
        body.event_timestamp = now.strftime("%Y-%m-%d %H:%M:%S")
        body.picked_at = now.strftime("%Y-%m-%d %H:%M:%S")
        body.expires_at = now + datetime.timedelta(hours=1)
        await make_scrape_request_log_event(str(row.id), body)
        
        logger.debug(f"TASK-{row.id}::Finalized Completion")
    except Exception as e:
        logger.debug(f"TASK-{row.id}::Failing Execution")
        msg = traceback.format_exc()
        error = f"{str(e)}"
        logger.debug(f"TASK-{row.id}::Error - {error}")
        logger.debug(f"TASK-{row.id}::Error Message - {msg}")
        statement = text("""
                                update scrape_request_log
                                    set status='FAILED', scraped_at=now(), data=:data
                                    where id = :id
                                """)
        await con.execute(statement, parameters={'data': error, 'id': row.id})
        
        now = datetime.datetime.now()
        body.status="FAILED"
        body.event_timestamp = now.strftime("%Y-%m-%d %H:%M:%S")
        body.picked_at = now.strftime("%Y-%m-%d %H:%M:%S")
        body.expires_at = now + datetime.timedelta(hours=1)
        await make_scrape_request_log_event(str(row.id), body)
        
        logger.debug(f"TASK-{row.id}::Finalized Failure")
    finally:
        logger.debug(f"TASK-{row.id}::Releasing Execution Locks")
        limit.release()


@sessionize
async def execute(scrape_log: ScrapeRequestLog, session=None) -> None | dict:
    scrape_request = ScrapeRequestLogModel(scrape_log.id,
                                           scrape_log.flow,
                                           scrape_log.params)
    return await perform_scrape_task(scrape_request, session=session)


def start_workers() -> None:
    workers = []
    start_listeners(workers)
    start_scrape_log_workers(workers)
    for w in workers:
        w.join()


def start_listeners(workers: list) -> None:
    for amqp_listener_config in amqp_listeners:
        for _ in range(amqp_listener_config.workers):
            w = multiprocessing.Process(target=listener, args=(amqp_listener_config,))
            w.start()
            workers.append(w)


def start_scrape_log_workers(workers: list) -> None:
    single_worker = bool(int(os.environ["SINGLE_WORKER"]))
    if single_worker:
        max_conc = int(os.environ["SCRAPE_LOG_MAX_CONCURRENCY_PER_LOOP"])
        limit = asyncio.Semaphore(max_conc)
        looper(1, limit, list(_whitelist.keys()), max_conc)
    else:
        for flow in _whitelist:
            num_process = _whitelist[flow]['no_of_workers']
            for i in range(num_process):
                max_conc = int(_whitelist[flow]['no_of_concurrency'])
                limit = asyncio.Semaphore(max_conc)
                w = multiprocessing.Process(target=looper, args=(i + 1, limit, [flow], max_conc))
                w.start()
                workers.append(w)

# def write_pid():
#     with open('beat.pid', 'w', encoding='utf-8') as f:
#         f.write(str(os.getpid()))

if __name__ == "__main__":
    load_dotenv()
    # write_pid()
    logger.remove()
    logger.add(sys.stderr, format="{time} |  {level}   | {process}   | {name}:{function}:{line} - {message}", level=os.environ["LOG_LEVEL"])
    start_workers()
