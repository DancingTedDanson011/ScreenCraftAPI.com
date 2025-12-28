"""
ScreenCraft SDK - Type Definitions

This module defines all types, dataclasses, and type aliases used by the SDK.
"""

from dataclasses import dataclass, field
from typing import Optional, List, Dict, Any, Union, Literal
from enum import Enum


# Type aliases
ImageFormat = Literal["png", "jpeg", "webp"]
PdfFormat = Literal["A0", "A1", "A2", "A3", "A4", "A5", "A6", "Letter", "Legal", "Tabloid"]
ScrollPosition = Literal["top", "bottom"]


class ImageFormatEnum(str, Enum):
    """Supported image formats for screenshots."""
    PNG = "png"
    JPEG = "jpeg"
    WEBP = "webp"


class PdfFormatEnum(str, Enum):
    """Supported page formats for PDF generation."""
    A0 = "A0"
    A1 = "A1"
    A2 = "A2"
    A3 = "A3"
    A4 = "A4"
    A5 = "A5"
    A6 = "A6"
    LETTER = "Letter"
    LEGAL = "Legal"
    TABLOID = "Tabloid"


@dataclass
class Viewport:
    """Viewport configuration for screenshot capture."""
    width: int = 1920
    height: int = 1080
    device_scale_factor: float = 1.0
    is_mobile: bool = False
    has_touch: bool = False
    is_landscape: bool = True

    def to_dict(self) -> Dict[str, Any]:
        """Convert to API-compatible dictionary."""
        return {
            "width": self.width,
            "height": self.height,
            "deviceScaleFactor": self.device_scale_factor,
            "isMobile": self.is_mobile,
            "hasTouch": self.has_touch,
            "isLandscape": self.is_landscape,
        }


@dataclass
class Clip:
    """Clip region for partial screenshots."""
    x: int
    y: int
    width: int
    height: int

    def to_dict(self) -> Dict[str, int]:
        """Convert to API-compatible dictionary."""
        return {
            "x": self.x,
            "y": self.y,
            "width": self.width,
            "height": self.height,
        }


@dataclass
class Cookie:
    """Cookie to be set before capturing."""
    name: str
    value: str
    domain: Optional[str] = None
    path: str = "/"
    secure: bool = False
    http_only: bool = False
    same_site: Optional[Literal["Strict", "Lax", "None"]] = None
    expires: Optional[int] = None

    def to_dict(self) -> Dict[str, Any]:
        """Convert to API-compatible dictionary."""
        result: Dict[str, Any] = {
            "name": self.name,
            "value": self.value,
            "path": self.path,
            "secure": self.secure,
            "httpOnly": self.http_only,
        }
        if self.domain:
            result["domain"] = self.domain
        if self.same_site:
            result["sameSite"] = self.same_site
        if self.expires:
            result["expires"] = self.expires
        return result


@dataclass
class Header:
    """Custom HTTP header."""
    name: str
    value: str

    def to_dict(self) -> Dict[str, str]:
        """Convert to API-compatible dictionary."""
        return {"name": self.name, "value": self.value}


@dataclass
class WebhookConfig:
    """Configuration for webhook callbacks."""
    url: str
    headers: Optional[Dict[str, str]] = None
    secret: Optional[str] = None
    retry_count: int = 3
    timeout: int = 30

    def to_dict(self) -> Dict[str, Any]:
        """Convert to API-compatible dictionary."""
        result: Dict[str, Any] = {
            "url": self.url,
            "retryCount": self.retry_count,
            "timeout": self.timeout,
        }
        if self.headers:
            result["headers"] = self.headers
        if self.secret:
            result["secret"] = self.secret
        return result


