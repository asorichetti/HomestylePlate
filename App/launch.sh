#!/bin/bash
# Homestyle Plate Launcher
# Double-click to start the meal randomizer

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "🍽️  Starting Homestyle Plate..."
echo "📁 Application directory: $SCRIPT_DIR"
echo ""
echo "Opening in browser..."

# Start the server in background
"$SCRIPT_DIR/HomestylePlate" &
SERVER_PID=$!

# Wait for server to be ready
sleep 2

# Open browser
open "http://localhost:8080"

echo ""
echo "✅ Homestyle Plate is running!"
echo "   Server PID: $SERVER_PID"
echo "   Open http://localhost:8080 in your browser"
echo ""
echo "To stop, press Ctrl+C or close this window"
echo "To quit the app, run: kill $SERVER_PID"

# Wait for the server process
wait $SERVER_PID
