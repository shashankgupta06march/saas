#!/bin/bash

# Stop the Go backend server
echo "Stopping backend server..."

# Kill the process
pkill -f "go run cmd/api/main.go"

if [ $? -eq 0 ]; then
    echo "✅ Backend server stopped successfully"
else
    echo "⚠️ No running backend server found"
fi


