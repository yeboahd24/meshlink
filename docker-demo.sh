#!/bin/bash

echo "Starting MeshLink Demo Environment..."
echo

echo "Building Docker images..."
docker compose -f deployments/docker-compose.yml build

echo
echo "Starting services..."
docker compose -f deployments/docker-compose.yml up -d

echo
echo "Demo environment is starting up..."
echo
echo "Services:"
echo "- Broadcaster: Running in background"
echo "- Viewer 1: Running in background"
echo "- Viewer 2: Running in background"
echo "- Demo UI: http://localhost:3000"
echo

sleep 10

echo "Opening demo dashboard..."
if command -v xdg-open > /dev/null; then
    xdg-open http://localhost:3000
elif command -v open > /dev/null; then
    open http://localhost:3000
else
    echo "Please open http://localhost:3000 in your browser"
fi

echo
echo "Demo is ready! Press Enter to stop all services..."
read

echo
echo "Stopping demo environment..."
docker compose -f deployments/docker-compose.yml down

echo "Demo stopped."