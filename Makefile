.PHONY: test
test:
	@go test ./...

.PHONY: test-e2e
test-e2e: ## Run end-to-end tests (requires git and gh)
	@go test -tags e2e ./e2e/...

.PHONY: install-tools
install-tools: ## Install tools
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.5.0
	@go install golang.org/x/vuln/cmd/govulncheck@latest

check-tool-%:
	@which $* > /dev/null || (echo "Install $* with 'make install-tools'"; exit 1 )

.PHONY: lint
lint: check-tool-golangci-lint
	@golangci-lint run ./...

.PHONY: vulncheck
vulncheck: check-tool-govulncheck
	@govulncheck ./...

BIN_DIR ?= $(shell pwd)/bin
VERSION ?= dev
LDFLAGS := -ldflags "-X github.com/jcchavezs/gh-auth/internal/cli.version=$(VERSION)"

.PHONY: build
build:
	@mkdir -p $(BIN_DIR)
	@go build $(LDFLAGS) -o $(BIN_DIR) ./cmd/gh-auth

.PHONY: install
install:
	@BIN_DIR=$(shell go env GOPATH)/bin $(MAKE) build

.PHONY: generate
generate: ## Generate code
	@go generate ./...
