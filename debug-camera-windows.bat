@echo off
echo 🔍 Debugging Windows Camera Detection...

echo.
echo 1. Testing FFmpeg DirectShow device list:
ffmpeg -list_devices true -f dshow -i dummy

echo.
echo 2. Checking if camera is accessible:
ffmpeg -f dshow -i video="Integrated Camera" -frames:v 1 -f null - 2>&1 | findstr /i "video error input"

echo.
echo 3. Alternative camera names to try:
echo - Integrated Camera
echo - USB Camera
echo - HD Webcam
echo - Laptop Camera

echo.
echo 4. Testing with default video device:
timeout /t 2 /nobreak > nul
ffmpeg -f dshow -list_options true -i video=0 2>&1 | findstr /i "video"

echo.
echo ✅ Debug complete!
echo.
echo Copy the camera name from above and update the code.
pause