package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/meshlink/church-streaming/internal/config"
	"github.com/meshlink/church-streaming/internal/p2p"
	"github.com/meshlink/church-streaming/internal/ui"
	"github.com/meshlink/church-streaming/pkg/streaming"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load configuration
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize P2P node with config
	node, err := p2p.NewNode(ctx)
	if err != nil {
		log.Fatalf("Failed to create P2P node: %v", err)
	}
	defer node.Close()

	log.Printf("Broadcaster started with ID: %s", node.Host.ID())
	log.Printf("Using quality: %s, bitrate: %d", cfg.Media.VideoCodec, cfg.Media.Bitrate)

	// Initialize broadcaster with config
	broadcaster, err := streaming.NewBroadcasterWithConfig(ctx, node.PubSub, cfg)
	if err != nil {
		log.Fatalf("Failed to create broadcaster: %v", err)
	}

	// Check if running in headless mode
	if os.Getenv("DISPLAY_MODE") == "headless" {
		// Headless mode for Docker
		headlessUI := ui.NewHeadlessUI("Broadcaster")
		headlessUI.Start()
		
		// Auto-start broadcasting
		if err := broadcaster.StartStreaming(); err != nil {
			log.Fatalf("Failed to start streaming: %v", err)
		}
		
		// Handle graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down broadcaster...")
		broadcaster.Stop()
		headlessUI.Stop()
	} else {
		// GUI mode
		broadcasterUI := ui.NewBroadcasterUI()
		broadcasterUI.SetCallbacks(
			func() error {
				return broadcaster.StartStreaming()
			},
			func() {
				broadcaster.Stop()
			},
		)
		
		// Set quality change callback
		broadcasterUI.SetOnQualityChange(func(quality string) error {
			return broadcaster.SetQuality(quality)
		})
		
		// Connect real statistics
		broadcasterUI.SetStatsCallbacks(
			func() (uint64, uint64, bool) {
				return broadcaster.GetStats()
			},
			func() int {
				return broadcaster.GetViewerCount()
			},
		)

		// Handle graceful shutdown
		go func() {
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
			<-sigChan
			log.Println("Shutting down broadcaster...")
			cancel()
		}()

		// Run UI (blocking)
		broadcasterUI.Run()
	}
}