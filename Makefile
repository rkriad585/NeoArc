APP_NAME    := neoarc
VERSION     := $(shell cat .version 2>/dev/null || echo "v0.0.0")
COMMIT      := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME  := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
CLI_DIR     := neoarc-cli
SERVER_DIR  := neoarc-server
OUTPUT_DIR  := bin

LDFLAGS := -s -w \
	-X neoarc/internal/cli.Version=$(VERSION) \
	-X neoarc/internal/cli.Commit=$(COMMIT) \
	-X neoarc/internal/cli.BuildTime=$(BUILD_TIME) \
	-X neoarc/internal/version.Version=$(VERSION) \
	-X neoarc/internal/version.Commit=$(COMMIT) \
	-X neoarc/internal/version.BuildTime=$(BUILD_TIME)

.PHONY: help build run clean test lint format install release all

help:
	@echo "$(APP_NAME) $(VERSION) — Makefile"
	@echo ""
	@echo "Targets:"
	@echo "  make build         Build CLI for current platform"
	@echo "  make build-all     Cross-compile for all 6 platforms"
	@echo "  make run           Run the CLI (alias: neoarc help)"
	@echo "  make clean         Remove build artifacts"
	@echo "  make test          Run all tests"
	@echo "  make lint          Lint Go and Python code"
	@echo "  make format        Format Go and Python code"
	@echo "  make install       Install CLI locally (to ~/.config/...)"
	@echo "  make release       Tag and prepare a release"
	@echo "  make docker        Build Docker image for CLI"
	@echo "  make docker-server Build Docker image for server"

build:
	@mkdir -p $(OUTPUT_DIR)
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(OUTPUT_DIR)/$(APP_NAME) $(CLI_DIR)/cmd/neoarc/

build-all:
	@$(SHELL) build.sh

run:
	@go run $(CLI_DIR)/cmd/neoarc/

clean:
	@rm -rf $(OUTPUT_DIR)
	@rm -f $(CLI_DIR)/$(APP_NAME)
	@rm -f $(CLI_DIR)/$(APP_NAME).exe
	@echo "Cleaned build artifacts."

test:
	@echo "--- CLI unit tests ---"
	cd $(CLI_DIR) && go test ./internal/... -v -count=1
	@echo ""
	@echo "--- Server tests ---"
	cd $(SERVER_DIR) && python -m pytest tests/ -v 2>/dev/null || echo "Server tests require 'pip install -e \".[test]\"' first"

lint:
	@echo "--- Go lint ---"
	cd $(CLI_DIR) && go vet ./...
	@echo ""
	@echo "--- Python lint ---"
	cd $(SERVER_DIR) && python -m flake8 . --count --select=E9,F63,F7,F82 --show-source --statistics 2>/dev/null || echo "flake8 not installed"

format:
	@echo "--- Go format ---"
	cd $(CLI_DIR) && gofmt -l -w .
	@echo "--- Python format ---"
	cd $(SERVER_DIR) && python -m black --check . 2>/dev/null && python -m black . 2>/dev/null || echo "black not installed"

install:
	@echo ">>> Installing $(APP_NAME) $(VERSION)"
	@mkdir -p $(OUTPUT_DIR)
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(OUTPUT_DIR)/$(APP_NAME) $(CLI_DIR)/cmd/neoarc/
	$(OUTPUT_DIR)/$(APP_NAME) --install

release:
	@echo "Current version: $(VERSION)"
	@echo ""
	@echo "Steps to release:"
	@echo "  1. Update .version if needed:  echo \"vX.Y.Z\" > .version"
	@echo "  2. Commit:                     git add .version && git commit -m \"Release \$$(cat .version)\""
	@echo "  3. Tag:                        git tag \$$(cat .version)"
	@echo "  4. Push:                       git push --tags"
	@echo "  5. CI builds all binaries and publishes the release."

docker:
	@docker build -t $(APP_NAME)-cli \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		-f Dockerfile .

docker-server:
	@cd $(SERVER_DIR) && docker build -t $(APP_NAME)-server .

all: build test
