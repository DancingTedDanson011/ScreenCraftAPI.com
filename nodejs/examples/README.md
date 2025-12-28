# ScreenCraft SDK Examples

This directory contains example code demonstrating various features of the ScreenCraft SDK.

## Prerequisites

1. Install dependencies:
```bash
npm install
```

2. Build the SDK:
```bash
npm run build
```

3. Set your API key as an environment variable:
```bash
export SCREENCRAFT_API_KEY="your-api-key-here"
```

Or on Windows:
```cmd
set SCREENCRAFT_API_KEY=your-api-key-here
```

## Running Examples

You can run any example using `ts-node` or compile with TypeScript:

```bash
# Using ts-node (install with: npm install -g ts-node)
ts-node examples/basic-screenshot.ts

# Or compile and run
npx tsc examples/basic-screenshot.ts
node examples/basic-screenshot.js
```

## Available Examples

### basic-screenshot.ts
Simple screenshot capture and save to file.
```bash
ts-node examples/basic-screenshot.ts
```

### fullpage-screenshot.ts
Full-page screenshot with custom viewport and settings.
```bash
ts-node examples/fullpage-screenshot.ts
```

### generate-pdf.ts
Generate a PDF from a web page with custom formatting.
```bash
ts-node examples/generate-pdf.ts
```

### async-webhook.ts
Asynchronous screenshot capture with webhook notification.
```bash
ts-node examples/async-webhook.ts
```

### error-handling.ts
Comprehensive error handling demonstration.
```bash
ts-node examples/error-handling.ts
```

### mobile-screenshot.ts
Mobile device emulation for screenshots.
```bash
ts-node examples/mobile-screenshot.ts
```

## Output

All examples save their output files to the `examples/` directory:
- `screenshot.png` - Basic screenshot
- `fullpage-screenshot.png` - Full-page screenshot
- `document.pdf` - Generated PDF
- `mobile-screenshot.png` - Mobile screenshot

## Customization

Feel free to modify the examples to test different:
- URLs
- Viewport sizes
- Image formats
- PDF settings
- Webhook configurations
- Error handling scenarios

## Support

For more information, see the main [README.md](../README.md) or visit [https://screencraftapi.com/docs](https://screencraftapi.com/docs)
