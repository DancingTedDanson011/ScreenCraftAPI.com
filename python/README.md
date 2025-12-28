# ScreenCraft Python SDK

The official Python SDK for the [ScreenCraft API](https://screencraftapi.com) - capture screenshots and generate PDFs from web pages with ease.

[![PyPI version](https://badge.fury.io/py/screencraft.svg)](https://badge.fury.io/py/screencraft)
[![Python Versions](https://img.shields.io/pypi/pyversions/screencraft.svg)](https://pypi.org/project/screencraft/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- **Screenshot Capture** - Capture full-page or viewport screenshots in PNG, JPEG, or WebP format
- **PDF Generation** - Generate PDFs with customizable page sizes, margins, and headers/footers
- **Sync & Async** - Both synchronous and asynchronous clients available
- **Type Hints** - Full type annotations for excellent IDE support
- **Automatic Retries** - Built-in retry logic with exponential backoff
- **Viewport Presets** - Common device viewport presets included
- **Webhook Support** - Async processing with webhook callbacks

## Installation

```bash
pip install screencraft
```

For async support:

```bash
pip install screencraft[async]
```

## Quick Start

### Screenshot Capture

```python
from screencraft import ScreenCraft

# Initialize the client
client = ScreenCraft(api_key='your-api-key')

# Capture a screenshot
screenshot = client.screenshot(
    url='https://example.com',
    format='png',
    full_page=True,
    accept_cookies=True
)

# Save to file
with open('screenshot.png', 'wb') as f:
    f.write(screenshot.data)

print(f"Credits remaining: {screenshot.credits_remaining}")
```

### PDF Generation

```python
from screencraft import ScreenCraft, PdfMargins

client = ScreenCraft(api_key='your-api-key')

# Generate a PDF
pdf = client.pdf(
    url='https://example.com',
    format='A4',
    landscape=False,
    print_background=True,
    margins=PdfMargins(top='20mm', bottom='20mm', left='15mm', right='15mm')
)

# Save to file
with open('document.pdf', 'wb') as f:
    f.write(pdf.data)

print(f"Page count: {pdf.page_count}")
```

### Async Usage

```python
import asyncio
from screencraft import AsyncScreenCraft

async def main():
    async with AsyncScreenCraft(api_key='your-api-key') as client:
        # Capture multiple screenshots concurrently
        urls = [
            'https://example.com',
            'https://github.com',
            'https://python.org'
        ]

        tasks = [client.screenshot(url=url) for url in urls]
        screenshots = await asyncio.gather(*tasks)

        for i, screenshot in enumerate(screenshots):
            with open(f'screenshot_{i}.png', 'wb') as f:
                f.write(screenshot.data)

asyncio.run(main())
```

## Configuration Options

### Client Configuration

```python
from screencraft import ScreenCraft

client = ScreenCraft(
    api_key='your-api-key',
    base_url='https://screencraftapi.com/api/v1',  # Custom API URL
    timeout=60,           # Request timeout in seconds
    max_retries=3,        # Maximum retry attempts
    retry_delay=1.0,      # Initial retry delay
    retry_max_delay=30.0, # Maximum retry delay
    retry_backoff=2.0,    # Backoff multiplier
)
```

### Screenshot Options

```python
from screencraft import ScreenCraft, Viewport, Clip, Cookie

client = ScreenCraft(api_key='your-api-key')

screenshot = client.screenshot(
    url='https://example.com',

    # Image settings
    format='png',              # 'png', 'jpeg', 'webp'
    quality=80,                # 1-100 (jpeg/webp only)

    # Page capture
    full_page=True,            # Capture full scrollable page
    scroll_position='top',     # 'top' or 'bottom'

    # Viewport
    viewport=Viewport(
        width=1920,
        height=1080,
        device_scale_factor=2,
        is_mobile=False
    ),
    # Or use a preset: viewport='mobile', 'tablet', 'desktop', etc.

    # Clip region (partial screenshot)
    clip=Clip(x=0, y=0, width=800, height=600),

    # Timing
    delay=2000,                # Wait 2 seconds before capture
    timeout=30000,             # Navigation timeout
    wait_until='networkidle0', # 'load', 'domcontentloaded', 'networkidle0', 'networkidle2'
    wait_for_selector='#content',  # Wait for element

    # Interaction
    accept_cookies=True,       # Auto-dismiss cookie banners
    click_selector='#load-more',   # Click element before capture
    hide_selectors=['.ads', '.popup'],  # Hide elements

    # Browser settings
    javascript_enabled=True,
    dark_mode=False,
    block_ads=True,
    bypass_csp=False,

    # Custom headers and cookies
    headers={'Authorization': 'Bearer token'},
    cookies=[
        Cookie(name='session', value='abc123', domain='example.com')
    ],
    user_agent='Custom User Agent',
)
```

### PDF Options

```python
from screencraft import ScreenCraft, PdfMargins

client = ScreenCraft(api_key='your-api-key')

pdf = client.pdf(
    url='https://example.com',

    # Page settings
    format='A4',               # 'A0'-'A6', 'Letter', 'Legal', 'Tabloid'
    landscape=False,
    scale=1.0,                 # 0.1-2.0

    # Margins
    margins=PdfMargins(
        top='20mm',
        bottom='20mm',
        left='15mm',
        right='15mm'
    ),

    # Background
    print_background=True,

    # Page ranges
    page_ranges='1-5, 8, 11-13',

    # Header/Footer
    display_header_footer=True,
    header_template='<div style="font-size:10px">Header</div>',
    footer_template='<div style="font-size:10px">Page <span class="pageNumber"></span></div>',

    # Other options (same as screenshot)
    accept_cookies=True,
    delay=1000,
    timeout=30000,
    # ...
)
```

### Viewport Presets

The SDK includes common device viewport presets:

```python
from screencraft import ScreenCraft, VIEWPORT_PRESETS

client = ScreenCraft(api_key='your-api-key')

# Use a preset by name
screenshot = client.screenshot(url='https://example.com', viewport='mobile')

# Available presets:
# - desktop (1920x1080)
# - desktop_hd (2560x1440)
# - laptop (1366x768)
# - tablet (768x1024)
# - tablet_landscape (1024x768)
# - mobile (375x812)
# - mobile_landscape (812x375)
# - iphone_14 (390x844)
# - iphone_14_pro_max (430x932)
# - pixel_7 (412x915)
# - ipad_pro (1024x1366)

# Or access presets directly
print(VIEWPORT_PRESETS['mobile'])
```

### Webhook Support

For long-running operations, use webhooks to receive results asynchronously:

```python
from screencraft import ScreenCraft, WebhookConfig

client = ScreenCraft(api_key='your-api-key')

screenshot = client.screenshot(
    url='https://example.com',
    full_page=True,
    webhook=WebhookConfig(
        url='https://your-server.com/webhook',
        headers={'Authorization': 'Bearer webhook-secret'},
        secret='hmac-secret',  # For signature verification
        retry_count=3,
        timeout=30
    )
)
```

## Error Handling

The SDK provides specific exception classes for different error types:

```python
from screencraft import (
    ScreenCraft,
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

client = ScreenCraft(api_key='your-api-key')

try:
    screenshot = client.screenshot(url='https://example.com')
except AuthenticationError as e:
    print(f"Invalid API key: {e}")
except RateLimitError as e:
    print(f"Rate limit exceeded. Retry after: {e.retry_after}s")
except ValidationError as e:
    print(f"Invalid parameter '{e.field}': {e.message}")
except TimeoutError as e:
    print(f"Request timed out: {e}")
except ConnectionError as e:
    print(f"Connection failed: {e}")
except ServerError as e:
    print(f"Server error: {e}")
except RetryExhaustedError as e:
    print(f"All {e.attempts} retries failed: {e.last_exception}")
except ScreenCraftError as e:
    print(f"API error [{e.status_code}]: {e.message}")
```

## Account Information

```python
from screencraft import ScreenCraft

client = ScreenCraft(api_key='your-api-key')

account = client.get_account_info()
print(f"Plan: {account.plan}")
print(f"Credits remaining: {account.credits_remaining}")
print(f"Credits used: {account.credits_used}")
print(f"Reset date: {account.reset_date}")
```

## Context Manager

Both sync and async clients support context managers for automatic cleanup:

```python
# Sync
with ScreenCraft(api_key='your-api-key') as client:
    screenshot = client.screenshot(url='https://example.com')

# Async
async with AsyncScreenCraft(api_key='your-api-key') as client:
    screenshot = await client.screenshot(url='https://example.com')
```

## Type Hints

The SDK is fully typed and includes a `py.typed` marker for PEP 561 compliance:

```python
from screencraft import ScreenCraft, ScreenshotResponse, ImageFormat

def capture_page(url: str, format: ImageFormat = "png") -> bytes:
    client = ScreenCraft(api_key='your-api-key')
    response: ScreenshotResponse = client.screenshot(url=url, format=format)
    return response.data
```

## Requirements

- Python 3.8+
- `requests` (sync client)
- `aiohttp` (async client, optional)

## License

MIT License - see [LICENSE](LICENSE) for details.

## Links

- [Documentation](https://screencraftapi.com/docs)
- [API Reference](https://screencraftapi.com/api-reference)
- [GitHub Repository](https://github.com/screencraft/screencraft-python)
- [PyPI Package](https://pypi.org/project/screencraft/)
- [Support](mailto:support@screencraftapi.com)
