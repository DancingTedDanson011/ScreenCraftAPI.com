/**
 * Async Screenshot with Webhook Example
 *
 * This example demonstrates how to capture a screenshot asynchronously
 * and receive the result via webhook.
 */

import { ScreenCraftClient } from '../src';

async function main() {
  const client = new ScreenCraftClient({
    apiKey: process.env.SCREENCRAFT_API_KEY || 'your-api-key-here',
  });

  try {
    console.log('Starting async screenshot job...');

    const job = await client.screenshotAsync({
      url: 'https://example.com',
      format: 'png',
      fullPage: true,
      viewport: {
        width: 1920,
        height: 1080,
      },
      webhook: {
        url: 'https://myapp.com/webhooks/screenshot',
        method: 'POST',
        headers: {
          'X-Webhook-Secret': 'my-secret-key',
          'Content-Type': 'application/json',
        },
      },
    });

    console.log('Job created successfully!');
    console.log(`Job ID: ${job.jobId}`);
    console.log(`Status: ${job.status}`);
    console.log(`Webhook URL: ${job.webhookUrl}`);
    if (job.estimatedTime) {
      console.log(`Estimated completion: ${job.estimatedTime} seconds`);
    }

    console.log('\nThe screenshot will be sent to your webhook URL when complete.');
  } catch (error) {
    console.error('Error creating async job:', error);
    process.exit(1);
  }
}

main();
