/*
WebRTC Unified Server - Complete Educational Implementation
==========================================================

This file implements a complete WebRTC infrastructure server that combines:
1. STUN Server (Session Traversal Utilities for NAT) - Helps clients discover their public IP
2. TURN Server (Traversal Using Relays around NAT) - Provides relay services when direct connection fails
3. WebRTC Signaling Server - Handles WebSocket connections for peer-to-peer signaling

WHAT IS WEBRTC?
===============
WebRTC (Web Real-Time Communication) enables peer-to-peer communication between browsers/devices
without requiring plugins or native apps. It's used for video calls, voice chat, file sharing, etc.

Think of WebRTC as a "direct phone line" between two devices:
- Traditional video calls: Device A → Server → Device B (server in the middle)
- WebRTC calls: Device A ←→ Device B (direct connection, server only helps establish)

WHY DO WE NEED STUN/TURN SERVERS?
=================================
Most devices are behind NAT (Network Address Translation) or firewalls, which makes
direct peer-to-peer connections impossible. STUN and TURN servers solve this:

STUN (Session Traversal Utilities for NAT):
STUNTURN (STUN + TURN) - IMPORTANT CONCEPT:
===========================================
Every TURN server is also a STUN server. This is a fundamental aspect of the TURN protocol
specification (RFC 5766). When you run a TURN server, it automatically handles both STUN and
TURN requests on the same port.

STUN functionality (NAT discovery):
- Helps clients discover their public IP address and NAT type
- Like asking "What's my address from the internet's perspective?"
- Free service, no relay of actual data
- Works in ~70% of cases

TURN (Traversal Using Relays around NAT):
- Acts as a relay server when direct peer-to-peer connection is impossible
- Like a "middleman" that forwards data between devices
- Paid service (uses bandwidth), but works in 100% of cases
- Required when STUN fails or for enterprise networks

REAL-WORLD ANALOGY:
===================
Imagine two people trying to meet:
- STUN: "I'm at 123 Main St, you're at 456 Oak Ave, let's meet directly"
- TURN: "We can't meet directly due to roadblocks, let's meet at the coffee shop (relay)"
- STUNTURN: "Let's try to meet directly first, but if that fails, we'll use the coffee shop"

ARCHITECTURE OVERVIEW:
=====================
This server provides multiple transport protocols to handle different network environments:

- UDP STUN/TURN: Standard protocol, works with most NATs (fastest)
- TCP STUN/TURN: Fallback for restrictive networks that block UDP (slower but more reliable)
- TLS STUN/TURN: Secure encrypted connections for enterprise environments (most secure)
- WebSocket Signaling: Coordinates connection establishment between peers

LEARNING OBJECTIVES:
===================
1. Understand WebRTC server architecture and NAT traversal
2. Learn how STUN and TURN servers work together. Learn how STUNTURN servers provide both STUN and TURN services
3. Implement secure authentication for TURN services
4. Handle multiple transport protocols (UDP/TCP/TLS)
5. Set up proper SSL/TLS encryption for production use
6. Create a scalable WebRTC infrastructure

HOW TO USE THIS SERVER:
=======================
1. Run the server with your public IP and TURN credentials
2. Configure your frontend clients with the ICE server information
3. Use the WebSocket signaling endpoint for peer coordination
4. Monitor logs to understand connection patterns

CLIENT CONFIGURATION:
====================
Clients should configure their ICE servers like this:
{
  urls: 'stun:YOUR_SERVER_IP:3478',        // STUN discovery (same server!)
  urls: 'turn:YOUR_SERVER_IP:3478',        // TURN relay (same server!)
  username: 'your_username',
  credential: 'your_password'
}

FRONTEND INTEGRATION:
====================
See the comprehensive examples at the bottom of this file for:
- JavaScript/WebRTC
- React Native
- Native Android
- Native iOS
- Go WebRTC client
- Protocol testing and troubleshooting

COMPREHENSIVE FRONTEND INTEGRATION EXAMPLES
===========================================

This section provides complete examples of how to integrate this WebRTC server
with different frontend platforms. Each example includes the complete ICE server
configuration and basic WebRTC setup.

IMPORTANT NOTES:
===============
- Replace 'YOUR_SERVER_IP' with your actual server IP address
- Replace 'YOUR_TURN_USERNAME' and 'YOUR_TURN_PASSWORD' with your TURN credentials
- For production, use HTTPS/WSS endpoints
- Test all protocol variants in your target environments
- Remember: TURN servers inherently support STUN functionality

1. JAVASCRIPT/WEBRTC CLIENT CONFIGURATION
=========================================

This is the most common use case for WebRTC applications in browsers.

// Complete ICE Server Configuration for WebRTC
const ICE_SERVERS = [
  // Primary STUN server for local network discovery (UDP)
  {
    urls: 'stun:YOUR_SERVER_IP:3478',
  },
  // STUN over TCP (fallback for UDP-blocked networks)
  {
    urls: 'stun:YOUR_SERVER_IP:3478?transport=tcp',
  },
  // STUN over TLS (secure encrypted discovery)
  {
    urls: 'stuns:YOUR_SERVER_IP:5349',
  },
  // TURN servers for NAT traversal when STUN fails
  {
    urls: 'turn:YOUR_SERVER_IP:3478',
    username: 'YOUR_TURN_USERNAME',
    credential: 'YOUR_TURN_PASSWORD',
  },
  // TURN over TCP (fallback for UDP-blocked networks)
  {
    urls: 'turn:YOUR_SERVER_IP:3478?transport=tcp',
    username: 'YOUR_TURN_USERNAME',
    credential: 'YOUR_TURN_PASSWORD',
  },
  // TURN over TLS (secure encrypted relay)
  {
    urls: 'turns:YOUR_SERVER_IP:5349',
    username: 'YOUR_TURN_USERNAME',
    credential: 'YOUR_TURN_PASSWORD',
  },
];


2. WHEN TO USE EACH PROTOCOL
============================

STUN UDP (stun:):
- Default choice for most networks
- Fastest and most efficient
- Works with most NAT types
- Use when: General purpose, good network conditions

STUN TCP (stun:?transport=tcp):
- Fallback when UDP is blocked
- Common in corporate networks
- Slightly slower than UDP
- Use when: UDP is blocked by firewall/NAT

STUN TLS (stuns:):
- Encrypted STUN communication
- Required for secure environments
- Works through strict firewalls
- Use when: Security is required, enterprise networks

TURN UDP (turn:):
- Standard relay service
- Most efficient relay option
- Works with most NAT types
- Use when: STUN fails, need relay

TURN TCP (turn:?transport=tcp):
- Relay fallback when UDP is blocked
- Slower but more reliable in restricted networks
- Use when: UDP is blocked, need relay

TURN TLS (turns:):
- Encrypted relay service
- Most secure option
- Required for WebRTC in browsers (HTTPS requirement)
- Use when: Security required, browser compatibility needed

3. NETWORK ENVIRONMENT CONSIDERATIONS
====================================

Corporate Networks:
- Often block UDP traffic
- Require TCP fallbacks
- May need TLS for security
- Recommended: Include TCP and TLS variants

Home Networks:
- Usually allow UDP
- STUN UDP works well
- TURN UDP for relay
- Recommended: UDP variants first, TCP as backup

Mobile Networks:
- Carrier-grade NATs
- Often block certain ports
- May require TCP/TLS
- Recommended: All protocol variants

Public WiFi:
- Often restrictive
- May block UDP
- Security concerns
- Recommended: TCP and TLS variants

4. TROUBLESHOOTING COMMON ISSUES
================================

a) Connection Fails:
   - Check if ports 3478, 5349, and 443 are open on the server
   - Verify SSL certificates are valid for HTTPS/TLS
   - Ensure public IP is correctly configured
   - Check firewall rules on both client and server

b) TURN Authentication Fails:
   - Verify username and password match server configuration
   - Check realm setting matches server
   - Ensure credentials are properly formatted

c) ICE Gathering Fails:
   - Test STUN connectivity first
   - Verify TURN server is accessible
   - Check network restrictions (corporate firewalls, etc.)

d) Signaling Connection Fails:
   - Verify WebSocket endpoint is accessible
   - Check SSL certificate validity
   - Ensure proper WebSocket URL format (wss:// for secure)

e) Protocol-Specific Issues:
   - UDP blocked: Try TCP variants
   - TLS fails: Check SSL certificates
   - TCP slow: Normal for TCP vs UDP
   - Mixed protocols: Ensure all variants are configured

5. PRODUCTION CONSIDERATIONS
============================

- Use strong, unique credentials for TURN authentication
- Implement proper SSL/TLS certificates (Let's Encrypt recommended)
- Monitor server performance and bandwidth usage
- Implement rate limiting to prevent abuse
- Consider using multiple TURN servers for redundancy
- Implement proper logging and monitoring
- Use load balancers for high-traffic scenarios
- Consider geographic distribution for global applications
- Test all protocol variants in your target environments
- Monitor which protocols are most successful in your use case

This server provides a complete WebRTC infrastructure that can handle
real-world networking challenges and provide reliable peer-to-peer communication
across all common network environments and security requirements.

For more information, see:
- RFC 5389: STUN Protocol
- RFC 5766: TURN Protocol
- RFC 5245: ICE Protocol
- WebRTC.org documentation
- Pion WebRTC library documentation
*/

package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"go-server/webrtc"

	"github.com/pion/turn/v4" // Pion TURN library - popular Go WebRTC implementation
)

// ============================================================================
// CONFIGURATION CONSTANTS
// ============================================================================

// HTTP/HTTPS server port for WebRTC signaling
// Port 443 is standard for HTTPS, but you can change this
// Note: Modern browsers require HTTPS for WebRTC, so HTTP is mainly for development
const (
	httpPort = 443
)

// Standard ports for different WebRTC protocols
// These are IANA-assigned standard ports that most WebRTC clients expect
// Using standard ports ensures maximum compatibility across different networks
const (
	turnUDPPort = 3478 // Standard TURN UDP - most common for WebRTC
	turnTCPPort = 3478 // Standard TURN TCP - same port as UDP
	turnTLSPort = 5349 // Standard TURNS (TURN over TLS) - secure TURN
	stunUDPPort = 3478 // Standard STUN UDP - same port as TURN UDP
	stunTCPPort = 3478 // Standard STUN TCP - same port as TURN TCP
	stunTLSPort = 5349 // Standard STUNS (STUN over TLS) - secure STUN
)

// ============================================================================
// GLOBAL VARIABLES
// ============================================================================

// Server instances for different protocols
// Each server handles a specific transport protocol (UDP/TCP/TLS)
// Having multiple servers allows us to handle different network environments
// NOTE: TURN servers inherently support STUN functionality - they are STUN/TURN servers
var (
	stunturnServer    *turn.Server      // UDP STUN/TURN server - handles both STUN discovery and TURN relay
	stunturnTCPServer *turn.Server      // TCP STUN/TURN server - fallback for UDP-blocked networks
	stunturnTLSServer *turn.Server      // TLS STUN/TURN server - secure encrypted discovery and relay
	usersMap          map[string][]byte // Authentication credentials (username -> auth key)
	stunturnPort      int               // STUN/TURN server port - configurable via command line

	// Loggers for different services
	// Separate loggers help with debugging and monitoring
	stunTurnLogger  *log.Logger // Logger for STUN/TURN services
	signalingLogger *log.Logger // Logger for WebRTC signaling

	// Monitoring processes for log windows
	// These help with real-time monitoring during development
	stunturnMonitor  *os.Process // Process for STUN/TURN log monitoring window
	signalingMonitor *os.Process // Process for signaling log monitoring window
)

// ============================================================================
// MAIN FUNCTION - SERVER ENTRY POINT
// ============================================================================

