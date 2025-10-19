package main

import (
	"fmt"
	"os/exec"
	"strings"
	"log"
)

func main() {
	fmt.Println("🔍 Windows Camera Detection Test")
	fmt.Println("===============================")
	
	// Test 1: List all DirectShow devices
	fmt.Println("\n1. Listing DirectShow devices...")
	cmd := exec.Command("ffmpeg", "-list_devices", "true", "-f", "dshow", "-i", "dummy")
	output, err := cmd.CombinedOutput()
	
	fmt.Printf("Raw output:\n%s\n", string(output))
	
	// Test 2: Parse camera names
	fmt.Println("\n2. Parsing camera names...")
	cameras := parseWindowsCameras(string(output))
	
	if len(cameras) == 0 {
		fmt.Println("❌ No cameras found!")
		fmt.Println("\n💡 Possible issues:")
		fmt.Println("   - Camera drivers not installed")
		fmt.Println("   - Camera blocked by Windows privacy settings")
		fmt.Println("   - Camera in use by another application")
		fmt.Println("   - DirectShow drivers missing")
		return
	}
	
	// Test 3: Try each camera
	fmt.Printf("\n3. Found %d camera(s), testing each...\n", len(cameras))
	for i, camera := range cameras {
		fmt.Printf("\n📹 Testing camera %d: %s\n", i+1, camera)
		if testCamera(camera) {
			fmt.Printf("✅ Camera %d works: %s\n", i+1, camera)
		} else {
			fmt.Printf("❌ Camera %d failed: %s\n", i+1, camera)
		}
	}
}

func parseWindowsCameras(output string) []string {
	var cameras []string
	lines := strings.Split(output, "\n")
	inVideoSection := false
	
	for _, line := range lines {
		// Detect video devices section
		if strings.Contains(line, "DirectShow video devices") {
			inVideoSection = true
			continue
		}
		if strings.Contains(line, "DirectShow audio devices") {
			inVideoSection = false
			break
		}
		
		// Extract camera name from quotes
		if inVideoSection && strings.Contains(line, "\"") {
			start := strings.Index(line, "\"")
			end := strings.LastIndex(line, "\"")
			if start != -1 && end != -1 && start < end {
				cameraName := line[start+1 : end]
				if cameraName != "" {
					cameras = append(cameras, cameraName)
					fmt.Printf("Found camera: %s\n", cameraName)
				}
			}
		}
	}
	
	return cameras
}

func testCamera(cameraName string) bool {
	fmt.Printf("🔍 Testing camera access: %s\n", cameraName)
	
	// Quick test if camera is accessible
	cmd := exec.Command("ffmpeg",
		"-f", "dshow",
		"-i", fmt.Sprintf("video=%s", cameraName),
		"-frames:v", "1",
		"-f", "null",
		"-",
		"-y")
	
	output, err := cmd.CombinedOutput()
	outputStr := string(output)
	
	// Check for specific error patterns
	if err != nil {
		if strings.Contains(outputStr, "Could not find") {
			fmt.Printf("❌ Camera not found: %s\n", cameraName)
		} else if strings.Contains(outputStr, "Cannot open") {
			fmt.Printf("❌ Camera cannot be opened (in use?): %s\n", cameraName)
		} else if strings.Contains(outputStr, "Permission denied") {
			fmt.Printf("❌ Camera permission denied: %s\n", cameraName)
		} else {
			fmt.Printf("❌ Camera test failed: %s\n", outputStr)
		}
		return false
	}
	
	fmt.Printf("✅ Camera test passed: %s\n", cameraName)
	return true
}