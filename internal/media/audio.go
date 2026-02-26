package media

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"meshlink/pkg/media"
)

type AudioCapture struct {
	deviceID    string
	sampleRate  int
	channels    int
	isCapturing bool

	cmd         *exec.Cmd
	pipeReader  io.ReadCloser
	latestChunk []byte
	chunkMu     sync.Mutex
	stopChan    chan struct{}
}

type AudioPlayer struct {
	sampleRate int
	channels   int
	isPlaying  bool
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

	a.stopChan = make(chan struct{})

	device := a.detectAudioDevice()
	if device == "" {
		log.Println("Audio: No microphone found — will produce silence")
		a.isCapturing = true
		return nil
	}

	args := a.buildFFmpegArgs(device)
	ffmpegPath := media.GetFFmpegPath()

	log.Printf("Audio: Starting FFmpeg with: %s %v", ffmpegPath, args)
	a.cmd = exec.Command(ffmpegPath, args...)

	stdout, err := a.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("audio: failed to get stdout pipe: %w", err)
	}
	a.pipeReader = stdout

	if err := a.cmd.Start(); err != nil {
		return fmt.Errorf("audio: failed to start FFmpeg: %w", err)
	}

	a.isCapturing = true
	go a.readAudioChunks()

	log.Printf("Audio: Capturing from device: %s", device)
	return nil
}

func (a *AudioCapture) detectAudioDevice() string {
	switch runtime.GOOS {
	case "windows":
		return a.detectWindowsAudioDevice()
	case "darwin":
		return ":0"
	case "linux":
		// Check if arecord sees any capture devices
		cmd := exec.Command("arecord", "-l")
		if err := cmd.Run(); err == nil {
			return "default"
		}
		return ""
	default:
		return ""
	}
}

func (a *AudioCapture) detectWindowsAudioDevice() string {
	ffmpegPath := media.GetFFmpegPath()
	cmd := exec.Command(ffmpegPath, "-list_devices", "true", "-f", "dshow", "-i", "dummy")
	output, _ := cmd.CombinedOutput()

	lines := strings.Split(string(output), "\n")
	inAudioSection := false

	for _, line := range lines {
		if strings.Contains(line, "DirectShow audio devices") {
			inAudioSection = true
			continue
		}
		if inAudioSection && strings.Contains(line, "DirectShow video devices") {
			break
		}
		if inAudioSection && strings.Contains(line, "\"") {
			start := strings.Index(line, "\"")
			end := strings.LastIndex(line, "\"")
			if start != -1 && end != -1 && start < end {
				name := line[start+1 : end]
				if name != "" {
					log.Printf("Audio: Found device: %s", name)
					return name
				}
			}
		}
	}
	return ""
}

func (a *AudioCapture) buildFFmpegArgs(device string) []string {
	switch runtime.GOOS {
	case "windows":
		return []string{
			"-f", "dshow",
			"-i", fmt.Sprintf("audio=%s", device),
			"-ac", fmt.Sprintf("%d", a.channels),
			"-ar", fmt.Sprintf("%d", a.sampleRate),
			"-f", "s16le",
			"-acodec", "pcm_s16le",
			"pipe:1",
		}
	case "darwin":
		return []string{
			"-f", "avfoundation",
			"-i", device,
			"-ac", fmt.Sprintf("%d", a.channels),
			"-ar", fmt.Sprintf("%d", a.sampleRate),
			"-f", "s16le",
			"-acodec", "pcm_s16le",
			"pipe:1",
		}
	default: // linux
		return []string{
			"-f", "alsa",
			"-i", device,
			"-ac", fmt.Sprintf("%d", a.channels),
			"-ar", fmt.Sprintf("%d", a.sampleRate),
			"-f", "s16le",
			"-acodec", "pcm_s16le",
			"pipe:1",
		}
	}
}

// readAudioChunks reads fixed-size PCM chunks (~33ms each) from the FFmpeg pipe.
// 44100 Hz * 2 channels * 2 bytes/sample * 33ms/1000 ≈ 5,821 bytes per chunk.
func (a *AudioCapture) readAudioChunks() {
	chunkSize := a.sampleRate * a.channels * 2 * 33 / 1000 // ~5,821 bytes
	buf := make([]byte, chunkSize)

	for {
		select {
		case <-a.stopChan:
			return
		default:
		}

		_, err := io.ReadFull(a.pipeReader, buf)
		if err != nil {
			select {
			case <-a.stopChan:
				return
			default:
			}
			log.Printf("Audio: read error: %v", err)
			return
		}

		chunk := make([]byte, len(buf))
		copy(chunk, buf)

		a.chunkMu.Lock()
		a.latestChunk = chunk
		a.chunkMu.Unlock()
	}
}

func (a *AudioCapture) CaptureAudio() ([]byte, error) {
	if !a.isCapturing {
		return nil, fmt.Errorf("not capturing audio")
	}

	a.chunkMu.Lock()
	chunk := a.latestChunk
	a.chunkMu.Unlock()

	if chunk != nil {
		return chunk, nil
	}

	return a.generateSilence(), nil
}

func (a *AudioCapture) Stop() {
	if !a.isCapturing {
		return
	}

	a.isCapturing = false
	close(a.stopChan)

	if a.cmd != nil && a.cmd.Process != nil {
		a.cmd.Process.Kill()
		a.cmd.Wait()
	}

	if a.pipeReader != nil {
		a.pipeReader.Close()
	}
}

func (a *AudioCapture) generateSilence() []byte {
	// Generate 33ms of silence (for 30fps sync)
	samples := a.sampleRate * a.channels * 33 / 1000
	return make([]byte, samples*2) // 16-bit samples, all zeros
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

	if len(audioData) > 0 {
		fmt.Printf("Audio: Playing %d bytes of audio data\n", len(audioData))
	}

	return nil
}
