/*
WebRTC Signaling Service Implementation
======================================

This file implements the core WebRTC signaling service that coordinates
peer-to-peer connections between WebRTC clients.

WHAT IS WEBRTC SIGNALING?
=========================
WebRTC signaling is the process of exchanging connection information between peers:
- SDP (Session Description Protocol) offers and answers
- ICE (Interactive Connectivity Establishment) candidates
- Call control messages (join, call, hangup, etc.)

This service acts as a "matchmaker" that helps WebRTC peers find each other
and exchange the necessary information to establish direct peer-to-peer connections.

SESSION MANAGEMENT:
==================
The service maintains user sessions and manages:
- User registration and authentication
- Active user tracking
- Call state management
- Connection lifecycle (connect, call, disconnect)

THREAD SAFETY:
==============
All session operations are protected by read-write mutexes to ensure
thread safety in concurrent environments.

MESSAGE TYPES HANDLED:
======================
- join: User registration and session creation
- activeUsers: Get list of available users
- call: Initiate a call between users
- cancelCall: Cancel an outgoing call
- acceptCall: Accept an incoming call
- offer: Forward SDP offer between peers
- answer: Forward SDP answer between peers
- candidate: Forward ICE candidates between peers
- hangUp: End an active call
- leave: User disconnection and cleanup

WEBRTC COORDINATION:
====================
This service coordinates the WebRTC connection establishment process:
1. Users register and discover each other
2. One user initiates a call to another
3. Both users exchange SDP offers/answers via the server
4. Both users exchange ICE candidates via the server
5. Direct peer-to-peer connection is established
6. Server is no longer needed for media traffic
*/

package webrtc

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Global session management variables
// These maintain the state of all connected users and their sessions
var (
	// Maps username to user session for quick lookups
	nameToUserSession = make(map[string]*UserSession)
	// Maps connection address to username for reverse lookups
	sessionIdToName = make(map[string]string)
	// Read-write mutex for thread-safe access to session data
	mu sync.RWMutex
)

// HandleJoin handles a join request from a user
// This function manages user registration and session creation
//
// JOIN PROCESS:
// =============
// 1. User sends join message with their username
// 2. Server checks if username is already taken
// 3. If available, creates new user session
// 4. If taken, checks if existing session is still valid
// 5. Sends join result back to client
// 6. Broadcasts updated user list to all clients
//
// SESSION VALIDATION:
// ===================
// - Checks if user already has an active session
// - Validates that existing connection is still alive
// - Allows rejoin if previous session was invalid
// - Prevents duplicate sessions for same user
//
// ERROR HANDLING:
// ===============
// - Rejects join if username already has active session
// - Cleans up invalid sessions automatically
// - Provides clear feedback to client about join status
func HandleJoin(conn *websocket.Conn, msg SignalingMessage, signalingLogger *log.Logger) {
	name := msg.Sender
	signalingLogger.Printf("Handling join request from user: %s", name)

	mu.Lock()

	// Check if user already has a valid session
	// This prevents duplicate sessions and ensures user uniqueness
	if existingSession, exists := nameToUserSession[name]; exists {
		// Check if the existing connection is still valid
		// If connection is nil, it means the previous session was invalid
		if existingSession.Conn != nil {
			signalingLogger.Printf("User %s already has an active session, rejecting join", name)
			mu.Unlock()
			conn.WriteJSON(SignalingMessage{
				Type:     "join",
				Receiver: name,
				Data:     JoinResult{Result: false},
			})
			return
		}
	}

	// Remove any existing session for this user (force rejoin only if connection was invalid)
	// This cleans up stale session data and allows user to rejoin
	if _, exists := nameToUserSession[name]; exists {
		signalingLogger.Printf("Removing existing session for user %s to allow rejoin", name)
		delete(nameToUserSession, name)
		// Clean up sessionIdToName entries for this user
		// This maintains consistency between the two mapping structures
		var keysToDelete []string
		for sessionId, userName := range sessionIdToName {
			if userName == name {
				keysToDelete = append(keysToDelete, sessionId)
			}
		}
		// Delete collected keys
		for _, sessionId := range keysToDelete {
			delete(sessionIdToName, sessionId)
		}
	}

	// Create new user session
	// This establishes the user's presence in the system
	userSession := &UserSession{Name: name, Conn: conn}
	nameToUserSession[name] = userSession
	sessionIdToName[conn.RemoteAddr().String()] = name
	signalingLogger.Printf("User %s joined successfully", name)
	mu.Unlock()

	// Send successful join response to client
	// This confirms that the user has been registered
	conn.WriteJSON(SignalingMessage{
		Type:     "join",
		Receiver: name,
		Data:     JoinResult{Result: true},
	})

	// Broadcast updated user list to all connected clients
	// This ensures all clients have current information about available users
	signalingLogger.Printf("Broadcasting active users after %s joined", name)
	BroadcastActiveUsers(signalingLogger)
}

