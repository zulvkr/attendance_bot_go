# Attendance Bot Go - GitHub Copilot Instructions

This is a Go-based Telegram attendance bot with TOTP authentication. Always follow these instructions first and only search for additional information if these instructions are incomplete or incorrect.

## Working Effectively

### Bootstrap and Build the Repository
- Check Go installation: `go version` (requires Go 1.22+)
- Download dependencies: `go mod download` -- takes ~6 seconds. NEVER CANCEL.
- Tidy dependencies: `go mod tidy` -- takes ~13 seconds. NEVER CANCEL.
- Build main bot: `go build -o attendance-bot ./cmd/bot` -- takes ~18 seconds first time, ~1 second subsequent builds. NEVER CANCEL. Set timeout to 60+ seconds.
- Build setup utility: `go build -o setup-totp ./cmd/setup-totp` -- takes ~1 second. Set timeout to 30+ seconds.

### Alternative Build Systems
- **Make (RECOMMENDED)**: `make all` -- builds both binaries in ~1 second after initial setup
- **PowerShell (Windows only)**: `.\build.ps1 build`
- **Batch script (Windows only)**: `.\build.bat build`
- **Just**: NOT AVAILABLE in environment - do not use

### Testing
- Run all tests: `go test -v ./...` -- takes ~4 seconds. NEVER CANCEL. Set timeout to 30+ minutes.
- **IMPORTANT**: No test files exist in this codebase - go test will show "no test files" for all packages
- Run with coverage: `go test -cover ./...`
- Run vet: `go vet ./...` -- takes ~1 second
- Format code: `go fmt ./...` -- takes ~1 second

### Environment Setup and Configuration
- **CRITICAL**: The bot requires environment variables, NOT .env files by default
- Generate TOTP secret: `./setup-totp` -- creates new TOTP secret and shows setup instructions
- **Required environment variables**:
  ```bash
  export BOT_TOKEN="your_telegram_bot_token_here"  # Must be 10+ characters
  export TOTP_SECRET="generated_32_char_secret"     # Must be 16+ characters  
  export ADMIN_PASSWORD="your_admin_password"      # Must be 8+ characters
  export NODE_ENV="development"                     # Optional, defaults to development
  export DATABASE_PATH="data/attendance.db"        # Optional, defaults to data/attendance.db
  ```

### Running the Application
- **ALWAYS set environment variables first**
- **ALWAYS run TOTP setup before first use**: `./setup-totp`
- Start the bot: `./attendance-bot`
- Expected startup sequence:
  1. "Configuration loaded" 
  2. "Database initialized"
  3. "Starting bot..."
  4. Will fail with network error if BOT_TOKEN is invalid (this is expected in testing)

### Docker Operations  
- **WARNING**: Docker builds fail due to network restrictions in this environment
- Build command: `docker build -t attendance-bot-go .` -- normally takes ~2 minutes
- Run with docker-compose: `docker-compose up -d`
- View logs: `docker-compose logs -f`
- Stop: `docker-compose down`

## Validation and Testing Scenarios

### Manual Validation Requirements
After making changes, ALWAYS run through these validation scenarios:

1. **Build Validation**:
   ```bash
   make clean  # or go clean
   make all    # or go build -o attendance-bot ./cmd/bot && go build -o setup-totp ./cmd/setup-totp
   ```

2. **TOTP Setup Validation**:
   ```bash
   ./setup-totp
   # Verify it outputs: Generated TOTP Secret, OTP Auth URL, Setup Instructions, Current TOTP token
   ```

3. **Configuration Validation**:
   ```bash
   export BOT_TOKEN="1234567890:ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890123456789"
   export TOTP_SECRET="MRDSVPGNARDWTZZYHOZG3SM7KEKP2FGX"  
   export ADMIN_PASSWORD="testpassword123"
   ./attendance-bot
   # Should show "Configuration loaded" and "Database initialized" before network error
   ```

