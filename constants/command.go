package constants

type Command struct {
	Name        string
	Args        []string
	Description string
}

var (
	InitCommand = Command{
		Name:        "init",
		Args:        []string{"init", "<app-name>", "--framework", "<framework>", "--template", "<template>"},
		Description: "Initialize a new Velo project example: velo init -n my-app --framework react --template nextjs",
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
		Args:        []string{"version"},
		Description: "Show the version of a Velo project",
	}
)

func GetCommand(name string) Command {
	switch name {
	case "init":
		return InitCommand
	case "show":
		return ShowCommand
	case "build":
		return BuildCommand
	case "dev":
		return DevCommand
	case "help":
		return HelpCommand
	case "doctor":
		return DoctorCommand
	case "version":
		return VersionCommand
	}

	return Command{}
}
