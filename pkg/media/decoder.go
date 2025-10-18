package media

import (
	"context"
	"io"
	"log"
	"os/exec"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type FFmpegDecoder struct {
	input  io.ReadCloser
	output io.WriteCloser
	ctx    context.Context
	cancel context.CancelFunc
	cmd    *exec.Cmd
}

func NewFFmpegDecoder(input io.ReadCloser, output io.WriteCloser) *FFmpegDecoder {
	ctx, cancel := context.WithCancel(context.Background())
	return &FFmpegDecoder{
		input:  input,
		output: output,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (d *FFmpegDecoder) Start() error {
	errChan := make(chan error, 1)

	go func() {
		defer func() {
			if err := d.output.Close(); err != nil {
				log.Printf("Failed to close decoder output: %v", err)
			}
			if err := d.input.Close(); err != nil {
				log.Printf("Failed to close decoder input: %v", err)
			}
		}()

		// Build the command
		cmd := ffmpeg.Input("pipe:",
			ffmpeg.KwArgs{
				"f": "mjpeg",
			}).
			Output("pipe:",
				ffmpeg.KwArgs{
					"f":       "rawvideo",
					"pix_fmt": "rgb24",
					"s":       "640x480",
				}).
			WithInput(d.input).
			WithOutput(d.output).
			Compile()

		d.cmd = cmd

		// Signal startup
		errChan <- nil

		// Start the process
		if err := cmd.Start(); err != nil {
			log.Printf("FFmpeg decoder start error: %v", err)
			return
		}

		// Wait for either context cancellation or process completion
		done := make(chan error)
		go func() {
			done <- cmd.Wait()
		}()

		select {
		case <-d.ctx.Done():
			log.Printf("Decoder context cancelled, stopping FFmpeg")
			if cmd.Process != nil {
				cmd.Process.Kill()
			}
		case err := <-done:
			if err != nil {
				log.Printf("FFmpeg decode error: %v", err)
			}
		}
	}()

	return <-errChan
}

func (d *FFmpegDecoder) Stop() {
	d.cancel()
	if d.cmd != nil && d.cmd.Process != nil {
		d.cmd.Process.Kill()
	}
}