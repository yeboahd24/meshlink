# Cross-Compilation Limitations

## ❌ What WSL CANNOT Build:

### Windows GUI Applications:
- **Fyne UI** requires Windows graphics libraries (OpenGL, DirectX)
- **CGO dependencies** need Windows C compilers
- **GUI frameworks** don't cross-compile with CGO

### Why GUI Cross-Compilation Fails:
```bash
# This WON'T work from WSL:
GOOS=windows GOARCH=amd64 go build cmd/broadcaster/main.go
# Error: CGO not supported when cross compiling
```

## ✅ What WSL CAN Build:

### Headless Applications:
```bash
# This WORKS from WSL:
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o broadcaster-headless.exe cmd/broadcaster-headless/main.go
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o viewer-headless.exe cmd/viewer-headless/main.go
```

### Docker Images:
```bash
# This WORKS from WSL:
docker build -f Dockerfile.simple -t meshlink .
```

## Solutions for Distribution:

### Option 1: Build Headless in WSL (Shareable)
```bash
# In WSL - creates Windows .exe files
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o meshlink-broadcaster-windows.exe cmd/broadcaster-headless/main.go
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o meshlink-viewer-windows.exe cmd/viewer-headless/main.go
```

### Option 2: Use GitHub Actions (Automated)
- Build GUI on Windows runners
- Build headless on Linux runners
- Automatic releases with binaries

### Option 3: Docker Distribution
- Share Docker images instead of binaries
- Works on any platform with Docker
- No compilation needed by end users

### Option 4: Windows VM/Machine
- Build GUI binaries on actual Windows
- Full GUI support with video windows

## Recommendation:
**Build headless versions in WSL for now** - they prove P2P streaming works and are fully shareable. Add GUI later when you have Windows build environment.