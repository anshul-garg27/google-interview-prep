from dataclasses import dataclass
from sqlite3 import Date
from typing import List

from core.models.models import Dimension, Metric


@dataclass
class YoutubeProfileLog:
    channel_id: str
    source: str
    dimensions: List[Dimension]
    metrics: List[Metric]


@dataclass
class YoutubePostLog:
    channel_id: str
    shortcode: str
    source: str
    dimensions: List[Dimension]
    metrics: List[Metric]


@dataclass
class YoutubePostActivityLog:
    activity_type: str  # comment, like
    shortcode: str
    actor_profile_id: str
    source: str
    dimensions: List[Dimension]
    metrics: List[Metric]


@dataclass
class YoutubeActivityLog:
    activity_id: str
    activity_type: str  # comment, like
    actor_channel_id: str
    source: str
    dimensions: List[Dimension]
    metrics: List[Metric]


@dataclass
class YoutubeProfileRelationshipLog:
    relationship_id: str
    relationship_type: str
    source: str
    source_channel_id: str
    target_channel_id: str
    source_dimensions: List[Dimension]
    source_metrics: List[Metric]
    target_dimensions: List[Dimension]
    target_metrics: List[Metric]
    subscribed_on: str
