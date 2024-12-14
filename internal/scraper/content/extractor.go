package content

import (
	"fmt"
	"log"
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

// CodeBlock represents a parsed code block with its metadata
type CodeBlock struct {
	Language string
	Content  string
	Original string
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
		text, err := selection.Html()
		if err != nil {
			log.Printf("Error getting HTML for title extraction: %v", err)
			return fallbackText
		}
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
		html, err := selection.Html()
		if err != nil {
			log.Printf("Error getting HTML for content extraction: %v", err)
			return ""
		}
		for _, pattern := range patterns.ContentPatterns {
			regex := regexp.MustCompile(pattern)
			matches := regex.FindAllString(html, -1)
			for _, match := range matches {
				content.WriteString(match + "\n")
			}
		}

		// Apply exclude patterns
		if len(patterns.ExcludePatterns) > 0 {
			contentStr := content.String()
			for _, pattern := range patterns.ExcludePatterns {
				regex := regexp.MustCompile(pattern)
				contentStr = regex.ReplaceAllString(contentStr, "")
			}
			return e.cleanupContent(contentStr)
		}

		return e.cleanupContent(content.String())
	}

	// Find the main content area, prioritizing article content
	mainContent := selection.Find("article, [role='main'], main, .main-content, #main-content").First()
	if mainContent.Length() == 0 {
		mainContent = selection.Find(".content, #content, .post, #post").First()
	}
	if mainContent.Length() == 0 {
		mainContent = selection.Find("body")
	}

	// Extract content while preserving structure
	e.extractStructuredContent(&content, mainContent)

	return e.cleanupContent(content.String())
}

// detectCodeLanguage attempts to detect the programming language from element classes
func (e *Extractor) detectCodeLanguage(s *goquery.Selection) string {
	// Common class patterns for code language identification
	languagePatterns := []string{
		"language-([\\w+-]+)",
		"lang-([\\w+-]+)",
		"brush:\\s*([\\w+-]+)",
		"highlight-([\\w+-]+)",
		"([\\w+-]+)-highlight",
	}

	for _, attr := range []string{"class", "data-lang", "data-language"} {
		if val, exists := s.Attr(attr); exists {
			for _, pattern := range languagePatterns {
				if match := regexp.MustCompile(pattern).FindStringSubmatch(val); len(match) > 1 {
					return match[1]
				}
			}
		}
	}

	return ""
}

// preserveIndentation maintains the original code indentation
func (e *Extractor) preserveIndentation(code string) string {
	lines := strings.Split(code, "\n")
	if len(lines) <= 1 {
		return code
	}

	// Find the minimum indentation level
	minIndent := -1
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " \t"))
		if minIndent == -1 || indent < minIndent {
			minIndent = indent
		}
	}

	// Remove common indentation while preserving relative indentation
	if minIndent > 0 {
		for i, line := range lines {
			if len(line) >= minIndent {
				lines[i] = line[minIndent:]
			}
		}
	}

	return strings.Join(lines, "\n")
}

