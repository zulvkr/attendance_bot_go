# Attendance Bot Go Port - Specification ✅ COMPLETED

## Overview

This document outlined the plan to port the TypeScript Attendance Bot to Go, maintaining all functionality while following Go idioms and best practices using only the standard library where possible.

**Status: ✅ IMPLEMENTATION COMPLETED**

The Go port has been successfully implemented with full feature parity to the TypeScript version.

## Current TypeScript Bot Analysis

### Key Features

1. **TOTP Authentication**: Uses time-based one-time passwords for secure attendance marking
2. **Dual Attendance System**: Supports check-in and check-out with automatic detection
3. **Telegram Bot Integration**: Full Telegram bot with commands and message handling
4. **SQLite Database**: Persistent storage with proper schema and indexing
5. **Reporting System**: Daily reports, user history, and CSV export
6. **User Aliases**: Support for custom display names
7. **Time Zone Support**: Indonesian locale formatting
8. **Docker Support**: Containerized deployment

### Core Components

- **Bot Handler**: Telegram bot commands and message processing
- **Attendance Service**: Core business logic for attendance management
- **Database Layer**: SQLite with repository pattern
- **TOTP Service**: Time-based OTP verification
- **Utilities**: Date/time formatting, CSV generation
- **Configuration**: Environment variable management with validation

## Go Implementation Plan

### 1. Project Structure

```
attendance-bot-go/
├── cmd/
│   ├── bot/
│   │   └── main.go           # Bot entry point
│   └── setup-totp/
│       └── main.go           # TOTP setup utility
├── internal/
│   ├── config/
│   │   └── config.go         # Environment configuration
│   ├── database/
│   │   ├── sqlite.go         # SQLite connection and schema
│   │   └── repository.go     # Data access layer
│   ├── attendance/
│   │   ├── service.go        # Core attendance business logic
│   │   └── totp.go           # TOTP verification
│   ├── bot/
│   │   ├── handlers.go       # Command handlers
│   │   └── bot.go            # Bot setup and middleware
│   ├── reports/
│   │   └── csv.go            # CSV report generation
│   └── utils/
│       ├── date.go           # Date/time utilities
│       └── validation.go     # Input validation
├── pkg/
│   └── models/
│       └── attendance.go     # Data models
├── data/                     # SQLite database directory
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
└── README.md
```

### 2. Technology Stack & Dependencies

#### Standard Library Usage

- `net/http` - HTTP client for Telegram API
- `database/sql` - Database operations with SQLite driver
- `encoding/json` - JSON parsing/encoding
- `crypto/hmac`, `crypto/sha1` - TOTP implementation
- `time` - Time handling and formatting
- `log/slog` - Structured logging
- `os` - Environment variables and file operations
- `encoding/csv` - CSV report generation

#### External Dependencies (Minimal)

- `modernc.org/sqlite` - Pure Go SQLite driver (CGo-free)
- No Telegram bot framework - implement using standard HTTP client

### 3. Core Components Implementation

#### 3.1 Configuration Management

```go
// internal/config/config.go
type Config struct {
    BotToken     string
    TOTPSecret   string
    AdminPassword string
    Environment  string
    DatabasePath string
}

func Load() (*Config, error) {
    // Load from environment variables with validation
    // Default values for development
}
```

#### 3.2 Database Layer

```go
// pkg/models/attendance.go
type AttendanceRecord struct {
    ID        int64     `json:"id"`
    UserID    int64     `json:"user_id"`
    Username  string    `json:"username"`
    FirstName string    `json:"first_name"`
    LastName  *string   `json:"last_name,omitempty"`
    Timestamp time.Time `json:"timestamp"`
    Type      string    `json:"type"` // "check_in" or "check_out"
    Date      string    `json:"date"` // YYYY-MM-DD format
}

type UserAlias struct {
    UserID    int64   `json:"user_id"`
    FirstName string  `json:"first_name"`
    LastName  *string `json:"last_name,omitempty"`
}
```

```go
// internal/database/repository.go
type Repository struct {
    db *sql.DB
}

func (r *Repository) InsertAttendance(record *models.AttendanceRecord) error
func (r *Repository) GetUserAttendanceToday(userID int64, date string) ([]models.AttendanceRecord, error)
func (r *Repository) GetUserAttendanceHistory(userID int64, days int) ([]models.AttendanceRecord, error)
func (r *Repository) GetDailyReport(date string) ([]models.AttendanceRecord, error)
```

#### 3.3 TOTP Implementation

```go
// internal/attendance/totp.go
type TOTPService struct {
    secret string
}

func NewTOTPService(secret string) *TOTPService
func (t *TOTPService) Verify(token string) bool
func (t *TOTPService) Generate() string
func generateSecret() string
```

#### 3.4 Attendance Service

```go
// internal/attendance/service.go
type Service struct {
    repo *database.Repository
    totp *TOTPService
}

type AttendanceResult struct {
    Success bool
    Message string
    Record  *models.AttendanceRecord
}

func (s *Service) MarkAttendance(userID int64, username, firstName string, lastName *string, otp string) (*AttendanceResult, error)
func (s *Service) GetUserStatus(userID int64, date string) (*UserStatus, error)
func (s *Service) GenerateReport(date string) (string, error)
```