func main() {
	// ========================================================================
	// COMMAND LINE ARGUMENT PARSING
	// ========================================================================
	// These flags allow users to configure the server without modifying code
	// This is a common pattern in Go applications for flexibility
	// Users can customize the server behavior without touching the source code

	publicIP := flag.String("public-ip", "", "IP Address that TURN can be contacted by.")
	// ^ This is CRITICAL - TURN server must know its public IP for relay allocation
	//   Clients will connect to this IP address for relay services
	//   Example: "203.0.113.1" or "api.yourdomain.com"

	turnUsers := flag.String("turn-users", "", "List of username and password (e.g. \"user=pass,user=pass\")")
	// ^ TURN authentication credentials - prevents unauthorized relay usage
	//   Format: "username1=password1,username2=password2"
	//   Example: "alice=secret123,bob=secret456"
	//   In production, use strong, unique credentials

	realm := flag.String("realm", "pion.ly", "Realm (defaults to \"pion.ly\")")
	// ^ TURN realm - identifies the authentication domain
	//   Think of it as the "domain" for your TURN server
	//   Example: "yourcompany.com" or "webrtc.example.com"

	threadNum := flag.Int("thread-num", 1, "Number of server threads (defaults to 1)")
	// ^ Number of concurrent listeners - increases throughput for high-traffic scenarios
	//   Each thread handles connections independently
	//   Recommended: 1-4 threads depending on your server's CPU cores

	turnPortFlag := flag.Int("turn-port", 3478, "TURN server port (defaults to 3478)")
	// ^ Custom TURN port - useful if 3478 is blocked or in use
	//   Standard port 3478 is recommended for maximum compatibility

	enableTCP := flag.Bool("enable-tcp", true, "Enable TURN/STUN over TCP (defaults to true)")
	// ^ Enable TCP fallback - some networks block UDP, so TCP is essential
	//   Corporate networks often block UDP, making TCP necessary

	enableTLS := flag.Bool("enable-tls", true, "Enable TURN/STUN over TLS (defaults to true)")
	// ^ Enable TLS encryption - required for secure enterprise environments
	//   Also needed for WebRTC in browsers (HTTPS requirement)

	// New logging flags for better monitoring and debugging
	stunturnLogFile := flag.String("stun-turn-log", "stun-turn.log", "Log file for STUN/TURN services (defaults to stdout)")
	signalingLogFile := flag.String("signaling-log", "signaling.log", "Log file for WebRTC signaling (defaults to stdout)")
	separateLogs := flag.Bool("separate-logs", true, "Separate STUN/TURN and signaling logs (defaults to false)")

	flag.Parse() // Parse all command line arguments

	// ========================================================================
	// LOGGING SETUP
	// ========================================================================
	// Set up separate loggers for different services
	// This helps with debugging and monitoring by separating concerns
	setupLogging(*separateLogs, *stunturnLogFile, *signalingLogFile)

	// Set global turn port for use throughout the application
	stunturnPort = *turnPortFlag

	// ========================================================================
	// DEFAULT CONFIGURATION
	// ========================================================================
	// If no TURN users are provided, use a default credential
	// In production, you should always provide your own credentials
	// These default credentials are for educational purposes only
	if len(*turnUsers) == 0 {
		*turnUsers = "1ac96ad0a8374103e5c58441=drTJQZjbVFKpcXfn"
		stunTurnLogger.Println("Using default TURN credentials - NOT recommended for production!")
		stunTurnLogger.Println("For production, use: -turn-users \"youruser=yourpassword\"")
	}

	// ========================================================================
	// VALIDATION
	// ========================================================================
	// Public IP is required because TURN server needs to know its external address
	// This is used when allocating relay addresses to clients
	// Without this, clients won't be able to connect to the relay
	if len(*publicIP) == 0 {
		stunTurnLogger.Println("No public IP provided. Attempting to auto-detect...")

		// Try multiple methods to detect public IP
		detectedIP := ""

		// Method 1: Try HTTP-based detection (most reliable)
		if ip, err := detectPublicIPViaHTTP(); err == nil {
			detectedIP = ip
			stunTurnLogger.Printf("Detected public IP via HTTP service: %s", detectedIP)
		} else {
			stunTurnLogger.Printf("HTTP-based IP detection failed: %v", err)

			// Method 2: Try DNS-based detection (fallback)
			// Try ipify.org
			if ips, err := net.LookupIP("api.ipify.org"); err == nil && len(ips) > 0 {
				detectedIP = ips[0].String()
				stunTurnLogger.Printf("Detected public IP via ipify.org: %s", detectedIP)
			} else {
				stunTurnLogger.Printf("Failed to detect IP via ipify.org: %v", err)

				// Try icanhazip.com
				if ips, err := net.LookupIP("icanhazip.com"); err == nil && len(ips) > 0 {
					detectedIP = ips[0].String()
					stunTurnLogger.Printf("Detected public IP via icanhazip.com: %s", detectedIP)
				} else {
					stunTurnLogger.Printf("Failed to detect IP via icanhazip.com: %v", err)

					// Try checkip.amazonaws.com
					if ips, err := net.LookupIP("checkip.amazonaws.com"); err == nil && len(ips) > 0 {
						detectedIP = ips[0].String()
						stunTurnLogger.Printf("Detected public IP via checkip.amazonaws.com: %s", detectedIP)
					} else {
						stunTurnLogger.Printf("Failed to detect IP via checkip.amazonaws.com: %v", err)
					}
				}
			}
		}

		// If we successfully detected an IP, use it
		if detectedIP != "" {
			*publicIP = detectedIP
			stunTurnLogger.Printf("Using auto-detected public IP: %s", *publicIP)
			stunTurnLogger.Println("Note: Auto-detected IP may not be accurate in all network configurations.")
			stunTurnLogger.Println("For production use, explicitly specify your public IP with -public-ip flag.")
		} else {
			// Last resort: Try to detect local IP for development
			stunTurnLogger.Println("External IP detection failed. Trying local IP detection for development...")
			if localIP, err := detectLocalIP(); err == nil {
				*publicIP = localIP
				stunTurnLogger.Printf("Using local IP for development: %s", *publicIP)
				stunTurnLogger.Println("WARNING: This is a local IP address. TURN relay may not work properly.")
				stunTurnLogger.Println("For production use, explicitly specify your public IP with -public-ip flag.")
			} else {
				// If all methods failed, provide helpful error message
				stunTurnLogger.Fatalf("Failed to auto-detect public IP and no public IP provided.")
				stunTurnLogger.Fatalf("Please provide your public IP address using the -public-ip flag.")
				stunTurnLogger.Fatalf("Example: -public-ip 203.0.113.1")
				stunTurnLogger.Fatalf("")
				stunTurnLogger.Fatalf("Common reasons for auto-detection failure:")
				stunTurnLogger.Fatalf("- No internet connection")
				stunTurnLogger.Fatalf("- Firewall blocking outbound HTTP/DNS queries")
				stunTurnLogger.Fatalf("- Running behind a corporate proxy")
				stunTurnLogger.Fatalf("- DNS resolution issues")
				stunTurnLogger.Fatalf("- HTTPS certificate validation issues")
			}
		}
	} else {
		stunTurnLogger.Printf("Using provided public IP: %s", *publicIP)
	}

	// ========================================================================
	// SERVER INITIALIZATION
	// ========================================================================
	// Initialize all STUNTURN servers with the provided configuration
	// This sets up UDP, TCP, and TLS variants based on the flags
	// Each protocol serves different network environments
	if err := initializeSTUNTurnServer(*publicIP, *turnUsers, *realm, *threadNum, *enableTCP, *enableTLS); err != nil {
		stunTurnLogger.Fatalf("Failed to initialize STUN/TURN server: %v", err)
	}

	// ========================================================================
	// WEBSOCKET SIGNALING SETUP
	// ========================================================================
	// WebSocket endpoint for WebRTC signaling
	// This is where clients exchange connection information (SDP, ICE candidates)
	// Signaling is the "coordination" part of WebRTC - it helps peers find each other
	http.HandleFunc("/signal", func(w http.ResponseWriter, r *http.Request) {
		webrtc.HandleWebSocket(w, r, signalingLogger)
	})
	// The WebSocket handler manages:
	// - User registration and session management
	// - SDP offer/answer exchange
	// - ICE candidate sharing
	// - Call state management (join, call, hangup, etc.)

	// ========================================================================
	// CONNECTION MONITORING SETUP
	// ========================================================================
	// Start monitoring for connection statistics and debugging
	startConnectionMonitoring()

	// ========================================================================
	// GRACEFUL SHUTDOWN SETUP
	// ========================================================================
	// Create a channel to listen for shutdown signals (Ctrl+C, kill command)
	// This allows the server to shut down cleanly without dropping connections
	// Graceful shutdown is important for production servers
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// ========================================================================
	// HTTP/HTTPS SERVER STARTUP
	// ========================================================================
	// Start the HTTP/HTTPS server in a separate goroutine
	// This allows the main thread to handle shutdown signals
	// Goroutines are Go's lightweight threads for concurrent execution
	go startWebRTC_SignallingServer()

	// ========================================================================
	// SERVER STATUS LOGGING
	// ========================================================================
	// Log all the services that are now running
	// This helps with debugging and monitoring
	// Users can see exactly what's available and on which ports
	stunTurnLogger.Printf("=== STUN/TURN SERVER STATUS ===")
	stunTurnLogger.Printf("Unified WebRTC server started:")
	stunTurnLogger.Printf("- STUN/TURN server UDP: :%d (STUN discovery + TURN relay)", stunturnPort)
	if *enableTCP {
		stunTurnLogger.Printf("- STUN/TURN server TCP: :%d (STUN discovery + TURN relay)", stunturnPort)
	}
	if *enableTLS {
		stunTurnLogger.Printf("- STUN/TURN server TLS: :%d (STUN discovery + TURN relay)", turnTLSPort)
	}
	stunTurnLogger.Printf("- Public IP: %s", *publicIP)
	stunTurnLogger.Printf("- Realm: %s", *realm)
	stunTurnLogger.Printf("=== STUN/TURN SERVER READY ===")

	signalingLogger.Printf("=== WEBRTC SIGNALING SERVER STATUS ===")
	signalingLogger.Printf("- Signaling server: :%d (HTTP/HTTPS)", httpPort)
	signalingLogger.Printf("- WebSocket endpoint: /signal")
	signalingLogger.Printf("=== SIGNALING SERVER READY ===\n\n\n")

	// ========================================================================
	// MAIN EVENT LOOP
	// ========================================================================
	// Block until user sends SIGINT (Ctrl+C) or SIGTERM (kill command)
	// This keeps the server running until explicitly stopped
	// The server will continue running and handling requests until shutdown
	<-sigs

	// ========================================================================
	// GRACEFUL SHUTDOWN
	// ========================================================================
	// When shutdown signal is received, close all servers cleanly
	// This ensures no data is lost and connections are properly closed
	stunTurnLogger.Println("Shutting down STUN/TURN servers...")
	signalingLogger.Println("Shutting down signaling server...")

	// Close all TURN/STUN servers to free resources and close connections
	// This prevents resource leaks and ensures clean shutdown
	servers := []*turn.Server{stunturnServer, stunturnTCPServer, stunturnTLSServer}
	for _, server := range servers {
		if server != nil {
			if err := server.Close(); err != nil {
				stunTurnLogger.Printf("Failed to close server: %v", err)
			}
		}
	}

	// Close monitoring windows
	// Clean up any monitoring processes we started
	if runtime.GOOS == "windows" {
		killBatchWindow("stun-turn-monitor.ps1")
		killBatchWindow("signaling-monitor.ps1")
		// On Windows, create shutdown signal file to tell PowerShell files to close
		os.WriteFile("shutdown-signal.txt", []byte("shutdown"), 0644)
		stunTurnLogger.Printf("Monitoring windows will close automatically")

		// Clean up temporary PowerShell files and shutdown signal
		os.Remove("stun-turn-monitor.ps1")
		os.Remove("signaling-monitor.ps1")
		os.Remove("shutdown-signal.txt")
	} else {
		// On Unix systems, use SIGTERM for graceful shutdown
		if stunturnMonitor != nil {
			if err := stunturnMonitor.Signal(syscall.SIGTERM); err != nil {
				stunTurnLogger.Printf("Failed to close STUN/TURN monitoring window: %v", err)
			} else {
				stunTurnLogger.Printf("STUN/TURN monitoring window closed")
			}
		}

		if signalingMonitor != nil {
			if err := signalingMonitor.Signal(syscall.SIGTERM); err != nil {
				stunTurnLogger.Printf("Failed to close signaling monitoring window: %v", err)
			} else {
				stunTurnLogger.Printf("Signaling monitoring window closed")
			}
		}
	}

	stunTurnLogger.Println("STUN/TURN servers shut down successfully")
	signalingLogger.Println("Signaling server shut down successfully")
}

// ============================================================================
// UTILITY FUNCTIONS
// ============================================================================

// detectPublicIPViaHTTP attempts to detect public IP using HTTP services
// This is more reliable than DNS-based methods in some network environments
func detectPublicIPViaHTTP() (string, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
		},
	}

	// List of IP detection services to try
	services := []string{
		"https://api.ipify.org",
		"https://icanhazip.com",
		"https://checkip.amazonaws.com",
		"https://ifconfig.me/ip",
		"https://ipecho.net/plain",
	}

	for _, service := range services {
		resp, err := client.Get(service)
		if err != nil {
			continue // Try next service
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				continue
			}

			// Clean the response (remove whitespace, newlines)
			ip := strings.TrimSpace(string(body))

			// Validate that it looks like an IP address
			if net.ParseIP(ip) != nil {
				return ip, nil
			}
		}
	}

	return "", fmt.Errorf("all HTTP IP detection services failed")
}

