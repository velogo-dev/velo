package commands

import (
	"fmt"

	"github.com/velogo-dev/velo/constants"
)

// ShowCommand implements the 'show' command to display various information
func ShowCommand(cli CLI) error {
	// Extract data from the CLI instance
	app, ok := cli.(*App)
	if !ok {
		return fmt.Errorf("invalid CLI implementation")
	}

	if len(app.Args) < 2 {
		fmt.Println("Error: Missing argument for 'show' command")
		fmt.Println("Usage: velo show [frameworks|templates|config]")
		return fmt.Errorf("missing argument")
	}

	switch app.Args[1] {
	case "frameworks":
		fmt.Println("Available Frameworks:")
		fmt.Println("--------------------")
		for _, fw := range app.AvailableFrameworks {
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
			for fw, templates := range app.FrameworkTemplates {
				fmt.Printf("- %s:\n", fw)
				for _, tmpl := range templates {
					fmt.Printf("  - %s\n", tmpl.Name)
				}
			}
		} else {
			fw := constants.Framework(framework)
			if templates, ok := app.FrameworkTemplates[fw]; ok {
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
		fmt.Printf("Version: %s\n", cli.Version())

	default:
		fmt.Printf("Unknown argument for 'show' command: %s\n", app.Args[1])
		fmt.Println("Usage: velo show [frameworks|templates|config]")
		return fmt.Errorf("unknown argument")
	}

	return nil
}
