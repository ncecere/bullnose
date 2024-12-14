# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- None

### Changed
- None

### Deprecated
- None

### Removed
- None

### Fixed
- None

### Security
- None

## [1.0.1] - 2024-01-09

### Changed
- Updated GitHub Actions release workflow to use Makefile targets for building and testing
- Simplified CI/CD pipeline by consolidating build steps

## [1.0.0] - 2023-12-13

### Added
- Initial release of Bullnose web scraper
- Command-line interface with cobra
- Configuration management with viper
- Web scraping functionality with colly
- Markdown conversion of web pages
- Recursive link following with configurable depth
- Domain restriction option
- Smart rescraping with time intervals
- Force scrape option
- Parallel processing
- URL ignore patterns
- Debug logging
- YAML configuration support
- Cross-platform support (Windows, macOS, Linux)
- Docker support
- GitHub Actions CI/CD pipeline
- Documentation (README, LICENSE, CONTRIBUTING)

### Fixed
- Enhanced code block handling in markdown generation
  - Improved language detection from multiple class patterns
  - Preserved original code indentation
  - Better handling of nested code blocks
  - Proper spacing around code blocks
  - More reliable markdown formatting

## [0.1.0] - 2024-12-01
- Initial beta release

[Unreleased]: https://github.com/yourusername/bullnose/compare/v1.0.1...HEAD
[1.0.1]: https://github.com/yourusername/bullnose/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/yourusername/bullnose/compare/v0.1.0...v1.0.0
[0.1.0]: https://github.com/yourusername/bullnose/releases/tag/v0.1.0
