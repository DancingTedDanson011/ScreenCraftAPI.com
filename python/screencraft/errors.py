"""
ScreenCraft SDK - Custom Exceptions

This module defines all custom exceptions used by the ScreenCraft SDK.
"""

from typing import Optional, Dict, Any


class ScreenCraftError(Exception):
    """Base exception for all ScreenCraft SDK errors."""

    def __init__(
        self,
        message: str,
        status_code: Optional[int] = None,
        response_body: Optional[Dict[str, Any]] = None,
    ) -> None:
        super().__init__(message)
        self.message = message
        self.status_code = status_code
        self.response_body = response_body

    def __str__(self) -> str:
        if self.status_code:
            return f"[{self.status_code}] {self.message}"
        return self.message

    def __repr__(self) -> str:
        return (
            f"{self.__class__.__name__}("
            f"message={self.message!r}, "
            f"status_code={self.status_code!r}, "
            f"response_body={self.response_body!r})"
        )


class AuthenticationError(ScreenCraftError):
    """Raised when API key is invalid or missing."""

    def __init__(
        self,
        message: str = "Invalid or missing API key",
        status_code: Optional[int] = 401,
        response_body: Optional[Dict[str, Any]] = None,
    ) -> None:
        super().__init__(message, status_code, response_body)


class RateLimitError(ScreenCraftError):
    """Raised when API rate limit is exceeded."""

    def __init__(
        self,
        message: str = "Rate limit exceeded",
        status_code: Optional[int] = 429,
        response_body: Optional[Dict[str, Any]] = None,
        retry_after: Optional[int] = None,
    ) -> None:
        super().__init__(message, status_code, response_body)
        self.retry_after = retry_after


class ValidationError(ScreenCraftError):
    """Raised when request parameters are invalid."""

    def __init__(
        self,
        message: str = "Invalid request parameters",
        status_code: Optional[int] = 400,
        response_body: Optional[Dict[str, Any]] = None,
        field: Optional[str] = None,
    ) -> None:
        super().__init__(message, status_code, response_body)
        self.field = field


class NotFoundError(ScreenCraftError):
    """Raised when the requested resource is not found."""

    def __init__(
        self,
        message: str = "Resource not found",
        status_code: Optional[int] = 404,
        response_body: Optional[Dict[str, Any]] = None,
    ) -> None:
        super().__init__(message, status_code, response_body)


class ServerError(ScreenCraftError):
    """Raised when the API server encounters an error."""

    def __init__(
        self,
        message: str = "Internal server error",
        status_code: Optional[int] = 500,
        response_body: Optional[Dict[str, Any]] = None,
    ) -> None:
        super().__init__(message, status_code, response_body)


class TimeoutError(ScreenCraftError):
    """Raised when a request times out."""

    def __init__(
        self,
        message: str = "Request timed out",
        status_code: Optional[int] = None,
        response_body: Optional[Dict[str, Any]] = None,
    ) -> None:
        super().__init__(message, status_code, response_body)


class ConnectionError(ScreenCraftError):
    """Raised when connection to the API fails."""

    def __init__(
        self,
        message: str = "Failed to connect to ScreenCraft API",
        status_code: Optional[int] = None,
        response_body: Optional[Dict[str, Any]] = None,
    ) -> None:
        super().__init__(message, status_code, response_body)


class WebhookError(ScreenCraftError):
    """Raised when webhook delivery or processing fails."""

    def __init__(
        self,
        message: str = "Webhook processing failed",
        status_code: Optional[int] = None,
        response_body: Optional[Dict[str, Any]] = None,
    ) -> None:
        super().__init__(message, status_code, response_body)


class RetryExhaustedError(ScreenCraftError):
    """Raised when all retry attempts have been exhausted."""

    def __init__(
        self,
        message: str = "All retry attempts exhausted",
        status_code: Optional[int] = None,
        response_body: Optional[Dict[str, Any]] = None,
        attempts: int = 0,
        last_exception: Optional[Exception] = None,
    ) -> None:
        super().__init__(message, status_code, response_body)
        self.attempts = attempts
        self.last_exception = last_exception
