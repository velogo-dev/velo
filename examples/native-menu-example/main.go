package main

import (
	"flag"
	"log"
	"os"
	"runtime"

	"github.com/velogo-dev/velo/internal/desktop"
)

func main() {
	// Parse command line flags
	var staticDir string
	var debug bool
	flag.StringVar(&staticDir, "dir", "frontend/dist", "Directory with static files")
	flag.BoolVar(&debug, "debug", false, "Enable debug mode")
	flag.Parse()

	log.Println("Starting Native Menu Example")
	log.Println("Debug mode:", debug)
	log.Println("Serving static files from:", staticDir)

	// Check if static directory exists
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		log.Fatalf("Static directory '%s' does not exist", staticDir)
	}

	// Create a new native menu
	nativeMenu := desktop.NewNativeMenu()

	// Add a File menu
	nativeMenu.AddFileMenu(
		func() { log.Println("New clicked") },
		func() { log.Println("Open clicked") },
		func() { log.Println("Save clicked") },
		func() { log.Println("Exit clicked"); os.Exit(0) },
	)

	// Add an Edit menu
	nativeMenu.AddEditMenu()

	// Add a custom Help menu
	helpMenu := &desktop.CommonMenuItem{
		ID:       "help_menu",
		Label:    "Help",
		SubItems: make([]*desktop.CommonMenuItem, 0),
	}

	// Add a Documentation item
	helpMenu.SubItems = append(helpMenu.SubItems, &desktop.CommonMenuItem{
		ID:    "help_docs",
		Label: "Documentation",
		Action: func() {
			log.Println("Documentation clicked")
		},
	})

	// Add an About item
	helpMenu.SubItems = append(helpMenu.SubItems, &desktop.CommonMenuItem{
		ID:    "help_about",
		Label: "About",
		Action: func() {
			log.Println("About clicked")
		},
	})

	// Add the Help menu to the menu bar
	nativeMenu.AddItem(helpMenu)

	// Build the app configuration
	config := desktop.DefaultConfig()
	config.Title = "Native Menu Example"
	config.Width = 800
	config.Height = 600
	config.Debug = debug
	config.StaticDir = staticDir
	config.NativeMenu = nativeMenu

	// Display a simple HTML page
	content := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Native Menu Example</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				margin: 0;
				padding: 20px;
				text-align: center;
			}
			h1 {
				color: #333;
			}
			.info {
				background-color: #f0f0f0;
				border-radius: 5px;
				padding: 20px;
				margin: 20px auto;
				max-width: 600px;
				text-align: left;
			}
			.note {
				color: #666;
				font-style: italic;
			}
		</style>
	</head>
	<body>
		<h1>Native Menu Example</h1>
		<div class="info">
			<p>This example demonstrates the usage of native menus with webview.</p>
			<p>The following menus have been added:</p>
			<ul>
				<li><strong>File Menu:</strong> New, Open, Save, Exit</li>
				<li><strong>Edit Menu:</strong> Cut, Copy, Paste, Select All</li>
				<li><strong>Help Menu:</strong> Documentation, About</li>
			</ul>
			<p class="note">Note: Check the console for menu click events.</p>
		</div>
		<p>Running on: <strong>` + runtime.GOOS + `</strong></p>
	</body>
	</html>
	`

	// Setup the webserver
	tmpDir, err := os.MkdirTemp("", "velo-menu-example")
	if err != nil {
		log.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Save HTML content to a file
	htmlPath := tmpDir + "/index.html"
	err = os.WriteFile(htmlPath, []byte(content), 0644)
	if err != nil {
		log.Fatalf("Failed to write HTML file: %v", err)
	}

	// Set the static directory to our temporary directory
	config.StaticDir = tmpDir

	// Run the app
	if err := desktop.Run(config); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
