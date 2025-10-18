@echo off
echo Building Portable GUI Binaries with Dependencies...

echo.
echo Setting CGO environment for static linking...
set CGO_ENABLED=1
set CGO_LDFLAGS=-static -static-libgcc -static-libstdc++

echo.
echo Building GUI Broadcaster (with camera preview)...
go build -ldflags "-s -w -extldflags '-static'" -o meshlink-broadcaster-gui.exe cmd/broadcaster/main.go

echo.
echo Building GUI Viewer (with video display)...
go build -ldflags "-s -w -extldflags '-static'" -o meshlink-viewer-gui.exe cmd/viewer/main.go

echo.
echo Creating distribution folder...
mkdir dist 2>nul
copy meshlink-broadcaster-gui.exe dist\
copy meshlink-viewer-gui.exe dist\
copy USER_GUIDE.txt dist\

echo.
echo Portable GUI binaries created in dist\ folder:
echo - meshlink-broadcaster-gui.exe (with camera preview)
echo - meshlink-viewer-gui.exe (with video display window)
echo - USER_GUIDE.txt (instructions)

echo.
echo These can be distributed to users - GUI included!
pause