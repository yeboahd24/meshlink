@echo off
echo Building MeshLink Simple Docker Image...
docker build -f Dockerfile.simple -t meshlink-simple .

echo.
echo Running MeshLink Simple Broadcaster...
docker run -it --rm -p 8080:8080 meshlink-simple