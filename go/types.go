// Package screencraft provides a Go SDK for the ScreenCraft API.
// It enables screenshot capture and PDF generation from web pages.
package screencraft

import "time"

// Format represents the output format for screenshots.
type Format string

const (
	// FormatPNG represents PNG image format.
	FormatPNG Format = "png"
	// FormatJPEG represents JPEG image format.
	FormatJPEG Format = "jpeg"
	// FormatWebP represents WebP image format.
	FormatWebP Format = "webp"
)

// PDFFormat represents the paper format for PDF generation.
type PDFFormat string

const (
	// A4 represents A4 paper format (210mm x 297mm).
	A4 PDFFormat = "A4"
	// A3 represents A3 paper format (297mm x 420mm).
	A3 PDFFormat = "A3"
	// A5 represents A5 paper format (148mm x 210mm).
	A5 PDFFormat = "A5"
	// Letter represents US Letter paper format (8.5in x 11in).
	Letter PDFFormat = "Letter"
	// Legal represents US Legal paper format (8.5in x 14in).
	Legal PDFFormat = "Legal"
	// Tabloid represents Tabloid paper format (11in x 17in).
	Tabloid PDFFormat = "Tabloid"
)

// PDFOrientation represents the page orientation for PDF generation.
type PDFOrientation string

const (
	// Portrait represents portrait orientation.
	Portrait PDFOrientation = "portrait"
	// Landscape represents landscape orientation.
	Landscape PDFOrientation = "landscape"
)

// Viewport represents the browser viewport dimensions.
type Viewport struct {
	// Width of the viewport in pixels.
	Width int `json:"width,omitempty"`
	// Height of the viewport in pixels.
	Height int `json:"height,omitempty"`
}

// ScrollPosition represents the scroll position before capture.
type ScrollPosition struct {
	// X is the horizontal scroll position in pixels.
	X int `json:"x,omitempty"`
	// Y is the vertical scroll position in pixels.
	Y int `json:"y,omitempty"`
}

// Clip represents a rectangular region to clip from the page.
type Clip struct {
	// X is the horizontal offset from the left edge.
	X int `json:"x"`
	// Y is the vertical offset from the top edge.
	Y int `json:"y"`
	// Width is the width of the clip region.
	Width int `json:"width"`
	// Height is the height of the clip region.
	Height int `json:"height"`
}

// Cookie represents a browser cookie to set before navigation.
type Cookie struct {
	// Name is the cookie name.
	Name string `json:"name"`
	// Value is the cookie value.
	Value string `json:"value"`
	// Domain is the cookie domain.
	Domain string `json:"domain,omitempty"`
	// Path is the cookie path.
	Path string `json:"path,omitempty"`
	// Secure indicates if the cookie is secure-only.
	Secure bool `json:"secure,omitempty"`
	// HTTPOnly indicates if the cookie is HTTP-only.
	HTTPOnly bool `json:"httpOnly,omitempty"`
	// SameSite is the cookie SameSite attribute.
	SameSite string `json:"sameSite,omitempty"`
	// Expires is the cookie expiration time.
	Expires *time.Time `json:"expires,omitempty"`
}

// Header represents a custom HTTP header.
type Header struct {
	// Name is the header name.
	Name string `json:"name"`
	// Value is the header value.
	Value string `json:"value"`
}

// WaitUntil represents the page load event to wait for.
type WaitUntil string

const (
	// WaitLoad waits for the 'load' event.
	WaitLoad WaitUntil = "load"
	// WaitDOMContentLoaded waits for the 'DOMContentLoaded' event.
	WaitDOMContentLoaded WaitUntil = "domcontentloaded"
	// WaitNetworkIdle waits for network to be idle (no more than 2 connections for 500ms).
	WaitNetworkIdle WaitUntil = "networkidle"
	// WaitNetworkIdle0 waits for network to be idle (no connections for 500ms).
	WaitNetworkIdle0 WaitUntil = "networkidle0"
)

// WebhookConfig represents webhook configuration for async operations.
type WebhookConfig struct {
	// URL is the webhook endpoint to call when the operation completes.
	URL string `json:"url"`
	// Headers are custom headers to include in the webhook request.
	Headers map[string]string `json:"headers,omitempty"`
	// Secret is an optional secret for webhook signature verification.
	Secret string `json:"secret,omitempty"`
}

