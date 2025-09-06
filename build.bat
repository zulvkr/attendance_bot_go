@echo off
REM Attendance Bot Go - Simple Windows Batch Script
REM Alternative for environments where PowerShell execution is restricted

if "%1"=="" goto help
if "%1"=="help" goto help
if "%1"=="build" goto build
if "%1"=="clean" goto clean
if "%1"=="test" goto test
if "%1"=="setup" goto setup
if "%1"=="run" goto run
goto unknown

:help
echo Attendance Bot Go - Windows Batch Script
echo.
echo Usage: build.bat ^<command^>
echo.
echo Available commands:
echo   build  - Build the bot and setup utility
echo   clean  - Clean build artifacts
echo   test   - Run tests
echo   setup  - Run TOTP setup utility
echo   run    - Build and run the bot
echo   help   - Show this help message
echo.
echo For more advanced features, use build.ps1 instead
goto end

:build
echo Building attendance bot...
go build -o attendance-bot.exe ./cmd/bot
go build -o setup-totp.exe ./cmd/setup-totp
if errorlevel 1 (
    echo Build failed!
    goto end
)
echo Build completed successfully!
goto end

:clean
echo Cleaning build artifacts...
del /Q attendance-bot.exe 2>nul
del /Q setup-totp.exe 2>nul
go clean
echo Clean completed!
goto end

:test
echo Running tests...
go test -v ./...
goto end

:setup
echo Running TOTP setup...
go build -o setup-totp.exe ./cmd/setup-totp
if errorlevel 1 (
    echo Build failed!
    goto end
)
setup-totp.exe
goto end

:run
echo Building and running bot...
call :build
if exist attendance-bot.exe (
    echo Starting attendance bot...
    attendance-bot.exe
) else (
    echo Build failed - cannot run bot!
)
goto end

:unknown
echo Unknown command: %1
echo Run "build.bat help" for available commands
goto end

:end