@dataclass
class ScreenshotOptions:
    """Options for screenshot capture."""
    url: str
    format: ImageFormat = "png"
    quality: int = 80
    full_page: bool = False
    viewport: Optional[Viewport] = None
    clip: Optional[Clip] = None
    scroll_position: Optional[ScrollPosition] = None
    accept_cookies: bool = False
    delay: int = 0
    wait_until: Literal["load", "domcontentloaded", "networkidle0", "networkidle2"] = "load"
    timeout: int = 30000
    cookies: Optional[List[Cookie]] = None
    headers: Optional[Dict[str, str]] = None
    user_agent: Optional[str] = None
    bypass_csp: bool = False
    javascript_enabled: bool = True
    dark_mode: bool = False
    block_ads: bool = False
    hide_selectors: Optional[List[str]] = None
    click_selector: Optional[str] = None
    wait_for_selector: Optional[str] = None
    webhook: Optional[WebhookConfig] = None

    def to_dict(self) -> Dict[str, Any]:
        """Convert to API-compatible dictionary."""
        result: Dict[str, Any] = {
            "url": self.url,
            "format": self.format,
            "quality": self.quality,
            "fullPage": self.full_page,
            "acceptCookies": self.accept_cookies,
            "delay": self.delay,
            "waitUntil": self.wait_until,
            "timeout": self.timeout,
            "bypassCsp": self.bypass_csp,
            "javascriptEnabled": self.javascript_enabled,
            "darkMode": self.dark_mode,
            "blockAds": self.block_ads,
        }

        if self.viewport:
            result["viewport"] = self.viewport.to_dict()
        if self.clip:
            result["clip"] = self.clip.to_dict()
        if self.scroll_position:
            result["scrollPosition"] = self.scroll_position
        if self.cookies:
            result["cookies"] = [c.to_dict() for c in self.cookies]
        if self.headers:
            result["headers"] = self.headers
        if self.user_agent:
            result["userAgent"] = self.user_agent
        if self.hide_selectors:
            result["hideSelectors"] = self.hide_selectors
        if self.click_selector:
            result["clickSelector"] = self.click_selector
        if self.wait_for_selector:
            result["waitForSelector"] = self.wait_for_selector
        if self.webhook:
            result["webhook"] = self.webhook.to_dict()

        return result


@dataclass
class PdfMargins:
    """Margins for PDF generation."""
    top: str = "0"
    right: str = "0"
    bottom: str = "0"
    left: str = "0"

    def to_dict(self) -> Dict[str, str]:
        """Convert to API-compatible dictionary."""
        return {
            "top": self.top,
            "right": self.right,
            "bottom": self.bottom,
            "left": self.left,
        }


@dataclass
class PdfOptions:
    """Options for PDF generation."""
    url: str
    format: PdfFormat = "A4"
    landscape: bool = False
    print_background: bool = True
    margins: Optional[PdfMargins] = None
    scale: float = 1.0
    page_ranges: Optional[str] = None
    header_template: Optional[str] = None
    footer_template: Optional[str] = None
    display_header_footer: bool = False
    prefer_css_page_size: bool = False
    accept_cookies: bool = False
    delay: int = 0
    wait_until: Literal["load", "domcontentloaded", "networkidle0", "networkidle2"] = "load"
    timeout: int = 30000
    cookies: Optional[List[Cookie]] = None
    headers: Optional[Dict[str, str]] = None
    user_agent: Optional[str] = None
    javascript_enabled: bool = True
    wait_for_selector: Optional[str] = None
    webhook: Optional[WebhookConfig] = None

    def to_dict(self) -> Dict[str, Any]:
        """Convert to API-compatible dictionary."""
        result: Dict[str, Any] = {
            "url": self.url,
            "format": self.format,
            "landscape": self.landscape,
            "printBackground": self.print_background,
            "scale": self.scale,
            "displayHeaderFooter": self.display_header_footer,
            "preferCssPageSize": self.prefer_css_page_size,
            "acceptCookies": self.accept_cookies,
            "delay": self.delay,
            "waitUntil": self.wait_until,
            "timeout": self.timeout,
            "javascriptEnabled": self.javascript_enabled,
        }

        if self.margins:
            result["margins"] = self.margins.to_dict()
        if self.page_ranges:
            result["pageRanges"] = self.page_ranges
        if self.header_template:
            result["headerTemplate"] = self.header_template
        if self.footer_template:
            result["footerTemplate"] = self.footer_template
        if self.cookies:
            result["cookies"] = [c.to_dict() for c in self.cookies]
        if self.headers:
            result["headers"] = self.headers
        if self.user_agent:
            result["userAgent"] = self.user_agent
        if self.wait_for_selector:
            result["waitForSelector"] = self.wait_for_selector
        if self.webhook:
            result["webhook"] = self.webhook.to_dict()

        return result


