#!/bin/bash

# WebRTC Server Monitoring Script
# This script demonstrates how to run the WebRTC server with separate logging
# and monitor STUN/TURN and signaling services in different terminals

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PUBLIC_IP=${1:-"YOUR_PUBLIC_IP"}
TURN_USERS=${2:-"1ac96ad0a8374103e5c58441=drTJQZjbVFKpcXfn"}
LOG_DIR="../logs"
STUN_TURN_LOG="$LOG_DIR/stun-turn.log"
SIGNALING_LOG="$LOG_DIR/signaling.log"

echo -e "${BLUE}WebRTC Server Monitoring Setup${NC}"
echo "=================================="
echo -e "Public IP: ${GREEN}$PUBLIC_IP${NC}"
echo -e "TURN Users: ${GREEN}$TURN_USERS${NC}"
echo -e "Log Directory: ${GREEN}$LOG_DIR${NC}"
echo ""

# Check if public IP is set
if [ "$PUBLIC_IP" = "YOUR_PUBLIC_IP" ]; then
    echo -e "${RED}Error: Please provide your public IP address as the first argument${NC}"
    echo "Usage: $0 <PUBLIC_IP> [TURN_USERS] [METHOD]"
    echo "Example: $0 203.0.113.1 'user1=pass1,user2=pass2' logs"
    echo ""
    echo "Available methods: logs, pipes, tmux, screen"
    exit 1
fi

# Create log directory
mkdir -p "$LOG_DIR"

# Function to show usage
show_usage() {
    echo ""
    echo -e "${YELLOW}Available monitoring methods:${NC}"
    echo ""
    echo "1. ${GREEN}logs${NC} - Separate Log Files (default)"
    echo "   Run: $0 $PUBLIC_IP '$TURN_USERS' logs"
    echo "   Then in separate terminals:"
    echo "   Terminal 1: tail -f $STUN_TURN_LOG"
    echo "   Terminal 2: tail -f $SIGNALING_LOG"
    echo ""
    echo "2. ${GREEN}pipes${NC} - Named Pipes (Real-time)"
    echo "   Run: $0 $PUBLIC_IP '$TURN_USERS' pipes"
    echo "   Then in separate terminals:"
    echo "   Terminal 1: cat $LOG_DIR/stun-turn-pipe"
    echo "   Terminal 2: cat $LOG_DIR/signaling-pipe"
    echo ""
    echo "3. ${GREEN}tmux${NC} - tmux Session"
    echo "   Run: $0 $PUBLIC_IP '$TURN_USERS' tmux"
    echo ""
    echo "4. ${GREEN}screen${NC} - GNU screen"
    echo "   Run: $0 $PUBLIC_IP '$TURN_USERS' screen"
    echo ""
}

# Function to run with log files
run_with_log_files() {
    echo -e "${GREEN}Starting WebRTC server with separate log files...${NC}"
    echo ""
    echo "Server will be started in the background."
    echo "Logs will be written to:"
    echo "  STUN/TURN: $STUN_TURN_LOG"
    echo "  Signaling: $SIGNALING_LOG"
    echo ""
    echo -e "${YELLOW}To monitor in separate terminals, run:${NC}"
    echo "  Terminal 1: tail -f $STUN_TURN_LOG"
    echo "  Terminal 2: tail -f $SIGNALING_LOG"
    echo ""
    echo -e "${YELLOW}To stop the server:${NC}"
    echo "  pkill -f go-server"
    echo ""
    
    # Start the server in background
    cd ..
    ./go-server \
        -public-ip "$PUBLIC_IP" \
        -turn-users "$TURN_USERS" \
        -separate-logs=true \
        -stun-turn-log="$STUN_TURN_LOG" \
        -signaling-log="$SIGNALING_LOG" &
    
    SERVER_PID=$!
    echo -e "${GREEN}Server started with PID: $SERVER_PID${NC}"
    echo "Server is running in the background."
    echo ""
    echo -e "${YELLOW}Press Enter to stop the server...${NC}"
    read
    echo -e "${RED}Stopping server...${NC}"
    kill $SERVER_PID 2>/dev/null || true
    echo -e "${GREEN}Server stopped.${NC}"
}

# Function to run with named pipes
run_with_pipes() {
    echo -e "${GREEN}Setting up named pipes for real-time monitoring...${NC}"
    
    # Create named pipes
    mkfifo "$LOG_DIR/stun-turn-pipe" 2>/dev/null || true
    mkfifo "$LOG_DIR/signaling-pipe" 2>/dev/null || true
    
    echo "Named pipes created:"
    echo "  STUN/TURN: $LOG_DIR/stun-turn-pipe"
    echo "  Signaling: $LOG_DIR/signaling-pipe"
    echo ""
    echo -e "${YELLOW}In separate terminals, run:${NC}"
    echo "  Terminal 1: cat $LOG_DIR/stun-turn-pipe"
    echo "  Terminal 2: cat $LOG_DIR/signaling-pipe"
    echo ""
    echo -e "${YELLOW}Press Enter when ready to start the server...${NC}"
    read
    
    # Start the server
    cd ..
    ./go-server \
        -public-ip "$PUBLIC_IP" \
        -turn-users "$TURN_USERS" \
        -separate-logs=true \
        -stun-turn-log="$LOG_DIR/stun-turn-pipe" \
        -signaling-log="$LOG_DIR/signaling-pipe"
}

