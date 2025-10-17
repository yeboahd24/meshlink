# MeshLink Architecture

## Overview
MeshLink is designed as a decentralized, offline-first streaming platform specifically for churches and religious communities.

## Core Principles

### 1. Offline-First Design
- No internet dependency for core functionality
- Works on local WiFi networks and hotspots
- Graceful degradation when connectivity is limited

### 2. Peer-to-Peer Architecture
- Direct communication between broadcaster and viewers
- No central server required for streaming
- Efficient bandwidth utilization through mesh networking

### 3. Scalability
- Support for 1 broadcaster to 50+ viewers
- Optimized for low-resource environments
- Raspberry Pi compatible for broadcaster nodes

## Technical Architecture

### Network Layer
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Broadcaster   │    │   P2P Network   │    │     Viewers     │
│                 │    │                 │    │                 │
│ ┌─────────────┐ │    │ ┌─────────────┐ │    │ ┌─────────────┐ │
│ │   libp2p    │◄┼────┼►│    mDNS     │◄┼────┼►│   libp2p    │ │
│ │   Host      │ │    │ │  Discovery  │ │    │ │   Host      │ │
│ └─────────────┘ │    │ └─────────────┘ │    │ └─────────────┘ │
│                 │    │                 │    │                 │
│ ┌─────────────┐ │    │ ┌─────────────┐ │    │ ┌─────────────┐ │
│ │  PubSub     │◄┼────┼►│   Topics    │◄┼────┼►│  PubSub     │ │
│ │ Publisher   │ │    │ │             │ │    │ │ Subscriber  │ │
│ └─────────────┘ │    │ └─────────────┘ │    │ └─────────────┘ │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Application Layer
```
┌─────────────────────────────────────────────────────────────┐
│                    Application Layer                        │
├─────────────────────────────────────────────────────────────┤
│  Broadcaster App           │           Viewer App           │
│  ┌─────────────────────┐   │   ┌─────────────────────────┐  │
│  │   Fyne UI           │   │   │      Fyne UI            │  │
│  │   - Start/Stop      │   │   │   - Connect/Disconnect  │  │
│  │   - Status Display  │   │   │   - Video Display       │  │
│  │   - Settings        │   │   │   - Audio Controls      │  │
│  └─────────────────────┘   │   └─────────────────────────┘  │
├─────────────────────────────────────────────────────────────┤
│                    Streaming Layer                          │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────────────┐   │   ┌─────────────────────────┐  │
│  │   Media Capture     │   │   │   Media Playback        │  │
│  │   - Camera Input    │   │   │   - Video Decode        │  │
│  │   - Audio Input     │   │   │   - Audio Decode        │  │
│  │   - Encoding        │   │   │   - Rendering           │  │
│  └─────────────────────┘   │   └─────────────────────────┘  │
├─────────────────────────────────────────────────────────────┤
│                      P2P Layer                              │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────────────────────────────────────────────────┐ │
│  │                    libp2p Stack                         │ │
│  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────────┐   │ │
│  │  │   mDNS      │ │   PubSub    │ │   Connection    │   │ │
│  │  │  Discovery  │ │  Messaging  │ │   Management    │   │ │
│  │  └─────────────┘ └─────────────┘ └─────────────────┘   │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## Data Flow

### Broadcasting Flow
1. **Media Capture**: Camera/microphone input via system APIs
2. **Encoding**: H.264/AAC encoding for efficient transmission
3. **Chunking**: Split encoded stream into manageable chunks
4. **P2P Publishing**: Broadcast chunks via libp2p PubSub
5. **Discovery**: Advertise stream availability via mDNS

### Viewing Flow
1. **Discovery**: Scan for available streams via mDNS
2. **Connection**: Establish P2P connection to broadcaster
3. **Subscription**: Subscribe to stream topic via PubSub
4. **Reception**: Receive and buffer stream chunks
5. **Playback**: Decode and render audio/video

## Scalability Considerations

### Network Efficiency
- **Gossip Protocol**: Efficient message propagation
- **Bandwidth Adaptation**: Dynamic quality adjustment
- **Peer Relay**: Viewers can relay to other viewers

### Resource Management
- **Memory Buffering**: Configurable buffer sizes
- **CPU Optimization**: Hardware-accelerated encoding/decoding
- **Storage**: Minimal local storage requirements

## Security Model

### P2P Security
- **Encrypted Connections**: All P2P traffic encrypted via libp2p
- **Peer Authentication**: Identity verification for trusted networks
- **Topic Access Control**: Stream access via discovery keys

### Privacy Protection
- **Local Processing**: No data sent to external servers
- **Ephemeral Connections**: No persistent user tracking
- **Configurable Logging**: Minimal data retention

## Future Extensions

### Cloud Integration
- **Analytics API**: Optional usage statistics
- **User Management**: Church registration and management
- **Content Delivery**: Hybrid P2P/CDN for larger audiences

### Advanced Features
- **Multi-Camera Support**: Multiple video sources
- **Interactive Features**: Chat, polls, donations
- **Recording**: Local and cloud recording options