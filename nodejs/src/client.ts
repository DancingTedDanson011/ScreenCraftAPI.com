import axios, { AxiosInstance, AxiosError, AxiosResponse } from 'axios';
import type {
  ScreenCraftConfig,
  ScreenshotOptions,
  PDFOptions,
  AsyncScreenshotOptions,
  AsyncPDFOptions,
  SyncResponse,
  AsyncResponse,
  RateLimitInfo,
} from './types';
import {
  ScreenCraftError,
  NetworkError,
  createErrorFromResponse,
} from './errors';

/**
 * Default configuration values
 */
const DEFAULT_CONFIG = {
  baseUrl: 'https://screencraftapi.com/v1',
  timeout: 30000,
  maxRetries: 3,
  retryDelay: 1000,
};

/**
 * ScreenCraft API Client
 *
 * Official TypeScript/Node.js client for the ScreenCraft API.
 * Provides methods to capture screenshots and generate PDFs from URLs.
 */
export class ScreenCraftClient {
  private readonly config: Required<ScreenCraftConfig>;
  private readonly httpClient: AxiosInstance;

  /**
   * Creates a new ScreenCraft client instance
   *
   * @param config - Client configuration
   * @throws {ValidationError} If API key is missing
   *
   * @example
   * ```typescript
   * const client = new ScreenCraftClient({
   *   apiKey: 'your-api-key',
   *   timeout: 60000,
   *   maxRetries: 5
   * });
   * ```
   */
  constructor(config: ScreenCraftConfig) {
    if (!config.apiKey) {
      throw new Error('API key is required');
    }

    this.config = {
      ...DEFAULT_CONFIG,
      ...config,
    };

    this.httpClient = axios.create({
      baseURL: this.config.baseUrl,
      timeout: this.config.timeout,
      headers: {
        'Authorization': `Bearer ${this.config.apiKey}`,
        'Content-Type': 'application/json',
        'User-Agent': 'screencraft-nodejs-sdk/1.0.0',
      },
    });
  }

  /**
   * Captures a screenshot synchronously
   *
   * @param options - Screenshot options
   * @returns Screenshot data as base64-encoded string
   * @throws {ValidationError} If options are invalid
   * @throws {RateLimitError} If rate limit is exceeded
   * @throws {CaptureError} If screenshot capture fails
   *
   * @example
   * ```typescript
   * const screenshot = await client.screenshot({
   *   url: 'https://example.com',
   *   format: 'png',
   *   fullPage: true,
   *   viewport: { width: 1920, height: 1080 }
   * });
   *
   * // Save to file
   * const buffer = Buffer.from(screenshot.data, 'base64');
   * fs.writeFileSync('screenshot.png', buffer);
   * ```
   */
  async screenshot(options: ScreenshotOptions): Promise<SyncResponse> {
    this.validateUrl(options.url);

    return this.executeWithRetry(async () => {
      const response = await this.httpClient.post<SyncResponse>('/screenshots', options);
      return response.data;
    });
  }

  /**
   * Captures a screenshot asynchronously with webhook notification
   *
   * @param options - Async screenshot options including webhook configuration
   * @returns Job information
   * @throws {ValidationError} If options are invalid
   * @throws {RateLimitError} If rate limit is exceeded
   *
   * @example
   * ```typescript
   * const job = await client.screenshotAsync({
   *   url: 'https://example.com',
   *   format: 'png',
   *   fullPage: true,
   *   webhook: {
   *     url: 'https://myapp.com/webhook',
   *     method: 'POST',
   *     headers: { 'X-Custom-Header': 'value' }
   *   }
   * });
   *
   * console.log('Job ID:', job.jobId);
   * ```
   */
  async screenshotAsync(options: AsyncScreenshotOptions): Promise<AsyncResponse> {
    this.validateUrl(options.url);
    this.validateWebhook(options.webhook);

    return this.executeWithRetry(async () => {
      const response = await this.httpClient.post<AsyncResponse>('/screenshots', { ...options, async: true });
      return response.data;
    });
  }

  /**
   * Generates a PDF synchronously
   *
   * @param options - PDF generation options
   * @returns PDF data as base64-encoded string
   * @throws {ValidationError} If options are invalid
   * @throws {RateLimitError} If rate limit is exceeded
   * @throws {CaptureError} If PDF generation fails
   *
   * @example
   * ```typescript
   * const pdf = await client.pdf({
   *   url: 'https://example.com',
   *   format: 'A4',
   *   printBackground: true,
   *   margin: {
   *     top: '1cm',
   *     right: '1cm',
   *     bottom: '1cm',
   *     left: '1cm'
   *   }
   * });
   *
   * // Save to file
   * const buffer = Buffer.from(pdf.data, 'base64');
   * fs.writeFileSync('document.pdf', buffer);
   * ```
   */
  async pdf(options: PDFOptions): Promise<SyncResponse> {
    this.validateUrl(options.url);

    return this.executeWithRetry(async () => {
      const response = await this.httpClient.post<SyncResponse>('/pdfs', options);
      return response.data;
    });
  }

