# Windows GUI Build Instructions

## Prerequisites

1. **Install Go**: Download from https://golang.org/dl/
2. **Install FFmpeg**: Download from https://ffmpeg.org/download.html
   - Add FFmpeg to your PATH
3. **Install pkg-config**: 
   - Download from http://ftp.gnome.org/pub/gnome/binaries/win32/dependencies/
   - Or use: `choco install pkgconfiglite` (if you have Chocolatey)
4. **Install TDM-GCC**: For CGO compilation
   - Download from https://jmeubank.github.io/tdm-gcc/

## Build Commands

### GUI Broadcaster (with camera preview)
```bash
go build -o broadcaster-gui.exe cmd/broadcaster/main.go
```

### GUI Viewer (with video display)
```bash
go build -o viewer-gui.exe cmd/viewer/main.go
```

### Headless versions (no GUI dependencies)
```bash
go build -o broadcaster-headless.exe cmd/broadcaster-headless/main.go
go build -o viewer-headless.exe cmd/viewer-headless/main.go
```

## Run GUI Applications

```bash
# Start broadcaster with GUI
./broadcaster-gui.exe

# Start viewer with GUI  
./viewer-gui.exe
```

## Troubleshooting

If you get OpenGL/GUI errors, use the headless versions:
```bash
./broadcaster-headless.exe
./viewer-headless.exe
```