# WebRTC Server Monitoring Guide

This guide explains how to monitor and debug the WebRTC server using the built-in logging and monitoring features.

## Overview

The WebRTC server provides comprehensive monitoring capabilities with:

- **Separate Logging**: Independent log files for STUN/TURN and signaling services
- **Real-time Monitoring**: Live log viewing in separate terminal windows
- **Structured Logging**: Clear prefixes and timestamps for easy filtering
- **Cross-platform Support**: Windows PowerShell and Unix terminal monitoring

## Logging Features

### Separate Log Files

The server can log STUN/TURN and signaling activities to separate files:

```cmd
# Enable separate logging
go-server.exe -public-ip=YOUR_IP -separate-logs=true -stun-turn-log="stun-turn.log" -signaling-log="signaling.log"
```

### Log File Contents

**STUN/TURN Log (`stun-turn.log`):**

- TURN authentication attempts (success/failure)
- Relay allocation and deallocation
- Connection events and statistics
- Protocol-specific events (UDP/TCP/TLS)

**Signaling Log (`signaling.log`):**

- WebSocket connection events
- User registration and session management
- SDP offer/answer exchange
- ICE candidate sharing
- Call state changes (join, call, hangup)

### Log Format

```
[STUN/TURN] 2024-01-15 10:30:45 TURN authentication successful for user: alice from 192.168.1.100:12345
[SIGNALING] 2024-01-15 10:30:46 WebSocket connection established from 192.168.1.100:54321
```

## Monitoring Methods

### Method 1: Separate Log Files (Default)

1. **Start the server with separate logging:**

   ```cmd
   go-server.exe -public-ip=YOUR_IP -separate-logs=true -stun-turn-log="stun-turn.log" -signaling-log="signaling.log"
   ```

2. **Monitor in separate terminals:**

   **Terminal 1 (STUN/TURN):**

   ```cmd
   type stun-turn.log
   ```

   **Terminal 2 (Signaling):**

   ```cmd
   type signaling.log
   ```

### Method 2: PowerShell Windows (Windows Only)

Use the monitoring script to open separate PowerShell windows:

```cmd
helpful-scripts\monitor-webrtc.bat YOUR_IP "username=password" powershell
```

This opens:

- **STUN/TURN Monitor Window**: Real-time STUN/TURN log viewing
- **Signaling Monitor Window**: Real-time signaling log viewing

### Method 3: Unix Terminal Windows (Linux/macOS)

The server automatically opens `xterm` windows for monitoring:

```bash
./go-server -public-ip=YOUR_IP -separate-logs=true
```

## Monitoring Scripts

### monitor-webrtc.bat

Windows batch script for easy monitoring setup:

```cmd
# Basic usage
helpful-scripts\monitor-webrtc.bat YOUR_IP

# With custom TURN users
helpful-scripts\monitor-webrtc.bat YOUR_IP "user1=pass1,user2=pass2"

# With specific method
helpful-scripts\monitor-webrtc.bat YOUR_IP "user1=pass1" powershell
```

**Available Methods:**

- `logs` - Separate log files (default)
- `powershell` - PowerShell monitoring windows

### monitor-webrtc.sh

Unix shell script for Linux/macOS monitoring:

```bash
# Make executable
chmod +x helpful-scripts/monitor-webrtc.sh

# Run monitoring
./helpful-scripts/monitor-webrtc.sh YOUR_IP "user1=pass1"
```

## Real-time Monitoring Features

### Automatic Window Management

- **Windows**: PowerShell windows with custom titles
- **Unix**: xterm windows with tail -f commands
- **Graceful Shutdown**: Automatic window cleanup on server exit

### Log Rotation

- **Fresh Start**: Clears existing log files on server startup
- **Continuous Writing**: Real-time log updates
- **File Management**: Proper file permissions and handling

### Cross-platform Compatibility

- **Windows**: PowerShell with custom monitoring scripts
- **Linux/macOS**: xterm with tail commands
- **Fallback**: Single logger to stdout if monitoring fails

