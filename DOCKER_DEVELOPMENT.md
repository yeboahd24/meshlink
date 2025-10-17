# Docker Development Environment

## Quick Start for Teams

### 1. Clone and Start Development Environment
```bash
git clone <repository>
cd MeshLink
make docker-dev
```

### 2. Enter Development Container
```bash
make docker-dev-shell
```

### 3. Build and Test Inside Container
```bash
# Inside container
make build
make test
```

## Development Workflow

### Start Development Environment
```bash
# Start dev container with all dependencies
make docker-dev

# Enter the container
make docker-dev-shell
```

### Build Applications
```bash
# Build inside container (recommended)
make docker-dev-build

# Or enter container and build manually
make docker-dev-shell
# Inside container:
make build
```

### Run Applications
```bash
# Run broadcaster in container
make docker-dev-run-broadcaster

# Run viewer in container (separate terminal)
make docker-dev-run-viewer
```

### Testing
```bash
# Run tests in container
make docker-dev-test

# Or enter container for interactive testing
make docker-dev-shell
# Inside container:
make test
go test ./...
```

## Full Demo Environment

### Start Everything (Dev + Production)
```bash
# Starts dev environment + production demo
make docker-full-demo
```

This provides:
- **Development Container**: `/workspace` with all tools
- **Production Demo**: Broadcaster + 2 Viewers + Web UI
- **Web Dashboard**: http://localhost:3000

### Access Points
- **Dev Environment**: `docker compose -f deployments/docker-compose.dev.yml exec dev-env bash`
- **Demo Dashboard**: http://localhost:3000
- **Production Broadcaster**: Port 9080
- **Production Viewers**: Ports 9081, 9082

## Container Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Docker Host                              │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │   Dev Container │  │  Prod Containers │  │  Demo UI    │ │
│  │                 │  │                  │  │             │ │
│  │ - Go 1.21       │  │ - Broadcaster    │  │ - Dashboard │ │
│  │ - GStreamer     │  │ - Viewer1        │  │ - Monitoring│ │
│  │ - All deps      │  │ - Viewer2        │  │             │ │
│  │ - Source code   │  │                  │  │             │ │
│  │ - Build tools   │  │                  │  │             │ │
│  └─────────────────┘  └─────────────────┘  └─────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## Team Sharing

### Share with Team
1. **Commit code**: `git push`
2. **Team pulls**: `git pull`
3. **Team runs**: `make docker-dev`

### No Local Installation Required
- ✅ Go compiler in container
- ✅ GStreamer in container  
- ✅ All dependencies in container
- ✅ Build tools in container
- ✅ Consistent environment across team

### Development Features
- **Volume Mounting**: Live code changes reflected
- **Port Forwarding**: Access apps from host
- **Persistent Cache**: Go modules cached between runs
- **Shell Access**: Full development environment

## Commands Reference

### Development Container
```bash
make docker-dev           # Start dev environment
make docker-dev-shell     # Enter container shell
make docker-dev-build     # Build applications
make docker-dev-test      # Run tests
make docker-dev-stop      # Stop dev environment
```

### Inside Container
```bash
make build               # Build broadcaster + viewer
make test                # Run all tests
make dev-broadcaster     # Run broadcaster with hot reload
make dev-viewer          # Run viewer with hot reload
make clean               # Clean build artifacts
```

### Production Demo
```bash
make docker-run          # Simple production demo
make docker-full-demo    # Full demo (dev + prod)
```

## Troubleshooting

### Container Won't Start
```bash
# Check Docker status
docker ps

# View logs
docker compose -f deployments/docker-compose.dev.yml logs dev-env

# Rebuild container
docker compose -f deployments/docker-compose.dev.yml build --no-cache dev-env
```

### Build Failures
```bash
# Enter container and debug
make docker-dev-shell

# Inside container, check Go environment
go version
go env
go mod tidy
```

### Port Conflicts
```bash
# Check what's using ports
netstat -tulpn | grep :8080

# Use different ports in docker-compose.dev.yml
ports:
  - "8090:8080"  # Change host port
```

## IDE Integration

### VS Code with Remote Containers
1. Install "Remote - Containers" extension
2. Open project in VS Code
3. Command: "Remote-Containers: Reopen in Container"
4. VS Code runs inside the dev container

### JetBrains GoLand
1. Configure Docker as remote interpreter
2. Point to dev container Go binary
3. Full IDE features with containerized Go

This setup ensures every team member has identical development environment without local installations.