# Docker Build Fixed

## Issues Resolved

1. **Removed unused FFmpeg import** from `internal/media/capture.go`
2. **Fixed undefined methods** in CaptureFrame - now uses fallback simulation
3. **Simplified camera capture** for Docker/headless environment

## Build Command

```bash
docker build -f Dockerfile.simple -t meshlink-simple .
```

## Run Command

```bash
docker run -it --rm -p 8080:8080 meshlink-simple
```

The broadcaster will now compile successfully in Docker and run in headless mode with simulated video frames.