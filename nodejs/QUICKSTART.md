# ScreenCraft SDK Quick Start Guide

Get started with the ScreenCraft SDK in under 5 minutes.

## Installation

```bash
npm install screencraft
```

## Basic Usage

### 1. Initialize the Client

```typescript
import { ScreenCraftClient } from 'screencraft';

const client = new ScreenCraftClient({
  apiKey: 'your-api-key-here'
});
```

### 2. Capture a Screenshot

```typescript
import fs from 'fs';

// Capture screenshot
const screenshot = await client.screenshot({
  url: 'https://example.com',
  format: 'png'
});

// Save to file
const buffer = Buffer.from(screenshot.data, 'base64');
fs.writeFileSync('screenshot.png', buffer);
```

### 3. Generate a PDF

```typescript
// Generate PDF
const pdf = await client.pdf({
  url: 'https://example.com',
  format: 'A4'
});

// Save to file
const buffer = Buffer.from(pdf.data, 'base64');
fs.writeFileSync('document.pdf', buffer);
```

## Complete Example

```typescript
import { ScreenCraftClient } from 'screencraft';
import fs from 'fs';

async function main() {
  // Initialize client
  const client = new ScreenCraftClient({
    apiKey: process.env.SCREENCRAFT_API_KEY
  });

  try {
    // Capture full-page screenshot
    const screenshot = await client.screenshot({
      url: 'https://example.com',
      format: 'png',
      fullPage: true,
      viewport: {
        width: 1920,
        height: 1080
      },
      acceptCookies: true
    });

    // Save screenshot
    const buffer = Buffer.from(screenshot.data, 'base64');
    fs.writeFileSync('screenshot.png', buffer);

    console.log('Screenshot saved successfully!');
  } catch (error) {
    console.error('Error:', error.message);
  }
}

main();
```

## Environment Setup

### Using Environment Variables

Create a `.env` file:
```
SCREENCRAFT_API_KEY=your-api-key-here
```

Install dotenv:
```bash
npm install dotenv
```

Use in your code:
```typescript
import 'dotenv/config';
import { ScreenCraftClient } from 'screencraft';

const client = new ScreenCraftClient({
  apiKey: process.env.SCREENCRAFT_API_KEY
});
```

## Common Use Cases

### Full-Page Screenshot

```typescript
await client.screenshot({
  url: 'https://example.com',
  fullPage: true,
  format: 'png'
});
```

### Mobile Screenshot

```typescript
await client.screenshot({
  url: 'https://example.com',
  viewport: {
    width: 375,
    height: 812,
    deviceScaleFactor: 3,
    isMobile: true
  }
});
```

### High-Quality PDF

```typescript
await client.pdf({
  url: 'https://example.com',
  format: 'A4',
  printBackground: true,
  margin: {
    top: '1cm',
    bottom: '1cm',
    left: '1cm',
    right: '1cm'
  }
});
```

### Async with Webhook

```typescript
await client.screenshotAsync({
  url: 'https://example.com',
  format: 'png',
  webhook: {
    url: 'https://myapp.com/webhook',
    method: 'POST'
  }
});
```

## Error Handling

```typescript
import { RateLimitError, ValidationError } from 'screencraft';

try {
  const screenshot = await client.screenshot({
    url: 'https://example.com'
  });
} catch (error) {
  if (error instanceof RateLimitError) {
    console.error('Rate limit exceeded');
  } else if (error instanceof ValidationError) {
    console.error('Invalid parameters:', error.details);
  } else {
    console.error('Unexpected error:', error);
  }
}
```

## Next Steps

- Read the [full documentation](README.md)
- Explore [examples](examples/)
- Check the [API reference](README.md#api-reference)
- Review [error handling guide](README.md#error-handling)

## Support

- Documentation: https://screencraftapi.com/docs
- GitHub: https://github.com/screencraft/screencraft-nodejs-sdk
- Email: support@screencraftapi.com
