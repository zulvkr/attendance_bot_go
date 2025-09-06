# Attendance Bot Go

A Telegram bot for employee attendance tracking using TOTP (Time-based One-Time Password) authentication with SQLite database storage, written in idiomatic Go.

## Features

- 🔐 TOTP-based attendance marking
- ⏰ Automatic late detection (after 9:00 AM)
- 📊 Daily attendance reports
- 📈 Personal attendance history
- 🚫 Prevents duplicate attendance marking per day
- 💾 SQLite database for persistent data storage
- 🔍 Efficient database operations with proper indexing

## Tech Stack

- **Go 1.22+** - Modern Go with generics support
- **Standard Library HTTP Client** - For Telegram Bot API (no external bot framework)
- **SQLite3** - Lightweight embedded database (using pure Go driver)
- **Pure Go TOTP** - Time-based OTP implementation using standard crypto
- **Structured Logging** - Using log/slog for better observability

## Setup

### 1. Generate TOTP Secret

Run the setup utility to generate a TOTP secret and QR code information:

```bash
go run cmd/setup-totp/main.go
```

This will:

- Generate a new TOTP secret
- Create setup instructions
- Generate a sample .env file
- Show the current TOTP token for testing

### 2. Create Telegram Bot

1. Message [@BotFather](https://t.me/botfather) on Telegram
2. Create a new bot with `/newbot`
3. Copy the bot token

### 3. Configure Environment

Create a `.env` file based on the generated `.env.example`:

```env
BOT_TOKEN=your_telegram_bot_token_here
TOTP_SECRET=your_generated_totp_secret_here
ADMIN_PASSWORD=your_admin_password_here
NODE_ENV=development
DATABASE_PATH=data/attendance.db
```

### 4. Setup Authenticator App

1. Install an authenticator app (Google Authenticator, Authy, etc.)
2. Use the OTP Auth URL from the setup script to add the account
3. Or manually enter the TOTP secret

### 5. Run the Bot

**Using Windows Batch Script (Windows - No dependencies, works everywhere):**

```batch
# Simple commands for any Windows environment
.\build.bat help      # Show available commands
.\build.bat build     # Build the bot
.\build.bat run       # Build and run
.\build.bat setup     # Run TOTP setup
```

**Using PowerShell Script (Windows - Advanced features):**

```powershell
# Full-featured script for Windows
.\build.ps1 help      # Show all commands
.\build.ps1 start     # Quick start (builds, sets up TOTP, and runs)
.\build.ps1 dev       # Build and run for development
.\build.ps1 init      # Initialize project
```

**Using Just (Cross-platform - Most flexible):**

Install Just command runner:

```bash
# Windows (using Scoop)
scoop install just

# Windows (using Chocolatey)
choco install just

# Or download from: https://github.com/casey/just/releases
```

Then use Just commands:

```bash
# Quick start (builds everything and runs setup)
just start

# Development mode
just dev-windows    # Windows
just dev           # Unix/Linux/macOS

# View all available commands
just --list
```

**Manual Development mode:**

```bash
go run cmd/bot/main.go
```

**Manual Production build:**

```bash
go build -o attendance-bot cmd/bot/main.go
./attendance-bot
```

**Using Docker:**

```bash
docker-compose up -d
```

## Database Structure

The bot uses SQLite with the following tables:

### `attendance` table

| Column     | Type    | Description                  |
| ---------- | ------- | ---------------------------- |
| id         | INTEGER | Primary key (auto-increment) |
| user_id    | INTEGER | Telegram user ID             |
| username   | TEXT    | Telegram username            |
| first_name | TEXT    | User's first name            |
| last_name  | TEXT    | User's last name (nullable)  |
| timestamp  | TEXT    | ISO timestamp of attendance  |
| type       | TEXT    | 'check_in' or 'check_out'    |
| date       | TEXT    | Date in YYYY-MM-DD format    |

### `alias` table

| Column     | Type    | Description                   |
| ---------- | ------- | ----------------------------- |
| user_id    | INTEGER | Primary key, Telegram user ID |
| first_name | TEXT    | Custom first name             |
| last_name  | TEXT    | Custom last name (nullable)   |

**Indexes:**

- `idx_user_date` on (user_id, date) for fast user attendance lookups
- `idx_date` on date for daily reports
- `idx_user_id` on user_id for user-specific queries
- Unique constraint on (user_id, date, type) to prevent duplicate attendance

## Usage

### For Employees

1. Open the bot in Telegram
2. Send `/start` to see instructions
3. Get your 6-digit code from your authenticator app
4. Send the code to mark attendance

### Available Commands

- 📝 **Send OTP** - Mark attendance with 6-digit code
- 📊 `/report` - View today's attendance report
- 📈 `/history` - View your attendance history (30 days)
- 🔄 `/status` - Check if you've marked attendance today
- 🏷️ `/alias` - Set custom display name
- ❓ `/help` - Show help message

### Attendance Rules

- ✅ **On Time**: Attendance marked before 9:00 AM
- ⚠️ **Late**: Attendance marked after 9:00 AM
- 🚫 **Once per day**: Each employee can only mark each type (check-in/check-out) once per day

## Architecture

### Project Structure

```
attendance-bot-go/
├── cmd/
│   ├── bot/main.go           # Main bot application
│   └── setup-totp/main.go    # TOTP setup utility
├── internal/
│   ├── config/config.go      # Configuration management
│   ├── database/             # Database layer
│   │   ├── sqlite.go         # SQLite connection and schema
│   │   └── repository.go     # Data access layer
│   ├── attendance/           # Business logic
│   │   ├── service.go        # Core attendance logic
│   │   └── totp.go           # TOTP implementation
│   ├── bot/                  # Telegram bot
│   │   ├── telegram.go       # Telegram API client
│   │   └── handlers.go       # Command handlers
│   ├── reports/csv.go        # CSV report generation
│   └── utils/                # Utilities
│       ├── date.go           # Date/time functions
│       └── validation.go     # Input validation
├── pkg/models/attendance.go  # Data models
└── data/                     # SQLite database directory
```

### Key Design Principles

1. **Standard Library First**: Uses Go's standard library wherever possible
2. **Clear Separation**: Business logic separated from transport layer
3. **Dependency Injection**: Services are injected for better testing
4. **Error Handling**: Comprehensive error handling with context
5. **Logging**: Structured logging for better observability
6. **Idiomatic Go**: Follows Go conventions and best practices

## Security Features

- TOTP authentication prevents unauthorized attendance
- Time-based codes expire every 30 seconds
- Input validation and sanitization
- No storage of sensitive authentication data
- User identification through Telegram IDs

## Performance Features

- Efficient SQLite queries with proper indexing
- Connection pooling handled by Go's database/sql
- Minimal memory allocations in hot paths
- Long polling with configurable timeouts
- Graceful shutdown handling

## Development

### Build Automation

This project includes multiple build automation options:

**PowerShell Script (Windows - No dependencies required):**

```powershell
# All-in-one script for Windows users
.\build.ps1 help          # Show available commands
.\build.ps1 start         # Quick start (init + run)
.\build.ps1 build         # Build bot and setup utility
.\build.ps1 dev          # Build and run for development
.\build.ps1 test         # Run tests
.\build.ps1 release      # Build for all platforms
```

**Using Just (Recommended for cross-platform):**

```bash
# Install Just: https://github.com/casey/just
# Windows: scoop install just  OR  choco install just

# View all available commands
just

# Quick start (Windows)
just start

# Build and develop
just build
just dev-windows    # Windows
just dev           # Unix/Linux/macOS

# Cross-platform builds
just build-windows
just build-linux
just build-darwin

# Testing and quality
just test
just test-coverage
just check         # fmt + vet + test

# Release builds
just release       # All platforms

# Docker operations
just docker-build
just docker-run
just docker-logs
```

**Using Make (Traditional):**

```bash
# All commands from Makefile
make help          # Show available targets
make all           # Build everything
make build         # Build bot and setup utility
make test          # Run tests
make docker-build  # Build Docker image
```

### Manual Building

```bash
# Build for current platform
go build -o attendance-bot cmd/bot/main.go

# Build for Linux (for Docker)
GOOS=linux GOARCH=amd64 go build -o attendance-bot cmd/bot/main.go

# Build setup utility
go build -o setup-totp cmd/setup-totp/main.go
```

### Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...
```

### Docker

```bash
# Build image
docker build -t attendance-bot-go .

# Run with docker-compose
docker-compose up -d

# View logs
docker-compose logs -f
```

## Migration from TypeScript Version

This Go implementation maintains full compatibility with the TypeScript version:

- Same database schema
- Same Telegram commands and responses
- Same TOTP implementation (RFC 6238)
- Same business logic and rules
- Compatible Docker setup

You can migrate by:

1. Copying the existing SQLite database to the `data/` directory
2. Using the same `.env` configuration (BOT_TOKEN and TOTP_SECRET)
3. Replacing the TypeScript container with the Go container

## Future Enhancements

- Webhook support for better performance
- Admin web interface
- Bulk user management
- Advanced reporting and analytics
- Integration with HR systems
- Multi-language support

## License

MIT License - see LICENSE file for details.