#### 3.5 Telegram Bot Implementation

```go
// internal/bot/bot.go
type Bot struct {
    token      string
    attendance *attendance.Service
    client     *http.Client
}

type Update struct {
    UpdateID int64   `json:"update_id"`
    Message  *Message `json:"message,omitempty"`
}

type Message struct {
    MessageID int64  `json:"message_id"`
    From      *User  `json:"from,omitempty"`
    Chat      *Chat  `json:"chat"`
    Text      string `json:"text,omitempty"`
}

func (b *Bot) Start() error
func (b *Bot) handleUpdate(update *Update) error
func (b *Bot) sendMessage(chatID int64, text string) error
```

#### 3.6 Command Handlers

```go
// internal/bot/handlers.go
func (b *Bot) handleStart(msg *Message)
func (b *Bot) handleHelp(msg *Message)
func (b *Bot) handleReport(msg *Message)
func (b *Bot) handleHistory(msg *Message)
func (b *Bot) handleStatus(msg *Message)
func (b *Bot) handleAlias(msg *Message)
func (b *Bot) handleOTP(msg *Message, otp string)
```

### 4. Key Implementation Details

#### 4.1 TOTP Implementation (RFC 6238)

- Implement TOTP using standard library crypto functions
- 30-second time step
- SHA-1 HMAC (standard for most authenticator apps)
- 6-digit codes

#### 4.2 Database Schema

```sql
CREATE TABLE IF NOT EXISTS attendance (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    username TEXT NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT,
    timestamp TEXT NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('check_in', 'check_out')),
    date TEXT NOT NULL,
    UNIQUE(user_id, date, type)
);

CREATE INDEX IF NOT EXISTS idx_user_date ON attendance(user_id, date);
CREATE INDEX IF NOT EXISTS idx_date ON attendance(date);

CREATE TABLE IF NOT EXISTS alias (
    user_id INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT
);
```

#### 4.3 Telegram Bot API Integration

- Use `net/http` to implement Telegram Bot API calls
- Long polling for receiving updates
- Proper error handling and retry logic
- Rate limiting awareness

#### 4.4 Error Handling

- Comprehensive error wrapping with context
- Structured logging with `log/slog`
- Graceful error recovery for non-critical operations
- User-friendly error messages in Indonesian

#### 4.5 Time Zone Handling

- Default to Asia/Jakarta timezone
- Proper date parsing and formatting
- Consistent date string format (YYYY-MM-DD)

### 5. Migration Considerations

#### 5.1 Database Compatibility

- Maintain exact same SQLite schema
- Support migration from existing TypeScript bot database
- Preserve all existing data integrity

#### 5.2 Feature Parity

- All Telegram commands maintained
- Same business logic for attendance rules
- Identical TOTP compatibility
- Same CSV export format

#### 5.3 Configuration

- Same environment variable names
- Compatible Docker setup
- Same data volume structure

### 6. Development Phases

#### Phase 1: Core Infrastructure

1. Project setup with Go modules
2. Configuration management
3. Database layer and models
4. TOTP implementation
5. Basic tests

#### Phase 2: Attendance Logic

1. Attendance service implementation
2. Business rules for check-in/check-out
3. User status tracking
4. Data validation

#### Phase 3: Telegram Bot

1. Telegram API client implementation
2. Command handlers
3. Message processing
4. Session management

#### Phase 4: Reporting & Utilities

1. Daily reports generation
2. CSV export functionality
3. User history tracking
4. Date/time utilities

#### Phase 5: Deployment & Testing

1. Docker configuration
2. Integration testing
3. Performance testing
4. Documentation

### 7. Go Idioms and Best Practices

#### 7.1 Code Organization

- Clear separation of concerns
- Internal packages for implementation details
- Public API in pkg directory
- Command-line tools in cmd directory

#### 7.2 Error Handling

- Explicit error returns
- Error wrapping with context
- Sentinel errors for common cases
- No panics in normal operation

#### 7.3 Concurrency

- Goroutines for handling multiple updates
- Proper channel usage for coordination
- Context for cancellation and timeouts
- Worker pool pattern for request handling

#### 7.4 Testing

- Table-driven tests
- Testable interfaces
- Mock implementations for testing
- Integration tests with test database

#### 7.5 Performance

- Efficient database queries
- Connection pooling
- Minimal allocations in hot paths
- Proper resource cleanup

### 8. Success Criteria

1. **Functional Compatibility**: All features from TypeScript version work identically
2. **Performance**: Comparable or better performance with lower resource usage
3. **Maintainability**: Clean, idiomatic Go code following standard conventions
4. **Deployability**: Easy Docker deployment with same interface
5. **Reliability**: Robust error handling and graceful degradation
6. **Documentation**: Clear documentation and examples

This specification ensures a complete, idiomatic Go implementation that maintains full compatibility with the existing TypeScript bot while leveraging Go's strengths in performance, concurrency, and deployment.