@dataclass
class ScreenshotResponse:
    """Response from screenshot capture."""
    success: bool
    data: Optional[bytes] = None
    url: Optional[str] = None
    content_type: Optional[str] = None
    request_id: Optional[str] = None
    credits_used: Optional[int] = None
    credits_remaining: Optional[int] = None
    metadata: Optional[Dict[str, Any]] = None


@dataclass
class PdfResponse:
    """Response from PDF generation."""
    success: bool
    data: Optional[bytes] = None
    url: Optional[str] = None
    content_type: Optional[str] = None
    request_id: Optional[str] = None
    page_count: Optional[int] = None
    credits_used: Optional[int] = None
    credits_remaining: Optional[int] = None
    metadata: Optional[Dict[str, Any]] = None


@dataclass
class WebhookPayload:
    """Payload received from webhook callback."""
    request_id: str
    status: Literal["success", "error"]
    url: str
    result_url: Optional[str] = None
    error_message: Optional[str] = None
    error_code: Optional[str] = None
    timestamp: Optional[str] = None
    metadata: Optional[Dict[str, Any]] = None

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "WebhookPayload":
        """Create WebhookPayload from dictionary."""
        return cls(
            request_id=data.get("requestId", ""),
            status=data.get("status", "error"),
            url=data.get("url", ""),
            result_url=data.get("resultUrl"),
            error_message=data.get("errorMessage"),
            error_code=data.get("errorCode"),
            timestamp=data.get("timestamp"),
            metadata=data.get("metadata"),
        )


@dataclass
class AccountInfo:
    """User account information."""
    email: str
    plan: str
    credits_remaining: int
    credits_used: int
    credits_total: int
    reset_date: Optional[str] = None
    api_calls_this_month: int = 0

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "AccountInfo":
        """Create AccountInfo from dictionary."""
        return cls(
            email=data.get("email", ""),
            plan=data.get("plan", ""),
            credits_remaining=data.get("creditsRemaining", 0),
            credits_used=data.get("creditsUsed", 0),
            credits_total=data.get("creditsTotal", 0),
            reset_date=data.get("resetDate"),
            api_calls_this_month=data.get("apiCallsThisMonth", 0),
        )


# Predefined viewport presets
VIEWPORT_PRESETS: Dict[str, Viewport] = {
    "desktop": Viewport(width=1920, height=1080),
    "desktop_hd": Viewport(width=2560, height=1440),
    "laptop": Viewport(width=1366, height=768),
    "tablet": Viewport(width=768, height=1024, is_mobile=True, has_touch=True),
    "tablet_landscape": Viewport(width=1024, height=768, is_mobile=True, has_touch=True),
    "mobile": Viewport(width=375, height=812, is_mobile=True, has_touch=True, is_landscape=False),
    "mobile_landscape": Viewport(width=812, height=375, is_mobile=True, has_touch=True),
    "iphone_14": Viewport(width=390, height=844, is_mobile=True, has_touch=True, is_landscape=False, device_scale_factor=3),
    "iphone_14_pro_max": Viewport(width=430, height=932, is_mobile=True, has_touch=True, is_landscape=False, device_scale_factor=3),
    "pixel_7": Viewport(width=412, height=915, is_mobile=True, has_touch=True, is_landscape=False, device_scale_factor=2.625),
    "ipad_pro": Viewport(width=1024, height=1366, is_mobile=True, has_touch=True, is_landscape=False, device_scale_factor=2),
}
