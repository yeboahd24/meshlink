package media

import (
	"context"
	"io"
	"log"
	"runtime"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type FFmpegStreamer struct {
	output io.WriteCloser
	ctx    context.Context
	cancel context.CancelFunc
}

func NewFFmpegStreamer(output io.WriteCloser) *FFmpegStreamer {
	ctx, cancel := context.WithCancel(context.Background())
	return &FFmpegStreamer{
		output: output,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (f *FFmpegStreamer) Start() error {
	go func() {
		defer f.output.Close()

		var stream *ffmpeg.Stream

		switch runtime.GOOS {
		case "windows":
			log.Printf("Using Windows DirectShow camera")
			stream = ffmpeg.Input(`video=Integrated Camera`, ffmpeg.KwArgs{"f": "dshow"})
		case "darwin":
			log.Printf("Using macOS AVFoundation camera")
			stream = ffmpeg.Input("0", ffmpeg.KwArgs{"f": "avfoundation"})
		default:
			log.Printf("Using Linux V4L2 camera")
			stream = ffmpeg.Input("/dev/video0", ffmpeg.KwArgs{"f": "v4l2"})
		}

		err := stream.
			Output("pipe:",
				ffmpeg.KwArgs{
					"f":   "mjpeg",
					"q:v": "5",
					"s":   "640x480",
					"r":   "30",
				}).
			WithContext(f.ctx).
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
