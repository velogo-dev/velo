package main

import (
	"log"

	"github.com/velogo-dev/velo/pkg/app"
	"github.com/velogo-dev/velo/pkg/config"
)

func main() {
	// Initialize configuration
	cfg := config.New()

	// Create and run the application
	app := app.New(cfg)

	if err := app.Run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
