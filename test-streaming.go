package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

func main() {
	fmt.Println("🎬 Testing P2P-style H.264 Streaming...")

	// Create a pipe to simulate P2P data transfer
	reader, writer := io.Pipe()

	// Start broadcaster (FFmpeg encoder)
	go startBroadcaster(writer)

	// Start viewer (FFmpeg decoder) 
	go startViewer(reader)

	// Let it stream for 10 seconds
	fmt.Println("📡 Streaming for 10 seconds...")
	time.Sleep(10 * time.Second)

	fmt.Println("🛑 Stopping stream...")
	writer.Close()
	reader.Close()

	fmt.Println("✅ Streaming test completed!")
}

func startBroadcaster(output io.WriteCloser) {
	defer output.Close()
	
	fmt.Println("📹 Starting broadcaster (H.264 encoder)...")
	
	cmd := exec.Command("ffmpeg",
		"-f", "lavfi",
		"-i", "testsrc=size=640x480:rate=30",
		"-c:v", "libx264",
		"-preset", "ultrafast", 
		"-tune", "zerolatency",
		"-f", "h264",
		"-",
	)
	
	cmd.Stdout = output
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Broadcaster error: %v\n", err)
	}
}

func startViewer(input io.ReadCloser) {
	defer input.Close()
	
	fmt.Println("📺 Starting viewer (H.264 decoder)...")
	
	cmd := exec.Command("ffmpeg",
		"-f", "h264",
		"-i", "-",
		"-f", "rawvideo",
		"-pix_fmt", "rgb24",
		"-s", "640x480",
		"-r", "30",
		"/tmp/decoded_frames.raw",
	)
	
	cmd.Stdin = input
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Viewer error: %v\n", err)
	}
	
	// Check decoded output
	if stat, err := os.Stat("/tmp/decoded_frames.raw"); err == nil {
		fmt.Printf("📊 Decoded %d bytes of video data\n", stat.Size())
	}
}