package parser

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
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

// browserUA is a real-browser User-Agent. Colly's default UA is non-browser-like
// and gets served stripped-down content by some sites' anti-bot defenses.
const browserUA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"

// maxTransientRetries is how many times a page that fails transiently (rate
// limiting, gateway/timeout errors) is retried before being counted as failed.
const maxTransientRetries = 3

// nonHTMLExtensions are skipped during link discovery since they can't
// contain crawlable text content.
var nonHTMLExtensions = []string{
	".pdf", ".jpg", ".jpeg", ".png", ".gif", ".svg", ".ico", ".webp",
	".zip", ".rar", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx",
	".css", ".js", ".mp4", ".mp3", ".avi", ".mov", ".woff", ".woff2", ".ttf",
}

// ScrapeWebsite crawls a website starting at websiteURL, following internal
// links up to `depth` levels deep (0 = the given page only), and returns the
// extracted text of every page visited (capped at maxPages) plus a count of
// pages that failed to load even after retries (e.g. due to rate limiting).
func ScrapeWebsite(websiteURL string, depth int, maxPages int) ([]ScrapedPage, int, error) {
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
		return nil, 0, fmt.Errorf("invalid URL: %w", err)
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
		retries   = map[string]int{}
		pageCount int
		failed    int
		scrapeErr error
	)

	c := colly.NewCollector(
		colly.AllowedDomains(allowedHosts...),
		// Colly counts the starting page as depth 1, so +1 maps the UI's
		// "0 = this page only" onto Colly's "1 = this page only".
		colly.MaxDepth(depth+1),
		colly.Async(true),
		colly.UserAgent(browserUA),
	)
	c.SetRequestTimeout(30 * time.Second)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		RandomDelay: 400 * time.Millisecond,
	})

	// queueURL normalizes a candidate link and queues it for crawling if it's
	// a new, in-scope, HTML page and we're still under the page budget.
	queueURL := func(req *colly.Request, raw string) {
		if raw == "" || strings.HasPrefix(raw, "#") ||
			strings.HasPrefix(raw, "mailto:") || strings.HasPrefix(raw, "tel:") ||
			strings.HasPrefix(raw, "javascript:") {
			return
		}

		absolute := req.AbsoluteURL(raw)
		if absolute == "" {
			return
		}

		norm := NormalizeURL(absolute)
		lower := strings.ToLower(norm)
		for _, ext := range nonHTMLExtensions {
			if strings.HasSuffix(lower, ext) {
				return
			}
		}

		mu.Lock()
		alreadyVisited := visited[norm]
		visited[norm] = true
		atLimit := pageCount >= maxPages
		mu.Unlock()

		if alreadyVisited || atLimit {
			return
		}

		req.Visit(norm)
	}

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

	// A single html handler so ordering is explicit and guaranteed:
	//   1. discover links  ->  2. strip boilerplate  ->  3. extract text
	//
	// Discovery MUST happen before stripping: most internal links live inside
	// the <nav>/<header>/<footer> (top menu, quick links, footer). Removing
	// those elements first — as the previous version did via a separately
	// registered handler — meant those links were never followed, so whole
	// sections of the site were silently missed.
	c.OnHTML("html", func(e *colly.HTMLElement) {
		// 1. Discover and queue links from the full, untouched DOM.
		e.ForEach("a[href]", func(_ int, el *colly.HTMLElement) {
			queueURL(el.Request, el.Attr("href"))
		})

		// 2. Strip elements that just repeat site-wide boilerplate so they
		// don't crowd out each page's actual unique content.
		e.DOM.Find("script, style, noscript, iframe, nav, header, footer").Remove()

		// 3. Extract title and main textual content (incl. table cells, which
		// often hold structured data like fee tables).
		pageTitle := strings.TrimSpace(e.ChildText("head title"))

		var content strings.Builder
		e.ForEach("p, h1, h2, h3, h4, h5, h6, li, article, section, div.content, div.main, main, table, td, th", func(_ int, el *colly.HTMLElement) {
			text := strings.TrimSpace(el.Text)
			if len(text) > 0 {
				content.WriteString(text)
				content.WriteString("\n")
			}
		})
		if content.Len() == 0 {
			content.WriteString(strings.TrimSpace(e.Text))
		}

		// Capture in-content links as "label: url" lines. Plain text extraction
		// throws away hrefs, so resources like "Download Brochure", application
		// forms, prospectuses and PDF documents would otherwise be invisible to
		// the chatbot (it would know the link exists by name but not its URL).
		// Boilerplate (nav/header/footer) was already removed above, so these
		// are real content links, not site-wide menu chrome.
		var linkLines []string
		seenLinks := map[string]bool{}
		e.ForEach("a[href]", func(_ int, el *colly.HTMLElement) {
			href := strings.TrimSpace(el.Attr("href"))
			if href == "" || strings.HasPrefix(href, "#") ||
				strings.HasPrefix(href, "mailto:") || strings.HasPrefix(href, "tel:") ||
				strings.HasPrefix(href, "javascript:") {
				return
			}
			abs := el.Request.AbsoluteURL(href)
			if !strings.HasPrefix(abs, "http") {
				return
			}
			label := strings.Join(strings.Fields(el.Text), " ")
			if label == "" {
				return
			}
			key := label + "|" + abs
			if seenLinks[key] {
				return
			}
			seenLinks[key] = true
			if len(linkLines) < 80 {
				linkLines = append(linkLines, label+": "+abs)
			}
		})
		if len(linkLines) > 0 {
			content.WriteString("\nLinks:\n")
			content.WriteString(strings.Join(linkLines, "\n"))
			content.WriteString("\n")
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

	rootNorm := NormalizeURL(websiteURL)

	c.OnError(func(r *colly.Response, err error) {
		urlStr := r.Request.URL.String()
		status := r.StatusCode

		// Retry transient failures (rate limiting / gateway errors / timeouts)
		// with a simple linear backoff before giving up on the page.
		if status == 429 || status == 502 || status == 503 || status == 504 || status == 0 {
			mu.Lock()
			n := retries[urlStr]
			canRetry := n < maxTransientRetries
			if canRetry {
				retries[urlStr] = n + 1
				// Offset the pageCount++ the retried request will trigger so
				// retries don't eat into the page budget.
				if pageCount > 0 {
					pageCount--
				}
			}
			mu.Unlock()

			if canRetry {
				time.Sleep(time.Duration(n+1) * time.Second)
				_ = r.Request.Retry()
				return
			}
		}

		// Permanent failure. Only abort the whole crawl if the very first page
		// fails; otherwise just tally it so the caller can report partial runs.
		if urlStr == websiteURL || NormalizeURL(urlStr) == rootNorm {
			mu.Lock()
			scrapeErr = fmt.Errorf("failed to scrape URL: %w", err)
			mu.Unlock()
			return
		}

		mu.Lock()
		failed++
		mu.Unlock()
	})

	// Seed the crawl from sitemap.xml for pages that aren't reachable via
	// in-page links (or are only linked from JavaScript-rendered menus).
	sitemapURLs := discoverSitemapURLs(parsedURL, allowedHosts)

	mu.Lock()
	visited[rootNorm] = true
	mu.Unlock()

	if err := c.Visit(websiteURL); err != nil {
		return nil, 0, fmt.Errorf("failed to visit URL: %w", err)
	}

	for _, u := range sitemapURLs {
		norm := NormalizeURL(u)

		mu.Lock()
		alreadyVisited := visited[norm]
		visited[norm] = true
		atLimit := pageCount >= maxPages
		mu.Unlock()

		if alreadyVisited || atLimit {
			continue
		}
		_ = c.Visit(norm)
	}

	c.Wait()

	if scrapeErr != nil {
		return nil, 0, scrapeErr
	}

	if len(pages) == 0 {
		return nil, 0, fmt.Errorf("no text content found on the website")
	}

	return pages, failed, nil
}

// NormalizeURL canonicalizes a URL for de-duplication: it drops the fragment
// (e.g. tab anchors like #nav-44 that point at the same page) and trims a
// trailing slash so "/admissions" and "/admissions/" aren't crawled/stored twice.
func NormalizeURL(raw string) string {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return strings.TrimSpace(raw)
	}
	u.Fragment = ""
	if u.Path != "/" && u.Path != "" {
		u.Path = strings.TrimRight(u.Path, "/")
	}
	if u.Path == "" {
		u.Path = "/"
	}
	return u.String()
}

// sitemapURLSet mirrors the <url><loc> entries of a urlset sitemap.
type sitemapURLSet struct {
	URLs []sitemapLoc `xml:"url"`
}

// sitemapIndexFile mirrors the <sitemap><loc> entries of a sitemap index.
type sitemapIndexFile struct {
	Sitemaps []sitemapLoc `xml:"sitemap"`
}

type sitemapLoc struct {
	Loc string `xml:"loc"`
}

// discoverSitemapURLs fetches /sitemap.xml (following one level of sitemap
// index nesting) and returns all in-scope page URLs it lists. Best-effort:
// any failure just yields no extra seeds.
func discoverSitemapURLs(base *url.URL, allowedHosts []string) []string {
	client := &http.Client{Timeout: 15 * time.Second}
	candidate := base.Scheme + "://" + base.Host + "/sitemap.xml"
	return fetchSitemap(client, candidate, allowedHosts, 0)
}

func fetchSitemap(client *http.Client, sitemapURL string, allowedHosts []string, depth int) []string {
	if depth > 2 {
		return nil
	}

	req, err := http.NewRequest(http.MethodGet, sitemapURL, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("User-Agent", browserUA)

	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 10*1024*1024))
	if err != nil {
		return nil
	}

	// A sitemap index points at child sitemaps — recurse into them.
	var idx sitemapIndexFile
	if err := xml.Unmarshal(body, &idx); err == nil && len(idx.Sitemaps) > 0 {
		var out []string
		for _, s := range idx.Sitemaps {
			loc := strings.TrimSpace(s.Loc)
			if loc == "" {
				continue
			}
			out = append(out, fetchSitemap(client, loc, allowedHosts, depth+1)...)
			if len(out) > 2000 {
				break
			}
		}
		return out
	}

	var set sitemapURLSet
	if err := xml.Unmarshal(body, &set); err == nil {
		var out []string
		for _, u := range set.URLs {
			loc := strings.TrimSpace(u.Loc)
			if loc != "" && hostAllowed(loc, allowedHosts) {
				out = append(out, loc)
			}
		}
		return out
	}

	return nil
}

func hostAllowed(rawURL string, allowedHosts []string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	for _, h := range allowedHosts {
		if strings.EqualFold(u.Host, h) {
			return true
		}
	}
	return false
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
