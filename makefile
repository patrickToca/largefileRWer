# Makefile for Large File Processing System
# Version: 1.0.0

PROJECT_NAME := largefile-copy
VERSION := 1.0.0
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Go configuration
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
CGO_ENABLED ?= 0

# Directories
BUILD_DIR := build
DIST_DIR := dist
CONFIG_DIR := config

# Binary names
BINARY_NAME := largefile-copy

# Go flags
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"
GOFLAGS := -mod=readonly -trimpath
TESTFLAGS := -v -race -count=1 -timeout=30m

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
setup: ## Setup development environment
	@echo "$(BLUE)Installing dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)✓ Setup complete$(NC)"

.PHONY: validate-cue
validate-cue: ## Validate all CUE configuration files
	@echo "$(BLUE)Validating CUE configurations...$(NC)"
	@cd $(CONFIG_DIR) && ./validate.sh
	@echo "$(GREEN)✓ CUE validation passed$(NC)"

.PHONY: generate
generate: validate-cue ## Generate code (currently no generation needed, just validation)
	@echo "$(BLUE)Checking generated config...$(NC)"
	@test -f $(CONFIG_DIR)/config_gen.go || { echo "$(RED)config_gen.go not found!$(NC)"; exit 1; }
	@echo "$(GREEN)✓ Generated code verified$(NC)"

.PHONY: fmt
fmt: ## Format code
	@echo "$(BLUE)Formatting code...$(NC)"
	@go fmt ./...
	@go mod tidy
	@echo "$(GREEN)✓ Code formatted$(NC)"

.PHONY: vet
vet: ## Run go vet
	@echo "$(BLUE)Running go vet...$(NC)"
	@go vet ./...
	@echo "$(GREEN)✓ Go vet passed$(NC)"

.PHONY: test
test: validate-cue ## Run tests
	@echo "$(BLUE)Running tests...$(NC)"
	@go test $(TESTFLAGS) ./...
	@echo "$(GREEN)✓ Tests passed$(NC)"

.PHONY: build
build: validate-cue fmt vet ## Build the application
	@echo "$(BLUE)Building $(BINARY_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "$(GREEN)✓ Binary built: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"
	@ls -lh $(BUILD_DIR)/$(BINARY_NAME)

.PHONY: build-quick
build-quick: validate-cue ## Quick build without vet
	@echo "$(BLUE)Quick building $(BINARY_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "$(GREEN)✓ Binary built: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

.PHONY: clean
clean: ## Clean build artifacts
	@echo "$(BLUE)Cleaning...$(NC)"
	@rm -rf $(BUILD_DIR) $(DIST_DIR)
	@go clean -cache -testcache
	@echo "$(GREEN)✓ Clean complete$(NC)"

.PHONY: all
all: clean setup validate-cue fmt vet test build ## Full build
	@echo "$(GREEN)✓ Complete build successful$(NC)"

.PHONY: run
run: build ## Run with default configuration
	@$(BUILD_DIR)/$(BINARY_NAME) -input test.dat -output output.dat -verbose

.PHONY: run-profile
run-profile: build ## Run with auto-detected profile
	@$(BUILD_DIR)/$(BINARY_NAME) -input test.dat -output output.dat -profile auto -verbose

.PHONY: run-nvme
run-nvme: build ## Run with NVMe profile
	@$(BUILD_DIR)/$(BINARY_NAME) -input test.dat -output output.dat -profile nvme -verbose

.PHONY: dev
dev: fmt vet test build ## Development workflow
	@echo "$(GREEN)✓ Development build complete$(NC)"

.PHONY: version
version: ## Show version information
	@echo "$(BLUE)Version:$(NC) $(VERSION)"
	@echo "$(BLUE)Build Time:$(NC) $(BUILD_TIME)"
	@echo "$(BLUE)Git Commit:$(NC) $(GIT_COMMIT)"