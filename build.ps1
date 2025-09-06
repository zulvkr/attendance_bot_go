# Attendance Bot Go - Build Script for Windows
# PowerShell alternative to Makefile/Justfile for Windows users

param(
    [Parameter(Position=0)]
    [string]$Command = "help"
)

$BinaryName = "attendance-bot"
$SetupBinary = "setup-totp"

function Show-Help {
    Write-Host "Attendance Bot Go - Build Script" -ForegroundColor Green
    Write-Host ""
    Write-Host "Usage: .\build.ps1 <command>" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Available commands:" -ForegroundColor Cyan
    Write-Host "  build         - Build the bot and setup utility"
    Write-Host "  clean         - Clean build artifacts"
    Write-Host "  test          - Run tests"
    Write-Host "  test-coverage - Run tests with coverage"
    Write-Host "  deps          - Download and tidy dependencies"
    Write-Host "  setup-totp    - Build and run TOTP setup utility"
    Write-Host "  dev           - Build and run for development"
    Write-Host "  docker-build  - Build Docker image"
    Write-Host "  docker-run    - Run with docker-compose"
    Write-Host "  docker-stop   - Stop docker-compose"
    Write-Host "  docker-logs   - View docker logs"
    Write-Host "  init          - Initialize project (deps + setup)"
    Write-Host "  start         - Quick start (init + run)"
    Write-Host "  release       - Build release binaries for all platforms"
    Write-Host "  help          - Show this help message"
}

function Invoke-Build {
    Write-Host "Building attendance bot and setup utility..." -ForegroundColor Green
    go build -o "$BinaryName.exe" ./cmd/bot
    go build -o "$SetupBinary.exe" ./cmd/setup-totp
    Write-Host "Build completed!" -ForegroundColor Green
}

function Invoke-Clean {
    Write-Host "Cleaning build artifacts..." -ForegroundColor Green
    go clean
    Remove-Item "$BinaryName.exe" -ErrorAction SilentlyContinue
    Remove-Item "$SetupBinary.exe" -ErrorAction SilentlyContinue
    Remove-Item "$BinaryName" -ErrorAction SilentlyContinue
    Remove-Item "$SetupBinary" -ErrorAction SilentlyContinue
    Write-Host "Clean completed!" -ForegroundColor Green
}

function Invoke-Test {
    Write-Host "Running tests..." -ForegroundColor Green
    go test -v ./...
}

function Invoke-TestCoverage {
    Write-Host "Running tests with coverage..." -ForegroundColor Green
    go test -cover ./...
}

function Invoke-Deps {
    Write-Host "Downloading and tidying dependencies..." -ForegroundColor Green
    go mod download
    go mod tidy
    Write-Host "Dependencies updated!" -ForegroundColor Green
}

function Invoke-SetupTotp {
    Write-Host "Building and running TOTP setup utility..." -ForegroundColor Green
    go build -o "$SetupBinary.exe" ./cmd/setup-totp
    & ".\$SetupBinary.exe"
}

function Invoke-Dev {
    Write-Host "Building and running for development..." -ForegroundColor Green
    Invoke-Build
    if (Test-Path "$BinaryName.exe") {
        Write-Host "Starting attendance bot..." -ForegroundColor Yellow
        & ".\$BinaryName.exe"
    } else {
        Write-Host "Build failed!" -ForegroundColor Red
        exit 1
    }
}

function Invoke-DockerBuild {
    Write-Host "Building Docker image..." -ForegroundColor Green
    docker build -t attendance-bot-go .
}

function Invoke-DockerRun {
    Write-Host "Running with docker-compose..." -ForegroundColor Green
    docker-compose up -d
}

function Invoke-DockerStop {
    Write-Host "Stopping docker-compose..." -ForegroundColor Green
    docker-compose down
}

function Invoke-DockerLogs {
    Write-Host "Viewing docker logs..." -ForegroundColor Green
    docker-compose logs -f
}

function Invoke-Init {
    Write-Host "Initializing project..." -ForegroundColor Green
    Invoke-Deps
    
    # Create .env from .env.example if it doesn't exist
    if (!(Test-Path ".env")) {
        if (Test-Path ".env.example") {
            Copy-Item ".env.example" ".env"
            Write-Host "Created .env from .env.example - please edit with your values" -ForegroundColor Yellow
        }
    }
    
    Invoke-SetupTotp
    Write-Host "Project initialized! Edit .env file with your bot token and run '.\build.ps1 dev' to start" -ForegroundColor Green
}

function Invoke-Start {
    Write-Host "Quick start..." -ForegroundColor Green
    Invoke-Init
    Invoke-Dev
}

function Invoke-Release {
    Write-Host "Building release binaries for all platforms..." -ForegroundColor Green
    
    # Create release directory
    New-Item -ItemType Directory -Force -Path "release" | Out-Null
    
    # Build for Windows
    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
    go build -ldflags "-w -s" -o "release/$BinaryName-windows-amd64.exe" ./cmd/bot
    go build -ldflags "-w -s" -o "release/$SetupBinary-windows-amd64.exe" ./cmd/setup-totp
    
    # Build for Linux
    $env:GOOS = "linux"
    $env:GOARCH = "amd64"
    go build -ldflags "-w -s" -o "release/$BinaryName-linux-amd64" ./cmd/bot
    go build -ldflags "-w -s" -o "release/$SetupBinary-linux-amd64" ./cmd/setup-totp
    
    # Build for macOS
    $env:GOOS = "darwin"
    $env:GOARCH = "amd64"
    go build -ldflags "-w -s" -o "release/$BinaryName-darwin-amd64" ./cmd/bot
    go build -ldflags "-w -s" -o "release/$SetupBinary-darwin-amd64" ./cmd/setup-totp
    
    # Reset environment
    Remove-Item Env:\GOOS -ErrorAction SilentlyContinue
    Remove-Item Env:\GOARCH -ErrorAction SilentlyContinue
    
    Write-Host "Release binaries built in release/ directory" -ForegroundColor Green
}

# Main script logic
switch ($Command.ToLower()) {
    "build" { Invoke-Build }
    "clean" { Invoke-Clean }
    "test" { Invoke-Test }
    "test-coverage" { Invoke-TestCoverage }
    "deps" { Invoke-Deps }
    "setup-totp" { Invoke-SetupTotp }
    "dev" { Invoke-Dev }
    "docker-build" { Invoke-DockerBuild }
    "docker-run" { Invoke-DockerRun }
    "docker-stop" { Invoke-DockerStop }
    "docker-logs" { Invoke-DockerLogs }
    "init" { Invoke-Init }
    "start" { Invoke-Start }
    "release" { Invoke-Release }
    "help" { Show-Help }
    default { 
        Write-Host "Unknown command: $Command" -ForegroundColor Red
        Write-Host ""
        Show-Help 
    }
}