// detectLocalIP attempts to find a suitable local IP address for development
func detectLocalIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no suitable local IP address found")
}

// ============================================================================
// WINDOWS PROCESS MANAGEMENT
// ============================================================================

// killBatchWindow terminates PowerShell monitoring windows on Windows systems
// This function is called during graceful shutdown to clean up monitoring processes
//
// WHY IS THIS NEEDED?
// ===================
// When we start separate monitoring windows for STUN/TURN and signaling logs,
// these windows run as separate PowerShell processes. During shutdown, we need
// to terminate these processes to prevent orphaned windows and ensure clean exit.
//
// HOW IT WORKS:
// =============
// 1. Uses Windows Management Instrumentation (WMI) to find PowerShell processes
// 2. Searches for processes containing the specific batch filename
// 3. Extracts the Process ID (PID) from the WMI output
// 4. Uses taskkill to forcefully terminate each process
//
// WINDOWS-SPECIFIC CONSIDERATIONS:
// ================================
// - WMI queries are Windows-specific and won't work on Unix systems
// - CSV format parsing is used because WMI output is structured
// - taskkill /f forces termination even if process is unresponsive
// - This ensures no monitoring windows are left running after server shutdown
func killBatchWindow(batchFileName string) {
	fmt.Printf("Attempting to kill batch window: %s\n", batchFileName)

	// Try to find specific CMD processes with batch filename
	// WMI (Windows Management Instrumentation) allows us to query system processes
	// We search for processes whose command line contains our batch filename
	escapedName := strings.ReplaceAll(batchFileName, "\\", "\\\\")
	query := fmt.Sprintf("commandline like '%%%s%%'", escapedName)
	findCmd := exec.Command("wmic", "process", "where", query, "get", "processid", "/format:csv")

	output, err := findCmd.Output()
	if err != nil {
		fmt.Printf("Failed to find CMD batch processes: %v\n", err)
		return
	}

	// Parse the CSV output to get PIDs
	// WMI returns data in CSV format: Node,ProcessId
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || line == "Node,ProcessId" {
			continue
		}

		// CSV format: Node,ProcessId
		parts := strings.Split(line, ",")
		if len(parts) >= 2 {
			pid := strings.TrimSpace(parts[1])
			if pid != "" && pid != "ProcessId" {
				fmt.Printf("Found CMD batch process PID: %s\n", pid)

				// Kill this specific process
				// taskkill /f forces termination even if process is unresponsive
				killCmd := exec.Command("taskkill", "/pid", pid, "/f")
				if err := killCmd.Run(); err != nil {
					fmt.Printf("Failed to kill process %s: %v\n", pid, err)
				} else {
					fmt.Printf("Successfully killed process %s\n", pid)
				}
			}
		}
	}
}

// ============================================================================
// LOGGING AND MONITORING SETUP
// ============================================================================

// setupLogging configures separate loggers for STUN/TURN and signaling services
// This function creates a sophisticated logging system with real-time monitoring capabilities
//
// WHY SEPARATE LOGGING?
// =====================
// WebRTC servers handle multiple types of traffic:
// - STUN/TURN: Network traversal and relay services (UDP/TCP/TLS)
// - Signaling: WebSocket connections for peer coordination
//
// Separate logging allows developers to:
// - Debug network issues independently from signaling issues
// - Monitor different aspects of the system in real-time
// - Filter logs by service type for better analysis
// - Identify which component is causing problems
//
// MONITORING WINDOWS:
// ===================
// This function can open separate terminal windows to monitor logs in real-time.
// This is especially useful during development and debugging.
//
// CROSS-PLATFORM SUPPORT:
// =======================
// - Windows: Uses PowerShell with custom monitoring scripts
// - Unix/Linux: Uses xterm with tail -f command
// - Fallback: Single logger to stdout if monitoring fails
//
// LOG FILE MANAGEMENT:
// ====================
// - Clears existing log files to start fresh
// - Creates new log files with proper permissions
// - Handles both file and stdout logging
// - Provides structured log prefixes for easy filtering
func setupLogging(separateLogs bool, stunturnLogFile, signalingLogFile string) {
	if separateLogs {
		// Clear existing log files to start fresh
		// This prevents log files from growing indefinitely and ensures clean logs
		if stunturnLogFile != "" {
			os.Remove(stunturnLogFile) // Remove existing file
		}
		if signalingLogFile != "" {
			os.Remove(signalingLogFile) // Remove existing file
		}

		// Set up STUN/TURN logger
		// This logger handles all STUN and TURN server activities
		// STUN/TURN logs include: authentication, relay allocation, connection events
		if stunturnLogFile != "" {
			file, err := os.OpenFile(stunturnLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				log.Fatalf("Failed to open STUN/TURN log file: %v", err)
			}
			stunTurnLogger = log.New(file, "[STUN/TURN] ", log.LstdFlags|log.Lshortfile)
		} else {
			stunTurnLogger = log.New(os.Stdout, "[STUN/TURN] ", log.LstdFlags|log.Lshortfile)
		}

		// Set up signaling logger
		// This logger handles all WebSocket signaling activities
		// Signaling logs include: user connections, SDP exchange, call management
		if signalingLogFile != "" {
			file, err := os.OpenFile(signalingLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				log.Fatalf("Failed to open signaling log file: %v", err)
			}
			signalingLogger = log.New(file, "[SIGNALING] ", log.LstdFlags|log.Lshortfile)
		} else {
			signalingLogger = log.New(os.Stdout, "[SIGNALING] ", log.LstdFlags|log.Lshortfile)
		}

		// Open a new terminal window to monitor STUN/TURN logs in real-time
		// This helps with debugging and monitoring server activity
		// Real-time monitoring is crucial for understanding connection patterns
		if runtime.GOOS == "windows" {
			// For Windows, create batch files for monitoring
			// PowerShell is used because it provides better process control than CMD
			// Create PowerShell file content for STUN/TURN monitoring
			// This script continuously monitors the log file and displays new entries
			stunturnPS := fmt.Sprintf(`$Host.UI.RawUI.WindowTitle = "STUN/TURN Log Monitor"
$logFile = "%s"
$lastLineCount = 0

while (-not (Test-Path "shutdown-signal.txt")) {
    if (Test-Path $logFile) {
        $currentLineCount = (Get-Content $logFile).Count
        if ($currentLineCount -gt $lastLineCount) {
            $newLines = Get-Content $logFile | Select-Object -Skip $lastLineCount
            $newLines | ForEach-Object { Write-Host $_ }
            $lastLineCount = $currentLineCount
        }
    }
    Start-Sleep -Seconds 1
}
exit`, stunturnLogFile)

			// Create PowerShell file content for signaling monitoring
			// Similar script but for signaling logs
			signalingPS := fmt.Sprintf(`$Host.UI.RawUI.WindowTitle = "Signaling Log Monitor"
$logFile = "%s"
$lastLineCount = 0

while (-not (Test-Path "shutdown-signal.txt")) {
    if (Test-Path $logFile) {
        $currentLineCount = (Get-Content $logFile).Count
        if ($currentLineCount -gt $lastLineCount) {
            $newLines = Get-Content $logFile | Select-Object -Skip $lastLineCount
            $newLines | ForEach-Object { Write-Host $_ }
            $lastLineCount = $currentLineCount
        }
    }
    Start-Sleep -Seconds 1
}
exit`, signalingLogFile)

			// Write PowerShell files
			// These temporary files contain the monitoring scripts
			os.WriteFile("stun-turn-monitor.ps1", []byte(stunturnPS), 0644)
			os.WriteFile("signaling-monitor.ps1", []byte(signalingPS), 0644)

			// Start the PowerShell files in new windows
			// Each monitoring window runs independently
			cmd1 := exec.Command("cmd", "/c", "start", "powershell", "-ExecutionPolicy", "Bypass", "-File", "stun-turn-monitor.ps1")
			cmd2 := exec.Command("cmd", "/c", "start", "powershell", "-ExecutionPolicy", "Bypass", "-File", "signaling-monitor.ps1")

			// Start both monitoring processes
			// Store process references for graceful shutdown
			if err := cmd1.Start(); err != nil {
				stunTurnLogger.Printf("Failed to open STUN/TURN log monitor window: %v", err)
			} else {
				stunturnMonitor = cmd1.Process
				stunTurnLogger.Printf("STUN/TURN log monitor window opened successfully")
			}

			if err := cmd2.Start(); err != nil {
				stunTurnLogger.Printf("Failed to open signaling log monitor window: %v", err)
			} else {
				signalingMonitor = cmd2.Process
				stunTurnLogger.Printf("Signaling log monitor window opened successfully")
			}
		} else {
			// For Linux/Unix, use 'xterm' to open new terminal
			// The -e flag executes the tail command in the new window
			cmd1 := exec.Command("xterm", "-e", "tail", "-f", stunturnLogFile)
			cmd2 := exec.Command("xterm", "-e", "tail", "-f", signalingLogFile)

			// Start both monitoring processes
			// Unix systems use different process management than Windows
			if err := cmd1.Start(); err != nil {
				stunTurnLogger.Printf("Failed to open STUN/TURN log monitor window: %v", err)
			} else {
				stunturnMonitor = cmd1.Process
				stunTurnLogger.Printf("STUN/TURN log monitor window opened successfully")
			}

			if err := cmd2.Start(); err != nil {
				stunTurnLogger.Printf("Failed to open signaling log monitor window: %v", err)
			} else {
				signalingMonitor = cmd2.Process
				stunTurnLogger.Printf("Signaling log monitor window opened successfully")
			}
		}
	} else {
		// Use single logger for all services
		// This is the fallback option when separate logging is disabled
		// All logs go to stdout with a generic [WEBRTC] prefix
		logger := log.New(os.Stdout, "[WEBRTC] ", log.LstdFlags|log.Lshortfile)
		stunTurnLogger = logger
		signalingLogger = logger
	}
}

// ============================================================================
// TURN SERVER INITIALIZATION
// ============================================================================

// initializeTURNServer sets up the STUN and TURN servers with all protocol variants
// This is the main initialization function that creates all server instances
//
// WHAT IS A TURN SERVER?
// ======================
// TURN (Traversal Using Relays around NAT) servers help WebRTC clients establish
// peer-to-peer connections when direct connection is impossible due to NAT/firewall.
//
// TURN servers act as "relay points" that forward data between clients when they
// cannot connect directly. This is essential for WebRTC to work in all network
// environments.
//
// AUTHENTICATION SYSTEM:
// =====================
// TURN servers require authentication to prevent abuse and unauthorized relay usage.
// This implementation uses the standard TURN authentication mechanism:
// - Username/password pairs are provided via command line
// - Cryptographic keys are generated using TURN protocol specification
// - Authentication is validated on every relay request
//
// PROTOCOL VARIANTS:
// ==================
// This server supports multiple transport protocols to handle different network
// environments:
// - UDP: Standard protocol, fastest, works with most NATs
// - TCP: Fallback for UDP-blocked networks (corporate firewalls)
// - TLS: Encrypted connections for secure environments
//
// RELAY ADDRESS GENERATION:
// =========================
// TURN servers allocate relay addresses for clients. These addresses must be
// reachable from the internet, so the server needs to know its public IP.
//
// THREADING MODEL:
// ================
// Multiple threads can be used to handle concurrent connections. Each thread
// gets its own listener, improving performance under high load.
func initializeTURNServer(publicIP, users, realm string, threadNum int, enableTCP, enableTLS bool) error {
	// ========================================================================
	// USER AUTHENTICATION SETUP
	// ========================================================================
	// Parse TURN user credentials from the command line argument
	// Format: "user1=pass1,user2=pass2"
	// This creates a map of username -> cryptographic auth key
	usersMap = make(map[string][]byte)

	// Use regex to parse username=password pairs
	// This regex finds all patterns like "username=password"
	// The regex (\w+)=(\w+) captures:
	// - Group 1: username (word characters)
	// - Group 2: password (word characters)
	for _, kv := range regexp.MustCompile(`(\w+)=(\w+)`).FindAllStringSubmatch(users, -1) {
		// Generate authentication key using TURN protocol specification
		// This creates a cryptographic key from username, realm, and password
		// The key is used to validate TURN requests from clients
		usersMap[kv[1]] = turn.GenerateAuthKey(kv[1], realm, kv[2])
		stunTurnLogger.Printf("Added TURN user: %s", kv[1])
	}

	// ========================================================================
	// RELAY ADDRESS GENERATOR
	// ========================================================================
	// This tells the TURN server what IP address to use for relay allocation
	// When a client requests a relay, the server will allocate an address on this IP
	// The publicIP must be reachable from the internet for relay to work
	relayAddressGenerator := &turn.RelayAddressGeneratorStatic{
		RelayAddress: net.ParseIP(publicIP), // Public IP for relay allocation
		Address:      "0.0.0.0",             // Listen on all interfaces
	}

	// ========================================================================
	// AUTHENTICATION HANDLER
	// ========================================================================
	// This function is called whenever a client tries to authenticate
	// It validates the username and returns the corresponding auth key
	// If authentication fails, the client cannot use relay services
	authHandler := createEnhancedAuthHandler(usersMap)

	// ========================================================================
	// SERVER INITIALIZATION SEQUENCE
	// ========================================================================
	// Initialize servers in order of importance and dependency
	// Each protocol variant serves different network environments

	// 2. UDP TURN server - main relay service, handles most WebRTC traffic
	// UDP is the standard protocol for TURN and works with most NAT types
	// It's the fastest and most efficient option
	if err := initializeUDPTURNServer(relayAddressGenerator, authHandler, realm, threadNum); err != nil {
		return fmt.Errorf("failed to initialize UDP STUN/TURN server: %w", err)
	}

	// 4. TCP TURN server (if enabled) - fallback relay service
	// TCP is used when UDP is blocked by firewalls or NATs
	// Common in corporate networks that block UDP traffic
	if enableTCP {
		if err := initializeTCPTURNServer(relayAddressGenerator, authHandler, realm, threadNum); err != nil {
			return fmt.Errorf("failed to initialize TCP STUN/TURN server: %w", err)
		}
	}

	// 6. TLS TURN server (if enabled) - secure relay service
	// TLS provides encrypted relay connections
	// Required for secure enterprise environments and browser compatibility
	if enableTLS {
		if err := initializeTLSTURNServer(relayAddressGenerator, authHandler, realm, threadNum); err != nil {
			return fmt.Errorf("failed to initialize TLS STUN/TURN server: %w", err)
		}
	}

	return nil
}

