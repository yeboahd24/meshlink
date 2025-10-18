package main

import (
	"fmt"
	"meshlink/internal/media"
)

func main() {
	fmt.Println("Testing camera capture...")
	
	camera := media.NewCameraCapture()
	
	err := camera.Start()
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	
	fmt.Println("Camera started successfully!")
	
	frame, err := camera.CaptureFrame()
	if err != nil {
		fmt.Printf("ERROR capturing frame: %v\n", err)
		return
	}
	
	fmt.Printf("SUCCESS: Captured frame of %d bytes\n", len(frame))
	camera.Stop()
}