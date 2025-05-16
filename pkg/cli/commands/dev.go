package commands

import (
	"fmt"
)

// DevCommand implements the 'dev' command to run the application in development mode
func DevCommand(cli CLI) error {
	fmt.Println("Starting development server...")

	// Extract data from the CLI instance
	app, ok := cli.(*App)
	if !ok {
		return fmt.Errorf("invalid CLI implementation")
	}

	// Process dev arguments
	var (
		port = "3000"
		host = "localhost"
	)

	for i := 1; i < len(app.Args); i++ {
		if app.Args[i] == "--port" || app.Args[i] == "-p" {
			if i+1 < len(app.Args) {
				port = app.Args[i+1]
				i++
			}
		} else if app.Args[i] == "--host" || app.Args[i] == "-h" {
			if i+1 < len(app.Args) {
				host = app.Args[i+1]
				i++
			}
		}
	}

	fmt.Printf("Dev server running at: http://%s:%s\n", host, port)
	fmt.Println("Press Ctrl+C to stop the server")

	// TODO: Implement actual dev server logic
	select {} // Keep the application running
}
