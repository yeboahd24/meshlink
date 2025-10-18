# GUI Distribution Guide

## Option 1: Static Linking (Recommended)
```bash
build-gui-portable.bat
```
- Creates fully self-contained GUI binaries
- No DLL dependencies needed
- Larger file size (~20-50MB each)
- Works on any Windows machine

## Option 2: Bundle with DLLs
```bash
build-gui-bundle.bat
```
- Creates folder with GUI apps + required DLLs
- Smaller executable size
- Users copy entire folder
- More compatible across different Windows versions

## What Users Get:

### Broadcaster GUI Features:
- ✅ **Camera preview window** - see what you're streaming
- ✅ **Start/Stop button** - easy control
- ✅ **Quality selector** - 480p/720p/1080p
- ✅ **Viewer counter** - see who's watching
- ✅ **Status indicators** - streaming/stopped

### Viewer GUI Features:
- ✅ **Video display window** - watch the stream
- ✅ **Connect/Disconnect button** - easy control
- ✅ **Auto-discovery** - finds broadcaster automatically
- ✅ **Quality indicator** - shows stream quality
- ✅ **Connection status** - connected/searching

## Distribution Methods:

### Method 1: ZIP Download
```
meshlink-gui-v1.0.zip
├── meshlink-broadcaster-gui.exe
├── meshlink-viewer-gui.exe
└── README.txt
```

### Method 2: Installer (Future)
- Create NSIS installer
- One-click installation
- Desktop shortcuts
- Start menu entries

### Method 3: Portable Folder
```
MeshLink-Portable/
├── meshlink-broadcaster-gui.exe
├── meshlink-viewer-gui.exe
├── opengl32.dll (if needed)
└── README.txt
```

## User Experience:
1. **Download** → Extract → Double-click
2. **Broadcaster** → Camera preview appears → Click "Start"
3. **Viewer** → Video window appears → Click "Connect"
4. **Streaming** → Real video display in window

Perfect for churches - visual interface with actual video display!