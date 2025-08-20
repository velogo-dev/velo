package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/velogo-dev/velo/pkg/cli"
)

func main() {
	// Create a new CLI instance with all commands registered
	err := cli.New().Run()
	if err != nil {
		fmt.Println(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000")).Render(err.Error()))
		os.Exit(1)
	}
}
