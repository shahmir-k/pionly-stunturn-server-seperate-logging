# Unified WebRTC Server

[![Go Version](https://img.shields.io/badge/go-1.21%2B-blue?logo=go)](https://golang.org/dl/)
[![License: MIT](https://img.shields.io/badge/license-MIT-green.svg)](./LICENSE)
[![Build Status](https://img.shields.io/badge/build-manual-lightgrey)](#building-from-source)
[![WebRTC](https://img.shields.io/badge/WebRTC-Native-blue?logo=webrtc)](https://webrtc.org/)
[![Pion](https://img.shields.io/badge/Pion-TURN%2FSTUN-orange?logo=go)](https://pion.ly/)
[![WebSocket](https://img.shields.io/badge/WebSocket-Signaling-yellow?logo=websocket)](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket)
[![STUN](https://img.shields.io/badge/STUN-NAT%20Discovery-cyan)](https://tools.ietf.org/html/rfc5389)
[![TURN](https://img.shields.io/badge/TURN-Relay%20Server-purple)](https://tools.ietf.org/html/rfc5766)
[![Gorilla WebSocket](https://img.shields.io/badge/Gorilla-WebSocket-red?logo=go)](https://github.com/gorilla/websocket)
[![Cross Platform](https://img.shields.io/badge/Cross%20Platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey)](https://golang.org/)
[![Protocols](https://img.shields.io/badge/Protocols-UDP%20%7C%20TCP%20%7C%20TLS-green)](https://tools.ietf.org/html/rfc5766)

A unified WebRTC server that combines signaling and STUN/TURN functionality in a single application, with enhanced logging, monitoring, and cross-platform support.

---

## 🚀 Features

- **STUN/TURN Server**: Unified server for STUN discovery and TURN relay
- **WebSocket Signaling**: Handles WebRTC signaling for call setup
- **Multi-Protocol**: UDP, TCP, and TLS variants
- **HTTPS Support**: Secure WebSocket connections with SSL certificates
- **Multi-threaded**: Multiple UDP listeners for load balancing
- **Authentication**: TURN server with username/password
- **Separate Logging**: Optional separate log files for STUN/TURN and signaling
- **Real-time Monitoring**: Live log viewing in separate terminal windows
- **Graceful Shutdown**: Proper cleanup on server termination

---

## 📦 Prerequisites

- [Go 1.21+](https://golang.org/dl/)
- Public IP address accessible from the internet
- SSL certificates (optional, for HTTPS)

---

## ⚡ Quick Start

1. **Edit the startup script:**

    If you don't set the PUBLIC_IP then the server will try to automatically detect your IP. So you can skip this step!

   - Windows: Edit `helpful-scripts\start-server.bat` and set `PUBLIC_IP=YOUR_ACTUAL_IP`
   - Linux/macOS: Edit `helpful-scripts/start-server.sh` and set `PUBLIC_IP=YOUR_ACTUAL_IP`

2. **Run the server:**


   - Windows: 
   
     Open the project folder in Powershell and run
   ```shell
   helpful-scripts\start-server.bat
   ```
   - Linux/macOS: 
     
     Open the project in Terminal and run
   ```shell
   chmod +x helpful-scripts/start-server.sh
   ./helpful-scripts/start-server.sh
   ```

---

## 🛠️ Manual Usage

```bash
# Build the server
go build -o go-server.exe .

# Start with custom parameters
./go-server.exe -public-ip=YOUR_PUBLIC_IP -turn-users="user1=pass1,user2=pass2" -realm="yourdomain.com" -thread-num=2
```

---

## ⚙️ Configuration

### Required Parameters

- `-public-ip`: Your server's public IP address (required)
- `-turn-users`: TURN users in format "username=password" (optional, has default)
- `-realm`: TURN server realm (optional, defaults to "pion.ly")
- `-thread-num`: Number of UDP listener threads (optional, defaults to 1)

### Optional Parameters

- `-enable-tcp`: Enable TCP fallback (default: true)
- `-enable-tls`: Enable TLS encryption (default: true)
- `-separate-logs`: Enable separate logging (default: true)
- `-stun-turn-log`: Custom STUN/TURN log file (default: "stun-turn.log")
- `-signaling-log`: Custom signaling log file (default: "signaling.log")

### SSL Certificates (Optional)

Place your SSL certificates in the `certs/` directory:

- `certs/fullchain.pem`: Your SSL certificate chain
- `certs/privkey.pem`: Your private key

---

## 🌐 Endpoints

- **Signaling Server:**
  - HTTP: `http://your-domain:443/signal`
  - HTTPS: `https://your-domain:443/signal` (if SSL certificates are present)
- **STUN/TURN Server:**
  - UDP: `your-domain:3478` (STUN discovery + TURN relay)
  - TCP: `your-domain:3478` (fallback)
  - TLS: `your-domain:5349` (secure)

---

## 🧊 WebRTC Client Configuration

```js
const iceServers = [
  { urls: "stun:your-domain:3478" },
  {
    urls: "turn:your-domain:3478",
    username: "1ac96ad0a8374103e5c58441",
    credential: "drTJQZjbVFKpcXfn",
  },
  {
    urls: "turn:your-domain:3478?transport=tcp",
    username: "1ac96ad0a8374103e5c58441",
    credential: "drTJQZjbVFKpcXfn",
  },
  {
    urls: "turns:your-domain:5349",
    username: "1ac96ad0a8374103e5c58441",
    credential: "drTJQZjbVFKpcXfn",
  },
];
const signalingUrl = "wss://your-domain:443/signal";
```

---

## 📊 Monitoring & Logging

- **Separate Logging:**
  ```sh
  ./go-server.exe -public-ip=YOUR_IP -separate-logs=true -stun-turn-log="stun-turn.log" -signaling-log="signaling.log"
  ```
- **Real-time Monitoring:**
  - Windows: `helpful-scripts\monitor-webrtc.bat YOUR_IP "username=password" powershell`
  - Linux/macOS: `./helpful-scripts/monitor-webrtc.sh YOUR_IP "username=password"`

---

## 🔒 Firewall

- **Required Ports:**
  - 443 (TCP): HTTPS/WebSocket Signaling
  - 3478 (UDP/TCP): STUN/TURN
  - 5349 (TCP): STUN/TURN TLS
- **Scripts:**
  - Windows: `helpful-scripts\configure-firewall.bat`
  - PowerShell: `helpful-scripts\configure-firewall.ps1`

---

## 📁 Project Structure

```
go-backend-seperateLogging/
├── main.go
├── webrtc/
│   ├── handler.go
│   ├── service.go
│   └── models.go
├── helpful-scripts/
│   ├── build-server.bat
│   ├── build-server.ps1
│   ├── start-server.bat
│   ├── monitor-webrtc.bat
│   ├── configure-firewall.bat
│   └── generate-certs.bat
├── certs/
│   ├── fullchain.pem
│   └── privkey.pem
├── go.mod
├── go.sum
└── README.md
```

---

## 🧰 Troubleshooting

- **"public-ip is required"**: Set the `-public-ip` flag to your server's public IP
- **Port already in use**: Ensure ports 443, 3478, and 5349 are not used by other services
- **TURN authentication fails**: Verify username/password in client configuration
- **SSL certificate errors**: Ensure certificate files are in the `certs/` directory
- **Monitoring windows don't open**: Check if PowerShell/xterm is available

---

## 🏗️ Building from Source

```sh
# Windows
helpful-scripts\build-server.bat
# PowerShell
helpful-scripts\build-server.ps1
# Manual
cd go-backend-seperateLogging
go mod tidy
go build -o go-server.exe -ldflags="-s -w" .
```

---

## 🔐 Security & Performance

- **Use HTTPS** in production
- **Strong TURN passwords**
- **Restrict firewall ports**
- **Keep Go and dependencies updated**
- **Restrict access to log files**
- **Increase `-thread-num` for high load**
- **Ensure adequate CPU/memory**
- **Use high-bandwidth for TURN relay**

---

## 📚 Documentation

- [BUILD_INSTRUCTIONS.md](./BUILD_INSTRUCTIONS.md)
- [MONITORING_README.md](./MONITORING_README.md)
- [LOGGING_IMPROVEMENTS.md](./LOGGING_IMPROVEMENTS.md)
- [FIREWALL_SETUP.md](./FIREWALL_SETUP.md)

---

## 📄 License

MIT

---

**Note:** This server provides both STUN and TURN services on the same port (3478). Every TURN server inherently supports STUN functionality according to RFC 5766. Configure your clients to use the same server for both STUN discovery and TURN relay services.
