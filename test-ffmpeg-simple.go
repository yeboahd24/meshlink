package main

import (
	"fmt"
	"os/exec"
	"time"
)

func main() {
	fmt.Println("🎥 Testing FFmpeg H.264 Streaming...")

	// Test 1: Generate test pattern and encode to H.264
	fmt.Println("📡 Starting H.264 test pattern stream...")
	
	ffmpegCmd := exec.Command("ffmpeg",
		"-f", "lavfi",                           // Use libavfilter
		"-i", "testsrc=duration=10:size=640x480:rate=30", // 10 second test pattern
		"-c:v", "libx264",                       // H.264 encoder
		"-preset", "ultrafast",                  // Fast encoding
		"-tune", "zerolatency",                  // Low latency
		"-f", "h264",                           // Raw H.264 output
		"-y", "/tmp/test_output.h264",          // Output file
	)

	// Start FFmpeg
	err := ffmpegCmd.Start()
	if err != nil {
		fmt.Printf("❌ Error starting FFmpeg: %v\n", err)
		return
	}

	fmt.Println("✅ FFmpeg started - generating H.264 stream...")
	
	// Let it run for a few seconds
	time.Sleep(3 * time.Second)
	
	// Check if process is still running
	if ffmpegCmd.Process != nil {
		fmt.Println("📊 FFmpeg process running, PID:", ffmpegCmd.Process.Pid)
	}

	// Wait for completion
	err = ffmpegCmd.Wait()
	if err != nil {
		fmt.Printf("⚠️  FFmpeg finished with: %v\n", err)
	} else {
		fmt.Println("✅ FFmpeg completed successfully!")
	}

	// Test 2: Check output file
	checkCmd := exec.Command("ls", "-la", "/tmp/test_output.h264")
	output, err := checkCmd.Output()
	if err != nil {
		fmt.Printf("❌ Could not check output: %v\n", err)
	} else {
		fmt.Printf("📁 Output file info:\n%s", output)
	}

	// Test 3: Probe the H.264 file
	fmt.Println("🔍 Probing H.264 output...")
	probeCmd := exec.Command("ffprobe", 
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		"/tmp/test_output.h264",
	)
	
	probeOutput, err := probeCmd.Output()
	if err != nil {
		fmt.Printf("❌ Probe failed: %v\n", err)
	} else {
		fmt.Printf("📋 H.264 stream info:\n%s\n", probeOutput)
	}

	fmt.Println("🎉 FFmpeg H.264 test completed!")
}