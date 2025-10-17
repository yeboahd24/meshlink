# MeshLink Docker Demo Environment

## Quick Start

### Windows
```cmd
docker-demo.bat
```

### Linux/macOS
```bash
chmod +x docker-demo.sh
./docker-demo.sh
```

## Manual Setup

### 1. Build and Start
```bash
cd MeshLink
docker compose -f deployments/docker-compose.yml up --build
```

### 2. Access Demo
- **Demo Dashboard**: http://localhost:3000
- **Broadcaster Logs**: `docker logs meshlink-broadcaster-1`
- **Viewer Logs**: `docker logs meshlink-viewer1-1`

### 3. Stop Demo
```bash
docker compose -f deployments/docker-compose.yml down
```

## Demo Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Broadcaster   │    │   Docker Net    │    │    Viewers      │
│   Port: 8080    │◄──►│   172.18.0.0/16 │◄──►│  Ports: 8081-82 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │    Demo UI      │
                    │   Port: 3000    │
                    └─────────────────┘
```

## Services

### Broadcaster Container
- **Image**: Built from `Dockerfile.broadcaster`
- **Function**: Streams data to P2P network
- **Mode**: Headless (no GUI)
- **Auto-start**: Begins broadcasting immediately

### Viewer Containers (2x)
- **Image**: Built from `Dockerfile.viewer`  
- **Function**: Receives P2P stream data
- **Mode**: Headless (no GUI)
- **Auto-connect**: Joins stream automatically

### Demo UI Container
- **Image**: Built from `Dockerfile.demo`
- **Function**: Web dashboard showing demo status
- **Access**: http://localhost:3000
- **Features**: Live logs, status monitoring

## Demo Features

### ✅ P2P Networking
- Automatic peer discovery via mDNS
- libp2p mesh network formation
- Direct container-to-container communication

### ✅ Offline Operation
- No external internet required
- Isolated Docker network simulates local WiFi
- All communication stays within containers

### ✅ Scalability Demo
- 1 broadcaster → 2 viewers (easily scalable)
- Real-time data distribution
- Low latency local networking

### ✅ Monitoring
- Live status dashboard
- Container logs and metrics
- Network topology visualization

## Investor Demo Points

### 1. Zero Configuration
- Containers start automatically
- No manual network setup
- Plug-and-play operation

### 2. Offline-First
- Disconnect internet → demo continues
- Local network only
- No cloud dependencies

### 3. Scalability
- Add viewers: `docker-compose up --scale viewer1=5`
- Handles multiple concurrent streams
- Efficient P2P distribution

### 4. Production Ready
- Containerized deployment
- Health monitoring
- Graceful shutdown handling

## Troubleshooting

### Port Conflicts
```bash
# Check if ports are in use
netstat -an | findstr "3000 8080 8081 8082"

# Use different ports
docker-compose -f deployments/docker-compose.yml up -p 3001:3000
```

### Container Issues
```bash
# View logs
docker compose -f deployments/docker-compose.yml logs broadcaster
docker compose -f deployments/docker-compose.yml logs viewer1

# Restart services
docker compose -f deployments/docker-compose.yml restart
```

### Network Problems
```bash
# Inspect Docker network
docker network ls
docker network inspect meshlink-meshlink

# Reset network
docker compose -f deployments/docker-compose.yml down
docker network prune
```

## Extending the Demo

### Add More Viewers
```bash
# Scale existing viewers
docker compose -f deployments/docker-compose.yml up --scale viewer1=3

# Or add new service in docker-compose.yml
viewer3:
  build:
    context: ..
    dockerfile: deployments/Dockerfile.viewer
  ports:
    - "8083:8080"
  # ... same config as viewer1
```

### Enable Real Video
1. Add camera access to broadcaster container
2. Install GStreamer in Docker images
3. Update media pipeline in Go code

### Cloud Integration
1. Add API server container
2. Connect to external database
3. Enable analytics and monitoring