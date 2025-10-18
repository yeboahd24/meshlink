# GUI Testing Instructions

## Option 1: Windows Local GUI Test

### Prerequisites:
- Windows machine with Go installed
- FFmpeg installed and in PATH
- pkg-config and TDM-GCC (for GUI dependencies)

### Run GUI Test:
```bash
test-gui-local.bat
```

**What happens:**
1. Builds broadcaster and viewer GUI applications
2. Opens broadcaster in separate window with camera preview
3. Opens viewer with video display window
4. You should see live video streaming between windows

## Option 2: Simple Terminal Test (Easier)

### Run Simple Test:
```bash
test-gui-simple.bat
```

**What happens:**
1. Broadcaster terminal shows: "Streaming H.264 frames..."
2. Viewer terminal shows: "Received H.264 frame: XXXX bytes"
3. Proves P2P streaming is working

## Option 3: Docker + Local Hybrid

### Step 1: Run broadcaster in Docker
```bash
docker run -it --rm -p 8080:8080 --network host meshlink-simple
```

### Step 2: Run viewer locally
```bash
go run cmd/viewer/main.go
```

## Expected Results:

### ✅ Success Indicators:
- **Broadcaster**: Shows "Streaming: true, Viewers: 1+"
- **Viewer**: Shows "Received frame: XXXX bytes" messages
- **GUI Version**: Video preview window appears
- **Network**: Both apps discover each other automatically

### ❌ Troubleshooting:
- **No GUI**: Use headless versions (still proves streaming works)
- **No connection**: Check both apps are on same network
- **Build errors**: Missing dependencies (use Docker instead)

## Bottom Line:
Even without GUI, if viewer receives H.264 frames from broadcaster, your streaming is working perfectly!