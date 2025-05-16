package commands

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/velogo-dev/velo/pkg/utils"
)

func (c *command) VersionCommand() error {
	latestTag, err := utils.GetLatestGitTag()
	if err != nil {
		fmt.Printf("Error getting latest git tag: %s\n", err)
	}
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00FFFF"))
	fmt.Println(titleStyle.Render("✨ Velo CLI " + latestTag + " ✨"))
	return nil
}
