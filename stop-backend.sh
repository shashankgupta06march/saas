#!/bin/bash

# Stop the Go backend server
echo "Stopping backend server..."

stopped=1

# 1. Kill the `go run` wrapper process.
if pkill -f "go run cmd/api/main.go"; then
    stopped=0
fi

# 2. `go run` compiles and launches a CHILD binary (e.g.
#    /tmp/go-build.../exe/main) that is the process actually listening on the
#    port. Killing only the wrapper leaves this child alive and holding the
#    port, which makes the next start fail with "address already in use".
if pkill -f "go-build.*/exe/main"; then
    stopped=0
fi

# 3. Belt-and-suspenders: kill whatever is still listening on the server port
#    so a restart can always bind.
PORT="${SERVER_PORT:-8081}"
PORT_PIDS=$(ss -ltnp 2>/dev/null | grep ":${PORT}\b" | grep -oP 'pid=\K[0-9]+' | sort -u)
if [ -n "$PORT_PIDS" ]; then
    for pid in $PORT_PIDS; do
        kill "$pid" 2>/dev/null && stopped=0
    done
fi

# Give processes a moment to release the port.
sleep 1

if [ "$stopped" -eq 0 ]; then
    echo "✅ Backend server stopped successfully"
else
    echo "⚠️ No running backend server found"
fi



