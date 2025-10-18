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

	"meshlink/internal/config"
	"meshlink/internal/p2p"
	"meshlink/pkg/streaming"
)

func main() {
	var (
		configPath = flag.String("config", "", "Path to config file (optional)")
		port       = flag.Int("port", 8080, "Port to connect to")
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

	// Create viewer with data callback
	viewer, err := streaming.NewViewer(ctx, node.PubSub, func(data []byte) {
		log.Printf("Received frame: %d bytes", len(data))
	})
	if err != nil {
		log.Fatalf("Failed to create viewer: %v", err)
	}

	fmt.Printf("Starting MeshLink Viewer (Headless)\n")
	fmt.Printf("Node ID: %s\n", node.Host.ID())
	fmt.Printf("Discovering streams...\n")
	fmt.Printf("Press Ctrl+C to stop\n\n")

	if err := viewer.StartViewing(); err != nil {
		log.Fatalf("Failed to start viewing: %v", err)
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
				frames, bytes, viewing, lastFrame := viewer.GetStats()
				bitrate := float64(bytes*8) / (1024*1024) // Mbps estimate
				fmt.Printf("Stats - Viewing: %t, Frames: %d, Bitrate: %.1f Mbps, Last: %s\n", 
					viewing, frames, bitrate, lastFrame.Format("15:04:05"))
			}
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nShutting down...")
	viewer.Stop()
}