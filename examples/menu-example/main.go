package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/velogo-dev/velo/internal/desktop"
)

func main() {
	// Parse command line flags
	appName := flag.String("name", "Velo Menu Example", "Application name")
	width := flag.Int("width", 1024, "Window width")
	height := flag.Int("height", 768, "Window height")
	debug := flag.Bool("debug", false, "Enable debug mode")
	staticDir := flag.String("dir", "frontend/dist", "Static files directory")
	useHTMLMenu := flag.Bool("html-menu", false, "Use HTML-based menu instead of native")
	flag.Parse()

	// Ensure static directory is absolute
	absStaticDir, err := filepath.Abs(*staticDir)
	if err != nil {
		log.Fatalf("Failed to get absolute path for static directory: %v", err)
	}

	// Check if we have the menu-example.html file
	menuExamplePath := filepath.Join(absStaticDir, "menu-example.html")
	customHTMLExists := fileExists(menuExamplePath)

	// Create a menu bar
	menuBar := desktop.NewMenuBar()

	// Add File menu
	menuBar.AddStandardFileMenu(
		// New
		func() {
			log.Println("New menu item clicked")
		},
		// Open
		func() {
			log.Println("Open menu item clicked")
		},
		// Save
		func() {
			log.Println("Save menu item clicked")
		},
		// Exit
		func() {
			log.Println("Exit menu item clicked")
			os.Exit(0)
		},
	)

	// Add Edit menu
	menuBar.AddStandardEditMenu()

	// Add custom Help menu
	helpMenu := desktop.MenuItem{
		ID:    "help",
		Label: "Help",
		SubItems: []desktop.MenuItem{
			{
				ID:       "about",
				Label:    "About",
				Shortcut: "F1",
				Action: func() {
					log.Println("About menu item clicked")
					// Show an about dialog using JavaScript
					javascript := `
						alert("Velo Menu Example\\n\\nA demonstration of menu functionality in Velo desktop applications.");
					`
					// This will be executed via the webview
					menuBar.ExecuteJS(javascript)
				},
			},
			{
				ID:       "documentation",
				Label:    "Documentation",
				Shortcut: "F2",
				Action: func() {
					log.Println("Documentation menu item clicked")
					javascript := `
						window.open("https://github.com/velogo-dev/velo", "_blank");
					`
					menuBar.ExecuteJS(javascript)
				},
			},
		},
	}
	menuBar.AddMenuItem(helpMenu)

	// Create custom tool menu
	toolMenu := desktop.MenuItem{
		ID:    "tools",
		Label: "Tools",
		SubItems: []desktop.MenuItem{
			{
				ID:       "console",
				Label:    "Show Console",
				Shortcut: "Ctrl+Shift+I",
				Action: func() {
					log.Println("Show console clicked")
					javascript := `
						console.log("Console activated from menu");
						alert("Check browser console for logs");
					`
					menuBar.ExecuteJS(javascript)
				},
			},
			{
				ID:       "refresh",
				Label:    "Refresh Page",
				Shortcut: "F5",
				Action: func() {
					log.Println("Refresh page clicked")
					javascript := `
						window.location.reload();
					`
					menuBar.ExecuteJS(javascript)
				},
			},
		},
	}
	menuBar.AddMenuItem(toolMenu)

	// Configure the app builder
	app := desktop.NewAppBuilder().
		WithTitle(*appName).
		WithSize(*width, *height).
		WithDebug(*debug).
		WithStaticDir(absStaticDir).
		WithLocalServer(true).
		WithMenu(menuBar).
		WithHTMLMenu(*useHTMLMenu)

	// Add custom init script if menu-example.html exists
	if customHTMLExists {
		// Navigate to the custom HTML file
		app.WithInitScript(`
			console.log("Loading menu example page...");
			window.addEventListener('DOMContentLoaded', function() {
				console.log("Menu example page loaded");
			});
		`)
	} else {
		app.WithInitScript(`
			console.log("Menu example HTML not found, using default page");
		`)
	}

	app.WithOnClose(func() {
		log.Println("Window closed, performing cleanup...")
	})

	// Run the application
	log.Printf("Starting desktop app with %s menu", map[bool]string{true: "HTML", false: "native"}[*useHTMLMenu])
	log.Printf("Using %s HTML file", map[bool]string{true: "custom", false: "default"}[customHTMLExists])

	// Custom URL to load the menu example HTML if it exists
	if customHTMLExists {
		app.WithURL(fmt.Sprintf("http://localhost:%d/menu-example.html", app.Build().ServerPort))
	}

	if err := app.Run(); err != nil {
		log.Fatalf("Failed to run app: %v", err)
	}
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
