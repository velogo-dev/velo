// Package commands contains the implementations for all CLI commands.
//
// This package contains the command handlers, argument processing, and
// business logic for each command available in the Velo CLI.
package commands

import (
	"github.com/velogo-dev/velo/constants"
)

// RegisterCommands registers all command handlers with the provided CLI
func RegisterCommands(cli CLI) {
	// Convert the CLI to our internal App type
	app, ok := cli.(*App)
	if !ok {
		// Initialize default frameworks data if not available
		app = &App{
			CLI: cli,
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
	}

	// Register all available commands
	app.RegisterCommand("init", "Initializes a new Velo project",
		"init [project-name] [--framework|-f <framework>] [--template|-t <template>] [--interactive]",
		InitCommand)

	app.RegisterCommand("build", "Builds the application",
		"build [--env|-e <environment>] [--output|-o <directory>]",
		BuildCommand)

	app.RegisterCommand("dev", "Runs the application in development mode",
		"dev [--port|-p <port>] [--host|-h <hostname>]",
		DevCommand)

	app.RegisterCommand("doctor", "Diagnose your environment",
		"doctor",
		DoctorCommand)

	app.RegisterCommand("update", "Update the Velo CLI",
		"update [--channel|-c <channel>] [--force]",
		UpdateCommand)

	app.RegisterCommand("show", "Shows various information",
		"show [frameworks|templates|config]",
		ShowCommand)

	app.RegisterCommand("generate", "Code Generation Tools",
		"generate [component|page|api|model] <name>",
		GenerateCommand)
}
