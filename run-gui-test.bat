@echo off
echo Starting MeshLink GUI Test...

echo.
echo Starting Broadcaster GUI in new window...
start "MeshLink Broadcaster" broadcaster-gui.exe

echo.
echo Waiting 3 seconds for broadcaster to start...
timeout /t 3 /nobreak > nul

echo.
echo Starting Viewer GUI...
viewer-gui.exe