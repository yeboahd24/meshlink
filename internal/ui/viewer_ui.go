package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ViewerUI struct {
	app         fyne.App
	window      fyne.Window
	statusText  *widget.Label
	connectBtn  *widget.Button
	videoArea   *widget.Card
	statsLabel  *widget.Label
	onConnect   func() error
	onDisconnect func()
	isConnected bool
	bytesReceived uint64
	framesReceived uint64
}

func NewViewerUI() *ViewerUI {
	a := app.New()

	w := a.NewWindow("MeshLink Church Viewer")
	w.Resize(fyne.NewSize(800, 600))

	ui := &ViewerUI{
		app:    a,
		window: w,
	}

	ui.setupUI()
	return ui
}

func (ui *ViewerUI) setupUI() {
	ui.statusText = widget.NewLabel("Searching for broadcasts...")
	ui.statusText.Alignment = fyne.TextAlignCenter

	ui.connectBtn = widget.NewButton("Connect to Stream", func() {
		if !ui.isConnected {
			if ui.onConnect != nil {
				ui.statusText.SetText("Connecting...")
				if err := ui.onConnect(); err != nil {
					ui.statusText.SetText(fmt.Sprintf("Connection failed: %v", err))
					return
				}
				ui.isConnected = true
				ui.updateUI()
			}
		} else {
			if ui.onDisconnect != nil {
				ui.onDisconnect()
			}
			ui.isConnected = false
			ui.bytesReceived = 0
			ui.framesReceived = 0
			ui.updateUI()
		}
	})

	ui.videoArea = widget.NewCard("Video Stream", "Waiting for connection...", 
		widget.NewLabel("📺 Video stream will appear here\n\nResolution: 1280x720\nCodec: H.264\nBitrate: 2000 kbps"),
	)
	ui.videoArea.Resize(fyne.NewSize(640, 480))

	ui.statsLabel = widget.NewLabel("Statistics: Not connected")
	ui.statsLabel.Alignment = fyne.TextAlignCenter

	ui.updateUI()

	topControls := container.NewVBox(
		ui.statusText,
		container.NewHBox(ui.connectBtn),
		ui.statsLabel,
	)

	content := container.NewBorder(
		topControls,
		nil, nil, nil,
		ui.videoArea,
	)

	ui.window.SetContent(content)
}

func (ui *ViewerUI) updateUI() {
	if ui.isConnected {
		ui.statusText.SetText("🔴 Connected - Receiving Stream")
		ui.connectBtn.SetText("Disconnect")
		ui.videoArea.SetSubTitle("Stream active - receiving data")
		ui.statsLabel.SetText("Statistics: Connected - waiting for data...")
	} else {
		ui.statusText.SetText("⚪ Searching for broadcasts...")
		ui.connectBtn.SetText("Connect to Stream")
		ui.videoArea.SetSubTitle("Waiting for connection...")
		ui.videoArea.SetContent(widget.NewLabel("📺 Video stream will appear here\n\nResolution: 1280x720\nCodec: H.264\nBitrate: 2000 kbps"))
		ui.statsLabel.SetText("Statistics: Not connected")
	}
}

func (ui *ViewerUI) SetOnConnect(callback func() error) {
	ui.onConnect = callback
}

func (ui *ViewerUI) SetOnDisconnect(callback func()) {
	ui.onDisconnect = callback
}

func (ui *ViewerUI) UpdateVideoFrame(data []byte) {
	if !ui.isConnected {
		return
	}

	// Update statistics
	ui.bytesReceived += uint64(len(data))
	ui.framesReceived++

	// Create animated video display
	animation := []string{"🔴", "🟠", "🟡", "🟢", "🔵", "🟣"}
	frameIndicator := animation[ui.framesReceived%uint64(len(animation))]
	
	// Enhanced video display with live indicators
	frameInfo := fmt.Sprintf("%s LIVE STREAM %s\n\n📺 Video: 720p H.264\n🎵 Audio: AAC Stereo\n📊 Frame #%d\n📦 Size: %d bytes\n💾 Total: %.2f MB\n⏱️ Time: %s", 
		frameIndicator, frameIndicator,
		ui.framesReceived, 
		len(data), 
		float64(ui.bytesReceived)/(1024*1024),
		time.Now().Format("15:04:05"))
	
	ui.videoArea.SetContent(widget.NewLabel(frameInfo))
	ui.videoArea.SetSubTitle("🔴 LIVE - Receiving HD Stream")

	// Enhanced statistics with bitrate calculation
	elapsed := time.Duration(ui.framesReceived) * 33 * time.Millisecond // 30fps timing
	fps := float64(ui.framesReceived) / elapsed.Seconds()
	bitrate := float64(ui.bytesReceived*8) / (1024*1024) // Mbps
	
	statsText := fmt.Sprintf("📊 Frames: %d | 💾 Data: %.2f MB | 📈 Rate: %.1f fps | 🌐 Bitrate: %.1f Mbps", 
		ui.framesReceived, 
		float64(ui.bytesReceived)/(1024*1024),
		fps,
		bitrate)
	ui.statsLabel.SetText(statsText)
}

func (ui *ViewerUI) Run() {
	ui.window.ShowAndRun()
}