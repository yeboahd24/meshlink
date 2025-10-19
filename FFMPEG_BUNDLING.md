# FFmpeg Bundling for Distribution

## Problem:
Currently, users need FFmpeg installed separately. This is NOT user-friendly for churches.

## Solution: Bundle FFmpeg

### Option 1: Include FFmpeg Executable

#### Windows Distribution:
```
MeshLink-Windows/
├── meshlink-broadcaster.exe
├── meshlink-viewer.exe
├── ffmpeg.exe              ← Bundle this
├── ffprobe.exe             ← Bundle this
└── README.txt
```

#### Download FFmpeg:
- Windows: https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-essentials.zip
- Extract `ffmpeg.exe` and `ffprobe.exe`
- Include in your distribution ZIP

#### Update Code to Use Bundled FFmpeg:
```go
// In pkg/media/ffmpeg.go
func getFFmpegPath() string {
    // Check if ffmpeg is in same directory as executable
    exePath, _ := os.Executable()
    exeDir := filepath.Dir(exePath)
    bundledFFmpeg := filepath.Join(exeDir, "ffmpeg.exe")
    
    if _, err := os.Stat(bundledFFmpeg); err == nil {
        return bundledFFmpeg
    }
    
    // Fallback to system FFmpeg
    return "ffmpeg"
}
```

### Option 2: Static Linking (Advanced)

Use Go bindings that include FFmpeg statically:
- More complex build process
- Larger binary size (~50MB)
- No external dependencies

### Option 3: Installer Package

Create installer that includes FFmpeg:
- **NSIS** (Windows installer)
- **MSI** (Windows installer)
- **DMG** (macOS)
- **DEB/RPM** (Linux)

## Recommended Approach for Churches:

### Simple ZIP Distribution:
```
MeshLink-v1.0-Windows.zip
├── meshlink-broadcaster.exe
├── meshlink-viewer.exe
├── ffmpeg.exe              ← 70MB
├── ffprobe.exe             ← 70MB
├── README.txt
└── INSTALL.txt
```

**INSTALL.txt:**
```
MeshLink Church Streaming - Installation

1. Extract all files to a folder (e.g., C:\MeshLink)
2. Run meshlink-broadcaster.exe (church staff)
3. Run meshlink-viewer.exe (congregation)

No installation needed!
All files must stay in the same folder.
```

## File Sizes:
- meshlink-broadcaster.exe: ~10MB
- meshlink-viewer.exe: ~10MB
- ffmpeg.exe: ~70MB
- ffprobe.exe: ~70MB
- **Total**: ~160MB ZIP file

## Legal:
FFmpeg is LGPL/GPL licensed. You can bundle it freely, but:
- Include FFmpeg license file
- Mention FFmpeg in your credits
- Provide link to FFmpeg source

## Alternative: Lightweight Solution

Use **ffmpeg-static** npm package approach:
- Download minimal FFmpeg build
- Only include needed codecs
- Reduce size to ~20MB

## Implementation:

1. Download FFmpeg essentials build
2. Extract ffmpeg.exe and ffprobe.exe
3. Add to your distribution folder
4. Update code to check for bundled FFmpeg first
5. Create ZIP with all files
6. Distribute to churches

**Result**: Churches get single ZIP file, extract, and run. No installation, no dependencies!