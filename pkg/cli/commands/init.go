// Package commands contains the implementations for all CLI commands.
package commands

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/velogo-dev/velo/constants"
	"github.com/velogo-dev/velo/internal"
)

// commandLineFlags for the init command
var (
	appName   string
	library   string
	framework string
)

// InitCommand implements the 'init' command to create a new Velo project.
//
// This command guides users through the process of creating a new Velo project
// either interactively or using command-line arguments. It supports specifying
// the application name, UI library, and framework (template) options.
//
// Parameters:
//   - ctx: A context.Context for cancellation support
//   - args: Command-line arguments passed to the init command
//
// Command syntax:
//   - velo init
//   - velo init <app-name>
//   - velo init <app-name> --library|-l <library-name>
//   - velo init <app-name> --library|-l <library-name> --framework|-f <framework-name>
//
// Returns:
//   - error: nil on successful completion, otherwise an error describing what went wrong
func (c *command) InitCommand(ctx context.Context, args []string) error {
	// Reset global variables to avoid state persistence between command invocations
	appName = ""
	library = ""
	framework = ""

	// Interactive mode (no arguments)
	if len(args) == 0 {
		if err := withAppName(); err != nil {
			return fmt.Errorf("failed to get application name: %w", err)
		}

		selectLibrary()
		selectFramework()
		return install()
	}

	// Process first argument as app name if provided
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		appName = args[0]
	}

	// Process command line flags
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--library", "-l":
			if i+1 < len(args) {
				library = args[i+1]
				i++ // Skip next argument as it's the value
			} else {
				return fmt.Errorf("missing value for %s flag", args[i])
			}
		case "--framework", "-f":
			if i+1 < len(args) {
				framework = args[i+1]
				i++ // Skip next argument as it's the value
			} else {
				return fmt.Errorf("missing value for %s flag", args[i])
			}
		case "--name", "-n":
			if i+1 < len(args) {
				appName = args[i+1]
				i++ // Skip next argument as it's the value
			} else {
				return fmt.Errorf("missing value for %s flag", args[i])
			}
		}
	}

	// Validate that we have an app name
	if appName == "" {
		if err := withAppName(); err != nil {
			return fmt.Errorf("failed to get application name: %w", err)
		}
	}

	// If library is provided but framework isn't, prompt for framework
	if library != "" && framework == "" {
		// Validate the provided library is supported
		if !isValidLibrary(library) {
			return fmt.Errorf("unsupported library: %s. Available libraries: %s",
				library, strings.Join(getLibraryNames(constants.AvailableLibraries), ", "))
		}
		selectFramework()
	}

	// If neither library nor framework is provided, prompt for both
	if library == "" {
		selectLibrary()
		selectFramework()
	}

	// Proceed with installation
	return install()
}

// withAppName prompts the user to enter an application name if not provided
// via command-line arguments.
//
// Returns:
//   - error: nil on successful completion, otherwise an error if the prompt fails
func withAppName() error {
	err := huh.NewInput().
		Title("Enter the name of your application").
		Placeholder("my-app").
		Validate(func(s string) error {
			if strings.TrimSpace(s) == "" {
				return fmt.Errorf("application name cannot be empty")
			}
			if strings.Contains(s, " ") {
				return fmt.Errorf("application name cannot contain spaces")
			}
			return nil
		}).
		Value(&appName).
		Run()

	if err != nil {
		return err
	}
	return nil
}

// selectLibrary displays an interactive prompt for the user to select
// a UI library from the available options.
func selectLibrary() {
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select the library you want to use").
				Options(makeStringOptions(getLibraryNames(constants.AvailableLibraries))...).
				Value(&library),
		),
	).WithTheme(huh.ThemeDracula()).Run()

	if err != nil {
		fmt.Println("Error selecting library:", err)
		os.Exit(1)
	}
}

// selectFramework displays an interactive prompt for the user to select
// a framework/template for the previously selected library.
func selectFramework() {
	// Get available frameworks for the selected library
	frameworks := constants.LibraryFrameworks[constants.Library(library)]
	if len(frameworks) == 0 {
		fmt.Printf("No frameworks available for library: %s\n", library)
		os.Exit(1)
	}

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(fmt.Sprintf("Select the framework for %s", library)).
				Options(makeStringOptions(getFrameworkNames(frameworks))...).
				Value(&framework),
		),
	).WithTheme(huh.ThemeDracula()).Run()

	if err != nil {
		fmt.Println("Error selecting framework:", err)
		os.Exit(1)
	}
}

// install creates and initializes a new project using the selected library and framework.
//
// Returns:
//   - error: nil on successful completion, otherwise an error if the installation fails
func install() error {
	// Validate that all required parameters are set
	if appName == "" {
		return fmt.Errorf("application name not specified")
	}
	if library == "" {
		return fmt.Errorf("library not selected")
	}
	if framework == "" {
		return fmt.Errorf("framework not selected")
	}

	fmt.Printf("Creating new %s project with %s framework in directory: %s\n",
		library, framework, appName)

	installer := internal.NewFrameworkInstaller(constants.Framework{
		Parent: constants.Library(library),
		Name:   framework,
	}, appName)

	return installer.Install()
}

// isValidLibrary checks if the provided library name is supported.
//
// Parameters:
//   - lib: The library name to validate
//
// Returns:
//   - bool: true if the library is supported, false otherwise
func isValidLibrary(lib string) bool {
	for _, validLib := range constants.AvailableLibraries {
		if string(validLib) == lib {
			return true
		}
	}
	return false
}

// makeStringOptions converts a string slice to a slice of huh.Option
//
// Parameters:
//   - items: A slice of string values to convert
//
// Returns:
//   - []huh.Option[string]: A slice of options for use with the huh library
func makeStringOptions(items []string) []huh.Option[string] {
	options := make([]huh.Option[string], len(items))
	for i, item := range items {
		options[i] = huh.NewOption(item, item)
	}
	return options
}

// getLibraryNames extracts the string names of libraries from a slice of constants.Library
//
// Parameters:
//   - libraries: A slice of constants.Library values
//
// Returns:
//   - []string: A slice of library names as strings
func getLibraryNames(libraries []constants.Library) []string {
	names := make([]string, len(libraries))
	for i, lib := range libraries {
		names[i] = string(lib)
	}
	return names
}

// getFrameworkNames extracts the names of frameworks from a slice of constants.Framework
//
// Parameters:
//   - frameworks: A slice of constants.Framework values
//
// Returns:
//   - []string: A slice of framework names
func getFrameworkNames(frameworks []constants.Framework) []string {
	names := make([]string, len(frameworks))
	for i, fw := range frameworks {
		names[i] = fw.Name
	}
	return names
}
