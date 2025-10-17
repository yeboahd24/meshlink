package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	app := "MeshLink"
	if len(os.Args) > 1 {
		app = os.Args[1]
	}
	
	fmt.Printf("[%s] Starting in demo mode...\n", app)
	
	for {
		fmt.Printf("[%s] Running (PID: %d)\n", app, os.Getpid())
		time.Sleep(30 * time.Second)
	}
}