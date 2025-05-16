// Package commands provides the command-line interface functionality for the Velo application.
package commands

import "context"

// Command represents a CLI command with its metadata and action function
type command struct {
	Name   string
	Args   []string
	Desc   string
	Action func(ctx context.Context) error
}

func NewCommand(options ...func(*command)) *command {
	cmd := &command{}
	for _, option := range options {
		option(cmd)
	}
	return cmd
}

func WithName(name string) func(*command) {
	return func(cmd *command) {
		cmd.Name = name
	}
}

func WithAction(action func() error) func(*command) {
	return func(cmd *command) {
		action := func(ctx context.Context) error {
			return action()
		}
		cmd.Action = action
	}
}
func WithActionContext(action func(ctx context.Context) error) func(*command) {
	return func(cmd *command) {
		cmd.Action = action
	}
}
