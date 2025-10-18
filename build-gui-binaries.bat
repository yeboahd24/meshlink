@echo off
echo 🖥️ Building MeshLink GUI Binaries...

echo.
echo Building GUI versions (with video windows)...
go build -ldflags "-s -w" -o meshlink-broadcaster-gui.exe cmd/broadcaster/main.go
if %errorlevel% neq 0 (
    echo ❌ GUI broadcaster build failed - missing dependencies
    echo Building headless fallback...
    go build -ldflags "-s -w" -o meshlink-broadcaster-gui.exe cmd/broadcaster-headless/main.go
)

go build -ldflags "-s -w" -o meshlink-viewer-gui.exe cmd/viewer/main.go
if %errorlevel% neq 0 (
    echo ❌ GUI viewer build failed - missing dependencies
    echo Building headless fallback...
    go build -ldflags "-s -w" -o meshlink-viewer-gui.exe cmd/viewer-headless/main.go
)

echo.
echo ✅ GUI Binaries built:
echo - meshlink-broadcaster-gui.exe (with camera preview window)
echo - meshlink-viewer-gui.exe (with video display window)

echo.
echo 🚀 To test GUI streaming:
echo 1. Run: meshlink-broadcaster-gui.exe
echo 2. Run: meshlink-viewer-gui.exe (in another terminal)
echo 3. Look for video windows to appear!
echo.
pause