## Debugging with Logs

### Common Log Patterns

**Successful TURN Authentication:**

```
[STUN/TURN] TURN authentication successful for user: alice from 192.168.1.100:12345
```

**Failed TURN Authentication:**

```
[STUN/TURN] TURN authentication failed for user: unknown from 192.168.1.100:12345
```

**WebSocket Connection:**

```
[SIGNALING] WebSocket connection established from 192.168.1.100:54321
```

**Relay Allocation:**

```
[STUN/TURN] UDP STUN/TURN server 0 listening on 0.0.0.0:3478
```

### Troubleshooting with Logs

1. **Authentication Issues:**

   - Check STUN/TURN logs for authentication failures
   - Verify username/password in client configuration

2. **Connection Problems:**

   - Check signaling logs for WebSocket connection errors
   - Verify server IP and port configuration

3. **Protocol Issues:**

   - Check STUN/TURN logs for protocol-specific errors
   - Verify firewall configuration for required ports

4. **Performance Issues:**
   - Monitor log frequency for connection patterns
   - Check for repeated authentication failures

## Server Statistics

The server logs periodic statistics:

```
[STUN/TURN] === SERVER STATISTICS ===
[STUN/TURN] Time: 2024-01-15 10:30:45
[STUN/TURN] Active STUN/TURN servers: 3
[STUN/TURN] ========================
```

**Statistics Include:**

- Current timestamp
- Number of active STUN/TURN servers (UDP/TCP/TLS)
- Server health indicators

## Configuration Options

### Logging Flags

```cmd
# Enable separate logging
-separate-logs=true

# Custom log file paths
-stun-turn-log="custom-stun-turn.log"
-signaling-log="custom-signaling.log"

# Disable separate logging (single logger)
-separate-logs=false
```

### Monitoring Behavior

- **Auto-start**: Monitoring windows open automatically when `-separate-logs=true`
- **Graceful shutdown**: Windows close automatically on server exit
- **Error handling**: Fallback to single logger if monitoring fails

## Best Practices

### Production Monitoring

1. **Log Rotation**: Implement log rotation for production environments
2. **Log Analysis**: Use log analysis tools for pattern recognition
3. **Alerting**: Set up alerts for authentication failures or connection issues
4. **Backup**: Regularly backup log files for analysis

### Development Monitoring

1. **Real-time Viewing**: Use PowerShell/xterm windows for immediate feedback
2. **Separate Concerns**: Use separate logs to isolate issues
3. **Verbose Logging**: Enable detailed logging for debugging
4. **Pattern Recognition**: Learn common log patterns for quick issue identification

## Troubleshooting Monitoring

### Common Issues

1. **Monitoring Windows Don't Open:**

   - Check if PowerShell/xterm is available
   - Verify file permissions
   - Check for antivirus blocking

2. **Log Files Not Created:**

   - Verify write permissions in directory
   - Check disk space
   - Ensure server started successfully

3. **High CPU Usage:**
   - Monitor log file sizes
   - Check for excessive logging
   - Consider log rotation

### Performance Considerations

- **Log File Size**: Monitor log file growth
- **I/O Impact**: Separate logging has minimal performance impact
- **Memory Usage**: Logging uses minimal additional memory
- **Network Impact**: Logging doesn't affect network performance

## Integration with External Tools

### Log Aggregation

- **ELK Stack**: Send logs to Elasticsearch/Logstash/Kibana
- **Splunk**: Forward logs to Splunk for analysis
- **Cloud Logging**: Send to cloud logging services

### Monitoring Dashboards

- **Grafana**: Create dashboards from log data
- **Prometheus**: Export metrics for monitoring
- **Custom Tools**: Parse logs for custom metrics

## Security Considerations

- **Log Sensitivity**: Logs may contain IP addresses and connection patterns
- **Access Control**: Restrict access to log files
- **Retention Policy**: Implement log retention policies
- **Encryption**: Consider encrypting sensitive log data
