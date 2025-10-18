package streaming

import (
	"context"
	"fmt"
	"time"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/sirupsen/logrus"
	"meshlink/internal/media"
)

type Viewer struct {
	subscription   *pubsub.Subscription
	logger         *logrus.Logger
	ctx            context.Context
	onData         func([]byte)
	onFrameReceived func(*media.DecodedFrame)
	isViewing      bool
	framesReceived uint64
	bytesReceived  uint64
	lastFrameTime  time.Time
	stopChan       chan struct{}
	decoder        *media.H264Decoder
}

func NewViewer(ctx context.Context, ps *pubsub.PubSub, onData func([]byte)) (*Viewer, error) {
	topic, err := ps.Join(StreamTopic)
	if err != nil {
		return nil, fmt.Errorf("failed to join topic: %w", err)
	}

	sub, err := topic.Subscribe()
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe: %w", err)
	}

	// Initialize decoder
	decoder := media.NewH264Decoder()

	return &Viewer{
		subscription: sub,
		logger:       logrus.New(),
		ctx:          ctx,
		onData:       onData,
		stopChan:     make(chan struct{}),
		decoder:      decoder,
	}, nil
}

func (v *Viewer) StartViewing() error {
	if v.isViewing {
		return fmt.Errorf("already viewing")
	}
	
	v.logger.Info("Starting stream viewer...")
	
	// Start decoder
	if err := v.decoder.Start(); err != nil {
		return fmt.Errorf("failed to start decoder: %w", err)
	}
	
	v.isViewing = true
	v.framesReceived = 0
	v.bytesReceived = 0
	v.lastFrameTime = time.Now()
	
	go v.receiveLoop()
	return nil
}

func (v *Viewer) receiveLoop() {
	for {
		select {
		case <-v.ctx.Done():
			v.logger.Info("Viewer stopped - context cancelled")
			return
		case <-v.stopChan:
			v.logger.Info("Viewer stopped - stop signal received")
			return
		default:
			msg, err := v.subscription.Next(v.ctx)
			if err != nil {
				if v.ctx.Err() != nil {
					return // Context cancelled
				}
				v.logger.Errorf("Failed to receive message: %v", err)
				continue
			}

			// Process received frame
			v.processFrame(msg.Data)
		}
	}
}

func (v *Viewer) processFrame(data []byte) {
	// Update statistics
	v.framesReceived++
	v.bytesReceived += uint64(len(data))
	v.lastFrameTime = time.Now()
	
	// Decode frame using H.264 decoder
	decodedFrame, err := v.decoder.DecodeFrame(data)
	if err != nil {
		v.logger.Errorf("Failed to decode frame: %v", err)
		// Still call legacy callback with raw data
		if v.onData != nil {
			v.onData(data)
		}
		return
	}
	
	// Call frame callback if set
	if v.onFrameReceived != nil {
		v.onFrameReceived(decodedFrame)
	}
	
	// Call legacy data callback
	if v.onData != nil {
		v.onData(data)
	}
	
	// Log statistics periodically
	if v.framesReceived%30 == 0 { // Every second at 30fps
		v.logger.Infof("Received %d frames, %d bytes total, quality: %s, last frame: %v", 
			v.framesReceived, v.bytesReceived, decodedFrame.GetQuality(), v.lastFrameTime.Format("15:04:05.000"))
	}
}

func (v *Viewer) SetOnFrameReceived(callback func(*media.DecodedFrame)) {
	v.onFrameReceived = callback
}

func (v *Viewer) Stop() {
	if !v.isViewing {
		return
	}
	
	v.logger.Info("Stopping stream viewer...")
	v.isViewing = false
	
	// Stop decoder
	v.decoder.Stop()
	
	// Signal stop to receive loop
	select {
	case v.stopChan <- struct{}{}:
	default:
	}
	
	v.subscription.Cancel()
}

func (v *Viewer) GetStats() (framesReceived uint64, bytesReceived uint64, isViewing bool, lastFrameTime time.Time) {
	return v.framesReceived, v.bytesReceived, v.isViewing, v.lastFrameTime
}

func (v *Viewer) GetFrameRate() float64 {
	if v.framesReceived == 0 {
		return 0
	}
	// Simple frame rate calculation
	duration := time.Since(v.lastFrameTime.Add(-time.Duration(v.framesReceived) * 33 * time.Millisecond))
	return float64(v.framesReceived) / duration.Seconds()
}