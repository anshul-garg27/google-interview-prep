from enum import auto, StrEnum


class Platform(StrEnum):
    INSTAGRAM = auto()
    YOUTUBE = auto()
    SHOPIFY = auto()


class ScrapeLogRequestStatus(StrEnum):
    PENDING = auto()
    RETRIEVING = auto()
    PARSING = auto()
    COMPLETE = auto()
    PROCESSING = auto()
    FAILED = auto()
