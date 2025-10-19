#!/bin/bash
echo "📹 Testing Camera on WSL..."

echo ""
echo "1. Checking for video devices..."
ls -la /dev/video* 2>/dev/null || echo "❌ No /dev/video* devices found"

echo ""
echo "2. Listing video devices with v4l2-ctl..."
v4l2-ctl --list-devices 2>/dev/null || echo "❌ v4l2-ctl not installed. Run: sudo apt-get install v4l-utils"

echo ""
echo "3. Testing FFmpeg camera access..."
timeout 3 ffmpeg -f v4l2 -i /dev/video0 -frames:v 1 -f null - 2>&1 | grep -i "video\|error\|input"

echo ""
echo "4. Checking FFmpeg version..."
ffmpeg -version | head -n 1

echo ""
echo "✅ Camera test complete!"
echo ""
echo "If camera works, run:"
echo "./meshlink-broadcaster-headless"
