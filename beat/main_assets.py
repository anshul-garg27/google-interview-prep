from random import randint
import asyncpg
import asyncio
import multiprocessing
import os
import sys
import traceback
from asyncio import Semaphore
import time
import uvloop
from dotenv import load_dotenv
from loguru import logger
from sqlalchemy import text
from sqlalchemy.ext.asyncio import create_async_engine, AsyncConnection

from core.amqp.amqp import listener
from core.amqp.models import AmqpListener
from core.flows.scraper import perform_scrape_task
from core.helpers.session import sessionize
from core.models.models import ScrapeRequestLog
from core.models.models import ScrapeRequestLog as ScrapeRequestLogModel
from credentials.listener import credential_validate, upsert_credential_from_identity

_whitelist = {
              'asset_upload_flow': {'no_of_workers': 50, 'no_of_concurrency':100},
             }
amqp_listeners = [
]


# def set_cpu_affinity(cpu_id: int):
#     proc = psutil.Process()
#     n_cpus = psutil.cpu_count()
#     if cpu_id >= n_cpus:
#         cpu_id = floor(random() * n_cpus)
#     proc.cpu_affinity([cpu_id])


def looper(id: int, limit: Semaphore, flows: list) -> None:
    # TODO: Only do this on non-mac env / stage/prod
    # set_cpu_affinity(id - 1) # Worker ID starts at 1, CPU ID starts at 0
    while True:
        try:
            with asyncio.Runner(loop_factory=uvloop.new_event_loop) as runner:
                runner.run(poller(id, limit, flows))
        except Exception as e:
            logger.error(f"Error: {e}")
            msg = traceback.format_exc()
            logger.error(f"Error Message - {msg}")


async def poller(id: int, limit: Semaphore, flows: list) -> None:
    logger.debug("Starting Worker")
    engine = create_async_engine(os.environ["PGBOUNCER_URL"], isolation_level="AUTOCOMMIT", echo=False, pool_size=100, max_overflow=5, pool_recycle=500)
    while True:
        await poll(id, limit, flows, engine)
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

async def check_and_cancel_long_running_task(background_task):
    logger.debug("Checking and cancelling long running tasks")
    current_time = time.time()
    task_timeout = 60 * 10 * 1000  # Timeout tasks after 10 minutes
    start_time = background_task[0]
    task = background_task[1]
    if current_time - start_time > task_timeout:
        task.cancel()
        logger.debug(f"Canceled tasks.")

async def poll(id: int, limit: Semaphore, flows: list, engine) -> None:
    background_tasks = set()
    while limit.locked():
        logger.debug("Workers are busy, will try again in a bit.")
        await asyncio.sleep(5)
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
        conn = await engine.connect()
        rs = await conn.execute(statement, values)
        for row in rs:
            logger.debug(str(id) + ": Found Work - " + str(row.id))
            task = asyncio.create_task(perform_task(row, limit, session_engine=engine))
            background_tasks.add((int(time.time()*1000), task))
        if randint(0, 1000) < 10:
            await check_and_cancel_long_running_tasks(background_tasks)

@sessionize
async def perform_task(row: any, limit: Semaphore, session=None, session_engine=None) -> None:
    logger.debug(f"TASK-{row.id}::Acquiring Locks")
    await limit.acquire()
    logger.debug(f"TASK-{row.id}::Starting Execution")
    try:
        logger.debug(f"TASK-{row.id}::Creating Task")
        await execute(row, session=session)
        logger.debug(f"TASK-{row.id}::Completed Task")
        statement = text("""
                        update scrape_request_log
                            set status='COMPLETE', scraped_at=now()
                            where id = %s
                    """ % row.id)
        conn = await session_engine.connect()
        await conn.execute(statement)
        await conn.close()
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
        conn = await session_engine.connect()
        await conn.execute(statement, parameters={'data': error, 'id': row.id})
        await conn.close()
        logger.debug(f"TASK-{row.id}::Finalized Failure")
    finally:
        logger.debug(f"TASK-{row.id}::Releasing Execution Locks")
        limit.release()


@sessionize
async def execute(scrape_log: ScrapeRequestLog, session=None) -> None:
    scrape_request = ScrapeRequestLogModel(scrape_log.id,
                                           scrape_log.flow,
                                           scrape_log.params)
    await perform_scrape_task(scrape_request, session=session)


def start_workers() -> None:
    workers = []
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
        limit = asyncio.Semaphore(int(os.environ["SCRAPE_LOG_MAX_CONCURRENCY_PER_LOOP"]))
        looper(1, limit, list(_whitelist.keys()))
    else:
        for flow in _whitelist:
            num_process = _whitelist[flow]['no_of_workers']
            for i in range(num_process):
                limit = asyncio.Semaphore(int(_whitelist[flow]['no_of_concurrency']))
                w = multiprocessing.Process(target=looper, args=(i + 1, limit, [flow]))
                w.start()
                workers.append(w)

if __name__ == "__main__":
    load_dotenv()
    logger.remove()
    logger.add(sys.stderr, level=os.environ["LOG_LEVEL"])
    start_workers()
