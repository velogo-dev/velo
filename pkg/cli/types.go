// Package cli provides the command-line interface functionality for the Velo application.
package cli

// Configuration stores CLI configuration settings
type Configuration struct {
	// Version of the CLI application
	Version string
	// ConfigDir is the directory for configuration files
	ConfigDir string
}

// Command represents a CLI command with its metadata and action function
type Command struct {
	// Name of the command as it should be typed in the terminal
	Name string
	// Description provides a short explanation of what the command does
	Description string
	// Action is the function that will be executed when the command is invoked
	Action func(*CLI) error
	// Usage provides detailed usage information for the command
	Usage string
}

// CLI represents the main command-line interface application
type CLI struct {
	// AppName is the name of the application from command line arguments
	AppName string
	// Framework is the selected frontend framework
	Framework string
	// SubFramework is the selected framework template
	SubFramework string
	// Args contains the command line arguments
	Args []string
	// Config stores configuration settings
	Config *Configuration
	// Commands maps command names to their implementations
	Commands map[string]Command
}
