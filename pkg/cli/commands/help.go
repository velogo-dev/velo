package commands

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/velogo-dev/velo/constants"
	"github.com/velogo-dev/velo/pkg/utils"
)

// HelpCommand displays help information for the CLI
// velo -h
// velo --help
func (c *command) HelpCommand() error {
	// Define styles
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00FFFF")).
		MarginBottom(1)

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFF00")).
		PaddingLeft(2).
		MarginTop(1)

	commandStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")).
		PaddingLeft(2)

	descriptionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#CCCCCC")).
		PaddingLeft(8)

	// Get version
	latestTag, err := utils.GetLatestGitTag()
	if err != nil {
		fmt.Printf("Error getting latest git tag: %s\n", err)
	}

	// Render title with version
	fmt.Println(titleStyle.Render("✨ Velo CLI " + latestTag + " ✨"))

	// Render header
	fmt.Println(headerStyle.Render("Available commands:"))

	// Render commands
	for _, command := range constants.AllCommands() {
		fmt.Println(commandStyle.Render("➜ " + command.Name))
		fmt.Println(descriptionStyle.Render(command.Description))
	}

	return nil
}
