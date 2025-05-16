// Package cli provides the command-line interface functionality for the Velo application.
//
// It handles command parsing, execution, and provides a structured way to interact
// with the application through the terminal. The package is organized into:
// - Core CLI functionality (this file)
// - Command implementations (commands package)
// - Command registration and management
package cli

import (
	"context"
	"os"

	"github.com/velogo-dev/velo/constants"
	"github.com/velogo-dev/velo/pkg/cli/commands"
)

// VeloCLI represents the main command-line interface application
type VeloCLI struct {
	AppName string
	// Version is the version of the CLI application
	Version string
	// Commands is a map of registered commands
}

// New creates a new CLI instance with registered commands
func New(options ...func(*VeloCLI)) *VeloCLI {
	cli := &VeloCLI{
		Version: "0.0.1",
	}

	for _, option := range options {
		option(cli)
	}
	return cli
}

func (c *VeloCLI) Run() error {
	if len(os.Args) == 0 {
		return commands.NewCommand().HelpCommand()
	}
	command := constants.GetCommand(os.Args[1])
	switch command.Name {
	case constants.InitCommand.Name:
		return commands.NewCommand().InitCommand(context.Background(), os.Args[2:])
	case constants.BuildCommand.Name:
		return commands.NewCommand().BuildCommand()
	case constants.DevCommand.Name:
		return commands.NewCommand().DevCommand()
	case constants.HelpCommand.Name:
		return commands.NewCommand().HelpCommand()
	case constants.DoctorCommand.Name:
		return commands.NewCommand().DoctorCommand()
	case constants.VersionCommand.Name:
		return commands.NewCommand().VersionCommand()
	default:
		return commands.NewCommand().HelpCommand()
	}
}
