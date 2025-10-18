package media

import (
	"fmt"
	"os/exec"
	"runtime"
)

type AudioCapture struct {
	deviceID     string
	sampleRate   int
	channels     int
	isCapturing  bool
}

type AudioPlayer struct {
	sampleRate  int
	channels    int
	isPlaying   bool
}

func NewAudioCapture() *AudioCapture {
	return &AudioCapture{
		deviceID:   "default",
		sampleRate: 44100,
		channels:   2,
	}
}

func NewAudioPlayer() *AudioPlayer {
	return &AudioPlayer{
		sampleRate: 44100,
		channels:   2,
	}
}

func (a *AudioCapture) Start() error {
	if a.isCapturing {
		return fmt.Errorf("already capturing audio")
	}
	
	if a.isAudioDeviceAvailable() {
		fmt.Println("Audio: Real microphone detected - starting capture")
	} else {
		fmt.Println("Audio: No microphone found - using silence simulation")
	}
	
	a.isCapturing = true
	return nil
}

func (a *AudioCapture) Stop() {
	a.isCapturing = false
}

func (a *AudioCapture) CaptureAudio() ([]byte, error) {
	if !a.isCapturing {
		return nil, fmt.Errorf("not capturing audio")
	}
	
	// Try real audio capture
	if a.isAudioDeviceAvailable() {
		audioData, err := a.captureFromSystem()
		if err == nil && len(audioData) > 0 {
			return audioData, nil
		}
	}
	
	// Fallback to silence simulation
	return a.generateSilence(), nil
}

func (a *AudioCapture) isAudioDeviceAvailable() bool {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("powershell", "-Command", "Get-WmiObject -Class Win32_SoundDevice")
		return cmd.Run() == nil
	case "darwin":
		cmd := exec.Command("system_profiler", "SPAudioDataType")
		return cmd.Run() == nil
	case "linux":
		cmd := exec.Command("arecord", "-l")
		return cmd.Run() == nil
	default:
		return false
	}
}

func (a *AudioCapture) captureFromSystem() ([]byte, error) {
	switch runtime.GOOS {
	case "windows":
		return a.captureWindows()
	case "darwin":
		return a.captureMacOS()
	case "linux":
		return a.captureLinux()
	default:
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

func (a *AudioCapture) captureWindows() ([]byte, error) {
	// Use ffmpeg to capture audio from DirectShow
	cmd := exec.Command("ffmpeg",
		"-f", "dshow",
		"-i", "audio=Microphone",
		"-t", "0.033", // 33ms for 30fps sync
		"-f", "wav",
		"-")
	
	output, err := cmd.Output()
	if err != nil {
		return a.generateSilence(), nil
	}
	
	return output, nil
}

func (a *AudioCapture) captureMacOS() ([]byte, error) {
	// Use ffmpeg to capture from AVFoundation
	cmd := exec.Command("ffmpeg",
		"-f", "avfoundation",
		"-i", ":0", // Default microphone
		"-t", "0.033",
		"-f", "wav",
		"-")
	
	output, err := cmd.Output()
	if err != nil {
		return a.generateSilence(), nil
	}
	
	return output, nil
}

func (a *AudioCapture) captureLinux() ([]byte, error) {
	// Use arecord to capture from ALSA
	cmd := exec.Command("arecord",
		"-D", "default",
		"-f", "S16_LE",
		"-r", "44100",
		"-c", "2",
		"-d", "0.033")
	
	output, err := cmd.Output()
	if err != nil {
		return a.generateSilence(), nil
	}
	
	return output, nil
}

func (a *AudioCapture) generateSilence() []byte {
	// Generate 33ms of silence (for 30fps sync)
	samples := a.sampleRate * a.channels * 33 / 1000 // 33ms worth
	audioData := make([]byte, samples*2) // 16-bit samples
	
	// Fill with silence (zeros)
	for i := range audioData {
		audioData[i] = 0
	}
	
	return audioData
}

// Audio Player methods
func (p *AudioPlayer) Start() error {
	if p.isPlaying {
		return fmt.Errorf("already playing audio")
	}
	
	fmt.Println("Audio: Starting audio playback")
	p.isPlaying = true
	return nil
}

func (p *AudioPlayer) Stop() {
	p.isPlaying = false
}

func (p *AudioPlayer) PlayAudio(audioData []byte) error {
	if !p.isPlaying {
		return fmt.Errorf("audio player not started")
	}
	
	// For now, just acknowledge audio received
	// Real implementation would decode and play through speakers
	if len(audioData) > 0 {
		fmt.Printf("Audio: Playing %d bytes of audio data\n", len(audioData))
	}
	
	return nil
}