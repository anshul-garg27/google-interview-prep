from dataclasses import dataclass
from typing import Callable


@dataclass
class AmqpListener:
    exchange: str
    exchange_type: str
    routingKey: str
    queue: str
    workers: int
    prefetch_count: int
    fn: Callable
