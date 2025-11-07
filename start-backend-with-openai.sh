#!/bin/bash

# Chatbot Backend Startup Script with OpenAI API Key
# Replace YOUR_OPENAI_KEY_HERE with your actual OpenAI API key

echo "=========================================="
echo "Starting Chatbot Backend Server"
echo "=========================================="
echo ""

# Check if OPENAI_API_KEY is provided as argument
if [ -n "$1" ]; then
    OPENAI_KEY="$1"
    echo "Using OpenAI API key from argument"
else
    # Prompt user for API key
    echo "Please enter your OpenAI API Key (starts with sk-):"
    read -s OPENAI_KEY
    echo ""
fi

# Validate key format
if [[ ! "$OPENAI_KEY" =~ ^sk- ]]; then
    echo "ERROR: Invalid API key format. OpenAI keys start with 'sk-'"
    exit 1
fi

echo "Starting backend on port 8081..."
echo ""

cd /var/www/html/chatbot/backend

export PATH=$HOME/go/bin:$PATH

SERVER_PORT=8081 \
DB_USER=admin \
DB_PASSWORD=Admin@123 \
DB_NAME=chatbot_saas \
DB_HOST=localhost \
DB_PORT=3306 \
JWT_SECRET=chatbot-saas-secret-key \
OPENAI_API_KEY="$OPENAI_KEY" \
ALLOWED_ORIGINS=http://localhost:3000 \
go run cmd/api/main.go

