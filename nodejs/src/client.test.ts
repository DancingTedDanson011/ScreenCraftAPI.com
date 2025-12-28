/**
 * Unit tests for ScreenCraftClient
 *
 * These tests demonstrate the testing patterns for the SDK.
 * In a real implementation, you would use a mocking library like nock
 * to mock HTTP requests.
 */

import { ScreenCraftClient } from './client';
import {
  AuthenticationError,
  ValidationError,
  RateLimitError,
  CaptureError,
  NetworkError,
} from './errors';

describe('ScreenCraftClient', () => {
  describe('constructor', () => {
    it('should create client with valid config', () => {
      const client = new ScreenCraftClient({
        apiKey: 'test-api-key',
      });

      expect(client).toBeInstanceOf(ScreenCraftClient);
    });

    it('should throw error if API key is missing', () => {
      expect(() => {
        new ScreenCraftClient({ apiKey: '' });
      }).toThrow('API key is required');
    });

    it('should use default configuration values', () => {
      const client = new ScreenCraftClient({
        apiKey: 'test-api-key',
      });

      // Internal config should have defaults
      expect(client['config'].baseUrl).toBe('https://screencraftapi.com/api/v1');
      expect(client['config'].timeout).toBe(30000);
      expect(client['config'].maxRetries).toBe(3);
      expect(client['config'].retryDelay).toBe(1000);
    });

    it('should allow custom configuration values', () => {
      const client = new ScreenCraftClient({
        apiKey: 'test-api-key',
        baseUrl: 'https://custom.api.com',
        timeout: 60000,
        maxRetries: 5,
        retryDelay: 2000,
      });

      expect(client['config'].baseUrl).toBe('https://custom.api.com');
      expect(client['config'].timeout).toBe(60000);
      expect(client['config'].maxRetries).toBe(5);
      expect(client['config'].retryDelay).toBe(2000);
    });
  });

  describe('screenshot', () => {
    it('should validate URL', async () => {
      const client = new ScreenCraftClient({ apiKey: 'test-api-key' });

      await expect(
        client.screenshot({ url: 'invalid-url' })
      ).rejects.toThrow('Invalid URL provided');
    });

    it('should accept valid screenshot options', () => {
      const client = new ScreenCraftClient({ apiKey: 'test-api-key' });

      const options = {
        url: 'https://example.com',
        format: 'png' as const,
        fullPage: true,
        viewport: {
          width: 1920,
          height: 1080,
        },
      };

      // Should not throw
      expect(() => client['validateUrl'](options.url)).not.toThrow();
    });
  });

  describe('screenshotAsync', () => {
    it('should validate webhook URL', async () => {
      const client = new ScreenCraftClient({ apiKey: 'test-api-key' });

      await expect(
        client.screenshotAsync({
          url: 'https://example.com',
          webhook: {
            url: 'invalid-webhook-url',
          },
        })
      ).rejects.toThrow('Invalid webhook URL provided');
    });

    it('should accept valid webhook configuration', () => {
      const client = new ScreenCraftClient({ apiKey: 'test-api-key' });

      const webhook = {
        url: 'https://example.com/webhook',
        method: 'POST' as const,
        headers: { 'X-Custom': 'value' },
      };

      // Should not throw
      expect(() => client['validateWebhook'](webhook)).not.toThrow();
    });
  });

  describe('pdf', () => {
    it('should validate URL', async () => {
      const client = new ScreenCraftClient({ apiKey: 'test-api-key' });

      await expect(
        client.pdf({ url: 'not-a-url' })
      ).rejects.toThrow('Invalid URL provided');
    });

    it('should accept valid PDF options', () => {
      const client = new ScreenCraftClient({ apiKey: 'test-api-key' });

      const options = {
        url: 'https://example.com',
        format: 'A4' as const,
        printBackground: true,
        margin: {
          top: '1cm',
          right: '1cm',
          bottom: '1cm',
          left: '1cm',
        },
      };

      // Should not throw
      expect(() => client['validateUrl'](options.url)).not.toThrow();
    });
  });

  describe('retry logic', () => {
    it('should calculate exponential backoff delay', () => {
      const client = new ScreenCraftClient({ apiKey: 'test-api-key' });

      const delay0 = client['calculateRetryDelay'](0);
      const delay1 = client['calculateRetryDelay'](1);
      const delay2 = client['calculateRetryDelay'](2);

      // Should increase exponentially (with jitter)
      expect(delay1).toBeGreaterThan(delay0);
      expect(delay2).toBeGreaterThan(delay1);
    });

    it('should not retry on client errors', () => {
      const client = new ScreenCraftClient({ apiKey: 'test-api-key' });

      const error = {
        response: {
          status: 400,
        },
      };

      expect(client['shouldRetry'](error, 0)).toBe(false);
    });

    it('should retry on server errors', () => {
      const client = new ScreenCraftClient({ apiKey: 'test-api-key' });

      const error = {
        response: {
          status: 500,
        },
      };

      expect(client['shouldRetry'](error, 0)).toBe(true);
    });

    it('should retry on rate limit errors', () => {
      const client = new ScreenCraftClient({ apiKey: 'test-api-key' });

      const error = {
        response: {
          status: 429,
        },
      };

      expect(client['shouldRetry'](error, 0)).toBe(true);
    });

    it('should not retry if max retries exceeded', () => {
      const client = new ScreenCraftClient({
        apiKey: 'test-api-key',
        maxRetries: 3,
      });

      const error = {
        response: {
          status: 500,
        },
      };

      expect(client['shouldRetry'](error, 3)).toBe(false);
    });
  });

  describe('error handling', () => {
    it('should extract rate limit info from headers', () => {
      const client = new ScreenCraftClient({ apiKey: 'test-api-key' });

      const response = {
        headers: {
          'x-ratelimit-limit': '100',
          'x-ratelimit-remaining': '50',
          'x-ratelimit-reset': '1640000000',
        },
      };

      const rateLimitInfo = client['extractRateLimitInfo'](response as any);

      expect(rateLimitInfo).toEqual({
        limit: 100,
        remaining: 50,
        reset: 1640000000,
      });
    });

    it('should return undefined if rate limit headers missing', () => {
      const client = new ScreenCraftClient({ apiKey: 'test-api-key' });

      const response = {
        headers: {},
      };

      const rateLimitInfo = client['extractRateLimitInfo'](response as any);

      expect(rateLimitInfo).toBeUndefined();
    });
  });
});
