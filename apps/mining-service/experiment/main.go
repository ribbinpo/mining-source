package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"golang.org/x/net/publicsuffix"
)

func main() {
	// Command line flags
	start := time.Now()
	targetURL := flag.String("url", "https://example.com", "URL to scrape")
	sameDomain := flag.Bool("same-domain", true, "Filter to only include URLs from the same domain")
	includeSubdomains := flag.Bool("include-subdomains", false, "Include URLs from subdomains")
	excludeImages := flag.Bool("exclude-images", true, "Exclude image files from results")
	excludeResources := flag.Bool("exclude-resources", true, "Exclude CSS, JS, and other resource files")
	flag.Parse()

	fmt.Printf("Scraping: %s\n", *targetURL)
	fmt.Printf("Filters: same-domain=%v, include-subdomains=%v, exclude-images=%v, exclude-resources=%v\n", *sameDomain, *includeSubdomains, *excludeImages, *excludeResources)
	fmt.Println("Found URLs and paths:")

	urls := scrapeURLs(*targetURL, *sameDomain, *includeSubdomains, *excludeImages, *excludeResources)
	if len(urls) == 0 {
		fmt.Println("No URLs found")
		os.Exit(1)
	}

	elapsed := time.Since(start)
	fmt.Printf("Time taken: %s\n", elapsed)

	for _, u := range urls {
		fmt.Println(u)
	}
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
			urls := extractURLsFromJS(onclick)
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
		urls := extractURLsFromJS(scriptContent)
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
			urls := extractURLsFromMetaRefresh(content)
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
	if excludeImages && isImageFile(parsedURL.Path) {
		return false
	}

	// Check if it's a resource file (if excludeResources is true)
	if excludeResources && isResourceFile(parsedURL.Path) {
		return false
	}

	// Check if it's from the same domain (if sameDomainOnly is true)
	if sameDomainOnly {
		return isSameDomain(parsedURL, baseURL, includeSubdomains)
	}

	return true
}

func isSameDomain(url1, url2 *url.URL, includeSubdomains bool) bool {
	// Extract base domains (remove subdomains)
	domain1 := extractBaseDomain(url1.Host)
	domain2 := extractBaseDomain(url2.Host)

	// Compare base domains (case-insensitive)
	if !strings.EqualFold(domain1, domain2) {
		return false
	}

	// If we don't want subdomains, check if the hostnames are exactly the same
	if !includeSubdomains {
		return strings.EqualFold(url1.Host, url2.Host)
	}

	return true
}

func extractBaseDomain(host string) string {
	// Remove port if present
	if colonIndex := strings.Index(host, ":"); colonIndex != -1 {
		host = host[:colonIndex]
	}

	// Use the publicsuffix library to get the effective top-level domain
	etld, err := publicsuffix.EffectiveTLDPlusOne(host)
	if err != nil {
		// Fallback to the original host if there's an error
		return host
	}

	return etld
}

func isImageFile(path string) bool {
	// Get file extension
	ext := strings.ToLower(filepath.Ext(path))

	// List of common image extensions
	imageExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".svg":  true,
		".webp": true,
		".ico":  true,
		".tiff": true,
		".tif":  true,
	}

	return imageExtensions[ext]
}

func isResourceFile(path string) bool {
	// Get file extension
	ext := strings.ToLower(filepath.Ext(path))

	// List of common resource file extensions
	resourceExtensions := map[string]bool{
		".css":   true,
		".js":    true,
		".json":  true,
		".xml":   true,
		".txt":   true,
		".pdf":   true,
		".doc":   true,
		".docx":  true,
		".xls":   true,
		".xlsx":  true,
		".ppt":   true,
		".pptx":  true,
		".zip":   true,
		".rar":   true,
		".7z":    true,
		".tar":   true,
		".gz":    true,
		".mp3":   true,
		".mp4":   true,
		".avi":   true,
		".mov":   true,
		".wmv":   true,
		".flv":   true,
		".webm":  true,
		".woff":  true,
		".woff2": true,
		".ttf":   true,
		".eot":   true,
		".otf":   true,
		".map":   true, // Source maps
		".min":   true, // Minified files
	}

	return resourceExtensions[ext]
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

func extractURLsFromJS(jsCode string) []string {
	var urls []string

	// Pattern for window.location.href = "url"
	locationPattern := regexp.MustCompile(`window\.location\.href\s*=\s*["']([^"']+)["']`)
	matches := locationPattern.FindAllStringSubmatch(jsCode, -1)
	for _, match := range matches {
		if len(match) > 1 {
			urls = append(urls, match[1])
		}
	}

	// Pattern for window.location = "url"
	locationPattern2 := regexp.MustCompile(`window\.location\s*=\s*["']([^"']+)["']`)
	matches2 := locationPattern2.FindAllStringSubmatch(jsCode, -1)
	for _, match := range matches2 {
		if len(match) > 1 {
			urls = append(urls, match[1])
		}
	}

	// Pattern for window.open("url")
	openPattern := regexp.MustCompile(`window\.open\s*\(\s*["']([^"']+)["']`)
	matches3 := openPattern.FindAllStringSubmatch(jsCode, -1)
	for _, match := range matches3 {
		if len(match) > 1 {
			urls = append(urls, match[1])
		}
	}

	// Pattern for document.location.href = "url"
	docLocationPattern := regexp.MustCompile(`document\.location\.href\s*=\s*["']([^"']+)["']`)
	matches4 := docLocationPattern.FindAllStringSubmatch(jsCode, -1)
	for _, match := range matches4 {
		if len(match) > 1 {
			urls = append(urls, match[1])
		}
	}

	// Pattern for href = "url" in onclick
	hrefPattern := regexp.MustCompile(`href\s*=\s*["']([^"']+)["']`)
	matches5 := hrefPattern.FindAllStringSubmatch(jsCode, -1)
	for _, match := range matches5 {
		if len(match) > 1 {
			urls = append(urls, match[1])
		}
	}

	// Pattern for location.href = "url"
	locationHrefPattern := regexp.MustCompile(`location\.href\s*=\s*["']([^"']+)["']`)
	matches6 := locationHrefPattern.FindAllStringSubmatch(jsCode, -1)
	for _, match := range matches6 {
		if len(match) > 1 {
			urls = append(urls, match[1])
		}
	}

	// Pattern for location.replace("url")
	locationReplacePattern := regexp.MustCompile(`location\.replace\s*\(\s*["']([^"']+)["']`)
	matches7 := locationReplacePattern.FindAllStringSubmatch(jsCode, -1)
	for _, match := range matches7 {
		if len(match) > 1 {
			urls = append(urls, match[1])
		}
	}

	return urls
}

func extractURLsFromMetaRefresh(content string) []string {
	var urls []string

	// Pattern for meta refresh: content="5;url=http://example.com"
	refreshPattern := regexp.MustCompile(`url=([^;]+)`)
	matches := refreshPattern.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 1 {
			urls = append(urls, strings.TrimSpace(match[1]))
		}
	}

	return urls
}
