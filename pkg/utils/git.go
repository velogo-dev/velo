package utils

import (
	"os/exec"
	"strings"
)

func GetLatestGitTag() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func GetEmail() (string, error) {
	cmd := exec.Command("git", "config", "--get", "user.email")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func GetUsername() (string, error) {
	cmd := exec.Command("git", "config", "--get", "user.username")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func GetName() (string, error) {
	cmd := exec.Command("git", "config", "--get", "user.name")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func GitInit() error {
	cmd := exec.Command("git", "init")
	return cmd.Run()
}

func GitAdd() error {
	cmd := exec.Command("git", "add", ".")
	return cmd.Run()
}

func GitCommit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	return cmd.Run()
}

func GitPush() error {
	cmd := exec.Command("git", "push")
	return cmd.Run()
}

func GitBranch(name string) error {
	cmd := exec.Command("git", "branch", "-m", name)
	return cmd.Run()
}

func GitPull() error {
	cmd := exec.Command("git", "pull")
	return cmd.Run()
}
