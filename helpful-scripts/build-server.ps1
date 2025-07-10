# WebRTC Server Build Script
# This script builds the go-server.exe from Go source files

Write-Host "Building WebRTC Server..." -ForegroundColor Green
Write-Host ""

# Check if Go is installed
try {
    $goVersion = go version 2>$null
    if ($LASTEXITCODE -ne 0) {
        throw "Go not found"
    }
    Write-Host "Go version: $goVersion" -ForegroundColor Cyan
} catch {
    Write-Host "ERROR: Go is not installed or not in PATH" -ForegroundColor Red
    Write-Host "Please install Go from https://golang.org/dl/" -ForegroundColor Yellow
    pause
    exit 1
}

Write-Host ""

# Clean previous builds
Write-Host "Cleaning previous builds..." -ForegroundColor Yellow
if (Test-Path "go-server.exe") {
    Remove-Item "go-server.exe" -Force
    Write-Host "✓ Removed old go-server.exe" -ForegroundColor Green
}
Write-Host ""

# Download dependencies
Write-Host "Downloading dependencies..." -ForegroundColor Yellow
go mod tidy
if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: Failed to download dependencies" -ForegroundColor Red
    pause
    exit 1
}
Write-Host "✓ Dependencies downloaded" -ForegroundColor Green
Write-Host ""

# Build the server
Write-Host "Building go-server.exe..." -ForegroundColor Yellow
$buildStart = Get-Date

# Build with optimizations
go build -o go-server.exe -ldflags="-s -w" .

if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: Build failed" -ForegroundColor Red
    pause
    exit 1
}

$buildEnd = Get-Date
$buildTime = $buildEnd - $buildStart
Write-Host "✓ Build completed in $($buildTime.TotalSeconds.ToString('F2')) seconds" -ForegroundColor Green
Write-Host ""

# Verify the build
if (Test-Path "go-server.exe") {
    $fileInfo = Get-Item "go-server.exe"
    $fileSize = [math]::Round($fileInfo.Length / 1MB, 2)
    
    Write-Host "✓ Build successful!" -ForegroundColor Green
    Write-Host ""
    Write-Host "File Details:" -ForegroundColor Cyan
    Write-Host "  Name: $($fileInfo.Name)" -ForegroundColor White
    Write-Host "  Size: $fileSize MB ($($fileInfo.Length) bytes)" -ForegroundColor White
    Write-Host "  Created: $($fileInfo.CreationTime)" -ForegroundColor White
    Write-Host ""
    Write-Host "Usage Examples:" -ForegroundColor Cyan
    Write-Host "  Basic usage:" -ForegroundColor White
    Write-Host "    .\go-server.exe -public-ip=YOUR_PUBLIC_IP -turn-users=`"username=password`"" -ForegroundColor Gray
    Write-Host ""
    Write-Host "  With separate logging:" -ForegroundColor White
    Write-Host "    .\go-server.exe -public-ip=YOUR_PUBLIC_IP -separate-logs=true -stun-turn-log=`"stun-turn.log`" -signaling-log=`"signaling.log`"" -ForegroundColor Gray
    Write-Host ""
    Write-Host "  Full options:" -ForegroundColor White
    Write-Host "    .\go-server.exe -public-ip=YOUR_PUBLIC_IP -turn-port=3478 -turn-users=`"username=password`" -realm=`"yourdomain.com`" -thread-num=4 -enable-tcp=true -enable-tls=true" -ForegroundColor Gray
    Write-Host ""
    Write-Host "  Or use the start-server.bat script with your configuration." -ForegroundColor Yellow
} else {
    Write-Host "ERROR: Build failed - go-server.exe not found" -ForegroundColor Red
}

Write-Host ""
Write-Host "Build process completed!" -ForegroundColor Green
pause 