package app

import (
	"fmt"
	"os"

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
	ConfigDir string
}

// Config represents the application configuration
type App struct {
	AppName      string `json:"app_name" yaml:"app_name"`           // Name of the application from command line arguments
	Framework    string `json:"framework" yaml:"framework"`         // Selected framework
	SubFramework string `json:"sub_framework" yaml:"sub_framework"` // Selected sub-framework
	Config       *configuration
}

// New creates a new configuration from command-line flags
func New(options ...func(*App)) *App {
	app := &App{}
	for _, option := range options {
		option(app)
	}
	return app
}

// GetAssetsDir returns the path to the assets directory
// func (a *App) GetAssetsDir() string {
// return filepath.Join(a.RootDir, "mobile-shell", "assets")
// }

// GetFrontendDir returns the path to the frontend directory
// func (a *App) GetFrontendDir() string {
// 	return filepath.Join(a.RootDir, "frontend")
// }

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
