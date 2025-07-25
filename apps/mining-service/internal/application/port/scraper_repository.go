package port

type ScraperRepository interface {
	ScrapeURLs(targetURL string, sameDomainOnly, includeSubdomains, excludeImages, excludeResources bool) []string
}
