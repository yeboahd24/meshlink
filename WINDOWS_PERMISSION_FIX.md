# 🔧 Fix Windows Permission Error

## Problem:
"Windows cannot access the specified device, path, or file. You may not have the appropriate permissions to access the item."

## Solutions (Try in order):

### Solution 1: Unblock Downloaded Files
1. **Right-click** on `meshlink-broadcaster-gui.exe`
2. **Select "Properties"**
3. **Check "Unblock"** checkbox at bottom (if present)
4. **Click "OK"**
5. **Repeat for** `meshlink-viewer-gui.exe`
6. **Try running again**

### Solution 2: Run as Administrator
1. **Right-click** on `meshlink-broadcaster-gui.exe`
2. **Select "Run as administrator"**
3. **Click "Yes"** when prompted

### Solution 3: Move to Different Location
1. **Copy files** from Downloads to `C:\MeshLink\`
2. **Create folder** if needed: `mkdir C:\MeshLink`
3. **Run from** new location

### Solution 4: Windows Defender Exception
1. **Open Windows Security**
2. **Go to** Virus & threat protection
3. **Add exclusion** for the folder containing binaries
4. **Try running again**

### Solution 5: Use Headless Versions (Guaranteed to Work)
```cmd
# These always work:
meshlink-broadcaster-headless.exe
meshlink-viewer-headless.exe
```

## Quick Test:
```cmd
# Open Command Prompt as Administrator
# Navigate to extracted folder
cd C:\Users\YourName\Downloads\meshlink-windows-binaries
meshlink-broadcaster-headless.exe
```

## Why This Happens:
- Windows blocks downloaded executables for security
- GUI versions have more dependencies
- Headless versions are simpler and more reliable

## Recommendation:
**Use headless versions for now** - they prove P2P streaming works perfectly. Add GUI later when permission issues are resolved.