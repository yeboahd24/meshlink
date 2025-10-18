package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/meshlink/church-streaming/internal/config"
	"github.com/meshlink/church-streaming/internal/p2p"
	"github.com/meshlink/church-streaming/pkg/streaming"
)

func main() {
	var (
		configPath = flag.String("config", "", "Path to config file (optional)")
		port       = flag.Int("port", 8080, "Port to listen on")
		quality    = flag.String("quality", "720p", "Video quality (480p, 720p, 1080p)")
	)
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Printf("Warning: Could not load config: %v. Using defaults.", err)
		cfg = config.DefaultConfig()
	}

	// Override with command line flags
	if *port != 8080 {
		cfg.Network.Port = *port
	}

	// Initialize P2P node
	node, err := p2p.NewNode(ctx)
	if err != nil {
		log.Fatalf("Failed to create P2P node: %v", err)
	}
	defer node.Close()

	// Create broadcaster
	broadcaster, err := streaming.NewBroadcasterWithConfig(ctx, node.PubSub, cfg)
	if err != nil {
		log.Fatalf("Failed to create broadcaster: %v", err)
	}

	fmt.Printf("Starting MeshLink Broadcaster (Headless)\n")
	fmt.Printf("Node ID: %s\n", node.Host.ID())
	fmt.Printf("Quality: %s\n", *quality)
	fmt.Printf("Press Ctrl+C to stop\n\n")

	if err := broadcaster.StartStreaming(); err != nil {
		log.Fatalf("Failed to start streaming: %v", err)
	}

	// Print statistics periodically
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				frames, bytes, streaming := broadcaster.GetStats()
				viewers := broadcaster.GetViewerCount()
				bitrate := float64(bytes*8) / (1024*1024) // Mbps estimate
				fmt.Printf("Stats - Streaming: %t, Viewers: %d, Frames: %d, Bitrate: %.1f Mbps\n", 
					streaming, viewers, frames, bitrate)
			}
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nShutting down...")
	broadcaster.Stop()
}