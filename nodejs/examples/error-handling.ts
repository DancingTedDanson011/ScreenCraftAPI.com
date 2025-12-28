/**
 * Error Handling Example
 *
 * This example demonstrates how to handle different types of errors
 * that can occur when using the ScreenCraft SDK.
 */

import {
  ScreenCraftClient,
  AuthenticationError,
  ValidationError,
  RateLimitError,
  CaptureError,
  NetworkError,
  TimeoutError,
  ServerError,
} from '../src';

async function main() {
  const client = new ScreenCraftClient({
    apiKey: process.env.SCREENCRAFT_API_KEY || 'your-api-key-here',
    timeout: 30000,
    maxRetries: 3,
  });

  try {
    console.log('Attempting to capture screenshot...');

    const screenshot = await client.screenshot({
      url: 'https://example.com',
      format: 'png',
    });

    console.log('Screenshot captured successfully!');
    console.log(`Size: ${screenshot.size} bytes`);
  } catch (error) {
    // Handle specific error types
    if (error instanceof AuthenticationError) {
      console.error('Authentication failed - check your API key');
      console.error(`Error code: ${error.code}`);
      console.error(`Message: ${error.message}`);
      console.error(`Status code: ${error.statusCode}`);
    } else if (error instanceof ValidationError) {
      console.error('Invalid request parameters');
      console.error(`Message: ${error.message}`);
      console.error('Details:', error.details);
    } else if (error instanceof RateLimitError) {
      console.error('Rate limit exceeded');
      console.error(`Message: ${error.message}`);

      if (error.rateLimitInfo) {
        const resetDate = new Date(error.rateLimitInfo.reset * 1000);
        console.error(`Limit: ${error.rateLimitInfo.limit}`);
        console.error(`Remaining: ${error.rateLimitInfo.remaining}`);
        console.error(`Resets at: ${resetDate.toISOString()}`);

        const waitTime = Math.max(0, error.rateLimitInfo.reset - Date.now() / 1000);
        console.error(`Wait ${Math.ceil(waitTime)} seconds before retrying`);
      }
    } else if (error instanceof CaptureError) {
      console.error('Failed to capture screenshot');
      console.error(`Message: ${error.message}`);
      console.error(`Status code: ${error.statusCode}`);
      console.error('Details:', error.details);
    } else if (error instanceof NetworkError) {
      console.error('Network error occurred');
      console.error(`Message: ${error.message}`);
      console.error('Check your internet connection and try again');
    } else if (error instanceof TimeoutError) {
      console.error('Request timed out');
      console.error(`Message: ${error.message}`);
      console.error('Try increasing the timeout or checking the target URL');
    } else if (error instanceof ServerError) {
      console.error('Server error occurred');
      console.error(`Message: ${error.message}`);
      console.error(`Status code: ${error.statusCode}`);
      console.error('The service may be experiencing issues. Try again later.');
    } else {
      console.error('Unknown error occurred');
      console.error(error);
    }

    process.exit(1);
  }
}

main();
