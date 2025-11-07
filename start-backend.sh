#!/bin/bash

# Start the Go backend server
cd "$(dirname "$0")/backend"

# Add Go to PATH
export PATH=$HOME/go/bin:$PATH

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go first."
    echo "Visit: https://golang.org/doc/install"
    exit 1
fi

# Check if .env file exists
if [ ! -f ".env" ]; then
    echo "Error: .env file not found in backend directory"
    echo "Please create a .env file with your configuration"
    exit 1
fi

echo "Starting backend server..."
echo "Logs are available at: /tmp/backend.log"
nohup go run cmd/api/main.go > /tmp/backend.log 2>&1 &

sleep 2
if ps aux | grep "go run cmd/api/main.go" | grep -v grep > /dev/null; then
    echo "✅ Backend server started successfully on port 8081"
    echo "View logs: tail -f /tmp/backend.log"
else
    echo "❌ Failed to start backend server"
    echo "Check logs: cat /tmp/backend.log"
    exit 1
fi
