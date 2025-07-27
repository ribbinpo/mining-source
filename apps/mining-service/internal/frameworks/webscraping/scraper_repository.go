package webscraping

import (
	"log"
	"net/url"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/ribbinpo/mining-service/internal/application/port"
	"github.com/ribbinpo/mining-service/internal/util"
)

type scraperRepository struct {
}

func NewScraperRepository() port.ScraperRepository {
	return &scraperRepository{}
}

func (r *scraperRepository) ScrapeURLs(targetURL string, sameDomainOnly, includeSubdomains, excludeImages, excludeResources bool) []string {
	urls := scrapeURLs(targetURL, sameDomainOnly, includeSubdomains, excludeImages, excludeResources)
	if len(urls) == 0 {
		return nil
	}

	return urls
}

func scrapeURLs(targetURL string, sameDomainOnly, includeSubdomains, excludeImages, excludeResources bool) []string {
	var foundURLs []string
	seen := make(map[string]bool)

	// Parse base URL for domain filtering
	baseURL, err := url.Parse(targetURL)
	if err != nil {
		log.Printf("Error parsing target URL: %v", err)
		return foundURLs
	}

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
		colly.MaxDepth(1), // Only scrape the initial page, not follow links
	)

	// Find all <a> tags with href attributes
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if href != "" {
			url := normalizeURL(href, targetURL)
			if url != "" && !seen[url] && shouldIncludeURL(url, baseURL, sameDomainOnly, includeSubdomains, excludeImages, excludeResources) {
				seen[url] = true
				foundURLs = append(foundURLs, url)
			}
		}
	})

	// Find onclick handlers in buttons and other elements
	c.OnHTML("*", func(e *colly.HTMLElement) {
		onclick := e.Attr("onclick")
		if onclick != "" {
			urls := util.ExtractURLsFromJS(onclick)
			for _, u := range urls {
				url := normalizeURL(u, targetURL)
				if url != "" && !seen[url] && shouldIncludeURL(url, baseURL, sameDomainOnly, includeSubdomains, excludeImages, excludeResources) {
					seen[url] = true
					foundURLs = append(foundURLs, url)
				}
			}
		}
	})

	// Find URLs in script tags
	c.OnHTML("script", func(e *colly.HTMLElement) {
		scriptContent := e.Text
		urls := util.ExtractURLsFromJS(scriptContent)
		for _, u := range urls {
			url := normalizeURL(u, targetURL)
			if url != "" && !seen[url] && shouldIncludeURL(url, baseURL, sameDomainOnly, includeSubdomains, excludeImages, excludeResources) {
				seen[url] = true
				foundURLs = append(foundURLs, url)
			}
		}
	})

	// Find form action URLs
	c.OnHTML("form[action]", func(e *colly.HTMLElement) {
		action := e.Attr("action")
		if action != "" {
			url := normalizeURL(action, targetURL)
			if url != "" && !seen[url] && shouldIncludeURL(url, baseURL, sameDomainOnly, includeSubdomains, excludeImages, excludeResources) {
				seen[url] = true
				foundURLs = append(foundURLs, url)
			}
		}
	})

	// Find meta refresh URLs
	c.OnHTML("meta[http-equiv='refresh']", func(e *colly.HTMLElement) {
		content := e.Attr("content")
		if content != "" {
			urls := util.ExtractURLsFromMetaRefresh(content)
			for _, u := range urls {
				url := normalizeURL(u, targetURL)
				if url != "" && !seen[url] && shouldIncludeURL(url, baseURL, sameDomainOnly, includeSubdomains, excludeImages, excludeResources) {
					seen[url] = true
					foundURLs = append(foundURLs, url)
				}
			}
		}
	})

	// Find iframe src URLs
	c.OnHTML("iframe[src]", func(e *colly.HTMLElement) {
		src := e.Attr("src")
		if src != "" {
			url := normalizeURL(src, targetURL)
			if url != "" && !seen[url] && shouldIncludeURL(url, baseURL, sameDomainOnly, includeSubdomains, excludeImages, excludeResources) {
				seen[url] = true
				foundURLs = append(foundURLs, url)
			}
		}
	})

	// Find img src URLs (only if not excluding images)
	if !excludeImages {
		c.OnHTML("img[src]", func(e *colly.HTMLElement) {
			src := e.Attr("src")
			if src != "" {
				url := normalizeURL(src, targetURL)
				if url != "" && !seen[url] && shouldIncludeURL(url, baseURL, sameDomainOnly, includeSubdomains, excludeImages, excludeResources) {
					seen[url] = true
					foundURLs = append(foundURLs, url)
				}
			}
		})
	}

	// Find link href URLs (only if not excluding resources)
	if !excludeResources {
		c.OnHTML("link[href]", func(e *colly.HTMLElement) {
			href := e.Attr("href")
			if href != "" {
				url := normalizeURL(href, targetURL)
				if url != "" && !seen[url] && shouldIncludeURL(url, baseURL, sameDomainOnly, includeSubdomains, excludeImages, excludeResources) {
					seen[url] = true
					foundURLs = append(foundURLs, url)
				}
			}
		})
	}

	err = c.Visit(targetURL)
	if err != nil {
		log.Printf("Error visiting %s: %v", targetURL, err)
	}

	return foundURLs
}

func shouldIncludeURL(urlStr string, baseURL *url.URL, sameDomainOnly, includeSubdomains, excludeImages, excludeResources bool) bool {
	// Parse the URL to check domain and file type
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	// Check if it's an image file (if excludeImages is true)
	if excludeImages && util.IsImageFile(parsedURL.Path) {
		return false
	}

	// Check if it's a resource file (if excludeResources is true)
	if excludeResources && util.IsResourceFile(parsedURL.Path) {
		return false
	}

	// Check if URL contains /. which could indicate directory traversal
	if strings.Contains(parsedURL.Path, "/.") {
		return false
	}

	if strings.Contains(parsedURL.Path, "#") {
		return false
	}

	// Check if it's from the same domain (if sameDomainOnly is true)
	if sameDomainOnly {
		return util.IsSameDomain(parsedURL, baseURL, includeSubdomains)
	}

	return true
}

func normalizeURL(href, baseURL string) string {
	href = strings.TrimSpace(href)

	// Skip javascript: and other non-http protocols
	if strings.HasPrefix(href, "javascript:") ||
		strings.HasPrefix(href, "mailto:") ||
		strings.HasPrefix(href, "tel:") ||
		strings.HasPrefix(href, "#") ||
		strings.HasPrefix(href, "data:") {
		return ""
	}

	// If it's already a full URL, return as is
	if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
		return href
	}

	// If it's a relative path, make it absolute
	if strings.HasPrefix(href, "/") {
		parsed, err := url.Parse(baseURL)
		if err != nil {
			return ""
		}
		return parsed.Scheme + "://" + parsed.Host + href
	}

	// If it's a relative path without leading slash, resolve against base
	if !strings.HasPrefix(href, "http") && !strings.HasPrefix(href, "//") {
		parsed, err := url.Parse(baseURL)
		if err != nil {
			return ""
		}
		// Simple path resolution
		if strings.HasSuffix(parsed.Path, "/") {
			return parsed.Scheme + "://" + parsed.Host + parsed.Path + href
		} else {
			// Remove filename from path
			lastSlash := strings.LastIndex(parsed.Path, "/")
			if lastSlash > 0 {
				basePath := parsed.Path[:lastSlash+1]
				return parsed.Scheme + "://" + parsed.Host + basePath + href
			}
		}
	}

	return href
}
