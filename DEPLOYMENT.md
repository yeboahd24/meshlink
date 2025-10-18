# MeshLink Deployment Guide

## 🚀 Quick Deployment

### For Cross-Platform Distribution
```bash
make build-all
```
**Creates:** Headless versions for all platforms (Windows, macOS, Linux, ARM)

### For GUI Experience  
```bash
make build-gui
```
**Creates:** GUI versions for current platform (requires display libraries)

## 📦 Distribution Files

### Cross-Platform (Headless)
- `broadcaster-windows-amd64.exe` - Windows command-line
- `viewer-windows-amd64.exe` - Windows command-line  
- `broadcaster-darwin-amd64` - macOS command-line
- `viewer-darwin-amd64` - macOS command-line
- `broadcaster-linux-amd64` - Linux command-line
- `viewer-linux-amd64` - Linux command-line
- `broadcaster-linux-arm7` - Raspberry Pi
- `viewer-linux-arm7` - Raspberry Pi

### GUI Versions (Platform-Specific)
- `broadcaster-gui` - Full GUI interface
- `viewer-gui` - Full GUI with video display

## 🎯 Usage Instructions

### For End Users (Churches)

#### Windows Users:
1. **GUI Experience**: Use `broadcaster-gui` and `viewer-gui`
2. **Command-line**: Use `.exe` files

#### Other Platforms:
1. Use platform-specific headless versions
2. Run in terminal for full functionality

### Example Usage:
```bash
# Windows GUI (double-click or run)
./broadcaster-gui
./viewer-gui

# Windows Command-line
./broadcaster-windows-amd64.exe
./viewer-windows-amd64.exe

# macOS/Linux
./broadcaster-darwin-amd64
./viewer-darwin-amd64
```

## ✅ Features in All Versions

- **Real Camera & Microphone Support**
- **P2P Network Discovery** 
- **H.264 Streaming at 30 FPS**
- **Cross-Platform Compatibility**
- **Single-File Deployment**
- **No Internet Required**

## 🏗️ Build Requirements

### For Cross-Platform (Headless)
- Go 1.21+
- No additional dependencies

### For GUI Versions
- Go 1.21+
- Display libraries (X11 on Linux, native on Windows/macOS)
- OpenGL support

## 📋 Deployment Checklist

- [ ] Run `make build-all` for cross-platform versions
- [ ] Run `make build-gui` for GUI versions (optional)
- [ ] Test on target platforms
- [ ] Distribute appropriate binaries
- [ ] Provide usage instructions

## 🎉 Success!

Your MeshLink system is ready for deployment with:
- Complete P2P streaming functionality
- Real hardware integration
- Professional user experience
- Zero-configuration setup