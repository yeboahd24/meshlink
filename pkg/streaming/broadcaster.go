package streaming

import (
	"context"
	"fmt"
	"time"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/sirupsen/logrus"
	"github.com/meshlink/church-streaming/internal/config"
	"github.com/meshlink/church-streaming/internal/media"
)

const StreamTopic = "meshlink/church/stream"

type Broadcaster struct {
	topic       *pubsub.Topic
	logger      *logrus.Logger
	ctx         context.Context
	isStreaming bool
	viewerCount int
	bytesSent   uint64
	frameCount  uint64
	stopChan    chan struct{}
	camera      *media.CameraCapture
	encoder     *media.H264Encoder
	quality     string
}

func NewBroadcaster(ctx context.Context, ps *pubsub.PubSub) (*Broadcaster, error) {
	return NewBroadcasterWithConfig(ctx, ps, nil)
}

func NewBroadcasterWithConfig(ctx context.Context, ps *pubsub.PubSub, cfg *config.Config) (*Broadcaster, error) {
	topic, err := ps.Join(StreamTopic)
	if err != nil {
		return nil, fmt.Errorf("failed to join topic: %w", err)
	}

	// Use config or defaults
	quality := "720p"
	if cfg != nil {
		switch {
		case cfg.Media.Resolution == "1920x1080":
			quality = "1080p"
		case cfg.Media.Resolution == "854x480":
			quality = "480p"
		default:
			quality = "720p"
		}
	}

	// Initialize media components with config
	camera := media.NewCameraCapture()
	encoder := media.NewH264Encoder(quality)

	return &Broadcaster{
		topic:    topic,
		logger:   logrus.New(),
		ctx:      ctx,
		stopChan: make(chan struct{}),
		camera:   camera,
		encoder:  encoder,
		quality:  quality,
	}, nil
}

type StreamFrame struct {
	FrameID   uint64    `json:"frame_id"`
	Timestamp time.Time `json:"timestamp"`
	FrameType string    `json:"frame_type"` // "video", "audio", "metadata"
	Data      []byte    `json:"data"`
	Size      int       `json:"size"`
	Quality   string    `json:"quality"`
}

func (b *Broadcaster) StartStreaming() error {
	if b.isStreaming {
		return fmt.Errorf("already streaming")
	}
	
	b.logger.Info("Starting broadcast stream...")
	
	// Start camera capture
	if err := b.camera.Start(); err != nil {
		return fmt.Errorf("failed to start camera: %w", err)
	}
	
	// Start encoder
	if err := b.encoder.Start(); err != nil {
		b.camera.Stop()
		return fmt.Errorf("failed to start encoder: %w", err)
	}
	
	b.isStreaming = true
	b.frameCount = 0
	b.bytesSent = 0
	
	// Start viewer count monitoring
	b.UpdateViewerCount()
	
	// Start streaming loop
	go b.streamLoop()
	
	return nil
}

func (b *Broadcaster) streamLoop() {
	// Production 30 FPS video streaming
	ticker := time.NewTicker(33 * time.Millisecond) // 30 FPS timing
	defer ticker.Stop()

	for {
		select {
		case <-b.ctx.Done():
			b.logger.Info("Stream stopped - context cancelled")
			return
		case <-b.stopChan:
			b.logger.Info("Stream stopped - stop signal received")
			return
		case <-ticker.C:
			if !b.isStreaming {
				return
			}
			
			// Capture frame from camera
			rawFrame, err := b.camera.CaptureFrame()
			if err != nil {
				b.logger.Errorf("Failed to capture frame: %v", err)
				continue
			}
			
			// Encode frame with H.264
			frameData, err := b.encoder.EncodeFrame(rawFrame, b.frameCount+1)
			if err != nil {
				b.logger.Errorf("Failed to encode frame: %v", err)
				continue
			}
			
			// Publish frame to P2P network
			if err := b.topic.Publish(b.ctx, frameData); err != nil {
				b.logger.Errorf("Failed to publish frame %d: %v", b.frameCount+1, err)
				continue
			}
			
			// Update statistics
			b.frameCount++
			b.bytesSent += uint64(len(frameData))
			
			if b.frameCount%30 == 0 { // Log every second
				b.logger.Infof("Streamed %d frames, %d bytes total", b.frameCount, b.bytesSent)
			}
		}
	}
}

func (b *Broadcaster) SetQuality(quality string) error {
	if b.isStreaming {
		return fmt.Errorf("cannot change quality while streaming")
	}
	
	b.quality = quality
	b.encoder = media.NewH264Encoder(quality)
	return nil
}

func (b *Broadcaster) GetQuality() string {
	return b.quality
}

func (b *Broadcaster) Stop() {
	if !b.isStreaming {
		return
	}
	
	b.logger.Info("Stopping broadcast stream...")
	b.isStreaming = false
	
	// Stop media components
	b.encoder.Stop()
	b.camera.Stop()
	
	// Signal stop to streaming loop
	select {
	case b.stopChan <- struct{}{}:
	default:
	}
	
	b.topic.Close()
}

func (b *Broadcaster) GetStats() (frameCount uint64, bytesSent uint64, isStreaming bool) {
	return b.frameCount, b.bytesSent, b.isStreaming
}

func (b *Broadcaster) GetViewerCount() int {
	// Query actual P2P network for subscriber count
	peers := b.topic.ListPeers()
	b.viewerCount = len(peers)
	return b.viewerCount
}

func (b *Broadcaster) UpdateViewerCount() {
	// Continuously update viewer count from P2P network
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		
		for b.isStreaming {
			select {
			case <-b.ctx.Done():
				return
			case <-ticker.C:
				peers := b.topic.ListPeers()
				b.viewerCount = len(peers)
				b.logger.Debugf("Updated viewer count: %d peers", b.viewerCount)
			}
		}
	}()
}