@echo off
setlocal enabledelayedexpansion

REM WebRTC Server Monitoring Script for Windows
REM This script demonstrates how to run the WebRTC server with separate logging
REM and monitor STUN/TURN and signaling services in different terminals

REM Configuration
set "PUBLIC_IP=%~1"
set "TURN_USERS=%~2"
set "METHOD=%~3"

if "%PUBLIC_IP%"=="" set "PUBLIC_IP=YOUR_PUBLIC_IP"
if "%TURN_USERS%"=="" set "TURN_USERS=1ac96ad0a8374103e5c58441=drTJQZjbVFKpcXfn"
if "%METHOD%"=="" set "METHOD=logs"

set "LOG_DIR=..\logs"
set "STUN_TURN_LOG=%LOG_DIR%\stun-turn.log"
set "SIGNALING_LOG=%LOG_DIR%\signaling.log"

echo WebRTC Server Monitoring Setup
echo ==================================
echo Public IP: %PUBLIC_IP%
echo TURN Users: %TURN_USERS%
echo Log Directory: %LOG_DIR%
echo Method: %METHOD%
echo.

REM Check if public IP is set
if "%PUBLIC_IP%"=="YOUR_PUBLIC_IP" (
    echo Error: Please provide your public IP address as the first argument
    echo Usage: %0 ^<PUBLIC_IP^> [TURN_USERS] [METHOD]
    echo Example: %0 203.0.113.1 "user1=pass1,user2=pass2" logs
    echo.
    echo Available methods: logs, powershell
    exit /b 1
)

REM Create log directory
if not exist "%LOG_DIR%" mkdir "%LOG_DIR%"

REM Function to show usage
if "%METHOD%"=="help" goto :show_usage
if "%METHOD%"=="-h" goto :show_usage
if "%METHOD%"=="--help" goto :show_usage

REM Main script logic
if "%METHOD%"=="logs" goto :run_with_log_files
if "%METHOD%"=="powershell" goto :run_with_powershell

echo Unknown method: %METHOD%
goto :show_usage

:show_usage
echo.
echo Available monitoring methods:
echo.
echo 1. logs - Separate Log Files (default)
echo    Run: %0 %PUBLIC_IP% "%TURN_USERS%" logs
echo    Then in separate terminals:
echo    Terminal 1: type %STUN_TURN_LOG%
echo    Terminal 2: type %SIGNALING_LOG%
echo.
echo 2. powershell - PowerShell Windows
echo    Run: %0 %PUBLIC_IP% "%TURN_USERS%" powershell
echo    This will open PowerShell windows for monitoring
echo.
exit /b 0

:run_with_log_files
echo Starting WebRTC server with separate log files...
echo.
echo Server will be started in the background.
echo Logs will be written to:
echo   STUN/TURN: %STUN_TURN_LOG%
echo   Signaling: %SIGNALING_LOG%
echo.
echo To monitor in separate terminals, run:
echo   Terminal 1: type %STUN_TURN_LOG%
echo   Terminal 2: type %SIGNALING_LOG%
echo.
echo To stop the server:
echo   taskkill /f /im go-server.exe
echo.

REM Start the server in background
start /b go-server.exe -public-ip "%PUBLIC_IP%" -turn-users "%TURN_USERS%" -separate-logs=true -stun-turn-log="%STUN_TURN_LOG%" -signaling-log="%SIGNALING_LOG%"

echo Server started in the background.
echo.
echo Press Enter to stop the server...
pause >nul
echo Stopping server...
taskkill /f /im go-server.exe >nul 2>&1
echo Server stopped.
exit /b 0

:run_with_powershell
echo Setting up PowerShell windows for monitoring...
echo.

REM Start the server in background
start /b go-server.exe -public-ip "%PUBLIC_IP%" -turn-users "%TURN_USERS%" -separate-logs=true -stun-turn-log="%STUN_TURN_LOG%" -signaling-log="%SIGNALING_LOG%"

echo Server started in the background.
echo.

REM Open PowerShell windows for monitoring
echo Opening PowerShell windows for monitoring...
start powershell -NoExit -Command "Write-Host 'STUN/TURN Monitor - Press Ctrl+C to exit' -ForegroundColor Green; Get-Content '%STUN_TURN_LOG%' -Wait"
start powershell -NoExit -Command "Write-Host 'Signaling Monitor - Press Ctrl+C to exit' -ForegroundColor Green; Get-Content '%SIGNALING_LOG%' -Wait"

echo PowerShell windows opened for monitoring.
echo.
echo To stop the server:
echo   taskkill /f /im go-server.exe
echo.
echo Press Enter to stop the server and close monitoring windows...
pause >nul
echo Stopping server...
taskkill /f /im go-server.exe >nul 2>&1
taskkill /f /im powershell.exe >nul 2>&1
echo Server and monitoring windows stopped.
exit /b 0 