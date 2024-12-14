package config

import "time"

// DomainConfig holds domain-specific configuration
type DomainConfig struct {
	Headers map[string]string `mapstructure:"headers"`
	Cookies map[string]string `mapstructure:"cookies"`
}

// ContentExtraction holds configuration for content extraction
type ContentExtraction struct {
	TitlePattern    string   `mapstructure:"title-pattern"`
	ContentPatterns []string `mapstructure:"content-patterns"`
	ExcludePatterns []string `mapstructure:"exclude-patterns"`
}

// Config holds all configuration for the scraper
type Config struct {
	Output          string                       `mapstructure:"output"`
	Depth           int                          `mapstructure:"depth"`
	Parallel        int                          `mapstructure:"parallel"`
	RestrictDomain  bool                         `mapstructure:"restrict-domain"`
	RescrapeAfter   time.Duration                `mapstructure:"rescrape-after"`
	Force           bool                         `mapstructure:"force"`
	Debug           bool                         `mapstructure:"debug"`
	Ignore          []string                     `mapstructure:"ignore"`
	URLs            []string                     `mapstructure:"urls"`
	ParseSitemaps   bool                         `mapstructure:"parse-sitemaps"`
	DomainConfig    map[string]*DomainConfig     `mapstructure:"domain-config"`
	ContentPatterns map[string]ContentExtraction `mapstructure:"content-patterns"`
}
