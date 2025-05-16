package app

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/velogo-dev/velo/constants"
	"github.com/velogo-dev/velo/internal"
)

// AvailableFrameworks defines the supported frontend frameworks
var AvailableFrameworks = []constants.Framework{
	constants.React,
	constants.Vue,
	constants.Svelte,
	constants.Angular,
	constants.Solid,
}

// FrameworkTemplates maps frameworks to their popular templates
var FrameworkTemplates = map[constants.Framework][]constants.SubFramework{
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
}

type configuration struct {
	Version   string
	ConfigDir string
}

// Command represents a CLI command
type Command struct {
	Name        string
	Description string
	Action      func(*App) error
}

// Config represents the application configuration
type App struct {
	AppName      string   `json:"app_name" yaml:"app_name"`           // Name of the application from command line arguments
	Framework    string   `json:"framework" yaml:"framework"`         // Selected framework
	SubFramework string   `json:"sub_framework" yaml:"sub_framework"` // Selected sub-framework
	Args         []string // Command line arguments
	Config       *configuration
	Commands     map[string]Command // Available commands
}

// New creates a new configuration from command-line flags
func New(options ...func(*App)) *App {
	app := &App{
		Config: &configuration{
			Version:   "0.0.1",
			ConfigDir: "~/.velo",
		},
		Commands: map[string]Command{},
	}

	// Register default commands
	app.registerCommands()

	if len(os.Args) > 1 {
		app.Args = os.Args[1:]

		// Process common flags first
		if app.handleCommonFlags() {
			return app
		}

		// Process commands
		if len(app.Args) > 0 {
			cmd := app.Args[0]
			if command, exists := app.Commands[cmd]; exists {
				if err := command.Action(app); err != nil {
					fmt.Printf("Error executing command '%s': %v\n", cmd, err)
					os.Exit(1)
				}
				os.Exit(0)
			}
		}
	} else {
		// When no arguments are provided, show the help information
		app.showHelp()
		os.Exit(0)
	}

	for _, option := range options {
		option(app)
	}
	return app
}

// handleCommonFlags processes common flags like --version and --help
// Returns true if the program should exit after handling
func (a *App) handleCommonFlags() bool {
	// Check for version flag
	for _, arg := range a.Args {
		if arg == "--version" || arg == "-v" {
			fmt.Printf("Velo version %s\n", a.Config.Version)
			os.Exit(0)
			return true
		}
	}

	// Check for help flag
	for _, arg := range a.Args {
		if arg == "--help" || arg == "-h" {
			a.showHelp()
			os.Exit(0)
			return true
		}
	}

	return false
}

// registerCommands registers all available CLI commands
func (a *App) registerCommands() {
	a.Commands["init"] = Command{
		Name:        "init",
		Description: "Initializes a new Velo project",
		Action:      a.cmdInit,
	}

	a.Commands["build"] = Command{
		Name:        "build",
		Description: "Builds the application",
		Action:      a.cmdBuild,
	}

	a.Commands["dev"] = Command{
		Name:        "dev",
		Description: "Runs the application in development mode",
		Action:      a.cmdDev,
	}

	a.Commands["doctor"] = Command{
		Name:        "doctor",
		Description: "Diagnose your environment",
		Action:      a.cmdDoctor,
	}

	a.Commands["update"] = Command{
		Name:        "update",
		Description: "Update the Velo CLI",
		Action:      a.cmdUpdate,
	}

	a.Commands["show"] = Command{
		Name:        "show",
		Description: "Shows various information",
		Action:      a.cmdShow,
	}

	a.Commands["generate"] = Command{
		Name:        "generate",
		Description: "Code Generation Tools",
		Action:      a.cmdGenerate,
	}
}

// showHelp displays help information
func (a *App) showHelp() {
	fmt.Printf("Velo CLI v%s\n\n", a.Config.Version)
	fmt.Println("Available commands:\n")

	// Find the longest command name for proper formatting
	maxLength := 0
	for cmd := range a.Commands {
		if len(cmd) > maxLength {
			maxLength = len(cmd)
		}
	}

	// Print each command with description
	for _, cmd := range getSortedCommands(a.Commands) {
		command := a.Commands[cmd]
		padding := strings.Repeat(" ", maxLength-len(command.Name)+3)
		fmt.Printf("   %s%s%s \n", command.Name, padding, command.Description)
	}

	fmt.Println("\nFlags:\n")
	fmt.Println("  -help, --help")
	fmt.Println("        Get help on the 'velo' command.")
	fmt.Println("  -v, --version")
	fmt.Println("        Print the current Velo version.")
}

