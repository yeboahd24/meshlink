package media

import (
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"
	"runtime"
	"strings"

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

// ListVideoDevices lists available video devices (helpful for debugging)
func ListVideoDevices() error {
	var cmd *exec.Cmd
	
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("ffmpeg", "-list_devices", "true", "-f", "dshow", "-i", "dummy")
	case "darwin":
		cmd = exec.Command("ffmpeg", "-f", "avfoundation", "-list_devices", "true", "-i", "")
	default:
		cmd = exec.Command("v4l2-ctl", "--list-devices")
	}
	
	output, err := cmd.CombinedOutput()
	log.Printf("Available devices:\n%s", string(output))
	return err
}

func (f *FFmpegStreamer) Start() error {
	// List devices for debugging
	log.Println("Listing available video devices...")
	ListVideoDevices()
	
	errChan := make(chan error, 1)

	go func() {
		defer f.output.Close()

		var stream *ffmpeg.Stream
		var inputDevice string

		switch runtime.GOOS {
		case "windows":
			// Try to find the first video device
			inputDevice = f.findWindowsCamera()
			if inputDevice == "" {
				errChan <- fmt.Errorf("no video device found on Windows")
				return
			}
			log.Printf("Using Windows DirectShow camera: %s", inputDevice)
			stream = ffmpeg.Input(fmt.Sprintf("video=%s", inputDevice), ffmpeg.KwArgs{
				"f": "dshow",
				"video_size": "640x480",
				"framerate": "30",
			})
			
		case "darwin":
			log.Printf("Using macOS AVFoundation camera")
			stream = ffmpeg.Input("0", ffmpeg.KwArgs{
				"f":          "avfoundation",
				"framerate":  "30",
				"video_size": "640x480",
			})
			
		default:
			log.Printf("Using Linux V4L2 camera")
			stream = ffmpeg.Input("/dev/video0", ffmpeg.KwArgs{
				"f":          "v4l2",
				"framerate":  "30",
				"video_size": "640x480",
			})
		}

		cmd := stream.
			Output("pipe:",
				ffmpeg.KwArgs{
					"f":   "mjpeg",
					"q:v": "5",
				}).
			WithOutput(f.output).
			Compile()

		f.cmd = cmd
		errChan <- nil

		log.Printf("Starting FFmpeg with command: %v", cmd.Args)
		
		if err := cmd.Start(); err != nil {
			log.Printf("FFmpeg start error: %v", err)
			return
		}

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

	return <-errChan
}

func (f *FFmpegStreamer) findWindowsCamera() string {
	cmd := exec.Command("ffmpeg", "-list_devices", "true", "-f", "dshow", "-i", "dummy")
	output, _ := cmd.CombinedOutput()
	
	lines := strings.Split(string(output), "\n")
	inVideoSection := false
	
	for _, line := range lines {
		if strings.Contains(line, "DirectShow video devices") {
			inVideoSection = true
			continue
		}
		if strings.Contains(line, "DirectShow audio devices") {
			inVideoSection = false
			break
		}
		
		if inVideoSection && strings.Contains(line, "\"") {
			// Extract device name between quotes
			start := strings.Index(line, "\"")
			end := strings.LastIndex(line, "\"")
			if start != -1 && end != -1 && start < end {
				deviceName := line[start+1 : end]
				log.Printf("Found video device: %s", deviceName)
				return deviceName
			}
		}
	}
	
	return ""
}

func (f *FFmpegStreamer) Stop() {
	f.cancel()
	if f.cmd != nil && f.cmd.Process != nil {
		f.cmd.Process.Kill()
	}
}