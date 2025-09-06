# Build Options Summary

This project provides multiple ways to build and run the attendance bot to accommodate different development environments and preferences.

## Quick Start Options

### Option 1: Windows Batch Script (Windows - Universal)

**Works on any Windows system - no PowerShell execution policy issues**

```batch
.\build.bat run
```

### Option 2: PowerShell Script (Windows - Advanced)

**Full-featured script with advanced options**

```powershell
.\build.ps1 start
```

### Option 3: Just Command Runner (Cross-platform)

**Requires Just installation - most flexible**

```bash
# Install Just first (see JUST.md)
just start
```

### Option 4: Traditional Make (Unix-style)

**Requires make command - familiar for Unix developers**

```bash
make all
```

### Option 5: Direct Go Commands (Universal)

**Always works if Go is installed**

```bash
go run cmd/setup-totp/main.go  # Setup TOTP first
go run cmd/bot/main.go         # Run the bot
```

## Comparison

| Method      | Windows      | Linux/macOS  | Dependencies | Features                        |
| ----------- | ------------ | ------------ | ------------ | ------------------------------- |
| `build.bat` | ✅ Excellent | ❌ No        | None         | Basic automation                |
| `build.ps1` | ✅ Excellent | ❌ No        | PowerShell   | Full automation                 |
| `justfile`  | ✅ Good      | ✅ Excellent | Just command | Full automation, cross-platform |
| `Makefile`  | ⚠️ Limited   | ✅ Excellent | make command | Traditional, familiar           |
| Direct Go   | ✅ Basic     | ✅ Basic     | Go only      | Manual steps                    |

## Feature Matrix

| Feature             | Batch     | PowerShell  | Just      | Make      | Direct Go |
| ------------------- | --------- | ----------- | --------- | --------- | --------- |
| Build automation    | ✅        | ✅          | ✅        | ✅        | ❌        |
| Cross-compilation   | ❌        | ✅          | ✅        | ✅        | Manual    |
| Testing integration | ✅        | ✅          | ✅        | ✅        | Manual    |
| Docker integration  | ❌        | ✅          | ✅        | ✅        | Manual    |
| Environment setup   | ❌        | ✅          | ✅        | ❌        | Manual    |
| Release builds      | ❌        | ✅          | ✅        | ✅        | Manual    |
| Help/Documentation  | ✅        | ✅          | ✅        | ✅        | N/A       |
| PowerShell issues   | ✅ Immune | ⚠️ May fail | ✅ Immune | ✅ Immune | ✅ Immune |

## Recommendations

- **Windows corporate environments**: Use `.\build.bat` (no admin rights needed)
- **Windows developers**: Use `.\build.ps1` for full features
- **Cross-platform teams**: Use `just` for consistency
- **Unix traditionalists**: Use `make`
- **Minimalists**: Use direct Go commands
- **CI/CD pipelines**: Use direct Go commands or Docker

## All Available Commands

### PowerShell Script (`.\build.ps1`)

- `help` - Show available commands
- `build` - Build bot and setup utility
- `clean` - Clean build artifacts
- `test` - Run tests
- `deps` - Download dependencies
- `setup-totp` - Run TOTP setup
- `dev` - Build and run for development
- `init` - Initialize project
- `start` - Quick start (init + run)
- `release` - Build for all platforms
- `docker-build/run/stop/logs` - Docker operations

### Just Commands (`just <command>`)

All PowerShell commands plus:

- `build-windows/linux/darwin` - Platform-specific builds
- `test-coverage/race` - Advanced testing
- `check` - Format + vet + test
- `fmt/vet/lint` - Code quality tools

### Make Commands (`make <target>`)

Traditional Unix-style targets:

- `all` - Build everything
- `build/clean/test` - Basic operations
- `build-prod` - Production build
- `docker-*` - Docker operations

Choose the method that best fits your workflow and environment!
