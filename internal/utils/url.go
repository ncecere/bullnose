package utils

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// GetAllowedDomains extracts domains from URLs
func GetAllowedDomains(urls []string) ([]string, error) {
	domains := make([]string, 0, len(urls))
	for _, u := range urls {
		parsed, err := url.Parse(u)
		if err != nil {
			return nil, err
		}
		domains = append(domains, parsed.Host)
	}
	return domains, nil
}

// GetCommonSitemapURLs returns common sitemap URLs for a given base URL
func GetCommonSitemapURLs(baseURL string) ([]string, error) {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return []string{
		fmt.Sprintf("%s://%s/sitemap.xml", parsed.Scheme, parsed.Host),
		fmt.Sprintf("%s://%s/sitemap_index.xml", parsed.Scheme, parsed.Host),
	}, nil
}

// GlobToRegex converts a glob pattern to a regular expression pattern
func GlobToRegex(pattern string) string {
	// Escape special regex characters
	pattern = regexp.QuoteMeta(pattern)

	// Convert glob patterns to regex patterns
	pattern = strings.ReplaceAll(pattern, "\\*", ".*")
	pattern = strings.ReplaceAll(pattern, "\\?", ".")

	// Add start and end anchors if not present
	if !strings.HasPrefix(pattern, ".*") {
		pattern = "^" + pattern
	}
	if !strings.HasSuffix(pattern, ".*") {
		pattern = pattern + "$"
	}

	return pattern
}

// CreateURLFilters creates regex filters from glob patterns
func CreateURLFilters(patterns []string) ([]*regexp.Regexp, error) {
	filters := make([]*regexp.Regexp, 0, len(patterns))
	for _, pattern := range patterns {
		regexPattern := GlobToRegex(pattern)
		regex, err := regexp.Compile(regexPattern)
		if err != nil {
			return nil, err
		}
		filters = append(filters, regex)
	}
	return filters, nil
}
