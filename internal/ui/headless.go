package ui

import (
	"fmt"
	"os"
	"time"
)

type HeadlessUI struct {
	name      string
	isRunning bool
}

func NewHeadlessUI(name string) *HeadlessUI {
	return &HeadlessUI{name: name}
}

func (h *HeadlessUI) Start() {
	fmt.Printf("[%s] Starting in headless mode...\n", h.name)
	h.isRunning = true
	
	// Keep running and show status
	go h.statusLoop()
}

func (h *HeadlessUI) statusLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for h.isRunning {
		select {
		case <-ticker.C:
			fmt.Printf("[%s] Status: Running (PID: %d)\n", h.name, os.Getpid())
		}
	}
}

func (h *HeadlessUI) Stop() {
	h.isRunning = false
	fmt.Printf("[%s] Stopping...\n", h.name)
}

func (h *HeadlessUI) Wait() {
	// Keep process alive
	select {}
}