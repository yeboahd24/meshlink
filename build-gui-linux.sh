#!/bin/bash
echo "Building MeshLink GUI for Linux (WSL alternative)..."

echo ""
echo "Installing dependencies..."
sudo apt-get update
sudo apt-get install -y pkg-config libgl1-mesa-dev libxrandr-dev libxcursor-dev libxinerama-dev libxi-dev libglfw3-dev

echo ""
echo "Building GUI applications..."
GOOS=linux GOARCH=amd64 go build -o meshlink-broadcaster-linux cmd/broadcaster/main.go
GOOS=linux GOARCH=amd64 go build -o meshlink-viewer-linux cmd/viewer/main.go

echo ""
echo "Building Windows executables (cross-compile)..."
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o meshlink-broadcaster-windows.exe cmd/broadcaster-headless/main.go
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o meshlink-viewer-windows.exe cmd/viewer-headless/main.go

echo ""
echo "Built files:"
echo "- meshlink-broadcaster-linux (GUI for Linux/WSL)"
echo "- meshlink-viewer-linux (GUI for Linux/WSL)" 
echo "- meshlink-broadcaster-windows.exe (headless for Windows)"
echo "- meshlink-viewer-windows.exe (headless for Windows)"

echo ""
echo "For Windows GUI, you need to build on actual Windows machine."