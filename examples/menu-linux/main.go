package main

import (
	"log"
	"os"
)

func main() {
	// Create a new GTK4 application
	app := NewApp("org.example.gtk4webview")
	if app == nil {
		log.Fatal("Failed to create GTK application")
	}

	// Ensure cleanup on exit
	defer app.Cleanup()

	// Run the application
	status := app.Run()

	os.Exit(status)
}
