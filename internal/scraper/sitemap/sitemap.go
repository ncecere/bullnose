package sitemap

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

// SitemapIndex represents a sitemap index file
type SitemapIndex struct {
	XMLName  xml.Name `xml:"sitemapindex"`
	Sitemaps []struct {
		Loc string `xml:"loc"`
	} `xml:"sitemap"`
}

// Sitemap represents a sitemap file
type Sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	URLs    []struct {
		Loc string `xml:"loc"`
	} `xml:"url"`
}

// Parser handles sitemap parsing
type Parser struct {
	client *http.Client
}

// NewParser creates a new sitemap parser
func NewParser() *Parser {
	return &Parser{
		client: &http.Client{},
	}
}

// Parse parses a sitemap URL and returns a list of page URLs
func (p *Parser) Parse(url string) ([]string, error) {
	resp, err := p.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching sitemap: %w", err)
	}

	// Ensure the response body is closed after we're done
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading sitemap body: %w", err)
	}

	// Try parsing as sitemap index first
	var index SitemapIndex
	if err := xml.Unmarshal(body, &index); err == nil && len(index.Sitemaps) > 0 {
		var urls []string
		for _, s := range index.Sitemaps {
			subUrls, err := p.Parse(s.Loc)
			if err != nil {
				// Log error but continue with other sitemaps
				fmt.Printf("Error parsing sub-sitemap %s: %v\n", s.Loc, err)
				continue
			}
			urls = append(urls, subUrls...)
		}
		return urls, nil
	}

	// Try parsing as regular sitemap
	var sitemap Sitemap
	if err := xml.Unmarshal(body, &sitemap); err != nil {
		return nil, fmt.Errorf("error parsing sitemap XML: %w", err)
	}

	var urls []string
	for _, u := range sitemap.URLs {
		urls = append(urls, u.Loc)
	}

	return urls, nil
}
