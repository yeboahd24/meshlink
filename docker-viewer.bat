@echo off
echo Building MeshLink Viewer Docker Image...
docker build -f Dockerfile.viewer -t meshlink-viewer .

echo.
echo Running MeshLink Viewer...
docker run -it --rm --network host meshlink-viewer