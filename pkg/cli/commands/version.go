package commands

import (
	"fmt"

	"github.com/velogo-dev/velo/pkg/utils"
)

func (c *command) VersionCommand() error {
	latestTag, err := utils.GetLatestGitTag()
	if err != nil {
		fmt.Printf("Error getting latest git tag: %s\n", err)
	}
	fmt.Println("Velo CLI v" + latestTag)
	return nil
}
