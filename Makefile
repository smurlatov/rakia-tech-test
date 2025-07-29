# Variables
APP_NAME=blog-api
DOCKER_IMAGE=blog-api:latest
PORT=8080

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GOFMT=gofmt

# Binary name
BINARY_NAME=blog-api

.PHONY: help build build-only run test test-coverage test-race clean deps fmt lint lint-fast docker-build docker-build-only docker-run docker-stop docker-logs

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: lint-fast test build-only ## Build the application (with lint + tests)

build-only: ## Build the application without tests/lint
	$(GOBUILD) -o $(BINARY_NAME) cmd/main.go

run: ## Run the application locally
	$(GOBUILD) -o $(BINARY_NAME) cmd/main.go
	PORT=$(PORT) ./$(BINARY_NAME)

test: ## Run all tests with verbose output
	$(GOTEST) -v ./...

test-coverage: ## Run tests with coverage report (generates coverage.html)
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-race: ## Run tests with race detection
	$(GOTEST) -v -race ./...

clean: ## Clean build artifacts and coverage files
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f coverage.out
	rm -f coverage.html

deps: ## Download and tidy Go module dependencies
	$(GOMOD) download
	$(GOMOD) tidy

fmt: ## Format Go code with gofmt
	$(GOFMT) -s -w .

docker-build: lint-fast test docker-build-only ## Build Docker image (with lint + tests)

docker-build-only: ## Build Docker image without tests/lint
	docker build -t $(DOCKER_IMAGE) .

docker-run: ## Run Docker container on port 8080
	docker run -d -p $(PORT):8080 --name $(APP_NAME) $(DOCKER_IMAGE)

docker-stop: ## Stop and remove Docker container
	-docker stop $(APP_NAME)
	-docker rm $(APP_NAME)

docker-logs: ## Show Docker container logs (live tail)
	docker logs -f $(APP_NAME)

lint: ## Lint Go code (auto-install if needed)
	@GOLANGCI_LINT_PATH="$(shell go env GOPATH)/bin/golangci-lint"; \
	if [ -x "$$GOLANGCI_LINT_PATH" ]; then \
		$$GOLANGCI_LINT_PATH run; \
	elif command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		$(shell go env GOPATH)/bin/golangci-lint run; \
	fi

lint-fast: ## Fast lint (essential checks only, used in build)
	@echo "Running linter..."
	@GOLANGCI_LINT_PATH="$(shell go env GOPATH)/bin/golangci-lint"; \
	if [ -x "$$GOLANGCI_LINT_PATH" ]; then \
		$$GOLANGCI_LINT_PATH run --fast || (echo "❌ Linting failed! Build aborted." && exit 1); \
	elif command -v golangci-lint > /dev/null; then \
		golangci-lint run --fast || (echo "❌ Linting failed! Build aborted." && exit 1); \
	else \
		echo "golangci-lint not installed. Using go vet..."; \
		go vet ./... || (echo "❌ Go vet failed! Build aborted." && exit 1); \
	fi
	@echo "✅ Linting passed!"
