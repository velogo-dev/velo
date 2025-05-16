// Package commands contains the implementations for all CLI commands.
package commands

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/velogo-dev/velo/constants"
	"github.com/velogo-dev/velo/internal"
)

// CLI is an interface that command implementations can use
// without creating import cycles
type CLI interface {
	Version() string
	ShowHelp()
	RegisterCommand(name, description, usage string, action func(CLI) error)
}

// InitCommand implements the 'init' command to create a new Velo project
func InitCommand(cli CLI) error {
	fmt.Println("Initializing a new Velo project...")

	// Extract data from the CLI instance
	app, ok := cli.(*App)
	if !ok {
		return fmt.Errorf("invalid CLI implementation")
	}

	// If no app name is provided, run the interactive init
	if len(app.Args) == 1 || (len(app.Args) > 1 && app.Args[1] == "--interactive") {
		return runInteractiveInit(app)
	}

	// Process init arguments
	if len(app.Args) > 1 {
		app.AppName = app.Args[1]
		fmt.Printf("Creating new project: %s\n", app.AppName)

		// Check for framework flag
		for i := 2; i < len(app.Args); i++ {
			if app.Args[i] == "--framework" || app.Args[i] == "-f" {
				if i+1 < len(app.Args) {
					app.Framework = app.Args[i+1]
					i++
				}
			} else if app.Args[i] == "--template" || app.Args[i] == "-t" {
				if i+1 < len(app.Args) {
					app.SubFramework = app.Args[i+1]
					i++
				}
			}
		}
	}

	// Default values if not provided
	if app.Framework == "" {
		app.Framework = string(constants.React)
	}

	if app.SubFramework == "" {
		if templates, ok := app.FrameworkTemplates[constants.Framework(app.Framework)]; ok && len(templates) > 0 {
			app.SubFramework = templates[0].Name
		}
	}

	// Install the project
	installer := internal.NewFrameworkInstaller(constants.Framework(app.Framework),
		constants.SubFramework{Parent: constants.Framework(app.Framework), Name: app.SubFramework}, app.AppName)
	return installer.Install()
}

// App represents the CLI application with all required properties for commands
type App struct {
	CLI
	AppName             string
	Framework           string
	SubFramework        string
	Args                []string
	AvailableFrameworks []constants.Framework
	FrameworkTemplates  map[constants.Framework][]constants.SubFramework
}

// runInteractiveInit runs the interactive initialization process
func runInteractiveInit(app *App) error {
	withAppName(app)
	selectFramework(app)
	selectSubFramework(app)
	return install(app)
}

// withAppName sets the app name from command line arguments
func withAppName(app *App) {
	huh.NewInput().
		Title("Enter the name of your application").
		Value(&app.AppName).
		Run()
}

// selectFramework allows the user to select a framework
func selectFramework(app *App) {
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select the framework you want to use").
				Options(makeStringOptions(getFrameworkNames(app.AvailableFrameworks))...).
				Value(&app.Framework),
		),
	).WithTheme(huh.ThemeDracula()).Run()
	if err != nil {
		fmt.Println("Error selecting framework:", err)
		os.Exit(1)
	}
}

// selectSubFramework allows the user to select a subframework/template
func selectSubFramework(app *App) {
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(fmt.Sprintf("Select the template you want to use for %s", app.Framework)).
				Options(makeStringOptions(getTemplateNames(app.FrameworkTemplates[constants.Framework(app.Framework)]))...).
				Value(&app.SubFramework),
		),
	).WithTheme(huh.ThemeDracula()).Run()
	if err != nil {
		fmt.Println("Error selecting sub-framework:", err)
		os.Exit(1)
	}
}

// install handles the installation of the selected framework
func install(app *App) error {
	installer := internal.NewFrameworkInstaller(constants.Framework(app.Framework),
		constants.SubFramework{Parent: constants.Framework(app.Framework), Name: app.SubFramework},
		app.AppName)
	// Run synchronously to maintain stdin/stdout/stderr connections for interactive prompts
	return installer.Install()
}

// helper function to convert string slice to huh options
func makeStringOptions(items []string) []huh.Option[string] {
	options := make([]huh.Option[string], len(items))
	for i, item := range items {
		options[i] = huh.NewOption(item, item)
	}
	return options
}

// helper function to get framework names
func getFrameworkNames(frameworks []constants.Framework) []string {
	names := make([]string, len(frameworks))
	for i, fw := range frameworks {
		names[i] = string(fw)
	}
	return names
}

// helper function to get template names
func getTemplateNames(templates []constants.SubFramework) []string {
	names := make([]string, len(templates))
	for i, tmpl := range templates {
		names[i] = tmpl.Name
	}
	return names
}
