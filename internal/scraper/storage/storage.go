package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Storage handles content storage operations
type Storage struct {
	outputDir     string
	rescrapeAfter time.Duration
	force         bool
	visited       map[string]bool
}

// New creates a new Storage instance
func New(outputDir string, rescrapeAfter time.Duration, force bool) *Storage {
	return &Storage{
		outputDir:     outputDir,
		rescrapeAfter: rescrapeAfter,
		force:         force,
		visited:       make(map[string]bool),
	}
}

// IsVisited checks if a URL has been visited
func (s *Storage) IsVisited(url string) bool {
	if s.force {
		return false
	}

	if visited := s.visited[url]; visited {
		return true
	}

	outputPath := s.getOutputPath(url)
	info, err := os.Stat(outputPath)
	if err != nil {
		return false
	}

	// Check if enough time has passed since last scrape
	if s.rescrapeAfter > 0 && time.Since(info.ModTime()) < s.rescrapeAfter {
		return true
	}

	return false
}

// MarkVisited marks a URL as visited
func (s *Storage) MarkVisited(url string) {
	s.visited[url] = true
}

// SaveContent saves content to a file
func (s *Storage) SaveContent(domain, title, content string) (string, error) {
	// Create domain directory
	outputDir := filepath.Join(s.outputDir, domain)
	if err := os.MkdirAll(outputDir, 0700); err != nil {
		return "", fmt.Errorf("error creating output directory: %w", err)
	}

	// Generate filename from title
	filename := s.sanitizeFilename(title)
	if filename == "" {
		filename = s.hashString(title)
	}
	filename += ".md"

	// Create full output path
	outputPath := filepath.Join(outputDir, filename)

	// Write content to file
	if err := os.WriteFile(outputPath, []byte(content), 0600); err != nil {
		return "", fmt.Errorf("error writing content: %w", err)
	}

	return outputPath, nil
}

// sanitizeFilename creates a safe filename from a string
func (s *Storage) sanitizeFilename(name string) string {
	// Replace invalid characters with hyphens
	name = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '-' || r == '_' {
			return r
		}
		return '-'
	}, name)

	// Remove consecutive hyphens
	name = strings.Join(strings.FieldsFunc(name, func(r rune) bool {
		return r == '-'
	}), "-")

	// Trim hyphens from ends
	name = strings.Trim(name, "-")

	// Limit length
	if len(name) > 100 {
		name = name[:100]
	}

	return strings.ToLower(name)
}

// hashString creates a SHA-256 hash of a string
func (s *Storage) hashString(str string) string {
	hash := sha256.Sum256([]byte(str))
	return hex.EncodeToString(hash[:])[:20] // Use first 20 chars of hash
}

// getOutputPath gets the output path for a URL
func (s *Storage) getOutputPath(url string) string {
	filename := s.hashString(url) + ".md"
	domain := strings.Split(strings.TrimPrefix(url, "http://"), "//")[0]
	domain = strings.Split(strings.TrimPrefix(domain, "https://"), "/")[0]
	return filepath.Join(s.outputDir, domain, filename)
}
