# Makefile for K8sMed

# Variables
NAME := k8smed
VERSION := 0.1.0-alpha
LDFLAGS := -ldflags "-X main.version=$(VERSION)"
BINDIR := bin
BIN := kubectl-$(NAME)
INSTALLDIR := /usr/local/bin
GO := go
GOFLAGS :=
PKG := github.com/k8smed/k8smed
SOURCES := $(shell find . -name "*.go" -type f)
GOFMT ?= gofmt -s
GOFILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")
TESTPKGS := $(shell go list ./... | grep -v vendor)
COVERAGE_DIR := coverage

# Detect OS and architecture
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# Build targets
.PHONY: all build clean install uninstall fmt lint vet test help tools coverage integration-test docs docker

all: build

build: $(BINDIR)/$(BIN)

$(BINDIR)/$(BIN): $(SOURCES)
	@mkdir -p $(BINDIR)
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BINDIR)/$(BIN) $(PKG)/cmd/kubectl-k8smed

install: build
	cp $(BINDIR)/$(BIN) $(INSTALLDIR)/$(BIN)
	@echo "Installed $(BIN) to $(INSTALLDIR)/$(BIN)"
	@echo "You can now use 'kubectl $(NAME)' command"

uninstall:
	rm -f $(INSTALLDIR)/$(BIN)
	@echo "Uninstalled $(BIN) from $(INSTALLDIR)/$(BIN)"

clean:
	rm -rf $(BINDIR)
	rm -rf $(COVERAGE_DIR)

fmt:
	$(GOFMT) -w $(GOFILES)

lint:
	@which golint > /dev/null; if [ $$? -ne 0 ]; then \
		$(GO) install golang.org/x/lint/golint@latest; \
	fi
	golint -set_exit_status $(TESTPKGS)

vet:
	$(GO) vet $(TESTPKGS)

test:
	$(GO) test -v $(TESTPKGS)

run: build
	$(BINDIR)/$(BIN)

# Tools installation
tools:
	@echo "Installing development tools..."
	go install golang.org/x/lint/golint@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/goreleaser/goreleaser@latest
	go install github.com/mikefarah/yq/v4@latest
	@echo "Development tools installed successfully!"

# Integration testing
integration-test:
	@echo "Running integration tests..."
	go test -v -tags=integration ./tests/...

# Code coverage
coverage:
	@mkdir -p $(COVERAGE_DIR)
	go test -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "Coverage report generated at $(COVERAGE_DIR)/coverage.html"

# Generate documentation
docs:
	@echo "Generating documentation..."
	@mkdir -p docs/api
	go doc -all $(PKG)/pkg > docs/api/API.md
	@echo "API documentation generated at docs/api/API.md"

# Docker build
docker:
	docker build -t k8smed/k8smed:$(VERSION) -t k8smed/k8smed:latest .
	@echo "Docker image built: k8smed/k8smed:$(VERSION)"

# Full development setup (run this first after cloning)
dev-setup: tools

# Run all checks
check: fmt lint vet test

# Cross-compilation
.PHONY: cross-build
cross-build:
	@mkdir -p $(BINDIR)/$(GOOS)_$(GOARCH)
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BINDIR)/$(GOOS)_$(GOARCH)/$(BIN) $(PKG)/cmd/kubectl-k8smed

# Build for multiple platforms
.PHONY: release
release:
	@mkdir -p $(BINDIR)/release
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BINDIR)/release/$(BIN)_linux_amd64 $(PKG)/cmd/kubectl-k8smed
	GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BINDIR)/release/$(BIN)_darwin_amd64 $(PKG)/cmd/kubectl-k8smed
	GOOS=darwin GOARCH=arm64 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BINDIR)/release/$(BIN)_darwin_arm64 $(PKG)/cmd/kubectl-k8smed
	GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BINDIR)/release/$(BIN)_windows_amd64.exe $(PKG)/cmd/kubectl-k8smed

help:
	@echo "Available targets:"
	@echo "  all           - Build the kubectl plugin (default)"
	@echo "  build         - Build the kubectl plugin"
	@echo "  install       - Install the kubectl plugin to $(INSTALLDIR)"
	@echo "  uninstall     - Remove the kubectl plugin from $(INSTALLDIR)"
	@echo "  clean         - Remove build artifacts"
	@echo "  fmt           - Format Go source code"
	@echo "  lint          - Run golint"
	@echo "  vet           - Run go vet"
	@echo "  test          - Run tests"
	@echo "  tools         - Install development tools"
	@echo "  integration-test - Run integration tests"
	@echo "  coverage      - Generate test coverage report"
	@echo "  docs          - Generate documentation"
	@echo "  docker        - Build Docker image"
	@echo "  dev-setup     - Set up development environment"
	@echo "  check         - Run all code quality checks"
	@echo "  release       - Build for multiple platforms"
	@echo "  cross-build   - Build for a specific OS/arch (use GOOS and GOARCH variables)"
	@echo "  help          - Show this help" 