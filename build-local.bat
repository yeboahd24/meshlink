@echo off
echo Building MeshLink locally (headless versions)...

go mod download
go mod tidy

mkdir build 2>nul

echo Building headless versions...
go build -buildvcs=false -o build\broadcaster.exe .\cmd\broadcaster-headless
go build -buildvcs=false -o build\viewer.exe .\cmd\viewer-headless

echo Build complete! Binaries are in build\ folder
dir build