// extractStructuredContent extracts content with markdown formatting
func (e *Extractor) extractStructuredContent(content *strings.Builder, selection *goquery.Selection) {
	// Track headings and links to avoid duplicates
	seenHeadings := make(map[string]bool)
	seenLinks := make(map[string]string)
	codeBlockStack := make([]*CodeBlock, 0)
	inCodeBlock := false

	selection.Find("h1, h2, h3, h4, h5, h6, p, ul, ol, li, pre, code, blockquote, a").Each(func(_ int, s *goquery.Selection) {
		// Skip navigation elements
		if s.ParentsFiltered("nav, .nav, .navigation, .menu, .sidebar, aside").Length() > 0 {
			return
		}

		switch s.Get(0).Data {
		case "h1", "h2", "h3", "h4", "h5", "h6":
			text := strings.TrimSpace(s.Text())
			if text != "" && !seenHeadings[text] {
				seenHeadings[text] = true
				content.WriteString("\n" + strings.Repeat("#", int(s.Get(0).Data[1]-'0')) + " " + text + "\n\n")
			}
		case "p":
			if !inCodeBlock {
				text := strings.TrimSpace(s.Text())
				if text != "" {
					content.WriteString(text + "\n\n")
				}
			}
		case "ul":
			content.WriteString("\n")
			s.Find("li").Each(func(_ int, li *goquery.Selection) {
				text := strings.TrimSpace(li.Text())
				if text != "" {
					content.WriteString("- " + text + "\n")
				}
			})
			content.WriteString("\n")
		case "ol":
			content.WriteString("\n")
			listIndex := 1
			s.Find("li").Each(func(_ int, li *goquery.Selection) {
				text := strings.TrimSpace(li.Text())
				if text != "" {
					content.WriteString(fmt.Sprintf("%d. %s\n", listIndex, text))
					listIndex++
				}
			})
			content.WriteString("\n")
		case "pre":
			lang := e.detectCodeLanguage(s)
			if s.Find("code").Length() > 0 {
				// Handle pre > code structure
				s.Find("code").Each(func(_ int, code *goquery.Selection) {
					if codeLang := e.detectCodeLanguage(code); codeLang != "" {
						lang = codeLang
					}
					text := e.preserveIndentation(code.Text())
					if text != "" {
						block := &CodeBlock{
							Language: lang,
							Content:  text,
							Original: text,
						}
						codeBlockStack = append(codeBlockStack, block)
						inCodeBlock = true
						content.WriteString("\n```" + lang + "\n")
						content.WriteString(text)
						content.WriteString("\n```\n\n")
						inCodeBlock = false
					}
				})
			} else {
				// Handle pre without code
				text := e.preserveIndentation(s.Text())
				if text != "" {
					block := &CodeBlock{
						Language: lang,
						Content:  text,
						Original: text,
					}
					codeBlockStack = append(codeBlockStack, block)
					inCodeBlock = true
					content.WriteString("\n```" + lang + "\n")
					content.WriteString(text)
					content.WriteString("\n```\n\n")
					inCodeBlock = false
				}
			}
		case "code":
			if !inCodeBlock && s.Parent().Get(0).Data != "pre" {
				text := strings.TrimSpace(s.Text())
				if text != "" {
					content.WriteString("`" + text + "`")
				}
			}
		case "blockquote":
			text := strings.TrimSpace(s.Text())
			if text != "" {
				content.WriteString("> " + text + "\n\n")
			}
		case "a":
			if href, exists := s.Attr("href"); exists {
				text := strings.TrimSpace(s.Text())
				if text != "" && href != "" && !strings.HasPrefix(href, "#") {
					seenLinks[text] = href
				}
			}
		}
	})

	// Add unique links at the end
	if len(seenLinks) > 0 {
		content.WriteString("\n## Links\n\n")
		for text, href := range seenLinks {
			content.WriteString(fmt.Sprintf("[%s](%s)\n", text, href))
		}
	}
}

// cleanupContent performs final cleanup of the extracted content
func (e *Extractor) cleanupContent(content string) string {
	// Split content into lines for processing
	lines := strings.Split(content, "\n")
	var processedLines []string
	var inCodeBlock bool

	// Process each line
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Handle code block boundaries
		if strings.HasPrefix(trimmedLine, "```") {
			if !inCodeBlock {
				// Starting a code block
				inCodeBlock = true
				processedLines = append(processedLines, "", line) // Add empty line before code block
			} else {
				// Ending a code block
				inCodeBlock = false
				processedLines = append(processedLines, line, "") // Add empty line after code block
			}
			continue
		}

		// Handle content inside code blocks
		if inCodeBlock {
			processedLines = append(processedLines, line) // Preserve original indentation
			continue
		}

		// Handle non-code content
		if trimmedLine != "" {
			processedLines = append(processedLines, trimmedLine)
		} else if len(processedLines) > 0 && processedLines[len(processedLines)-1] != "" {
			processedLines = append(processedLines, "") // Add empty line for spacing
		}
	}

	// Join lines and perform final cleanup
	content = strings.Join(processedLines, "\n")

	// Remove any remaining multiple consecutive empty lines
	content = regexp.MustCompile(`\n{3,}`).ReplaceAllString(content, "\n\n")

	// Ensure proper spacing around headers
	content = regexp.MustCompile(`([^\n])(#{1,6}\s)`).ReplaceAllString(content, "$1\n\n$2")
	content = regexp.MustCompile(`(#{1,6}[^\n]+)\n([^#\n])`).ReplaceAllString(content, "$1\n\n$2")

	// Clean up list formatting
	content = regexp.MustCompile(`\n{2,}(-|\d+\.) `).ReplaceAllString(content, "\n- ")

	// Final trim
	return strings.TrimSpace(content)
}
