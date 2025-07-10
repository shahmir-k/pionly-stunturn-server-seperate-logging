/*
WebRTC WebSocket Signaling Handler
==================================

This file implements the WebSocket handler for WebRTC signaling.
WebSocket signaling is the communication channel that allows WebRTC peers
to exchange connection information and coordinate their peer-to-peer connections.

WHAT IS WEBSOCKET SIGNALING?
============================
WebRTC signaling is the process of exchanging connection information between peers:
- SDP (Session Description Protocol) offers and answers
- ICE (Interactive Connectivity Establishment) candidates
- Call control messages (join, call, hangup, etc.)

WebSocket provides a persistent, bidirectional communication channel for this exchange.

SIGNALING MESSAGE TYPES:
========================
This handler supports the following message types:
- join: User joins the signaling server
- activeUsers: Get list of currently active users
- call: Initiate a call to another user
- cancelCall: Cancel an outgoing call
- acceptCall: Accept an incoming call
- offer: Send SDP offer to peer
- answer: Send SDP answer to peer
- candidate: Send ICE candidate to peer
- hangUp: End an active call
- leave: User leaves the signaling server

CONNECTION LIFECYCLE:
=====================
1. Client connects via WebSocket upgrade
2. Client sends 'join' message to register
3. Client can send/receive signaling messages
4. Client sends 'leave' message or connection closes
5. Server cleans up user session

ERROR HANDLING:
==============
- WebSocket upgrade failures
- JSON parsing errors
- Connection read/write errors
- Unknown message types
- Graceful disconnection handling
*/

package webrtc

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// WebSocket upgrader configuration
// This handles the HTTP to WebSocket protocol upgrade
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024, // Buffer size for reading messages
	WriteBufferSize: 1024, // Buffer size for writing messages
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
		// In production, you should implement proper origin checking
		// Example: return r.Header.Get("Origin") == "https://yourdomain.com"
	},
}

// HandleWebSocket manages WebSocket connections for WebRTC signaling
// This function is the main entry point for all WebSocket signaling traffic
//
// CONNECTION MANAGEMENT:
// =====================
// - Handles WebSocket protocol upgrade from HTTP
// - Manages connection lifecycle (connect, message handling, disconnect)
// - Routes incoming messages to appropriate handlers
// - Provides comprehensive error handling and logging
//
// MESSAGE ROUTING:
// ================
// Each message type is routed to a specific handler function
// This modular design makes the code maintainable and extensible
// New message types can be easily added by extending the switch statement
//
// ERROR HANDLING:
// ===============
// - WebSocket upgrade failures are logged and handled gracefully
// - JSON parsing errors are logged and connection is closed
// - Unknown message types are logged for debugging
// - Connection errors trigger cleanup and disconnection
//
// LOGGING:
// ========
// All messages and errors are logged for debugging and monitoring
// This helps with troubleshooting connection issues
// Logs include message content and connection details
func HandleWebSocket(w http.ResponseWriter, r *http.Request, signalingLogger *log.Logger) {
	// Upgrade HTTP connection to WebSocket
	// This performs the WebSocket handshake and establishes the connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		signalingLogger.Println("Upgrade error:", err)
		return
	}

	// Ensure connection is closed when function exits
	// This prevents resource leaks and ensures proper cleanup
	defer func() {
		// Handle disconnection
		HandleDisconnect(conn, signalingLogger)
		conn.Close()
	}()

	// Main message handling loop
	// This loop continuously reads messages from the WebSocket connection
	// Each message is parsed and routed to the appropriate handler
	for {
		var msg SignalingMessage
		if err := conn.ReadJSON(&msg); err != nil {
			signalingLogger.Println("Read error:", err)
			break
		}

		// Add debug logging for all messages
		// This helps with debugging and understanding message flow
		//signalingLogger.Printf("Received message: %+v", msg)

		// Route message to appropriate handler based on message type
		// Each message type has its own handler function for modularity
		switch msg.Type {
		case "join":
			signalingLogger.Printf("Received: join From: %s To: %s", msg.Sender, msg.Receiver)
			// User joins the signaling server
			// Registers user and adds to active users list
			HandleJoin(conn, msg, signalingLogger)
		case "activeUsers":
			signalingLogger.Printf("Received: activeUsers From: %s To: %s", msg.Sender, msg.Receiver)
			// Get list of currently active users
			// Sends current user list to requesting client
			HandleActiveUsers(conn, msg, signalingLogger)
		case "call":
			signalingLogger.Printf("Received: call From: %s To: %s", msg.Sender, msg.Receiver)
			// Initiate a call to another user
			// Sends call request to target user
			HandleCall(conn, msg, signalingLogger)
		case "cancelCall":
			signalingLogger.Printf("Received: cancelCall From: %s To: %s", msg.Sender, msg.Receiver)
			// Cancel an outgoing call
			// Notifies target user that call was cancelled
			HandleCancelCall(conn, msg, signalingLogger)
		case "acceptCall":
			signalingLogger.Printf("Received: acceptCall From: %s To: %s", msg.Sender, msg.Receiver)
			// Accept an incoming call
			// Establishes call connection between users
			HandleAcceptCall(conn, msg, signalingLogger)
		case "offer":
			signalingLogger.Printf("Received: offer From: %s To: %s", msg.Sender, msg.Receiver)
			// Send SDP offer to peer
			// Initiates WebRTC connection establishment
			HandleOffer(conn, msg, signalingLogger)
		case "answer":
			signalingLogger.Printf("Received: answer From: %s To: %s", msg.Sender, msg.Receiver)
			// Send SDP answer to peer
			// Completes WebRTC connection establishment
			HandleAnswer(conn, msg, signalingLogger)
		case "candidate":
			signalingLogger.Printf("Received: candidate From: %s To: %s", msg.Sender, msg.Receiver)
			// Send ICE candidate to peer
			// Helps establish optimal peer-to-peer connection
			HandleIceCandidate(conn, msg, signalingLogger)
		case "hangUp":
			signalingLogger.Printf("Received: hangUp From: %s To: %s", msg.Sender, msg.Receiver)
			// End an active call
			// Terminates WebRTC connection and notifies both users
			HandleHangUp(conn, msg, signalingLogger)
		case "leave":
			signalingLogger.Printf("Received: leave From: %s To: %s", msg.Sender, msg.Receiver)
			// User leaves the signaling server
			// Cleans up user session and removes from active users
			HandleDisconnect(conn, signalingLogger)
			conn.Close()
		default:
			// Unknown message type
			// Log for debugging but don't break the connection
			signalingLogger.Printf("Unknown message type: %s From: %s To: %s", msg.Type, msg.Sender, msg.Receiver)
		}
	}
}
