#!/bin/bash

# Script to sync the corrected main.go to Oracle Cloud server
# Usage: ./sync-to-oracle.sh <oracle-host-or-ip>

if [ -z "$1" ]; then
    echo "Error: Please provide Oracle Cloud server hostname or IP"
    echo "Usage: ./sync-to-oracle.sh <oracle-host-or-ip>"
    echo "Example: ./sync-to-oracle.sh instance-20251107-1609"
    echo "     or: ./sync-to-oracle.sh 123.45.67.89"
    exit 1
fi

ORACLE_HOST=$1
LOCAL_FILE="/var/www/html/chatbot/backend/cmd/api/main.go"
REMOTE_PATH="~/saas/backend/cmd/api/main.go"

echo "Copying main.go to Oracle Cloud server..."
scp "$LOCAL_FILE" ubuntu@"$ORACLE_HOST":"$REMOTE_PATH"

if [ $? -eq 0 ]; then
    echo "✓ File copied successfully!"
    echo ""
    echo "Now SSH into your server and run:"
    echo "  cd ~/saas/backend"
    echo "  SERVER_PORT=8081 go run cmd/api/main.go"
else
    echo "✗ Failed to copy file. Please check:"
    echo "  - Your Oracle Cloud hostname/IP is correct"
    echo "  - You have SSH access configured"
    echo "  - The remote directory exists"
fi


