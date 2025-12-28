"""
ScreenCraft SDK

A Python SDK for the ScreenCraft API - capture screenshots and generate PDFs from web pages.

Basic Usage:
    >>> from screencraft import ScreenCraft
    >>> client = ScreenCraft(api_key='your-api-key')
    >>> screenshot = client.screenshot(url='https://example.com', full_page=True)
    >>> with open('screenshot.png', 'wb') as f:
    ...     f.write(screenshot.data)

Async Usage:
    >>> import asyncio
    >>> from screencraft import AsyncScreenCraft
    >>> async def main():
    ...     async with AsyncScreenCraft(api_key='your-api-key') as client:
    ...         screenshot = await client.screenshot(url='https://example.com')
    ...         with open('screenshot.png', 'wb') as f:
    ...             f.write(screenshot.data)
    >>> asyncio.run(main())

For more information, visit: https://screencraftapi.com/docs
"""

__version__ = "1.0.0"
__author__ = "ScreenCraft"
__license__ = "MIT"

from .client import ScreenCraft, AsyncScreenCraft
from .errors import (
    ScreenCraftError,
    AuthenticationError,
    RateLimitError,
    ValidationError,
    NotFoundError,
    ServerError,
    TimeoutError,
    ConnectionError,
    WebhookError,
    RetryExhaustedError,
)
from .types import (
    Viewport,
    Clip,
    Cookie,
    Header,
    WebhookConfig,
    ScreenshotOptions,
    PdfOptions,
    PdfMargins,
    ScreenshotResponse,
    PdfResponse,
    WebhookPayload,
    AccountInfo,
    ImageFormat,
    PdfFormat,
    ScrollPosition,
    ImageFormatEnum,
    PdfFormatEnum,
    VIEWPORT_PRESETS,
)

__all__ = [
    # Version
    "__version__",
    # Clients
    "ScreenCraft",
    "AsyncScreenCraft",
    # Errors
    "ScreenCraftError",
    "AuthenticationError",
    "RateLimitError",
    "ValidationError",
    "NotFoundError",
    "ServerError",
    "TimeoutError",
    "ConnectionError",
    "WebhookError",
    "RetryExhaustedError",
    # Types
    "Viewport",
    "Clip",
    "Cookie",
    "Header",
    "WebhookConfig",
    "ScreenshotOptions",
    "PdfOptions",
    "PdfMargins",
    "ScreenshotResponse",
    "PdfResponse",
    "WebhookPayload",
    "AccountInfo",
    "ImageFormat",
    "PdfFormat",
    "ScrollPosition",
    "ImageFormatEnum",
    "PdfFormatEnum",
    "VIEWPORT_PRESETS",
]
