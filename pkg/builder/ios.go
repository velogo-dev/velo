package builder

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/velogo-dev/velo/pkg/utils"
)

// IOS represents the iOS app builder
type IOS struct {
	RootDir          string
	ShellDir         string
	XcodeProjectPath string
	BuildPath        string
}

// NewIOS creates a new iOS builder
func NewIOS(rootDir string) *IOS {
	shellDir := filepath.Join(rootDir, "mobile-shell", "ios")

	return &IOS{
		RootDir:          rootDir,
		ShellDir:         shellDir,
		XcodeProjectPath: filepath.Join(shellDir, "GolangMobile.xcodeproj"),
		BuildPath:        filepath.Join(rootDir, "build"),
	}
}

// Build builds the iOS app
func (i *IOS) Build() error {
	fmt.Println("Building iOS app...")

	if runtime.GOOS != "darwin" {
		return fmt.Errorf("iOS builds are only supported on macOS")
	}

	return utils.RunCmd(
		"xcodebuild",
		"-project", i.XcodeProjectPath,
		"-scheme", "GolangMobile",
		"-configuration", "Debug",
		"-derivedDataPath", i.BuildPath,
	)
}

// InstallApp installs the app on the simulator or device
func (i *IOS) InstallApp(deviceID string) error {
	fmt.Println("Installing iOS app on simulator/device...")

	if runtime.GOOS != "darwin" {
		return fmt.Errorf("iOS app installation is only supported on macOS")
	}

	args := []string{"simctl", "install"}

	if deviceID == "" {
		args = append(args, "-simulator")
	} else {
		args = append(args, "-device", deviceID)
	}

	// Build path for the iOS app
	appPath := filepath.Join(i.BuildPath, "Build", "Products", "Debug-iphonesimulator", "GolangMobile.app")
	args = append(args, appPath)

	return utils.RunCmd("xcrun", args...)
}

// LaunchApp launches the app on the simulator or device
func (i *IOS) LaunchApp(deviceID string) error {
	fmt.Println("Launching iOS app...")

	if runtime.GOOS != "darwin" {
		return fmt.Errorf("iOS app launch is only supported on macOS")
	}

	args := []string{"simctl", "launch"}

	if deviceID == "" {
		args = append(args, "-simulator")
	} else {
		args = append(args, "-device", deviceID)
	}

	// Bundle ID
	args = append(args, "com.example.golangmobile")

	return utils.RunCmd("xcrun", args...)
}
