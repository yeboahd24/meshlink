@echo off
echo Rebuilding web broadcaster with P2P integration...

docker compose -f deployments/docker-compose.dev.yml exec dev-env go build -buildvcs=false -o dist/broadcaster-web-fixed ./cmd/broadcaster-web

echo.
echo ✅ Fixed web broadcaster built!
echo.
echo Now run: dist/broadcaster-web-fixed
echo Then visit: http://localhost:8080
echo.
echo You should now see:
echo - Real camera and microphone detected
echo - P2P streaming active
echo - Web frames being sent to browser
echo - Live video display updating every 200ms