// ============================================================================
// UDP TURN SERVER IMPLEMENTATION
// ============================================================================

// initializeUDPTURNServer sets up UDP TURN server
// TURN (Traversal Using Relays around NAT) provides relay services when direct connection fails
//
// WHY UDP FOR TURN?
// ================
// UDP is the standard transport protocol for TURN servers because:
// - It's faster than TCP (no connection establishment overhead)
// - It works better with most NAT types
// - It's more efficient for real-time media (WebRTC's primary use case)
// - It handles packet loss gracefully (important for video/audio)
//
// RELAY MECHANISM:
// ================
// When two WebRTC clients cannot connect directly (due to NAT/firewall):
// 1. Client A connects to TURN server and requests relay
// 2. TURN server allocates a relay address for Client A
// 3. Client B connects to TURN server and requests relay
// 4. TURN server allocates a relay address for Client B
// 5. TURN server forwards data between the two relay addresses
//
// THREADING MODEL:
// ================
// Multiple threads are used to handle concurrent connections:
// - Each thread gets its own UDP listener
// - This improves performance under high load
// - Prevents one slow connection from blocking others
// - Allows better CPU utilization on multi-core systems
//
// SOCKET OPTIONS:
// ===============
// SO_REUSEADDR: Allows multiple listeners to bind to the same port
// SO_BROADCAST: Enables broadcast capabilities for UDP
// These options are essential for proper UDP server operation
func initializeUDPTURNServer(relayGen *turn.RelayAddressGeneratorStatic, authHandler func(string, string, net.Addr) ([]byte, bool), realm string, threadNum int) error {
	// Create UDP address for the server
	// "0.0.0.0" means listen on all network interfaces
	// Port 3478 is the standard TURN UDP port (IANA assigned)
	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:"+strconv.Itoa(turnUDPPort))
	if err != nil {
		return fmt.Errorf("failed to parse server address: %w", err)
	}

	// Create listener configuration with proper socket options for multithreading
	// SO_REUSEADDR allows multiple listeners to bind to the same port
	// SO_BROADCAST enables broadcast capabilities for UDP
	// These socket options are crucial for proper UDP server operation
	listenerConfig := &net.ListenConfig{
		Control: func(network, address string, conn syscall.RawConn) error {
			var operr error
			if err := conn.Control(func(fd uintptr) {
				// Set SO_REUSEADDR to allow multiple listeners on same port
				// This is essential when using multiple threads
				operr = syscall.SetsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
				if operr != nil {
					return
				}
				// Set SO_BROADCAST for UDP broadcast capabilities
				// This allows the server to handle broadcast packets
				operr = syscall.SetsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1)
			}); err != nil {
				return err
			}
			return operr
		},
	}

	// Create multiple UDP listeners for better performance
	// Each thread gets its own listener to handle concurrent connections
	// This prevents connection bottlenecks and improves throughput
	packetConnConfigs := make([]turn.PacketConnConfig, threadNum)
	stunTurnLogger.Printf("")
	for i := 0; i < threadNum; i++ {
		// Create UDP listener with proper socket options
		// Each listener runs on the same port but in a separate thread
		conn, err := listenerConfig.ListenPacket(context.Background(), addr.Network(), addr.String())
		if err != nil {
			return fmt.Errorf("failed to create UDP TURN listener %d: %w", i, err)
		}

		// Configure the packet connection with relay capabilities
		// Each listener is configured with the same relay address generator
		// This ensures consistent relay allocation across all threads
		packetConnConfigs[i] = turn.PacketConnConfig{
			PacketConn:            conn,     // UDP connection
			RelayAddressGenerator: relayGen, // How to allocate relay addresses
		}
		stunTurnLogger.Printf("UDP TURN server %d listening on %s", i, conn.LocalAddr().String())
	}

	// Create TURN server with authentication and relay capabilities
	// The server combines all UDP listeners into a single TURN server instance
	// This provides unified authentication and relay management
	stunturnServer, err = turn.NewServer(turn.ServerConfig{
		Realm:             realm,             // Authentication realm
		AuthHandler:       authHandler,       // Authentication function
		PacketConnConfigs: packetConnConfigs, // UDP listeners
	})
	if err != nil {
		return fmt.Errorf("failed to create UDP TURN server: %w", err)
	}
	return nil
}

// ============================================================================
// TCP TURN SERVER IMPLEMENTATION
// ============================================================================

// initializeTCPTURNServer sets up TCP TURN server
// TCP TURN provides relay services over TCP when UDP is blocked
//
// WHY TCP FOR TURN?
// ================
// TCP TURN is used as a fallback when UDP is blocked:
// - Corporate firewalls often block UDP traffic
// - Some NATs don't handle UDP well
// - TCP is more reliable in restrictive network environments
// - Required for enterprise networks with strict security policies
//
// TCP vs UDP PERFORMANCE:
// =======================
// TCP TURN is slower than UDP TURN because:
// - TCP has connection establishment overhead (3-way handshake)
// - TCP has congestion control and flow control
// - TCP retransmits lost packets (good for reliability, bad for latency)
// - TCP has higher latency due to buffering
//
// WHEN TO USE TCP TURN:
// ====================
// Use TCP TURN when:
// - UDP is blocked by firewall/NAT
// - Network has strict security policies
// - Reliability is more important than latency
// - Working in corporate/enterprise environments
//
// RELAY MECHANISM (TCP):
// ======================
// TCP relay works similarly to UDP relay:
// 1. Client A establishes TCP connection to TURN server
// 2. TURN server allocates relay address for Client A
// 3. Client B establishes TCP connection to TURN server
// 4. TURN server allocates relay address for Client B
// 5. TURN server forwards data between the two TCP connections
//
// THREADING MODEL:
// ================
// Similar to UDP, multiple threads handle concurrent connections
// Each thread gets its own TCP listener for better performance
func initializeTCPTURNServer(relayGen *turn.RelayAddressGeneratorStatic, authHandler func(string, string, net.Addr) ([]byte, bool), realm string, threadNum int) error {
	// Create TCP address for the server
	// Same port as UDP (3478) but different protocol
	// "0.0.0.0" means listen on all network interfaces
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:"+strconv.Itoa(turnTCPPort))
	if err != nil {
		return fmt.Errorf("failed to parse server address: %w", err)
	}

	// Create listener configuration with proper socket options for multithreading
	// SO_REUSEADDR allows multiple listeners to bind to the same port
	// This is essential for multi-threaded TCP servers
	listenerConfig := &net.ListenConfig{
		Control: func(network, address string, conn syscall.RawConn) error {
			var operr error
			if err := conn.Control(func(fd uintptr) {
				// Set SO_REUSEADDR to allow multiple listeners on same port
				// This enables multiple threads to share the same port
				operr = syscall.SetsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
			}); err != nil {
				return err
			}
			return operr
		},
	}

	// Create multiple TCP listeners for better performance
	// Each thread gets its own listener to handle concurrent connections
	// This prevents connection bottlenecks and improves throughput
	listenerConfigs := make([]turn.ListenerConfig, threadNum)
	stunTurnLogger.Printf("")
	for i := 0; i < threadNum; i++ {
		// Create TCP listener with proper socket options
		// Each listener runs on the same port but in a separate thread
		listener, err := listenerConfig.Listen(context.Background(), addr.Network(), addr.String())
		if err != nil {
			return fmt.Errorf("failed to create TCP TURN listener %d: %w", i, err)
		}

		// Configure the TCP listener with relay capabilities
		// Each listener is configured with the same relay address generator
		// This ensures consistent relay allocation across all threads
		listenerConfigs[i] = turn.ListenerConfig{
			Listener:              listener, // TCP connection
			RelayAddressGenerator: relayGen, // How to allocate relay addresses
		}
		stunTurnLogger.Printf("TCP TURN server %d listening on %s", i, listener.Addr().String())
	}

	// Create TURN server with TCP listeners
	// The server combines all TCP listeners into a single TURN server instance
	// This provides unified authentication and relay management
	stunturnTCPServer, err = turn.NewServer(turn.ServerConfig{
		Realm:           realm,           // Authentication realm
		AuthHandler:     authHandler,     // Authentication function
		ListenerConfigs: listenerConfigs, // TCP listeners
	})
	if err != nil {
		return fmt.Errorf("failed to create TCP TURN server: %w", err)
	}
	return nil
}

// ============================================================================
// TLS TURN SERVER IMPLEMENTATION
// ============================================================================

