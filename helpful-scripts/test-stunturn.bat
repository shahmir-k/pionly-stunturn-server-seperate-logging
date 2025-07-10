@echo off
echo Testing STUN/TURN Server Connectivity
echo =====================================

echo.
echo 1. Testing STUN discovery (UDP)...
echo Using: stun:localhost:3478
echo.

REM Test STUN binding request using curl (if available)
curl -s --stun-host localhost --stun-port 3478 >nul 2>&1
if %errorlevel% equ 0 (
    echo STUN test completed
) else (
    echo STUN test failed or curl not available
)

echo.
echo 2. Testing TURN allocation (UDP)...
echo Using: turn:localhost:3478
echo Username: 1ac96ad0a8374103e5c58441
echo Password: drTJQZjbVFKpcXfn
echo.

REM Test TURN allocation using curl (if available)
curl -s --turn-host localhost --turn-port 3478 --turn-username 1ac96ad0a8374103e5c58441 --turn-password drTJQZjbVFKpcXfn >nul 2>&1
if %errorlevel% equ 0 (
    echo TURN test completed
) else (
    echo TURN test failed or curl not available
)

echo.
echo 3. Testing TCP STUN/TURN...
echo Using: stun:localhost:3478?transport=tcp
echo.

REM Test TCP STUN
curl -s --stun-host localhost --stun-port 3478 --stun-transport tcp >nul 2>&1
if %errorlevel% equ 0 (
    echo TCP STUN test completed
) else (
    echo TCP STUN test failed or curl not available
)

echo.
echo 4. Manual testing instructions:
echo.
echo To test manually, you can use:
echo - WebRTC test tools: https://webrtc.github.io/samples/src/content/peerconnection/constraints/
echo - STUN/TURN testing tools: https://webrtc.github.io/samples/src/content/peerconnection/ice/
echo - Browser developer tools to monitor WebRTC connections
echo.
echo Check the stun-turn.log file for detailed logging information.
echo.
echo Press any key to exit...
pause >nul 