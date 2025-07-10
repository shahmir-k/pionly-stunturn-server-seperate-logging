# STUN/TURN Server Enhanced Logging System

## Overview

This document explains the enhanced logging system implemented for the STUN/TURN server to provide comprehensive visibility into user connections and protocol activity.

## Why You Weren't Seeing Logs Before

### 1. **STUN vs TURN Logging Difference**

- **STUN requests** (NAT discovery) are typically **unauthenticated** and don't trigger the `AuthHandler`
- **TURN requests** (relay allocation) require authentication and trigger the `AuthHandler`
- Your server was only logging **TURN authentication events**, not STUN discovery events

### 2. **Missing STUN Request Logging**

- The Pion TURN library doesn't provide built-in callbacks for STUN binding requests
- STUN requests were happening but not being logged

### 3. **Raw Packet Data vs Protocol Messages**

- Previous logging only showed raw packet data without identifying STUN/TURN message types
- No distinction between STUN_BINDING_REQUEST, TURN_ALLOCATE_REQUEST, etc.

## What's New - Enhanced Logging System

### 1. **STUNTurnLogger - Comprehensive Protocol Logging**

- **`STUNTurnLogger`**: Centralized logging system with emoji-based visual indicators
- **`LogSTUNRequest/Response`**: Logs STUN binding requests and responses
- **`LogTURNRequest/Response`**: Logs TURN allocation, refresh, send, and data requests
- **`LogAuthentication`**: Logs authentication attempts with success/failure indicators
- **`LogConnection`**: Logs new TCP connections
- **`LogRelayAllocation`**: Logs relay address allocation events
- **`LogDataTransfer`**: Logs data transfer between clients

### 2. **Enhanced Packet Handlers**

- **`LoggingPacketConn`**: Wraps UDP connections to log all STUN/TURN packets with protocol identification
- **`LoggingListener`**: Wraps TCP listeners to log connection events
- **`LoggingConn`**: Wraps TCP connections to log data transfer

### 3. **STUN/TURN Message Parsing**

- **`parseSTUNTURNMessage`**: Parses raw packet data to identify STUN/TURN message types
- **`getMessageTypeName`**: Converts message type codes to human-readable names
- **`isSTUNMessage/isTURNMessage`**: Distinguishes between STUN and TURN messages

### 4. **Enhanced Authentication Handler**

- **`createEnhancedAuthHandler`**: Provides detailed authentication logging with emoji indicators
- Logs both successful and failed authentication attempts
- Shows username, source address, and realm information

## Logging Output Examples

### STUN Discovery Logs

```
🔍 STUN STUN_BINDING_REQUEST from 192.168.1.100:54321
🔍 STUN STUN_BINDING_RESPONSE to 192.168.1.100:54321
```

### TURN Relay Logs

```
🔄 TURN TURN_ALLOCATE_REQUEST from 192.168.1.100:54321 (user: testuser)
✅ Authentication SUCCESS for user 'testuser' from 192.168.1.100:54321
🔄 TURN TURN_ALLOCATE_RESPONSE to 192.168.1.100:54321 (user: testuser)
🌐 Relay allocated for user 'testuser' from 192.168.1.100:54321 -> 203.0.113.1:49152
```

### Connection Logs

```
🔗 New TCP connection from 192.168.1.100:54321
📥 Received 32 bytes from 192.168.1.100:54321
📤 Sent 28 bytes to 192.168.1.100:54321
```

### Authentication Logs

```
🔐 Authentication attempt for user: testuser from 192.168.1.100:54321 (realm: pion.ly)
✅ Authentication SUCCESS for user 'testuser' from 192.168.1.100:54321
```

## How to Use the Enhanced Logging

### 1. **Start the Server with Enhanced Logging**

```bash
go run main.go -public-ip YOUR_IP -turn-users "user=pass" -separate-logs
```

### 2. **Monitor the Logs**

- The server will open separate log windows for STUN/TURN and signaling
- Watch for the emoji indicators to quickly identify different types of activity
- Look for STUN_BINDING_REQUEST/RESPONSE for NAT discovery
- Look for TURN_ALLOCATE_REQUEST/RESPONSE for relay allocation

### 3. **Test the Logging**

```bash
# Run the enhanced logging test
test-enhanced-logging.bat
```

## Troubleshooting

### If You Still Don't See STUN Logs

1. **Check if clients are actually sending STUN requests**

   - Use the test HTML page to verify STUN discovery
   - Check browser developer tools for WebRTC activity

2. **Verify server is listening on correct ports**

   - UDP: 3478 (STUN/TURN)
   - TCP: 3478 (STUN/TURN)
   - TLS: 5349 (STUN/TURN)

3. **Check firewall settings**
   - Ensure ports 3478 and 5349 are open
   - Allow both UDP and TCP traffic

### If You Don't See TURN Logs

1. **Verify authentication credentials**

   - Check username/password match server configuration
   - Ensure realm setting is correct

2. **Check client configuration**
   - Verify TURN server URL format: `turn:server:port`
   - Ensure username and credential are provided

## Message Types Identified

### STUN Messages

- `STUN_BINDING_REQUEST` (0x0001)
- `STUN_BINDING_RESPONSE` (0x0101)
- `STUN_BINDING_ERROR_RESPONSE` (0x0111)

### TURN Messages

- `TURN_ALLOCATE_REQUEST` (0x0003)
- `TURN_ALLOCATE_RESPONSE` (0x0103)
- `TURN_ALLOCATE_ERROR_RESPONSE` (0x0113)
- `TURN_REFRESH_REQUEST` (0x0004)
- `TURN_REFRESH_RESPONSE` (0x0104)
- `TURN_REFRESH_ERROR_RESPONSE` (0x0114)
- `TURN_SEND_REQUEST` (0x0006)
- `TURN_SEND_RESPONSE` (0x0106)
- `TURN_SEND_ERROR_RESPONSE` (0x0116)
- `TURN_DATA_REQUEST` (0x0007)
- `TURN_DATA_RESPONSE` (0x0107)
- `TURN_DATA_ERROR_RESPONSE` (0x0117)
- `TURN_CREATE_PERMISSION_REQUEST` (0x0008)
- `TURN_CREATE_PERMISSION_RESPONSE` (0x0108)
- `TURN_CREATE_PERMISSION_ERROR_RESPONSE` (0x0118)
- `TURN_CHANNEL_BIND_REQUEST` (0x0009)
- `TURN_CHANNEL_BIND_RESPONSE` (0x0109)
- `TURN_CHANNEL_BIND_ERROR_RESPONSE` (0x0119)

## Benefits of Enhanced Logging

1. **Clear Protocol Identification**: Distinguishes between STUN and TURN messages
2. **Visual Indicators**: Emoji-based logging for quick identification
3. **Comprehensive Coverage**: Logs both authenticated and unauthenticated requests
4. **Detailed Authentication**: Shows success/failure with user information
5. **Connection Tracking**: Monitors TCP connections and data transfer
6. **Relay Monitoring**: Tracks relay allocation and usage

## Next Steps

1. **Run the enhanced logging test** to verify everything is working
2. **Monitor the logs** during actual WebRTC calls
3. **Use the logging** to debug connection issues
4. **Extend the logging** if you need additional information

The enhanced logging system now provides complete visibility into STUN/TURN server activity, making it much easier to debug and monitor WebRTC connections.
