/**
 * ScreenCraft Node.js/TypeScript SDK
 *
 * Official SDK for the ScreenCraft API - Professional screenshot and PDF generation service
 *
 * @packageDocumentation
 */

export { ScreenCraftClient } from './client';

export type {
  ScreenCraftConfig,
  ImageFormat,
  PDFFormat,
  Viewport,
  ScrollPosition,
  BaseOptions,
  ScreenshotOptions,
  PDFOptions,
  WebhookConfig,
  AsyncScreenshotOptions,
  AsyncPDFOptions,
  SyncResponse,
  AsyncResponse,
  ErrorResponse,
  RateLimitInfo,
} from './types';

export {
  ScreenCraftError,
  AuthenticationError,
  ValidationError,
  RateLimitError,
  CaptureError,
  NetworkError,
  TimeoutError,
  ServerError,
  UnknownError,
} from './errors';

// Default export for convenience
export { ScreenCraftClient as default } from './client';