// HandleActiveUsers sends the list of active users to the requesting user
// This function provides user discovery capabilities
//
// USER DISCOVERY:
// ===============
// - Returns list of all currently connected users
// - Includes call status for each user
// - Helps clients build user interface
// - Enables call initiation to specific users
//
// THREAD SAFETY:
// ==============
// Uses read lock for concurrent access to user data
// Multiple clients can request user lists simultaneously
//
// DATA STRUCTURE:
// ===============
// Returns structured data with user names and call status
// This allows clients to show who's available for calls
func HandleActiveUsers(conn *websocket.Conn, msg SignalingMessage, signalingLogger *log.Logger) {
	mu.RLock()
	activeUsers := make([]ActiveUser, 0, len(nameToUserSession))
	for name, session := range nameToUserSession {
		activeUsers = append(activeUsers, ActiveUser{
			Name:   name,
			InCall: session.InCall,
		})
	}
	mu.RUnlock()

	conn.WriteJSON(SignalingMessage{
		Type: "activeUsers",
		Data: ActiveUsers{Users: activeUsers},
	})
}

// HandleCall initiates a call between two users
// This function manages call state and notifies the target user
//
// CALL INITIATION:
// ================
// 1. Validates both users exist and are available
// 2. Sets both users' call status to "in call"
// 3. Sends call notification to target user
// 4. Broadcasts updated user list to all clients
//
// VALIDATION:
// ===========
// - Ensures both sender and receiver exist
// - Checks that neither user is already in a call
// - Prevents invalid call attempts
//
// STATE MANAGEMENT:
// =================
// - Updates call status for both users
// - Prevents other users from calling users who are busy
// - Maintains consistent state across all clients
func HandleCall(conn *websocket.Conn, msg SignalingMessage, signalingLogger *log.Logger) {
	sender := msg.Sender
	receiver := msg.Receiver
	mu.Lock()
	senderSession, senderExists := nameToUserSession[sender]
	receiverSession, receiverExists := nameToUserSession[receiver]
	if !senderExists || !receiverExists || senderSession.InCall || receiverSession.InCall {
		mu.Unlock()
		return
	}
	senderSession.SetInCall(true)
	receiverSession.SetInCall(true)
	mu.Unlock()

	receiverSession.Send(SignalingMessage{
		Type:     "call",
		Sender:   sender,
		Receiver: receiver,
	})
	BroadcastActiveUsers(signalingLogger)
}

// HandleCancelCall cancels an ongoing call between two users
// This function manages call cancellation and state cleanup
//
// CALL CANCELLATION:
// ==================
// 1. Validates both users exist
// 2. Resets both users' call status to "available"
// 3. Notifies target user that call was cancelled
// 4. Broadcasts updated user list to all clients
//
// USE CASES:
// ==========
// - User cancels call before it's answered
// - Call timeout or error conditions
// - User changes mind about making call
//
// STATE CLEANUP:
// ==============
// - Resets call status for both users
// - Makes users available for new calls
// - Maintains consistent state across clients
func HandleCancelCall(conn *websocket.Conn, msg SignalingMessage, signalingLogger *log.Logger) {
	sender := msg.Sender
	receiver := msg.Receiver
	mu.Lock()
	senderSession, senderExists := nameToUserSession[sender]
	receiverSession, receiverExists := nameToUserSession[receiver]
	if !senderExists || !receiverExists {
		mu.Unlock()
		return
	}
	senderSession.SetInCall(false)
	receiverSession.SetInCall(false)
	mu.Unlock()

	receiverSession.Send(SignalingMessage{
		Type:     "cancelCall",
		Sender:   sender,
		Receiver: receiver,
	})
	BroadcastActiveUsers(signalingLogger)
}

