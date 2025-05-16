package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/velogo-dev/velo/internal/desktop"
)

// Example of a function to be exposed to JavaScript
func helloFromGo(name string) string {
	return fmt.Sprintf("Hello, %s! This message was generated from Go code!", name)
}

// A more complex function to demonstrate binding capabilities
func calculateSum(numbers []int) (int, error) {
	sum := 0
	for _, num := range numbers {
		sum += num
	}
	return sum, nil
}

// Gets the current system time and returns it to JavaScript
func getCurrentTime() string {
	return time.Now().Format(time.RFC1123)
}

// Example of a function to exit the application
func exitApplication(w *desktop.AppBuilder) func() {
	return func() {
		log.Println("Exit requested from JavaScript")
		// In a real application you might want to perform cleanup here
	}
}

func main() {
	// Parse command line flags
	appName := flag.String("name", "Velo Desktop Example", "Application name")
	width := flag.Int("width", 1024, "Window width")
	height := flag.Int("height", 768, "Window height")
	debug := flag.Bool("debug", false, "Enable debug mode")
	staticDir := flag.String("dir", "frontend/dist", "Static files directory")
	url := flag.String("url", "", "URL to load (if empty, will serve local static files)")
	port := flag.Int("port", 0, "Port for local server (0 for random port)")
	flag.Parse()

	// Ensure static directory is absolute if using local static files
	var absStaticDir string
	if *url == "" {
		var err error
		absStaticDir, err = filepath.Abs(*staticDir)
		if err != nil {
			log.Fatalf("Failed to get absolute path for static directory: %v", err)
		}
	}

	// Create and configure the app builder
	builder := desktop.NewAppBuilder().
		WithTitle(*appName).
		WithSize(*width, *height).
		WithDebug(*debug)

	// Add navigation callback
	builder.WithOnNavigate(func(url string) {
		log.Printf("Navigation changed to: %s", url)
	})

	// Add close callback
	builder.WithOnClose(func() {
		log.Println("Window closed, performing cleanup...")
	})

	// Add init script for custom functionality
	builder.WithInitScript(`
		console.log("Initializing custom scripts...");
		// Add global functions that can interact with our bound Go functions
		window.getCurrentTimeFromGo = function() {
			getCurrentTime().then(function(time) {
				document.getElementById('time-display').textContent = time;
			});
		};
	`)

	// Either use a custom URL or local server
	if *url != "" {
		builder.WithURL(*url)
	} else {
		builder.WithStaticDir(absStaticDir).
			WithLocalServer(true).
			WithPort(*port)
	}

	// Add JavaScript bindings
	builder.WithBinding("helloFromGo", helloFromGo).
		WithBinding("calculateSum", calculateSum).
		WithBinding("getCurrentTime", getCurrentTime).
		WithBinding("exitApp", exitApplication(builder))

	// Run the app
	log.Printf("Starting desktop app with config: %+v", builder.Build())
	if err := builder.Run(); err != nil {
		log.Fatalf("Failed to run desktop app: %v", err)
	}
}