  /**
   * Generates a PDF asynchronously with webhook notification
   *
   * @param options - Async PDF options including webhook configuration
   * @returns Job information
   * @throws {ValidationError} If options are invalid
   * @throws {RateLimitError} If rate limit is exceeded
   *
   * @example
   * ```typescript
   * const job = await client.pdfAsync({
   *   url: 'https://example.com',
   *   format: 'A4',
   *   webhook: {
   *     url: 'https://myapp.com/webhook',
   *     method: 'POST'
   *   }
   * });
   *
   * console.log('Job ID:', job.jobId);
   * ```
   */
  async pdfAsync(options: AsyncPDFOptions): Promise<AsyncResponse> {
    this.validateUrl(options.url);
    this.validateWebhook(options.webhook);

    return this.executeWithRetry(async () => {
      const response = await this.httpClient.post<AsyncResponse>('/pdfs', { ...options, async: true });
      return response.data;
    });
  }

  /**
   * Validates a URL
   * @private
   */
  private validateUrl(url: string): void {
    try {
      new URL(url);
    } catch {
      throw new Error('Invalid URL provided');
    }
  }

  /**
   * Validates webhook configuration
   * @private
   */
  private validateWebhook(webhook: { url: string }): void {
    try {
      new URL(webhook.url);
    } catch {
      throw new Error('Invalid webhook URL provided');
    }
  }

  /**
   * Executes a request with exponential backoff retry logic
   * @private
   */
  private async executeWithRetry<T>(
    fn: () => Promise<T>,
    attempt: number = 0
  ): Promise<T> {
    try {
      return await fn();
    } catch (error) {
      const shouldRetry = this.shouldRetry(error, attempt);

      if (shouldRetry && attempt < this.config.maxRetries) {
        const delay = this.calculateRetryDelay(attempt);
        await this.sleep(delay);
        return this.executeWithRetry(fn, attempt + 1);
      }

      throw this.handleError(error);
    }
  }

  /**
   * Determines if a request should be retried
   * @private
   */
  private shouldRetry(error: unknown, attempt: number): boolean {
    if (attempt >= this.config.maxRetries) {
      return false;
    }

    if (axios.isAxiosError(error)) {
      const status = error.response?.status;

      // Don't retry client errors (except 429 rate limit)
      if (status && status >= 400 && status < 500 && status !== 429) {
        return false;
      }

      // Retry on network errors, timeouts, and 5xx errors
      return (
        !error.response || // Network error
        error.code === 'ECONNABORTED' || // Timeout
        (status !== undefined && status >= 500) || // Server error
        status === 429 // Rate limit
      );
    }

    return false;
  }

  /**
   * Calculates retry delay with exponential backoff
   * @private
   */
  private calculateRetryDelay(attempt: number): number {
    const exponentialDelay = this.config.retryDelay * Math.pow(2, attempt);
    const jitter = Math.random() * 1000; // Add jitter to prevent thundering herd
    return exponentialDelay + jitter;
  }

  /**
   * Sleeps for the specified duration
   * @private
   */
  private sleep(ms: number): Promise<void> {
    return new Promise((resolve) => setTimeout(resolve, ms));
  }

  /**
   * Handles errors and converts them to ScreenCraft error types
   * @private
   */
  private handleError(error: unknown): ScreenCraftError {
    if (error instanceof ScreenCraftError) {
      return error;
    }

    if (axios.isAxiosError(error)) {
      const axiosError = error as AxiosError;

      // Network error
      if (!axiosError.response) {
        return new NetworkError(
          axiosError.message || 'Network error occurred',
          { originalError: axiosError.code }
        );
      }

      // Extract rate limit info if available
      const rateLimitInfo = this.extractRateLimitInfo(axiosError.response);

      // API error response
      const statusCode = axiosError.response.status;
      const errorResponse = axiosError.response.data as any;

      return createErrorFromResponse(statusCode, errorResponse, rateLimitInfo);
    }

    // Unknown error
    return new NetworkError(
      error instanceof Error ? error.message : 'An unknown error occurred',
      { originalError: error }
    );
  }

  /**
   * Extracts rate limit information from response headers
   * @private
   */
  private extractRateLimitInfo(response: AxiosResponse): RateLimitInfo | undefined {
    const headers = response.headers;

    if (
      headers['x-ratelimit-limit'] &&
      headers['x-ratelimit-remaining'] &&
      headers['x-ratelimit-reset']
    ) {
      return {
        limit: parseInt(headers['x-ratelimit-limit'], 10),
        remaining: parseInt(headers['x-ratelimit-remaining'], 10),
        reset: parseInt(headers['x-ratelimit-reset'], 10),
      };
    }

    return undefined;
  }
}
