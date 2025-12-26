# Build configuration
APP_NAME := container-network
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DIR := build
CMD_PATH := ./cmd/$(APP_NAME)
LDFLAGS := -ldflags "-X $(APP_NAME)/pkg/config.Version=$(VERSION) -s -w"

# Default target
.PHONY: all
all: clean build

# Build for current platform
.PHONY: build
build:
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(CMD_PATH)

# Build for all platforms
.PHONY: build-all
build-all: build-linux-amd64 build-linux-386 build-darwin-amd64 build-darwin-arm64

# Linux AMD64
.PHONY: build-linux-amd64
build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 $(CMD_PATH)

# Linux x86 (386)
.PHONY: build-linux-386
build-linux-386:
	GOOS=linux GOARCH=386 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux-386 $(CMD_PATH)

# macOS AMD64 (Intel)
.PHONY: build-darwin-amd64
build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 $(CMD_PATH)

# macOS ARM64 (Apple Silicon)
.PHONY: build-darwin-arm64
build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 $(CMD_PATH)

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

# Run the application
.PHONY: run
run:
	go run $(CMD_PATH)

# Run tests
.PHONY: test
test:
	go test -v ./...

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Lint code
.PHONY: lint
lint:
	go vet ./...

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build            - Build for current platform"
	@echo "  build-all        - Build for all platforms"
	@echo "  build-linux-amd64 - Build for Linux AMD64"
	@echo "  build-linux-386   - Build for Linux x86 (32-bit)"
	@echo "  build-darwin-amd64 - Build for macOS Intel"
	@echo "  build-darwin-arm64 - Build for macOS Apple Silicon"
	@echo "  clean            - Remove build artifacts"
	@echo "  run              - Run the application"
	@echo "  test             - Run tests"
	@echo "  fmt              - Format code"
	@echo "  lint             - Lint code"
	@echo ""
	@echo "Variables:"
	@echo "  VERSION          - Set version (default: git tag or 'dev')"
	@echo ""
	@echo "Example:"
	@echo "  make build-all VERSION=1.0.0"
