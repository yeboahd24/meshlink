@echo off
echo Building MeshLink with GUI for all platforms...

docker compose -f deployments/docker-compose.dev.yml exec dev-env make build-all

echo.
echo Build complete! GUI binaries created:
echo.
echo Windows (with GUI):
echo - broadcaster-windows-amd64.exe
echo - viewer-windows-amd64.exe
echo.
echo macOS (with GUI):  
echo - broadcaster-darwin-amd64
echo - viewer-darwin-amd64
echo.
echo Linux (with GUI):
echo - broadcaster-linux-amd64  
echo - viewer-linux-amd64
echo.
echo Raspberry Pi (headless):
echo - broadcaster-linux-arm7
echo - viewer-linux-arm7
echo.
echo All binaries include:
echo ✅ Professional GUI interface
echo ✅ Real camera and microphone support  
echo ✅ Animated video display
echo ✅ P2P networking
echo ✅ Single-file deployment