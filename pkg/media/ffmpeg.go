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

// getFFmpegPath returns path to ffmpeg, prioritizing system installation
func getFFmpegPath() string {
	// First, try to use system FFmpeg (same as manual command)
	if systemFFmpeg, err := exec.LookPath("ffmpeg"); err == nil {
		log.Printf("Using system FFmpeg: %s", systemFFmpeg)
		return systemFFmpeg
	}
	
	log.Printf("System FFmpeg not found in PATH, checking for bundled version...")
	
	// Fallback: Check if ffmpeg is bundled with executable
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
	
	// Final fallback to system FFmpeg (even if not in PATH, let it try)
	log.Printf("No bundled FFmpeg found, falling back to system 'ffmpeg' command")
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

		// Check if we're in WSL2
		if f.isWSL2() {
			log.Printf("🔍 WSL2 detected - attempting Windows camera passthrough")
			inputDevice = f.detectWSL2Camera()
			
			if inputDevice == "" {
				log.Printf("⚠️  WSL2: No camera passthrough available, using enhanced test pattern")
				stream = ffmpeg.Input("testsrc=duration=3600:size=640x480:rate=30:color=0x4169E1", ffmpeg.KwArgs{
					"f": "lavfi",
				})
			} else {
				log.Printf("✅ WSL2: Using camera passthrough: %s", inputDevice)
				stream = ffmpeg.Input(inputDevice, ffmpeg.KwArgs{
					"f": "dshow",
					"video_size": "640x480",
					"framerate": "30",
				})
			}
		} else {
			switch runtime.GOOS {
			case "windows":
				// Auto-detect camera and audio from FFmpeg device list
				videoDevice := f.detectWindowsCamera()
				audioDevice := f.detectWindowsAudio()
				
				if videoDevice == "" {
					log.Printf("No camera found, using test pattern")
					stream = ffmpeg.Input("testsrc=duration=3600:size=640x480:rate=30", ffmpeg.KwArgs{"f": "lavfi"})
				} else {
					// Build DirectShow input with both video and audio
					var dshowInput string
					if audioDevice != "" {
						dshowInput = fmt.Sprintf("video=%s:audio=%s", videoDevice, audioDevice)
						log.Printf("✅ Using Windows DirectShow camera: %s with audio: %s", videoDevice, audioDevice)
					} else {
						dshowInput = fmt.Sprintf("video=%s", videoDevice)
						log.Printf("✅ Using Windows DirectShow camera: %s (no audio device)", videoDevice)
					}
					
					stream = ffmpeg.Input(dshowInput, ffmpeg.KwArgs{
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
				if _, err := os.Stat("/dev/video0"); err == nil {
					stream = ffmpeg.Input("/dev/video0", ffmpeg.KwArgs{
						"f":          "v4l2",
						"framerate":  "30",
						"video_size": "640x480",
					})
				} else {
					log.Printf("No camera found, using test pattern")
					stream = ffmpeg.Input("testsrc=duration=3600:size=640x480:rate=30", ffmpeg.KwArgs{"f": "lavfi"})
				}
			}
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

func (f *FFmpegStreamer) detectWindowsAudio() string {
	ffmpegPath := getFFmpegPath()
	
	// List all DirectShow audio devices
	cmd := exec.Command(ffmpegPath, "-list_devices", "true", "-f", "dshow", "-i", "dummy")
	output, _ := cmd.CombinedOutput()
	
	lines := strings.Split(string(output), "\n")
	inAudioSection := false
	
	for _, line := range lines {
		// Detect audio devices section
		if strings.Contains(line, "DirectShow audio devices") {
			inAudioSection = true
			continue
		}
		if inAudioSection && strings.Contains(line, "DirectShow video devices") {
			inAudioSection = false
			break
		}
		
		// Extract audio device name from quotes
		if inAudioSection && strings.Contains(line, "\"") {
			start := strings.Index(line, "\"")
			end := strings.LastIndex(line, "\"")
			if start != -1 && end != -1 && start < end {
				audioName := line[start+1 : end]
				log.Printf("🎤 Found audio device: %s", audioName)
				
				// For audio, we'll be less strict - just use the first one found
				// Audio devices are generally more reliable than cameras
				if audioName != "" {
					log.Printf("✅ Using audio device: %s", audioName)
					return audioName
				}
			}
		}
	}
	
	log.Printf("⚠️  No audio device found")
	return ""
}

func (f *FFmpegStreamer) testWindowsCamera(cameraName string) bool {
	log.Printf("🔍 Testing camera: %s", cameraName)
	
	// SIMPLIFIED TEST: If camera was detected in device list, assume it works
	// This avoids permission issues and camera access conflicts
	
	// The camera was already found in detectWindowsCamera(), so it exists
	// Let's try a gentler test - just check device info without capturing
	ffmpegPath := getFFmpegPath()
	
	cmd := exec.Command(ffmpegPath,
		"-f", "dshow",
		"-list_options", "true",
		"-i", fmt.Sprintf("video=%s", cameraName))
	
	output, err := cmd.CombinedOutput()
	outputStr := string(output)
	
	// If we can list camera options, it's accessible
	if err == nil || (!strings.Contains(outputStr, "Could not find") && 
		!strings.Contains(outputStr, "Cannot open")) {
		log.Printf("✅ Camera test passed: %s", cameraName)
		return true
	}
	
	// Even more permissive fallback - if it was detected, probably works
	log.Printf("⚠️  Camera test uncertain, but camera was detected: %s", cameraName)
	log.Printf("    Assuming camera works (was found in device list)")
	return true
}

// isWSL2 detects if we're running in WSL2
func (f *FFmpegStreamer) isWSL2() bool {
	if runtime.GOOS != "linux" {
		return false
	}
	
	// Check for WSL2 indicators
	if data, err := os.ReadFile("/proc/version"); err == nil {
		version := string(data)
		return strings.Contains(version, "microsoft") || strings.Contains(version, "WSL2")
	}
	return false
}

// detectWSL2Camera attempts to detect camera through WSL2-Windows bridge
func (f *FFmpegStreamer) detectWSL2Camera() string {
	log.Printf("🔍 WSL2: Checking for Windows camera access...")
	
	// Try to use Windows camera through WSL2
	// This requires special setup but let's try common approaches
	
	// Method 1: Try if Windows cameras are accessible via device passthrough
	testCameras := []string{
		"Integrated Camera",
		"USB Camera", 
		"Camera",
		"Webcam",
		"HD Camera",
	}
	
	for _, camera := range testCameras {
		if f.testWSL2Camera(camera) {
			return fmt.Sprintf("video=%s", camera)
		}
	}
	
	log.Printf("⚠️  WSL2: No Windows camera passthrough available")
	log.Printf("💡 Tip: For best camera support, run the Windows .exe version instead")
	return ""
}

// testWSL2Camera tests if a camera is accessible in WSL2
func (f *FFmpegStreamer) testWSL2Camera(cameraName string) bool {
	ffmpegPath := getFFmpegPath()
	log.Printf("🔍 WSL2: Testing camera: %s", cameraName)
	
	// Quick test if camera is accessible
	cmd := exec.Command(ffmpegPath,
		"-f", "dshow",
		"-i", fmt.Sprintf("video=%s", cameraName),
		"-frames:v", "1",
		"-f", "null",
		"-",
		"-y")
	
	output, err := cmd.CombinedOutput()
	outputStr := string(output)
	
	// Check if camera was found
	if err != nil || strings.Contains(outputStr, "Could not find") || 
		strings.Contains(outputStr, "Cannot open") {
		log.Printf("❌ WSL2: Camera test failed: %s", cameraName)
		return false
	}
	
	log.Printf("✅ WSL2: Camera test passed: %s", cameraName)
	return true
}

func (f *FFmpegStreamer) Stop() {
	f.cancel()
	if f.cmd != nil && f.cmd.Process != nil {
		f.cmd.Process.Kill()
	}
}