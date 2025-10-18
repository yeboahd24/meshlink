//go:build gstreamer
// +build gstreamer

package media

import (
	"bytes"
	"fmt"
	"os/exec"
)

type GStreamerPipeline struct {
	isRunning bool
	quality   string
}

func NewGStreamerPipeline(quality string) *GStreamerPipeline {
	return &GStreamerPipeline{
		quality: quality,
	}
}

func (g *GStreamerPipeline) Start() error {
	if g.isRunning {
		return fmt.Errorf("pipeline already running")
	}
	
	g.isRunning = true
	fmt.Printf("GStreamer: Command-line pipeline ready for %s\n", g.quality)
	return nil
}

func (g *GStreamerPipeline) Stop() {
	g.isRunning = false
	fmt.Println("GStreamer: Pipeline stopped")
}

func (g *GStreamerPipeline) ConvertH264ToJPEG(h264Data []byte) ([]byte, error) {
	if !g.isRunning {
		return nil, fmt.Errorf("pipeline not running")
	}

	// Use command-line GStreamer for H.264 to JPEG conversion
	cmd := exec.Command("gst-launch-1.0",
		"fdsrc", "!",
		"h264parse", "!",
		"avdec_h264", "!",
		"videoconvert", "!",
		"jpegenc", "quality=90", "!",
		"filesink", "location=/dev/stdout",
		"-e")
	
	cmd.Stdin = bytes.NewReader(h264Data)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("GStreamer conversion failed: %w", err)
	}
	
	return output, nil
}