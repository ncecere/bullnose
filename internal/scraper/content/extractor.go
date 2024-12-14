package content

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ExtractionPatterns defines patterns for content extraction
type ExtractionPatterns struct {
	TitlePattern    string
	ContentPatterns []string
	ExcludePatterns []string
}

// Extractor handles content extraction from HTML
type Extractor struct {
	patterns map[string]ExtractionPatterns
}

// NewExtractor creates a new content extractor
func NewExtractor(domainPatterns map[string]ExtractionPatterns) *Extractor {
	return &Extractor{
		patterns: domainPatterns,
	}
}

// ExtractTitle extracts the title from HTML using domain-specific patterns if available
func (e *Extractor) ExtractTitle(domain string, selection *goquery.Selection, fallbackText string) string {
	// Try domain-specific pattern first
	if patterns, ok := e.patterns[domain]; ok && patterns.TitlePattern != "" {
		titleRegex := regexp.MustCompile(patterns.TitlePattern)
		text, _ := selection.Html()
		matches := titleRegex.FindStringSubmatch(text)
		if len(matches) > 1 {
			return matches[1]
		}
	}

	// Fallback to standard title extraction
	if title := selection.Find("title").Text(); title != "" {
		return title
	}
	if title := selection.Find("h1").First().Text(); title != "" {
		return title
	}

	return fallbackText
}

// ExtractContent extracts and formats content from HTML
func (e *Extractor) ExtractContent(domain string, selection *goquery.Selection) string {
	var content strings.Builder

	// Try domain-specific patterns first
	if patterns, ok := e.patterns[domain]; ok && len(patterns.ContentPatterns) > 0 {
		html, _ := selection.Html()
		for _, pattern := range patterns.ContentPatterns {
			regex := regexp.MustCompile(pattern)
			matches := regex.FindAllString(html, -1)
			for _, match := range matches {
				content.WriteString(match + "\n\n")
			}
		}

		// Apply exclude patterns
		if len(patterns.ExcludePatterns) > 0 {
			contentStr := content.String()
			for _, pattern := range patterns.ExcludePatterns {
				regex := regexp.MustCompile(pattern)
				contentStr = regex.ReplaceAllString(contentStr, "")
			}
			return contentStr
		}

		return content.String()
	}

	// Fallback to default content extraction
	mainContent := selection.Find("main, article, .content, #content, .post, #post").First()
	if mainContent.Length() == 0 {
		mainContent = selection.Find("body")
	}

	e.extractStructuredContent(&content, mainContent)
	e.extractLinks(&content, mainContent)

	return content.String()
}

// extractStructuredContent extracts content with markdown formatting
func (e *Extractor) extractStructuredContent(content *strings.Builder, selection *goquery.Selection) {
	selection.Find("h1, h2, h3, h4, h5, h6, p, ul, ol, li, pre, code, blockquote").Each(func(i int, s *goquery.Selection) {
		switch s.Get(0).Data {
		case "h1":
			content.WriteString("\n# " + s.Text() + "\n\n")
		case "h2":
			content.WriteString("\n## " + s.Text() + "\n\n")
		case "h3":
			content.WriteString("\n### " + s.Text() + "\n\n")
		case "h4":
			content.WriteString("\n#### " + s.Text() + "\n\n")
		case "h5":
			content.WriteString("\n##### " + s.Text() + "\n\n")
		case "h6":
			content.WriteString("\n###### " + s.Text() + "\n\n")
		case "p":
			content.WriteString(s.Text() + "\n\n")
		case "ul":
			s.Find("li").Each(func(i int, li *goquery.Selection) {
				content.WriteString("- " + li.Text() + "\n")
			})
			content.WriteString("\n")
		case "ol":
			s.Find("li").Each(func(i int, li *goquery.Selection) {
				content.WriteString(fmt.Sprintf("%d. %s\n", i+1, li.Text()))
			})
			content.WriteString("\n")
		case "pre", "code":
			content.WriteString("```\n" + s.Text() + "\n```\n\n")
		case "blockquote":
			content.WriteString("> " + s.Text() + "\n\n")
		}
	})
}

// extractLinks extracts links from the content
func (e *Extractor) extractLinks(content *strings.Builder, selection *goquery.Selection) {
	selection.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			content.WriteString(fmt.Sprintf("[%s](%s)\n", s.Text(), href))
		}
	})
}
