package web

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WebServer struct {
	port     int
	clients  map[*websocket.Conn]bool
	broadcast chan []byte
	mutex    sync.RWMutex
	upgrader websocket.Upgrader
}

func NewWebServer(port int) *WebServer {
	return &WebServer{
		port:      port,
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for local network
			},
		},
	}
}

func (ws *WebServer) Start() error {
	// Serve static files
	http.HandleFunc("/", ws.serveHome)
	http.HandleFunc("/ws", ws.handleWebSocket)

	// Start broadcast handler
	go ws.handleBroadcast()

	fmt.Printf("Web viewer available at: http://localhost:%d\n", ws.port)
	fmt.Printf("Share this URL with congregation on same WiFi\n")
	return http.ListenAndServe(fmt.Sprintf(":%d", ws.port), nil)
}

func (ws *WebServer) serveHome(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>MeshLink Church Viewer</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; background: #1a1a1a; color: white; }
        .container { max-width: 1200px; margin: 0 auto; }
        .video-area { background: #333; border-radius: 10px; padding: 20px; margin: 20px 0; min-height: 400px; text-align: center; }
        .stats { background: #444; padding: 15px; border-radius: 5px; margin: 10px 0; }
        .status { padding: 10px; border-radius: 5px; margin: 10px 0; text-align: center; font-weight: bold; }
        .connected { background: #2d5a2d; }
        .disconnected { background: #5a2d2d; }
        .frame-display { font-family: monospace; background: #222; padding: 20px; border-radius: 5px; font-size: 18px; line-height: 1.6; }
        .live-indicator { animation: pulse 2s infinite; }
        @keyframes pulse { 0% { opacity: 1; } 50% { opacity: 0.5; } 100% { opacity: 1; } }
    </style>
</head>
<body>
    <div class="container">
        <h1>MeshLink Church Viewer</h1>
        <div id="status" class="status disconnected">Connecting to stream...</div>

        <div class="video-area">
            <h3>Live Church Stream</h3>
            <canvas id="video-canvas" width="800" height="450" style="background: #000; border-radius: 5px; border: 2px solid #555;">
                Your browser doesn't support HTML5 canvas.
            </canvas>
            <div id="video-display" class="frame-display" style="margin-top: 20px;">
                Waiting for video stream...<br>
                Make sure broadcaster is running and you're on the same WiFi network.
            </div>
        </div>

        <div class="stats">
            <h3>Stream Statistics</h3>
            <div id="stats">Waiting for connection...</div>
        </div>
    </div>

    <script>
        const ws = new WebSocket('ws://' + window.location.host + '/ws');
        const status = document.getElementById('status');
        const videoDisplay = document.getElementById('video-display');
        const canvas = document.getElementById('video-canvas');
        const ctx = canvas.getContext('2d');
        const stats = document.getElementById('stats');

        let frameCount = 0;
        let totalBytes = 0;
        let startTime = Date.now();

        // Draw test pattern on canvas
        function drawTestPattern() {
            const gradient = ctx.createLinearGradient(0, 0, canvas.width, canvas.height);
            gradient.addColorStop(0, '#1a1a1a');
            gradient.addColorStop(0.5, '#333');
            gradient.addColorStop(1, '#1a1a1a');

            ctx.fillStyle = gradient;
            ctx.fillRect(0, 0, canvas.width, canvas.height);

            // Draw animated pattern
            const time = Date.now() / 1000;
            ctx.fillStyle = '#4CAF50';
            for (let i = 0; i < 10; i++) {
                const x = (canvas.width / 10) * i + Math.sin(time + i) * 20;
                const y = canvas.height / 2 + Math.cos(time + i) * 50;
                ctx.fillRect(x, y, 20, 20);
            }

            // Draw frame info
            ctx.fillStyle = 'white';
            ctx.font = '24px Arial';
            ctx.textAlign = 'center';
            ctx.fillText('LIVE CHURCH STREAM', canvas.width / 2, 50);
            ctx.fillText('Frame #' + frameCount, canvas.width / 2, canvas.height - 50);
        }

        ws.onopen = function() {
            status.textContent = 'LIVE - Connected to Church Stream';
            status.className = 'status connected';
            startTime = Date.now();

            // Draw initial test pattern
            drawTestPattern();
        };

        ws.onmessage = function(event) {
            const data = JSON.parse(event.data);

            if (data.type === 'frame') {
                frameCount++;
                totalBytes += data.size;

                // Update canvas with frame data
                if (data.imageData && data.hasImage) {
                    const img = new Image();
                    img.onload = function() {
                        ctx.clearRect(0, 0, canvas.width, canvas.height);

                        // Draw image scaled to fit canvas
                        ctx.drawImage(img, 0, 0, canvas.width, canvas.height);

                        // Overlay frame info
                        ctx.fillStyle = 'rgba(0,0,0,0.8)';
                        ctx.fillRect(0, 0, canvas.width, 80);
                        ctx.fillStyle = 'white';
                        ctx.font = '24px Arial';
                        ctx.textAlign = 'center';
                        ctx.fillText('LIVE CHURCH STREAM', canvas.width / 2, 30);
                        ctx.fillText('Frame #' + data.frame, canvas.width / 2, 60);
                    };
                    img.onerror = function() {
                        drawTestPattern();
                    };
                    img.src = data.imageData;
                } else {
                    drawTestPattern();
                }

                // Create frame info display
                const elapsed = (Date.now() - startTime) / 1000;
                const fps = (frameCount / elapsed).toFixed(1);

                const frameInfo = '<div class="live-indicator">LIVE CHURCH STREAM</div><br><br>' +
                                '<div style="font-size: 24px; margin: 20px 0;">STREAMING NOW</div>' +
                                '<div style="font-size: 18px; line-height: 2;">' +
                                'Video Quality: ' + (data.quality || '720p HD') + '<br>' +
                                'Audio: Stereo PCM<br>' +
                                'Frame #' + frameCount + ' of live stream<br>' +
                                'Data: ' + (data.size / 1024).toFixed(1) + ' KB per frame<br>' +
                                'Time: ' + new Date().toLocaleTimeString() + '<br>' +
                                'Total: ' + (totalBytes / 1024 / 1024).toFixed(2) + ' MB received<br>' +
                                'Frame Rate: ' + fps + ' FPS<br>' +
                                '</div><br>' +
                                '<div style="color: #4CAF50; font-weight: bold;">Connected to Church Stream</div>';

                videoDisplay.innerHTML = frameInfo;

                // Update stats
                const bitrate = (totalBytes * 8 / 1024 / 1024 / elapsed).toFixed(1);
                const duration = Math.floor(elapsed / 60) + ':' + String(Math.floor(elapsed % 60)).padStart(2, '0');

                stats.innerHTML = 'Frames: ' + frameCount +
                                ' | Data: ' + (totalBytes / 1024 / 1024).toFixed(2) + ' MB' +
                                ' | Bitrate: ' + bitrate + ' Mbps' +
                                ' | Duration: ' + duration;
            }
        };

        ws.onclose = function() {
            status.textContent = 'Disconnected from stream';
            status.className = 'status disconnected';
            videoDisplay.innerHTML = 'Connection lost<br>Trying to reconnect...';
        };

        ws.onerror = function() {
            status.textContent = 'Connection error';
            status.className = 'status disconnected';
        };
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func (ws *WebServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	ws.mutex.Lock()
	ws.clients[conn] = true
	ws.mutex.Unlock()

	fmt.Printf("New viewer connected (Total: %d)\n", len(ws.clients))

	// Remove client when done
	defer func() {
		ws.mutex.Lock()
		delete(ws.clients, conn)
		ws.mutex.Unlock()
		fmt.Printf("Viewer disconnected (Total: %d)\n", len(ws.clients))
	}()

	// Keep connection alive
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (ws *WebServer) handleBroadcast() {
	for {
		data := <-ws.broadcast

		ws.mutex.RLock()
		for client := range ws.clients {
			err := client.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				client.Close()
				delete(ws.clients, client)
			}
		}
		ws.mutex.RUnlock()
	}
}

func (ws *WebServer) BroadcastFrame(frameData []byte, frameNum uint64, quality string) {
	// Extract raw frame data from the metadata wrapper
	rawData := ws.extractFrameData(frameData)

	var imageData string
	hasImage := false

	// Check for JPEG SOI marker (0xFF 0xD8)
	if len(rawData) >= 2 && rawData[0] == 0xFF && rawData[1] == 0xD8 {
		imageData = fmt.Sprintf("data:image/jpeg;base64,%s", base64.StdEncoding.EncodeToString(rawData))
		hasImage = true
	} else {
		imageData = ws.generateTestImage(int(frameNum))
		hasImage = imageData != ""
	}

	message := map[string]interface{}{
		"type":      "frame",
		"frame":     frameNum,
		"size":      len(frameData),
		"quality":   quality,
		"imageData": imageData,
		"hasImage":  hasImage,
	}

	jsonData, _ := json.Marshal(message)

	select {
	case ws.broadcast <- jsonData:
	default:
		// Don't block if channel is full
	}
}

// extractFrameData extracts the raw payload from the [4-byte-length][metadata][data] wrapper.
func (ws *WebServer) extractFrameData(frameData []byte) []byte {
	if len(frameData) < 4 {
		return nil
	}

	metadataLen := int(frameData[0])<<24 | int(frameData[1])<<16 | int(frameData[2])<<8 | int(frameData[3])
	start := 4 + metadataLen
	if start >= len(frameData) {
		return nil
	}

	return frameData[start:]
}

func (ws *WebServer) generateTestImage(frame int) string {
	// Generate a visible test pattern that changes with each frame
	color := (frame * 5) % 360

	svg := fmt.Sprintf(`<svg width="320" height="240" xmlns="http://www.w3.org/2000/svg">
		<rect width="100%%" height="100%%" fill="hsl(%d,50%%,30%%)"/>
		<circle cx="160" cy="120" r="50" fill="hsl(%d,80%%,60%%)"/>
		<text x="160" y="130" text-anchor="middle" fill="white" font-size="16">Frame %d</text>
	</svg>`, color, (color+180)%360, frame)

	encoded := base64.StdEncoding.EncodeToString([]byte(svg))
	return fmt.Sprintf("data:image/svg+xml;base64,%s", encoded)
}

func (ws *WebServer) GetViewerCount() int {
	ws.mutex.RLock()
	defer ws.mutex.RUnlock()
	return len(ws.clients)
}
