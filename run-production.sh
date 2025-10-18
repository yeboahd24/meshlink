#!/bin/bash
echo "🚀 Starting MeshLink Production with FFmpeg Streaming..."

echo ""
echo "Building and starting services..."
docker compose -f docker-compose-production.yml up --build

echo ""
echo "Services started:"
echo "📡 Broadcaster: http://localhost:8080 (FFmpeg H.264 streaming)"
echo "📺 Viewer 1: http://localhost:8081 (Receiving stream)"
echo "📺 Viewer 2: http://localhost:8082 (Receiving stream)"

echo ""
echo "Press Ctrl+C to stop all services"