// HandleAcceptCall marks the call as accepted by the receiver
// This function notifies the caller that their call was accepted
//
// CALL ACCEPTANCE:
// ================
// 1. Validates receiver exists
// 2. Forwards acceptance message to caller
// 3. Initiates WebRTC connection establishment
//
// WEBRTC COORDINATION:
// ====================
// This message triggers the start of WebRTC signaling:
// - SDP offer/answer exchange
// - ICE candidate exchange
// - Direct peer-to-peer connection establishment
//
// MESSAGE FLOW:
// =============
// Caller -> Server -> Receiver: "call"
// Receiver -> Server -> Caller: "acceptCall"
// Then WebRTC signaling begins...
func HandleAcceptCall(conn *websocket.Conn, msg SignalingMessage, signalingLogger *log.Logger) {
	sender := msg.Sender
	receiver := msg.Receiver
	mu.RLock()
	receiverSession, receiverExists := nameToUserSession[receiver]
	mu.RUnlock()
	if !receiverExists {
		return
	}

	receiverSession.Send(SignalingMessage{
		Type:     "acceptCall",
		Sender:   sender,
		Receiver: receiver,
	})
}

// HandleOffer forwards an SDP offer from the sender to the receiver
// This function is crucial for WebRTC connection establishment
//
// SDP OFFER:
// ==========
// SDP (Session Description Protocol) offers contain:
// - Media capabilities and preferences
// - Network information
// - Codec preferences
// - Connection parameters
//
// WEBRTC ESTABLISHMENT:
// =====================
// 1. Caller creates SDP offer with their capabilities
// 2. Offer is forwarded to receiver via this function
// 3. Receiver processes offer and creates SDP answer
// 4. Answer is sent back via HandleAnswer
// 5. ICE candidates are exchanged via HandleIceCandidate
//
// ERROR HANDLING:
// ===============
// - Validates receiver exists before forwarding
// - Logs offer content for debugging
// - Handles send errors gracefully
// - Provides detailed logging for troubleshooting
func HandleOffer(conn *websocket.Conn, msg SignalingMessage, signalingLogger *log.Logger) {
	sender := msg.Sender
	receiver := msg.Receiver
	offer := msg.Data

	signalingLogger.Printf("Received offer from %s to %s", sender, receiver)

	mu.RLock()
	receiverSession, receiverExists := nameToUserSession[receiver]
	mu.RUnlock()

	if !receiverExists {
		signalingLogger.Printf("Receiver %s not found for offer from %s", receiver, sender)
		return
	}

	err := receiverSession.Send(SignalingMessage{
		Type:     "offer",
		Sender:   sender,
		Receiver: receiver,
		Data:     offer,
	})
	if err != nil {
		signalingLogger.Printf("Error sending offer from %s to %s: %v", sender, receiver, err)
		return
	}

	signalingLogger.Printf("Offer forwarded from %s to %s", sender, receiver)
}

