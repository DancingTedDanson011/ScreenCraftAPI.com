/**
 * Basic Screenshot Example
 *
 * This example demonstrates how to capture a simple screenshot
 * and save it to a file.
 */

import { ScreenCraftClient } from '../src';
import fs from 'fs';
import path from 'path';

async function main() {
  // Initialize the client with your API key
  const client = new ScreenCraftClient({
    apiKey: process.env.SCREENCRAFT_API_KEY || 'your-api-key-here',
  });

  try {
    console.log('Capturing screenshot...');

    // Capture a screenshot
    const screenshot = await client.screenshot({
      url: 'https://example.com',
      format: 'png',
    });

    console.log(`Screenshot captured successfully!`);
    console.log(`Processing time: ${screenshot.processingTime}ms`);
    console.log(`Size: ${screenshot.size} bytes`);

    // Convert base64 to buffer
    const buffer = Buffer.from(screenshot.data, 'base64');

    // Save to file
    const outputPath = path.join(__dirname, 'screenshot.png');
    fs.writeFileSync(outputPath, buffer);

    console.log(`Screenshot saved to: ${outputPath}`);
  } catch (error) {
    console.error('Error capturing screenshot:', error);
    process.exit(1);
  }
}

main();
