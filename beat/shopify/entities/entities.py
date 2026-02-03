from sqlalchemy import Column, TIMESTAMP, String, BigInteger
from sqlalchemy.dialects.postgresql import JSONB

from core.entities.entities import Base

class Order(Base):
    __tablename__ = 'order'

    id = Column(BigInteger, primary_key=True)
    platform = Column(String)
    platform_order_id = Column(String)
    status = Column(String)
    name = Column(String)
    store = Column(String)
    store_name = Column(String)
    landing_site = Column(String)
    utm_params = Column(String)
    order_data = Column(JSONB)
    created_at = Column(TIMESTAMP)
    updated_at = Column(TIMESTAMP)
