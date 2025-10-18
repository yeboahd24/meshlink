@echo off
echo 🖥️ Testing MeshLink GUI Locally...

echo.
echo Building GUI applications...
go build -o broadcaster-gui.exe cmd/broadcaster/main.go
if %errorlevel% neq 0 (
    echo ❌ Failed to build broadcaster GUI
    echo Trying headless version...
    go build -o broadcaster-headless.exe cmd/broadcaster-headless/main.go
    set BROADCASTER=broadcaster-headless.exe
) else (
    set BROADCASTER=broadcaster-gui.exe
)

go build -o viewer-gui.exe cmd/viewer/main.go
if %errorlevel% neq 0 (
    echo ❌ Failed to build viewer GUI
    echo Trying headless version...
    go build -o viewer-headless.exe cmd/viewer-headless/main.go
    set VIEWER=viewer-headless.exe
) else (
    set VIEWER=viewer-gui.exe
)

echo.
echo Starting broadcaster in new window...
start "MeshLink Broadcaster" %BROADCASTER%

echo.
echo Waiting 5 seconds for broadcaster to start...
timeout /t 5 /nobreak > nul

echo.
echo Starting viewer...
echo 📺 Look for video window to appear!
%VIEWER%