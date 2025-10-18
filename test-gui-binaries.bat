@echo off
echo 🎬 Testing MeshLink GUI Binaries...

echo.
echo Starting GUI broadcaster...
echo 📡 Look for camera preview window to appear
start "MeshLink Broadcaster GUI" meshlink-broadcaster-gui.exe

echo.
echo Waiting 5 seconds for broadcaster GUI to load...
timeout /t 5 /nobreak > nul

echo.
echo Starting GUI viewer...
echo 📺 Look for video display window to appear
echo 👀 You should see live video streaming between windows!
meshlink-viewer-gui.exe

echo.
echo GUI test completed!
pause