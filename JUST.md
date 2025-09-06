# Using Just Command Runner

This project includes a `justfile` for cross-platform build automation as an alternative to the traditional `Makefile`.

## Installation

### Windows

**Using Scoop (Recommended):**

```powershell
scoop install just
```

**Using Chocolatey:**

```powershell
choco install just
```

**Using Cargo (if you have Rust installed):**

```powershell
cargo install just
```

**Manual Installation:**

1. Download the latest release from [GitHub](https://github.com/casey/just/releases)
2. Extract `just.exe` to a directory in your PATH

### Linux/macOS

**Using Homebrew:**

```bash
brew install just
```

**Using Cargo:**

```bash
cargo install just
```

**Using package managers:**

```bash
# Arch Linux
sudo pacman -S just

# Ubuntu/Debian (via snap)
sudo snap install --edge just

# Fedora
sudo dnf install just
```

## Quick Start

Once Just is installed:

```bash
# View all available commands
just

# Quick start for Windows users
just start

# Build and run for development
just dev-windows     # Windows
just dev            # Unix/Linux/macOS

# Run tests
just test

# Build for production
just build-prod

# Docker operations
just docker-build
just docker-run
```

## Why Just?

- **Cross-platform**: Works consistently on Windows, Linux, and macOS
- **Simple syntax**: Easier to read and write than Makefiles
- **Built-in features**: Command listing, parameter handling, and more
- **No dependencies**: Single binary with no runtime dependencies
- **Better Windows support**: Handles Windows paths and commands properly

## Fallback

If you prefer not to install Just, you can always use the traditional `Makefile` (requires `make` command) or run Go commands directly:

```bash
# Using Make
make build

# Direct Go commands
go build -o attendance-bot cmd/bot/main.go
go run cmd/bot/main.go
```
