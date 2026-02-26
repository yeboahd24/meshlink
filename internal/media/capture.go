package media

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"meshlink/pkg/media"
)

type CameraCapture struct {
	deviceID    string
	resolution  string
	fps         int
	isCapturing bool

	streamer   *media.FFmpegStreamer
	pipeReader *io.PipeReader
	pipeWriter *io.PipeWriter

	latestFrame []byte
	frameMu     sync.Mutex
	frameReady  chan struct{}
	stopChan    chan struct{}
}

func NewCameraCapture() *CameraCapture {
	return &CameraCapture{
		deviceID:   "0",
		resolution: "1280x720",
		fps:        30,
	}
}

func (c *CameraCapture) Start() error {
	if c.isCapturing {
		return fmt.Errorf("already capturing")
	}

	c.frameReady = make(chan struct{}, 1)
	c.stopChan = make(chan struct{})

	// Create pipe: FFmpegStreamer writes MJPEG to pipeWriter, we read from pipeReader
	c.pipeReader, c.pipeWriter = io.Pipe()

	c.streamer = media.NewFFmpegStreamer(c.pipeWriter)
	if err := c.streamer.Start(); err != nil {
		c.pipeReader.Close()
		c.pipeWriter.Close()
		return fmt.Errorf("failed to start FFmpeg streamer: %w", err)
	}

	c.isCapturing = true

	// Background goroutine to parse MJPEG frames from the pipe
	go c.readFrames()

	return nil
}

// readFrames continuously reads MJPEG frames from the pipe by scanning for
// JPEG SOI (0xFF 0xD8) and EOI (0xFF 0xD9) markers.
func (c *CameraCapture) readFrames() {
	reader := bufio.NewReaderSize(c.pipeReader, 256*1024)

	for {
		select {
		case <-c.stopChan:
			return
		default:
		}

		frame, err := c.readOneJPEG(reader)
		if err != nil {
			select {
			case <-c.stopChan:
				return
			default:
			}
			log.Printf("Error reading MJPEG frame: %v", err)
			return
		}

		if len(frame) > 0 {
			c.frameMu.Lock()
			c.latestFrame = frame
			c.frameMu.Unlock()

			// Signal that first frame is ready (non-blocking)
			select {
			case c.frameReady <- struct{}{}:
			default:
			}
		}
	}
}

// readOneJPEG scans the reader for a complete JPEG image delimited by SOI and EOI markers.
func (c *CameraCapture) readOneJPEG(reader *bufio.Reader) ([]byte, error) {
	// Scan for SOI marker (0xFF 0xD8)
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		if b != 0xFF {
			continue
		}
		b2, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		if b2 == 0xD8 {
			break // Found SOI
		}
	}

	// We found SOI; accumulate the JPEG data
	buf := []byte{0xFF, 0xD8}

	for {
		b, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		buf = append(buf, b)

		if b == 0xFF {
			b2, err := reader.ReadByte()
			if err != nil {
				return nil, err
			}
			buf = append(buf, b2)

			if b2 == 0xD9 {
				// Found EOI — complete JPEG frame
				return buf, nil
			}
		}
	}
}

func (c *CameraCapture) CaptureFrame() ([]byte, error) {
	if !c.isCapturing {
		return nil, fmt.Errorf("not capturing")
	}

	// Check if we already have a frame
	c.frameMu.Lock()
	frame := c.latestFrame
	c.frameMu.Unlock()

	if frame != nil {
		return frame, nil
	}

	// First call — wait up to 500ms for the first frame
	select {
	case <-c.frameReady:
	case <-time.After(500 * time.Millisecond):
	}

	c.frameMu.Lock()
	frame = c.latestFrame
	c.frameMu.Unlock()

	if frame != nil {
		return frame, nil
	}

	// No frame yet, return fallback
	return c.generateFallbackFrame(), nil
}

func (c *CameraCapture) Stop() {
	if !c.isCapturing {
		return
	}

	c.isCapturing = false

	// Signal the reader goroutine to stop
	close(c.stopChan)

	// Stop the FFmpeg streamer (kills the process, closes pipeWriter)
	if c.streamer != nil {
		c.streamer.Stop()
	}

	// Close the pipe reader to unblock any pending reads
	if c.pipeReader != nil {
		c.pipeReader.Close()
	}
}

func (c *CameraCapture) generateFallbackFrame() []byte {
	// Return a placeholder; the web server will display a test pattern
	// if this isn't a valid JPEG.
	return []byte("no-camera")
}
