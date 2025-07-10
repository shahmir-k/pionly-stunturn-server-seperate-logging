@echo off
echo WebRTC Server Firewall Configuration
echo ====================================
echo.

echo This script will configure Windows Firewall to allow your WebRTC server
echo to function properly by opening the necessary ports.
echo.

echo Required ports:
echo • Port 443 (TCP) - HTTPS/WebSocket Signaling Server
echo • Port 3478 (UDP/TCP) - STUN/TURN Server  
echo • Port 5349 (TCP) - STUN/TURN TLS Server
echo.

set /p confirm="Do you want to continue? (y/n): "
if /i not "%confirm%"=="y" (
    echo Configuration cancelled.
    pause
    exit /b 0
)

echo.
echo Running firewall configuration...
echo.

REM Run PowerShell script with administrator privileges
powershell -ExecutionPolicy Bypass -File "%~dp0configure-firewall.ps1"

echo.
echo Firewall configuration completed!
echo You can now run your WebRTC server.
echo.
pause 