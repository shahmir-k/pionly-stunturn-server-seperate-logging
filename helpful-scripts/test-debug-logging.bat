@echo off
echo Testing Enhanced STUN/TURN Logging with Debug Output
echo ====================================================
echo.

echo Starting STUN/TURN server with debug logging...
echo.

REM Start the server with enhanced logging
start "STUN/TURN Server Debug" cmd /k "go run ..\main.go -public-ip 127.0.0.1 -turn-users testuser=testpass -separate-logs"

echo Waiting for server to start...
timeout /t 5 /nobreak >nul

echo.
echo Opening test page to generate STUN/TURN traffic...
echo.

REM Open the test HTML page
start test-stunturn.html

echo.
echo ===================================
echo Debug Test Complete
echo ===================================
echo.
echo Check the STUN/TURN log window for:
echo - DEBUG: Packet header lines showing raw packet data
echo - DEBUG: STUN magic cookie found/not found messages
echo - DEBUG: Message type code and parsed message type
echo.
echo This will help us understand why message types aren't being identified.
echo.
pause 