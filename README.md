# Bullnose

A powerful web scraper that converts web pages to clean markdown format, with support for recursive link following, content filtering, and smart rescraping.

## Features

- 🔄 Converts web pages to clean markdown format
- 🌲 Recursive link following with configurable depth
- 🔒 Domain restriction option
- ⏱️ Smart rescraping with time-based updates
- ⚡ Parallel processing for faster scraping
- 🎯 Configurable URL filtering
- 🔍 Sitemap.xml parsing support
- 📝 Detailed logging options
- ⚙️ YAML configuration
- 💻 Cross-platform support (Windows, macOS, Linux)

## Installation

### Linux

```bash
# Using curl
curl -L https://github.com/ncecere/bullnose/releases/download/v1.0.0/bullnose-linux-amd64 -o bullnose
chmod +x bullnose
sudo mv bullnose /usr/local/bin/

# Using wget
wget https://github.com/ncecere/bullnose/releases/download/v1.0.0/bullnose-linux-amd64 -O bullnose
chmod +x bullnose
sudo mv bullnose /usr/local/bin/
```

### macOS

```bash
# Using curl
curl -L https://github.com/ncecere/bullnose/releases/download/v1.0.0/bullnose-darwin-amd64 -o bullnose
chmod +x bullnose
sudo mv bullnose /usr/local/bin/

# Using Homebrew
brew tap ncecere/bullnose
brew install bullnose
```

### Windows

```powershell
# Using PowerShell
Invoke-WebRequest -Uri https://github.com/ncecere/bullnose/releases/download/v1.0.0/bullnose-windows-amd64.exe -OutFile bullnose.exe
Move-Item bullnose.exe $env:USERPROFILE\AppData\Local\Microsoft\WindowsApps\

# Using Scoop
scoop bucket add ncecere https://github.com/ncecere/scoop-bucket
scoop install bullnose
```

### Using Go

```bash
# Install latest version
go install github.com/ncecere/bullnose/cmd/bullnose@latest

# Install specific version
go install github.com/ncecere/bullnose/cmd/bullnose@v1.0.0
```

### Using Docker

```bash
# Pull the image
docker pull ncecere/bullnose:latest

# Basic usage
docker run --rm -v $(pwd)/output:/app/scraped-content ncecere/bullnose https://example.com

# Using docker-compose
docker-compose run bullnose https://example.com
```

## Quick Start

```bash
# Basic usage
bullnose https://example.com

# Use configuration file
bullnose -c config.yaml
```

## Documentation

- [Usage Guide](docs/USAGE.md) - Detailed instructions on using Bullnose
- [Build Guide](docs/BUILD.md) - Instructions for building and developing
- [Example Configurations](examples/)
  - [Simple Configuration](examples/config-simple.yaml)
  - [Advanced Configuration](examples/config-advanced.yaml)
  - [Full Configuration Reference](examples/config-full.yaml)
- [Roadmap](ROADMAP.md) - Future development plans

## Example Output

```markdown
# Page Title

## Metadata
- URL: https://example.com/page
- Scraped: 2024-01-01T12:00:00Z

## Content
[Clean, formatted content with preserved structure]
```

## Basic Configuration

```yaml
# Simple configuration example
urls:
  - "https://example.com"
output: "./scraped-content"
depth: 2
restrict-domain: true
```

See [example configurations](examples/) for more options.

## Command Line Options

```bash
bullnose [flags] [urls...]

Flags:
  -c, --config string         Config file path
  -o, --output string        Output directory (default "./scraped-content")
  -d, --depth int           Maximum link depth (default 3)
  -p, --parallel int        Parallel actions (default 8)
  -r, --restrict-domain     Only follow same-domain links (default true)
  -f, --force              Force rescrape
      --debug              Enable debug logging
```

See [Usage Guide](docs/USAGE.md) for complete details.

## Building from Source

```bash
# Clone repository
git clone https://github.com/ncecere/bullnose.git
cd bullnose

# Build
make build

# Run tests
make test
```

See [Build Guide](docs/BUILD.md) for detailed build instructions.

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) and [Code of Conduct](CODE_OF_CONDUCT.md).

1. Fork the repository
2. Create your feature branch
3. Make your changes
4. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
