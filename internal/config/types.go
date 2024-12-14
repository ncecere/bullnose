package config

import "time"

// ContentExtraction defines patterns for extracting content from specific domains
type ContentExtraction struct {
	TitlePattern    string   `yaml:"title_pattern"`
	ContentPatterns []string `yaml:"content_patterns"`
	ExcludePatterns []string `yaml:"exclude_patterns"`
}

// DomainConfig defines domain-specific configuration
type DomainConfig struct {
	Headers map[string]string `yaml:"headers"`
	Cookies map[string]string `yaml:"cookies"`
}

// Config represents the application configuration
type Config struct {
	// URLs to scrape
	URLs []string `yaml:"urls"`

	// Output directory for scraped content
	Output string `yaml:"output"`

	// Maximum depth to follow links (0 = unlimited)
	Depth int `yaml:"depth"`

	// Number of parallel scraping actions
	Parallel int `yaml:"parallel"`

	// Only follow links within starting domain
	RestrictDomain bool `yaml:"restrict_domain"`

	// Only rescrape after this duration
	RescrapeAfter time.Duration `yaml:"rescrape_after"`

	// Force rescrape regardless of time
	Force bool `yaml:"force"`

	// Enable debug logging
	Debug bool `yaml:"debug"`

	// Parse sitemap.xml files
	ParseSitemaps bool `yaml:"parse_sitemaps"`

	// URLs or patterns to ignore
	Ignore []string `yaml:"ignore"`

	// Domain-specific configuration
	DomainConfig map[string]DomainConfig `yaml:"domain_config"`

	// Content extraction patterns per domain
	ContentPatterns map[string]ContentExtraction `yaml:"content_patterns"`
}
