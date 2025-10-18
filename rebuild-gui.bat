@echo off
echo Rebuilding GUI with enhanced video display...

docker compose -f deployments/docker-compose.dev.yml exec dev-env make build-gui

echo Done! Test the enhanced video display:
echo 1. Run: ./dist/broadcaster-gui
echo 2. Run: ./dist/viewer-gui  
echo 3. Click "Connect to Stream" in viewer
echo 4. Click "Start Broadcasting" in broadcaster
echo 5. Watch animated video display with live stats!