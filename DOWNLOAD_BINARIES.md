# 📦 Download Your MeshLink Binaries

## Step 1: Go to GitHub Actions

1. **Open your repository**: https://github.com/yeboahd24/meshlink
2. **Click "Actions" tab** (next to Code, Issues, Pull requests)
3. **Find the latest successful build** (green checkmark ✅)
4. **Click on the build** to open details

## Step 2: Download Artifacts

### Windows Binaries:
- **Scroll down** to "Artifacts" section
- **Click "meshlink-windows-binaries"** to download ZIP
- **Extract** to get:
  - `meshlink-broadcaster-gui.exe` (with video windows)
  - `meshlink-viewer-gui.exe` (with video windows)
  - `meshlink-broadcaster-headless.exe` (terminal version)
  - `meshlink-viewer-headless.exe` (terminal version)

### Linux Binaries:
- **Click "meshlink-linux-binaries"** to download ZIP
- **Extract** to get:
  - `meshlink-broadcaster-linux` (GUI version)
  - `meshlink-viewer-linux` (GUI version)
  - `meshlink-broadcaster-linux-headless` (terminal version)
  - `meshlink-viewer-linux-headless` (terminal version)

## Step 3: Test Your Binaries

### Windows Test:
```cmd
# Terminal 1
meshlink-broadcaster-gui.exe

# Terminal 2  
meshlink-viewer-gui.exe
```

### Expected Result:
- **Broadcaster**: Shows camera preview window + "Viewers: 1"
- **Viewer**: Shows live video stream window
- **P2P Connection**: Automatic discovery and streaming

## Step 4: Distribute to Churches

### Ready-to-Share Files:
- ✅ **meshlink-broadcaster-gui.exe** - Church staff runs this
- ✅ **meshlink-viewer-gui.exe** - Congregation runs this
- ✅ **No installation needed** - Just copy and double-click!

## Alternative: Create Release

### For easier distribution:
```bash
git tag v1.0.0
git push origin v1.0.0
```

This creates a GitHub Release with downloadable ZIP packages automatically!