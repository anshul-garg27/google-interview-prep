from dataclasses import dataclass
from typing import List

from core.models.models import Dimension, Metric


@dataclass
class OrderLog:
    platform_order_id: str
    source: str
    dimensions: List[Dimension]
    metrics: List[Metric]