// ScreenshotOptions represents options for taking a screenshot.
type ScreenshotOptions struct {
	// URL is the target URL to capture.
	URL string `json:"url"`
	// Format is the output image format (png, jpeg, webp).
	Format Format `json:"format,omitempty"`
	// Quality is the image quality (0-100), applicable for JPEG and WebP.
	Quality int `json:"quality,omitempty"`
	// FullPage captures the full scrollable page if true.
	FullPage bool `json:"fullPage,omitempty"`
	// Viewport sets the browser viewport dimensions.
	Viewport *Viewport `json:"viewport,omitempty"`
	// ScrollPosition sets the scroll position before capture.
	ScrollPosition *ScrollPosition `json:"scrollPosition,omitempty"`
	// Clip defines a rectangular region to clip.
	Clip *Clip `json:"clip,omitempty"`
	// AcceptCookies automatically accepts cookie consent banners.
	AcceptCookies bool `json:"acceptCookies,omitempty"`
	// Delay is the time to wait after page load before capture (in milliseconds).
	Delay int `json:"delay,omitempty"`
	// WaitUntil specifies the page load event to wait for.
	WaitUntil WaitUntil `json:"waitUntil,omitempty"`
	// WaitForSelector waits for a specific CSS selector to appear.
	WaitForSelector string `json:"waitForSelector,omitempty"`
	// WaitForTimeout is an additional wait time in milliseconds.
	WaitForTimeout int `json:"waitForTimeout,omitempty"`
	// Cookies are cookies to set before navigation.
	Cookies []Cookie `json:"cookies,omitempty"`
	// Headers are custom HTTP headers to send.
	Headers []Header `json:"headers,omitempty"`
	// UserAgent sets a custom user agent string.
	UserAgent string `json:"userAgent,omitempty"`
	// DeviceScaleFactor sets the device scale factor (DPR).
	DeviceScaleFactor float64 `json:"deviceScaleFactor,omitempty"`
	// IsMobile emulates a mobile device.
	IsMobile bool `json:"isMobile,omitempty"`
	// HasTouch enables touch event emulation.
	HasTouch bool `json:"hasTouch,omitempty"`
	// IsLandscape sets the viewport to landscape orientation.
	IsLandscape bool `json:"isLandscape,omitempty"`
	// DarkMode enables dark mode emulation.
	DarkMode bool `json:"darkMode,omitempty"`
	// BlockAds blocks advertisements.
	BlockAds bool `json:"blockAds,omitempty"`
	// BlockTrackers blocks tracking scripts.
	BlockTrackers bool `json:"blockTrackers,omitempty"`
	// BypassCSP bypasses Content Security Policy.
	BypassCSP bool `json:"bypassCSP,omitempty"`
	// JavaScript enables or disables JavaScript (enabled by default).
	JavaScript *bool `json:"javascript,omitempty"`
	// Webhook configures async webhook delivery.
	Webhook *WebhookConfig `json:"webhook,omitempty"`
}

