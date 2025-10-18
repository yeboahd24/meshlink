# GitHub Actions Build Guide

## 🚀 Automated Binary Building

### What This Does:
- **Windows Runner**: Builds GUI binaries with video windows
- **Linux Runner**: Builds Linux GUI binaries  
- **Automatic**: Triggers on every push/PR
- **Release**: Creates downloadable packages

### Setup Steps:

#### 1. Push to GitHub
```bash
git add .
git commit -m "Add GitHub Actions workflow"
git push origin main
```

#### 2. Check Actions Tab
- Go to your GitHub repository
- Click "Actions" tab
- Watch the build process

#### 3. Download Binaries
- Click on completed workflow
- Download "meshlink-windows-binaries" artifact
- Extract and test!

### What You Get:

#### Windows Binaries:
- `meshlink-broadcaster-gui.exe` - **With camera preview window**
- `meshlink-viewer-gui.exe` - **With video display window**
- `meshlink-broadcaster-headless.exe` - Terminal fallback
- `meshlink-viewer-headless.exe` - Terminal fallback

#### Linux Binaries:
- `meshlink-broadcaster-linux` - GUI version
- `meshlink-viewer-linux` - GUI version

### For Releases:

#### Create Release:
```bash
git tag v1.0.0
git push origin v1.0.0
```

#### GitHub automatically:
- Builds all binaries
- Creates ZIP packages
- Attaches to release
- Ready for church distribution!

### Benefits:
- ✅ **No Windows machine needed** - GitHub provides it
- ✅ **Professional builds** - All dependencies included
- ✅ **Automatic releases** - Just tag and push
- ✅ **Cross-platform** - Windows + Linux binaries
- ✅ **Church ready** - Download and distribute

Perfect solution for GUI binary distribution! 🎉