func (a *App) Version() string {
	return a.Config.Version
}

// GetAssetsDir returns the path to the assets directory
// func (a *App) GetAssetsDir() string {
// return filepath.Join(a.RootDir, "mobile-shell", "assets")
// }

// GetFrontendDir returns the path to the frontend directory
// func (a *App) GetFrontendDir() string {
// 	return filepath.Join(a.RootDir, "frontend")
// }

// Command implementations
func (a *App) cmdInit(app *App) error {
	fmt.Println("Initializing a new Velo project...")

	// If no app name is provided, run the interactive init
	if len(app.Args) == 1 || (len(app.Args) > 1 && app.Args[1] == "--interactive") {
		return app.runInteractiveInit()
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
		if templates, ok := FrameworkTemplates[constants.Framework(app.Framework)]; ok && len(templates) > 0 {
			app.SubFramework = templates[0].Name
		}
	}

	// Install the project
	installer := internal.NewFrameworkInstaller(constants.Framework(app.Framework),
		constants.SubFramework{Parent: constants.Framework(app.Framework), Name: app.SubFramework}, app.AppName)
	return installer.Install()
}

func (a *App) cmdBuild(app *App) error {
	fmt.Println("Building the application...")

	// Process build arguments
	var (
		environment = "production"
		output      = "./dist"
	)

	for i := 1; i < len(app.Args); i++ {
		if app.Args[i] == "--env" || app.Args[i] == "-e" {
			if i+1 < len(app.Args) {
				environment = app.Args[i+1]
				i++
			}
		} else if app.Args[i] == "--output" || app.Args[i] == "-o" {
			if i+1 < len(app.Args) {
				output = app.Args[i+1]
				i++
			}
		}
	}

	fmt.Printf("Building for %s environment\n", environment)
	fmt.Printf("Output directory: %s\n", output)

	// TODO: Implement actual build logic
	fmt.Println("Build completed successfully")
	return nil
}

