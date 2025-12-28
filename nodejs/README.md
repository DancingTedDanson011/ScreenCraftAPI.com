# ScreenCraft Node.js/TypeScript SDK

Official Node.js and TypeScript SDK for the [ScreenCraft API](https://screencraftapi.com) - Professional screenshot and PDF generation service.

[![npm version](https://img.shields.io/npm/v/screencraft.svg)](https://www.npmjs.com/package/screencraft)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.0-blue.svg)](https://www.typescriptlang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- **Full TypeScript Support** - Complete type definitions for all API operations
- **Screenshot Capture** - High-quality screenshots in PNG, JPEG, or WebP formats
- **PDF Generation** - Convert web pages to PDF documents
- **Async Operations** - Webhook support for long-running captures
- **Automatic Retries** - Exponential backoff retry logic for resilient operations
- **Error Handling** - Typed error classes for precise error handling
- **Rate Limit Management** - Built-in rate limit detection and handling
- **Comprehensive Options** - Full control over viewport, cookies, delays, and more

## Installation

```bash
npm install screencraft
```

Or with yarn:

```bash
yarn add screencraft
```

## Quick Start

```typescript
import { ScreenCraftClient } from 'screencraft';
import fs from 'fs';

// Initialize the client
const client = new ScreenCraftClient({
  apiKey: 'your-api-key-here'
});

// Capture a screenshot
const screenshot = await client.screenshot({
  url: 'https://example.com',
  format: 'png',
  fullPage: true
});

// Save to file
const buffer = Buffer.from(screenshot.data, 'base64');
fs.writeFileSync('screenshot.png', buffer);
```

## Usage Examples

### Basic Screenshot

```typescript
const screenshot = await client.screenshot({
  url: 'https://example.com',
  format: 'png'
});

console.log(`Screenshot captured in ${screenshot.processingTime}ms`);
console.log(`Size: ${screenshot.size} bytes`);
```

### Full-Page Screenshot with Custom Viewport

```typescript
const screenshot = await client.screenshot({
  url: 'https://example.com',
  format: 'png',
  fullPage: true,
  viewport: {
    width: 1920,
    height: 1080,
    deviceScaleFactor: 2
  },
  acceptCookies: true,
  delay: 2000 // Wait 2 seconds before capturing
});
```

### High-Quality JPEG Screenshot

```typescript
const screenshot = await client.screenshot({
  url: 'https://example.com',
  format: 'jpeg',
  quality: 95,
  viewport: {
    width: 1920,
    height: 1080
  }
});
```

### Screenshot of Specific Element

```typescript
const screenshot = await client.screenshot({
  url: 'https://example.com',
  selector: '#main-content',
  format: 'png',
  omitBackground: true // Transparent background
});
```

### Mobile Device Screenshot

```typescript
const screenshot = await client.screenshot({
  url: 'https://example.com',
  viewport: {
    width: 375,
    height: 812,
    deviceScaleFactor: 3,
    isMobile: true,
    hasTouch: true
  },
  userAgent: 'Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X)'
});
```

### Generate PDF

```typescript
const pdf = await client.pdf({
  url: 'https://example.com',
  format: 'A4',
  printBackground: true,
  margin: {
    top: '1cm',
    right: '1cm',
    bottom: '1cm',
    left: '1cm'
  }
});

const buffer = Buffer.from(pdf.data, 'base64');
fs.writeFileSync('document.pdf', buffer);
```

### Landscape PDF with Custom Scale

```typescript
const pdf = await client.pdf({
  url: 'https://example.com',
  format: 'A4',
  landscape: true,
  scale: 0.8,
  printBackground: true
});
```

### Async Screenshot with Webhook

```typescript
const job = await client.screenshotAsync({
  url: 'https://example.com',
  format: 'png',
  fullPage: true,
  webhook: {
    url: 'https://myapp.com/webhooks/screenshot',
    method: 'POST',
    headers: {
      'X-Webhook-Secret': 'my-secret-key'
    }
  }
});

console.log('Job ID:', job.jobId);
console.log('Status:', job.status);
console.log('Estimated time:', job.estimatedTime, 'seconds');
```

### Async PDF with Webhook

```typescript
const job = await client.pdfAsync({
  url: 'https://example.com',
  format: 'A4',
  printBackground: true,
  webhook: {
    url: 'https://myapp.com/webhooks/pdf',
    method: 'POST'
  }
});
```

### Advanced Configuration

```typescript
const client = new ScreenCraftClient({
  apiKey: 'your-api-key',
  baseUrl: 'https://screencraftapi.com/api/v1', // Custom base URL
  timeout: 60000, // 60 second timeout
  maxRetries: 5, // Maximum retry attempts
  retryDelay: 2000 // Initial retry delay in ms
});
```

### Custom Headers and Geolocation

```typescript
const screenshot = await client.screenshot({
  url: 'https://example.com',
  headers: {
    'X-Custom-Header': 'value',
    'Authorization': 'Bearer token'
  },
  geolocation: {
    latitude: 51.5074,
    longitude: -0.1278,
    accuracy: 100
  },
  timezone: 'Europe/London',
  locale: 'en-GB'
});
```

### Block Ads and Track Position

```typescript
const screenshot = await client.screenshot({
  url: 'https://example.com',
  blockAds: true,
  scrollPosition: {
    x: 0,
    y: 500 // Scroll 500px down
  },
  delay: 1000
});
```

## Error Handling

The SDK provides typed error classes for precise error handling:

```typescript
import {
  ScreenCraftClient,
  AuthenticationError,
  ValidationError,
  RateLimitError,
  CaptureError,
  NetworkError,
  TimeoutError
} from 'screencraft';

const client = new ScreenCraftClient({ apiKey: 'your-api-key' });

try {
  const screenshot = await client.screenshot({
    url: 'https://example.com',
    format: 'png'
  });
} catch (error) {
  if (error instanceof AuthenticationError) {
    console.error('Invalid API key:', error.message);
  } else if (error instanceof ValidationError) {
    console.error('Invalid parameters:', error.message, error.details);
  } else if (error instanceof RateLimitError) {
    console.error('Rate limit exceeded');
    console.log('Limit:', error.rateLimitInfo?.limit);
    console.log('Remaining:', error.rateLimitInfo?.remaining);
    console.log('Reset at:', new Date(error.rateLimitInfo?.reset * 1000));
  } else if (error instanceof CaptureError) {
    console.error('Failed to capture:', error.message);
  } else if (error instanceof NetworkError) {
    console.error('Network error:', error.message);
  } else if (error instanceof TimeoutError) {
    console.error('Request timed out:', error.message);
  }
}
```

## API Reference

### `ScreenCraftClient`

#### Constructor

```typescript
new ScreenCraftClient(config: ScreenCraftConfig)
```

**Parameters:**
- `config.apiKey` (required): Your ScreenCraft API key
- `config.baseUrl` (optional): Custom API base URL (default: `https://screencraftapi.com/api/v1`)
- `config.timeout` (optional): Request timeout in milliseconds (default: `30000`)
- `config.maxRetries` (optional): Maximum retry attempts (default: `3`)
- `config.retryDelay` (optional): Initial retry delay in milliseconds (default: `1000`)

#### Methods

##### `screenshot(options: ScreenshotOptions): Promise<SyncResponse>`

Captures a screenshot synchronously.

**Options:**
- `url` (required): URL to capture
- `format`: Image format - `'png'`, `'jpeg'`, or `'webp'` (default: `'png'`)
- `quality`: Image quality 1-100 (JPEG/WebP only)
- `fullPage`: Capture full scrollable page (default: `false`)
- `viewport`: Viewport configuration
- `selector`: CSS selector for specific element
- `omitBackground`: Transparent background (default: `false`)
- `scrollPosition`: Scroll position before capture
- `acceptCookies`: Auto-accept cookie dialogs (default: `false`)
- `delay`: Delay before capture in milliseconds
- `headers`: Custom HTTP headers
- `userAgent`: Custom user agent string
- `blockAds`: Block advertisements (default: `false`)
- `geolocation`: Geolocation to emulate
- `timezone`: Timezone to emulate
- `locale`: Locale to emulate

##### `screenshotAsync(options: AsyncScreenshotOptions): Promise<AsyncResponse>`

Captures a screenshot asynchronously with webhook notification.

Same options as `screenshot()` plus:
- `webhook` (required): Webhook configuration
  - `url` (required): Webhook URL
  - `method`: HTTP method (default: `'POST'`)
  - `headers`: Custom headers

##### `pdf(options: PDFOptions): Promise<SyncResponse>`

Generates a PDF synchronously.

**Options:**
- `url` (required): URL to capture
- `format`: Page format - `'A4'`, `'Letter'`, `'Legal'`, `'Tabloid'`, `'A3'`, `'A5'` (default: `'A4'`)
- `landscape`: Landscape orientation (default: `false`)
- `printBackground`: Print background graphics (default: `false`)
- `scale`: Page scale (default: `1`)
- `margin`: Page margins (top, right, bottom, left)
- `pageRanges`: Page ranges to print (e.g., `'1-5, 8, 11-13'`)
- `preferCSSPageSize`: Use CSS-defined page size (default: `false`)
- Plus all base options from screenshot

##### `pdfAsync(options: AsyncPDFOptions): Promise<AsyncResponse>`

Generates a PDF asynchronously with webhook notification.

Same options as `pdf()` plus webhook configuration.

### Types

#### `SyncResponse`

```typescript
interface SyncResponse {
  data: string;              // Base64-encoded data
  contentType: string;       // Content type
  size: number;              // Size in bytes
  processingTime: number;    // Processing time in ms
}
```

#### `AsyncResponse`

```typescript
interface AsyncResponse {
  jobId: string;                                    // Unique job ID
  status: 'queued' | 'processing' | 'completed' | 'failed';
  webhookUrl: string;                               // Webhook URL
  estimatedTime?: number;                           // Estimated completion time
}
```

#### `Viewport`

```typescript
interface Viewport {
  width: number;
  height: number;
  deviceScaleFactor?: number;
  isMobile?: boolean;
  hasTouch?: boolean;
  isLandscape?: boolean;
}
```

## Rate Limiting

The SDK automatically handles rate limiting with exponential backoff. When a rate limit error occurs, you can access rate limit information:

```typescript
try {
  await client.screenshot({ url: 'https://example.com' });
} catch (error) {
  if (error instanceof RateLimitError) {
    const resetDate = new Date(error.rateLimitInfo.reset * 1000);
    console.log(`Rate limit resets at ${resetDate.toISOString()}`);
    console.log(`Remaining: ${error.rateLimitInfo.remaining}/${error.rateLimitInfo.limit}`);
  }
}
```

## Automatic Retries

The SDK implements automatic retry logic with exponential backoff for:
- Network errors
- Timeout errors
- Server errors (5xx)
- Rate limit errors (429)

Retries are **not** attempted for client errors (4xx) except rate limits.

## TypeScript Support

This SDK is written in TypeScript and provides complete type definitions. No additional `@types` packages are needed.

```typescript
import type { ScreenshotOptions, PDFOptions, SyncResponse } from 'screencraft';
```

## Requirements

- Node.js 14.0.0 or higher
- TypeScript 5.0 or higher (for TypeScript projects)

## License

MIT License - see [LICENSE](LICENSE) file for details

## Support

- Documentation: [https://screencraftapi.com/docs](https://screencraftapi.com/docs)
- Issues: [GitHub Issues](https://github.com/screencraft/screencraft-nodejs-sdk/issues)
- Email: support@screencraftapi.com

## Contributing

Contributions are welcome! Please read our contributing guidelines before submitting pull requests.

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for release history.
