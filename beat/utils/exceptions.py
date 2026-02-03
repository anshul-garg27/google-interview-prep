class BeatException(Exception):

    def __init__(self, status_code: int, message: str):
        self.status_code = status_code
        self.message = message
        super().__init__(self.message)


class ProfileScrapingFailed(BeatException):

    def __init__(self, status_code: int, message: str):
        self.status_code = status_code
        self.message = message


class ProfileMissingError(BeatException):

    def __init__(self, status_code: int, message: str):
        self.status_code = status_code
        self.message = message


class ProfileInsightsMissingError(BeatException):
    def __init__(self, status_code: int, message: str):
        self.status_code = status_code
        self.message = message


class InsightScrapingFailed(BeatException):

    def __init__(self, status_code: int, message: str):
        self.status_code = status_code
        self.message = message


class OrderScrapingFailed(Exception):

    def __init__(self, status_code: int, message: str):
        self.status_code = status_code
        self.message = message


class NotSupportedError(Exception):

    def __init__(self, message: str):
        self.message = message
        super().__init__(self.message)


class NoAvailableSources(Exception):
    def __init__(self, message: str):
        self.message = message
        super().__init__(self.message)


class CredentialNotFound(BeatException):
    def __init__(self, status_code: int, message: str):
        self.status_code = status_code
        self.message = message


class TokenValidationFalied(BeatException):
    def __init__(self, status_code: int, message: str):
        self.status_code = status_code
        self.message = message


class DataAccessExpired(BeatException):
    def __init__(self, status_code: int, message: str):
        self.status_code = status_code
        self.message = message


class QuotaExceeded(BeatException):
    def __init__(self, status_code: int, message: str):
        self.status_code = status_code
        self.message = message


class PostScrapingFailed(Exception):
    def __init__(self, message: str):
        self.message = message
        super().__init__(self.message)


class StartDateAfterEndDateError(Exception):
    def __init__(self, status_code: int, message: str):
        self.status_code = status_code
        self.message = message
        super().__init__(self.message)


class UrlSignatureExpiredError(Exception):
    def __init__(self, message: str):
        self.message = message
        super().__init__(self.message)


class InvalidCsvFormatError(Exception):
    def __init__(self, message: str):
        self.message = message
        super().__init__(self.message)


class SubscriptionNotAllowedError(Exception):
    def __init__(self, status_code: int, message: str):
        self.status_code = status_code
        self.message = message
        super().__init__(self.message)

class NoSubscriptionsFoundError(Exception):
    def __init__(self, status_code: int, message: str):
        self.status_code = status_code
        self.message = message
        super().__init__(self.message)
