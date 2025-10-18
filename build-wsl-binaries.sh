#!/bin/bash
echo "🔨 Building Windows Binaries from WSL..."

echo ""
echo "Building headless Windows executables (shareable)..."
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w" -o meshlink-broadcaster-windows.exe cmd/broadcaster-headless/main.go
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w" -o meshlink-viewer-windows.exe cmd/viewer-headless/main.go

echo ""
echo "Building Linux binaries..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w" -o meshlink-broadcaster-linux cmd/broadcaster-headless/main.go
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w" -o meshlink-viewer-linux cmd/viewer-headless/main.go

echo ""
echo "✅ Cross-compiled binaries created:"
echo "Windows (shareable):"
echo "- meshlink-broadcaster-windows.exe"
echo "- meshlink-viewer-windows.exe"
echo ""
echo "Linux:"
echo "- meshlink-broadcaster-linux"
echo "- meshlink-viewer-linux"

echo ""
echo "🚀 These Windows .exe files can be shared with churches!"
echo "No installation needed - just copy and run."

echo ""
echo "Note: These are headless (terminal) versions."
echo "For GUI with video windows, build on actual Windows machine."