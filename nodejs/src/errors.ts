import type { ErrorResponse, RateLimitInfo } from './types';

/**
 * Base error class for all ScreenCraft errors
 */
export class ScreenCraftError extends Error {
  public readonly code: string;
  public readonly statusCode?: number;
  public readonly details?: Record<string, unknown>;

  constructor(message: string, code: string, statusCode?: number, details?: Record<string, unknown>) {
    super(message);
    this.name = 'ScreenCraftError';
    this.code = code;
    this.statusCode = statusCode;
    this.details = details;

    // Maintains proper stack trace for where our error was thrown (only available on V8)
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, this.constructor);
    }
  }

  /**
   * Convert error to JSON representation
   */
  toJSON(): ErrorResponse {
    return {
      code: this.code,
      message: this.message,
      details: this.details,
    };
  }
}

/**
 * Authentication-related errors
 */
export class AuthenticationError extends ScreenCraftError {
  constructor(message: string, details?: Record<string, unknown>) {
    super(message, 'AUTHENTICATION_ERROR', 401, details);
    this.name = 'AuthenticationError';
  }
}

/**
 * Validation errors for invalid request parameters
 */
export class ValidationError extends ScreenCraftError {
  constructor(message: string, details?: Record<string, unknown>) {
    super(message, 'VALIDATION_ERROR', 400, details);
    this.name = 'ValidationError';
  }
}

/**
 * Rate limiting errors
 */
export class RateLimitError extends ScreenCraftError {
  public readonly rateLimitInfo?: RateLimitInfo;

  constructor(message: string, rateLimitInfo?: RateLimitInfo, details?: Record<string, unknown>) {
    super(message, 'RATE_LIMIT_ERROR', 429, details);
    this.name = 'RateLimitError';
    this.rateLimitInfo = rateLimitInfo;
  }
}

/**
 * Errors related to capturing screenshots or PDFs
 */
export class CaptureError extends ScreenCraftError {
  constructor(message: string, statusCode?: number, details?: Record<string, unknown>) {
    super(message, 'CAPTURE_ERROR', statusCode, details);
    this.name = 'CaptureError';
  }
}

/**
 * Network-related errors
 */
export class NetworkError extends ScreenCraftError {
  constructor(message: string, details?: Record<string, unknown>) {
    super(message, 'NETWORK_ERROR', undefined, details);
    this.name = 'NetworkError';
  }
}

/**
 * Timeout errors
 */
export class TimeoutError extends ScreenCraftError {
  constructor(message: string, details?: Record<string, unknown>) {
    super(message, 'TIMEOUT_ERROR', 408, details);
    this.name = 'TimeoutError';
  }
}

/**
 * Server errors (5xx status codes)
 */
export class ServerError extends ScreenCraftError {
  constructor(message: string, statusCode?: number, details?: Record<string, unknown>) {
    super(message, 'SERVER_ERROR', statusCode, details);
    this.name = 'ServerError';
  }
}

/**
 * Unknown or unexpected errors
 */
export class UnknownError extends ScreenCraftError {
  constructor(message: string, details?: Record<string, unknown>) {
    super(message, 'UNKNOWN_ERROR', undefined, details);
    this.name = 'UnknownError';
  }
}

/**
 * Factory function to create appropriate error from API response
 */
export function createErrorFromResponse(
  statusCode: number,
  errorResponse?: ErrorResponse,
  rateLimitInfo?: RateLimitInfo
): ScreenCraftError {
  const message = errorResponse?.message || 'An unknown error occurred';
  const details = errorResponse?.details;

  switch (statusCode) {
    case 401:
    case 403:
      return new AuthenticationError(message, details);

    case 400:
      return new ValidationError(message, details);

    case 429:
      return new RateLimitError(message, rateLimitInfo, details);

    case 408:
      return new TimeoutError(message, details);

    case 500:
    case 502:
    case 503:
    case 504:
      return new ServerError(message, statusCode, details);

    default:
      if (statusCode >= 400 && statusCode < 500) {
        return new CaptureError(message, statusCode, details);
      }
      return new UnknownError(message, details);
  }
}
