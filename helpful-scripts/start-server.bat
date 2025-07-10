@echo off
echo Starting Unified WebRTC Server...
echo.

REM Get the public IP address (you may need to modify this for your setup)
REM For development, you can use your local IP or a public IP
set PUBLIC_IP=YOUR_PUBLIC_IP_HERE

REM Optional: Set custom TURN port (default is 3478)
set TURN_PORT=3478

REM Check if public IP is set
if "%PUBLIC_IP%"=="YOUR_PUBLIC_IP_HERE" (
    echo ERROR: Please set your public IP address in this script
    echo Edit this file and replace YOUR_PUBLIC_IP_HERE with your actual public IP
    pause
    exit /b 1
)

echo Using public IP: %PUBLIC_IP%
echo Using TURN port: %TURN_PORT%
echo.

REM Download dependencies
echo Downloading dependencies...
go mod tidy

REM Build the server
echo Building server...
go build -o go-server.exe .

REM Start the server
echo Starting server...
echo.
echo Server will be available at:
echo - Signaling: https://%PUBLIC_IP%:443/signal (HTTPS)
echo - STUN/TURN: %PUBLIC_IP%:%TURN_PORT% (UDP/TCP)
echo - STUN/TURN TLS: %PUBLIC_IP%:5349 (TLS)
echo.
echo Press Ctrl+C to stop the server
echo.

go-server.exe -public-ip=%PUBLIC_IP% -turn-port=%TURN_PORT% -turn-users="1ac96ad0a8374103e5c58441=drTJQZjbVFKpcXfn" -realm="yourdomain.com" -thread-num=4 -separate-logs=true

pause 