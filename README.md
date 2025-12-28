# ScreenCraft API SDKs

Official SDKs for the [ScreenCraft API](https://screencraftapi.com) - a powerful screenshot and PDF generation service.

## Available SDKs

| Language | Package | Documentation |
|----------|---------|---------------|
| **Node.js** | [screencraft](./nodejs) | [README](./nodejs/README.md) |
| **Python** | [screencraft](./python) | [README](./python/README.md) |
| **Go** | [screencraft](./go) | [README](./go/README.md) |

## Quick Start

### Node.js

```bash
npm install screencraft
```

```typescript
import { ScreenCraft } from 'screencraft';

const client = new ScreenCraft({ apiKey: 'your-api-key' });

// Take a screenshot
const screenshot = await client.screenshots.capture({
  url: 'https://example.com',
  format: 'png',
  viewport: { width: 1920, height: 1080 }
});

// Save to file
await screenshot.toFile('screenshot.png');
```

### Python

```bash
pip install screencraft
```

```python
from screencraft import ScreenCraft

client = ScreenCraft(api_key='your-api-key')

# Take a screenshot
screenshot = client.screenshots.capture(
    url='https://example.com',
    format='png',
    viewport={'width': 1920, 'height': 1080}
)

# Save to file
screenshot.to_file('screenshot.png')
```

### Go

```bash
go get github.com/DancingTedDanson011/ScreenCraftAPI.com/go
```

```go
package main

import (
    "github.com/DancingTedDanson011/ScreenCraftAPI.com/go"
)

func main() {
    client := screencraft.New("your-api-key")

    // Take a screenshot
    screenshot, err := client.Screenshots.Capture(
        screencraft.WithURL("https://example.com"),
        screencraft.WithFormat("png"),
        screencraft.WithViewport(1920, 1080),
    )
    if err != nil {
        panic(err)
    }

    // Save to file
    err = screenshot.ToFile("screenshot.png")
}
```

## Features

All SDKs support the full ScreenCraft API:

- **Screenshots** - Capture web pages as PNG, JPEG, or WebP
- **PDFs** - Generate PDFs from URLs or HTML content
- **Full Page Capture** - Capture entire scrollable pages
- **Scroll Position** - Capture at specific scroll positions
- **Cookie Consent** - Automatic cookie banner acceptance
- **Custom Viewports** - Desktop, tablet, and mobile sizes
- **Resource Blocking** - Block images, scripts, fonts, etc.
- **Custom Headers/Cookies** - Set custom HTTP headers and cookies
- **Wait Conditions** - Wait for selectors, network idle, or delays

## API Documentation

Full API documentation is available at [screencraftapi.com/docs](https://screencraftapi.com/docs)

## Getting an API Key

1. Sign up at [screencraftapi.com](https://screencraftapi.com)
2. Go to Dashboard > API Keys
3. Create a new API key

## License

MIT License - see individual SDK directories for details.

## Support

- **Documentation**: [screencraftapi.com/docs](https://screencraftapi.com/docs)
- **Issues**: [GitHub Issues](https://github.com/DancingTedDanson011/ScreenCraftAPI.com/issues)
- **Email**: support@screencraftapi.com