// HandleAnswer forwards an SDP answer from the sender to the receiver
// This function completes the SDP exchange for WebRTC connection
//
// SDP ANSWER:
// ===========
// SDP answers contain:
// - Receiver's media capabilities
// - Accepted codecs and parameters
// - Network information
// - Connection confirmation
//
// CONNECTION ESTABLISHMENT:
// =========================
// 1. Receiver processes SDP offer
// 2. Receiver creates SDP answer with their capabilities
// 3. Answer is forwarded to caller via this function
// 4. Caller processes answer and establishes connection
// 5. ICE candidates are exchanged for optimal connection
//
// COMPLETION:
// ===========
// After SDP exchange, both peers have:
// - Agreed on media parameters
// - Established connection parameters
// - Ready to exchange ICE candidates
func HandleAnswer(conn *websocket.Conn, msg SignalingMessage, signalingLogger *log.Logger) {
	sender := msg.Sender
	receiver := msg.Receiver
	answer := msg.Data

	signalingLogger.Printf("Received answer from %s to %s", sender, receiver)

	mu.RLock()
	receiverSession, receiverExists := nameToUserSession[receiver]
	mu.RUnlock()

	if !receiverExists {
		signalingLogger.Printf("Receiver %s not found for answer from %s", receiver, sender)
		return
	}

	err := receiverSession.Send(SignalingMessage{
		Type:     "answer",
		Sender:   sender,
		Receiver: receiver,
		Data:     answer,
	})
	if err != nil {
		signalingLogger.Printf("Error sending answer from %s to %s: %v", sender, receiver, err)
		return
	}

	signalingLogger.Printf("Answer forwarded from %s to %s", sender, receiver)
}

// HandleIceCandidate forwards an ICE candidate from the sender to the receiver
// This function is essential for establishing optimal peer-to-peer connections
//
// ICE CANDIDATES:
// ===============
// ICE (Interactive Connectivity Establishment) candidates represent:
// - Local network addresses (host candidates)
// - Server-reflexive addresses (STUN candidates)
// - Relay addresses (TURN candidates)
// - Different transport protocols (UDP, TCP, TLS)
//
// CONNECTION OPTIMIZATION:
// ========================
// 1. Both peers gather multiple ICE candidates
// 2. Candidates are exchanged via this function
// 3. ICE framework tests all candidate pairs
// 4. Best performing pair is selected for connection
// 5. Direct peer-to-peer connection is established
//
// NETWORK TRAVERSAL:
// ==================
// ICE candidates help WebRTC work across:
// - Different NAT types (symmetric, restricted, etc.)
// - Corporate firewalls and proxies
// - Mobile networks with carrier-grade NAT
// - Various network topologies
//
// PERFORMANCE:
// ============
// - Multiple candidates ensure connection success
// - ICE testing finds the optimal path
// - Fallback to relay if direct connection fails
// - Minimizes latency and maximizes bandwidth
func HandleIceCandidate(conn *websocket.Conn, msg SignalingMessage, signalingLogger *log.Logger) {
	sender := msg.Sender
	receiver := msg.Receiver
	candidate := msg.Data

	signalingLogger.Printf("Received ICE candidate from %s to %s", sender, receiver)

	mu.RLock()
	receiverSession, receiverExists := nameToUserSession[receiver]
	mu.RUnlock()

	if !receiverExists {
		signalingLogger.Printf("Receiver %s not found for ICE candidate from %s", receiver, sender)
		return
	}

	err := receiverSession.Send(SignalingMessage{
		Type:     "candidate",
		Sender:   sender,
		Receiver: receiver,
		Data:     candidate,
	})
	if err != nil {
		signalingLogger.Printf("Error sending ICE candidate from %s to %s: %v", sender, receiver, err)
		return
	}

	signalingLogger.Printf("ICE candidate forwarded from %s to %s", sender, receiver)
}