// PDFOptions represents options for generating a PDF.
type PDFOptions struct {
	// URL is the target URL to convert to PDF.
	URL string `json:"url"`
	// Format is the paper format (A4, Letter, etc.).
	Format PDFFormat `json:"format,omitempty"`
	// Orientation is the page orientation (portrait or landscape).
	Orientation PDFOrientation `json:"orientation,omitempty"`
	// Width is a custom page width (e.g., "8.5in", "210mm").
	Width string `json:"width,omitempty"`
	// Height is a custom page height (e.g., "11in", "297mm").
	Height string `json:"height,omitempty"`
	// Scale is the PDF scale factor (0.1 to 2.0).
	Scale float64 `json:"scale,omitempty"`
	// DisplayHeaderFooter enables the header and footer.
	DisplayHeaderFooter bool `json:"displayHeaderFooter,omitempty"`
	// HeaderTemplate is the HTML template for the header.
	HeaderTemplate string `json:"headerTemplate,omitempty"`
	// FooterTemplate is the HTML template for the footer.
	FooterTemplate string `json:"footerTemplate,omitempty"`
	// PrintBackground prints background graphics.
	PrintBackground bool `json:"printBackground,omitempty"`
	// PreferCSSPageSize uses CSS @page size if specified.
	PreferCSSPageSize bool `json:"preferCSSPageSize,omitempty"`
	// PageRanges specifies which pages to include (e.g., "1-5, 8, 11-13").
	PageRanges string `json:"pageRanges,omitempty"`
	// Margin sets the page margins.
	Margin *PDFMargin `json:"margin,omitempty"`
	// Viewport sets the browser viewport dimensions.
	Viewport *Viewport `json:"viewport,omitempty"`
	// AcceptCookies automatically accepts cookie consent banners.
	AcceptCookies bool `json:"acceptCookies,omitempty"`
	// Delay is the time to wait after page load before PDF generation (in milliseconds).
	Delay int `json:"delay,omitempty"`
	// WaitUntil specifies the page load event to wait for.
	WaitUntil WaitUntil `json:"waitUntil,omitempty"`
	// WaitForSelector waits for a specific CSS selector to appear.
	WaitForSelector string `json:"waitForSelector,omitempty"`
	// WaitForTimeout is an additional wait time in milliseconds.
	WaitForTimeout int `json:"waitForTimeout,omitempty"`
	// Cookies are cookies to set before navigation.
	Cookies []Cookie `json:"cookies,omitempty"`
	// Headers are custom HTTP headers to send.
	Headers []Header `json:"headers,omitempty"`
	// UserAgent sets a custom user agent string.
	UserAgent string `json:"userAgent,omitempty"`
	// DarkMode enables dark mode emulation.
	DarkMode bool `json:"darkMode,omitempty"`
	// BlockAds blocks advertisements.
	BlockAds bool `json:"blockAds,omitempty"`
	// BlockTrackers blocks tracking scripts.
	BlockTrackers bool `json:"blockTrackers,omitempty"`
	// BypassCSP bypasses Content Security Policy.
	BypassCSP bool `json:"bypassCSP,omitempty"`
	// JavaScript enables or disables JavaScript (enabled by default).
	JavaScript *bool `json:"javascript,omitempty"`
	// Webhook configures async webhook delivery.
	Webhook *WebhookConfig `json:"webhook,omitempty"`
}

// PDFMargin represents page margins for PDF generation.
type PDFMargin struct {
	// Top margin (e.g., "1in", "25mm", "100px").
	Top string `json:"top,omitempty"`
	// Right margin.
	Right string `json:"right,omitempty"`
	// Bottom margin.
	Bottom string `json:"bottom,omitempty"`
	// Left margin.
	Left string `json:"left,omitempty"`
}

// ScreenshotResult represents the result of a screenshot operation.
type ScreenshotResult struct {
	// Data contains the screenshot image data.
	Data []byte
	// ContentType is the MIME type of the image.
	ContentType string
	// URL is the captured URL.
	URL string
	// Width is the image width in pixels.
	Width int
	// Height is the image height in pixels.
	Height int
	// JobID is the async job ID when using webhooks.
	JobID string
}

// PDFResult represents the result of a PDF generation operation.
type PDFResult struct {
	// Data contains the PDF data.
	Data []byte
	// ContentType is the MIME type (application/pdf).
	ContentType string
	// URL is the captured URL.
	URL string
	// Pages is the number of pages in the PDF.
	Pages int
	// JobID is the async job ID when using webhooks.
	JobID string
}

// APIResponse represents a generic API response.
type APIResponse struct {
	// Success indicates if the operation was successful.
	Success bool `json:"success"`
	// Message is an optional message from the API.
	Message string `json:"message,omitempty"`
	// JobID is the async job ID for webhook operations.
	JobID string `json:"jobId,omitempty"`
	// Error contains error details if success is false.
	Error *APIErrorDetails `json:"error,omitempty"`
}

// APIErrorDetails contains detailed error information from the API.
type APIErrorDetails struct {
	// Code is the error code.
	Code string `json:"code"`
	// Message is the error message.
	Message string `json:"message"`
	// Details contains additional error details.
	Details map[string]interface{} `json:"details,omitempty"`
}

// RateLimitInfo contains rate limit information from API responses.
type RateLimitInfo struct {
	// Limit is the maximum number of requests allowed.
	Limit int
	// Remaining is the number of requests remaining.
	Remaining int
	// Reset is when the rate limit resets.
	Reset time.Time
}
