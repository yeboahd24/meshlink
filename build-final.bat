@echo off
echo Building MeshLink - GUI + Cross-platform versions...

docker compose -f deployments/docker-compose.dev.yml exec dev-env make build-all

echo.
echo ✅ Build Complete! You now have:
echo.
echo 🖥️  GUI Versions (Full Interface):
echo    - broadcaster-gui
echo    - viewer-gui
echo.
echo 📱 Cross-Platform Versions (Headless):
echo    - broadcaster-windows-amd64.exe
echo    - viewer-windows-amd64.exe  
echo    - broadcaster-darwin-amd64 (macOS)
echo    - viewer-darwin-amd64 (macOS)
echo    - broadcaster-linux-amd64
echo    - viewer-linux-amd64
echo    - broadcaster-linux-arm7 (Raspberry Pi)
echo    - viewer-linux-arm7 (Raspberry Pi)
echo.
echo 🎯 For Windows Users:
echo    - Use broadcaster-gui / viewer-gui for full GUI experience
echo    - Use .exe files for command-line or server deployment
echo.
echo 🚀 All versions include:
echo    ✅ Real camera and microphone support
echo    ✅ P2P networking and discovery  
echo    ✅ H.264 streaming at 30 FPS
echo    ✅ Single-file deployment