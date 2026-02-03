import os
import threading
import traceback
from functools import wraps

from loguru import logger
from sqlalchemy.ext.asyncio import AsyncSession, create_async_engine, AsyncEngine

from utils.db import get_session_for_engine

engine = None


class SingletonDoubleChecked(object):
    # resources shared by each and every
    # instance

    __singleton_lock = threading.Lock()
    # __singleton_lock = asyncio.Lock()
    __singleton_instance = None

    # define the classmethod
    @classmethod
    def instance(cls):

        # check for the singleton instance
        if not cls.__singleton_instance:
            with cls.__singleton_lock:
                if not cls.__singleton_instance:
                    cls.__singleton_instance = cls()

        # return the singleton instance
        return cls.__singleton_instance


class SessionFactory(SingletonDoubleChecked):
    def __init__(self):
        logger.debug("Creating PG Engine")
        self.engine: AsyncEngine = create_async_engine(os.environ["PG_URL"], echo=False)

    def get_engine(self) -> AsyncEngine:
        return self.engine


def sessionize(func):
    @wraps(func)
    async def wrapper(*args, **kwargs):
        if not kwargs:
            kwargs = {}
        if 'session' not in kwargs or not kwargs['session']:
            if 'session_engine' in kwargs and kwargs['session_engine']:
                session_engine = kwargs['session_engine']
            else:
                logger.error(f"MISSING SESSION ENGINE - {func.__name__} with args: {args}, kwargs: {kwargs}")
                session_engine = SessionFactory().get_engine()
            session: AsyncSession = get_session_for_engine(session_engine)
            kwargs['session'] = session
            try:
                func_res = await func(*args, **kwargs)
                await session.commit()
                await session.close()
                return func_res
            except Exception as e:
                logger.error(f"Error while handling API - {e}")
                logger.error(traceback.format_exc())
                await session.rollback()
                await session.close()
                raise e
        else:
            return await func(*args, **kwargs)  # do nothing, we're already getting session from upstream

    return wrapper


def sessionize_api(func):
    @wraps(func)
    async def wrapper(*args, **kwargs):
        if 'session' not in kwargs or not kwargs['session']:
            if 'request' not in kwargs or not kwargs['request']:
                raise Exception("Invalid Request - Request Missing")
            request = kwargs['request']
            session_engine = request.app.state.db
            if not session_engine:
                raise Exception("Invalid Request - Session Missing")
            session: AsyncSession = get_session_for_engine(session_engine)
            kwargs['session'] = session
            try:
                func_res = await func(*args, **kwargs)
                await session.commit()
                await session.close()
                return func_res
            except Exception as e:
                logger.error(f"Error while handling API - {e}")
                logger.error(traceback.format_exc())
                await session.rollback()
                await session.close()
                raise e
        else:
            return await func(*args, **kwargs)  # do nothing, we're already getting session from upstream

    return wrapper


def sessionize_scl_api(func):
    @wraps(func)
    async def wrapper(*args, **kwargs):
        if 'session' not in kwargs or not kwargs['session']:
            if 'request' not in kwargs or not kwargs['request']:
                raise Exception("Invalid Request - Request Missing")
            request = kwargs['request']
            session_engine = request.app.state.scl_db
            if not session_engine:
                raise Exception("Invalid Request - Session Missing")
            session: AsyncSession = get_session_for_engine(session_engine)
            kwargs['session'] = session
            try:
                func_res = await func(*args, **kwargs)
                await session.commit()
                await session.close()
                return func_res
            except Exception as e:
                logger.error(f"Error while handling API - {e}")
                logger.error(traceback.format_exc())
                await session.rollback()
                await session.close()
                raise e
        else:
            return await func(*args, **kwargs)  # do nothing, we're already getting session from upstream

    return wrapper
