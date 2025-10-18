#!/bin/bash
echo "🧪 Testing FFmpeg in Docker..."

echo ""
echo "Building FFmpeg test image..."
docker build -f Dockerfile.ffmpeg-test -t meshlink-ffmpeg-test .

echo ""
echo "Running FFmpeg test..."
docker run --rm -it meshlink-ffmpeg-test

echo ""
echo "✅ FFmpeg test completed!"