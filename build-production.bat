@echo off
echo Building MeshLink for Production Distribution...
echo.

echo 1️⃣ Building Cross-Platform Versions (Headless)...
docker compose -f deployments/docker-compose.dev.yml exec dev-env make build-all

echo.
echo 2️⃣ Building GUI Versions (Docker)...
docker compose -f deployments/docker-compose.dev.yml exec dev-env make docker-build-gui

echo.
echo ✅ Production Build Complete!
echo.
echo 📦 Available Binaries:
echo.
echo 🖥️  GUI Versions (Full Interface):
echo    - broadcaster-gui
echo    - viewer-gui
echo.
echo 📱 Cross-Platform Versions:
echo    - broadcaster-windows-amd64.exe
echo    - viewer-windows-amd64.exe
echo    - broadcaster-darwin-amd64 (macOS)
echo    - viewer-darwin-amd64 (macOS)  
echo    - broadcaster-linux-amd64
echo    - viewer-linux-amd64
echo    - broadcaster-linux-arm7 (Raspberry Pi)
echo    - viewer-linux-arm7 (Raspberry Pi)
echo.
echo 🎯 Ready for Distribution!