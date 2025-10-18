# MeshLink - Offline P2P Church Streaming

A decentralized, offline-first streaming solution for churches and communities.

## 🎯 Vision
Enable churches to stream services without internet dependency, using local WiFi networks and peer-to-peer technology.

## ✨ Key Features
- **Offline-First**: No internet required - works on local WiFi hotspots
- **P2P Architecture**: Direct peer-to-peer streaming using libp2p
- **Low Resource**: Runs on Raspberry Pi and mobile devices
- **Scalable**: 1 broadcaster → 50+ viewers
- **Cross-Platform**: Windows, macOS, Linux, Android, iOS

## 🏗️ Architecture
```
[WiFi Hotspot] → [Broadcaster] → [P2P Network] → [Viewers]
```

## 📡 How Data is Transmitted

### Data Flow Process
```
Camera → FFmpeg H.264 Encode → P2P Publish → P2P Receive → FFmpeg H.264 Decode → Display
```

### Step-by-Step Transmission
1. **Camera Capture**: FFmpeg captures from system camera
2. **H.264 Encoding**: FFmpeg compresses to ~50KB frames (720p quality)
3. **P2P Discovery**: mDNS finds peers on local WiFi network automatically
4. **Topic Publishing**: Broadcaster publishes frames to "meshlink/church/stream" topic
5. **Mesh Distribution**: libp2p distributes frames to all subscribed viewers
6. **Frame Reception**: Viewers receive encrypted frames via direct P2P connections
7. **H.264 Decoding**: FFmpeg decompresses frames back to video data
8. **Display**: Render 30 FPS video stream in real-time

### Network Architecture
```
┌─────────────┐    WiFi     ┌─────────────┐    WiFi     ┌─────────────┐
│ Broadcaster │◄──────────►│   Router    │◄──────────►│   Viewer    │
│             │             │             │             │             │
│ Publishes   │             │ Local Net   │             │ Subscribes  │
│ H.264 Data  │             │192.168.1.x  │             │ H.264 Data  │
└─────────────┘             └─────────────┘             └─────────────┘
```

### Technical Specifications
- **Protocol**: libp2p PubSub over TCP/QUIC
- **Discovery**: mDNS (zero configuration)
- **Encryption**: Built-in libp2p security
- **Bandwidth**: ~2 Mbps per stream
- **Latency**: <100ms on local network
- **Quality**: 720p @ 30 FPS with H.264 compression

### Key Advantages
- **No Internet Required**: Works on isolated WiFi networks
- **Direct P2P**: No central server or cloud dependency  
- **Automatic Discovery**: Viewers find broadcaster instantly
- **Encrypted**: All data transmission is secure
- **Efficient**: Multicast distribution saves bandwidth
- **Resilient**: Mesh network has no single point of failure

## 📱 Deployment Platforms

### Broadcaster Applications
- **Desktop**: Windows, macOS, Linux (Go + Fyne)
- **Raspberry Pi**: ARM builds for $35 hardware
- **Features**: Camera preview, quality controls, viewer statistics

### Viewer Applications
- **Mobile Apps**: iOS App Store, Google Play Store
- **Desktop**: Windows, macOS, Linux applications
- **Web App**: Progressive Web App (PWA) for browsers
- **Features**: Auto-discovery, HD video, touch controls, offline operation

### Real-World Usage
```
Church Staff (Broadcaster)
├── Desktop/Laptop App
├── Camera preview & controls
├── One-click start/stop
└── Live viewer count

Congregation (Viewers)
├── Mobile Apps (Primary)
│   ├── iOS App Store
│   ├── Google Play Store
│   └── Auto-discovers streams
├── Desktop Apps (Secondary)
└── Web Browser (Fallback)
```

## 🚀 Quick Start

### Development Mode
```bash
# Broadcaster (Church Setup)
go run cmd/broadcaster/main.go

# Viewer (Congregation)
go run cmd/viewer/main.go
```

## 💻 Production Deployment

### Build for Production
```bash
# Build for all platforms
make build-all

# Creates binaries in dist/ folder:
# - broadcaster-windows-amd64.exe
# - viewer-windows-amd64.exe
# - broadcaster-darwin-amd64 (macOS)
# - viewer-linux-amd64
# - broadcaster-linux-arm7 (Raspberry Pi)
```

### Church Setup (Zero Installation)

#### Broadcaster (Church Staff)
1. **Copy file** to laptop: `broadcaster.exe`
2. **Connect camera** (USB webcam or built-in)
3. **Double-click** `broadcaster.exe`
4. **Configure quality** in UI (720p/1080p/480p)
5. **Click "Start Broadcasting"** → Live streaming begins

#### Viewers (Congregation)
1. **Copy file** to laptops: `viewer.exe`
2. **Connect to church WiFi**
3. **Double-click** `viewer.exe`
4. **Click "Connect"** → Automatically finds and joins stream

### Distribution Options

#### Option 1: Direct Download
```
Church downloads from website:
├── broadcaster.exe (Single file)
├── viewer.exe (Single file)
└── Quick setup guide
```

#### Option 2: Hardware Kit ($299)
```
Church Streaming Kit:
├── Raspberry Pi (broadcaster)
├── USB camera
├── WiFi router
├── USB drives with viewer apps
└── Setup instructions
```

#### Option 3: App Store (Future)
```
Microsoft Store / Mac App Store:
├── MeshLink Broadcaster (Church staff)
├── MeshLink Viewer (Congregation)
└── One-click install
```

### Production Benefits
- ✅ **Single executable** - no dependencies to install
- ✅ **Portable** - runs from USB stick or any folder
- ✅ **Self-contained** - all libraries bundled
- ✅ **Cross-platform** - same process for Windows/macOS/Linux
- ✅ **Zero IT support** - just copy and run
- ✅ **Instant setup** - ready in 30 seconds

### Mobile Development (Coming Soon)
```bash
# iOS/Android apps in development
# Will use same P2P core with native mobile UI
```

## 📁 Project Structure
```
├── cmd/                    # Application entry points
├── internal/               # Private application code
├── pkg/                    # Public libraries
├── api/                    # API definitions
├── deployments/            # Docker & deployment configs
└── docs/                   # Documentation
```

## 🛠️ Technology Stack
- **Backend**: Go 1.21+
- **P2P**: libp2p
- **Media**: FFmpeg integration
- **UI**: Fyne (cross-platform)
- **Deployment**: Docker, Kubernetes ready
