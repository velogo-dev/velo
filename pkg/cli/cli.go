// Package cli provides the command-line interface functionality for the Velo application.
//
// It handles command parsing, execution, and provides a structured way to interact
// with the application through the terminal. The package is organized into:
// - Core CLI functionality (this file)
// - Command implementations (commands package)
// - Command registration and management
package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/velogo-dev/velo/constants"
	"github.com/velogo-dev/velo/pkg/cli/commands"
)

// Configuration stores CLI configuration settings
type ConfigSettings struct {
	// Version of the CLI application
	Version string
	// ConfigDir is the directory for configuration files
	ConfigDir string
}

// Command represents a CLI command with its metadata and action function
type CommandDef struct {
	// Name of the command as it should be typed in the terminal
	Name string
	// Description provides a short explanation of what the command does
	Description string
	// Action is the function that will be executed when the command is invoked
	Action func(commands.CLI) error
	// Usage provides detailed usage information for the command
	Usage string
}

// CLI represents the main command-line interface application
type VeloCLI struct {
	// AppName is the name of the application from command line arguments
	AppName string
	// Framework is the selected frontend framework
	Framework string
	// SubFramework is the selected framework template
	SubFramework string
	// Args contains the command line arguments
	Args []string
	// Config stores configuration settings
	Config *ConfigSettings
	// Commands maps command names to their implementations
	Commands map[string]CommandDef
	// Available frameworks
	AvailableFrameworks []constants.Framework
	// Framework templates
	FrameworkTemplates map[constants.Framework][]constants.SubFramework
}

// New creates a new CLI instance with registered commands
func New(options ...func(*VeloCLI)) *VeloCLI {
	cli := &VeloCLI{
		Config: &ConfigSettings{
			Version:   "0.0.1",
			ConfigDir: "~/.velo",
		},
		Commands: map[string]CommandDef{},
		AvailableFrameworks: []constants.Framework{
			constants.React,
			constants.Vue,
			constants.Svelte,
			constants.Angular,
			constants.Solid,
		},
		FrameworkTemplates: map[constants.Framework][]constants.SubFramework{
			constants.React: {
				constants.CreateReactApp,
				constants.NextJS,
				constants.Remix,
				constants.ReactVite,
			},
			constants.Vue: {
				constants.Nuxt,
				constants.Quasar,
				constants.VueVite,
			},
			constants.Svelte: {
				constants.SvelteKit,
				constants.SvelteVite,
			},
			constants.Angular: {
				constants.AngularUniversal,
				constants.Nest,
			},
			constants.Solid: {
				constants.SolidStart,
				constants.SolidVite,
			},
		},
	}

	// Register available commands
	commands.RegisterCommands(cli)

	if len(os.Args) > 1 {
		cli.Args = os.Args[1:]

		// Process common flags first
		if cli.handleCommonFlags() {
			return cli
		}

		// Process commands
		if len(cli.Args) > 0 {
			cmd := cli.Args[0]
			if command, exists := cli.Commands[cmd]; exists {
				if err := command.Action(cli); err != nil {
					fmt.Printf("Error executing command '%s': %v\n", cmd, err)
					os.Exit(1)
				}
				os.Exit(0)
			}
		}
	} else {
		// When no arguments are provided, show the help information
		cli.ShowHelp()
		os.Exit(0)
	}

	for _, option := range options {
		option(cli)
	}
	return cli
}

// RegisterCommand registers a new command in the CLI
func (c *VeloCLI) RegisterCommand(name, description, usage string, action func(commands.CLI) error) {
	c.Commands[name] = CommandDef{
		Name:        name,
		Description: description,
		Action:      action,
		Usage:       usage,
	}
}

// handleCommonFlags processes common flags like --version and --help
// Returns true if the program should exit after handling
func (c *VeloCLI) handleCommonFlags() bool {
	// Check for version flag
	for _, arg := range c.Args {
		if arg == "--version" || arg == "-v" {
			fmt.Printf("Velo version %s\n", c.Config.Version)
			os.Exit(0)
			return true
		}
	}

	// Check for help flag
	for _, arg := range c.Args {
		if arg == "--help" || arg == "-h" {
			c.ShowHelp()
			os.Exit(0)
			return true
		}
	}

	return false
}

// ShowHelp displays help information for the CLI
func (c *VeloCLI) ShowHelp() {
	fmt.Printf("Velo CLI v%s\n\n", c.Config.Version)
	fmt.Println("Available commands:\n")

	// Find the longest command name for proper formatting
	maxLength := 0
	for cmd := range c.Commands {
		if len(cmd) > maxLength {
			maxLength = len(cmd)
		}
	}

	// Print each command with description
	for _, cmd := range getSortedCommands(c.Commands) {
		command := c.Commands[cmd]
		padding := strings.Repeat(" ", maxLength-len(command.Name)+3)
		fmt.Printf("   %s%s%s \n", command.Name, padding, command.Description)
	}

	fmt.Println("\nFlags:\n")
	fmt.Println("  -help, --help")
	fmt.Println("        Get help on the 'velo' command.")
	fmt.Println("  -v, --version")
	fmt.Println("        Print the current Velo version.")
}

// Version returns the CLI version
func (c *VeloCLI) Version() string {
	return c.Config.Version
}

// GetFrameworkByName returns a Framework constant from a string name
func (c *VeloCLI) GetFrameworkByName(name string) constants.Framework {
	return constants.Framework(name)
}

// GetSubFrameworkTemplate returns the SubFramework template for the given framework and template name
func (c *VeloCLI) GetSubFrameworkTemplate(framework constants.Framework, templateName string) constants.SubFramework {
	return constants.SubFramework{
		Parent: framework,
		Name:   templateName,
	}
}

// getSortedCommands returns a sorted list of command names
func getSortedCommands(commands map[string]CommandDef) []string {
	names := make([]string, 0, len(commands))
	for name := range commands {
		names = append(names, name)
	}
	// Basic alphabetical sort
	for i := 0; i < len(names)-1; i++ {
		for j := i + 1; j < len(names); j++ {
			if names[i] > names[j] {
				names[i], names[j] = names[j], names[i]
			}
		}
	}
	return names
}
