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
	"github.com/meshlink/church-streaming/internal/web"
	"github.com/meshlink/church-streaming/pkg/streaming"
)

func main() {
	var (
		configPath = flag.String("config", "", "Path to config file (optional)")
		webPort    = flag.Int("web", 8080, "Web server port for browser viewers")
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

	fmt.Printf("🎥 MeshLink Church Broadcaster (Web-Enabled)\n")
	fmt.Printf("🌐 Web viewer: http://localhost:%d\n", *webPort)
	fmt.Printf("📱 Share with congregation on same WiFi\n")
	fmt.Printf("Press Ctrl+C to stop\n\n")

	// Initialize P2P node
	node, err := p2p.NewNode(ctx)
	if err != nil {
		log.Fatalf("Failed to create P2P node: %v", err)
	}
	defer node.Close()

	// Initialize web server
	webServer := web.NewWebServer(*webPort)
	go func() {
		fmt.Printf("🌐 Starting web server on port %d...\n", *webPort)
		if err := webServer.Start(); err != nil {
			log.Printf("❌ Web server error: %v", err)
		} else {
			fmt.Printf("✅ Web server started successfully\n")
		}
	}()
	
	// Give web server time to start
	time.Sleep(2 * time.Second)

	// Create broadcaster
	broadcaster, err := streaming.NewBroadcasterWithConfig(ctx, node.PubSub, cfg)
	if err != nil {
		log.Fatalf("Failed to create broadcaster: %v", err)
	}

	fmt.Printf("Node ID: %s\n", node.Host.ID())
	fmt.Printf("Starting broadcast...\n\n")

	if err := broadcaster.StartStreaming(); err != nil {
		log.Fatalf("Failed to start streaming: %v", err)
	}

	// Forward P2P frames to web viewers
	go func() {
		ticker := time.NewTicker(200 * time.Millisecond) // 5 FPS for web
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// Get current P2P streaming stats
				p2pFrames, _, isStreaming := broadcaster.GetStats()
				
				if isStreaming {
					// Get actual camera frame from broadcaster
					frameData := broadcaster.GetCurrentFrame()
					if frameData == nil {
						// Fallback to simulated frame
						frameData = make([]byte, 50*1024)
						for i := range frameData {
							frameData[i] = byte((p2pFrames + uint64(i)) % 256)
						}
					}
					
					// Send frame to web viewers
					webServer.BroadcastFrame(frameData, p2pFrames, "720p")
					if webServer.GetViewerCount() > 0 {
						fmt.Printf("📡 Sent frame #%d (%d bytes) to %d web viewers\n", p2pFrames, len(frameData), webServer.GetViewerCount())
					}
				}
			}
		}
	}()

	// Print statistics
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				p2pFrames, p2pBytes, streaming := broadcaster.GetStats()
				p2pViewers := broadcaster.GetViewerCount()
				webViewers := webServer.GetViewerCount()
				totalViewers := p2pViewers + webViewers
				
				bitrate := float64(p2pBytes*8) / (1024*1024)
				
				fmt.Printf("📊 Streaming: %t | 👥 Viewers: %d (P2P: %d, Web: %d) | 📈 Frames: %d | 🌐 %.1f Mbps\n", 
					streaming, totalViewers, p2pViewers, webViewers, p2pFrames, bitrate)
			}
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\n🛑 Shutting down...")
	broadcaster.Stop()
}