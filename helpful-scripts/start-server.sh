#!/bin/bash

echo "Starting Unified WebRTC Server..."
echo

# Get the public IP address (you may need to modify this for your setup)
# For development, you can use your local IP or a public IP
PUBLIC_IP="YOUR_PUBLIC_IP_HERE"

# Check if public IP is set
if [ "$PUBLIC_IP" = "YOUR_PUBLIC_IP_HERE" ]; then
    echo "ERROR: Please set your public IP address in this script"
    echo "Edit this file and replace YOUR_PUBLIC_IP_HERE with your actual public IP"
    exit 1
fi

echo "Using public IP: $PUBLIC_IP"
echo

# Download dependencies
echo "Downloading dependencies..."
cd ..
go mod tidy

# Build the server
echo "Building server..."
go build -o go-server .

# Start the server
echo "Starting server..."
echo
echo "Server will be available at:"
echo "- Signaling: https://$PUBLIC_IP:443/signal (HTTPS)"
echo "- STUN/TURN: $PUBLIC_IP:3478 (UDP/TCP)"
echo "- STUN/TURN TLS: $PUBLIC_IP:5349 (TLS)"
echo
echo "Press Ctrl+C to stop the server"
echo

./go-server -public-ip="$PUBLIC_IP" -turn-users="1ac96ad0a8374103e5c58441=drTJQZjbVFKpcXfn" -realm="yourdomain.com" -thread-num=2 -separate-logs=true 