// initializeTLSTURNServer sets up TLS TURN server
// TLS TURN provides encrypted relay services for secure environments
//
// WHY TLS FOR TURN?
// ================
// TLS TURN provides encrypted relay connections:
// - All relay traffic is encrypted end-to-end
// - Required for secure enterprise environments
// - Needed for WebRTC in browsers (HTTPS requirement)
// - Protects against man-in-the-middle attacks
// - Complies with security policies and regulations
//
// SECURITY BENEFITS:
// ==================
// TLS encryption provides:
// - Confidentiality: Relay data cannot be intercepted
// - Integrity: Data cannot be modified in transit
// - Authentication: Server identity is verified
// - Forward secrecy: Past communications remain secure
//
// CERTIFICATE REQUIREMENTS:
// =========================
// TLS TURN requires SSL/TLS certificates:
// - Certificate must be valid and trusted
// - Domain name must match server hostname
// - Certificate must be in PEM format
// - Private key must be secure and accessible
//
// WHEN TO USE TLS TURN:
// ====================
// Use TLS TURN when:
// - Security is a primary concern
// - Working in enterprise environments
// - WebRTC is used in browsers (HTTPS requirement)
// - Compliance with security policies is required
// - Protecting sensitive communications
//
// PERFORMANCE CONSIDERATIONS:
// ==========================
// TLS TURN has additional overhead:
// - TLS handshake adds connection establishment time
// - Encryption/decryption uses CPU resources
// - Slightly higher latency than TCP TURN
// - Still faster than TCP for most use cases
//
// RELAY MECHANISM (TLS):
// ======================
// TLS relay works like TCP relay but with encryption:
// 1. Client A establishes TLS connection to TURN server
// 2. TURN server allocates relay address for Client A
// 3. Client B establishes TLS connection to TURN server
// 4. TURN server allocates relay address for Client B
// 5. TURN server forwards encrypted data between connections
func initializeTLSTURNServer(relayGen *turn.RelayAddressGeneratorStatic, authHandler func(string, string, net.Addr) ([]byte, bool), realm string, threadNum int) error {
	// Check if SSL certificates exist (same as TLS STUN)
	// Certificates must be in the certs/ directory
	// fullchain.pem contains the certificate chain
	// privkey.pem contains the private key
	certFile := "certs/fullchain.pem"
	keyFile := "certs/privkey.pem"
	var err error

	// If certificates don't exist, skip TLS server
	// This allows the server to run without TLS if certificates are not available
	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		stunTurnLogger.Printf("SSL certificates not found. Skipping TLS TURN server.")
		return nil
	}

	// Load TLS certificate and private key
	// The certificate must be valid and trusted by clients
	// The private key must be secure and accessible to the server
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return fmt.Errorf("failed to load TLS certificate: %w", err)
	}

	// Configure TLS settings
	// MinVersion ensures we use secure TLS versions
	// TLS 1.2 is the minimum recommended version for security
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert}, // Our SSL certificate
		MinVersion:   tls.VersionTLS12,        // Minimum TLS version (secure)
	}

	// Create TCP address for the server
	// Port 5349 is the standard TURNS (TURN over TLS) port
	// Different from standard TURN port (3478) to distinguish protocols
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:"+strconv.Itoa(turnTLSPort))
	if err != nil {
		return fmt.Errorf("failed to parse server address: %w", err)
	}

	// Create listener configuration with proper socket options for multithreading
	// SO_REUSEADDR allows multiple listeners to bind to the same port
	// This is essential for multi-threaded TLS servers
	listenerConfig := &net.ListenConfig{
		Control: func(network, address string, conn syscall.RawConn) error {
			var operr error
			if err := conn.Control(func(fd uintptr) {
				// Set SO_REUSEADDR to allow multiple listeners on same port
				// This enables multiple threads to share the same port
				operr = syscall.SetsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
			}); err != nil {
				return err
			}
			return operr
		},
	}

	// Create multiple TLS listeners for better performance
	// Each thread gets its own listener to handle concurrent connections
	// This prevents connection bottlenecks and improves throughput
	listenerConfigs := make([]turn.ListenerConfig, threadNum)
	stunTurnLogger.Printf("")
	for i := 0; i < threadNum; i++ {
		// Create TCP listener with proper socket options first
		// TLS is built on top of TCP, so we start with a TCP listener
		tcpListener, err := listenerConfig.Listen(context.Background(), addr.Network(), addr.String())
		if err != nil {
			return fmt.Errorf("failed to create TCP listener for TLS TURN %d: %w", i, err)
		}

		// Wrap TCP listener with TLS
		// This adds encryption to the TCP connection
		// All data transmitted through this listener will be encrypted
		tlsListener := tls.NewListener(tcpListener, tlsConfig)

		// Configure the TLS listener with relay capabilities
		// Each listener is configured with the same relay address generator
		// This ensures consistent relay allocation across all threads
		listenerConfigs[i] = turn.ListenerConfig{
			Listener:              tlsListener, // TLS connection
			RelayAddressGenerator: relayGen,    // How to allocate relay addresses
		}
		stunTurnLogger.Printf("TLS TURN server %d listening on %s", i, tlsListener.Addr().String())
	}

	// Create TURN server with TLS listeners
	// The server combines all TLS listeners into a single TURN server instance
	// This provides unified authentication and relay management
	stunturnTLSServer, err = turn.NewServer(turn.ServerConfig{
		Realm:           realm,           // Authentication realm
		AuthHandler:     authHandler,     // Authentication function
		ListenerConfigs: listenerConfigs, // TLS listeners
	})
	if err != nil {
		return fmt.Errorf("failed to create TLS TURN server: %w", err)
	}
	return nil
}

// ============================================================================
// HTTP/HTTPS SERVER FOR WEBSOCKET SIGNALING
// ============================================================================

// startWebRTC_SignallingServer starts the HTTP/HTTPS server for WebSocket signaling
// This server handles the WebSocket connections that clients use for signaling
//
// WHAT IS WEBSOCKET SIGNALING?
// ============================
// WebRTC signaling is the process of exchanging connection information between peers:
// - SDP (Session Description Protocol) offers and answers
// - ICE (Interactive Connectivity Establishment) candidates
// - Call control messages (join, call, hangup, etc.)
//
// WebSocket provides a persistent, bidirectional communication channel for this exchange.
//
// HTTPS REQUIREMENT:
// ==================
// Modern browsers require HTTPS for WebRTC to work:
// - getUserMedia() (camera/microphone access) requires secure context
// - WebRTC connections are blocked in non-secure contexts
// - This is a security measure to protect user privacy
//
// CERTIFICATE HANDLING:
// ====================
// The server automatically detects SSL certificates:
// - If certificates exist: Start HTTPS server
// - If no certificates: Start HTTP server (development only)
// - Certificates must be in certs/ directory
// - Supports Let's Encrypt and other certificate authorities
//
// WEBSOCKET ENDPOINT:
// ===================
// The /signal endpoint handles all WebRTC signaling:
// - Accepts WebSocket upgrade requests
// - Manages user sessions and connections
// - Routes signaling messages between peers
// - Handles connection lifecycle (connect, disconnect, etc.)
//
// ERROR HANDLING:
// ===============
// The server includes comprehensive error handling:
// - TLS configuration errors
// - Certificate loading failures
// - Port binding issues
// - Graceful fallback to HTTP when needed
func startWebRTC_SignallingServer() {
	// SSL certificate files for HTTPS
	// These files must be in PEM format
	// fullchain.pem contains the certificate chain (including intermediate certificates)
	// privkey.pem contains the private key (must be kept secure)
	certFile := "certs/fullchain.pem"
	keyFile := "certs/privkey.pem"

	// Check if certificates exist to decide between HTTP and HTTPS
	// This allows the server to run in both development and production environments
	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		// No SSL certificates found - start HTTP server
		// This is suitable for development and testing
		// Note: WebRTC may not work in browsers without HTTPS
		signalingLogger.Printf("SSL certificate not found. Starting HTTP server on :%d", httpPort)
		signalingLogger.Println("To enable HTTPS, place fullchain.pem and privkey.pem files in the certs/ directory")

		// Start HTTP server
		// Note: Modern browsers require HTTPS for WebRTC, so HTTP is mainly for development
		// HTTP can be used for testing with non-browser clients (mobile apps, etc.)
		signalingLogger.Printf("WebRTC signaling server starting on :%d (HTTP)", httpPort)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil); err != nil {
			signalingLogger.Fatal("Server error:", err)
		}
	} else {
		// SSL certificates found - start HTTPS server
		// This is the recommended configuration for production use
		signalingLogger.Printf("SSL certificates found. Starting HTTPS server on :%d", httpPort)
		signalingLogger.Printf("WebRTC signaling server starting on :%d (HTTPS)", httpPort)

		// Configure TLS settings for HTTPS
		// MinVersion ensures we use secure TLS versions
		// TLS 1.2 is the minimum recommended version for security
		tlsConfig := &tls.Config{
			MinVersion: tls.VersionTLS12, // Minimum TLS version (secure)
		}

		// Create HTTPS server with TLS configuration and custom error logging
		// The server includes proper error handling and logging
		// Custom error logger helps with debugging TLS issues
		server := &http.Server{
			Addr:      fmt.Sprintf(":%d", httpPort),
			TLSConfig: tlsConfig,
			ErrorLog:  signalingLogger,
		}

		// Start HTTPS server with SSL certificates
		// This provides secure WebSocket connections (WSS)
		// Required for WebRTC to work in modern browsers
		if err := server.ListenAndServeTLS(certFile, keyFile); err != nil {
			signalingLogger.Fatal("HTTPS Server error:", err)
		}
	}
}

// ============================================================================
// MONITORING AND STATISTICS
// ============================================================================

// startMonitoring starts a goroutine to monitor server statistics
// This function provides real-time insights into server performance and health
//
// WHY MONITORING IS IMPORTANT:
// ============================
// WebRTC servers need monitoring to:
// - Track server performance and resource usage
// - Identify connection patterns and bottlenecks
// - Monitor authentication success/failure rates
// - Detect potential issues before they become problems
// - Provide insights for capacity planning
//
// MONITORING FREQUENCY:
// ====================
// Statistics are logged every 30 seconds by default
// This provides a good balance between detail and performance
// More frequent monitoring can be enabled for debugging
//
// WHAT IS MONITORED:
// ==================
// - Active TURN servers (UDP, TCP, TLS variants)
// - Active STUN servers (UDP, TCP, TLS variants)
// - Server uptime and timestamp
// - Connection patterns and usage statistics
//
// EXTENSIBILITY:
// ==============
// This monitoring system can be extended to track:
// - Bandwidth usage per relay
// - Authentication success rates
// - Connection duration statistics
// - Geographic distribution of clients
// - Protocol preference patterns
func startMonitoring() {
	go func() {
		ticker := time.NewTicker(30 * time.Second) // Log stats every 30 seconds
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				logServerStats()
			}
		}
	}()
}

// logServerStats logs current server statistics
// This function provides a snapshot of server health and performance
//
// STATISTICS INCLUDED:
// ====================
// - Current timestamp for correlation with other logs
// - Number of active STUNTURN servers (UDP, TCP, TLS)
// - Server status and health indicators
//
// LOGGING FORMAT:
// ===============
// Statistics are logged in a structured format for easy parsing
// Each log entry includes timestamp and relevant metrics
// This facilitates automated monitoring and alerting
//
// USE CASES:
// ==========
// - Real-time server health monitoring
// - Performance trend analysis
// - Capacity planning and scaling decisions
// - Troubleshooting connection issues
// - Compliance and audit requirements
func logServerStats() {
	stunTurnLogger.Printf("=== SERVER STATISTICS ===")
	stunTurnLogger.Printf("Time: %s", time.Now().Format("2006-01-02 15:04:05"))
	stunTurnLogger.Printf("Active STUN/TURN servers: %d", countActiveSTUNTURNServers())
	stunTurnLogger.Printf("========================")
}

// countActiveSTUNTURNServers counts the number of active STUNTURN servers
// This function provides insight into which STUNTURN protocols are running
//
// STUNTURN SERVER TYPES:
// =====================
// - UDP STUNTURN: Standard discovery and relay service (port 3478)
// - TCP STUNTURN: Fallback discovery and relay service (port 3478)
// - TLS STUNTURN: Secure discovery and relay service (port 5349)
//
// DUAL FUNCTIONALITY:
// ==================
// Each STUNTURN server provides both STUN and TURN services:
// - STUN: Helps clients discover their public IP and NAT type
// - TURN: Provides relay services when direct connection fails
// - Combined: Single server handles all NAT traversal needs
//
// COUNTING LOGIC:
// ===============
// Only counts servers that were successfully initialized
// Nil servers (failed initialization) are not counted
// This provides accurate status of available services
//
// MONITORING VALUE:
// =================
// - Identifies which protocols are available to clients
// - Helps diagnose initialization failures
// - Provides insight into server configuration
// - Assists with troubleshooting connection issues
func countActiveSTUNTURNServers() int {
	count := 0
	servers := []*turn.Server{stunturnServer, stunturnTCPServer, stunturnTLSServer}
	for _, server := range servers {
		if server != nil {
			count++
		}
	}
	return count
}

// ============================================================================
// STUNTURN SERVER INITIALIZATION
// ============================================================================

