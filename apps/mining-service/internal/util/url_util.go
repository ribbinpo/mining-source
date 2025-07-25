package util

import (
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/publicsuffix"
)

func ExtractBaseDomain(host string) string {
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

func IsSameDomain(url1, url2 *url.URL, includeSubdomains bool) bool {
	// Extract base domains (remove subdomains)
	domain1 := ExtractBaseDomain(url1.Host)
	domain2 := ExtractBaseDomain(url2.Host)

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

func ExtractURLsFromJS(jsCode string) []string {
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

func ExtractURLsFromMetaRefresh(content string) []string {
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
