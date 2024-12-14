# Bullnose Usage Guide

This guide explains how to use Bullnose, including command-line arguments and configuration options.

## Table of Contents
- [Quick Start](#quick-start)
- [Command-Line Arguments](#command-line-arguments)
- [Configuration File](#configuration-file)
- [Examples](#examples)

## Quick Start

```bash
# Basic usage with a single URL
bullnose https://example.com

# Use a configuration file
bullnose -c config.yaml

# Scrape multiple URLs
bullnose https://example.com https://another-site.com
```

## Command-Line Arguments

### Basic Arguments

```bash
bullnose [flags] [urls...]
```

### Available Flags

| Flag | Short | Default | Description |
|------|--------|---------|-------------|
| `--config` | `-c` | `./bullnose.yaml` | Configuration file path |
| `--output` | `-o` | `./scraped-content` | Output directory for scraped content |
| `--depth` | `-d` | `3` | Maximum depth to follow links |
| `--parallel` | `-p` | `8` | Number of parallel scraping actions |
| `--restrict-domain` | `-r` | `true` | Only follow links within starting domain |
| `--rescrape-after` | | `12h` | Only rescrape after this duration |
| `--force` | `-f` | `false` | Force rescrape regardless of time |
| `--debug` | | `false` | Enable debug logging |
| `--ignore` | | `[]` | URLs or patterns to ignore |

### Flag Details

#### --config, -c
Specify the path to a YAML configuration file. If not provided, Bullnose looks for `bullnose.yaml` in:
1. Current directory
2. $HOME/.config/bullnose
3. /etc/bullnose

```bash
bullnose -c my-config.yaml
```

#### --output, -o
Set the directory where scraped content will be saved. The directory structure will be:
```
output-dir/
  domain1.com/
    page1.md
    page2.md
  domain2.com/
    page1.md
```

```bash
bullnose -o ./my-content https://example.com
```

#### --depth, -d
Control how deep the scraper follows links:
- `1`: Only scrape provided URLs
- `2`: Also scrape pages linked from initial URLs
- `3+`: Continue following links to specified depth

```bash
bullnose -d 2 https://example.com  # Only scrape homepage and direct links
```

#### --parallel, -p
Set the number of concurrent scraping operations. Higher numbers may improve speed but increase server load.

```bash
bullnose -p 4 https://example.com  # Use 4 parallel scrapers
```

#### --restrict-domain, -r
Control whether to follow external links:
- `true`: Only follow links within the starting domain
- `false`: Follow links to any domain

```bash
bullnose -r false https://example.com  # Follow external links
```

#### --rescrape-after
Set the minimum time before rescaping content. Uses Go duration format:
- `30m`: 30 minutes
- `24h`: 24 hours
- `7d`: 7 days

```bash
bullnose --rescrape-after 48h https://example.com
```

#### --force, -f
Force rescrape all content regardless of when it was last scraped.

```bash
bullnose -f https://example.com  # Rescrape everything
```

#### --debug
Enable detailed logging for troubleshooting.

```bash
bullnose --debug https://example.com
```

#### --ignore
Specify patterns for URLs to ignore. Supports glob patterns.

```bash
bullnose --ignore "login,*.pdf,private/*" https://example.com
```

## Configuration File

The configuration file offers more control than command-line arguments. See example configurations:

- [Simple Configuration](../examples/config-simple.yaml): Basic settings for getting started
- [Advanced Configuration](../examples/config-advanced.yaml): More features with explanations
- [Full Configuration](../examples/config-full.yaml): Complete reference of all options

### Configuration Priority

Settings are applied in this order (later overrides earlier):
1. Default values
2. Configuration file
3. Environment variables (prefixed with `BULLNOSE_`)
4. Command-line arguments

### Environment Variables

All settings can be set via environment variables:
- Prefix: `BULLNOSE_`
- Format: Uppercase with underscores
- Examples:
  ```bash
  BULLNOSE_OUTPUT="./content"
  BULLNOSE_DEPTH="3"
  BULLNOSE_RESTRICT_DOMAIN="true"
  ```

## Examples

### Basic Scraping

```bash
# Scrape a single site
bullnose https://example.com

# Scrape multiple sites
bullnose https://site1.com https://site2.com

# Scrape with depth limit
bullnose -d 2 https://example.com

# Scrape in parallel
bullnose -p 4 https://example.com
```

### Advanced Usage

```bash
# Use custom output directory and follow external links
bullnose -o ./docs -r false https://example.com

# Force rescrape with debug logging
bullnose -f --debug https://example.com

# Use configuration file with additional URLs
bullnose -c config.yaml https://extra-site.com

# Ignore specific patterns
bullnose --ignore "login,admin,*.pdf" https://example.com
```

### Content Organization

Scraped content is saved as Markdown files with metadata:

```markdown
# Page Title

## Metadata
- URL: https://example.com/page
- Scraped: 2024-01-01T12:00:00Z

## Content
[Page content with preserved structure]
```

### Error Handling

- Network errors: The scraper will retry failed requests
- Rate limiting: Respects server restrictions
- Invalid URLs: Skipped with warning
- Parse errors: Logged if debug enabled

## Tips and Best Practices

1. Start with a small depth (`-d 2`) to test scraping behavior
2. Use `--debug` when configuring new sites
3. Respect robots.txt and site terms of service
4. Use appropriate delays between requests
5. Monitor server response times and adjust parallel settings
6. Use domain-specific configurations for sites needing special handling
7. Regularly backup your configuration files
8. Test configuration changes on a small subset first