// initializeSTUNTurnServer sets up the STUN/TURN servers with all protocol variants
// This is the main initialization function that creates all server instances
//
// IMPORTANT: TURN SERVERS INHERENTLY SUPPORT STUN FUNCTIONALITY
// ============================================================
// Every TURN server is also a STUN server. This is a fundamental aspect of the
// TURN protocol specification (RFC 5766). When you run a TURN server, it automatically
// handles both STUN and TURN requests on the same port.
//
// HOW IT WORKS:
// =============
//
//  1. Client sends STUN_BINDING_REQUEST → Server responds with STUN_BINDING_RESPONSE
//     (This is standard STUN functionality for NAT discovery)
//
//  2. Client sends TURN_ALLOCATE_REQUEST → Server responds with TURN_ALLOCATE_RESPONSE
//     (This is TURN functionality for relay allocation)
//
// 3. Same server, same port handles both protocols seamlessly
//
// WHY WE DON'T NEED SEPARATE STUN SERVERS:
// ========================================
// - TURN servers are supersets of STUN servers
// - Running separate STUN servers is redundant and wasteful
// - Single STUN/TURN server provides both services efficiently
// - Clients can use the same server for both STUN discovery and TURN relay
//
// CLIENT CONFIGURATION:
// ====================
// Clients should configure their ICE servers like this:
//
//	{
//	  urls: 'stun:YOUR_SERVER_IP:3478',        // STUN discovery
//	  urls: 'turn:YOUR_SERVER_IP:3478',        // TURN relay (same server!)
//	  username: 'your_username',
//	  credential: 'your_password'
//	}
//
// WHAT IS A STUN/TURN SERVER?
// ===========================
// STUN/TURN servers provide both STUN and TURN services:
// - STUN: Helps clients discover their public IP and NAT type
// - TURN: Provides relay services when direct connection fails
// - Combined: Single server handles all NAT traversal needs
//
// AUTHENTICATION SYSTEM:
// =====================
// STUN/TURN servers require authentication for TURN services:
// - STUN requests (discovery) are typically unauthenticated
// - TURN requests (relay) require username/password authentication
// - This prevents abuse of relay services while allowing free discovery
//
// PROTOCOL VARIANTS:
// ==================
// This server supports multiple transport protocols:
// - UDP: Standard protocol, fastest, works with most NATs
// - TCP: Fallback for UDP-blocked networks (corporate firewalls)
// - TLS: Encrypted connections for secure environments
//
// RELAY ADDRESS GENERATION:
// =========================
// STUN/TURN servers allocate relay addresses for clients:
// - Relay addresses must be reachable from the internet
// - Server needs to know its public IP for relay allocation
// - Without this, clients won't be able to connect to the relay
//
// THREADING MODEL:
// ================
// Multiple threads can be used to handle concurrent connections:
// - Each thread gets its own listener
// - Improves performance under high load
// - Prevents connection bottlenecks
func initializeSTUNTurnServer(publicIP, users, realm string, threadNum int, enableTCP, enableTLS bool) error {
	// ========================================================================
	// USER AUTHENTICATION SETUP
	// ========================================================================
	// Parse TURN user credentials from the command line argument
	// Format: "user1=pass1,user2=pass2"
	// This creates a map of username -> cryptographic auth key
	usersMap = make(map[string][]byte)

	// Use regex to parse username=password pairs
	// This regex finds all patterns like "username=password"
	// The regex (\w+)=(\w+) captures:
	// - Group 1: username (word characters)
	// - Group 2: password (word characters)
	for _, kv := range regexp.MustCompile(`(\w+)=(\w+)`).FindAllStringSubmatch(users, -1) {
		// Generate authentication key using TURN protocol specification
		// This creates a cryptographic key from username, realm, and password
		// The key is used to validate TURN requests from clients
		usersMap[kv[1]] = turn.GenerateAuthKey(kv[1], realm, kv[2])
		stunTurnLogger.Printf("Added TURN user: %s", kv[1])
	}

	// ========================================================================
	// RELAY ADDRESS GENERATOR
	// ========================================================================
	// This tells the TURN server what IP address to use for relay allocation
	// When a client requests a relay, the server will allocate an address on this IP
	// The publicIP must be reachable from the internet for relay to work
	relayAddressGenerator := &turn.RelayAddressGeneratorStatic{
		RelayAddress: net.ParseIP(publicIP), // Public IP for relay allocation
		Address:      "0.0.0.0",             // Listen on all interfaces
	}

	// ========================================================================
	// AUTHENTICATION HANDLER
	// ========================================================================
	// This function is called whenever a client tries to authenticate
	// It validates the username and returns the corresponding auth key
	// If authentication fails, the client cannot use relay services
	authHandler := createEnhancedAuthHandler(usersMap)

	// ========================================================================
	// SERVER INITIALIZATION SEQUENCE
	// ========================================================================
	// Initialize servers in order of importance and dependency
	// Each protocol variant serves different network environments

	// 2. UDP STUN/TURN server - main relay service, handles most WebRTC traffic
	// UDP is the standard protocol for STUN/TURN and works with most NAT types
	// It's the fastest and most efficient option
	if err := initializeUDPSTUNTurnServer(relayAddressGenerator, authHandler, realm, threadNum); err != nil {
		return fmt.Errorf("failed to initialize UDP STUN/TURN server: %w", err)
	}

	// 4. TCP STUN/TURN server (if enabled) - fallback relay service
	// TCP is used when UDP is blocked by firewalls or NATs
	// Common in corporate networks that block UDP traffic
	if enableTCP {
		if err := initializeTCPSTUNTurnServer(relayAddressGenerator, authHandler, realm, threadNum); err != nil {
			return fmt.Errorf("failed to initialize TCP STUN/TURN server: %w", err)
		}
	}

	// 6. TLS STUN/TURN server (if enabled) - secure relay service
	// TLS provides encrypted relay connections
	// Required for secure enterprise environments and browser compatibility
	if enableTLS {
		if err := initializeTLSSTUNTurnServer(relayAddressGenerator, authHandler, realm, threadNum); err != nil {
			return fmt.Errorf("failed to initialize TLS STUN/TURN server: %w", err)
		}
	}

	return nil
}

// ============================================================================
// UDP STUNTURN SERVER IMPLEMENTATION
// ============================================================================

// initializeUDPSTUNTurnServer sets up UDP STUN/TURN server
// This server handles both STUN discovery and TURN relay services over UDP
//
// WHY UDP FOR STUN/TURN?
// =====================
// UDP is the standard transport protocol for STUN/TURN servers because:
// - It's faster than TCP (no connection establishment overhead)
// - It works better with most NAT types
// - It's more efficient for real-time media (WebRTC's primary use case)
// - It handles packet loss gracefully (important for video/audio)
//
// DUAL FUNCTIONALITY:
// ==================
// This single server handles both STUN and TURN requests:
// - STUN_BINDING_REQUEST → STUN_BINDING_RESPONSE (NAT discovery)
// - TURN_ALLOCATE_REQUEST → TURN_ALLOCATE_RESPONSE (Relay allocation)
// - Same port, same server, seamless operation
//
// RELAY MECHANISM:
// ================
// When two WebRTC clients cannot connect directly (due to NAT/firewall):
// 1. Client A connects to STUN/TURN server and requests relay
// 2. STUN/TURN server allocates a relay address for Client A
// 3. Client B connects to STUN/TURN server and requests relay
// 4. STUN/TURN server allocates a relay address for Client B
// 5. STUN/TURN server forwards data between the two relay addresses
//
// THREADING MODEL:
// ================
// Multiple threads are used to handle concurrent connections:
// - Each thread gets its own UDP listener
// - This improves performance under high load
// - Prevents one slow connection from blocking others
// - Allows better CPU utilization on multi-core systems
//
// SOCKET OPTIONS:
// ===============
// SO_REUSEADDR: Allows multiple listeners to bind to the same port
// SO_BROADCAST: Enables broadcast capabilities for UDP
// These options are essential for proper UDP server operation
func initializeUDPSTUNTurnServer(relayGen *turn.RelayAddressGeneratorStatic, authHandler func(string, string, net.Addr) ([]byte, bool), realm string, threadNum int) error {
	// Create UDP address for the server
	// "0.0.0.0" means listen on all network interfaces
	// Port 3478 is the standard STUNTURN UDP port (IANA assigned)
	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:"+strconv.Itoa(turnUDPPort))
	if err != nil {
		return fmt.Errorf("failed to parse server address: %w", err)
	}

	// Create listener configuration with proper socket options for multithreading
	// SO_REUSEADDR allows multiple listeners to bind to the same port
	// SO_BROADCAST enables broadcast capabilities for UDP
	// These socket options are crucial for proper UDP server operation
	listenerConfig := &net.ListenConfig{
		Control: func(network, address string, conn syscall.RawConn) error {
			var operr error
			if err := conn.Control(func(fd uintptr) {
				// Set SO_REUSEADDR to allow multiple listeners on same port
				// This is essential when using multiple threads
				operr = syscall.SetsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
				if operr != nil {
					return
				}
				// Set SO_BROADCAST for UDP broadcast capabilities
				// This allows the server to handle broadcast packets
				operr = syscall.SetsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1)
			}); err != nil {
				return err
			}
			return operr
		},
	}

	// Create multiple UDP listeners for better performance
	// Each thread gets its own listener to handle concurrent connections
	// This prevents connection bottlenecks and improves throughput
	packetConnConfigs := make([]turn.PacketConnConfig, threadNum)
	stunTurnLogger.Printf("")
	for i := 0; i < threadNum; i++ {
		// Create UDP listener with proper socket options
		// Each listener runs on the same port but in a separate thread
		conn, err := listenerConfig.ListenPacket(context.Background(), addr.Network(), addr.String())
		if err != nil {
			return fmt.Errorf("failed to create UDP STUNTURN listener %d: %w", i, err)
		}

		// Wrap the connection with custom logging
		logger := NewSTUNTurnLogger(stunTurnLogger)
		customConn := NewLoggingPacketConn(conn, logger, fmt.Sprintf("UDP-%d", i))

		// Configure the packet connection with relay capabilities
		// Each listener is configured with the same relay address generator
		// This ensures consistent relay allocation across all threads
		packetConnConfigs[i] = turn.PacketConnConfig{
			PacketConn:            customConn, // Custom UDP connection with logging
			RelayAddressGenerator: relayGen,   // How to allocate relay addresses
		}
		stunTurnLogger.Printf("UDP STUNTURN server %d listening on %s", i, conn.LocalAddr().String())
	}

	// Create STUN/TURN server with authentication and relay capabilities
	// The server combines all UDP listeners into a single STUN/TURN server instance
	// This provides unified authentication and relay management
	// NOTE: This server automatically handles both STUN and TURN requests
	stunturnServer, err = turn.NewServer(turn.ServerConfig{
		Realm:             realm,             // Authentication realm
		AuthHandler:       authHandler,       // Authentication function
		PacketConnConfigs: packetConnConfigs, // UDP listeners
	})
	if err != nil {
		return fmt.Errorf("failed to create UDP STUN/TURN server: %w", err)
	}
	return nil
}

// ============================================================================
// TCP STUNTURN SERVER IMPLEMENTATION
// ============================================================================

// initializeTCPSTUNTURNServer sets up TCP STUN/TURN server
// This server handles both STUN discovery and TURN relay services over TCP
//
// WHY TCP FOR STUN/TURN?
// =====================
// TCP STUN/TURN is used as a fallback when UDP is blocked:
// - Corporate firewalls often block UDP traffic
// - Some NATs don't handle UDP well
// - TCP is more reliable in restrictive network environments
// - Required for enterprise networks with strict security policies
//
// DUAL FUNCTIONALITY:
// ==================
// This single server handles both STUN and TURN requests over TCP:
// - STUN_BINDING_REQUEST → STUN_BINDING_RESPONSE (NAT discovery)
// - TURN_ALLOCATE_REQUEST → TURN_ALLOCATE_RESPONSE (Relay allocation)
// - Same port, same server, seamless operation
//
// TCP vs UDP PERFORMANCE:
// =======================
// TCP STUN/TURN is slower than UDP STUN/TURN because:
// - TCP has connection establishment overhead (3-way handshake)
// - TCP has congestion control and flow control
// - TCP retransmits lost packets (good for reliability, bad for latency)
// - TCP has higher latency due to buffering
//
// WHEN TO USE TCP STUN/TURN:
// =========================
// Use TCP STUN/TURN when:
// - UDP is blocked by firewall/NAT
// - Network has strict security policies
// - Reliability is more important than latency
// - Working in corporate/enterprise environments
//
// RELAY MECHANISM (TCP):
// ======================
// TCP relay works similarly to UDP relay:
// 1. Client A establishes TCP connection to STUN/TURN server
// 2. STUN/TURN server allocates relay address for Client A
// 3. Client B establishes TCP connection to STUN/TURN server
// 4. STUN/TURN server allocates relay address for Client B
// 5. STUN/TURN server forwards data between the two TCP connections
//
// THREADING MODEL:
// ================
// Similar to UDP, multiple threads handle concurrent connections
// Each thread gets its own TCP listener for better performance
func initializeTCPSTUNTurnServer(relayGen *turn.RelayAddressGeneratorStatic, authHandler func(string, string, net.Addr) ([]byte, bool), realm string, threadNum int) error {
	// Create TCP address for the server
	// Same port as UDP (3478) but different protocol
	// "0.0.0.0" means listen on all network interfaces
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:"+strconv.Itoa(turnTCPPort))
	if err != nil {
		return fmt.Errorf("failed to parse server address: %w", err)
	}

	// Create listener configuration with proper socket options for multithreading
	// SO_REUSEADDR allows multiple listeners to bind to the same port
	// This is essential for multi-threaded TCP servers
	listenerConfig := &net.ListenConfig{
		Control: func(network, address string, conn syscall.RawConn) error {
			var operr error
			if err := conn.Control(func(fd uintptr) {
				// Set SO_REUSEADDR to allow multiple listeners on same port
				// This enables multiple threads to share the same port
				operr = syscall.SetsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
			}); err != nil {
				return err
			}
			return operr
		},
	}

	// Create multiple TCP listeners for better performance
	// Each thread gets its own listener to handle concurrent connections
	// This prevents connection bottlenecks and improves throughput
	listenerConfigs := make([]turn.ListenerConfig, threadNum)
	stunTurnLogger.Printf("")
	for i := 0; i < threadNum; i++ {
		// Create TCP listener with proper socket options
		// Each listener runs on the same port but in a separate thread
		listener, err := listenerConfig.Listen(context.Background(), addr.Network(), addr.String())
		if err != nil {
			return fmt.Errorf("failed to create TCP STUNTURN listener %d: %w", i, err)
		}

		// Wrap the listener with custom logging
		logger := NewSTUNTurnLogger(stunTurnLogger)
		customListener := NewLoggingListener(listener, logger, fmt.Sprintf("TCP-%d", i))

		// Configure the TCP listener with relay capabilities
		// Each listener is configured with the same relay address generator
		// This ensures consistent relay allocation across all threads
		listenerConfigs[i] = turn.ListenerConfig{
			Listener:              customListener, // Custom TCP listener with logging
			RelayAddressGenerator: relayGen,       // How to allocate relay addresses
		}
		stunTurnLogger.Printf("TCP STUNTURN server %d listening on %s", i, listener.Addr().String())
	}

	// Create STUNTURN server with TCP listeners
	// The server combines all TCP listeners into a single STUNTURN server instance
	// This provides unified authentication and relay management
	// NOTE: This server automatically handles both STUN and TURN requests
	stunturnTCPServer, err = turn.NewServer(turn.ServerConfig{
		Realm:           realm,           // Authentication realm
		AuthHandler:     authHandler,     // Authentication function
		ListenerConfigs: listenerConfigs, // TCP listeners
	})
	if err != nil {
		return fmt.Errorf("failed to create TCP STUNTURN server: %w", err)
	}
	return nil
}

