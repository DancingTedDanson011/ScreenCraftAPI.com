# Changelog

All notable changes to the ScreenCraft Node.js SDK will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-12-28

### Added
- Initial release of the ScreenCraft Node.js/TypeScript SDK
- Screenshot capture with sync and async methods
- PDF generation with sync and async methods
- Full TypeScript support with complete type definitions
- Automatic retry logic with exponential backoff
- Comprehensive error handling with typed error classes
- Rate limit detection and handling
- Support for all ScreenCraft API parameters:
  - Viewport configuration
  - Full page capture
  - Image formats (PNG, JPEG, WebP)
  - PDF formats (A4, Letter, Legal, etc.)
  - Cookie acceptance
  - Custom delays
  - Scroll positioning
  - Element selection
  - Custom headers
  - Geolocation and timezone emulation
  - Ad blocking
- Webhook support for async operations
- Detailed documentation with usage examples

### Features
- `ScreenCraftClient` class for API interactions
- `screenshot()` method for synchronous screenshot capture
- `screenshotAsync()` method for asynchronous screenshot capture with webhooks
- `pdf()` method for synchronous PDF generation
- `pdfAsync()` method for asynchronous PDF generation with webhooks
- Typed error classes:
  - `AuthenticationError`
  - `ValidationError`
  - `RateLimitError`
  - `CaptureError`
  - `NetworkError`
  - `TimeoutError`
  - `ServerError`
  - `UnknownError`
- Complete TypeScript interfaces for all request and response types
