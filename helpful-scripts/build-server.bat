@echo off
echo Building WebRTC Server...
echo.

REM Check if Go is installed
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: Go is not installed or not in PATH
    echo Please install Go from https://golang.org/dl/
    pause
    exit /b 1
)

echo Go version:
go version
echo.

REM Clean previous builds
echo Cleaning previous builds...
if exist go-server.exe del go-server.exe
echo.

REM Download dependencies
echo Downloading dependencies...
go mod tidy
if %errorlevel% neq 0 (
    echo ERROR: Failed to download dependencies
    pause
    exit /b 1
)
echo.

REM Build the server
echo Building go-server.exe...
go build -o go-server.exe -ldflags="-s -w" .
if %errorlevel% neq 0 (
    echo ERROR: Build failed
    pause
    exit /b 1
)
echo.

REM Verify the build
if exist go-server.exe (
    echo ✓ Build successful!
    echo.
    echo File: go-server.exe
    echo Size: 
    for %%A in (go-server.exe) do echo %%~zA bytes
    echo.
    echo You can now run the server using:
    echo   go-server.exe -public-ip=YOUR_PUBLIC_IP -turn-users="username=password"
    echo.
    echo Or use the start-server.bat script with your configuration.
) else (
    echo ERROR: Build failed - go-server.exe not found
)

pause 