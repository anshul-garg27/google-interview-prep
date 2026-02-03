import os
from datetime import datetime
from typing import TypeVar
from loguru import logger
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession, create_async_engine, AsyncEngine
from sqlalchemy.orm import sessionmaker

T = TypeVar("T")


async def get_or_create(session: AsyncSession, model: T, defaults: dict = None, **kwargs: any) -> T:
    instance = await get(session, model, defaults, **kwargs)
    if not instance:
        return await create(session, model, defaults, **kwargs)
    else:
        return instance, True


async def get(session: AsyncSession, model: T, defaults: dict = None, **kwargs: any) -> T:
    result = await session.execute(select(model).filter_by(**kwargs))
    instance = None
    if result:
        instance = result.first()
    if instance:
        return instance[0]
    else:
        return None


async def create(session: AsyncSession, model: T, defaults: dict = None, **kwargs: any) -> T:
    kwargs |= defaults or {}
    kwargs['created_at'] = datetime.now()
    instance = model(**kwargs)
    session.add(instance)
    return instance, True


def get_session_for_engine(engine: AsyncEngine) -> AsyncSession:
    session: AsyncSession = sessionmaker(bind=engine, class_=AsyncSession)()
    return session


def get_engine() -> AsyncEngine:
    engine = create_async_engine(os.environ["PG_URL"], echo=bool(int(os.environ["DEBUG_SQL"])), pool_size=1)
    return engine
