package ui

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type BroadcasterUI struct {
	app           fyne.App
	window        fyne.Window
	statusText    *widget.Label
	startBtn      *widget.Button
	stopBtn       *widget.Button
	previewArea   *widget.Card
	statsLabel    *widget.Label
	qualitySelect *widget.Select
	onStart       func() error
	onStop        func()
	onQualityChange func(string) error
	getStats      func() (uint64, uint64, bool)
	getViewerCount func() int
	isStreaming   bool
	startTime     time.Time
}

func NewBroadcasterUI() *BroadcasterUI {
	a := app.New()

	w := a.NewWindow("MeshLink Church Broadcaster")
	w.Resize(fyne.NewSize(400, 300))

	ui := &BroadcasterUI{
		app:    a,
		window: w,
	}

	ui.setupUI()
	return ui
}

func (ui *BroadcasterUI) setupUI() {
	ui.statusText = widget.NewLabel("Ready to broadcast")
	ui.statusText.Alignment = fyne.TextAlignCenter

	ui.startBtn = widget.NewButton("Start Broadcasting", func() {
		if ui.onStart != nil {
			ui.statusText.SetText("Starting broadcast...")
			if err := ui.onStart(); err != nil {
				ui.statusText.SetText(fmt.Sprintf("Error: %v", err))
				return
			}
			ui.isStreaming = true
			ui.updateUI()
			go ui.updateStats()
		}
	})

	ui.stopBtn = widget.NewButton("Stop Broadcasting", func() {
		if ui.onStop != nil {
			ui.onStop()
			ui.isStreaming = false
			ui.updateUI()
		}
	})

	// Quality selection
	ui.qualitySelect = widget.NewSelect([]string{"720p (2Mbps)", "1080p (4Mbps)", "480p (1Mbps)"}, func(value string) {
		// Extract quality from selection
		var quality string
		switch {
		case strings.Contains(value, "1080p"):
			quality = "1080p"
		case strings.Contains(value, "720p"):
			quality = "720p"
		case strings.Contains(value, "480p"):
			quality = "480p"
		default:
			quality = "720p"
		}
		
		// Update broadcaster quality if callback is set
		if ui.onQualityChange != nil {
			if err := ui.onQualityChange(quality); err != nil {
				ui.statusText.SetText(fmt.Sprintf("Quality change failed: %v", err))
			}
		}
	})
	ui.qualitySelect.SetSelected("720p (2Mbps)")

	// Preview area
	ui.previewArea = widget.NewCard("Camera Preview", "Camera feed will appear here", 
		widget.NewLabel("📹 Camera Preview\n\nSource: Default camera\nResolution: 1280x720\nFPS: 30\n\nIn production:\n- Live camera feed\n- Audio level meters\n- Recording controls"),
	)
	ui.previewArea.Resize(fyne.NewSize(320, 240))

	// Statistics
	ui.statsLabel = widget.NewLabel("Statistics: Not broadcasting")
	ui.statsLabel.Alignment = fyne.TextAlignCenter

	ui.updateUI()

	// Layout
	controls := container.NewVBox(
		ui.statusText,
		container.NewHBox(ui.startBtn, ui.stopBtn),
		widget.NewLabel("Quality:"),
		ui.qualitySelect,
		ui.statsLabel,
	)

	content := container.NewHBox(
		widget.NewCard("MeshLink Church Broadcaster", "", controls),
		ui.previewArea,
	)

	ui.window.SetContent(content)
}

func (ui *BroadcasterUI) updateUI() {
	if ui.isStreaming {
		ui.statusText.SetText("🔴 Broadcasting Live")
		ui.startBtn.Disable()
		ui.stopBtn.Enable()
		ui.qualitySelect.Disable()
		ui.previewArea.SetSubTitle("Live - broadcasting to network")
	} else {
		ui.statusText.SetText("⚪ Ready to broadcast")
		ui.startBtn.Enable()
		ui.stopBtn.Disable()
		ui.qualitySelect.Enable()
		ui.previewArea.SetSubTitle("Camera feed will appear here")
		ui.statsLabel.SetText("Statistics: Not broadcasting")
	}
}

func (ui *BroadcasterUI) updateStats() {
	ui.startTime = time.Now()
	
	for ui.isStreaming {
		// Get real statistics from broadcaster
		var frameCount, bytesSent uint64
		var viewerCount int
		
		if ui.getStats != nil {
			frameCount, bytesSent, _ = ui.getStats()
		}
		
		if ui.getViewerCount != nil {
			viewerCount = ui.getViewerCount()
		}
		
		// Calculate actual uptime
		uptime := time.Since(ui.startTime).Truncate(time.Second)
		
		// Calculate frame rate
		var fps float64
		if uptime.Seconds() > 0 {
			fps = float64(frameCount) / uptime.Seconds()
		}
		
		statsText := fmt.Sprintf("Viewers: %d | Sent: %.2f MB | FPS: %.1f | Uptime: %s", 
			viewerCount, 
			float64(bytesSent)/(1024*1024),
			fps,
			uptime.String())
		ui.statsLabel.SetText(statsText)
		
		time.Sleep(1 * time.Second)
	}
}

func (ui *BroadcasterUI) SetStatsCallbacks(getStats func() (uint64, uint64, bool), getViewerCount func() int) {
	ui.getStats = getStats
	ui.getViewerCount = getViewerCount
}

func (ui *BroadcasterUI) SetCallbacks(onStart func() error, onStop func()) {
	ui.onStart = onStart
	ui.onStop = onStop
}

func (ui *BroadcasterUI) SetOnQualityChange(callback func(string) error) {
	ui.onQualityChange = callback
}

func (ui *BroadcasterUI) Run() {
	ui.window.ShowAndRun()
}