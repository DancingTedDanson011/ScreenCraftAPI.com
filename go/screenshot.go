package screencraft

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const (
	screenshotEndpoint = "/screenshot"
)

// Screenshot captures a screenshot of the specified URL.
//
// The function sends a request to the ScreenCraft API to capture a screenshot
// of the given URL with the specified options. It returns the screenshot data
// or an error if the operation fails.
//
// Example:
//
//	result, err := client.Screenshot(ctx, &screencraft.ScreenshotOptions{
//	    URL:      "https://example.com",
//	    Format:   screencraft.FormatPNG,
//	    FullPage: true,
//	    Viewport: &screencraft.Viewport{Width: 1920, Height: 1080},
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	os.WriteFile("screenshot.png", result.Data, 0644)
func (c *Client) Screenshot(ctx context.Context, opts *ScreenshotOptions) (*ScreenshotResult, error) {
	if err := ValidateScreenshotOptions(opts); err != nil {
		return nil, err
	}

	// Build request body
	reqBody := c.buildScreenshotRequest(opts)

	resp, err := c.doRequest(ctx, http.MethodPost, screenshotEndpoint, reqBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return c.parseScreenshotResponse(resp, opts)
}

// ScreenshotAsync captures a screenshot asynchronously using webhooks.
//
// The function sends a request to the ScreenCraft API to capture a screenshot
// and delivers the result to the specified webhook URL. It returns the job ID
// for tracking the operation.
//
// Example:
//
//	jobID, err := client.ScreenshotAsync(ctx, &screencraft.ScreenshotOptions{
//	    URL:    "https://example.com",
//	    Format: screencraft.FormatPNG,
//	    Webhook: &screencraft.WebhookConfig{
//	        URL: "https://yoursite.com/webhook",
//	    },
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Job ID: %s\n", jobID)
func (c *Client) ScreenshotAsync(ctx context.Context, opts *ScreenshotOptions) (string, error) {
	if err := ValidateScreenshotOptions(opts); err != nil {
		return "", err
	}

	if opts.Webhook == nil || opts.Webhook.URL == "" {
		return "", NewValidationError("webhook.url", "webhook URL is required for async operations", "required").Error
	}

	// Build request body
	reqBody := c.buildScreenshotRequest(opts)

	resp, err := c.doRequest(ctx, http.MethodPost, screenshotEndpoint, reqBody)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse async response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("screencraft: failed to read response: %w", err)
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", fmt.Errorf("screencraft: failed to parse response: %w", err)
	}

	if !apiResp.Success {
		return "", &Error{
			StatusCode: resp.StatusCode,
			Message:    apiResp.Message,
		}
	}

	return apiResp.JobID, nil
}

// buildScreenshotRequest builds the API request body for a screenshot.
func (c *Client) buildScreenshotRequest(opts *ScreenshotOptions) map[string]interface{} {
	req := map[string]interface{}{
		"url": opts.URL,
	}

	if opts.Format != "" {
		req["format"] = opts.Format
	}

	if opts.Quality > 0 {
		req["quality"] = opts.Quality
	}

	if opts.FullPage {
		req["fullPage"] = true
	}

	if opts.Viewport != nil {
		viewport := map[string]interface{}{}
		if opts.Viewport.Width > 0 {
			viewport["width"] = opts.Viewport.Width
		}
		if opts.Viewport.Height > 0 {
			viewport["height"] = opts.Viewport.Height
		}
		if len(viewport) > 0 {
			req["viewport"] = viewport
		}
	}

	if opts.ScrollPosition != nil {
		scrollPos := map[string]interface{}{
			"x": opts.ScrollPosition.X,
			"y": opts.ScrollPosition.Y,
		}
		req["scrollPosition"] = scrollPos
	}

	if opts.Clip != nil {
		req["clip"] = map[string]interface{}{
			"x":      opts.Clip.X,
			"y":      opts.Clip.Y,
			"width":  opts.Clip.Width,
			"height": opts.Clip.Height,
		}
	}

	if opts.AcceptCookies {
		req["acceptCookies"] = true
	}

	if opts.Delay > 0 {
		req["delay"] = opts.Delay
	}

	if opts.WaitUntil != "" {
		req["waitUntil"] = opts.WaitUntil
	}

	if opts.WaitForSelector != "" {
		req["waitForSelector"] = opts.WaitForSelector
	}

	if opts.WaitForTimeout > 0 {
		req["waitForTimeout"] = opts.WaitForTimeout
	}

	if len(opts.Cookies) > 0 {
		req["cookies"] = opts.Cookies
	}

	if len(opts.Headers) > 0 {
		req["headers"] = opts.Headers
	}

	if opts.UserAgent != "" {
		req["userAgent"] = opts.UserAgent
	}

	if opts.DeviceScaleFactor > 0 {
		req["deviceScaleFactor"] = opts.DeviceScaleFactor
	}

	if opts.IsMobile {
		req["isMobile"] = true
	}

	if opts.HasTouch {
		req["hasTouch"] = true
	}

	if opts.IsLandscape {
		req["isLandscape"] = true
	}

	if opts.DarkMode {
		req["darkMode"] = true
	}

	if opts.BlockAds {
		req["blockAds"] = true
	}

	if opts.BlockTrackers {
		req["blockTrackers"] = true
	}

	if opts.BypassCSP {
		req["bypassCSP"] = true
	}

	if opts.JavaScript != nil {
		req["javascript"] = *opts.JavaScript
	}

	if opts.Webhook != nil {
		webhook := map[string]interface{}{
			"url": opts.Webhook.URL,
		}
		if len(opts.Webhook.Headers) > 0 {
			webhook["headers"] = opts.Webhook.Headers
		}
		if opts.Webhook.Secret != "" {
			webhook["secret"] = opts.Webhook.Secret
		}
		req["webhook"] = webhook
	}

	return req
}

