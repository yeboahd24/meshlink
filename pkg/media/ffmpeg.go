package media

import (
	"context"
	"io"
	"log"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type FFmpegStreamer struct {
	input  string
	output io.WriteCloser
	ctx    context.Context
	cancel context.CancelFunc
}

func NewFFmpegStreamer(input string, output io.WriteCloser) *FFmpegStreamer {
	ctx, cancel := context.WithCancel(context.Background())
	return &FFmpegStreamer{
		input:  input,
		output: output,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (f *FFmpegStreamer) Start() error {
	go func() {
		defer f.output.Close()
		
		// Handle different input types
		var stream *ffmpeg.Stream
		if f.input == "testsrc=duration=3600:size=640x480:rate=30" {
			// Use lavfi (libavfilter) for test pattern
			log.Printf("Using FFmpeg test pattern")
			stream = ffmpeg.Input(f.input, ffmpeg.KwArgs{"f": "lavfi"})
		} else if f.input == "video=0" {
			// Windows DirectShow camera
			log.Printf("Using Windows DirectShow camera")
			stream = ffmpeg.Input("video="+f.input, ffmpeg.KwArgs{"f": "dshow"})
		} else if f.input == "0" {
			// macOS AVFoundation camera
			log.Printf("Using macOS AVFoundation camera")
			stream = ffmpeg.Input(f.input, ffmpeg.KwArgs{"f": "avfoundation"})
		} else {
			// Linux V4L2 camera
			log.Printf("Using Linux V4L2 camera: %s", f.input)
			stream = ffmpeg.Input(f.input, ffmpeg.KwArgs{"f": "v4l2"})
		}
		
		err := stream.
			Output("pipe:",
				ffmpeg.KwArgs{
					"f":    "mjpeg",
					"q:v":  "5",
					"s":    "640x480",
					"r":    "30",
				}).
			WithOutput(f.output).
			Run()
		
		if err != nil {
			log.Printf("FFmpeg error: %v", err)
		}
	}()
	
	return nil
}

func (f *FFmpegStreamer) Stop() {
	f.cancel()
}