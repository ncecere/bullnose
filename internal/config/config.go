package config

import (
	"fmt"
	"path/filepath"
	"regexp"
	"time"

	"github.com/spf13/viper"
)

// LoadConfig loads configuration from file and environment variables
func LoadConfig(configFile string) (*Config, error) {
	v := viper.New()

	// Set default values
	v.SetDefault("output", "./scraped-content")
	v.SetDefault("depth", 3)
	v.SetDefault("parallel", 8)
	v.SetDefault("restrict-domain", true)
	v.SetDefault("rescrape-after", "12h")
	v.SetDefault("force", false)
	v.SetDefault("debug", false)
	v.SetDefault("parse-sitemaps", true)
	v.SetDefault("ignore", []string{
		"login",
		"admin",
		"logout",
		"/api/",
		"*.pdf",
		"private/*",
	})
	v.SetDefault("urls", []string{})
	v.SetDefault("domain-config", map[string]*DomainConfig{})
	v.SetDefault("content-patterns", map[string]ContentExtraction{})

	// Environment variables
	v.SetEnvPrefix("BULLNOSE")
	v.AutomaticEnv()

	// Config file
	if configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		v.SetConfigName("bullnose")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("$HOME/.config/bullnose")
		v.AddConfigPath("/etc/bullnose")
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Parse rescrape-after duration if it's a string
	if v.GetString("rescrape-after") != "" {
		duration, err := time.ParseDuration(v.GetString("rescrape-after"))
		if err != nil {
			return nil, fmt.Errorf("invalid rescrape-after duration: %w", err)
		}
		config.RescrapeAfter = duration
	}

	// Ensure output path is absolute
	if !filepath.IsAbs(config.Output) {
		absPath, err := filepath.Abs(config.Output)
		if err != nil {
			return nil, fmt.Errorf("error converting output path to absolute: %w", err)
		}
		config.Output = absPath
	}

	// Validate the configuration
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func validateConfig(config *Config) error {
	if config.Depth < 1 {
		return fmt.Errorf("depth must be greater than 0")
	}

	if config.Parallel < 1 {
		return fmt.Errorf("parallel must be greater than 0")
	}

	if config.RescrapeAfter < 0 {
		return fmt.Errorf("rescrape-after must be non-negative")
	}

	// Validate regex patterns
	for domain, extraction := range config.ContentPatterns {
		if extraction.TitlePattern != "" {
			if _, err := regexp.Compile(extraction.TitlePattern); err != nil {
				return fmt.Errorf("invalid title pattern for domain %s: %w", domain, err)
			}
		}
		for _, pattern := range extraction.ContentPatterns {
			if _, err := regexp.Compile(pattern); err != nil {
				return fmt.Errorf("invalid content pattern for domain %s: %w", domain, err)
			}
		}
		for _, pattern := range extraction.ExcludePatterns {
			if _, err := regexp.Compile(pattern); err != nil {
				return fmt.Errorf("invalid exclude pattern for domain %s: %w", domain, err)
			}
		}
	}

	return nil
}
