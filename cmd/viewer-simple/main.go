package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"meshlink/internal/p2p"
	"meshlink/pkg/streaming"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Printf("📺 MeshLink Viewer (Simple)\n")
	fmt.Printf("🔍 Discovering streams...\n")

	node, err := p2p.NewNode(ctx)
	if err != nil {
		log.Fatalf("Failed to create P2P node: %v", err)
	}
	defer node.Close()

	viewer, err := streaming.NewViewer(ctx, node.PubSub, func(data []byte) {
		fmt.Printf("📥 Received H.264 frame: %d bytes\n", len(data))
	})
	if err != nil {
		log.Fatalf("Failed to create viewer: %v", err)
	}

	fmt.Printf("Node ID: %s\n", node.Host.ID())
	fmt.Printf("Starting viewer...\n")

	if err := viewer.StartViewing(); err != nil {
		log.Fatalf("Failed to start viewing: %v", err)
	}

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				frames, bytes, viewing, lastFrame := viewer.GetStats()
				fmt.Printf("📊 Frames: %d | Viewing: %t | Bytes: %d | Last: %s\n", 
					frames, viewing, bytes, lastFrame.Format("15:04:05"))
			}
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("Shutting down...")
	viewer.Stop()
}