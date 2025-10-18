@echo off
echo Testing web broadcaster...

echo Building web broadcaster...
docker compose -f deployments/docker-compose.dev.yml exec dev-env go build -buildvcs=false -o dist/broadcaster-web-test ./cmd/broadcaster-web

echo.
echo Web broadcaster built! Now run:
echo   dist/broadcaster-web-test
echo.
echo Then visit: http://localhost:8080
echo.
echo If it works, you should see:
echo - Web server startup messages
echo - HTML page in browser
echo - Websocket connection