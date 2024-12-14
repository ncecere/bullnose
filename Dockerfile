# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache make git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN make build

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN adduser -D -h /app bullnose

WORKDIR /app
USER bullnose

# Copy binary from builder
COPY --from=builder /app/build/bullnose /usr/local/bin/

# Create directory for output
RUN mkdir -p /app/scraped-content

# Set default config location
ENV CONFIG_FILE=/app/config/bullnose.yaml

# Document the port (if needed in future)
# EXPOSE 8080

# Set entrypoint
ENTRYPOINT ["bullnose"]

# Default command (can be overridden)
CMD ["--help"]

# Labels
LABEL org.opencontainers.image.source="https://github.com/ncecere/bullnose"
LABEL org.opencontainers.image.description="A web scraper that converts web pages to clean markdown"
LABEL org.opencontainers.image.licenses="MIT"
