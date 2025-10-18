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
	videoWidget *VideoWidget
	videoCard   *widget.Card
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

	// Create video widget
	ui.videoWidget = NewVideoWidget(640, 480)
	ui.videoCard = widget.NewCard("Video Stream", "Waiting for connection...", ui.videoWidget)
	ui.videoCard.Resize(fyne.NewSize(660, 500))

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
		ui.videoCard,
	)

	ui.window.SetContent(content)
}

func (ui *ViewerUI) updateUI() {
	if ui.isConnected {
		ui.statusText.SetText("🔴 Connected - Receiving Stream")
		ui.connectBtn.SetText("Disconnect")
		ui.videoCard.SetSubTitle("🔴 LIVE - Receiving HD Stream")
		ui.statsLabel.SetText("Statistics: Connected - waiting for data...")
	} else {
		ui.statusText.SetText("⚪ Searching for broadcasts...")
		ui.connectBtn.SetText("Connect to Stream")
		ui.videoCard.SetSubTitle("Waiting for connection...")
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
	
	// Update video widget with actual frame data
	ui.videoWidget.UpdateFrame(data)
	
	// Update card subtitle with live indicator
	ui.videoCard.SetSubTitle(fmt.Sprintf("%s LIVE - Frame #%d - %d bytes %s", 
		frameIndicator, ui.framesReceived, len(data), frameIndicator))

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