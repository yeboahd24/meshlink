@echo off
echo Building Portable MeshLink Binaries...

echo.
echo Building headless versions (no dependencies needed)...
go build -ldflags "-s -w" -o meshlink-broadcaster.exe cmd/broadcaster-headless/main.go
go build -ldflags "-s -w" -o meshlink-viewer.exe cmd/viewer-headless/main.go

echo.
echo Building simple versions (minimal dependencies)...
go build -ldflags "-s -w" -o meshlink-broadcaster-simple.exe cmd/broadcaster-simple/main.go
go build -ldflags "-s -w" -o meshlink-viewer-simple.exe cmd/viewer-simple/main.go

echo.
echo Portable binaries created:
echo - meshlink-broadcaster.exe (ready to share)
echo - meshlink-viewer.exe (ready to share)
echo - meshlink-broadcaster-simple.exe (Docker version)
echo - meshlink-viewer-simple.exe (Docker version)

echo.
echo These binaries can be shared with users - no installation needed!
echo Just copy the .exe files and run them.
pause