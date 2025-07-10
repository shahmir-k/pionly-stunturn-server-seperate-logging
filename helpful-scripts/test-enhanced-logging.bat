@echo off
echo Testing Enhanced STUN/TURN Logging
echo ===================================
echo.

echo Starting STUN/TURN server with enhanced logging...
echo.

REM Start the server with enhanced logging
start "STUN/TURN Server" cmd /k "go run ..\main.go -public-ip 127.0.0.1 -turn-users testuser=testpass -separate-logs"

echo Waiting for server to start...
timeout /t 3 /nobreak >nul

echo.
echo Testing STUN discovery...
echo Using: stun:127.0.0.1:3478
echo.

REM Test STUN binding request using curl (if available)
curl -s --stun-host 127.0.0.1 --stun-port 3478 >nul 2>&1
if %errorlevel% equ 0 (
    echo STUN test completed
) else (
    echo STUN test failed or curl not available
)

echo.
echo Testing TURN allocation...
echo Using: turn:127.0.0.1:3478
echo Username: testuser
echo Password: testpass
echo.

REM Test TURN allocation using curl (if available)
curl -s --turn-host 127.0.0.1 --turn-port 3478 --turn-username testuser --turn-password testpass >nul 2>&1
if %errorlevel% equ 0 (
    echo TURN test completed
) else (
    echo TURN test failed or curl not available
)

echo.
echo Testing with WebRTC test page...
echo Opening test page in browser...
echo.

REM Open the test HTML page
start test-stunturn.html

echo.
echo ===================================
echo Enhanced Logging Test Complete
echo ===================================
echo.
echo Check the STUN/TURN log window for detailed logging including:
echo - STUN_BINDING_REQUEST/RESPONSE messages
echo - TURN_ALLOCATE_REQUEST/RESPONSE messages  
echo - Authentication attempts and results
echo - Connection events with emojis and clear formatting
echo.
echo The enhanced logging should now show:
echo - 🔍 STUN messages (NAT discovery)
echo - 🔄 TURN messages (relay allocation)
echo - ✅ Authentication successes
echo - ❌ Authentication failures
echo - 🔗 New connections
echo - 📥📤 Data transfer events
echo.
pause 