#!/bin/bash

echo "Starting Unified WebRTC Server..."
echo

# Get the public IP address (you may need to modify this for your setup)
# For development, you can use your local IP or a public IP
PUBLIC_IP="YOUR_PUBLIC_IP_HERE"

# Check if public IP is set
# if [ "$PUBLIC_IP" = "YOUR_PUBLIC_IP_HERE" ]; then
#     echo "Autodetecting public IP..."
# fi


if ! command -v go &> /dev/null; then
    GREEN='\033[0;32m'
    NC='\033[0m' # No Color
    echo "Go is not installed. Please install Go (https://golang.org/dl/) and try again."
    echo
    echo
    echo "HINT: You can install Go by running the following commands:"
    echo
    echo -e "      ${GREEN}sudo apt install golang${NC}"
    echo
    echo
    echo
    echo
    echo "Or manually like this:"
    echo
    echo "Download latest version of Go from https://golang.org/dl/"
    echo "And open a terminal in the folder where you downloaded the Go package"
    echo
    echo "Remove any previous Go installation using the following command:"
    echo -e "      ${GREEN}sudo rm -rf /usr/local/go${NC}"
    echo
    echo "Then install latest version of Go (I used version 1.24.5 here) by running the following command:"
    echo -e "      ${GREEN}sudo tar -C /usr/local -xzf go1.24.5.linux-amd64.tar.gz${NC}"
    echo 
    echo "Update your PATH environment variable to include the Go bin directory:"
    echo -e "      ${GREEN}export PATH=\$PATH:/usr/local/go/bin${NC}"
    echo
    echo "Note: The above command only works for the current terminal session."
    echo "To make it permanent, add to your profile file using the following command:"
    echo -e "      ${GREEN}echo 'export PATH=\$PATH:/usr/local/go/bin' >> \$HOME/.profile${NC}"
    echo -e "      ${GREEN}source \$HOME/.profile${NC}"
    echo
    echo "Note: Changes may not apply until the next time you log in."
    echo "To apply changes immediately, run the source command above."
    echo
    echo "Verify the installation by running:"
    echo -e "      ${GREEN}go version${NC}"
    echo
    echo "If you see the version number, Go is installed correctly."
    echo "You may need to restart your computer for the changes to take effect across all terminals."
    echo
    exit 1
fi

# Download dependencies
echo "Downloading dependencies..."
go mod tidy

# Build the server
echo "Building server..."
go build -o go-server .
echo 

# Start the server


if [ "$PUBLIC_IP" = "YOUR_PUBLIC_IP_HERE" ]; then
    echo "Server will autodetect public IP..."
    echo
    echo "Starting server..."
    echo
    echo "Server will be available at:"
    echo "- Signaling: https://AUTODETECTED_IP:443/signal (HTTPS)"
    echo "- STUN/TURN: AUTODETECTED_IP:3478 (UDP/TCP)"
    echo "- STUN/TURN TLS: AUTODETECTED_IP:5349 (TLS)"
    echo
    echo "Press Ctrl+C to stop the server"
    echo
    ./go-server -turn-users="1ac96ad0a8374103e5c58441=drTJQZjbVFKpcXfn" -realm="yourdomain.com" -thread-num=2 -separate-logs=true 
else
    echo "Using public IP: $PUBLIC_IP"
    echo
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
fi