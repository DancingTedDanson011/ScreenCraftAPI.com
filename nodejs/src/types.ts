/**
 * Supported image formats for screenshots
 */
export type ImageFormat = 'png' | 'jpeg' | 'webp';

/**
 * Supported PDF page formats
 */
export type PDFFormat = 'A4' | 'Letter' | 'Legal' | 'Tabloid' | 'A3' | 'A5';

/**
 * Viewport configuration
 */
export interface Viewport {
  width: number;
  height: number;
  deviceScaleFactor?: number;
  isMobile?: boolean;
  hasTouch?: boolean;
  isLandscape?: boolean;
}

/**
 * Scroll position configuration
 */
export interface ScrollPosition {
  x: number;
  y: number;
}

/**
 * Client configuration options
 */
export interface ScreenCraftConfig {
  /**
   * API key for authentication
   */
  apiKey: string;

  /**
   * Base URL for the API (default: https://screencraftapi.com/api/v1)
   */
  baseUrl?: string;

  /**
   * Request timeout in milliseconds (default: 30000)
   */
  timeout?: number;

  /**
   * Maximum number of retry attempts (default: 3)
   */
  maxRetries?: number;

  /**
   * Initial retry delay in milliseconds (default: 1000)
   */
  retryDelay?: number;
}

/**
 * Base options shared between screenshot and PDF requests
 */
export interface BaseOptions {
  /**
   * URL to capture
   */
  url: string;

  /**
   * Viewport configuration
   */
  viewport?: Viewport;

  /**
   * Whether to capture the full scrollable page
   */
  fullPage?: boolean;

  /**
   * Whether to automatically accept cookie consent dialogs
   */
  acceptCookies?: boolean;

  /**
   * Delay before capture in milliseconds
   */
  delay?: number;

  /**
   * Custom HTTP headers to send with the request
   */
  headers?: Record<string, string>;

  /**
   * User agent string
   */
  userAgent?: string;

  /**
   * Whether to block ads
   */
  blockAds?: boolean;

  /**
   * Geolocation to emulate
   */
  geolocation?: {
    latitude: number;
    longitude: number;
    accuracy?: number;
  };

  /**
   * Timezone to emulate
   */
  timezone?: string;

  /**
   * Locale to emulate
   */
  locale?: string;
}

/**
 * Screenshot-specific options
 */
export interface ScreenshotOptions extends BaseOptions {
  /**
   * Image format (default: 'png')
   */
  format?: ImageFormat;

  /**
   * Image quality (1-100, only for JPEG and WebP)
   */
  quality?: number;

  /**
   * Scroll position before capture
   */
  scrollPosition?: ScrollPosition;

  /**
   * Whether to omit the background (transparent)
   */
  omitBackground?: boolean;

  /**
   * CSS selector to capture a specific element
   */
  selector?: string;
}

/**
 * PDF-specific options
 */
export interface PDFOptions extends BaseOptions {
  /**
   * Page format (default: 'A4')
   */
  format?: PDFFormat;

  /**
   * Print background graphics
   */
  printBackground?: boolean;

  /**
   * Page margins
   */
  margin?: {
    top?: string;
    right?: string;
    bottom?: string;
    left?: string;
  };

  /**
   * Page ranges to print (e.g., '1-5, 8, 11-13')
   */
  pageRanges?: string;

  /**
   * Scale of the webpage rendering (default: 1)
   */
  scale?: number;

  /**
   * Prefer CSS page size
   */
  preferCSSPageSize?: boolean;

  /**
   * Landscape orientation
   */
  landscape?: boolean;
}

/**
 * Webhook configuration for async operations
 */
export interface WebhookConfig {
  /**
   * Webhook URL to receive the result
   */
  url: string;

  /**
   * HTTP method (default: 'POST')
   */
  method?: 'POST' | 'PUT' | 'PATCH';

  /**
   * Custom headers to send with the webhook request
   */
  headers?: Record<string, string>;
}

/**
 * Options for async screenshot requests
 */
export interface AsyncScreenshotOptions extends ScreenshotOptions {
  /**
   * Webhook configuration
   */
  webhook: WebhookConfig;
}

/**
 * Options for async PDF requests
 */
export interface AsyncPDFOptions extends PDFOptions {
  /**
   * Webhook configuration
   */
  webhook: WebhookConfig;
}

/**
 * Synchronous API response containing the captured data
 */
export interface SyncResponse {
  /**
   * Base64-encoded image or PDF data
   */
  data: string;

  /**
   * Content type of the response
   */
  contentType: string;

  /**
   * Size of the response in bytes
   */
  size: number;

  /**
   * Time taken to generate in milliseconds
   */
  processingTime: number;
}

/**
 * Asynchronous API response with job information
 */
export interface AsyncResponse {
  /**
   * Unique job ID
   */
  jobId: string;

  /**
   * Current job status
   */
  status: 'queued' | 'processing' | 'completed' | 'failed';

  /**
   * Webhook URL that will receive the result
   */
  webhookUrl: string;

  /**
   * Estimated time to completion in seconds
   */
  estimatedTime?: number;
}

/**
 * Error response from the API
 */
export interface ErrorResponse {
  /**
   * Error code
   */
  code: string;

  /**
   * Human-readable error message
   */
  message: string;

  /**
   * Additional error details
   */
  details?: Record<string, unknown>;
}

/**
 * Rate limit information
 */
export interface RateLimitInfo {
  /**
   * Maximum requests allowed per window
   */
  limit: number;

  /**
   * Remaining requests in current window
   */
  remaining: number;

  /**
   * Unix timestamp when the rate limit resets
   */
  reset: number;
}
