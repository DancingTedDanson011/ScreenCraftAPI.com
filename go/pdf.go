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
	pdfEndpoint = "/pdf"
)

// PDF generates a PDF from the specified URL.
//
// The function sends a request to the ScreenCraft API to generate a PDF
// of the given URL with the specified options. It returns the PDF data
// or an error if the operation fails.
//
// Example:
//
//	result, err := client.PDF(ctx, &screencraft.PDFOptions{
//	    URL:    "https://example.com",
//	    Format: screencraft.A4,
//	    PrintBackground: true,
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	os.WriteFile("document.pdf", result.Data, 0644)
func (c *Client) PDF(ctx context.Context, opts *PDFOptions) (*PDFResult, error) {
	if err := ValidatePDFOptions(opts); err != nil {
		return nil, err
	}

	// Build request body
	reqBody := c.buildPDFRequest(opts)

	resp, err := c.doRequest(ctx, http.MethodPost, pdfEndpoint, reqBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return c.parsePDFResponse(resp, opts)
}

// PDFAsync generates a PDF asynchronously using webhooks.
//
// The function sends a request to the ScreenCraft API to generate a PDF
// and delivers the result to the specified webhook URL. It returns the job ID
// for tracking the operation.
//
// Example:
//
//	jobID, err := client.PDFAsync(ctx, &screencraft.PDFOptions{
//	    URL:    "https://example.com",
//	    Format: screencraft.A4,
//	    Webhook: &screencraft.WebhookConfig{
//	        URL: "https://yoursite.com/webhook",
//	    },
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Job ID: %s\n", jobID)
func (c *Client) PDFAsync(ctx context.Context, opts *PDFOptions) (string, error) {
	if err := ValidatePDFOptions(opts); err != nil {
		return "", err
	}

	if opts.Webhook == nil || opts.Webhook.URL == "" {
		return "", NewValidationError("webhook.url", "webhook URL is required for async operations", "required").Error
	}

	// Build request body
	reqBody := c.buildPDFRequest(opts)

	resp, err := c.doRequest(ctx, http.MethodPost, pdfEndpoint, reqBody)
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

// buildPDFRequest builds the API request body for PDF generation.
func (c *Client) buildPDFRequest(opts *PDFOptions) map[string]interface{} {
	req := map[string]interface{}{
		"url": opts.URL,
	}

	if opts.Format != "" {
		req["format"] = opts.Format
	}

	if opts.Orientation != "" {
		req["orientation"] = opts.Orientation
	}

	if opts.Width != "" {
		req["width"] = opts.Width
	}

	if opts.Height != "" {
		req["height"] = opts.Height
	}

	if opts.Scale != 0 {
		req["scale"] = opts.Scale
	}

	if opts.DisplayHeaderFooter {
		req["displayHeaderFooter"] = true
	}

	if opts.HeaderTemplate != "" {
		req["headerTemplate"] = opts.HeaderTemplate
	}

	if opts.FooterTemplate != "" {
		req["footerTemplate"] = opts.FooterTemplate
	}

	if opts.PrintBackground {
		req["printBackground"] = true
	}

	if opts.PreferCSSPageSize {
		req["preferCSSPageSize"] = true
	}

	if opts.PageRanges != "" {
		req["pageRanges"] = opts.PageRanges
	}

	if opts.Margin != nil {
		margin := map[string]interface{}{}
		if opts.Margin.Top != "" {
			margin["top"] = opts.Margin.Top
		}
		if opts.Margin.Right != "" {
			margin["right"] = opts.Margin.Right
		}
		if opts.Margin.Bottom != "" {
			margin["bottom"] = opts.Margin.Bottom
		}
		if opts.Margin.Left != "" {
			margin["left"] = opts.Margin.Left
		}
		if len(margin) > 0 {
			req["margin"] = margin
		}
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

// parsePDFResponse parses the PDF response from the API.
func (c *Client) parsePDFResponse(resp *http.Response, opts *PDFOptions) (*PDFResult, error) {
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
		return &PDFResult{
			URL:   opts.URL,
			JobID: apiResp.JobID,
		}, nil
	}

	// Binary PDF response
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("screencraft: failed to read PDF data: %w", err)
	}

	result := &PDFResult{
		Data:        data,
		ContentType: contentType,
		URL:         opts.URL,
	}

	// Parse page count header if available
	if p := resp.Header.Get("X-PDF-Pages"); p != "" {
		if pages, err := strconv.Atoi(p); err == nil {
			result.Pages = pages
		}
	}

	return result, nil
}

// PDFURL generates a PDF with minimal options.
//
// This is a convenience method for simple PDF generation.
//
// Example:
//
//	result, err := client.PDFURL(ctx, "https://example.com")
func (c *Client) PDFURL(ctx context.Context, url string) (*PDFResult, error) {
	return c.PDF(ctx, &PDFOptions{
		URL:    url,
		Format: A4,
	})
}

// PDFA4 generates an A4 PDF.
//
// Example:
//
//	result, err := client.PDFA4(ctx, "https://example.com")
func (c *Client) PDFA4(ctx context.Context, url string) (*PDFResult, error) {
	return c.PDF(ctx, &PDFOptions{
		URL:             url,
		Format:          A4,
		PrintBackground: true,
	})
}

// PDFLetter generates a US Letter PDF.
//
// Example:
//
//	result, err := client.PDFLetter(ctx, "https://example.com")
func (c *Client) PDFLetter(ctx context.Context, url string) (*PDFResult, error) {
	return c.PDF(ctx, &PDFOptions{
		URL:             url,
		Format:          Letter,
		PrintBackground: true,
	})
}

// PDFLandscape generates a PDF in landscape orientation.
//
// Example:
//
//	result, err := client.PDFLandscape(ctx, "https://example.com", screencraft.A4)
func (c *Client) PDFLandscape(ctx context.Context, url string, format PDFFormat) (*PDFResult, error) {
	return c.PDF(ctx, &PDFOptions{
		URL:             url,
		Format:          format,
		Orientation:     Landscape,
		PrintBackground: true,
	})
}

// PDFWithMargins generates a PDF with custom margins.
//
// Example:
//
//	result, err := client.PDFWithMargins(ctx, "https://example.com", &screencraft.PDFMargin{
//	    Top:    "1in",
//	    Right:  "1in",
//	    Bottom: "1in",
//	    Left:   "1in",
//	})
func (c *Client) PDFWithMargins(ctx context.Context, url string, margins *PDFMargin) (*PDFResult, error) {
	return c.PDF(ctx, &PDFOptions{
		URL:             url,
		Format:          A4,
		Margin:          margins,
		PrintBackground: true,
	})
}

// PDFWithHeaderFooter generates a PDF with custom header and footer.
//
// Example:
//
//	result, err := client.PDFWithHeaderFooter(ctx, "https://example.com",
//	    "<div style='font-size:10px;'>Header</div>",
//	    "<div style='font-size:10px;'>Page <span class='pageNumber'></span></div>",
//	)
func (c *Client) PDFWithHeaderFooter(ctx context.Context, url, headerHTML, footerHTML string) (*PDFResult, error) {
	return c.PDF(ctx, &PDFOptions{
		URL:                 url,
		Format:              A4,
		DisplayHeaderFooter: true,
		HeaderTemplate:      headerHTML,
		FooterTemplate:      footerHTML,
		PrintBackground:     true,
		Margin: &PDFMargin{
			Top:    "100px",
			Bottom: "100px",
		},
	})
}

// PDFPageRange generates a PDF with specific page ranges.
//
// Example:
//
//	result, err := client.PDFPageRange(ctx, "https://example.com", "1-5, 8")
func (c *Client) PDFPageRange(ctx context.Context, url, pageRanges string) (*PDFResult, error) {
	return c.PDF(ctx, &PDFOptions{
		URL:             url,
		Format:          A4,
		PageRanges:      pageRanges,
		PrintBackground: true,
	})
}

// PDFWithCookieConsent generates a PDF and auto-accepts cookie banners.
//
// Example:
//
//	result, err := client.PDFWithCookieConsent(ctx, "https://example.com")
func (c *Client) PDFWithCookieConsent(ctx context.Context, url string) (*PDFResult, error) {
	return c.PDF(ctx, &PDFOptions{
		URL:             url,
		Format:          A4,
		AcceptCookies:   true,
		PrintBackground: true,
	})
}