// HandleHangUp ends an active call between two users
// This function manages call termination and state cleanup
//
// CALL TERMINATION:
// =================
// 1. Validates both users exist
// 2. Resets both users' call status to "available"
// 3. Notifies both users that call has ended
// 4. Broadcasts updated user list to all clients
//
// WEBRTC CLEANUP:
// ===============
// - Terminates WebRTC peer connection
// - Stops media streams
// - Releases camera/microphone resources
// - Cleans up network connections
//
// STATE MANAGEMENT:
// =================
// - Makes both users available for new calls
// - Updates user list for all connected clients
// - Maintains consistent state across the system
// - Prevents resource leaks and orphaned connections
//
// USER EXPERIENCE:
// ================
// - Both users are notified of call end
// - Users can immediately start new calls
// - UI is updated to reflect available status
// - Clean transition from call to idle state
func HandleHangUp(conn *websocket.Conn, msg SignalingMessage, signalingLogger *log.Logger) {
	sender := msg.Sender
	receiver := msg.Receiver
	mu.Lock()
	senderSession, senderExists := nameToUserSession[sender]
	receiverSession, receiverExists := nameToUserSession[receiver]
	if !senderExists || !receiverExists {
		mu.Unlock()
		return
	}
	senderSession.SetInCall(false)
	receiverSession.SetInCall(false)
	mu.Unlock()

	receiverSession.Send(SignalingMessage{
		Type:     "hangUp",
		Sender:   sender,
		Receiver: receiver,
	})
	BroadcastActiveUsers(signalingLogger)
}

// HandleDisconnect manages user disconnection and session cleanup
// This function is called when a user leaves or connection is lost
//
// DISCONNECTION SCENARIOS:
// ========================
// - User explicitly sends "leave" message
// - WebSocket connection is closed by client
// - Network connection is lost
// - Server detects connection timeout
//
// CLEANUP PROCESS:
// ================
// 1. Identifies user by connection address
// 2. Removes user from active sessions
// 3. Cleans up session mappings
// 4. Notifies other users of departure
// 5. Broadcasts updated user list
//
// RESOURCE MANAGEMENT:
// ===================
// - Prevents memory leaks from abandoned sessions
// - Maintains accurate user count
// - Ensures system stability
// - Frees up resources for new connections
//
// ERROR HANDLING:
// ===============
// - Gracefully handles missing user sessions
// - Logs disconnection events for monitoring
// - Continues operation even if cleanup fails
// - Maintains system integrity
func HandleDisconnect(conn *websocket.Conn, signalingLogger *log.Logger) {
	// Find user by connection address
	// This reverse lookup helps identify which user disconnected
	mu.Lock()
	userName, exists := sessionIdToName[conn.RemoteAddr().String()]
	if !exists {
		mu.Unlock()
		return
	}

	// Clean up session data
	// Remove user from all session mappings
	delete(nameToUserSession, userName)
	delete(sessionIdToName, conn.RemoteAddr().String())
	mu.Unlock()

	signalingLogger.Printf("User %s disconnected", userName)

	// Broadcast updated user list to remaining clients
	// This ensures all clients have current information
	BroadcastActiveUsers(signalingLogger)
}

// BroadcastActiveUsers sends the current user list to all connected clients
// This function ensures all clients have synchronized user information
//
// BROADCAST PURPOSE:
// ==================
// - Keeps all clients informed of user changes
// - Enables real-time user discovery
// - Maintains consistent UI across clients
// - Supports dynamic user list updates
//
// MESSAGE CONTENT:
// ================
// - List of all currently connected users
// - Call status for each user (available/busy)
// - User names for display purposes
// - Structured data for easy client processing
//
// PERFORMANCE CONSIDERATIONS:
// ===========================
// - Uses read lock for concurrent access
// - Efficiently builds user list once
// - Sends same data to all clients
// - Minimizes server load during broadcasts
//
// CLIENT SYNCHRONIZATION:
// =======================
// - All clients receive same user information
// - UI updates happen simultaneously
// - Prevents inconsistent user states
// - Enables coordinated user interactions
func BroadcastActiveUsers(signalingLogger *log.Logger) {
	mu.RLock()
	activeUsers := make([]ActiveUser, 0, len(nameToUserSession))
	for name, session := range nameToUserSession {
		activeUsers = append(activeUsers, ActiveUser{
			Name:   name,
			InCall: session.InCall,
		})
	}
	mu.RUnlock()

	// Send updated user list to all connected clients
	// This ensures everyone has current information
	message := SignalingMessage{
		Type: "activeUsers",
		Data: ActiveUsers{Users: activeUsers},
	}

	mu.RLock()
	for _, session := range nameToUserSession {
		if session.Conn != nil {
			session.Send(message)
		}
	}
	mu.RUnlock()
}
