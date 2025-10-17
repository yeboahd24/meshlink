@echo off
echo 🚀 Starting MeshLink Development Environment...
echo.

echo 📦 Building development environment...
docker compose -f deployments/docker-compose.dev.yml build dev-env

echo.
echo 🔧 Starting development container...
docker compose -f deployments/docker-compose.dev.yml up -d dev-env

echo.
echo ✅ Development environment ready!
echo.
echo Available commands:
echo   make docker-dev-shell     - Enter development container
echo   make docker-dev-build     - Build applications in container
echo   make docker-dev-test      - Run tests in container
echo   make docker-full-demo     - Start full demo environment
echo.
echo 🔗 Quick access:
echo   Development shell: make docker-dev-shell
echo   Full demo: make docker-full-demo
echo.

pause

echo 🐚 Opening development shell...
docker compose -f deployments/docker-compose.dev.yml exec dev-env bash