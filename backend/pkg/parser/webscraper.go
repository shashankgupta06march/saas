package parser

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gocolly/colly/v2"
)

// ScrapeWebsite extracts text content from a website URL with specified depth
func ScrapeWebsite(websiteURL string, depth int) (string, error) {
	// Validate depth (0-5 for safety)
	if depth < 0 {
		depth = 0
	}
	if depth > 5 {
		depth = 5
	}


	// Parse base URL to get allowed domain
	parsedURL, err := url.Parse(websiteURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	var content strings.Builder
	var scrapeErr error

	c := colly.NewCollector(
		colly.AllowedDomains(parsedURL.Host),
		colly.MaxDepth(depth),
	)

	// Remove script and style tags
	c.OnHTML("script, style, noscript, iframe", func(e *colly.HTMLElement) {
		e.DOM.Remove()
	})

	// Extract text from common content areas
	c.OnHTML("body", func(e *colly.HTMLElement) {
		// Get text from paragraphs, headings, lists, etc.
		e.ForEach("p, h1, h2, h3, h4, h5, h6, li, article, section, div.content, div.main, main", func(_ int, el *colly.HTMLElement) {
			text := strings.TrimSpace(el.Text)
			if len(text) > 0 {
				content.WriteString(text)
				content.WriteString("\n")
			}
		})

		// Fallback: if no content found with specific tags, get all text
		if content.Len() == 0 {
			text := strings.TrimSpace(e.Text)
			content.WriteString(text)
		}
	})

	// Handle errors
	c.OnError(func(r *colly.Response, err error) {
		scrapeErr = fmt.Errorf("failed to scrape URL: %w", err)
	})

	// Visit the URL
	err = c.Visit(websiteURL)
	if err != nil {
		return "", fmt.Errorf("failed to visit URL: %w", err)
	}

	if scrapeErr != nil {
		return "", scrapeErr
	}

	result := content.String()
	if len(result) == 0 {
		return "", fmt.Errorf("no text content found on the website")
	}

	// Clean up the text
	result = cleanText(result)

	return result, nil
}

// cleanText removes excessive whitespace and empty lines
func cleanText(text string) string {
	lines := strings.Split(text, "\n")
	var cleaned []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			cleaned = append(cleaned, line)
		}
	}

	return strings.Join(cleaned, "\n")
}


