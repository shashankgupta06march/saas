package parser

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
)

// ScrapedPage holds the extracted content of a single crawled page.
type ScrapedPage struct {
	URL     string
	Title   string
	Content string
}

// nonHTMLExtensions are skipped during link discovery since they can't
// contain crawlable text content.
var nonHTMLExtensions = []string{
	".pdf", ".jpg", ".jpeg", ".png", ".gif", ".svg", ".ico", ".webp",
	".zip", ".rar", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx",
	".css", ".js", ".mp4", ".mp3", ".avi", ".mov", ".woff", ".woff2", ".ttf",
}

// ScrapeWebsite crawls a website starting at websiteURL, following internal
// links up to `depth` levels deep (0 = the given page only), and returns the
// extracted text of every page visited, capped at maxPages.
func ScrapeWebsite(websiteURL string, depth int, maxPages int) ([]ScrapedPage, error) {
	if depth < 0 {
		depth = 0
	}
	if depth > 5 {
		depth = 5
	}
	if maxPages <= 0 {
		maxPages = 30
	}

	parsedURL, err := url.Parse(websiteURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	allowedHosts := []string{parsedURL.Host}
	if strings.HasPrefix(parsedURL.Host, "www.") {
		allowedHosts = append(allowedHosts, strings.TrimPrefix(parsedURL.Host, "www."))
	} else {
		allowedHosts = append(allowedHosts, "www."+parsedURL.Host)
	}

	var (
		mu        sync.Mutex
		pages     []ScrapedPage
		visited   = map[string]bool{}
		pageCount int
		scrapeErr error
	)

	c := colly.NewCollector(
		colly.AllowedDomains(allowedHosts...),
		// Colly counts the starting page as depth 1, so +1 maps the UI's
		// "0 = this page only" onto Colly's "1 = this page only".
		colly.MaxDepth(depth+1),
		colly.Async(true),
		// Colly's default UA is non-browser-like and gets served stripped-down
		// content by some sites' anti-bot defenses. A real browser UA avoids that.
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		RandomDelay: 400 * time.Millisecond,
	})

	// Enforce the page cap; Colly has no built-in equivalent.
	c.OnRequest(func(r *colly.Request) {
		mu.Lock()
		defer mu.Unlock()
		if pageCount >= maxPages {
			r.Abort()
			return
		}
		pageCount++
	})

	// Strip elements that just repeat site-wide boilerplate (nav menus,
	// headers, footers) on every single page — without this, large multi-level
	// nav menus dominate the extracted text and crowd out each page's actual
	// unique content.
	c.OnHTML("script, style, noscript, iframe, nav, header, footer", func(e *colly.HTMLElement) {
		e.DOM.Remove()
	})

	c.OnHTML("html", func(e *colly.HTMLElement) {
		pageTitle := strings.TrimSpace(e.ChildText("head title"))

		var content strings.Builder
		e.ForEach("p, h1, h2, h3, h4, h5, h6, li, article, section, div.content, div.main, main", func(_ int, el *colly.HTMLElement) {
			text := strings.TrimSpace(el.Text)
			if len(text) > 0 {
				content.WriteString(text)
				content.WriteString("\n")
			}
		})
		if content.Len() == 0 {
			content.WriteString(strings.TrimSpace(e.Text))
		}

		text := cleanText(content.String())
		if len(text) == 0 {
			return
		}

		mu.Lock()
		pages = append(pages, ScrapedPage{
			URL:     e.Request.URL.String(),
			Title:   pageTitle,
			Content: text,
		})
		mu.Unlock()
	})

	// Discover and follow same-domain links — without this, MaxDepth has
	// nothing to apply to and only the starting page is ever visited.
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if link == "" || strings.HasPrefix(link, "#") ||
			strings.HasPrefix(link, "mailto:") || strings.HasPrefix(link, "tel:") ||
			strings.HasPrefix(link, "javascript:") {
			return
		}

		absolute := e.Request.AbsoluteURL(link)
		if absolute == "" {
			return
		}

		lower := strings.ToLower(absolute)
		for _, ext := range nonHTMLExtensions {
			if strings.HasSuffix(lower, ext) {
				return
			}
		}

		mu.Lock()
		alreadyVisited := visited[absolute]
		visited[absolute] = true
		atLimit := pageCount >= maxPages
		mu.Unlock()

		if alreadyVisited || atLimit {
			return
		}

		e.Request.Visit(absolute)
	})

	c.OnError(func(r *colly.Response, err error) {
		// Individual page failures (404s, timeouts) are expected on real
		// sites; only abort the whole crawl if the very first page fails.
		if r.Request.URL.String() == websiteURL {
			mu.Lock()
			scrapeErr = fmt.Errorf("failed to scrape URL: %w", err)
			mu.Unlock()
		}
	})

	if err := c.Visit(websiteURL); err != nil {
		return nil, fmt.Errorf("failed to visit URL: %w", err)
	}
	c.Wait()

	if scrapeErr != nil {
		return nil, scrapeErr
	}

	if len(pages) == 0 {
		return nil, fmt.Errorf("no text content found on the website")
	}

	return pages, nil
}

// cleanText removes excessive whitespace and empty lines.
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
