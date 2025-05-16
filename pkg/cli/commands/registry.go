// Package commands contains the implementations for all CLI commands.
//
// This package contains the command handlers, argument processing, and
// business logic for each command available in the Velo CLI.
package commands

// RegisterCommands registers all command handlers with the provided CLI
// func (c *command) RegisterCommands() error {

// 	// Register all available commands
// 	c.Commands["init"] = commands.NewCommand("init", "Initializes a new Velo project",
// 		"init [project-name] [--framework|-f <framework>] [--template|-t <template>] [--interactive]",
// 		commands.InitCommand)

// 	app.RegisterCommand("build", "Builds the application",
// 		"build [--env|-e <environment>] [--output|-o <directory>]",
// 		BuildCommand)

// 	app.RegisterCommand("dev", "Runs the application in development mode",
// 		"dev [--port|-p <port>] [--host|-h <hostname>]",
// 		DevCommand)

// 	app.RegisterCommand("doctor", "Diagnose your environment",
// 		"doctor",
// 		DoctorCommand)

// 	app.RegisterCommand("update", "Update the Velo CLI",
// 		"update [--channel|-c <channel>] [--force]",
// 		UpdateCommand)

// 	app.RegisterCommand("show", "Shows various information",
// 		"show [frameworks|templates|config]",
// 		ShowCommand)

// 	app.RegisterCommand("generate", "Code Generation Tools",
// 		"generate [component|page|api|model] <name>",
// 		GenerateCommand)

// 	app.RegisterCommand("help", "Show help",
// 		"help",
// 		HelpCommand)
// }

// ShowHelp displays help information for the CLI

// func (a *App) RegisterCommand(name, description, usage string, action func() error) {
// 	a.Commands[name] = CommandDef{
// 		Name:        name,
// 		Description: description,
// 		Action:      action,
// 		Usage:       usage,
// 	}
// }

// handleCommonFlags processes common flags like --version and --help
// Returns true if the program should exit after handling
// func (c *command) handleCommonFlags() bool {
// 	// Check for version flag
// 	for _, arg := range c.Args {
// 		if arg == "--version" || arg == "-v" {
// 			fmt.Printf("Velo version %s\n", c.Version)
// 			os.Exit(0)
// 			return true
// 		}
// 	}

// 	// Check for help flag
// 	for _, arg := range c.Args {
// 		if arg == "--help" || arg == "-h" {
// 			c.ShowHelp()
// 			os.Exit(0)
// 			return true
// 		}
// 	}

// 	return false
// }
