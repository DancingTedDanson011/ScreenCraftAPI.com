"""
ScreenCraft SDK - Client

This module provides the main ScreenCraft client for interacting with the API.
"""

import time
import random
import logging
from typing import Optional, Dict, Any, List, Union, Literal
from urllib.parse import urljoin

import requests

from .errors import (
    ScreenCraftError,
    AuthenticationError,
    RateLimitError,
    ValidationError,
    NotFoundError,
    ServerError,
    TimeoutError,
    ConnectionError,
    RetryExhaustedError,
)
from .types import (
    Viewport,
    Clip,
    Cookie,
    WebhookConfig,
    ScreenshotOptions,
    PdfOptions,
    PdfMargins,
    ScreenshotResponse,
    PdfResponse,
    AccountInfo,
    ImageFormat,
    PdfFormat,
    ScrollPosition,
    VIEWPORT_PRESETS,
)


logger = logging.getLogger("screencraft")


class ScreenCraft:
    """
    ScreenCraft API client for capturing screenshots and generating PDFs.

    Usage:
        >>> client = ScreenCraft(api_key='your-api-key')
        >>> screenshot = client.screenshot(url='https://example.com')
        >>> with open('screenshot.png', 'wb') as f:
        ...     f.write(screenshot.data)

    Args:
        api_key: Your ScreenCraft API key.
        base_url: Base URL for the API (default: https://screencraftapi.com/api/v1).
        timeout: Request timeout in seconds (default: 60).
        max_retries: Maximum number of retry attempts (default: 3).
        retry_delay: Initial delay between retries in seconds (default: 1).
        retry_max_delay: Maximum delay between retries in seconds (default: 30).
        retry_backoff: Backoff multiplier for retries (default: 2).
    """

    DEFAULT_BASE_URL = "https://screencraftapi.com/api/v1"
    DEFAULT_TIMEOUT = 60
    DEFAULT_MAX_RETRIES = 3
    DEFAULT_RETRY_DELAY = 1.0
    DEFAULT_RETRY_MAX_DELAY = 30.0
    DEFAULT_RETRY_BACKOFF = 2.0

    def __init__(
        self,
        api_key: str,
        base_url: str = DEFAULT_BASE_URL,
        timeout: int = DEFAULT_TIMEOUT,
        max_retries: int = DEFAULT_MAX_RETRIES,
        retry_delay: float = DEFAULT_RETRY_DELAY,
        retry_max_delay: float = DEFAULT_RETRY_MAX_DELAY,
        retry_backoff: float = DEFAULT_RETRY_BACKOFF,
    ) -> None:
        if not api_key:
            raise ValueError("API key is required")

        self.api_key = api_key
        self.base_url = base_url.rstrip("/")
        self.timeout = timeout
        self.max_retries = max_retries
        self.retry_delay = retry_delay
        self.retry_max_delay = retry_max_delay
        self.retry_backoff = retry_backoff

        self._session = requests.Session()
        self._session.headers.update({
            "Authorization": f"Bearer {self.api_key}",
            "Content-Type": "application/json",
            "User-Agent": "ScreenCraft-Python-SDK/1.0.0",
        })

    def _get_url(self, endpoint: str) -> str:
        """Build full URL for an endpoint."""
        return urljoin(self.base_url + "/", endpoint.lstrip("/"))

    def _handle_error_response(
        self,
        response: requests.Response,
    ) -> None:
        """Handle error responses from the API."""
        status_code = response.status_code

        try:
            body = response.json()
            message = body.get("error", body.get("message", response.text))
        except (ValueError, KeyError):
            body = {}
            message = response.text or f"HTTP {status_code}"

        if status_code == 401:
            raise AuthenticationError(message, status_code, body)
        elif status_code == 429:
            retry_after = response.headers.get("Retry-After")
            raise RateLimitError(
                message,
                status_code,
                body,
                retry_after=int(retry_after) if retry_after else None,
            )
        elif status_code == 400:
            field = body.get("field") if body else None
            raise ValidationError(message, status_code, body, field=field)
        elif status_code == 404:
            raise NotFoundError(message, status_code, body)
        elif status_code >= 500:
            raise ServerError(message, status_code, body)
        else:
            raise ScreenCraftError(message, status_code, body)

    def _should_retry(self, exception: Exception, attempt: int) -> bool:
        """Determine if a request should be retried."""
        if attempt >= self.max_retries:
            return False

        # Retry on connection errors and server errors
        if isinstance(exception, (ConnectionError, TimeoutError)):
            return True
        if isinstance(exception, ServerError):
            return True
        if isinstance(exception, RateLimitError):
            return True

        return False

    def _calculate_delay(self, attempt: int, rate_limit_error: Optional[RateLimitError] = None) -> float:
        """Calculate delay before next retry with exponential backoff."""
        if rate_limit_error and rate_limit_error.retry_after:
            return float(rate_limit_error.retry_after)

        delay = self.retry_delay * (self.retry_backoff ** attempt)
        # Add jitter to prevent thundering herd
        jitter = random.uniform(0, 0.1 * delay)
        delay = min(delay + jitter, self.retry_max_delay)
        return delay

    def _request(
        self,
        method: str,
        endpoint: str,
        json_data: Optional[Dict[str, Any]] = None,
        params: Optional[Dict[str, Any]] = None,
        stream: bool = False,
    ) -> requests.Response:
        """Make an HTTP request with automatic retries."""
        url = self._get_url(endpoint)
        last_exception: Optional[Exception] = None

        for attempt in range(self.max_retries + 1):
            try:
                response = self._session.request(
                    method=method,
                    url=url,
                    json=json_data,
                    params=params,
                    timeout=self.timeout,
                    stream=stream,
                )

                if response.status_code >= 400:
                    self._handle_error_response(response)

                return response

            except requests.exceptions.Timeout as e:
                last_exception = TimeoutError(str(e))
            except requests.exceptions.ConnectionError as e:
                last_exception = ConnectionError(str(e))
            except (RateLimitError, ServerError) as e:
                last_exception = e

            if not self._should_retry(last_exception, attempt):
                break

            delay = self._calculate_delay(
                attempt,
                last_exception if isinstance(last_exception, RateLimitError) else None,
            )
            logger.debug(f"Retrying in {delay:.2f}s (attempt {attempt + 1}/{self.max_retries})")
            time.sleep(delay)

        if last_exception:
            if isinstance(last_exception, ScreenCraftError):
                raise last_exception
            raise RetryExhaustedError(
                f"All {self.max_retries} retry attempts exhausted",
                attempts=self.max_retries,
                last_exception=last_exception,
            )

        # Should never reach here, but just in case
        raise ScreenCraftError("Unknown error occurred")

    def screenshot(
        self,
        url: str,
        *,
        format: ImageFormat = "png",
        quality: int = 80,
        full_page: bool = False,
        viewport: Optional[Union[Viewport, str]] = None,
        clip: Optional[Clip] = None,
        scroll_position: Optional[ScrollPosition] = None,
        accept_cookies: bool = False,
        delay: int = 0,
        wait_until: Literal["load", "domcontentloaded", "networkidle0", "networkidle2"] = "load",
        timeout: int = 30000,
        cookies: Optional[List[Cookie]] = None,
        headers: Optional[Dict[str, str]] = None,
        user_agent: Optional[str] = None,
        bypass_csp: bool = False,
        javascript_enabled: bool = True,
        dark_mode: bool = False,
        block_ads: bool = False,
        hide_selectors: Optional[List[str]] = None,
        click_selector: Optional[str] = None,
        wait_for_selector: Optional[str] = None,
        webhook: Optional[WebhookConfig] = None,
    ) -> ScreenshotResponse:
        """
        Capture a screenshot of a web page.

        Args:
            url: The URL to capture.
            format: Image format ('png', 'jpeg', 'webp').
            quality: Image quality (1-100, only for jpeg/webp).
            full_page: Capture the full scrollable page.
            viewport: Viewport configuration or preset name.
            clip: Clip region for partial screenshots.
            scroll_position: Scroll position ('top' or 'bottom').
            accept_cookies: Auto-accept cookie banners.
            delay: Delay in milliseconds before capture.
            wait_until: Page load event to wait for.
            timeout: Navigation timeout in milliseconds.
            cookies: Cookies to set before capture.
            headers: Custom HTTP headers.
            user_agent: Custom user agent string.
            bypass_csp: Bypass Content Security Policy.
            javascript_enabled: Enable JavaScript execution.
            dark_mode: Enable dark mode emulation.
            block_ads: Block advertisements.
            hide_selectors: CSS selectors to hide.
            click_selector: CSS selector to click before capture.
            wait_for_selector: CSS selector to wait for.
            webhook: Webhook configuration for async delivery.

        Returns:
            ScreenshotResponse with the captured image data.

        Raises:
            ValidationError: If parameters are invalid.
            AuthenticationError: If API key is invalid.
            RateLimitError: If rate limit is exceeded.
            ScreenCraftError: For other API errors.
        """
        # Handle viewport presets
        if isinstance(viewport, str):
            if viewport not in VIEWPORT_PRESETS:
                raise ValidationError(
                    f"Unknown viewport preset: {viewport}. "
                    f"Available presets: {', '.join(VIEWPORT_PRESETS.keys())}",
                    field="viewport",
                )
            viewport = VIEWPORT_PRESETS[viewport]

        options = ScreenshotOptions(
            url=url,
            format=format,
            quality=quality,
            full_page=full_page,
            viewport=viewport,
            clip=clip,
            scroll_position=scroll_position,
            accept_cookies=accept_cookies,
            delay=delay,
            wait_until=wait_until,
            timeout=timeout,
            cookies=cookies,
            headers=headers,
            user_agent=user_agent,
            bypass_csp=bypass_csp,
            javascript_enabled=javascript_enabled,
            dark_mode=dark_mode,
            block_ads=block_ads,
            hide_selectors=hide_selectors,
            click_selector=click_selector,
            wait_for_selector=wait_for_selector,
            webhook=webhook,
        )

        response = self._request(
            "POST",
            "/screenshots",
            json_data=options.to_dict(),
            stream=True,
        )

        # Extract metadata from headers
        request_id = response.headers.get("X-Request-Id")
        credits_used = response.headers.get("X-Credits-Used")
        credits_remaining = response.headers.get("X-Credits-Remaining")

        return ScreenshotResponse(
            success=True,
            data=response.content,
            url=url,
            content_type=response.headers.get("Content-Type"),
            request_id=request_id,
            credits_used=int(credits_used) if credits_used else None,
            credits_remaining=int(credits_remaining) if credits_remaining else None,
        )

    def pdf(
        self,
        url: str,
        *,
        format: PdfFormat = "A4",
        landscape: bool = False,
        print_background: bool = True,
        margins: Optional[PdfMargins] = None,
        scale: float = 1.0,
        page_ranges: Optional[str] = None,
        header_template: Optional[str] = None,
        footer_template: Optional[str] = None,
        display_header_footer: bool = False,
        prefer_css_page_size: bool = False,
        accept_cookies: bool = False,
        delay: int = 0,
        wait_until: Literal["load", "domcontentloaded", "networkidle0", "networkidle2"] = "load",
        timeout: int = 30000,
        cookies: Optional[List[Cookie]] = None,
        headers: Optional[Dict[str, str]] = None,
        user_agent: Optional[str] = None,
        javascript_enabled: bool = True,
        wait_for_selector: Optional[str] = None,
        webhook: Optional[WebhookConfig] = None,
    ) -> PdfResponse:
        """
        Generate a PDF from a web page.

        Args:
            url: The URL to convert to PDF.
            format: Page format ('A0'-'A6', 'Letter', 'Legal', 'Tabloid').
            landscape: Use landscape orientation.
            print_background: Print background graphics.
            margins: Page margins.
            scale: Scale of the webpage rendering (0.1-2).
            page_ranges: Paper ranges to print (e.g., '1-5, 8, 11-13').
            header_template: HTML template for the print header.
            footer_template: HTML template for the print footer.
            display_header_footer: Display header and footer.
            prefer_css_page_size: Use CSS page size from the content.
            accept_cookies: Auto-accept cookie banners.
            delay: Delay in milliseconds before capture.
            wait_until: Page load event to wait for.
            timeout: Navigation timeout in milliseconds.
            cookies: Cookies to set before capture.
            headers: Custom HTTP headers.
            user_agent: Custom user agent string.
            javascript_enabled: Enable JavaScript execution.
            wait_for_selector: CSS selector to wait for.
            webhook: Webhook configuration for async delivery.

        Returns:
            PdfResponse with the generated PDF data.

        Raises:
            ValidationError: If parameters are invalid.
            AuthenticationError: If API key is invalid.
            RateLimitError: If rate limit is exceeded.
            ScreenCraftError: For other API errors.
        """
        options = PdfOptions(
            url=url,
            format=format,
            landscape=landscape,
            print_background=print_background,
            margins=margins,
            scale=scale,
            page_ranges=page_ranges,
            header_template=header_template,
            footer_template=footer_template,
            display_header_footer=display_header_footer,
            prefer_css_page_size=prefer_css_page_size,
            accept_cookies=accept_cookies,
            delay=delay,
            wait_until=wait_until,
            timeout=timeout,
            cookies=cookies,
            headers=headers,
            user_agent=user_agent,
            javascript_enabled=javascript_enabled,
            wait_for_selector=wait_for_selector,
            webhook=webhook,
        )

        response = self._request(
            "POST",
            "/pdfs",
            json_data=options.to_dict(),
            stream=True,
        )

        # Extract metadata from headers
        request_id = response.headers.get("X-Request-Id")
        page_count = response.headers.get("X-Page-Count")
        credits_used = response.headers.get("X-Credits-Used")
        credits_remaining = response.headers.get("X-Credits-Remaining")

        return PdfResponse(
            success=True,
            data=response.content,
            url=url,
            content_type=response.headers.get("Content-Type"),
            request_id=request_id,
            page_count=int(page_count) if page_count else None,
            credits_used=int(credits_used) if credits_used else None,
            credits_remaining=int(credits_remaining) if credits_remaining else None,
        )

    def get_account_info(self) -> AccountInfo:
        """
        Get current account information and usage.

        Returns:
            AccountInfo with account details and credits.

        Raises:
            AuthenticationError: If API key is invalid.
            ScreenCraftError: For other API errors.
        """
        response = self._request("GET", "/account")
        return AccountInfo.from_dict(response.json())

    def close(self) -> None:
        """Close the HTTP session."""
        self._session.close()

    def __enter__(self) -> "ScreenCraft":
        return self

    def __exit__(self, exc_type: Any, exc_val: Any, exc_tb: Any) -> None:
        self.close()


