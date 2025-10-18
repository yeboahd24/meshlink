# MeshLink Build System

.PHONY: build clean test run-broadcaster run-viewer docker-build docker-run deps

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary names
BROADCASTER_BINARY=broadcaster
VIEWER_BINARY=viewer

# Build directories
BUILD_DIR=build
DIST_DIR=dist

all: build

# Install dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Build applications (headless versions - no GUI dependencies)
build: deps
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) -buildvcs=false -o $(BUILD_DIR)/$(BROADCASTER_BINARY) ./cmd/broadcaster-headless
	$(GOBUILD) -buildvcs=false -o $(BUILD_DIR)/$(VIEWER_BINARY) ./cmd/viewer-headless

# Build GUI versions for current platform (requires X11/OpenGL)
build-gui: deps
	mkdir -p $(DIST_DIR)
	$(GOBUILD) -buildvcs=false -o $(DIST_DIR)/$(BROADCASTER_BINARY)-gui ./cmd/broadcaster
	$(GOBUILD) -buildvcs=false -o $(DIST_DIR)/$(VIEWER_BINARY)-gui ./cmd/viewer

# Build for multiple platforms (headless versions for cross-compilation)
build-all: deps
	mkdir -p $(DIST_DIR)
	# Windows
	GOOS=windows GOARCH=amd64 $(GOBUILD) -buildvcs=false -o $(DIST_DIR)/$(BROADCASTER_BINARY)-windows-amd64.exe ./cmd/broadcaster-headless
	GOOS=windows GOARCH=amd64 $(GOBUILD) -buildvcs=false -o $(DIST_DIR)/$(VIEWER_BINARY)-windows-amd64.exe ./cmd/viewer-headless
	# macOS
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -buildvcs=false -o $(DIST_DIR)/$(BROADCASTER_BINARY)-darwin-amd64 ./cmd/broadcaster-headless
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -buildvcs=false -o $(DIST_DIR)/$(VIEWER_BINARY)-darwin-amd64 ./cmd/viewer-headless
	# Linux
	GOOS=linux GOARCH=amd64 $(GOBUILD) -buildvcs=false -o $(DIST_DIR)/$(BROADCASTER_BINARY)-linux-amd64 ./cmd/broadcaster-headless
	GOOS=linux GOARCH=amd64 $(GOBUILD) -buildvcs=false -o $(DIST_DIR)/$(VIEWER_BINARY)-linux-amd64 ./cmd/viewer-headless
	# ARM (Raspberry Pi)
	GOOS=linux GOARCH=arm GOARM=7 $(GOBUILD) -buildvcs=false -o $(DIST_DIR)/$(BROADCASTER_BINARY)-linux-arm7 ./cmd/broadcaster-headless
	GOOS=linux GOARCH=arm GOARM=7 $(GOBUILD) -buildvcs=false -o $(DIST_DIR)/$(VIEWER_BINARY)-linux-arm7 ./cmd/viewer-headless

# Run applications
run-broadcaster: build
	./$(BUILD_DIR)/$(BROADCASTER_BINARY)

run-viewer: build
	./$(BUILD_DIR)/$(VIEWER_BINARY)

# Testing
test:
	$(GOTEST) -v ./...

test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

# Docker operations
docker-build:
	docker build -f deployments/Dockerfile.broadcaster -t meshlink/broadcaster:latest .

docker-run:
	docker compose -f deployments/docker-compose.yml up

# Docker development environment
docker-dev:
	docker compose -f deployments/docker-compose.dev.yml up -d dev-env

docker-dev-shell:
	docker compose -f deployments/docker-compose.dev.yml exec dev-env bash

docker-dev-build:
	docker compose -f deployments/docker-compose.dev.yml exec dev-env make build

docker-dev-test:
	docker compose -f deployments/docker-compose.dev.yml exec dev-env make test

docker-dev-run-broadcaster:
	docker compose -f deployments/docker-compose.dev.yml exec dev-env make dev-broadcaster

docker-dev-run-viewer:
	docker compose -f deployments/docker-compose.dev.yml exec dev-env make dev-viewer

docker-dev-stop:
	docker compose -f deployments/docker-compose.dev.yml down

# Full demo environment (dev + production)
docker-full-demo:
	docker compose -f deployments/docker-compose.dev.yml up --build

# Development (headless)
dev-broadcaster:
	$(GOCMD) run ./cmd/broadcaster-headless

dev-viewer:
	$(GOCMD) run ./cmd/viewer-headless

# Development (GUI - requires X11/OpenGL)
dev-broadcaster-gui:
	$(GOCMD) run ./cmd/broadcaster

dev-viewer-gui:
	$(GOCMD) run ./cmd/viewer

# Cleanup
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -rf $(DIST_DIR)
	rm -f coverage.out

# Generate default config
config:
	@echo "Generating default configuration..."
	@echo '{"network":{"port":8080,"discovery_key":"meshlink-church","max_peers":50},"media":{"video_codec":"h264","audio_codec":"aac","bitrate":2000,"resolution":"1280x720","frame_rate":30},"ui":{"theme":"dark","fullscreen":false,"show_stats":true}}' > config.json

# Install system dependencies (Ubuntu/Debian)
install-deps-ubuntu:
	sudo apt-get update
	sudo apt-get install -y libgstreamer1.0-dev libgstreamer-plugins-base1.0-dev

# Install system dependencies (macOS)
install-deps-macos:
	brew install gstreamer gst-plugins-base

# Help
help:
	@echo "Available targets:"
	@echo "  Local Development:"
	@echo "    build          - Build applications"
	@echo "    build-all      - Build for all platforms"
	@echo "    run-broadcaster - Run broadcaster application"
	@echo "    run-viewer     - Run viewer application"
	@echo "    test           - Run tests"
	@echo ""
	@echo "  Docker Development:"
	@echo "    docker-dev     - Start dev environment"
	@echo "    docker-dev-shell - Enter dev container"
	@echo "    docker-dev-build - Build in container"
	@echo "    docker-dev-test  - Test in container"
	@echo "    docker-dev-run-broadcaster - Run broadcaster in container"
	@echo "    docker-dev-run-viewer - Run viewer in container"
	@echo "    docker-dev-stop - Stop dev environment"
	@echo ""
	@echo "  Production Demo:"
	@echo "    docker-run     - Run production demo"
	@echo "    docker-full-demo - Run full demo (dev + prod)"
	@echo ""
	@echo "  Utilities:"
	@echo "    clean          - Clean build artifacts"
	@echo "    config         - Generate default config"
	@echo "    help           - Show this help"