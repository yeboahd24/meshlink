package main

import (
	"context"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.DefaultConfig()
	
	fmt.Printf("🎥 MeshLink Broadcaster (Simple)\n")
	fmt.Printf("📡 Starting P2P streaming...\n")

	node, err := p2p.NewNode(ctx)
	if err != nil {
		log.Fatalf("Failed to create P2P node: %v", err)
	}
	defer node.Close()

	broadcaster, err := streaming.NewBroadcasterWithConfig(ctx, node.PubSub, cfg)
	if err != nil {
		log.Fatalf("Failed to create broadcaster: %v", err)
	}

	fmt.Printf("Node ID: %s\n", node.Host.ID())
	fmt.Printf("Starting broadcast...\n")

	if err := broadcaster.StartStreaming(); err != nil {
		log.Fatalf("Failed to start streaming: %v", err)
	}

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
				fmt.Printf("📊 Frames: %d | Viewers: %d | Streaming: %t | Bytes: %d\n", 
					frames, viewers, streaming, bytes)
			}
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("Shutting down...")
	broadcaster.Stop()
}