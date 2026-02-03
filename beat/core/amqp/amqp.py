import asyncio
import os
import sys
import traceback

import aio_pika
from kombu import Connection
from kombu import Message
from loguru import logger

from core.amqp.models import AmqpListener
from core.helpers.session import SessionFactory
from utils.db import get_session_for_engine


async def async_listener_wrapper(fn, engine, body, message: Message):
    session = get_session_for_engine(engine)
    try:
        func_res = await fn(session, body, message)
        await session.commit()
        await session.close()
        return func_res
    except Exception as e:
        logger.error(e)
        message.reject()
        await session.rollback()
        await session.close()


def listener_wrapper(fn):
    engine = SessionFactory().get_engine()
    return lambda body, message: asyncio.run(async_listener_wrapper(fn, engine, body, message))


async def async_listener(config):
    engine = SessionFactory().get_engine()
    connection = await aio_pika.connect_robust(
        os.environ["RMQ_URL"],
    )
    queue_name = config.queue
    async with connection:
        channel = await connection.channel()
        await channel.set_qos(prefetch_count=config.prefetch_count)
        queue = await channel.declare_queue(queue_name, durable=True)
        async with queue.iterator() as queue_iter:
            async for message in queue_iter:
                async with message.process():
                    session = get_session_for_engine(engine)
                    try:
                        await config.fn(session, message.body, message)
                        await session.commit()
                        await session.close()
                    except Exception as e:
                        logger.error(traceback.format_exc())
                        logger.error(e)
                        await session.rollback()
                        await session.close()


def listener(config: AmqpListener) -> None:
    logger.remove()
    logger.add(sys.stderr, level=os.environ["LOG_LEVEL"])
    asyncio.run(async_listener(config))


def publish(payload: dict, exchange: str, routing_key: str) -> None:
    logger.debug("Publishing message", payload, exchange, routing_key)
    with Connection(os.environ["RMQ_URL"]) as conn:
        producer = conn.Producer(serializer='json')
        producer.publish(payload,
                         exchange=exchange,
                         routing_key=routing_key)