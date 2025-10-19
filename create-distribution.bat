@echo off
echo 📦 Creating MeshLink Distribution Package...

echo.
echo Step 1: Creating distribution folder...
mkdir MeshLink-Distribution 2>nul
cd MeshLink-Distribution

echo.
echo Step 2: Download your binaries from GitHub Actions
echo Place them here:
echo - meshlink-broadcaster-gui.exe
echo - meshlink-viewer-gui.exe

echo.
echo Step 3: Download FFmpeg
echo Visit: https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-essentials.zip
echo Extract ffmpeg.exe and place it here

echo.
echo Step 4: Create README
echo Creating README.txt...
(
echo MeshLink - Church Streaming System
echo ===================================
echo.
echo BROADCASTER ^(Church Staff^):
echo 1. Run meshlink-broadcaster-gui.exe
echo 2. Click "Start Broadcasting"
echo 3. Camera preview will show
echo.
echo VIEWER ^(Congregation^):
echo 1. Connect to church WiFi
echo 2. Run meshlink-viewer-gui.exe  
echo 3. Click "Connect"
echo 4. Video will appear
echo.
echo REQUIREMENTS:
echo - Windows 7 or later
echo - WiFi connection
echo - No internet needed
echo.
echo TROUBLESHOOTING:
echo - Keep all files in same folder
echo - Run as administrator if needed
echo - Check Windows camera permissions
echo.
echo Support: https://github.com/yeboahd24/meshlink
) > README.txt

echo.
echo Step 5: Create LICENSE file for FFmpeg
(
echo FFmpeg License
echo ==============
echo.
echo This distribution includes FFmpeg ^(https://ffmpeg.org^)
echo FFmpeg is licensed under LGPL v2.1 or later
echo.
echo Full license: https://ffmpeg.org/legal.html
echo Source code: https://github.com/FFmpeg/FFmpeg
) > FFMPEG-LICENSE.txt

echo.
echo ✅ Distribution folder created!
echo.
echo 📋 Checklist:
echo [ ] Download binaries from GitHub Actions
echo [ ] Download ffmpeg.exe from https://www.gyan.dev/ffmpeg/builds/
echo [ ] Place all files in MeshLink-Distribution folder
echo [ ] Create ZIP file
echo [ ] Test on clean Windows machine
echo.
echo Final structure should be:
echo MeshLink-Distribution/
echo ├── meshlink-broadcaster-gui.exe
echo ├── meshlink-viewer-gui.exe
echo ├── ffmpeg.exe
echo ├── README.txt
echo └── FFMPEG-LICENSE.txt
echo.
pause