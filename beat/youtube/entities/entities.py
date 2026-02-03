from sqlalchemy import Column, TIMESTAMP, String, Boolean, BigInteger
from sqlalchemy.dialects.postgresql import JSONB

from core.entities.entities import Base


class YoutubeAccount(Base):
    __tablename__ = 'youtube_account'

    id = Column(BigInteger, primary_key=True)
    description = Column(String)
    updated_at = Column(TIMESTAMP)
    created_at = Column(TIMESTAMP)
    published_at = Column(TIMESTAMP)
    channel_id = Column(String)
    title = Column(String)
    custom_url = Column(String)
    subscribers = Column(BigInteger)
    uploads = Column(BigInteger)
    views = Column(BigInteger)
    country = Column(String)
    is_subscriber_count_hidden = Column(Boolean)
    is_verified = Column(Boolean)
    uploads_playlist_id = Column(String)
    thumbnails = Column(String)
    thumbnail = Column(String)
    external_urls = Column(String)

class YoutubeAccountTimeSeries(Base):
    __tablename__ = 'youtube_account_ts'

    id = Column(BigInteger, primary_key=True)
    description = Column(String)
    created_at = Column(TIMESTAMP)
    published_at = Column(TIMESTAMP)
    channel_id = Column(String)
    title = Column(String)
    custom_url = Column(String)
    subscribers = Column(BigInteger)
    uploads = Column(BigInteger)
    views = Column(BigInteger)
    country = Column(String)
    is_subscriber_count_hidden = Column(Boolean)
    uploads_playlist_id = Column(String)
    thumbnails = Column(String)
    thumbnail = Column(String)

class YoutubePost(Base):
    __tablename__ = 'youtube_post'

    id = Column(BigInteger, primary_key=True)
    shortcode = Column(String)
    description = Column(String)
    updated_at = Column(TIMESTAMP)
    created_at = Column(TIMESTAMP)
    published_at = Column(TIMESTAMP)
    channel_id = Column(String)
    title = Column(String)
    views = Column(BigInteger)
    likes = Column(BigInteger)
    comments = Column(BigInteger)
    favourites = Column(BigInteger)
    thumbnails = Column(String)
    thumbnail = Column(String)
    duration = Column(String)
    channel_title = Column(String)
    category_id = Column(String)
    category = Column(String)
    audio_language = Column(String)
    is_licensed = Column(Boolean)
    tags = Column(String)
    content_rating = Column(String)
    playlist_id = Column(String)
    privacy_status = Column(String)
    post_type = Column(String)
    keywords = Column(String)
    
class YoutubePostTimeSeries(Base):
    __tablename__ = 'youtube_post_ts'

    id = Column(BigInteger, primary_key=True)
    shortcode = Column(String)
    description = Column(String)
    created_at = Column(TIMESTAMP)
    published_at = Column(TIMESTAMP)
    channel_id = Column(String)
    title = Column(String)
    views = Column(BigInteger)
    likes = Column(BigInteger)
    favourites = Column(BigInteger)
    thumbnails = Column(String)
    thumbnail = Column(String)
    duration = Column(String)
    channel_title = Column(String)
    category_id = Column(String)
    audio_language = Column(String)
    is_licensed = Column(Boolean)
    tags = Column(String)
    content_rating = Column(String)
    playlist_id = Column(String)
    privacy_status = Column(String)
class YoutubeProfileInsights(Base):
    __tablename__ = 'youtube_profile_insights'

    id = Column(BigInteger, primary_key=True)
    channel_id = Column(String)
    created_at = Column(TIMESTAMP)
    updated_at = Column(TIMESTAMP)
    country = Column(JSONB)
    city = Column(JSONB)
    gender_age = Column(JSONB)
