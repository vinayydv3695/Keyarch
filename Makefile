# Makefile for Keyarch

.PHONY: build run clean install test lint fmt help

APP_NAME := keyarch
BUILD_DIR := build
MAIN_PATH := ./cmd/keyarch
INSTALL_PATH := /usr/local/bin

# Build flags
LDFLAGS := -s -w
GOFLAGS := -trimpath

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building $(APP_NAME)..."
	@go build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

run: ## Run the application
	@go run $(MAIN_PATH)

install: build ## Install the application to system
	@echo "Installing $(APP_NAME) to $(INSTALL_PATH)..."
	@sudo cp $(BUILD_DIR)/$(APP_NAME) $(INSTALL_PATH)/
	@echo "Installation complete!"

uninstall: ## Uninstall the application
	@echo "Uninstalling $(APP_NAME)..."
	@sudo rm -f $(INSTALL_PATH)/$(APP_NAME)
	@echo "Uninstall complete!"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@go clean
	@echo "Clean complete!"

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	@golangci-lint run ./...

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@gofmt -s -w .

tidy: ## Tidy dependencies
	@echo "Tidying dependencies..."
	@go mod tidy
	@go mod verify

dev: ## Run in development mode
	@go run $(MAIN_PATH)

# Cross-compilation targets
build-all: build-linux build-darwin build-windows ## Build for all platforms

build-linux: ## Build for Linux (amd64)
	@echo "Building for Linux (amd64)..."
	@GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 $(MAIN_PATH)

build-darwin: ## Build for macOS (amd64 and arm64)
	@echo "Building for macOS (amd64)..."
	@GOOS=darwin GOARCH=amd64 go build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 $(MAIN_PATH)
	@echo "Building for macOS (arm64)..."
	@GOOS=darwin GOARCH=arm64 go build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 $(MAIN_PATH)

build-windows: ## Build for Windows (amd64)
	@echo "Building for Windows (amd64)..."
	@GOOS=windows GOARCH=amd64 go build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe $(MAIN_PATH)

release: clean build-all ## Create release archives
	@echo "Creating release archives..."
	@mkdir -p $(BUILD_DIR)/release
	@cd $(BUILD_DIR) && tar -czf release/$(APP_NAME)-linux-amd64.tar.gz $(APP_NAME)-linux-amd64
	@cd $(BUILD_DIR) && tar -czf release/$(APP_NAME)-darwin-amd64.tar.gz $(APP_NAME)-darwin-amd64
	@cd $(BUILD_DIR) && tar -czf release/$(APP_NAME)-darwin-arm64.tar.gz $(APP_NAME)-darwin-arm64
	@cd $(BUILD_DIR) && zip -q release/$(APP_NAME)-windows-amd64.zip $(APP_NAME)-windows-amd64.exe
	@echo "Release archives created in $(BUILD_DIR)/release/"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(APP_NAME):latest .

docker-run: ## Run Docker container
	@docker run -it --rm $(APP_NAME):latest

.DEFAULT_GOAL := help
