//go:build !gstreamer
// +build !gstreamer

package media

import "fmt"

type GStreamerPipeline struct {
	quality string
}

func NewGStreamerPipeline(quality string) *GStreamerPipeline {
	return &GStreamerPipeline{quality: quality}
}

func (g *GStreamerPipeline) Start() error {
	fmt.Println("GStreamer: Not available (build without gstreamer tag)")
	return nil
}

func (g *GStreamerPipeline) Stop() {
	// No-op
}

func (g *GStreamerPipeline) ConvertH264ToJPEG(h264Data []byte) ([]byte, error) {
	return nil, fmt.Errorf("GStreamer not available")
}