# Async client using aiohttp
class AsyncScreenCraft:
    """
    Async ScreenCraft API client for capturing screenshots and generating PDFs.

    Usage:
        >>> import asyncio
        >>> async def main():
        ...     async with AsyncScreenCraft(api_key='your-api-key') as client:
        ...         screenshot = await client.screenshot(url='https://example.com')
        ...         with open('screenshot.png', 'wb') as f:
        ...             f.write(screenshot.data)
        >>> asyncio.run(main())

    Args:
        api_key: Your ScreenCraft API key.
        base_url: Base URL for the API (default: https://screencraftapi.com/api/v1).
        timeout: Request timeout in seconds (default: 60).
        max_retries: Maximum number of retry attempts (default: 3).
        retry_delay: Initial delay between retries in seconds (default: 1).
        retry_max_delay: Maximum delay between retries in seconds (default: 30).
        retry_backoff: Backoff multiplier for retries (default: 2).
    """

    DEFAULT_BASE_URL = "https://screencraftapi.com/api/v1"
    DEFAULT_TIMEOUT = 60
    DEFAULT_MAX_RETRIES = 3
    DEFAULT_RETRY_DELAY = 1.0
    DEFAULT_RETRY_MAX_DELAY = 30.0
    DEFAULT_RETRY_BACKOFF = 2.0

    def __init__(
        self,
        api_key: str,
        base_url: str = DEFAULT_BASE_URL,
        timeout: int = DEFAULT_TIMEOUT,
        max_retries: int = DEFAULT_MAX_RETRIES,
        retry_delay: float = DEFAULT_RETRY_DELAY,
        retry_max_delay: float = DEFAULT_RETRY_MAX_DELAY,
        retry_backoff: float = DEFAULT_RETRY_BACKOFF,
    ) -> None:
        if not api_key:
            raise ValueError("API key is required")

        self.api_key = api_key
        self.base_url = base_url.rstrip("/")
        self.timeout = timeout
        self.max_retries = max_retries
        self.retry_delay = retry_delay
        self.retry_max_delay = retry_max_delay
        self.retry_backoff = retry_backoff

        self._session: Optional["aiohttp.ClientSession"] = None
        self._headers = {
            "Authorization": f"Bearer {self.api_key}",
            "Content-Type": "application/json",
            "User-Agent": "ScreenCraft-Python-SDK/1.0.0",
        }

    async def _ensure_session(self) -> "aiohttp.ClientSession":
        """Ensure aiohttp session is created."""
        if self._session is None or self._session.closed:
            import aiohttp
            timeout = aiohttp.ClientTimeout(total=self.timeout)
            self._session = aiohttp.ClientSession(
                headers=self._headers,
                timeout=timeout,
            )
        return self._session

    def _get_url(self, endpoint: str) -> str:
        """Build full URL for an endpoint."""
        return urljoin(self.base_url + "/", endpoint.lstrip("/"))

    async def _handle_error_response(
        self,
        response: "aiohttp.ClientResponse",
    ) -> None:
        """Handle error responses from the API."""
        import aiohttp

        status_code = response.status

        try:
            body = await response.json()
            message = body.get("error", body.get("message", await response.text()))
        except (ValueError, KeyError, aiohttp.ContentTypeError):
            body = {}
            text = await response.text()
            message = text or f"HTTP {status_code}"

        if status_code == 401:
            raise AuthenticationError(message, status_code, body)
        elif status_code == 429:
            retry_after = response.headers.get("Retry-After")
            raise RateLimitError(
                message,
                status_code,
                body,
                retry_after=int(retry_after) if retry_after else None,
            )
        elif status_code == 400:
            field = body.get("field") if body else None
            raise ValidationError(message, status_code, body, field=field)
        elif status_code == 404:
            raise NotFoundError(message, status_code, body)
        elif status_code >= 500:
            raise ServerError(message, status_code, body)
        else:
            raise ScreenCraftError(message, status_code, body)

    def _should_retry(self, exception: Exception, attempt: int) -> bool:
        """Determine if a request should be retried."""
        if attempt >= self.max_retries:
            return False

        if isinstance(exception, (ConnectionError, TimeoutError)):
            return True
        if isinstance(exception, ServerError):
            return True
        if isinstance(exception, RateLimitError):
            return True

        return False

    def _calculate_delay(self, attempt: int, rate_limit_error: Optional[RateLimitError] = None) -> float:
        """Calculate delay before next retry with exponential backoff."""
        if rate_limit_error and rate_limit_error.retry_after:
            return float(rate_limit_error.retry_after)

        delay = self.retry_delay * (self.retry_backoff ** attempt)
        jitter = random.uniform(0, 0.1 * delay)
        delay = min(delay + jitter, self.retry_max_delay)
        return delay

    async def _request(
        self,
        method: str,
        endpoint: str,
        json_data: Optional[Dict[str, Any]] = None,
        params: Optional[Dict[str, Any]] = None,
    ) -> tuple[bytes, Dict[str, str]]:
        """Make an async HTTP request with automatic retries."""
        import aiohttp
        import asyncio as async_module

        url = self._get_url(endpoint)
        session = await self._ensure_session()
        last_exception: Optional[Exception] = None

        for attempt in range(self.max_retries + 1):
            try:
                async with session.request(
                    method=method,
                    url=url,
                    json=json_data,
                    params=params,
                ) as response:
                    if response.status >= 400:
                        await self._handle_error_response(response)

                    data = await response.read()
                    headers = dict(response.headers)
                    return data, headers

            except async_module.TimeoutError as e:
                last_exception = TimeoutError(str(e))
            except aiohttp.ClientConnectionError as e:
                last_exception = ConnectionError(str(e))
            except (RateLimitError, ServerError) as e:
                last_exception = e

            if not self._should_retry(last_exception, attempt):
                break

            delay = self._calculate_delay(
                attempt,
                last_exception if isinstance(last_exception, RateLimitError) else None,
            )
            logger.debug(f"Retrying in {delay:.2f}s (attempt {attempt + 1}/{self.max_retries})")
            await async_module.sleep(delay)

        if last_exception:
            if isinstance(last_exception, ScreenCraftError):
                raise last_exception
            raise RetryExhaustedError(
                f"All {self.max_retries} retry attempts exhausted",
                attempts=self.max_retries,
                last_exception=last_exception,
            )

        raise ScreenCraftError("Unknown error occurred")

    async def screenshot(
        self,
        url: str,
        *,
        format: ImageFormat = "png",
        quality: int = 80,
        full_page: bool = False,
        viewport: Optional[Union[Viewport, str]] = None,
        clip: Optional[Clip] = None,
        scroll_position: Optional[ScrollPosition] = None,
        accept_cookies: bool = False,
        delay: int = 0,
        wait_until: Literal["load", "domcontentloaded", "networkidle0", "networkidle2"] = "load",
        timeout: int = 30000,
        cookies: Optional[List[Cookie]] = None,
        headers: Optional[Dict[str, str]] = None,
        user_agent: Optional[str] = None,
        bypass_csp: bool = False,
        javascript_enabled: bool = True,
        dark_mode: bool = False,
        block_ads: bool = False,
        hide_selectors: Optional[List[str]] = None,
        click_selector: Optional[str] = None,
        wait_for_selector: Optional[str] = None,
        webhook: Optional[WebhookConfig] = None,
    ) -> ScreenshotResponse:
        """
        Capture a screenshot of a web page asynchronously.

        See ScreenCraft.screenshot for full documentation.
        """
        if isinstance(viewport, str):
            if viewport not in VIEWPORT_PRESETS:
                raise ValidationError(
                    f"Unknown viewport preset: {viewport}. "
                    f"Available presets: {', '.join(VIEWPORT_PRESETS.keys())}",
                    field="viewport",
                )
            viewport = VIEWPORT_PRESETS[viewport]

        options = ScreenshotOptions(
            url=url,
            format=format,
            quality=quality,
            full_page=full_page,
            viewport=viewport,
            clip=clip,
            scroll_position=scroll_position,
            accept_cookies=accept_cookies,
            delay=delay,
            wait_until=wait_until,
            timeout=timeout,
            cookies=cookies,
            headers=headers,
            user_agent=user_agent,
            bypass_csp=bypass_csp,
            javascript_enabled=javascript_enabled,
            dark_mode=dark_mode,
            block_ads=block_ads,
            hide_selectors=hide_selectors,
            click_selector=click_selector,
            wait_for_selector=wait_for_selector,
            webhook=webhook,
        )

        data, response_headers = await self._request(
            "POST",
            "/screenshots",
            json_data=options.to_dict(),
        )

        request_id = response_headers.get("X-Request-Id")
        credits_used = response_headers.get("X-Credits-Used")
        credits_remaining = response_headers.get("X-Credits-Remaining")

        return ScreenshotResponse(
            success=True,
            data=data,
            url=url,
            content_type=response_headers.get("Content-Type"),
            request_id=request_id,
            credits_used=int(credits_used) if credits_used else None,
            credits_remaining=int(credits_remaining) if credits_remaining else None,
        )

    async def pdf(
        self,
        url: str,
        *,
        format: PdfFormat = "A4",
        landscape: bool = False,
        print_background: bool = True,
        margins: Optional[PdfMargins] = None,
        scale: float = 1.0,
        page_ranges: Optional[str] = None,
        header_template: Optional[str] = None,
        footer_template: Optional[str] = None,
        display_header_footer: bool = False,
        prefer_css_page_size: bool = False,
        accept_cookies: bool = False,
        delay: int = 0,
        wait_until: Literal["load", "domcontentloaded", "networkidle0", "networkidle2"] = "load",
        timeout: int = 30000,
        cookies: Optional[List[Cookie]] = None,
        headers: Optional[Dict[str, str]] = None,
        user_agent: Optional[str] = None,
        javascript_enabled: bool = True,
        wait_for_selector: Optional[str] = None,
        webhook: Optional[WebhookConfig] = None,
    ) -> PdfResponse:
        """
        Generate a PDF from a web page asynchronously.

        See ScreenCraft.pdf for full documentation.
        """
        options = PdfOptions(
            url=url,
            format=format,
            landscape=landscape,
            print_background=print_background,
            margins=margins,
            scale=scale,
            page_ranges=page_ranges,
            header_template=header_template,
            footer_template=footer_template,
            display_header_footer=display_header_footer,
            prefer_css_page_size=prefer_css_page_size,
            accept_cookies=accept_cookies,
            delay=delay,
            wait_until=wait_until,
            timeout=timeout,
            cookies=cookies,
            headers=headers,
            user_agent=user_agent,
            javascript_enabled=javascript_enabled,
            wait_for_selector=wait_for_selector,
            webhook=webhook,
        )

        data, response_headers = await self._request(
            "POST",
            "/pdfs",
            json_data=options.to_dict(),
        )

        request_id = response_headers.get("X-Request-Id")
        page_count = response_headers.get("X-Page-Count")
        credits_used = response_headers.get("X-Credits-Used")
        credits_remaining = response_headers.get("X-Credits-Remaining")

        return PdfResponse(
            success=True,
            data=data,
            url=url,
            content_type=response_headers.get("Content-Type"),
            request_id=request_id,
            page_count=int(page_count) if page_count else None,
            credits_used=int(credits_used) if credits_used else None,
            credits_remaining=int(credits_remaining) if credits_remaining else None,
        )

    async def get_account_info(self) -> AccountInfo:
        """
        Get current account information and usage asynchronously.

        See ScreenCraft.get_account_info for full documentation.
        """
        import json
        data, _ = await self._request("GET", "/account")
        return AccountInfo.from_dict(json.loads(data))

    async def close(self) -> None:
        """Close the HTTP session."""
        if self._session and not self._session.closed:
            await self._session.close()

    async def __aenter__(self) -> "AsyncScreenCraft":
        return self

    async def __aexit__(self, exc_type: Any, exc_val: Any, exc_tb: Any) -> None:
        await self.close()
