package media

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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

// getFFmpegPath returns path to ffmpeg, checking bundled version first
func getFFmpegPath() string {
	// Check if ffmpeg is bundled with executable
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		
		// Check for bundled ffmpeg
		var bundledFFmpeg string
		if runtime.GOOS == "windows" {
			bundledFFmpeg = filepath.Join(exeDir, "ffmpeg.exe")
		} else {
			bundledFFmpeg = filepath.Join(exeDir, "ffmpeg")
		}
		
		if _, err := os.Stat(bundledFFmpeg); err == nil {
			log.Printf("Using bundled FFmpeg: %s", bundledFFmpeg)
			return bundledFFmpeg
		}
	}
	
	// Fallback to system FFmpeg
	log.Printf("Using system FFmpeg")
	return "ffmpeg"
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
	ffmpegPath := getFFmpegPath()
	var cmd *exec.Cmd
	
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command(ffmpegPath, "-list_devices", "true", "-f", "dshow", "-i", "dummy")
	case "darwin":
		cmd = exec.Command(ffmpegPath, "-f", "avfoundation", "-list_devices", "true", "-i", "")
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
			// Auto-detect camera from FFmpeg device list
			inputDevice = f.detectWindowsCamera()
			
			if inputDevice == "" {
				log.Printf("No camera found, using test pattern")
				stream = ffmpeg.Input("testsrc=duration=3600:size=640x480:rate=30", ffmpeg.KwArgs{"f": "lavfi"})
			} else {
				log.Printf("✅ Using Windows DirectShow camera: %s", inputDevice)
				stream = ffmpeg.Input(fmt.Sprintf("video=%s", inputDevice), ffmpeg.KwArgs{
					"f": "dshow",
					"video_size": "640x480",
					"framerate": "30",
				})
			}
			
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

func (f *FFmpegStreamer) detectWindowsCamera() string {
	ffmpegPath := getFFmpegPath()
	
	// List all DirectShow video devices
	cmd := exec.Command(ffmpegPath, "-list_devices", "true", "-f", "dshow", "-i", "dummy")
	output, _ := cmd.CombinedOutput()
	
	lines := strings.Split(string(output), "\n")
	inVideoSection := false
	
	for _, line := range lines {
		// Detect video devices section
		if strings.Contains(line, "DirectShow video devices") {
			inVideoSection = true
			continue
		}
		if strings.Contains(line, "DirectShow audio devices") {
			inVideoSection = false
			break
		}
		
		// Extract camera name from quotes
		if inVideoSection && strings.Contains(line, "\"") {
			start := strings.Index(line, "\"")
			end := strings.LastIndex(line, "\"")
			if start != -1 && end != -1 && start < end {
				cameraName := line[start+1 : end]
				log.Printf("📹 Found camera: %s", cameraName)
				
				// Test if this camera works
				if f.testWindowsCamera(cameraName) {
					return cameraName
				}
			}
		}
	}
	
	log.Printf("⚠️  No working camera found")
	return ""
}

func (f *FFmpegStreamer) testWindowsCamera(cameraName string) bool {
	ffmpegPath := getFFmpegPath()
	log.Printf("🔍 Testing camera: %s", cameraName)
	
	// Quick test if camera is accessible
	cmd := exec.Command(ffmpegPath,
		"-f", "dshow",
		"-i", fmt.Sprintf("video=%s", cameraName),
		"-frames:v", "1",
		"-f", "null",
		"-")
	
	output, err := cmd.CombinedOutput()
	outputStr := string(output)
	
	// Check if camera was found
	if err != nil && (strings.Contains(outputStr, "Could not find") || 
		strings.Contains(outputStr, "Cannot open")) {
		log.Printf("❌ Camera test failed: %s", cameraName)
		return false
	}
	
	log.Printf("✅ Camera test passed: %s", cameraName)
	return true
}

func (f *FFmpegStreamer) Stop() {
	f.cancel()
	if f.cmd != nil && f.cmd.Process != nil {
		f.cmd.Process.Kill()
	}
}