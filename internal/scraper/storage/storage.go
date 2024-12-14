package storage

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Storage handles file operations and URL tracking
type Storage struct {
	outputDir     string
	visitedURLs   sync.Map
	rescrapeAfter time.Duration
	force         bool
}

// New creates a new Storage instance
func New(outputDir string, rescrapeAfter time.Duration, force bool) *Storage {
	return &Storage{
		outputDir:     outputDir,
		rescrapeAfter: rescrapeAfter,
		force:         force,
	}
}

// MarkVisited marks a URL as visited
func (s *Storage) MarkVisited(url string) {
	s.visitedURLs.Store(url, true)
}

// IsVisited checks if a URL has been visited
func (s *Storage) IsVisited(url string) bool {
	_, visited := s.visitedURLs.Load(url)
	return visited
}

// ShouldRescrape determines if a file should be rescraped
func (s *Storage) ShouldRescrape(filepath string) bool {
	// Always scrape if force is enabled
	if s.force {
		return true
	}

	// Check if file exists
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return true
	}
	if err != nil {
		return true
	}

	// Check if enough time has passed since last scrape
	return time.Since(info.ModTime()) >= s.rescrapeAfter
}

// SaveContent saves content to a file
func (s *Storage) SaveContent(domain, title, content string) (string, error) {
	// Create safe filename from title
	filename := s.sanitizeFilename(title)
	if filename == "" {
		filename = s.hashURL(title) // Use title as URL if no valid filename can be created
	}

	// Determine output directory based on domain
	outputDir := filepath.Join(s.outputDir, domain)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create markdown file
	outputPath := filepath.Join(outputDir, filename+".md")

	// Write content to file
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return outputPath, nil
}

// sanitizeFilename creates a safe filename from a string
func (s *Storage) sanitizeFilename(name string) string {
	// Remove invalid characters
	name = strings.Map(func(r rune) rune {
		if r > 127 || strings.ContainsRune(`<>:"/\|?*`, r) {
			return '-'
		}
		return r
	}, name)

	// Trim spaces and dashes
	name = strings.Trim(name, " -")

	// Replace multiple dashes with single dash
	name = strings.Join(strings.Fields(name), "-")

	return strings.ToLower(name)
}

// hashURL creates a hash from a URL string
func (s *Storage) hashURL(url string) string {
	hash := md5.Sum([]byte(url))
	return hex.EncodeToString(hash[:])
}

// GetOutputPath returns the full output path for a domain and filename
func (s *Storage) GetOutputPath(domain, filename string) string {
	return filepath.Join(s.outputDir, domain, filename+".md")
}
