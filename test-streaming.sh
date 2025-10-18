#!/bin/bash
echo "🎬 Testing P2P-style Streaming..."

echo ""
echo "Building streaming test..."
docker build -f Dockerfile.streaming-test -t meshlink-streaming-test .

echo ""
echo "Running streaming test..."
docker run --rm -it meshlink-streaming-test

echo ""
echo "✅ Streaming test completed!"