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

## 🚀 Quick Start

### Broadcaster (Church Setup)
```bash
go run cmd/broadcaster/main.go
```

### Viewer (Congregation)
```bash
go run cmd/viewer/main.go
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
- **Media**: GStreamer integration
- **UI**: Fyne (cross-platform)
- **Deployment**: Docker, Kubernetes ready

## 📊 Market Opportunity
- 300,000+ churches in the US
- Growing demand for hybrid services
- Cost-effective alternative to cloud streaming
- Emerging markets with limited internet

## 🎯 Investment Highlights
- Proven P2P technology stack
- Minimal infrastructure costs
- Scalable SaaS model potential
- Strong technical foundation