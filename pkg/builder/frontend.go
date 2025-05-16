package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/velogo-dev/velo/pkg/utils"
)

// Frontend represents the frontend web application
type Frontend struct {
	RootDir   string
	AssetsDir string
}

// NewFrontend creates a new frontend builder
func NewFrontend(rootDir string) *Frontend {
	return &Frontend{
		RootDir:   filepath.Join(rootDir, "frontend"),
		AssetsDir: filepath.Join(rootDir, "mobile-shell", "assets"),
	}
}

// Build builds the frontend for production
func (f *Frontend) Build() error {
	fmt.Println("Building frontend...")
	return utils.RunCmdWithDir(f.RootDir, "npm", "run", "build")
}

// StartDevServer starts the development server
func (f *Frontend) StartDevServer() error {
	fmt.Println("Starting frontend dev server...")
	return utils.RunCmdInBackground(f.RootDir, "npm", "run", "dev")
}

// CopyBuildToMobile copies the build output to mobile shell assets
func (f *Frontend) CopyBuildToMobile() error {
	fmt.Println("Copying build output to mobile shell assets...")

	// Create assets directory if it doesn't exist
	if err := os.MkdirAll(f.AssetsDir, 0755); err != nil {
		return fmt.Errorf("failed to create assets directory: %w", err)
	}

	src := filepath.Join(f.RootDir, "dist") // Adjust based on framework output

	// Use different commands based on OS
	if runtime.GOOS == "windows" {
		return utils.RunCmd("xcopy", "/E", "/I", "/Y", src, f.AssetsDir)
	}

	return utils.RunCmd("cp", "-r", src+"/.", f.AssetsDir)
}

// InstallDependencies installs all frontend dependencies
func (f *Frontend) InstallDependencies() error {
	fmt.Println("Installing frontend dependencies...")
	return utils.RunCmdWithDir(f.RootDir, "npm", "install")
}
