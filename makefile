# Makefile for Large File Processing System
# Version: 1.0.0

PROJECT_NAME := largefile
VERSION := 1.0.0
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GO_VERSION := $(shell go version | awk '{print $$3}')

# Go configuration
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
CGO_ENABLED ?= 0

# Directories
BUILD_DIR := build
DIST_DIR := dist
COVERAGE_DIR := coverage
BENCH_DIR := bench
CONFIG_DIR := config

# Binary names
BINARY_NAME := largefile-copy

# CUE configuration
CUE := cue

# Go flags
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"
GOFLAGS := -mod=readonly -trimpath
TESTFLAGS := -v -race -count=1 -timeout=30m
BENCHFLAGS := -bench=. -benchmem -benchtime=10s

# Colors
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
NC := \033[0m

.PHONY: help
help:
	@echo "$(BLUE)Large File Processing System$(NC) - $(VERSION)"
	@echo ""
	@echo "$(GREEN)Targets:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

.PHONY: setup
setup: setup-cue setup-deps ## Setup development environment
	@echo "$(GREEN)✓ Setup complete$(NC)"

.PHONY: setup-cue
setup-cue:
	@echo "$(BLUE)Installing CUE...$(NC)"
	@go install cuelang.org/go/cmd/cue@latest
	@cue version
	@echo "$(GREEN)✓ CUE installed$(NC)"

.PHONY: setup-deps
setup-deps:
	@echo "$(BLUE)Installing dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)✓ Dependencies installed$(NC)"

.PHONY: generate
generate: validate-cue ## Generate code from CUE
	@echo "$(BLUE)Generating Go code from CUE...$(NC)"
	@cd $(CONFIG_DIR) && cue gen .
	@go fmt ./...
	@echo "$(GREEN)✓ Code generated$(NC)"

.PHONY: validate-cue
validate-cue: ## Validate all CUE configuration files
	@echo "$(BLUE)Validating CUE configurations...$(NC)"
	@cd $(CONFIG_DIR) && ./validate.sh
	@echo "$(GREEN)✓ CUE validation passed$(NC)"

.PHONY: fmt-cue
fmt-cue:
	@echo "$(BLUE)Formatting CUE files...$(NC)"
	@find . -name "*.cue" -type f -exec cue fmt {} \;
	@echo "$(GREEN)✓ CUE formatted$(NC)"

.PHONY: fmt
fmt: fmt-cue
	@echo "$(BLUE)Formatting Go code...$(NC)"
	@go fmt ./...
	@go mod tidy
	@echo "$(GREEN)✓ Code formatted$(NC)"

.PHONY: vet
vet:
	@echo "$(BLUE)Running go vet...$(NC)"
	@go vet ./...
	@echo "$(GREEN)✓ Go vet passed$(NC)"

.PHONY: lint
lint: fmt vet
	@echo "$(BLUE)Running golangci-lint...$(NC)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@golangci-lint run --timeout=5m ./...
	@echo "$(GREEN)✓ Linting passed$(NC)"

.PHONY: test
test: generate
	@echo "$(BLUE)Running tests...$(NC)"
	@mkdir -p $(COVERAGE_DIR)
	@go test $(TESTFLAGS) -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	@go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@go tool cover -func=$(COVERAGE_DIR)/coverage.out | tail -1
	@echo "$(GREEN)✓ Tests passed$(NC)"

.PHONY: bench
bench:
	@echo "$(BLUE)Running benchmarks...$(NC)"
	@mkdir -p $(BENCH_DIR)
	@go test $(BENCHFLAGS) -bench=. ./... | tee $(BENCH_DIR)/benchmark_$(shell date +%Y%m%d_%H%M%S).txt
	@echo "$(GREEN)✓ Benchmarks complete$(NC)"

.PHONY: build
build: generate
	@echo "$(BLUE)Building $(BINARY_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "$(GREEN)✓ Binary built: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"
	@ls -lh $(BUILD_DIR)/$(BINARY_NAME)

.PHONY: build-all
build-all: generate
	@echo "$(BLUE)Building for all platforms...$(NC)"
	@mkdir -p $(DIST_DIR)
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(GOFLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 .
	@GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build $(GOFLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 .
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(GOFLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 .
	@GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build $(GOFLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 .
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build $(GOFLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	@echo "$(GREEN)✓ Binaries built in $(DIST_DIR)/$(NC)"
	@ls -lh $(DIST_DIR)/

.PHONY: install
install: build
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)
	@echo "$(GREEN)✓ Installed to $(GOPATH)/bin/$(BINARY_NAME)$(NC)"

.PHONY: run
run: build
	@$(BUILD_DIR)/$(BINARY_NAME) -input test.dat -output output.dat -verbose

.PHONY: run-profile
run-profile: build
	@$(BUILD_DIR)/$(BINARY_NAME) -input test.dat -output output.dat -profile auto -verbose

.PHONY: clean
clean:
	@echo "$(BLUE)Cleaning...$(NC)"
	@rm -rf $(BUILD_DIR) $(DIST_DIR) $(COVERAGE_DIR) $(BENCH_DIR)
	@rm -f *.checkpoint *.tmp
	@go clean -cache -testcache
	@echo "$(GREEN)✓ Clean complete$(NC)"

.PHONY: clean-all
clean-all: clean
	@rm -rf vendor
	@go clean -modcache
	@echo "$(GREEN)✓ Full clean complete$(NC)"

.PHONY: release
release: clean test lint build-all
	@echo "$(BLUE)Creating release...$(NC)"
	@cd $(DIST_DIR) && sha256sum * > SHA256SUMS
	@echo "$(GREEN)✓ Release created in $(DIST_DIR)/$(NC)"

.PHONY: docker-build
docker-build:
	@echo "$(BLUE)Building Docker image...$(NC)"
	@docker build -t $(PROJECT_NAME):$(VERSION) -t $(PROJECT_NAME):latest .
	@echo "$(GREEN)✓ Docker image built$(NC)"

.PHONY: all
all: clean setup fmt lint test build
	@echo "$(GREEN)✓ Complete build successful$(NC)"

.PHONY: dev
dev: fmt lint test
	@echo "$(GREEN)✓ Development workflow complete$(NC)"

.PHONY: ci
ci: fmt-cue validate-cue generate fmt vet test build
	@echo "$(GREEN)✓ CI pipeline passed$(NC)"

.PHONY: version
version:
	@echo "$(BLUE)Version:$(NC) $(VERSION)"
	@echo "$(BLUE)Build Time:$(NC) $(BUILD_TIME)"
	@echo "$(BLUE)Git Commit:$(NC) $(GIT_COMMIT)"