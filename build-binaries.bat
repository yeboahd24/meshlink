@echo off
echo 🔨 Building MeshLink Binaries...

echo.
echo Building headless versions (no GUI dependencies)...
go build -ldflags "-s -w" -o meshlink-broadcaster.exe cmd/broadcaster-headless/main.go
go build -ldflags "-s -w" -o meshlink-viewer.exe cmd/viewer-headless/main.go

echo.
echo Building simple versions...
go build -ldflags "-s -w" -o meshlink-broadcaster-simple.exe cmd/broadcaster-simple/main.go
go build -ldflags "-s -w" -o meshlink-viewer-simple.exe cmd/viewer-simple/main.go

echo.
echo ✅ Binaries built successfully:
echo - meshlink-broadcaster.exe (ready to run)
echo - meshlink-viewer.exe (ready to run)
echo - meshlink-broadcaster-simple.exe (Docker version)
echo - meshlink-viewer-simple.exe (Docker version)

echo.
echo 🚀 To test streaming:
echo 1. Run: meshlink-broadcaster.exe
echo 2. Run: meshlink-viewer.exe (in another terminal)
echo.
pause