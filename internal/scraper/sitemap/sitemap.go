package sitemap

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Sitemap represents a standard XML sitemap
type Sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	URLs    []struct {
		Loc string `xml:"loc"`
	} `xml:"url"`
}

// SitemapIndex represents a sitemap index file
type SitemapIndex struct {
	XMLName  xml.Name `xml:"sitemapindex"`
	Sitemaps []struct {
		Loc string `xml:"loc"`
	} `xml:"sitemap"`
}

// Parser handles sitemap parsing operations
type Parser struct {
	client *http.Client
}

// NewParser creates a new sitemap parser
func NewParser() *Parser {
	return &Parser{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Parse attempts to parse a sitemap URL and returns all discovered URLs
func (p *Parser) Parse(sitemapURL string) ([]string, error) {
	resp, err := p.client.Get(sitemapURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sitemap: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read sitemap body: %w", err)
	}

	urls, err := p.parseContent(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse sitemap content: %w", err)
	}

	return urls, nil
}

// parseContent attempts to parse the XML content as either a sitemap or sitemap index
func (p *Parser) parseContent(content []byte) ([]string, error) {
	var urls []string

	// Try parsing as sitemap
	var sitemap Sitemap
	if err := xml.Unmarshal(content, &sitemap); err == nil {
		for _, url := range sitemap.URLs {
			urls = append(urls, url.Loc)
		}
		return urls, nil
	}

	// Try parsing as sitemap index
	var sitemapIndex SitemapIndex
	if err := xml.Unmarshal(content, &sitemapIndex); err == nil {
		for _, sitemap := range sitemapIndex.Sitemaps {
			subUrls, err := p.Parse(sitemap.Loc)
			if err != nil {
				// Log error but continue with other sitemaps
				continue
			}
			urls = append(urls, subUrls...)
		}
		return urls, nil
	}

	return nil, fmt.Errorf("content is neither a valid sitemap nor sitemap index")
}
