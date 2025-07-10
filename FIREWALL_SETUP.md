# Windows Firewall Configuration for WebRTC Server

Your WebRTC server requires specific ports to be open in Windows Firewall to function properly.

## Required Ports

| Port | Protocol | Service         | Description                                   |
| ---- | -------- | --------------- | --------------------------------------------- |
| 443  | TCP      | HTTPS/WebSocket | Signaling server for WebRTC peer coordination |
| 3478 | UDP      | STUN/TURN       | Standard WebRTC NAT traversal (UDP)           |
| 3478 | TCP      | STUN/TURN       | Standard WebRTC NAT traversal (TCP fallback)  |
| 5349 | TCP      | STUN/TURN TLS   | Secure WebRTC NAT traversal (TLS)             |

## Important Notes

- **STUN/TURN Server**: The server on port 3478 provides both STUN discovery and TURN relay services
- **TLS Support**: Port 5349 is used for encrypted STUN/TURN connections when SSL certificates are available
- **Signaling**: Port 443 handles WebSocket connections for WebRTC signaling

## Automatic Configuration

### Option 1: Use the provided script (Recommended)

1. Right-click on `helpful-scripts\configure-firewall.bat`
2. Select "Run as administrator"
3. Follow the prompts

### Option 2: Use PowerShell script directly

1. Right-click on PowerShell
2. Select "Run as administrator"
3. Navigate to this directory
4. Run: `.\helpful-scripts\configure-firewall.ps1`

## Manual Configuration

If you prefer to configure the firewall manually:

### Using Windows Defender Firewall GUI

1. **Open Windows Defender Firewall**

   - Press `Win + R`, type `wf.msc`, press Enter
   - Or search for "Windows Defender Firewall" in Start menu

2. **Create Inbound Rules**

   - Click "Inbound Rules" in the left panel
   - Click "New Rule..." in the right panel
   - Select "Port" and click Next
   - Select "TCP" and enter "443" for port, click Next
   - Select "Allow the connection", click Next
   - Select all profiles (Domain, Private, Public), click Next
   - Name: "WebRTC Signaling HTTPS", click Finish
   - Repeat for other ports (3478 UDP/TCP, 5349 TCP)

3. **Create Outbound Rules** (Optional but recommended)
   - Click "Outbound Rules" in the left panel
   - Follow the same process as above for outbound rules

### Using Command Line

Open Command Prompt as Administrator and run:

```cmd
REM Inbound Rules
netsh advfirewall firewall add rule name="WebRTC Signaling HTTPS" dir=in action=allow protocol=TCP localport=443
netsh advfirewall firewall add rule name="WebRTC STUN/TURN UDP" dir=in action=allow protocol=UDP localport=3478
netsh advfirewall firewall add rule name="WebRTC STUN/TURN TCP" dir=in action=allow protocol=TCP localport=3478
netsh advfirewall firewall add rule name="WebRTC STUN/TURN TLS" dir=in action=allow protocol=TCP localport=5349

REM Outbound Rules
netsh advfirewall firewall add rule name="WebRTC Signaling HTTPS Outbound" dir=out action=allow protocol=TCP localport=443
netsh advfirewall firewall add rule name="WebRTC STUN/TURN UDP Outbound" dir=out action=allow protocol=UDP localport=3478
netsh advfirewall firewall add rule name="WebRTC STUN/TURN TCP Outbound" dir=out action=allow protocol=TCP localport=3478
netsh advfirewall firewall add rule name="WebRTC STUN/TURN TLS Outbound" dir=out action=allow protocol=TCP localport=5349
```

## Verification

After configuration, you can verify the rules were created:

1. Open Windows Defender Firewall
2. Look for rules starting with "WebRTC" in both Inbound and Outbound Rules
3. Ensure they are enabled (green checkmark)

## Troubleshooting

### Common Issues

1. **"Access Denied" errors**

   - Make sure you're running as Administrator

2. **Port already in use**

   - Check if another application is using the ports
   - Use `netstat -an | findstr :443` to check port usage

3. **Firewall rules not working**
   - Ensure Windows Defender Firewall is enabled
   - Check if third-party antivirus is blocking the connections

### Testing Port Accessibility

You can test if ports are accessible using:

```cmd
telnet localhost 443
telnet localhost 3478
telnet localhost 5349
```

Or use online port checkers to test from external networks.

## Removing Firewall Rules

To remove the WebRTC firewall rules:

```cmd
netsh advfirewall firewall delete rule name="WebRTC*"
```

Or use PowerShell:

```powershell
Remove-NetFirewallRule -DisplayName "WebRTC*"
```

## Security Notes

- These rules allow incoming connections on specific ports
- Only use these rules on trusted networks
- Consider using more restrictive rules for production environments
- Monitor firewall logs for suspicious activity

## Server Configuration

The server automatically detects which protocols to enable:

- **UDP 3478**: Always enabled (standard STUN/TURN)
- **TCP 3478**: Enabled with `-enable-tcp=true` (fallback for UDP-blocked networks)
- **TLS 5349**: Enabled with `-enable-tls=true` (requires SSL certificates in `certs/` directory)
- **HTTPS 443**: Always enabled for WebSocket signaling

## Client Configuration

Clients should configure their ICE servers like this:

```javascript
const iceServers = [
  // STUN discovery (same server as TURN)
  {
    urls: "stun:YOUR_SERVER_IP:3478",
  },
  // TURN relay (same server as STUN)
  {
    urls: "turn:YOUR_SERVER_IP:3478",
    username: "your_username",
    credential: "your_password",
  },
  // TURN over TCP (fallback)
  {
    urls: "turn:YOUR_SERVER_IP:3478?transport=tcp",
    username: "your_username",
    credential: "your_password",
  },
  // TURN over TLS (secure)
  {
    urls: "turns:YOUR_SERVER_IP:5349",
    username: "your_username",
    credential: "your_password",
  },
];
```
