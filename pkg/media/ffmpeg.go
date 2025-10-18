package media

import (
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"
	"runtime"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type FFmpegStreamer struct {
	output io.WriteCloser
	ctx    context.Context
	cancel context.CancelFunc
	cmd    *exec.Cmd
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
	errChan := make(chan error, 1)

	go func() {
		defer f.output.Close()

		var stream *ffmpeg.Stream

		switch runtime.GOOS {
		case "windows":
			log.Printf("Using Windows DirectShow camera")
			stream = ffmpeg.Input(`video=Integrated Camera`, ffmpeg.KwArgs{
				"f": "dshow",
			})
		case "darwin":
			log.Printf("Using macOS AVFoundation camera")
			stream = ffmpeg.Input("0", ffmpeg.KwArgs{
				"f":           "avfoundation",
				"framerate":   "30",
				"video_size":  "640x480",
			})
		default:
			log.Printf("Using Linux V4L2 camera")
			stream = ffmpeg.Input("/dev/video0", ffmpeg.KwArgs{
				"f":          "v4l2",
				"framerate":  "30",
				"video_size": "640x480",
			})
		}

		// Build the command but don't run yet
		cmd := stream.
			Output("pipe:",
				ffmpeg.KwArgs{
					"f":      "mjpeg",
					"q:v":    "5",
					"s":      "640x480",
					"r":      "30",
				}).
			WithOutput(f.output).
			Compile()

		// Store cmd so we can kill it later
		f.cmd = cmd
		
		// Signal that we're starting
		errChan <- nil

		// Run with context
		if err := cmd.Start(); err != nil {
			log.Printf("FFmpeg start error: %v", err)
			return
		}

		// Wait for either context cancellation or process completion
		done := make(chan error)
		go func() {
			done <- cmd.Wait()
		}()

		select {
		case <-f.ctx.Done():
			log.Printf("Context cancelled, stopping FFmpeg")
			if cmd.Process != nil {
				cmd.Process.Kill()
			}
		case err := <-done:
			if err != nil {
				log.Printf("FFmpeg error: %v", err)
			}
		}
	}()

	// Wait for startup or error
	return <-errChan
}

func (f *FFmpegStreamer) Stop() {
	f.cancel()
	if f.cmd != nil && f.cmd.Process != nil {
		f.cmd.Process.Kill()
	}
}