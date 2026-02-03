from sqlalchemy import Column, Integer, String, TIMESTAMP, Boolean, BigInteger, Float
from sqlalchemy.dialects.postgresql import JSONB
from sqlalchemy.orm import declarative_base

Base = declarative_base()


class Credential(Base):
    __tablename__ = 'credential'

    id = Column(BigInteger, primary_key=True)
    idempotency_key = Column(String, unique=True)
    source = Column(String)
    handle = Column(String)
    credentials = Column(JSONB)
    created_at = Column(TIMESTAMP)
    updated_at = Column(TIMESTAMP)
    disabled_till = Column(TIMESTAMP)
    enabled = Column(Boolean)
    data_access_expired = Column(Boolean)


class ProfileLog(Base):
    __tablename__ = 'profile_log'

    id = Column(BigInteger, primary_key=True)
    platform = Column(String)
    profile_id = Column(String)
    dimensions = Column(JSONB)
    metrics = Column(JSONB)
    source = Column(String)
    timestamp = Column(TIMESTAMP)


class PostLog(Base):
    __tablename__ = 'post_log'

    id = Column(BigInteger, primary_key=True)
    platform = Column(String)
    profile_id = Column(String)
    platform_post_id = Column(String)
    dimensions = Column(JSONB)
    metrics = Column(JSONB)
    source = Column(String)
    timestamp = Column(TIMESTAMP)


class SentimentLog(Base):
    __tablename__ = 'sentiment_log'

    id = Column(BigInteger, primary_key=True)
    platform = Column(String)
    platform_post_id = Column(String)
    comment_id = Column(String)
    comment = Column(JSONB)
    sentiment = Column(JSONB)
    score = Column(Float)
    dimensions = Column(JSONB)
    metrics = Column(JSONB)
    source = Column(String)
    timestamp = Column(TIMESTAMP)


class ProfileRelationshipLog(Base):
    __tablename__ = 'profile_relationship_log'

    id = Column(BigInteger, primary_key=True)
    relationship_type = Column(String)
    platform = Column(String)
    source_profile_id = Column(String)
    target_profile_id = Column(String)
    source_dimensions = Column(JSONB)
    source_metrics = Column(JSONB)
    target_dimensions = Column(JSONB)
    target_metrics = Column(JSONB)
    source = Column(String)
    timestamp = Column(TIMESTAMP)


class PostActivityLog(Base):
    __tablename__ = 'post_activity_log'

    id = Column(BigInteger, primary_key=True)
    activity_type = Column(String)
    platform = Column(String)
    platform_post_id = Column(String)
    actor_profile_id = Column(String)
    dimensions = Column(JSONB)
    metrics = Column(JSONB)
    source = Column(String)
    timestamp = Column(TIMESTAMP)


class ScrapeRequestLog(Base):
    __tablename__ = 'scrape_request_log'

    id = Column(BigInteger, primary_key=True)
    idempotency_key = Column(String)
    platform = Column(String)
    source = Column(String)
    account_id = Column(String)
    flow = Column(String)
    status = Column(String)
    priority = Column(Integer)
    params = Column(JSONB)
    data = Column(String)
    created_at = Column(TIMESTAMP)
    scraped_at = Column(TIMESTAMP)
    expires_at = Column(TIMESTAMP)
    picked_at = Column(TIMESTAMP)
    retry_count = Column(Integer)


class TaskLog(Base):
    __tablename__ = 'task_log'
    id = Column(BigInteger, primary_key=True)
    name = Column(String)
    arguments = Column(JSONB)
    source = Column(String)
    status = Column(String)
    created_at = Column(TIMESTAMP)


class AssetLog(Base):
    __tablename__ = 'asset_log'
    id = Column(BigInteger, primary_key=True)
    entity_id = Column(String)
    entity_type = Column(String)
    platform = Column(String)
    asset_type = Column(String)
    asset_url = Column(String)
    original_url = Column(String)
    updated_at = Column(TIMESTAMP)
    created_at = Column(TIMESTAMP)


class OrderLogs(Base):
    __tablename__ = 'order_log'

    id = Column(BigInteger, primary_key=True)
    platform = Column(String)
    platform_order_id = Column(String)
    store = Column(String)
    dimensions = Column(JSONB)
    metrics = Column(JSONB)
    source = Column(String)
    timestamp = Column(TIMESTAMP)


class ProfileActivityLog(Base):
    __tablename__ = 'profile_activity_log'

    id = Column(BigInteger, primary_key=True)
    activity_id = Column(String)
    activity_type = Column(String)  # comment, like
    actor_channel_id = Column(String)
    platform = Column(String)
    source = Column(String)
    dimensions = Column(JSONB)
    metrics = Column(JSONB)


class YtProfileRelationshipLog(Base):
    __tablename__ = 'profile_subscription_log'

    id = Column(BigInteger, primary_key=True)
    relationship_id = Column(String)
    relationship_type = Column(String)  # comment, like
    source_channel_id = Column(String)
    target_channel_id = Column(String)
    platform = Column(String)
    source = Column(String)
    source_dimensions = Column(JSONB)
    source_metrics = Column(JSONB)
    target_dimensions = Column(JSONB)
    target_metrics = Column(JSONB)
    subscribed_on = Column(String)
