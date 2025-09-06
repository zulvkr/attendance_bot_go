# Attendance Bot Go - Justfile
# Cross-platform build automation using Just (https://github.com/casey/just)
# Install with: cargo install just

# Variables
binary_name := "attendance-bot"
setup_binary := "setup-totp"

# Default recipe to display available commands
default:
    @just --list

# Install dependencies and build everything
all: deps build

# Download and tidy Go dependencies
deps:
    go mod download
    go mod tidy

# Build both the bot and setup utility
build:
    go build -o {{binary_name}} ./cmd/bot
    go build -o {{setup_binary}} ./cmd/setup-totp

# Build for Windows (with .exe extension)
build-windows:
    $env:GOOS="windows"; $env:GOARCH="amd64"; go build -o {{binary_name}}.exe ./cmd/bot
    $env:GOOS="windows"; $env:GOARCH="amd64"; go build -o {{setup_binary}}.exe ./cmd/setup-totp

# Build for Linux
build-linux:
    $env:GOOS="linux"; $env:GOARCH="amd64"; go build -o {{binary_name}} ./cmd/bot
    $env:GOOS="linux"; $env:GOARCH="amd64"; go build -o {{setup_binary}} ./cmd/setup-totp

# Build for macOS
build-darwin:
    $env:GOOS="darwin"; $env:GOARCH="amd64"; go build -o {{binary_name}} ./cmd/bot
    $env:GOOS="darwin"; $env:GOARCH="amd64"; go build -o {{setup_binary}} ./cmd/setup-totp

# Production build with optimizations (Linux)
build-prod:
    $env:CGO_ENABLED="0"; $env:GOOS="linux"; go build -a -installsuffix cgo -ldflags "-w -s" -o {{binary_name}} ./cmd/bot

# Clean build artifacts (Windows compatible)
clean:
    go clean
    @echo "Cleaning build artifacts..."
    -Remove-Item {{binary_name}} -ErrorAction SilentlyContinue
    -Remove-Item {{binary_name}}.exe -ErrorAction SilentlyContinue
    -Remove-Item {{setup_binary}} -ErrorAction SilentlyContinue
    -Remove-Item {{setup_binary}}.exe -ErrorAction SilentlyContinue

# Clean build artifacts (Unix/Linux/macOS)
clean-unix:
    go clean
    rm -f {{binary_name}}
    rm -f {{setup_binary}}

# Run tests
test:
    go test -v ./...

# Run tests with coverage
test-coverage:
    go test -cover ./...

# Run tests with race detection
test-race:
    go test -race ./...

# Build and run TOTP setup utility
setup-totp:
    go build -o {{setup_binary}} ./cmd/setup-totp
    ./{{setup_binary}}

# Build and run TOTP setup utility (Windows)
setup-totp-windows:
    go build -o {{setup_binary}}.exe ./cmd/setup-totp
    ./{{setup_binary}}.exe

# Build and run the bot for development
dev: build
    ./{{binary_name}}

# Build and run the bot for development (Windows)
dev-windows: build-windows
    ./{{binary_name}}.exe

# Docker targets
docker-build:
    docker build -t attendance-bot-go .

# Run with docker-compose
docker-run:
    docker-compose up -d

# Stop docker-compose
docker-stop:
    docker-compose down

# View docker logs
docker-logs:
    docker-compose logs -f

# Format Go code
fmt:
    go fmt ./...

# Run Go linter (requires golangci-lint)
lint:
    golangci-lint run

# Vet Go code for potential issues
vet:
    go vet ./...

# Run security scanner (requires gosec)
security:
    gosec ./...

# Download Go modules to local cache
mod-download:
    go mod download

# Verify dependencies have expected content
mod-verify:
    go mod verify

# Generate Go documentation
docs:
    godoc -http=:6060

# Check for outdated dependencies (requires go-mod-outdated)
mod-outdated:
    go list -u -m all | grep "\["

# Run all checks (format, vet, test)
check: fmt vet test

# Release build for all platforms
release: clean
    @echo "Building release binaries..."
    $env:GOOS="windows"; $env:GOARCH="amd64"; go build -ldflags "-w -s" -o release/{{binary_name}}-windows-amd64.exe ./cmd/bot
    $env:GOOS="linux"; $env:GOARCH="amd64"; go build -ldflags "-w -s" -o release/{{binary_name}}-linux-amd64 ./cmd/bot
    $env:GOOS="darwin"; $env:GOARCH="amd64"; go build -ldflags "-w -s" -o release/{{binary_name}}-darwin-amd64 ./cmd/bot
    $env:GOOS="windows"; $env:GOARCH="amd64"; go build -ldflags "-w -s" -o release/{{setup_binary}}-windows-amd64.exe ./cmd/setup-totp
    $env:GOOS="linux"; $env:GOARCH="amd64"; go build -ldflags "-w -s" -o release/{{setup_binary}}-linux-amd64 ./cmd/setup-totp
    $env:GOOS="darwin"; $env:GOARCH="amd64"; go build -ldflags "-w -s" -o release/{{setup_binary}}-darwin-amd64 ./cmd/setup-totp
    @echo "Release binaries built in release/ directory"

# Install development dependencies
install-deps:
    @echo "Installing development dependencies..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
    @echo "Development dependencies installed"

# Show project information
info:
    @echo "Project: Attendance Bot Go"
    @echo "Go version: $(go version)"
    @echo "Module: $(go list -m)"
    @echo "Dependencies:"
    @go list -m all

# Initialize a new .env file from .env.example
init-env:
    @if (Test-Path .env) { echo ".env file already exists" } else { Copy-Item .env.example .env; echo "Created .env from .env.example - please edit with your values" }

# Run setup and initialization
init: deps init-env setup-totp-windows
    @echo "Project initialized! Edit .env file with your bot token and run 'just dev-windows' to start"

# Quick start for Windows users
start: init
    @echo "Starting attendance bot..."
    just dev-windows
