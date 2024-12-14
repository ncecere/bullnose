# Contributing to Bullnose

First off, thank you for considering contributing to Bullnose! It's people like you that make Bullnose such a great tool.

## Code of Conduct

This project and everyone participating in it is governed by our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the issue list as you might find out that you don't need to create one. When you are creating a bug report, please include as many details as possible:

* Use a clear and descriptive title
* Describe the exact steps which reproduce the problem
* Provide specific examples to demonstrate the steps
* Describe the behavior you observed after following the steps
* Explain which behavior you expected to see instead and why
* Include logs if relevant
* Include your Go version and OS

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, please include:

* A clear and descriptive title
* A detailed description of the proposed functionality
* Explain why this enhancement would be useful
* List any alternative solutions you've considered
* Include mockups or examples if applicable

### Pull Requests

* Fork the repo and create your branch from `main`
* If you've added code that should be tested, add tests
* Ensure the test suite passes
* Make sure your code lints (`make lint`)
* Update the documentation if needed
* Update the CHANGELOG.md

## Development Setup

1. Fork and clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Run tests:
   ```bash
   make test
   ```
4. Create a branch for your changes:
   ```bash
   git checkout -b feature/amazing-feature
   ```

## Project Structure

```
bullnose/
├── cmd/                    # Command line interface
│   └── bullnose/          # Main application
├── internal/              # Internal packages
│   ├── config/           # Configuration management
│   └── scraper/         # Web scraping logic
├── .github/              # GitHub specific files
├── bin/                  # Compiled binaries
└── test/                 # Test files
```

## Coding Style

* Follow standard Go formatting (`go fmt`)
* Use meaningful variable and function names
* Write comments for non-obvious code
* Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines

## Testing

* Write unit tests for new code
* Ensure all tests pass before submitting PR
* Include integration tests for new features
* Run tests with:
  ```bash
  make test
  make coverage
  ```

## Documentation

* Update README.md for user-facing changes
* Document new features in code
* Update CHANGELOG.md
* Keep documentation clear and concise

## Commit Messages

* Use the present tense ("Add feature" not "Added feature")
* Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
* Limit the first line to 72 characters or less
* Reference issues and pull requests liberally after the first line

## Review Process

1. Create a Pull Request with a clear title and description
2. Wait for review from maintainers
3. Make requested changes if any
4. Once approved, your PR will be merged

## Release Process

1. Update CHANGELOG.md
2. Create a new tag following semver
3. Push tag to trigger release workflow
4. GitHub Actions will build and publish releases

## Questions?

Feel free to create an issue labeled 'question' if you need help or clarification.

Thank you for contributing to Bullnose! 🎉
