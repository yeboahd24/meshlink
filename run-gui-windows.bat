@echo off
echo Starting MeshLink GUI on Windows...
echo.

echo Building GUI versions first...
docker compose -f deployments/docker-compose.dev.yml exec dev-env make docker-build-gui

echo.
echo GUI versions ready! You can now run:
echo.
echo For Church Staff (Broadcaster):
echo   dist\broadcaster-gui.exe
echo.
echo For Congregation (Viewers):  
echo   dist\viewer-gui.exe
echo.
echo Instructions:
echo 1. Make sure all devices are on same WiFi
echo 2. Run broadcaster first, then viewers
echo 3. Click "Start Broadcasting" in broadcaster
echo 4. Click "Connect to Stream" in viewers
echo.
pause