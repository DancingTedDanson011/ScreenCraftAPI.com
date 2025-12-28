/**
 * Mobile Screenshot Example
 *
 * This example demonstrates how to capture a screenshot
 * with mobile device emulation.
 */

import { ScreenCraftClient } from '../src';
import fs from 'fs';
import path from 'path';

async function main() {
  const client = new ScreenCraftClient({
    apiKey: process.env.SCREENCRAFT_API_KEY || 'your-api-key-here',
  });

  try {
    console.log('Capturing mobile screenshot (iPhone 13 Pro)...');

    const screenshot = await client.screenshot({
      url: 'https://example.com',
      format: 'png',
      fullPage: true,
      viewport: {
        width: 390,
        height: 844,
        deviceScaleFactor: 3,
        isMobile: true,
        hasTouch: true,
        isLandscape: false,
      },
      userAgent:
        'Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/604.1',
      acceptCookies: true,
      delay: 2000,
    });

    console.log(`Screenshot captured successfully!`);
    console.log(`Processing time: ${screenshot.processingTime}ms`);
    console.log(`Size: ${screenshot.size} bytes`);

    const buffer = Buffer.from(screenshot.data, 'base64');
    const outputPath = path.join(__dirname, 'mobile-screenshot.png');
    fs.writeFileSync(outputPath, buffer);

    console.log(`Mobile screenshot saved to: ${outputPath}`);
  } catch (error) {
    console.error('Error capturing mobile screenshot:', error);
    process.exit(1);
  }
}

main();
