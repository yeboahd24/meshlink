@echo off
echo Building MeshLink Windows GUI Applications...

echo.
echo Building Broadcaster GUI...
go build -o broadcaster-gui.exe cmd/broadcaster/main.go
if %errorlevel% neq 0 (
    echo Failed to build broadcaster GUI
    pause
    exit /b 1
)

echo.
echo Building Viewer GUI...
go build -o viewer-gui.exe cmd/viewer/main.go
if %errorlevel% neq 0 (
    echo Failed to build viewer GUI
    pause
    exit /b 1
)

echo.
echo Building Headless versions as backup...
go build -o broadcaster-headless.exe cmd/broadcaster-headless/main.go
go build -o viewer-headless.exe cmd/viewer-headless/main.go

echo.
echo Build complete! Files created:
echo - broadcaster-gui.exe (with camera preview)
echo - viewer-gui.exe (with video display)
echo - broadcaster-headless.exe (no GUI)
echo - viewer-headless.exe (no GUI)

echo.
echo To test:
echo 1. Run broadcaster-gui.exe
echo 2. Run viewer-gui.exe in another terminal
echo.
pause