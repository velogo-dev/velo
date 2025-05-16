package commands

import (
	"fmt"
)

// DevCommand implements the 'dev' command to run the application in development mode
func (c *command) DevCommand() error {
	fmt.Println("Starting development server...")

	// Process dev arguments
	var (
		port = "3000"
		host = "localhost"
	)

	for i := 1; i < len(c.Args); i++ {
		if c.Args[i] == "--port" || c.Args[i] == "-p" {
			if i+1 < len(c.Args) {
				port = c.Args[i+1]
				i++
			}
		} else if c.Args[i] == "--host" || c.Args[i] == "-h" {
			if i+1 < len(c.Args) {
				host = c.Args[i+1]
				i++
			}
		}
	}

	fmt.Printf("Dev server running at: http://%s:%s\n", host, port)
	fmt.Println("Press Ctrl+C to stop the server")

	// TODO: Implement actual dev server logic
	select {} // Keep the application running
}
