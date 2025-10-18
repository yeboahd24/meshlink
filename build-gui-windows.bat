@echo off
echo Building GUI versions for Windows distribution...

docker compose -f deployments/docker-compose.dev.yml exec dev-env make build-gui

echo Copying GUI binaries for Windows distribution...
copy dist\broadcaster-gui dist\broadcaster-gui-windows.exe
copy dist\viewer-gui dist\viewer-gui-windows.exe

echo Done! Windows GUI binaries ready:
echo - broadcaster-gui-windows.exe
echo - viewer-gui-windows.exe

echo.
echo These have full GUI with:
echo - Camera preview and controls
echo - Video display with animations  
echo - Real camera and microphone support
echo - Professional interface