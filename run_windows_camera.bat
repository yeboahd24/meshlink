@echo off
echo 🎥 MeshLink Windows Camera Fix
echo =============================
echo.
echo 🔍 Testing camera access...
echo.

REM Test if FFmpeg can see DirectShow devices
ffmpeg -list_devices true -f dshow -i dummy 2>nul
if errorlevel 1 (
    echo ❌ FFmpeg DirectShow not available
    echo 💡 Installing Windows FFmpeg...
    goto :install_ffmpeg
) else (
    echo ✅ FFmpeg DirectShow available
)

echo.
echo 📹 Listing available cameras...
ffmpeg -list_devices true -f dshow -i dummy 2>&1 | findstr "video devices\|Integrated\|Camera\|USB\|Webcam"

echo.
echo 🚀 Starting MeshLink Broadcaster with Windows camera support...
echo.
echo If you see "Camera: Found [camera name]" - SUCCESS! 🎉
echo If you see "Camera: Using test pattern" - Camera access issue ⚠️
echo.

.\meshlink-broadcaster-windows.exe

pause
goto :eof

:install_ffmpeg
echo.
echo ⚠️  Windows FFmpeg not found or missing DirectShow support
echo.
echo 💡 To fix this:
echo 1. Download FFmpeg for Windows from: https://ffmpeg.org/download.html
echo 2. Extract to C:\ffmpeg\bin\
echo 3. Add C:\ffmpeg\bin to your PATH environment variable
echo 4. Restart Command Prompt and try again
echo.
pause