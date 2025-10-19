#!/bin/bash
echo "🔨 Building MeshLink on WSL..."

# Install GUI dependencies for Linux
echo "Installing dependencies..."
sudo apt-get update
sudo apt-get install -y pkg-config libgl1-mesa-dev libxrandr-dev libxcursor-dev libxinerama-dev libxi-dev libglfw3-dev libxxf86vm-dev libx11-dev libxext-dev v4l-utils

echo ""
echo "Building broadcaster..."
go build -o meshlink-broadcaster cmd/broadcaster/main.go

echo ""
echo "Building viewer..."
go build -o meshlink-viewer cmd/viewer/main.go

echo ""
echo "Building headless versions..."
go build -o meshlink-broadcaster-headless cmd/broadcaster-headless/main.go
go build -o meshlink-viewer-headless cmd/viewer-headless/main.go

echo ""
echo "✅ Build complete!"
echo ""
echo "📹 Check your camera:"
echo "v4l2-ctl --list-devices"
echo ""
echo "🚀 To test:"
echo "Terminal 1: ./meshlink-broadcaster-headless"
echo "Terminal 2: ./meshlink-viewer-headless"
