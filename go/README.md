# ScreenCraft Go SDK

Official Go SDK for the [ScreenCraft API](https://screencraftapi.com) - capture screenshots and generate PDFs from web pages.

## Installation

```bash
go get github.com/DancingTedDanson011/screencraft-go
```

## Quick Start

```go
package main

import (
    "context"
    "log"
    "os"

    "github.com/DancingTedDanson011/screencraft-go"
)

func main() {
    // Create a client with your API key
    client := screencraft.New("your-api-key")

    // Capture a screenshot
    result, err := client.Screenshot(context.Background(), &screencraft.ScreenshotOptions{
        URL:      "https://example.com",
        Format:   screencraft.FormatPNG,
        FullPage: true,
    })
    if err != nil {
        log.Fatal(err)
    }

    // Save the screenshot
    if err := os.WriteFile("screenshot.png", result.Data, 0644); err != nil {
        log.Fatal(err)
    }
}
```

## Features

- Screenshot capture (PNG, JPEG, WebP)
- PDF generation with customizable options
- Viewport and device emulation
- Cookie consent auto-acceptance
- Custom delays and wait conditions
- Webhook support for async operations
- Automatic retries with exponential backoff
- Context support for cancellation
- Comprehensive error handling

## Configuration

### Client Options

```go
client := screencraft.New("your-api-key",
    screencraft.WithTimeout(120*time.Second),
    screencraft.WithMaxRetries(5),
    screencraft.WithRetryWait(2*time.Second, 60*time.Second),
    screencraft.WithDebug(true),
    screencraft.WithLogger(log.Default()),
)
```

### Available Options

| Option | Description |
|--------|-------------|
| `WithBaseURL(url)` | Set a custom API base URL |
| `WithHTTPClient(client)` | Use a custom HTTP client |
| `WithTimeout(duration)` | Set HTTP client timeout |
| `WithMaxRetries(n)` | Set maximum retry attempts |
| `WithRetryWait(min, max)` | Set retry wait bounds |
| `WithUserAgent(ua)` | Set custom User-Agent |
| `WithDebug(bool)` | Enable debug logging |
| `WithLogger(logger)` | Set custom logger |

## Screenshots

### Basic Screenshot

```go
result, err := client.Screenshot(ctx, &screencraft.ScreenshotOptions{
    URL:    "https://example.com",
    Format: screencraft.FormatPNG,
})
```

### Full Page Screenshot

```go
result, err := client.Screenshot(ctx, &screencraft.ScreenshotOptions{
    URL:      "https://example.com",
    Format:   screencraft.FormatPNG,
    FullPage: true,
})
```

### Custom Viewport

```go
result, err := client.Screenshot(ctx, &screencraft.ScreenshotOptions{
    URL:    "https://example.com",
    Format: screencraft.FormatPNG,
    Viewport: &screencraft.Viewport{
        Width:  1920,
        Height: 1080,
    },
})
```

### Mobile Emulation

```go
result, err := client.Screenshot(ctx, &screencraft.ScreenshotOptions{
    URL:    "https://example.com",
    Format: screencraft.FormatPNG,
    Viewport: &screencraft.Viewport{
        Width:  375,
        Height: 812,
    },
    IsMobile:          true,
    HasTouch:          true,
    DeviceScaleFactor: 3,
})

// Or use the convenience method
result, err := client.ScreenshotMobile(ctx, "https://example.com")
```

### With Cookie Consent

```go
result, err := client.Screenshot(ctx, &screencraft.ScreenshotOptions{
    URL:           "https://example.com",
    Format:        screencraft.FormatPNG,
    AcceptCookies: true,
})
```

### Wait for Selector

```go
result, err := client.Screenshot(ctx, &screencraft.ScreenshotOptions{
    URL:             "https://example.com",
    Format:          screencraft.FormatPNG,
    WaitForSelector: "#main-content",
    WaitUntil:       screencraft.WaitNetworkIdle,
})
```

### JPEG with Quality

```go
result, err := client.Screenshot(ctx, &screencraft.ScreenshotOptions{
    URL:     "https://example.com",
    Format:  screencraft.FormatJPEG,
    Quality: 85,
})
```

### Dark Mode

```go
result, err := client.Screenshot(ctx, &screencraft.ScreenshotOptions{
    URL:      "https://example.com",
    Format:   screencraft.FormatPNG,
    DarkMode: true,
})
```

### With Custom Headers and Cookies

```go
result, err := client.Screenshot(ctx, &screencraft.ScreenshotOptions{
    URL:    "https://example.com",
    Format: screencraft.FormatPNG,
    Headers: []screencraft.Header{
        {Name: "Authorization", Value: "Bearer token123"},
    },
    Cookies: []screencraft.Cookie{
        {Name: "session", Value: "abc123", Domain: "example.com"},
    },
})
```

## PDF Generation

### Basic PDF

```go
result, err := client.PDF(ctx, &screencraft.PDFOptions{
    URL:    "https://example.com",
    Format: screencraft.A4,
})
```

### Landscape PDF

```go
result, err := client.PDF(ctx, &screencraft.PDFOptions{
    URL:         "https://example.com",
    Format:      screencraft.A4,
    Orientation: screencraft.Landscape,
})
```

### Custom Margins

```go
result, err := client.PDF(ctx, &screencraft.PDFOptions{
    URL:    "https://example.com",
    Format: screencraft.A4,
    Margin: &screencraft.PDFMargin{
        Top:    "1in",
        Right:  "0.5in",
        Bottom: "1in",
        Left:   "0.5in",
    },
})
```

### Header and Footer

```go
result, err := client.PDF(ctx, &screencraft.PDFOptions{
    URL:                 "https://example.com",
    Format:              screencraft.A4,
    DisplayHeaderFooter: true,
    HeaderTemplate:      `<div style="font-size:10px; text-align:center; width:100%;">My Document</div>`,
    FooterTemplate:      `<div style="font-size:10px; text-align:center; width:100%;">Page <span class="pageNumber"></span> of <span class="totalPages"></span></div>`,
    Margin: &screencraft.PDFMargin{
        Top:    "100px",
        Bottom: "100px",
    },
})
```

### Page Ranges

```go
result, err := client.PDF(ctx, &screencraft.PDFOptions{
    URL:        "https://example.com",
    Format:     screencraft.A4,
    PageRanges: "1-5, 8, 11-13",
})
```

### Print Background

```go
result, err := client.PDF(ctx, &screencraft.PDFOptions{
    URL:             "https://example.com",
    Format:          screencraft.A4,
    PrintBackground: true,
})
```

## Async Operations with Webhooks

```go
// Screenshot with webhook
jobID, err := client.ScreenshotAsync(ctx, &screencraft.ScreenshotOptions{
    URL:    "https://example.com",
    Format: screencraft.FormatPNG,
    Webhook: &screencraft.WebhookConfig{
        URL: "https://yoursite.com/webhook",
        Headers: map[string]string{
            "Authorization": "Bearer your-webhook-secret",
        },
        Secret: "webhook-signature-secret",
    },
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Job ID: %s\n", jobID)

// PDF with webhook
jobID, err := client.PDFAsync(ctx, &screencraft.PDFOptions{
    URL:    "https://example.com",
    Format: screencraft.A4,
    Webhook: &screencraft.WebhookConfig{
        URL: "https://yoursite.com/webhook",
    },
})
```

## Error Handling

```go
result, err := client.Screenshot(ctx, opts)
if err != nil {
    switch {
    case screencraft.IsAuthenticationError(err):
        log.Fatal("Invalid API key")
    case screencraft.IsRateLimitError(err):
        var rateErr *screencraft.RateLimitError
        if errors.As(err, &rateErr) {
            log.Printf("Rate limited. Retry after: %s", rateErr.RetryAfter)
        }
    case screencraft.IsValidationError(err):
        var valErr *screencraft.ValidationError
        if errors.As(err, &valErr) {
            log.Printf("Validation failed for field: %s", valErr.Field)
        }
    case screencraft.IsTimeoutError(err):
        log.Println("Request timed out")
    case screencraft.IsNetworkError(err):
        log.Println("Network error occurred")
    case screencraft.IsServerError(err):
        log.Println("Server error occurred")
    default:
        log.Printf("Error: %v", err)
    }
}
```

### Checking Retryable Errors

```go
if screencraft.IsRetryable(err) {
    retryAfter := screencraft.GetRetryAfter(err)
    time.Sleep(retryAfter)
    // Retry the request
}
```

## Context Support

```go
// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

result, err := client.Screenshot(ctx, opts)

// With cancellation
ctx, cancel := context.WithCancel(context.Background())
go func() {
    time.Sleep(5 * time.Second)
    cancel() // Cancel after 5 seconds
}()

result, err := client.Screenshot(ctx, opts)
```

## Rate Limiting

The SDK automatically handles rate limiting with exponential backoff. You can also check rate limit information:

```go
info := client.GetRateLimitInfo()
if info != nil {
    fmt.Printf("Limit: %d, Remaining: %d, Reset: %s\n",
        info.Limit, info.Remaining, info.Reset)
}
```

## Convenience Methods

### Screenshot Convenience Methods

```go
// Simple URL screenshot
result, err := client.ScreenshotURL(ctx, "https://example.com")

// Full page screenshot
result, err := client.ScreenshotFullPage(ctx, "https://example.com", screencraft.FormatPNG)

// Mobile screenshot
result, err := client.ScreenshotMobile(ctx, "https://example.com")

// Desktop screenshot (1920x1080)
result, err := client.ScreenshotDesktop(ctx, "https://example.com")

// Screenshot with delay
result, err := client.ScreenshotWithDelay(ctx, "https://example.com", 2000)

// Screenshot with cookie consent
result, err := client.ScreenshotWithCookieConsent(ctx, "https://example.com")
```

### PDF Convenience Methods

```go
// Simple URL to PDF
result, err := client.PDFURL(ctx, "https://example.com")

// A4 PDF
result, err := client.PDFA4(ctx, "https://example.com")

// US Letter PDF
result, err := client.PDFLetter(ctx, "https://example.com")

// Landscape PDF
result, err := client.PDFLandscape(ctx, "https://example.com", screencraft.A4)

// PDF with custom margins
result, err := client.PDFWithMargins(ctx, "https://example.com", &screencraft.PDFMargin{
    Top: "1in", Right: "1in", Bottom: "1in", Left: "1in",
})

// PDF with header and footer
result, err := client.PDFWithHeaderFooter(ctx, "https://example.com",
    "<div>Header</div>",
    "<div>Footer</div>",
)

// PDF with page range
result, err := client.PDFPageRange(ctx, "https://example.com", "1-5")

// PDF with cookie consent
result, err := client.PDFWithCookieConsent(ctx, "https://example.com")
```

## Type Constants

### Image Formats

```go
screencraft.FormatPNG   // PNG format
screencraft.FormatJPEG  // JPEG format
screencraft.FormatWebP  // WebP format
```

### PDF Formats

```go
screencraft.A4      // A4 paper (210mm x 297mm)
screencraft.A3      // A3 paper (297mm x 420mm)
screencraft.A5      // A5 paper (148mm x 210mm)
screencraft.Letter  // US Letter (8.5in x 11in)
screencraft.Legal   // US Legal (8.5in x 14in)
screencraft.Tabloid // Tabloid (11in x 17in)
```

### PDF Orientation

```go
screencraft.Portrait  // Portrait orientation
screencraft.Landscape // Landscape orientation
```

### Wait Until Events

```go
screencraft.WaitLoad            // Wait for 'load' event
screencraft.WaitDOMContentLoaded // Wait for 'DOMContentLoaded' event
screencraft.WaitNetworkIdle     // Wait for network idle (2 connections)
screencraft.WaitNetworkIdle0    // Wait for network idle (0 connections)
```

## License

MIT License - see [LICENSE](LICENSE) for details.
