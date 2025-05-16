package constants

import "strings"


type Command struct {
	Name        string
	Args        []string
	Description string
}

var (
	InitCommand = Command{
		Name:        "init",
		Args:        []string{"init", "<app-name>", "--library", "<library>", "--framework", "<framework>", "--template", "<template>"},
		Description: "Initialize a new Velo project example: velo init -n my-app --library react --framework nextjs --template nextjs",
	}
	ShowCommand = Command{
		Name:        "show",
		Args:        []string{"show", "<app-name>"},
		Description: "Show a Velo project",
	}
	BuildCommand = Command{
		Name:        "build",
		Args:        []string{"build", "<app-name>"},
		Description: "Build a Velo project",
	}
	DevCommand = Command{
		Name:        "dev",
		Args:        []string{"dev", "<app-name>"},
		Description: "Run a Velo project in development mode",
	}
	HelpCommand = Command{
		Name:        "help",
		Args:        []string{"help"},
		Description: "Show help for a Velo project",
	}
	DoctorCommand = Command{
		Name:        "doctor",
		Args:        []string{"doctor"},
		Description: "Check the health of a Velo project",
	}
	VersionCommand = Command{
		Name:        "version",
		Args:        []string{"-v", "--version"},
		Description: "Show the version of a Velo project",
	}
)

func GetCommand(name string) Command {
	commands := map[string]Command{
		"init, -i, --init":       InitCommand,
		"show,  --show":          ShowCommand,
		"build,  --build":        BuildCommand,
		"dev,  --dev":            DevCommand,
		"help, -h, --help":       HelpCommand,
		"doctor, --doctor":       DoctorCommand,
		"version, -v, --version": VersionCommand,
	}
	for key, command := range commands {
		if strings.Contains(key, name) {
			return command
		}
	}
	return Command{}
}

func AllCommands() []Command {
	return []Command{
		InitCommand,
		ShowCommand,
		BuildCommand,
		DevCommand,
		HelpCommand,
		DoctorCommand,
		VersionCommand,
	}
}

func GetCommandArgs(name string) []string {
	return GetCommand(name).Args
}

// Commands array contains all available CLI commands
var Commands = []Command{
	InitCommand,
	ShowCommand,
	BuildCommand,
	DevCommand,
	HelpCommand,
	DoctorCommand,
	VersionCommand,
}