// ============================================================================
// TLS STUNTURN SERVER IMPLEMENTATION
// ============================================================================

// initializeTLSSTUNTurnServer sets up TLS STUN/TURN server
// This server handles both STUN discovery and TURN relay services over encrypted TLS
//
// WHY TLS FOR STUN/TURN?
// =====================
// TLS STUN/TURN provides encrypted relay connections:
// - All relay traffic is encrypted end-to-end
// - Required for secure enterprise environments
// - Needed for WebRTC in browsers (HTTPS requirement)
// - Protects against man-in-the-middle attacks
// - Complies with security policies and regulations
//
// DUAL FUNCTIONALITY:
// ==================
// This single server handles both STUN and TURN requests over TLS:
// - STUN_BINDING_REQUEST → STUN_BINDING_RESPONSE (NAT discovery)
// - TURN_ALLOCATE_REQUEST → TURN_ALLOCATE_RESPONSE (Relay allocation)
// - Same port, same server, seamless operation with encryption
//
// SECURITY BENEFITS:
// ==================
// TLS encryption provides:
// - Confidentiality: Relay data cannot be intercepted
// - Integrity: Data cannot be modified in transit
// - Authentication: Server identity is verified
// - Forward secrecy: Past communications remain secure
//
// CERTIFICATE REQUIREMENTS:
// =========================
// TLS STUN/TURN requires SSL/TLS certificates:
// - Certificate must be valid and trusted
// - Domain name must match server hostname
// - Certificate must be in PEM format
// - Private key must be secure and accessible
//
// WHEN TO USE TLS STUN/TURN:
// ==========================
// Use TLS STUN/TURN when:
// - Security is a primary concern
// - Working in enterprise environments
// - WebRTC is used in browsers (HTTPS requirement)
// - Compliance with security policies is required
// - Protecting sensitive communications
//
// PERFORMANCE CONSIDERATIONS:
// ==========================
// TLS STUN/TURN has additional overhead:
// - TLS handshake adds connection establishment time
// - Encryption/decryption uses CPU resources
// - Slightly higher latency than TCP STUN/TURN
// - Still faster than TCP for most use cases
//
// RELAY MECHANISM (TLS):
// ======================
// TLS relay works like TCP relay but with encryption:
// 1. Client A establishes TLS connection to STUN/TURN server
// 2. STUN/TURN server allocates relay address for Client A
// 3. Client B establishes TLS connection to STUN/TURN server
// 4. STUN/TURN server allocates relay address for Client B
// 5. STUN/TURN server forwards encrypted data between connections
func initializeTLSSTUNTurnServer(relayGen *turn.RelayAddressGeneratorStatic, authHandler func(string, string, net.Addr) ([]byte, bool), realm string, threadNum int) error {
	// Check if SSL certificates exist (same as TLS STUN)
	// Certificates must be in the certs/ directory
	// fullchain.pem contains the certificate chain
	// privkey.pem contains the private key
	certFile := "certs/fullchain.pem"
	keyFile := "certs/privkey.pem"
	var err error

	// If certificates don't exist, skip TLS server
	// This allows the server to run without TLS if certificates are not available
	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		stunTurnLogger.Printf("SSL certificates not found. Skipping TLS STUNTURN server.")
		return nil
	}

	// Load TLS certificate and private key
	// The certificate must be valid and trusted by clients
	// The private key must be secure and accessible to the server
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return fmt.Errorf("failed to load TLS certificate: %w", err)
	}

	// Configure TLS settings
	// MinVersion ensures we use secure TLS versions
	// TLS 1.2 is the minimum recommended version for security
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert}, // Our SSL certificate
		MinVersion:   tls.VersionTLS12,        // Minimum TLS version (secure)
	}

	// Create TCP address for the server
	// Port 5349 is the standard STUNTURNS (STUNTURN over TLS) port
	// Different from standard STUNTURN port (3478) to distinguish protocols
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:"+strconv.Itoa(turnTLSPort))
	if err != nil {
		return fmt.Errorf("failed to parse server address: %w", err)
	}

	// Create listener configuration with proper socket options for multithreading
	// SO_REUSEADDR allows multiple listeners to bind to the same port
	// This is essential for multi-threaded TLS servers
	listenerConfig := &net.ListenConfig{
		Control: func(network, address string, conn syscall.RawConn) error {
			var operr error
			if err := conn.Control(func(fd uintptr) {
				// Set SO_REUSEADDR to allow multiple listeners on same port
				// This enables multiple threads to share the same port
				operr = syscall.SetsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
			}); err != nil {
				return err
			}
			return operr
		},
	}

	// Create multiple TLS listeners for better performance
	// Each thread gets its own listener to handle concurrent connections
	// This prevents connection bottlenecks and improves throughput
	listenerConfigs := make([]turn.ListenerConfig, threadNum)
	stunTurnLogger.Printf("")
	for i := 0; i < threadNum; i++ {
		// Create TCP listener with proper socket options first
		// TLS is built on top of TCP, so we start with a TCP listener
		tcpListener, err := listenerConfig.Listen(context.Background(), addr.Network(), addr.String())
		if err != nil {
			return fmt.Errorf("failed to create TCP listener for TLS STUNTURN %d: %w", i, err)
		}

		// Wrap TCP listener with TLS
		// This adds encryption to the TCP connection
		// All data transmitted through this listener will be encrypted
		tlsListener := tls.NewListener(tcpListener, tlsConfig)

		// Configure the TLS listener with relay capabilities
		// Each listener is configured with the same relay address generator
		// This ensures consistent relay allocation across all threads
		listenerConfigs[i] = turn.ListenerConfig{
			Listener:              tlsListener, // TLS connection
			RelayAddressGenerator: relayGen,    // How to allocate relay addresses
		}
		stunTurnLogger.Printf("TLS STUNTURN server %d listening on %s", i, tlsListener.Addr().String())
	}

	// Create STUNTURN server with TLS listeners
	// The server combines all TLS listeners into a single STUNTURN server instance
	// This provides unified authentication and relay management
	// NOTE: This server automatically handles both STUN and TURN requests
	stunturnTLSServer, err = turn.NewServer(turn.ServerConfig{
		Realm:           realm,           // Authentication realm
		AuthHandler:     authHandler,     // Authentication function
		ListenerConfigs: listenerConfigs, // TLS listeners
	})
	if err != nil {
		return fmt.Errorf("failed to create TLS STUNTURN server: %w", err)
	}
	return nil
}

// ============================================================================
// CUSTOM LOGGING HANDLERS
// ============================================================================

// CustomPacketConn wraps a net.PacketConn to add logging for all STUN/TURN packets
type CustomPacketConn struct {
	net.PacketConn
	connID string
}

func (c *CustomPacketConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	n, addr, err = c.PacketConn.ReadFrom(p)
	if err == nil && n > 0 {
		// Log all incoming packets (both STUN and TURN)
		stunTurnLogger.Printf("[%s] Received %d bytes from %s", c.connID, n, addr.String())

		// Try to identify STUN/TURN message type
		if n >= 20 { // Minimum STUN message size
			messageType := getSTUNTURNMessageType(p[:n])
			if messageType != "" {
				stunTurnLogger.Printf("[%s] %s request from %s", c.connID, messageType, addr.String())
			}
		}
	}
	return n, addr, err
}

func (c *CustomPacketConn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	n, err = c.PacketConn.WriteTo(p, addr)
	if err == nil && n > 0 {
		stunTurnLogger.Printf("[%s] Sent %d bytes to %s", c.connID, n, addr.String())

		// Try to identify STUN/TURN message type
		if n >= 20 { // Minimum STUN message size
			messageType := getSTUNTURNMessageType(p[:n])
			if messageType != "" {
				stunTurnLogger.Printf("[%s] %s response to %s", c.connID, messageType, addr.String())
			}
		}
	}
	return n, err
}

// CustomListener wraps a net.Listener to add logging for TCP connections
type CustomListener struct {
	net.Listener
	connID string
}

func (c *CustomListener) Accept() (net.Conn, error) {
	conn, err := c.Listener.Accept()
	if err == nil {
		stunTurnLogger.Printf("[%s] New TCP connection from %s", c.connID, conn.RemoteAddr().String())

		// Wrap the connection to log data
		conn = &CustomConn{
			Conn:   conn,
			connID: c.connID,
		}
	}
	return conn, err
}

// CustomConn wraps a net.Conn to add logging for TCP data
type CustomConn struct {
	net.Conn
	connID string
}

func (c *CustomConn) Read(b []byte) (n int, err error) {
	n, err = c.Conn.Read(b)
	if err == nil && n > 0 {
		stunTurnLogger.Printf("[%s] Received %d bytes from %s", c.connID, n, c.RemoteAddr().String())

		// Try to identify STUN/TURN message type
		if n >= 20 { // Minimum STUN message size
			messageType := getSTUNTURNMessageType(b[:n])
			if messageType != "" {
				stunTurnLogger.Printf("[%s] %s request from %s", c.connID, messageType, c.RemoteAddr().String())
			}
		}
	}
	return n, err
}

func (c *CustomConn) Write(b []byte) (n int, err error) {
	n, err = c.Conn.Write(b)
	if err == nil && n > 0 {
		stunTurnLogger.Printf("[%s] Sent %d bytes to %s", c.connID, n, c.RemoteAddr().String())

		// Try to identify STUN/TURN message type
		if n >= 20 { // Minimum STUN message size
			messageType := getSTUNTURNMessageType(b[:n])
			if messageType != "" {
				stunTurnLogger.Printf("[%s] %s response to %s", c.connID, messageType, c.RemoteAddr().String())
			}
		}
	}
	return n, err
}

// getSTUNTURNMessageType attempts to identify STUN/TURN message types
func getSTUNTURNMessageType(data []byte) string {
	if len(data) < 20 {
		return ""
	}

	// Check STUN magic cookie (0x2112A442)
	if len(data) >= 4 && data[0] == 0x21 && data[1] == 0x12 && data[2] == 0xA4 && data[3] == 0x42 {
		// Extract message type (bytes 4-5)
		if len(data) >= 6 {
			messageType := uint16(data[4])<<8 | uint16(data[5])
			switch messageType {
			case 0x0001:
				return "STUN_BINDING_REQUEST"
			case 0x0101:
				return "STUN_BINDING_RESPONSE"
			case 0x0111:
				return "STUN_BINDING_ERROR_RESPONSE"
			case 0x0003:
				return "TURN_ALLOCATE_REQUEST"
			case 0x0103:
				return "TURN_ALLOCATE_RESPONSE"
			case 0x0113:
				return "TURN_ALLOCATE_ERROR_RESPONSE"
			case 0x0004:
				return "TURN_REFRESH_REQUEST"
			case 0x0104:
				return "TURN_REFRESH_RESPONSE"
			case 0x0114:
				return "TURN_REFRESH_ERROR_RESPONSE"
			case 0x0006:
				return "TURN_SEND_REQUEST"
			case 0x0106:
				return "TURN_SEND_RESPONSE"
			case 0x0116:
				return "TURN_SEND_ERROR_RESPONSE"
			case 0x0007:
				return "TURN_DATA_REQUEST"
			case 0x0107:
				return "TURN_DATA_RESPONSE"
			case 0x0117:
				return "TURN_DATA_ERROR_RESPONSE"
			case 0x0008:
				return "TURN_CREATE_PERMISSION_REQUEST"
			case 0x0108:
				return "TURN_CREATE_PERMISSION_RESPONSE"
			case 0x0118:
				return "TURN_CREATE_PERMISSION_ERROR_RESPONSE"
			case 0x0009:
				return "TURN_CHANNEL_BIND_REQUEST"
			case 0x0109:
				return "TURN_CHANNEL_BIND_RESPONSE"
			case 0x0119:
				return "TURN_CHANNEL_BIND_ERROR_RESPONSE"
			default:
				return fmt.Sprintf("UNKNOWN_STUNTURN_0x%04X", messageType)
			}
		}
	}

	return ""
}

