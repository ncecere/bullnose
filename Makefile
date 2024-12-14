BINARY_NAME=bullnose
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=${VERSION}"

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOSRC=./cmd/bullnose

# OS/ARCH specific builds
PLATFORMS=linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

# Ensure GOPATH
GOPATH?=$(shell go env GOPATH)

.PHONY: all build clean test coverage lint vet fmt help cross-build install uninstall

all: clean build test ## Build and run tests

build: ## Build the binary
	@echo "Building ${BINARY_NAME}..."
	@go build ${LDFLAGS} -o ${GOBIN}/${BINARY_NAME} ${GOSRC}

clean: ## Remove build artifacts
	@echo "Cleaning..."
	@rm -rf ${GOBIN}
	@rm -rf coverage.out

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

fmt: ## Run go fmt
	@echo "Running go fmt..."
	@go fmt ./...

cross-build: ## Build for all platforms
	@echo "Building for multiple platforms..."
	@mkdir -p ${GOBIN}
	@for platform in ${PLATFORMS}; do \
		GOOS=$${platform%/*} \
		GOARCH=$${platform#*/} \
		OUTPUT="${GOBIN}/${BINARY_NAME}-$${GOOS}-$${GOARCH}" ; \
		if [ "$${GOOS}" = "windows" ]; then \
			OUTPUT="$${OUTPUT}.exe" ; \
		fi ; \
		echo "Building $${OUTPUT}" ; \
		GOOS=$${GOOS} GOARCH=$${GOARCH} go build ${LDFLAGS} -o "$${OUTPUT}" ${GOSRC} || exit 1; \
	done

install: build ## Install binary to GOPATH/bin
	@echo "Installing to ${GOPATH}/bin"
	@cp ${GOBIN}/${BINARY_NAME} ${GOPATH}/bin/

uninstall: ## Remove binary from GOPATH/bin
	@echo "Removing from ${GOPATH}/bin"
	@rm -f ${GOPATH}/bin/${BINARY_NAME}

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# Default target
.DEFAULT_GOAL := help
