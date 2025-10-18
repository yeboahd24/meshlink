package media

import (
	"context"
	"io"
	"log"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type FFmpegDecoder struct {
	input  io.ReadCloser
	output io.WriteCloser
	ctx    context.Context
	cancel context.CancelFunc
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
	go func() {
		defer d.output.Close()
		
		err := ffmpeg.Input("pipe:",
			ffmpeg.KwArgs{
				"format": "mjpeg",
			}).
			Output("pipe:",
				ffmpeg.KwArgs{
					"format":  "rawvideo",
					"pix_fmt": "rgb24",
					"s":       "640x480",
				}).
			WithInput(d.input).
			WithOutput(d.output).
			Run()
		
		if err != nil {
			log.Printf("FFmpeg decode error: %v", err)
		}
	}()
	
	return nil
}

func (d *FFmpegDecoder) Stop() {
	d.cancel()
}