from dataclasses import dataclass
from datetime import datetime
from pydantic import BaseModel
from credentials.models import CredentialModel
from typing import Optional

@dataclass
class Dimension:
    key: str
    value: str


@dataclass
class Metric:
    key: str
    value: float

class ScrapeRequest(BaseModel):
    flow: str
    platform: str
    params: dict
    reason: Optional[str]
    status: Optional[str]
    retry_count: Optional[int]
    event_timestamp: Optional[datetime]
    expires_at: Optional[datetime]
    picked_at: Optional[datetime]
    priority: Optional[int]


@dataclass
class ScrapeRequestLog:
    id: int
    flow: str
    params: dict


@dataclass
class Context:
    now: datetime


@dataclass
class RequestContext:
    credential: CredentialModel
    
@dataclass
class PaginationContext:
    cursor: str
    source: str