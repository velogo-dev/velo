package utils

import (
	"os"
	"os/exec"
)

// RunCmd executes a shell command and connects it to stdout/stderr
func RunCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunCmdWithDir executes a shell command in the specified directory
func RunCmdWithDir(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunCmdInBackground executes a shell command in the background
func RunCmdInBackground(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}

// RunCmdWait executes a shell command and waits for it to complete
// while allowing for interactive input
func RunCmdWait(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin // Add stdin for interactive prompts
	return cmd.Run()
}
