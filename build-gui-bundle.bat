@echo off
echo Creating MeshLink GUI Bundle for Distribution...

echo.
echo Building GUI applications...
go build -o meshlink-broadcaster-gui.exe cmd/broadcaster/main.go
go build -o meshlink-viewer-gui.exe cmd/viewer/main.go

echo.
echo Creating distribution bundle...
mkdir MeshLink-Bundle 2>nul
copy meshlink-broadcaster-gui.exe MeshLink-Bundle\
copy meshlink-viewer-gui.exe MeshLink-Bundle\

echo.
echo Copying required DLLs (if available)...
copy C:\Windows\System32\opengl32.dll MeshLink-Bundle\ 2>nul
copy C:\Windows\System32\glu32.dll MeshLink-Bundle\ 2>nul

echo.
echo Creating user guide...
echo MeshLink - Church Video Streaming > MeshLink-Bundle\README.txt
echo. >> MeshLink-Bundle\README.txt
echo BROADCASTER (Church Staff): >> MeshLink-Bundle\README.txt
echo 1. Run meshlink-broadcaster-gui.exe >> MeshLink-Bundle\README.txt
echo 2. Click "Start Broadcasting" >> MeshLink-Bundle\README.txt
echo 3. Camera preview will show >> MeshLink-Bundle\README.txt
echo. >> MeshLink-Bundle\README.txt
echo VIEWER (Congregation): >> MeshLink-Bundle\README.txt
echo 1. Connect to church WiFi >> MeshLink-Bundle\README.txt
echo 2. Run meshlink-viewer-gui.exe >> MeshLink-Bundle\README.txt
echo 3. Click "Connect" to view stream >> MeshLink-Bundle\README.txt
echo 4. Video will appear in window >> MeshLink-Bundle\README.txt

echo.
echo Bundle created in MeshLink-Bundle\ folder
echo Ready for distribution - includes GUI with video display!
pause