@echo off
echo 🎬 Testing MeshLink Binaries...

echo.
echo Starting broadcaster...
echo 📡 Broadcasting FFmpeg H.264 stream
start "MeshLink Broadcaster" cmd /k "meshlink-broadcaster.exe"

echo.
echo Waiting 5 seconds for broadcaster to start...
timeout /t 5 /nobreak > nul

echo.
echo Starting viewer...
echo 📺 Connecting to broadcaster stream
echo 👀 Watch for "Received H.264 frame" messages
meshlink-viewer.exe

echo.
echo Test completed!
pause