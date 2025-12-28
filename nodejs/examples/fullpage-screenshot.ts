/**
 * Full-Page Screenshot Example
 *
 * This example demonstrates how to capture a full-page screenshot
 * with custom viewport settings and cookie acceptance.
 */

import { ScreenCraftClient } from '../src';
import fs from 'fs';
import path from 'path';

async function main() {
  const client = new ScreenCraftClient({
    apiKey: process.env.SCREENCRAFT_API_KEY || 'your-api-key-here',
  });

  try {
    console.log('Capturing full-page screenshot...');

    const screenshot = await client.screenshot({
      url: 'https://example.com',
      format: 'png',
      fullPage: true,
      viewport: {
        width: 1920,
        height: 1080,
        deviceScaleFactor: 2, // High DPI for retina displays
      },
      acceptCookies: true,
      delay: 2000, // Wait 2 seconds for page to fully load
      blockAds: true,
    });

    console.log(`Screenshot captured successfully!`);
    console.log(`Processing time: ${screenshot.processingTime}ms`);
    console.log(`Size: ${screenshot.size} bytes`);

    const buffer = Buffer.from(screenshot.data, 'base64');
    const outputPath = path.join(__dirname, 'fullpage-screenshot.png');
    fs.writeFileSync(outputPath, buffer);

    console.log(`Full-page screenshot saved to: ${outputPath}`);
  } catch (error) {
    console.error('Error capturing screenshot:', error);
    process.exit(1);
  }
}

main();
