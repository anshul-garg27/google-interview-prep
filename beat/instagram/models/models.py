from dataclasses import dataclass
from typing import List

from core.models.models import Dimension, Metric


@dataclass
class InstagramProfileLog:
    handle: str
    profile_id: str
    source: str
    dimensions: List[Dimension]
    metrics: List[Metric]


@dataclass
class InstagramPostLog:
    shortcode: str
    source: str
    dimensions: List[Dimension]
    metrics: List[Metric]


@dataclass
class InstagramPostActivityLog:
    activity_type: str  # comment, like
    shortcode: str
    actor_profile_id: str
    source: str
    dimensions: List[Dimension]
    metrics: List[Metric]


@dataclass
class InstagramRelationshipLog:
    relationship_type: str
    source_profile_id: str
    target_profile_id: str
    source: str
    source_dimensions: dict
    source_metrics: dict
    target_dimensions: dict
    target_metrics: dict