#!/bin/bash

echo "🚀 Starting Chatbot Test Server"
echo "================================"
echo ""

# Check if backend is running
if ! curl -s http://localhost:8081/widget.js > /dev/null; then
    echo "⚠️  Backend is not running on port 8081"
    echo "Starting backend..."
    /var/www/html/chatbot/start-backend.sh
    sleep 3
fi

echo "✅ Backend is running on port 8081"
echo ""

# Start Python HTTP server for testing
cd /var/www/html/chatbot

echo "🌐 Starting test web server..."
echo ""
echo "📝 Test page will be available at:"
echo "   → http://localhost:8000/test-chatbot.html"
echo "   → http://127.0.0.1:8000/test-chatbot.html"
echo ""
echo "Press Ctrl+C to stop the server"
echo ""

python3 -m http.server 8000

