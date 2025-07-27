package database

import (
	"net/url"
	"strings"

	"github.com/ribbinpo/mining-service/internal/application/domain"
)

// extractPathFromURL extracts the path portion from a full URL
// Example: "https://example.com/path/to/page" -> "/path/to/page"
func extractPathFromURL(fullURL string) string {
	// If it's already a path (starts with /), return as is
	if strings.HasPrefix(fullURL, "/") {
		return fullURL
	}

	// Try to parse as URL and extract path
	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		// If parsing fails, return the original string
		return fullURL
	}

	// Return the path portion
	path := parsedURL.Path
	if parsedURL.RawQuery != "" {
		path += "?" + parsedURL.RawQuery
	}
	if parsedURL.Fragment != "" {
		path += "#" + parsedURL.Fragment
	}

	return path
}

func StoneDomainToModel(domain *domain.StoneDomain) *StoneModel {
	paths := make([]PathModel, len(domain.Paths))
	for i, path := range domain.Paths {
		// Extract only the path portion from full URLs
		pathOnly := extractPathFromURL(path)
		paths[i] = PathModel{
			Name: pathOnly,
		}
	}
	return &StoneModel{
		ID:        domain.ID,
		Domain:    domain.DomainURL,
		Version:   domain.Version,
		Status:    string(domain.Status),
		Paths:     paths,
		CreatedAt: domain.CreatedAt,
		UpdatedAt: domain.UpdatedAt,
	}
}

func StoneModelToDomain(model *StoneModel) *domain.StoneDomain {
	paths := make([]string, len(model.Paths))
	for i, path := range model.Paths {
		paths[i] = path.Name
	}
	return &domain.StoneDomain{
		ID:        model.ID,
		DomainURL: model.Domain,
		Version:   model.Version,
		Status:    domain.StoneStatusEnum(model.Status),
		Paths:     paths,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
