@echo off
echo Starting MeshLink Demo Environment...
echo.

echo Building Docker images...
docker compose -f deployments/docker-compose.yml build

echo.
echo Starting services...
docker compose -f deployments/docker-compose.yml up -d

echo.
echo Demo environment is starting up...
echo.
echo Services:
echo - Broadcaster: Running in background
echo - Viewer 1: Running in background  
echo - Viewer 2: Running in background
echo - Demo UI: http://localhost:3000
echo.

timeout /t 10 /nobreak > nul

echo Opening demo dashboard...
start http://localhost:3000

echo.
echo Demo is ready! Press any key to stop all services...
pause > nul

echo.
echo Stopping demo environment...
docker compose -f deployments/docker-compose.yml down

echo Demo stopped.