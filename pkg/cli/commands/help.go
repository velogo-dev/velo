package commands

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/velogo-dev/velo/pkg/utils"
)

// HelpCommand displays help information for the CLI
// velo -h
// velo --help
func (c *command) HelpCommand() error {
	latestTag, err := utils.GetLatestGitTag()
	if err != nil {
		return err
	}
	fmt.Printf("Velo CLI v%s\n\n", latestTag)
	fmt.Println("Available commands:")
	fmt.Println("  init")
	fmt.Println("        Initializes a new Velo project.")
	fmt.Println("  build")
	fmt.Println("        Builds the application.")
	fmt.Println("  dev")
	fmt.Println("        Runs the application in development mode.")
	fmt.Println("  doctor")
	fmt.Println("        Diagnose your environment.")
	fmt.Println("  update")
	fmt.Println("        Update the Velo CLI.")
	fmt.Println("  show")
	fmt.Println("        Shows various information.")
	fmt.Println("  generate")
	fmt.Println("        Code Generation Tools.")
	fmt.Println("  help")
	fmt.Println("        Show help.")
	fmt.Println("\nFlags:")

	return nil
}

func getLatestGitTag() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
