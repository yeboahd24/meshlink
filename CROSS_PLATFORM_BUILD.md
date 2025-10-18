# Cross-Platform Build Guide

## Current Situation: You're in WSL

### ❌ What WSL CAN'T Do:
- Build Windows GUI applications (OpenGL/graphics dependencies)
- Create .exe files with GUI support
- Access Windows graphics libraries

### ✅ What WSL CAN Do:
- Build headless/terminal applications
- Cross-compile basic Go programs
- Build Linux GUI applications

## Solutions:

### Option 1: Switch to Windows (Recommended for GUI)
```cmd
# Open Windows Command Prompt or PowerShell
cd C:\MeshLink
build-gui-portable.bat
```

### Option 2: Use WSL for Linux GUI + Windows Headless
```bash
# In WSL
chmod +x build-gui-linux.sh
./build-gui-linux.sh
```

### Option 3: GitHub Actions (Automated)
Create `.github/workflows/build.yml` to build on Windows runners automatically.

## Recommended Workflow:

### For Development (WSL):
```bash
# Test functionality
go run cmd/broadcaster-simple/main.go
go run cmd/viewer-simple/main.go
```

### For Distribution (Windows):
```cmd
# Switch to Windows
cd C:\MeshLink
build-gui-portable.bat
```

### For Church Deployment:
1. **Develop** in WSL/Linux (faster iteration)
2. **Build GUI** on Windows (for distribution)
3. **Distribute** Windows GUI binaries to churches

## Quick Test in WSL:
```bash
# Build and test headless versions
go build -o broadcaster cmd/broadcaster-headless/main.go
go build -o viewer cmd/viewer-headless/main.go

# Test P2P functionality
./broadcaster &
./viewer
```

**Bottom Line:** For GUI binaries that churches can use, you need to build on Windows.