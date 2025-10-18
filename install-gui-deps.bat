@echo off
echo 🔧 Installing GUI Dependencies for Windows...

echo.
echo This script helps install dependencies for GUI builds.
echo You need these for video windows to work:

echo.
echo 1. TDM-GCC (C compiler for CGO):
echo    Download: https://jmeubank.github.io/tdm-gcc/
echo    Install and add to PATH

echo.
echo 2. pkg-config (for library linking):
echo    Option A: choco install pkgconfiglite
echo    Option B: Download from http://ftp.gnome.org/pub/gnome/binaries/win32/dependencies/

echo.
echo 3. FFmpeg (already installed if Docker tests worked):
echo    Download: https://ffmpeg.org/download.html
echo    Add to PATH

echo.
echo After installing dependencies, run:
echo build-gui-binaries.bat

echo.
echo Alternative: Use headless versions (no GUI dependencies needed):
echo build-binaries.bat

pause