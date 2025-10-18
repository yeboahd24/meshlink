@echo off
echo Building final MeshLink distribution with web broadcasting...

docker compose -f deployments/docker-compose.dev.yml exec dev-env make build-all

echo.
echo ✅ Final Build Complete!
echo.
echo 📦 Distribution Files:
echo.
echo 🌐 Web Broadcasting (RECOMMENDED):
echo    - broadcaster-web-windows.exe (Windows with web interface)
echo    - broadcaster-web-linux (Linux with web interface)
echo.
echo 📱 Cross-Platform Headless:
echo    - broadcaster-windows-amd64.exe (Windows command-line)
echo    - viewer-windows-amd64.exe (Windows command-line)
echo    - broadcaster-darwin-amd64 (macOS command-line)
echo    - viewer-darwin-amd64 (macOS command-line)
echo    - broadcaster-linux-amd64 (Linux command-line)
echo    - viewer-linux-amd64 (Linux command-line)
echo.
echo 🎯 For Churches - Use Web Broadcasting:
echo    1. Run: broadcaster-web-windows.exe
echo    2. Share URL: http://[YOUR-IP]:8080
echo    3. Congregation opens browsers to that URL
echo    4. Everyone sees live video stream!
echo.
echo 🚀 MeshLink is ready for deployment!