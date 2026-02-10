################################################################################
# Variables & Configuration
################################################################################

# Project paths
ROOT_MODULE      := .
CMD_MODULE       := ./cmd

# CLI binary
CLI_BINARY       := example
CLI_MAIN         := $(CMD_MODULE)/main.go

################################################################################
# PHONY Targets
################################################################################
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  dev              - Run CLI in development mode"
	@echo "  test             - Run tests on all modules"
	@echo "  test-root        - Run tests on root module only"
	@echo "  test-cmd         - Run tests on cmd module only"
	@echo "  coverage         - Run tests with coverage report"
	@echo "  coverage-html    - Run coverage and open HTML report"
	@echo "  lint             - Run linter on all modules"
	@echo "  lint-root        - Run linter on root module only"
	@echo "  lint-cmd         - Run linter on cmd module only"
	@echo "  tidy             - Run go mod tidy on all modules"
	@echo "  tidy-root        - Run go mod tidy on root module only"
	@echo "  tidy-cmd         - Run go mod tidy on cmd module only"
	@echo "  download         - Download dependencies for all modules"
	@echo "  clean            - Clean build artifacts and test cache"

################################################################################
# Development
################################################################################

.PHONY: dev
dev:
	@echo "Running CLI in development mode..."
	cd $(CMD_MODULE) && go run $(notdir $(CLI_MAIN))

################################################################################
# Testing
################################################################################

.PHONY: test
test: test-root test-cmd

.PHONY: test-root
test-root:
	@echo "Running tests on root module..."
	go test ./pkg/ads/... -v

.PHONY: test-cmd
test-cmd:
	@echo "Running tests on cmd module..."
	cd $(CMD_MODULE) && go test ./... -v

.PHONY: coverage
coverage:
	@echo "Running tests with coverage (root module)..."
	@go test -coverprofile=/tmp/coverage.out ./pkg/ads/...
	@echo ""
	@echo "Coverage Summary:"
	@go tool cover -func=/tmp/coverage.out | tail -1
	@echo ""
	@echo "Detailed coverage report saved to /tmp/coverage.out"
	@echo "View HTML report: make coverage-html"

.PHONY: coverage-html
coverage-html: coverage
	@echo "Opening HTML coverage report..."
	@go tool cover -html=/tmp/coverage.out

################################################################################
# Linting
################################################################################

.PHONY: lint
lint: lint-root lint-cmd

.PHONY: lint-root
lint-root:
	@echo "Running linter on root module..."
	golangci-lint run ./...

.PHONY: lint-cmd
lint-cmd:
	@echo "Running linter on cmd module..."
	cd $(CMD_MODULE) && golangci-lint run ./...

################################################################################
# Dependency Management
################################################################################

.PHONY: tidy
tidy: tidy-root tidy-cmd

.PHONY: tidy-root
tidy-root:
	@echo "Running go mod tidy on root module..."
	go mod tidy

.PHONY: tidy-cmd
tidy-cmd:
	@echo "Running go mod tidy on cmd module..."
	cd $(CMD_MODULE) && go mod tidy

.PHONY: download
download:
	@echo "Downloading dependencies for root module..."
	go mod download
	@echo "Downloading dependencies for cmd module..."
	cd $(CMD_MODULE) && go mod download

################################################################################
# Cleaning
################################################################################

.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	go clean
	cd $(CMD_MODULE) && go clean
	rm -f $(CMD_MODULE)/$(CLI_BINARY)
	@echo "Cleaning test cache..."
	go clean -testcache
	@echo "Clean complete."

################################################################################
# Default Goal
################################################################################
.DEFAULT_GOAL := help