4. **Cross-Platform Build Validation**:
   ```bash
   make build-linux   # Takes ~1 second
   make build-windows # Takes ~1 second  
   make build-darwin  # Takes ~1 second
   ```

## Repository Structure and Navigation

### Key Directories and Files
- `cmd/bot/main.go` - Main bot application entry point
- `cmd/setup-totp/main.go` - TOTP setup utility
- `internal/config/config.go` - Configuration management and validation
- `internal/bot/` - Telegram bot implementation (handlers.go, telegram.go)
- `internal/attendance/` - Core attendance logic and TOTP service
- `internal/database/` - SQLite database layer and repository pattern
- `internal/reports/` - CSV report generation
- `internal/utils/` - Utility functions for date/validation
- `pkg/models/` - Data models
- `data/` - SQLite database storage directory (created at runtime)

### Build Configuration Files
- `go.mod` - Go module definition with modernc.org/sqlite dependency
- `Makefile` - Unix-style build automation (WORKS)
- `justfile` - Just command runner configuration (NOT AVAILABLE)
- `build.ps1` - PowerShell build script (Windows only)
- `build.bat` - Batch build script (Windows only)
- `Dockerfile` - Docker configuration (builds fail due to network restrictions)

### Important Configuration Files
- `.env.example` - Template for environment variables
- `LoadEnv.ps1` - PowerShell script to load .env files (Windows only)
- `BUILD.md` - Comprehensive build documentation
- `specification.md` - Project architecture and implementation details

## Common Tasks and Commands

### Development Workflow
1. Make code changes
2. Run `go fmt ./...` to format code
3. Run `go vet ./...` to check for issues  
4. Run `make all` to build both binaries
5. Test configuration with environment variables and `./attendance-bot`
6. Validate TOTP setup with `./setup-totp`

### Frequent File Locations
- Configuration validation: `internal/config/config.go`
- Bot command handlers: `internal/bot/handlers.go`
- TOTP implementation: `internal/attendance/` directory
- Database operations: `internal/database/` directory
- Main bot logic: `cmd/bot/main.go`

### Debugging Common Issues
- "missing environment variables" error: Export required environment variables
- Build failures: Ensure Go 1.22+ is installed, run `go mod download`
- TOTP setup issues: Run `./setup-totp` to generate new secret
- Database errors: Check `data/` directory permissions
- Network errors on bot start: Expected with dummy BOT_TOKEN (actual Telegram token needed for real operation)

## Timing Expectations and Timeouts

**CRITICAL**: NEVER CANCEL long-running commands. Always use appropriate timeouts:

- `go mod download`: ~6 seconds - Set timeout to 300+ seconds
- `go mod tidy`: ~13 seconds - Set timeout to 300+ seconds  
- `go build ./cmd/bot`: ~18 seconds first time, ~1 second subsequent - Set timeout to 60+ seconds
- `go test ./...`: ~4 seconds - Set timeout to 1800+ seconds (30 minutes)
- `make all`: ~1 second after initial setup - Set timeout to 120+ seconds
- `docker build`: ~2 minutes when working - Set timeout to 3600+ seconds (60 minutes)

## Dependencies and External Tools

### Required
- Go 1.22+ (installed)
- make (installed and working)

### Optional/Not Available
- Just command runner (NOT AVAILABLE - do not use `just` commands)
- golangci-lint (NOT AVAILABLE - use `go vet` instead)
- Docker (available but builds fail due to network restrictions)

### Go Dependencies
- `modernc.org/sqlite` - Pure Go SQLite driver (no CGO required)
- Standard library only for HTTP client, JSON, crypto, logging

## Environment Considerations

This environment has these limitations:
- Network restrictions prevent Docker Alpine package installation
- Just command runner is not installed
- golangci-lint is not available
- .env files are not automatically loaded (use environment variables)

Always use direct Go commands or Make for reliable builds in this environment.