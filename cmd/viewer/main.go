package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"meshlink/internal/config"
	"meshlink/internal/p2p"
	"meshlink/internal/ui"
	"meshlink/pkg/streaming"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load configuration (uses defaults if config.json missing)
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Printf("Config warning: %v, using defaults", err)
		cfg = config.DefaultConfig()
	}

	// Initialize P2P node
	node, err := p2p.NewNode(ctx)
	if err != nil {
		log.Fatalf("Failed to create P2P node: %v", err)
	}
	defer node.Close()

	log.Printf("Viewer started with ID: %s", node.Host.ID())
	log.Printf("Expecting quality: %s, resolution: %s", cfg.Media.VideoCodec, cfg.Media.Resolution)

	// Check if running in headless mode
	if os.Getenv("DISPLAY_MODE") == "headless" {
		// Headless mode for Docker
		headlessUI := ui.NewHeadlessUI("Viewer")
		headlessUI.Start()
		
		// Auto-connect to stream
		viewer, err := streaming.NewViewer(ctx, node.PubSub, func(data []byte) {
			log.Printf("Received stream data: %d bytes", len(data))
		})
		if err != nil {
			log.Fatalf("Failed to create viewer: %v", err)
		}
		
		if err := viewer.StartViewing(); err != nil {
			log.Fatalf("Failed to start viewing: %v", err)
		}
		
		// Handle graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down viewer...")
		viewer.Stop()
		headlessUI.Stop()
	} else {
		// GUI mode
		viewerUI := ui.NewViewerUI()
		
		var viewer *streaming.Viewer
		
		viewerUI.SetOnConnect(func() error {
			if viewer == nil {
				v, err := streaming.NewViewer(ctx, node.PubSub, func(data []byte) {
					viewerUI.UpdateVideoFrame(data)
				})
				if err != nil {
					return err
				}
				viewer = v
			}
			return viewer.StartViewing()
		})
		
		viewerUI.SetOnDisconnect(func() {
			if viewer != nil {
				viewer.Stop()
			}
		})

		// Handle graceful shutdown
		go func() {
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
			<-sigChan
			log.Println("Shutting down viewer...")
			if viewer != nil {
				viewer.Stop()
			}
			cancel()
		}()

		// Run UI (blocking)
		viewerUI.Run()
	}
}