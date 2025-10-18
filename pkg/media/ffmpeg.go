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
		
		// Handle test pattern input differently
		var stream *ffmpeg.Stream
		if f.input == "testsrc=duration=3600:size=1280x720:rate=30" {
			// Use lavfi (libavfilter) for test pattern
			stream = ffmpeg.Input(f.input, ffmpeg.KwArgs{"f": "lavfi"})
		} else {
			// Regular camera input
			stream = ffmpeg.Input(f.input, ffmpeg.KwArgs{"f": "v4l2"})
		}
		
		err := stream.
			Output("pipe:",
				ffmpeg.KwArgs{
					"format":  "h264",
					"vcodec":  "libx264",
					"preset":  "ultrafast",
					"tune":    "zerolatency",
					"crf":     "23",
					"s":       "1280x720",
					"r":       "30",
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