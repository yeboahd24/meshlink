@echo off
echo Building MeshLink for all platforms...

go mod download
go mod tidy

mkdir dist 2>nul

echo Building Windows...
set GOOS=windows
set GOARCH=amd64
go build -buildvcs=false -o dist\broadcaster-windows-amd64.exe .\cmd\broadcaster-headless
go build -buildvcs=false -o dist\viewer-windows-amd64.exe .\cmd\viewer-headless

echo Building macOS...
set GOOS=darwin
set GOARCH=amd64
go build -buildvcs=false -o dist\broadcaster-darwin-amd64 .\cmd\broadcaster-headless
go build -buildvcs=false -o dist\viewer-darwin-amd64 .\cmd\viewer-headless

echo Building Linux...
set GOOS=linux
set GOARCH=amd64
go build -buildvcs=false -o dist\broadcaster-linux-amd64 .\cmd\broadcaster-headless
go build -buildvcs=false -o dist\viewer-linux-amd64 .\cmd\viewer-headless

echo Building ARM (Raspberry Pi)...
set GOOS=linux
set GOARCH=arm
set GOARM=7
go build -buildvcs=false -o dist\broadcaster-linux-arm7 .\cmd\broadcaster-headless
go build -buildvcs=false -o dist\viewer-linux-arm7 .\cmd\viewer-headless

echo Build complete! Binaries are in dist\ folder
dir dist