// ============================================================================
// ENHANCED AUTHENTICATION HANDLER
// ============================================================================

// createEnhancedAuthHandler creates an authentication handler with comprehensive logging
func createEnhancedAuthHandler(usersMap map[string][]byte) func(string, string, net.Addr) ([]byte, bool) {
	logger := NewSTUNTurnLogger(stunTurnLogger)

	return func(username string, realm string, srcAddr net.Addr) ([]byte, bool) {
		stunTurnLogger.Printf("Authentication attempt for user: %s from %s (realm: %s)", username, srcAddr.String(), realm)

		if key, ok := usersMap[username]; ok {
			logger.LogAuthentication(srcAddr, username, true)
			return key, true
		}

		logger.LogAuthentication(srcAddr, username, false)
		return nil, false
	}
}

// ============================================================================
// CONNECTION MONITORING
// ============================================================================

// startConnectionMonitoring starts a goroutine to monitor active connections
func startConnectionMonitoring() {
	go func() {
		ticker := time.NewTicker(60 * time.Second) // Log every minute
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				logConnectionStats()
			}
		}
	}()
}

// logConnectionStats logs current connection statistics
func logConnectionStats() {
	stunTurnLogger.Printf("=== CONNECTION STATISTICS ===")
	stunTurnLogger.Printf("Time: %s", time.Now().Format("2006-01-02 15:04:05"))
	stunTurnLogger.Printf("Active STUN/TURN servers: %d", countActiveSTUNTURNServers())
	stunTurnLogger.Printf("Server status: RUNNING")
	stunTurnLogger.Printf("=============================")
}

// ============================================================================
// ENHANCED STUN/TURN LOGGING
// ============================================================================

// STUNTurnLogger provides comprehensive logging for STUN/TURN server activity
type STUNTurnLogger struct {
	logger *log.Logger
}

// NewSTUNTurnLogger creates a new STUN/TURN logger
func NewSTUNTurnLogger(logger *log.Logger) *STUNTurnLogger {
	return &STUNTurnLogger{logger: logger}
}

// LogSTUNRequest logs STUN binding requests
func (l *STUNTurnLogger) LogSTUNRequest(srcAddr net.Addr, messageType string) {
	l.logger.Printf("STUN %s from %s", messageType, srcAddr.String())
}

// LogSTUNResponse logs STUN binding responses
func (l *STUNTurnLogger) LogSTUNResponse(dstAddr net.Addr, messageType string) {
	l.logger.Printf("STUN %s to %s", messageType, dstAddr.String())
}

// LogTURNRequest logs TURN requests (allocate, refresh, send, etc.)
func (l *STUNTurnLogger) LogTURNRequest(srcAddr net.Addr, messageType string, username string) {
	l.logger.Printf("TURN %s from %s (user: %s)", messageType, srcAddr.String(), username)
}

// LogTURNResponse logs TURN responses
func (l *STUNTurnLogger) LogTURNResponse(dstAddr net.Addr, messageType string, username string) {
	l.logger.Printf("TURN %s to %s (user: %s)", messageType, dstAddr.String(), username)
}

// LogAuthentication logs authentication attempts
func (l *STUNTurnLogger) LogAuthentication(srcAddr net.Addr, username string, success bool) {
	if success {
		l.logger.Printf("AUTH SUCCESS for user '%s' from %s", username, srcAddr.String())
	} else {
		l.logger.Printf("AUTH FAILED for user '%s' from %s", username, srcAddr.String())
	}
}

// LogConnection logs new connections
func (l *STUNTurnLogger) LogConnection(srcAddr net.Addr, protocol string) {
	l.logger.Printf("New %s connection from %s", protocol, srcAddr.String())
}

// LogRelayAllocation logs relay allocation events
func (l *STUNTurnLogger) LogRelayAllocation(srcAddr net.Addr, relayAddr net.Addr, username string) {
	l.logger.Printf("Relay allocated for user '%s' from %s -> %s", username, srcAddr.String(), relayAddr.String())
}

// LogDataTransfer logs data transfer events
func (l *STUNTurnLogger) LogDataTransfer(srcAddr net.Addr, dstAddr net.Addr, bytes int, protocol string) {
	l.logger.Printf("%s data transfer: %s -> %s (%d bytes)", protocol, srcAddr.String(), dstAddr.String(), bytes)
}

// ============================================================================
// CUSTOM PACKET HANDLERS
// ============================================================================

// LoggingPacketConn wraps a net.PacketConn to add comprehensive STUN/TURN logging
type LoggingPacketConn struct {
	net.PacketConn
	logger *STUNTurnLogger
	connID string
}

func NewLoggingPacketConn(conn net.PacketConn, logger *STUNTurnLogger, connID string) *LoggingPacketConn {
	return &LoggingPacketConn{
		PacketConn: conn,
		logger:     logger,
		connID:     connID,
	}
}

func (l *LoggingPacketConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	n, addr, err = l.PacketConn.ReadFrom(p)
	if err == nil && n > 0 {
		// Log the raw packet first
		l.logger.logger.Printf("[%s] Received %d bytes from %s", l.connID, n, addr.String())

		// Try to identify and log STUN/TURN message type
		if n >= 20 { // Minimum STUN message size
			messageType := parseSTUNTURNMessage(p[:n])
			if messageType != "" {
				if isSTUNMessage(messageType) {
					l.logger.LogSTUNRequest(addr, messageType)
				} else if isTURNMessage(messageType) {
					// For TURN messages, we'll log the request but username comes later in auth
					l.logger.LogTURNRequest(addr, messageType, "unknown")
				}
			}
		}
	}
	return n, addr, err
}

func (l *LoggingPacketConn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	n, err = l.PacketConn.WriteTo(p, addr)
	if err == nil && n > 0 {
		// Log the raw packet first
		l.logger.logger.Printf("[%s] Sent %d bytes to %s", l.connID, n, addr.String())

		// Try to identify and log STUN/TURN message type
		if n >= 20 { // Minimum STUN message size
			messageType := parseSTUNTURNMessage(p[:n])
			if messageType != "" {
				if isSTUNMessage(messageType) {
					l.logger.LogSTUNResponse(addr, messageType)
				} else if isTURNMessage(messageType) {
					l.logger.LogTURNResponse(addr, messageType, "unknown")
				}
			}
		}
	}
	return n, err
}

// LoggingListener wraps a net.Listener to add connection logging
type LoggingListener struct {
	net.Listener
	logger *STUNTurnLogger
	connID string
}

func NewLoggingListener(listener net.Listener, logger *STUNTurnLogger, connID string) *LoggingListener {
	return &LoggingListener{
		Listener: listener,
		logger:   logger,
		connID:   connID,
	}
}

func (l *LoggingListener) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err == nil {
		l.logger.LogConnection(conn.RemoteAddr(), "TCP")

		// Wrap the connection to log data transfer
		conn = &LoggingConn{
			Conn:   conn,
			logger: l.logger,
			connID: l.connID,
		}
	}
	return conn, err
}

// LoggingConn wraps a net.Conn to add data transfer logging
type LoggingConn struct {
	net.Conn
	logger *STUNTurnLogger
	connID string
}

func (l *LoggingConn) Read(b []byte) (n int, err error) {
	n, err = l.Conn.Read(b)
	if err == nil && n > 0 {
		l.logger.logger.Printf("[%s] Received %d bytes from %s", l.connID, n, l.RemoteAddr().String())

		// Try to identify STUN/TURN message type
		if n >= 20 {
			messageType := parseSTUNTURNMessage(b[:n])
			if messageType != "" {
				if isSTUNMessage(messageType) {
					l.logger.LogSTUNRequest(l.RemoteAddr(), messageType)
				} else if isTURNMessage(messageType) {
					l.logger.LogTURNRequest(l.RemoteAddr(), messageType, "unknown")
				}
			}
		}
	}
	return n, err
}

func (l *LoggingConn) Write(b []byte) (n int, err error) {
	n, err = l.Conn.Write(b)
	if err == nil && n > 0 {
		l.logger.logger.Printf("[%s] Sent %d bytes to %s", l.connID, n, l.RemoteAddr().String())

		// Try to identify STUN/TURN message type
		if n >= 20 {
			messageType := parseSTUNTURNMessage(b[:n])
			if messageType != "" {
				if isSTUNMessage(messageType) {
					l.logger.LogSTUNResponse(l.RemoteAddr(), messageType)
				} else if isTURNMessage(messageType) {
					l.logger.LogTURNResponse(l.RemoteAddr(), messageType, "unknown")
				}
			}
		}
	}
	return n, err
}

// ============================================================================
// STUN/TURN MESSAGE PARSING
// ============================================================================

// parseSTUNTURNMessage attempts to parse STUN/TURN message types from raw data
func parseSTUNTURNMessage(data []byte) string {
	if len(data) < 20 {
		return ""
	}

	// Check STUN magic cookie at bytes 4-7 (RFC 5389)
	if len(data) >= 8 && data[4] == 0x21 && data[5] == 0x12 && data[6] == 0xA4 && data[7] == 0x42 {
		// Extract message type (bytes 0-1)
		messageType := uint16(data[0])<<8 | uint16(data[1])
		return getMessageTypeName(messageType)
	}

	return ""
}

// getMessageTypeName returns the human-readable name for STUN/TURN message types
func getMessageTypeName(messageType uint16) string {
	switch messageType {
	// STUN Methods
	case 0x0001:
		return "STUN_BINDING_REQUEST"
	case 0x0101:
		return "STUN_BINDING_RESPONSE"
	case 0x0111:
		return "STUN_BINDING_ERROR_RESPONSE"

	// TURN Methods
	case 0x0003:
		return "TURN_ALLOCATE_REQUEST"
	case 0x0103:
		return "TURN_ALLOCATE_RESPONSE"
	case 0x0113:
		return "TURN_ALLOCATE_ERROR_RESPONSE"
	case 0x0004:
		return "TURN_REFRESH_REQUEST"
	case 0x0104:
		return "TURN_REFRESH_RESPONSE"
	case 0x0114:
		return "TURN_REFRESH_ERROR_RESPONSE"
	case 0x0006:
		return "TURN_SEND_REQUEST"
	case 0x0106:
		return "TURN_SEND_RESPONSE"
	case 0x0116:
		return "TURN_SEND_ERROR_RESPONSE"
	case 0x0007:
		return "TURN_DATA_REQUEST"
	case 0x0107:
		return "TURN_DATA_RESPONSE"
	case 0x0117:
		return "TURN_DATA_ERROR_RESPONSE"
	case 0x0008:
		return "TURN_CREATE_PERMISSION_REQUEST"
	case 0x0108:
		return "TURN_CREATE_PERMISSION_RESPONSE"
	case 0x0118:
		return "TURN_CREATE_PERMISSION_ERROR_RESPONSE"
	case 0x0009:
		return "TURN_CHANNEL_BIND_REQUEST"
	case 0x0109:
		return "TURN_CHANNEL_BIND_RESPONSE"
	case 0x0119:
		return "TURN_CHANNEL_BIND_ERROR_RESPONSE"
	default:
		return fmt.Sprintf("UNKNOWN_STUNTURN_0x%04X", messageType)
	}
}

// isSTUNMessage checks if a message type is a STUN message
func isSTUNMessage(messageType string) bool {
	return strings.HasPrefix(messageType, "STUN_")
}

// isTURNMessage checks if a message type is a TURN message
func isTURNMessage(messageType string) bool {
	return strings.HasPrefix(messageType, "TURN_")
}