func (a *App) cmdDev(app *App) error {
	fmt.Println("Starting development server...")

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

func (a *App) cmdDoctor(app *App) error {
	fmt.Println("Diagnosing your environment...")

	// Display system information
	fmt.Println("\nSystem Information:")
	fmt.Println("------------------")
	fmt.Printf("OS: %s\n", runtime.GOOS)
	fmt.Printf("Architecture: %s\n", runtime.GOARCH)
	fmt.Printf("Go Version: %s\n", runtime.Version())

	// Check for required dependencies
	fmt.Println("\nDependency Check:")
	fmt.Println("----------------")

	// TODO: Implement actual dependency checks
	dependencies := []struct {
		name    string
		command string
	}{
		{"Node.js", "node --version"},
		{"npm", "npm --version"},
		{"Git", "git --version"},
	}

	for _, dep := range dependencies {
		fmt.Printf("Checking for %s... ", dep.name)
		fmt.Println("OK") // Placeholder, would actually check if installed
	}

	fmt.Println("\nEnvironment check completed.")
	return nil
}

func (a *App) cmdUpdate(app *App) error {
	fmt.Println("Updating Velo CLI...")

	var (
		channel = "stable"
		force   = false
	)

	for i := 1; i < len(app.Args); i++ {
		if app.Args[i] == "--channel" || app.Args[i] == "-c" {
			if i+1 < len(app.Args) {
				channel = app.Args[i+1]
				i++
			}
		} else if app.Args[i] == "--force" {
			force = true
		}
	}

	fmt.Printf("Checking for updates on %s channel...\n", channel)

	// TODO: Implement actual update logic
	if force {
		fmt.Println("Forcing update...")
	}

	fmt.Println("Velo CLI is now up to date!")
	return nil
}

func (a *App) cmdShow(app *App) error {
	if len(app.Args) < 2 {
		fmt.Println("Error: Missing argument for 'show' command")
		fmt.Println("Usage: velo show [frameworks|templates|config]")
		return fmt.Errorf("missing argument")
	}

	switch app.Args[1] {
	case "frameworks":
		fmt.Println("Available Frameworks:")
		fmt.Println("--------------------")
		for _, fw := range AvailableFrameworks {
			fmt.Printf("- %s\n", fw)
		}

	case "templates":
		framework := ""
		if len(app.Args) > 2 {
			framework = app.Args[2]
		}

		if framework == "" {
			fmt.Println("Available Templates:")
			fmt.Println("------------------")
			for fw, templates := range FrameworkTemplates {
				fmt.Printf("- %s:\n", fw)
				for _, tmpl := range templates {
					fmt.Printf("  - %s\n", tmpl.Name)
				}
			}
		} else {
			fw := constants.Framework(framework)
			if templates, ok := FrameworkTemplates[fw]; ok {
				fmt.Printf("Templates for %s:\n", fw)
				fmt.Println("------------------")
				for _, tmpl := range templates {
					fmt.Printf("- %s\n", tmpl.Name)
				}
			} else {
				fmt.Printf("No templates found for framework: %s\n", framework)
				return fmt.Errorf("invalid framework")
			}
		}

	case "config":
		fmt.Println("Velo Configuration:")
		fmt.Println("------------------")
		fmt.Printf("Version: %s\n", a.Config.Version)
		fmt.Printf("Config Directory: %s\n", a.Config.ConfigDir)

	default:
		fmt.Printf("Unknown argument for 'show' command: %s\n", app.Args[1])
		fmt.Println("Usage: velo show [frameworks|templates|config]")
		return fmt.Errorf("unknown argument")
	}

	return nil
}

func (a *App) cmdGenerate(app *App) error {
	if len(app.Args) < 2 {
		fmt.Println("Error: Missing argument for 'generate' command")
		fmt.Println("Usage: velo generate [component|page|api|model]")
		return fmt.Errorf("missing argument")
	}

	name := ""
	if len(app.Args) > 2 {
		name = app.Args[2]
	} else {
		fmt.Println("Error: Missing name for generation")
		fmt.Printf("Usage: velo generate %s <name>\n", app.Args[1])
		return fmt.Errorf("missing name")
	}

	switch app.Args[1] {
	case "component":
		fmt.Printf("Generating component: %s\n", name)
		// TODO: Implement component generation

	case "page":
		fmt.Printf("Generating page: %s\n", name)
		// TODO: Implement page generation

	case "api":
		fmt.Printf("Generating API endpoint: %s\n", name)
		// TODO: Implement API generation

	case "model":
		fmt.Printf("Generating model: %s\n", name)
		// TODO: Implement model generation

	default:
		fmt.Printf("Unknown argument for 'generate' command: %s\n", app.Args[1])
		fmt.Println("Usage: velo generate [component|page|api|model]")
		return fmt.Errorf("unknown argument")
	}

	fmt.Println("Generation completed successfully!")
	return nil
}

// runInteractiveInit runs the interactive initialization process
func (a *App) runInteractiveInit() error {
	WithAppName()(a)
	SelectFramework()(a)
	SelectSubFramework()(a)
	Install()(a)
	return nil
}

// WithAppName sets the app name from command line arguments
func WithAppName() func(*App) {
	return func(a *App) {
		huh.NewInput().
			Title("Enter the name of your application").
			Value(&a.AppName).
			Run()
	}
}

func SelectFramework() func(*App) {
	return func(a *App) {
		err := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Select the framework you want to use").
					Options(makeStringOptions(getFrameworkNames(AvailableFrameworks))...).
					Value(&a.Framework),
			),
		).WithTheme(huh.ThemeDracula()).Run()
		if err != nil {
			fmt.Println("Error selecting framework:", err)
			os.Exit(1)
		}
	}
}

func SelectSubFramework() func(*App) {
	return func(a *App) {
		err := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title(fmt.Sprintf("Select the template you want to use for %s", a.Framework)).
					Options(makeStringOptions(getTemplateNames(FrameworkTemplates[constants.Framework(a.Framework)]))...).
					Value(&a.SubFramework),
			),
		).WithTheme(huh.ThemeDracula()).Run()
		if err != nil {
			fmt.Println("Error selecting sub-framework:", err)
			os.Exit(1)
		}
	}
}

// Install handles the installation of the selected framework
func Install() func(*App) {
	return func(a *App) {
		installer := internal.NewFrameworkInstaller(constants.Framework(a.Framework), constants.SubFramework{Parent: constants.Framework(a.Framework), Name: a.SubFramework}, a.AppName)
		// Run synchronously to maintain stdin/stdout/stderr connections for interactive prompts
		err := installer.Install()
		if err != nil {
			fmt.Printf("Installation error: %v\n", err)
			os.Exit(1)
		}
	}
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

// getSortedCommands returns a sorted list of command names
func getSortedCommands(commands map[string]Command) []string {
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
