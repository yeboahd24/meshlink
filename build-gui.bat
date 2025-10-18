@echo off
echo Building MeshLink GUI versions...

go mod download
go mod tidy

mkdir build 2>nul

echo Building GUI versions (requires display libraries)...
go build -buildvcs=false -o build\broadcaster-gui.exe .\cmd\broadcaster
go build -buildvcs=false -o build\viewer-gui.exe .\cmd\viewer

echo Build complete! GUI binaries are in build\ folder
dir build