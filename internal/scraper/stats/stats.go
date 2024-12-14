package stats

import (
	"fmt"
	"sync"
	"time"
)

// Stats tracks scraping statistics
type Stats struct {
	URLsScanned int
	URLsScraped int
	URLsSkipped int
	StartTime   time.Time
	mutex       sync.Mutex
}

// New creates a new Stats tracker
func New() *Stats {
	return &Stats{
		StartTime: time.Now(),
	}
}

// IncrementScanned increments the number of URLs scanned
func (s *Stats) IncrementScanned() {
	s.mutex.Lock()
	s.URLsScanned++
	s.mutex.Unlock()
}

// IncrementScraped increments the number of URLs scraped
func (s *Stats) IncrementScraped() {
	s.mutex.Lock()
	s.URLsScraped++
	s.mutex.Unlock()
}

// IncrementSkipped increments the number of URLs skipped
func (s *Stats) IncrementSkipped() {
	s.mutex.Lock()
	s.URLsSkipped++
	s.mutex.Unlock()
}

// GetSummary returns a formatted summary of the statistics
func (s *Stats) GetSummary() string {
	duration := time.Since(s.StartTime)
	return fmt.Sprintf(`
Scraping Statistics:
URLs Scanned: %d
URLs Scraped: %d
URLs Skipped: %d
Total Time: %s
`, s.URLsScanned, s.URLsScraped, s.URLsSkipped, duration.Round(time.Second))
}

// GetStats returns the current statistics
func (s *Stats) GetStats() (scanned, scraped, skipped int, duration time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.URLsScanned, s.URLsScraped, s.URLsSkipped, time.Since(s.StartTime)
}
