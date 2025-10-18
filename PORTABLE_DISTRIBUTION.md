# Portable Distribution Guide

## ✅ YES - You CAN share just the binaries!

### Build Portable Executables:
```bash
build-portable.bat
```

### What Users Get:
- **meshlink-broadcaster.exe** - Single file, no installation needed
- **meshlink-viewer.exe** - Single file, no installation needed

### User Instructions (Zero Installation):

#### Church Staff (Broadcaster):
1. Copy `meshlink-broadcaster.exe` to any folder
2. Double-click to run
3. Streaming starts automatically
4. Share WiFi network name with congregation

#### Congregation (Viewers):
1. Copy `meshlink-viewer.exe` to any folder  
2. Connect to church WiFi
3. Double-click to run
4. Automatically finds and connects to stream

### Distribution Options:

#### Option 1: Direct Download
```
Website download:
├── meshlink-broadcaster.exe (5-10MB)
├── meshlink-viewer.exe (5-10MB)
└── Quick Start Guide.pdf
```

#### Option 2: USB Distribution
```
Church USB drives:
├── meshlink-broadcaster.exe
├── meshlink-viewer.exe
├── README.txt
└── "Just copy to your laptop and run!"
```

#### Option 3: Email/Cloud
```
Email attachment or cloud link:
- Small file size (under 20MB total)
- No installation required
- Works on any Windows computer
```

## Key Benefits:
- ✅ **No installation** - just copy and run
- ✅ **No dependencies** - everything bundled
- ✅ **Portable** - runs from any folder/USB
- ✅ **Small size** - under 10MB each
- ✅ **Zero IT support** - congregation can self-serve

## Technical Details:
- Headless versions avoid GUI dependency issues
- Built with `-ldflags "-s -w"` for smaller size
- All libraries statically linked
- Works on Windows 7+ (no additional runtime needed)