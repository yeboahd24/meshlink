@echo off
echo Testing web broadcaster build...

docker compose -f deployments/docker-compose.dev.yml exec dev-env make build-all

echo.
echo Build complete! Check for:
echo - broadcaster-web-windows.exe (Web-enabled broadcaster)
echo - broadcaster-windows-amd64.exe (Headless broadcaster)  
echo - viewer-windows-amd64.exe (Headless viewer)

dir dist