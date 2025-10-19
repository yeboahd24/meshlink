# WSL Testing Guide

## ⚠️ WSL Camera Limitation

**WSL cannot access Windows cameras directly** - this is a WSL limitation, not a MeshLink issue.

### What This Means:
- ❌ `/dev/video0` doesn't exist in WSL
- ❌ WSL can't access Windows hardware cameras
- ✅ **MeshLink will use FFmpeg test pattern automatically**
- ✅ **P2P streaming still works perfectly**

## 🎯 Testing Strategy

### Option 1: Test P2P with Test Pattern (WSL)
```bash
# Build
chmod +x build-test-wsl.sh
./build-test-wsl.sh

# Test streaming (uses test pattern)
# Terminal 1
./meshlink-broadcaster-headless

# Terminal 2
./meshlink-viewer-headless
```

**Expected**: Colorful test pattern streams via P2P ✅

### Option 2: Test Real Camera (Windows)
```bash
# Use GitHub Actions binaries (already built)
# Download from Actions artifacts
# Run on Windows to test real camera
```

**Expected**: Real camera video streams via P2P ✅

## 🧪 What WSL Testing Proves:

### ✅ Works in WSL:
- P2P network discovery
- FFmpeg encoding (test pattern)
- MJPEG compression
- Real-time streaming
- Multi-viewer support
- Frame statistics

### ❌ Doesn't Work in WSL:
- Real camera access (WSL limitation)
- GUI windows (no X server by default)

## 🚀 Quick Test:

```bash
# Build headless versions
go build -o broadcaster cmd/broadcaster-headless/main.go
go build -o viewer cmd/viewer-headless/main.go

# Terminal 1: Start broadcaster
./broadcaster

# Terminal 2: Start viewer
./viewer
```

## 📊 Expected Output:

### Broadcaster:
```
Camera: No camera found - using test pattern
Using FFmpeg test pattern
Streaming: true | Viewers: 1 | Frames: 150
```

### Viewer:
```
Received H.264 frame: 51234 bytes
Frames: 45 | Viewing: true | Bytes: 2301440
```

## ✅ Success Criteria:

Even without camera, if you see:
- ✅ **Frames increasing** on both sides
- ✅ **Viewer count: 1** on broadcaster
- ✅ **Bytes transferring** continuously
- ✅ **30 FPS** frame rate

**Then MeshLink P2P streaming is working perfectly!**

## 🎬 For Real Camera Testing:

Use the **Windows GUI binaries** from GitHub Actions:
1. Download artifacts
2. Run on Windows
3. Real camera will be detected automatically
4. Same P2P streaming, real video instead of test pattern

**Bottom Line**: WSL is perfect for testing P2P functionality. Real camera testing needs Windows/Linux with actual hardware.