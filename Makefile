# Variables
BINARY_NAME=streaming-platform
DOCKER_IMAGE=go-streaming-platform
DOCKER_TAG=latest
MAIN_PATH=cmd/api/main.go

# Go related variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_UNIX=$(BINARY_NAME)_unix

# Build the application
.PHONY: build
build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)

# Build for multiple platforms
.PHONY: build-all
build-all:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)_linux_amd64 $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)_windows_amd64.exe $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)_darwin_amd64 $(MAIN_PATH)

# Clean build artifacts
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)_*

# Run tests
.PHONY: test
test:
	$(GOTEST) -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Run the application
.PHONY: run
run:
	$(GOCMD) run $(MAIN_PATH)

# Run with hot reload (requires air)
.PHONY: dev
dev:
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not found. Installing air..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

# Install dependencies
.PHONY: deps
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Update dependencies
.PHONY: deps-update
deps-update:
	$(GOMOD) get -u ./...
	$(GOMOD) tidy

# Format code
.PHONY: fmt
fmt:
	$(GOCMD) fmt ./...

# Run linter
.PHONY: lint
lint:
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run; \
	fi

# Generate documentation
.PHONY: docs
docs:
	@if command -v swag > /dev/null; then \
		swag init -g $(MAIN_PATH); \
	else \
		echo "swag not found. Installing..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
		swag init -g $(MAIN_PATH); \
	fi

# Docker commands
.PHONY: docker-build
docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

.PHONY: docker-run
docker-run:
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE):$(DOCKER_TAG)

.PHONY: docker-stop
docker-stop:
	docker stop $$(docker ps -q --filter ancestor=$(DOCKER_IMAGE):$(DOCKER_TAG))

.PHONY: docker-clean
docker-clean:
	docker rmi $(DOCKER_IMAGE):$(DOCKER_TAG)

# Docker Compose commands
.PHONY: compose-up
compose-up:
	docker-compose up -d

.PHONY: compose-down
compose-down:
	docker-compose down

.PHONY: compose-logs
compose-logs:
	docker-compose logs -f

# Database commands
.PHONY: db-migrate
db-migrate:
	@echo "Database migration will be handled automatically by GORM"

.PHONY: db-seed
db-seed:
	@echo "Add database seeding commands here if needed"

# Security checks
.PHONY: security-check
security-check:
	@if command -v gosec > /dev/null; then \
		gosec ./...; \
	else \
		echo "gosec not found. Installing..."; \
		go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; \
		gosec ./...; \
	fi

# Performance testing
.PHONY: bench
bench:
	$(GOTEST) -bench=. -benchmem ./...

# Install development tools
.PHONY: install-tools
install-tools:
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

# Check for vulnerabilities
.PHONY: vuln-check
vuln-check:
	$(GOCMD) list -json -deps ./... | nancy sleuth

# Create release
.PHONY: release
release: clean build-all
	@echo "Creating release..."
	@mkdir -p release
	@cp $(BINARY_NAME)_* release/
	@cp README.md LICENSE release/
	@echo "Release created in ./release/"

# Help
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  build          - Build the application"
	@echo "  build-all      - Build for multiple platforms"
	@echo "  clean          - Clean build artifacts"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  run            - Run the application"
	@echo "  dev            - Run with hot reload (requires air)"
	@echo "  deps           - Install dependencies"
	@echo "  deps-update    - Update dependencies"
	@echo "  fmt            - Format code"
	@echo "  lint           - Run linter"
	@echo "  docs           - Generate documentation"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  docker-stop    - Stop Docker container"
	@echo "  docker-clean   - Remove Docker image"
	@echo "  compose-up     - Start with Docker Compose"
	@echo "  compose-down   - Stop Docker Compose"
	@echo "  compose-logs   - View Docker Compose logs"
	@echo "  security-check - Run security checks"
	@echo "  bench          - Run benchmarks"
	@echo "  install-tools  - Install development tools"
	@echo "  vuln-check     - Check for vulnerabilities"
	@echo "  release        - Create release package"
	@echo "  help           - Show this help message"

# Default target
.DEFAULT_GOAL := help 