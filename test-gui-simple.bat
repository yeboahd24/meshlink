@echo off
echo 🎬 Simple GUI Test - MeshLink Live Streaming

echo.
echo Step 1: Building applications...
go build -o test-broadcaster.exe cmd/broadcaster-headless/main.go
go build -o test-viewer.exe cmd/viewer-headless/main.go

echo.
echo Step 2: Starting broadcaster...
echo 📡 Broadcaster will stream FFmpeg test pattern
start "Broadcaster" cmd /k "test-broadcaster.exe"

echo.
echo Step 3: Waiting for broadcaster to initialize...
timeout /t 3 /nobreak > nul

echo.
echo Step 4: Starting viewer...
echo 📺 Viewer will receive and display stream data
echo 👀 Watch the terminal for "Received H.264 frame" messages
test-viewer.exe

echo.
echo Test completed! Check if viewer received frames from broadcaster.
pause