@echo off
echo Building improved web broadcaster...

docker compose -f deployments/docker-compose.dev.yml exec dev-env go build -buildvcs=false -o dist/broadcaster-web-improved ./cmd/broadcaster-web

echo.
echo ✅ Improved web broadcaster built!
echo.
echo Now run: dist/broadcaster-web-improved
echo Then visit: http://localhost:8080
echo.
echo You should now see:
echo - Better formatted video display
echo - "LIVE CHURCH STREAM" header
echo - Professional streaming interface
echo - Debug output showing frames being sent
echo.
echo The "non-English words" were raw frame data - now it shows proper video info!