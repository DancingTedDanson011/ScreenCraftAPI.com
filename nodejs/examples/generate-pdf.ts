/**
 * PDF Generation Example
 *
 * This example demonstrates how to generate a PDF from a web page
 * with custom formatting options.
 */

import { ScreenCraftClient } from '../src';
import fs from 'fs';
import path from 'path';

async function main() {
  const client = new ScreenCraftClient({
    apiKey: process.env.SCREENCRAFT_API_KEY || 'your-api-key-here',
  });

  try {
    console.log('Generating PDF...');

    const pdf = await client.pdf({
      url: 'https://example.com',
      format: 'A4',
      printBackground: true,
      margin: {
        top: '1cm',
        right: '1cm',
        bottom: '1cm',
        left: '1cm',
      },
      scale: 1.0,
      landscape: false,
      acceptCookies: true,
      delay: 2000,
    });

    console.log(`PDF generated successfully!`);
    console.log(`Processing time: ${pdf.processingTime}ms`);
    console.log(`Size: ${pdf.size} bytes`);

    const buffer = Buffer.from(pdf.data, 'base64');
    const outputPath = path.join(__dirname, 'document.pdf');
    fs.writeFileSync(outputPath, buffer);

    console.log(`PDF saved to: ${outputPath}`);
  } catch (error) {
    console.error('Error generating PDF:', error);
    process.exit(1);
  }
}

main();
