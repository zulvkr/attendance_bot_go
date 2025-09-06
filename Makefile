# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary names
BINARY_NAME=attendance-bot
SETUP_BINARY=setup-totp

# Build targets
.PHONY: all build clean test deps setup-totp docker

all: deps build

build:
	$(GOBUILD) -o $(BINARY_NAME) ./cmd/bot
	$(GOBUILD) -o $(SETUP_BINARY) ./cmd/setup-totp

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(SETUP_BINARY)

test:
	$(GOTEST) -v ./...

test-coverage:
	$(GOTEST) -cover ./...

deps:
	$(GOMOD) download
	$(GOMOD) tidy

setup-totp:
	$(GOBUILD) -o $(SETUP_BINARY) ./cmd/setup-totp
	./$(SETUP_BINARY)

# Docker targets
docker-build:
	docker build -t attendance-bot-go .

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Development
dev: build
	./$(BINARY_NAME)

# Production build with optimizations
build-prod:
	CGO_ENABLED=0 GOOS=linux $(GOBUILD) -a -installsuffix cgo -ldflags '-w -s' -o $(BINARY_NAME) ./cmd/bot

# Cross-compilation
build-windows:
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME).exe ./cmd/bot

build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) ./cmd/bot

build-darwin:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) ./cmd/bot

# Help
help:
	@echo "Available targets:"
	@echo "  build         - Build the bot and setup utility"
	@echo "  clean         - Clean build artifacts"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  deps          - Download and tidy dependencies"
	@echo "  setup-totp    - Build and run TOTP setup utility"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run with docker-compose"
	@echo "  docker-stop   - Stop docker-compose"
	@echo "  docker-logs   - View docker logs"
	@echo "  dev           - Build and run for development"
	@echo "  build-prod    - Production build with optimizations"
	@echo "  build-windows - Cross-compile for Windows"
	@echo "  build-linux   - Cross-compile for Linux"
	@echo "  build-darwin  - Cross-compile for macOS"
