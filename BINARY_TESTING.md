# Binary Testing Guide

## Step 1: Build Binaries
```bash
build-binaries.bat
```

**Creates:**
- `meshlink-broadcaster.exe` - Streams H.264 video
- `meshlink-viewer.exe` - Receives and displays stream

## Step 2: Test Streaming
```bash
test-binaries.bat
```

**What happens:**
1. Broadcaster starts in separate window
2. Viewer connects and shows received frames
3. You see real P2P streaming in action!

## Manual Testing:

### Terminal 1 (Broadcaster):
```bash
meshlink-broadcaster.exe
```
**Expected output:**
```
📡 Starting P2P streaming...
Node ID: 12D3KooW...
Starting broadcast...
📊 Frames: 150 | Viewers: 1 | Streaming: true
```

### Terminal 2 (Viewer):
```bash
meshlink-viewer.exe
```
**Expected output:**
```
📺 Discovering streams...
Node ID: 12D3KooW...
📥 Received H.264 frame: 51234 bytes
📊 Frames: 45 | Viewing: true | Bytes: 2301440
```

## Success Indicators:
- ✅ **Broadcaster**: Shows "Viewers: 1" (or more)
- ✅ **Viewer**: Shows "Received H.264 frame: XXXX bytes"
- ✅ **Network**: Both apps discover each other automatically
- ✅ **Streaming**: Frame count increases continuously

## Distribution Ready:
Once tested, you can share these `.exe` files with churches:
- No installation needed
- No dependencies required
- Just copy and run!

Perfect for church deployment! 🎉