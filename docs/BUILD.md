# Building and Development Guide

This guide explains how to set up your development environment and work on the Bullnose project.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Development Setup](#development-setup)
- [Building](#building)
- [Testing](#testing)
- [Development Workflow](#development-workflow)
- [Project Structure](#project-structure)
- [Contributing](#contributing)

## Prerequisites

### Required Software
- Go 1.21 or later
- Make
- Git
- A code editor (VSCode recommended)

### Recommended VSCode Extensions
- Go extension (ms-vscode.go)
- YAML extension (redhat.vscode-yaml)
- GitLens (eamodio.gitlens)
- EditorConfig (editorconfig.editorconfig)

### Optional Tools
- golangci-lint (for code linting)
- mockgen (for generating mocks)
- godoc (for documentation)

## Development Setup

1. Clone the repository:
```bash
git clone https://github.com/yourusername/bullnose.git
cd bullnose
```

2. Install dependencies:
```bash
go mod download
```

3. Install development tools:
```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install mockgen
go install github.com/golang/mock/mockgen@latest

# Install godoc
go install golang.org/x/tools/cmd/godoc@latest
```

## Building

### Basic Build
```bash
# Build for current platform
make build

# Run the built binary
./build/bullnose
```

### Cross-Platform Building
```bash
# Build for all platforms
make build-all

# Build for specific platforms
make build-linux
make build-windows
make build-macos
```

### Development Build
```bash
# Build with debug information
make build-debug

# Build and run with race detection
make build-race
```

### Clean Build
```bash
# Remove build artifacts
make clean

# Clean and rebuild
make clean build
```

## Testing

### Running Tests
```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests for a specific package
go test ./internal/scraper/...

# Run tests with race detection
make test-race
```

### Generating Mocks
```bash
# Generate mocks for interfaces
make generate-mocks
```

### Linting
```bash
# Run linter
make lint

# Fix auto-fixable issues
make lint-fix
```

## Development Workflow

1. Create a new branch for your feature/fix:
```bash
git checkout -b feature/your-feature-name
```

2. Make your changes following these steps:
   - Write tests first (TDD approach)
   - Implement your changes
   - Run tests and linting
   - Update documentation if needed

3. Commit your changes:
```bash
git add .
git commit -m "feat: add your feature description"
```

4. Push and create a pull request:
```bash
git push origin feature/your-feature-name
```

### Code Style
- Follow standard Go conventions
- Use gofmt for formatting
- Follow project structure conventions
- Write meaningful commit messages (conventional commits)

## Project Structure

```
bullnose/
├── cmd/                    # Command line applications
│   └── bullnose/          # Main application
├── docs/                   # Documentation
├── examples/              # Example configurations
├── internal/              # Internal packages
│   ├── config/           # Configuration handling
│   ├── scraper/          # Core scraping functionality
│   │   ├── content/     # Content extraction
│   │   ├── sitemap/     # Sitemap parsing
│   │   ├── stats/       # Statistics tracking
│   │   └── storage/     # File operations
│   └── utils/           # Common utilities
├── test/                 # Integration tests
└── scripts/             # Build and development scripts
```

### Package Guidelines

1. cmd/bullnose/
   - Contains main.go
   - Handles CLI setup and configuration
   - Minimal business logic

2. internal/config/
   - Configuration structs and loading
   - Validation logic
   - Environment variable handling

3. internal/scraper/
   - Core scraping logic
   - Modular components in subpackages
   - Clean interfaces between components

4. test/
   - Integration tests
   - Test fixtures and helpers
   - Performance tests

## Contributing

1. Check existing issues or create a new one
2. Fork the repository
3. Create a feature branch
4. Write tests for your changes
5. Implement your changes
6. Ensure all tests pass
7. Update documentation
8. Submit a pull request

### Pull Request Guidelines

1. Keep changes focused and atomic
2. Include tests for new functionality
3. Update relevant documentation
4. Follow the code style guidelines
5. Write clear commit messages
6. Respond to review comments

### Common Tasks

#### Adding a New Feature
1. Create an issue describing the feature
2. Create a new branch
3. Add tests in relevant package
4. Implement the feature
5. Update documentation
6. Submit pull request

#### Fixing a Bug
1. Create an issue if none exists
2. Create a branch
3. Add failing test that demonstrates the bug
4. Fix the bug
5. Verify all tests pass
6. Submit pull request

#### Adding Documentation
1. Update relevant .md files
2. Keep examples up to date
3. Run spell checker
4. Verify links work
5. Submit pull request

### Release Process

1. Update version in relevant files
2. Update CHANGELOG.md
3. Create release tag
4. Build release binaries
5. Create GitHub release
6. Update documentation

## Troubleshooting

### Common Issues

1. Build Errors
   - Ensure Go version is 1.21+
   - Run `go mod tidy`
   - Check for missing dependencies

2. Test Failures
   - Run with -v flag for verbose output
   - Check test dependencies
   - Verify test environment

3. Linting Issues
   - Run `make lint-fix`
   - Check .golangci.yml settings
   - Update code to follow style guide

### Getting Help

1. Check existing issues
2. Review documentation
3. Ask in discussions
4. Join developer chat