// parseScreenshotResponse parses the screenshot response from the API.
func (c *Client) parseScreenshotResponse(resp *http.Response, opts *ScreenshotOptions) (*ScreenshotResult, error) {
	contentType := resp.Header.Get("Content-Type")

	// Check if this is a JSON response (async or error)
	if contentType == "application/json" {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("screencraft: failed to read response: %w", err)
		}

		var apiResp APIResponse
		if err := json.Unmarshal(body, &apiResp); err != nil {
			return nil, fmt.Errorf("screencraft: failed to parse response: %w", err)
		}

		if !apiResp.Success {
			return nil, &Error{
				StatusCode: resp.StatusCode,
				Message:    apiResp.Message,
			}
		}

		// Async response
		return &ScreenshotResult{
			URL:   opts.URL,
			JobID: apiResp.JobID,
		}, nil
	}

	// Binary image response
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("screencraft: failed to read image data: %w", err)
	}

	result := &ScreenshotResult{
		Data:        data,
		ContentType: contentType,
		URL:         opts.URL,
	}

	// Parse dimension headers if available
	if w := resp.Header.Get("X-Image-Width"); w != "" {
		if width, err := strconv.Atoi(w); err == nil {
			result.Width = width
		}
	}

	if h := resp.Header.Get("X-Image-Height"); h != "" {
		if height, err := strconv.Atoi(h); err == nil {
			result.Height = height
		}
	}

	return result, nil
}

// ScreenshotURL captures a screenshot with minimal options.
//
// This is a convenience method for simple screenshot captures.
//
// Example:
//
//	result, err := client.ScreenshotURL(ctx, "https://example.com")
func (c *Client) ScreenshotURL(ctx context.Context, url string) (*ScreenshotResult, error) {
	return c.Screenshot(ctx, &ScreenshotOptions{
		URL:    url,
		Format: FormatPNG,
	})
}

// ScreenshotFullPage captures a full-page screenshot.
//
// This is a convenience method for full-page screenshot captures.
//
// Example:
//
//	result, err := client.ScreenshotFullPage(ctx, "https://example.com", screencraft.FormatPNG)
func (c *Client) ScreenshotFullPage(ctx context.Context, url string, format Format) (*ScreenshotResult, error) {
	return c.Screenshot(ctx, &ScreenshotOptions{
		URL:      url,
		Format:   format,
		FullPage: true,
	})
}

// ScreenshotMobile captures a screenshot with mobile emulation.
//
// This method sets appropriate viewport and mobile device settings.
//
// Example:
//
//	result, err := client.ScreenshotMobile(ctx, "https://example.com")
func (c *Client) ScreenshotMobile(ctx context.Context, url string) (*ScreenshotResult, error) {
	return c.Screenshot(ctx, &ScreenshotOptions{
		URL:    url,
		Format: FormatPNG,
		Viewport: &Viewport{
			Width:  375,
			Height: 812,
		},
		IsMobile:          true,
		HasTouch:          true,
		DeviceScaleFactor: 3,
	})
}

// ScreenshotDesktop captures a screenshot with desktop viewport.
//
// This method sets a standard desktop viewport size.
//
// Example:
//
//	result, err := client.ScreenshotDesktop(ctx, "https://example.com")
func (c *Client) ScreenshotDesktop(ctx context.Context, url string) (*ScreenshotResult, error) {
	return c.Screenshot(ctx, &ScreenshotOptions{
		URL:    url,
		Format: FormatPNG,
		Viewport: &Viewport{
			Width:  1920,
			Height: 1080,
		},
	})
}

// ScreenshotWithDelay captures a screenshot after waiting for a specified delay.
//
// This is useful for pages with animations or dynamic content.
//
// Example:
//
//	result, err := client.ScreenshotWithDelay(ctx, "https://example.com", 2000)
func (c *Client) ScreenshotWithDelay(ctx context.Context, url string, delayMs int) (*ScreenshotResult, error) {
	return c.Screenshot(ctx, &ScreenshotOptions{
		URL:    url,
		Format: FormatPNG,
		Delay:  delayMs,
	})
}

// ScreenshotWithCookieConsent captures a screenshot and auto-accepts cookie banners.
//
// Example:
//
//	result, err := client.ScreenshotWithCookieConsent(ctx, "https://example.com")
func (c *Client) ScreenshotWithCookieConsent(ctx context.Context, url string) (*ScreenshotResult, error) {
	return c.Screenshot(ctx, &ScreenshotOptions{
		URL:           url,
		Format:        FormatPNG,
		AcceptCookies: true,
	})
}
