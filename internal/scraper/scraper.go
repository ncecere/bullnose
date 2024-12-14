package scraper

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"

	"github.com/ncecere/bullnose/internal/config"
	"github.com/ncecere/bullnose/internal/scraper/content"
	"github.com/ncecere/bullnose/internal/scraper/sitemap"
	"github.com/ncecere/bullnose/internal/scraper/stats"
	"github.com/ncecere/bullnose/internal/scraper/storage"
	"github.com/ncecere/bullnose/internal/utils"
)

// Scraper coordinates the web scraping process
type Scraper struct {
	config    *config.Config
	collector *colly.Collector
	stats     *stats.Stats
	storage   *storage.Storage
	extractor *content.Extractor
}

// New creates a new Scraper instance
func New(cfg *config.Config) (*Scraper, error) {
	// Create collector with configuration
	c := colly.NewCollector(
		colly.MaxDepth(cfg.Depth),
		colly.Async(cfg.Parallel > 1),
		colly.URLFilters(
			regexp.MustCompile(`^https?://[^/]+(?:/.*)?$`), // Allow base domain and paths
		),
	)

	// Set parallel limit if enabled
	if cfg.Parallel > 1 {
		if err := c.Limit(&colly.LimitRule{
			DomainGlob:  "*",
			Parallelism: cfg.Parallel,
			RandomDelay: 1 * time.Second,
		}); err != nil {
			return nil, fmt.Errorf("error setting rate limit: %w", err)
		}
	}

	// Configure domain restriction if enabled
	if cfg.RestrictDomain && len(cfg.URLs) > 0 {
		domains, err := utils.GetAllowedDomains(cfg.URLs)
		if err != nil {
			return nil, fmt.Errorf("error getting allowed domains: %w", err)
		}
		c.AllowedDomains = domains
	}

	// Configure URL filters
	filters, err := utils.CreateURLFilters(cfg.Ignore)
	if err != nil {
		return nil, fmt.Errorf("error creating URL filters: %w", err)
	}
	c.DisallowedURLFilters = filters

	// Initialize components
	s := &Scraper{
		config:    cfg,
		collector: c,
		stats:     stats.New(),
		storage:   storage.New(cfg.Output, cfg.RescrapeAfter, cfg.Force),
		extractor: content.NewExtractor(convertContentPatterns(cfg.ContentPatterns)),
	}

	s.setupCallbacks()
	return s, nil
}

// Start begins the scraping process
func (s *Scraper) Start() error {
	// Process sitemaps if enabled
	if s.config.ParseSitemaps {
		sitemapParser := sitemap.NewParser()
		for _, baseURL := range s.config.URLs {
			sitemapURLs, err := utils.GetCommonSitemapURLs(baseURL)
			if err != nil {
				log.Printf("Error getting sitemap URLs for %s: %v", baseURL, err)
				continue
			}

			for _, sitemapURL := range sitemapURLs {
				urls, err := sitemapParser.Parse(sitemapURL)
				if err != nil {
					if s.config.Debug {
						log.Printf("Error parsing sitemap %s: %v", sitemapURL, err)
					}
					continue
				}

				// Add discovered URLs to the scraping queue
				for _, u := range urls {
					if !s.storage.IsVisited(u) {
						if err := s.collector.Visit(u); err != nil &&
							err != colly.ErrAlreadyVisited &&
							s.config.Debug {
							log.Printf("Error visiting %s: %v", u, err)
						}
					}
				}
			}
		}
	}

	// Process regular URLs
	for _, u := range s.config.URLs {
		if err := s.collector.Visit(u); err != nil && err != colly.ErrAlreadyVisited {
			return fmt.Errorf("error visiting %s: %w", u, err)
		}
	}

	s.collector.Wait()

	// Print statistics
	fmt.Print(s.stats.GetSummary())

	return nil
}

func (s *Scraper) setupCallbacks() {
	// Set up custom headers and cookies for each request
	s.collector.OnRequest(func(r *colly.Request) {
		if s.storage.IsVisited(r.URL.String()) {
			r.Abort()
			return
		}
		s.storage.MarkVisited(r.URL.String())

		s.stats.IncrementScanned()
		if s.config.Debug {
			log.Printf("Visiting %s", r.URL)
		}

		// Add domain-specific headers and cookies
		if domainCfg, ok := s.config.DomainConfig[r.URL.Host]; ok {
			for key, value := range domainCfg.Headers {
				r.Headers.Set(key, value)
			}
			for key, value := range domainCfg.Cookies {
				r.Headers.Set("Cookie", fmt.Sprintf("%s=%s", key, value))
			}
		}
	})

	// Set up link following
	s.collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if s.config.Debug {
			log.Printf("Found link: %s", link)
		}
		if !s.storage.IsVisited(e.Request.AbsoluteURL(link)) {
			if err := e.Request.Visit(link); err != nil &&
				err != colly.ErrAlreadyVisited &&
				s.config.Debug {
				log.Printf("Error visiting link %s: %v", link, err)
			}
		}
	})

	s.collector.OnHTML("html", func(e *colly.HTMLElement) {
		// Extract title
		title := s.extractor.ExtractTitle(
			e.Request.URL.Host,
			e.DOM,
			e.Request.URL.Path,
		)

		// Extract content
		content := s.extractor.ExtractContent(e.Request.URL.Host, e.DOM)

		// Generate markdown content with consistent formatting
		var markdown strings.Builder
		markdown.WriteString(fmt.Sprintf("# %s\n", title))
		markdown.WriteString("\n## Metadata\n")
		markdown.WriteString(fmt.Sprintf("- URL: %s\n", e.Request.URL.String()))
		markdown.WriteString(fmt.Sprintf("- Scraped: %s\n", time.Now().UTC().Format(time.RFC3339)))
		markdown.WriteString("\n## Content\n")
		markdown.WriteString(content)

		// Save content
		outputPath, err := s.storage.SaveContent(e.Request.URL.Host, title, markdown.String())
		if err != nil {
			log.Printf("Error saving content for %s: %v", e.Request.URL, err)
			return
		}

		s.stats.IncrementScraped()

		if s.config.Debug {
			log.Printf("Scraped %s -> %s", e.Request.URL, outputPath)
		}
	})

	s.collector.OnError(func(r *colly.Response, err error) {
		if err != colly.ErrAlreadyVisited {
			log.Printf("Error scraping %s: %v", r.Request.URL, err)
		}
	})

	s.collector.OnResponse(func(r *colly.Response) {
		if s.config.Debug {
			log.Printf("Got response from %s: %d bytes", r.Request.URL, len(r.Body))
		}
	})
}

// convertContentPatterns converts config content patterns to extractor patterns
func convertContentPatterns(configPatterns map[string]config.ContentExtraction) map[string]content.ExtractionPatterns {
	patterns := make(map[string]content.ExtractionPatterns)
	for domain, p := range configPatterns {
		patterns[domain] = content.ExtractionPatterns{
			TitlePattern:    p.TitlePattern,
			ContentPatterns: p.ContentPatterns,
			ExcludePatterns: p.ExcludePatterns,
		}
	}
	return patterns
}
