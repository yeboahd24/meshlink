#!/bin/bash
echo "🚀 MeshLink Quick Start..."

# Try modern docker compose first
if command -v docker &> /dev/null; then
    echo "Using: docker compose"
    docker compose -f docker-compose-production.yml up --build
else
    echo "❌ Docker not found. Please install Docker first."
fi