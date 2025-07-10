# Building go-server.exe

This guide explains how to build the `go-server.exe` file from the Go source code.

## Prerequisites

1. **Go Programming Language** (version 1.21 or later)

   - Download from: https://golang.org/dl/
   - Install and add to PATH

2. **Git** (optional, for version control)
   - Download from: https://git-scm.com/

## Quick Build

### Option 1: Use the provided build script (Recommended)

**Windows Batch:**

```cmd
helpful-scripts\build-server.bat
```

**PowerShell:**

```powershell
.\helpful-scripts\build-server.ps1
```

### Option 2: Manual build commands

1. **Open Command Prompt or PowerShell** in the `go-backend-seperateLogging` directory

2. **Download dependencies:**

   ```cmd
   go mod tidy
   ```

3. **Build the executable:**

   ```cmd
   go build -o go-server.exe .
   ```

4. **Build with optimizations (smaller file size):**
   ```cmd
   go build -o go-server.exe -ldflags="-s -w" .
   ```

## Build Options

### Basic Build

```cmd
go build -o go-server.exe .
```

### Optimized Build (Recommended)

```cmd
go build -o go-server.exe -ldflags="-s -w" .
```

- `-s`: Strips debug information
- `-w`: Strips DWARF symbol table
- Results in smaller executable size

### Cross-Platform Build

```cmd
# For Windows
GOOS=windows GOARCH=amd64 go build -o go-server.exe .

# For Linux
GOOS=linux GOARCH=amd64 go build -o go-server .

# For macOS
GOOS=darwin GOARCH=amd64 go build -o go-server .
```

### Debug Build

```cmd
go build -o go-server.exe -gcflags="-N -l" .
```

- Includes debug information
- Larger file size but better for debugging

## Source Files

The main source files for the build:

- `main.go` - Main entry point with comprehensive STUN/TURN server implementation
- `webrtc/handler.go` - WebSocket handler for signaling
- `webrtc/service.go` - WebRTC service implementation
- `webrtc/models.go` - Data structures
- `go.mod` - Go module definition
- `go.sum` - Dependency checksums

## Server Features

The built server provides:

- **STUN/TURN Server**: Unified server providing both STUN discovery and TURN relay services
- **WebSocket Signaling**: HTTP/HTTPS server for WebRTC peer coordination
- **Multi-Protocol Support**: UDP, TCP, and TLS variants for different network environments
- **Authentication**: TURN server with username/password authentication
- **Separate Logging**: Optional separate log files for STUN/TURN and signaling services
- **Graceful Shutdown**: Proper cleanup on server termination

## Troubleshooting

### Common Issues

1. **"go: command not found"**

   - Go is not installed or not in PATH
   - Reinstall Go and restart terminal

2. **"module not found"**

   - Run `go mod tidy` to download dependencies
   - Check internet connection

3. **"build failed"**

   - Check for syntax errors in Go files
   - Ensure all required files are present
   - Check Go version compatibility

4. **"permission denied"**
   - Run as administrator (Windows)
   - Check file permissions

### Build Verification

After building, verify the executable:

```cmd
# Check if file exists
dir go-server.exe

# Check file size (should be ~10MB)
# Check if it's a valid executable
go-server.exe --help
```

## File Structure After Build

```
go-backend-seperateLogging/
├── go-server.exe          # Built executable
├── main.go               # Main server implementation
├── webrtc/
│   ├── handler.go        # WebSocket handler
│   ├── service.go        # WebRTC service
│   └── models.go         # Data structures
├── go.mod
├── go.sum
├── helpful-scripts/
│   ├── build-server.bat  # Build script
│   ├── build-server.ps1  # PowerShell build script
│   └── start-server.bat  # Server startup script
└── certs/                # SSL certificates (optional)
    ├── fullchain.pem
    └── privkey.pem
```

## Running the Server

After building, you can run the server:

```cmd
# Basic usage
go-server.exe -public-ip=YOUR_PUBLIC_IP -turn-users="username=password"

# Or use the provided startup script
helpful-scripts\start-server.bat
```

## Build Scripts Details

### build-server.bat

- Windows batch script
- Checks Go installation
- Cleans previous builds
- Downloads dependencies
- Builds with optimizations
- Shows build status

### build-server.ps1

- PowerShell script
- More detailed output
- Build time measurement
- File size information
- Usage examples
- Better error handling

## Performance Notes

- **Build time**: ~5-30 seconds depending on system
- **Executable size**: ~10MB (optimized)
- **Memory usage**: ~50-100MB when running
- **CPU usage**: Low when idle, scales with connections

## Security Considerations

- The built executable contains all dependencies
- No external dependencies required at runtime
- Consider code signing for production deployment
- Verify the executable with antivirus software

## Server Configuration

The server supports various command-line options:

```cmd
# Basic configuration
go-server.exe -public-ip=YOUR_PUBLIC_IP

# Full configuration
go-server.exe -public-ip=YOUR_PUBLIC_IP -turn-users="user1=pass1,user2=pass2" -realm="yourdomain.com" -thread-num=4 -enable-tcp=true -enable-tls=true

# With separate logging
go-server.exe -public-ip=YOUR_PUBLIC_IP -separate-logs=true -stun-turn-log="stun-turn.log" -signaling-log="signaling.log"
```

## STUN/TURN Server Features

- **Dual Functionality**: Single server provides both STUN discovery and TURN relay
- **Protocol Variants**: UDP (standard), TCP (fallback), TLS (secure)
- **Authentication**: Username/password authentication for TURN services
- **Multi-threading**: Configurable number of listener threads
- **Relay Allocation**: Static relay address generation with public IP