# Function to run with tmux
run_with_tmux() {
    echo -e "${GREEN}Setting up tmux session for monitoring...${NC}"
    
    # Check if tmux is installed
    if ! command -v tmux &> /dev/null; then
        echo -e "${RED}Error: tmux is not installed${NC}"
        echo "Please install tmux or use another method."
        exit 1
    fi
    
    # Create tmux session
    tmux new-session -d -s webrtc-monitor
    
    # Split the window horizontally
    tmux split-window -h -t webrtc-monitor:0
    
    # Set up the left pane (STUN/TURN monitoring)
    tmux send-keys -t webrtc-monitor:0.0 "echo 'STUN/TURN Monitor - Press Ctrl+C to exit'; tail -f $STUN_TURN_LOG" Enter
    
    # Set up the right pane (signaling monitoring)
    tmux send-keys -t webrtc-monitor:0.1 "echo 'Signaling Monitor - Press Ctrl+C to exit'; tail -f $SIGNALING_LOG" Enter
    
    # Start the server in the background
    cd ..
    ./go-server \
        -public-ip "$PUBLIC_IP" \
        -turn-users "$TURN_USERS" \
        -separate-logs=true \
        -stun-turn-log="$STUN_TURN_LOG" \
        -signaling-log="$SIGNALING_LOG" &
    
    SERVER_PID=$!
    echo -e "${GREEN}Server started with PID: $SERVER_PID${NC}"
    echo -e "${GREEN}tmux session created: webrtc-monitor${NC}"
    echo ""
    echo -e "${YELLOW}To attach to the tmux session:${NC}"
    echo "  tmux attach-session -t webrtc-monitor"
    echo ""
    echo -e "${YELLOW}tmux commands:${NC}"
    echo "  Ctrl+B, Arrow Keys: Switch between panes"
    echo "  Ctrl+B, D: Detach from session"
    echo "  Ctrl+B, &: Kill current pane"
    echo ""
    echo -e "${YELLOW}To stop the server:${NC}"
    echo "  kill $SERVER_PID"
    echo "  tmux kill-session -t webrtc-monitor"
    echo ""
    echo -e "${YELLOW}Attaching to tmux session now...${NC}"
    tmux attach-session -t webrtc-monitor
}

# Function to run with screen
run_with_screen() {
    echo -e "${GREEN}Setting up GNU screen session for monitoring...${NC}"
    
    # Check if screen is installed
    if ! command -v screen &> /dev/null; then
        echo -e "${RED}Error: GNU screen is not installed${NC}"
        echo "Please install screen or use another method."
        exit 1
    fi
    
    # Create screen session
    screen -dmS webrtc-monitor
    
    # Create the first window for STUN/TURN monitoring
    screen -S webrtc-monitor -X screen -t stun-turn
    screen -S webrtc-monitor -p stun-turn -X stuff "echo 'STUN/TURN Monitor'; tail -f $STUN_TURN_LOG
"
    
    # Create the second window for signaling monitoring
    screen -S webrtc-monitor -X screen -t signaling
    screen -S webrtc-monitor -p signaling -X stuff "echo 'Signaling Monitor'; tail -f $SIGNALING_LOG
"
    
    # Start the server in the background
    cd ..
    ./go-server \
        -public-ip "$PUBLIC_IP" \
        -turn-users "$TURN_USERS" \
        -separate-logs=true \
        -stun-turn-log="$STUN_TURN_LOG" \
        -signaling-log="$SIGNALING_LOG" &
    
    SERVER_PID=$!
    echo -e "${GREEN}Server started with PID: $SERVER_PID${NC}"
    echo -e "${GREEN}screen session created: webrtc-monitor${NC}"
    echo ""
    echo -e "${YELLOW}To attach to the screen session:${NC}"
    echo "  screen -r webrtc-monitor"
    echo ""
    echo -e "${YELLOW}screen commands:${NC}"
    echo "  Ctrl+A, N: Next window"
    echo "  Ctrl+A, P: Previous window"
    echo "  Ctrl+A, D: Detach from session"
    echo "  Ctrl+A, K: Kill current window"
    echo ""
    echo -e "${YELLOW}To stop the server:${NC}"
    echo "  kill $SERVER_PID"
    echo "  screen -S webrtc-monitor -X quit"
    echo ""
    echo -e "${YELLOW}Attaching to screen session now...${NC}"
    screen -r webrtc-monitor
}

# Main script logic
case "${3:-logs}" in
    "logs")
        run_with_log_files
        ;;
    "pipes")
        run_with_pipes
        ;;
    "tmux")
        run_with_tmux
        ;;
    "screen")
        run_with_screen
        ;;
    "help"|"-h"|"--help")
        show_usage
        ;;
    *)
        echo -e "${RED}Unknown method: $3${NC}"
        show_usage
        exit 1
        ;;
esac 