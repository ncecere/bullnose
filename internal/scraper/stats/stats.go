package stats

import (
	"fmt"
	"sync/atomic"
	"time"
)

// Stats tracks scraping statistics
type Stats struct {
	startTime time.Time
	scanned   atomic.Int64
	scraped   atomic.Int64
}

// New creates a new Stats instance
func New() *Stats {
	return &Stats{
		startTime: time.Now(),
	}
}

// IncrementScanned increments the scanned counter
func (s *Stats) IncrementScanned() {
	s.scanned.Add(1)
}

// IncrementScraped increments the scraped counter
func (s *Stats) IncrementScraped() {
	s.scraped.Add(1)
}

// GetSummary returns a summary of the scraping statistics
func (s *Stats) GetSummary() string {
	duration := time.Since(s.startTime)
	return fmt.Sprintf("\nScraping Summary:\n"+
		"Duration: %v\n"+
		"URLs Scanned: %d\n"+
		"Pages Scraped: %d\n",
		duration.Round(time.Second),
		s.scanned.Load(),
		s.scraped.Load())
}
