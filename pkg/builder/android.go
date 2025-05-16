package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/velogo-dev/velo/pkg/utils"
)

// Android represents the Android app builder
type Android struct {
	RootDir     string
	ShellDir    string
	GradlewPath string
}

// NewAndroid creates a new Android builder
func NewAndroid(rootDir string) *Android {
	shellDir := filepath.Join(rootDir, "mobile-shell", "android")
	var gradlewPath string

	if runtime.GOOS == "windows" {
		gradlewPath = filepath.Join(shellDir, "gradlew.bat")
	} else {
		gradlewPath = filepath.Join(shellDir, "gradlew")
	}

	return &Android{
		RootDir:     rootDir,
		ShellDir:    shellDir,
		GradlewPath: gradlewPath,
	}
}

// Build builds the Android app
func (a *Android) Build() error {
	fmt.Println("Building Android app...")

	if _, err := os.Stat(a.GradlewPath); err != nil {
		return fmt.Errorf("Android build tools not found. Make sure the Android project is set up correctly: %w", err)
	}

	return utils.RunCmdWithDir(a.ShellDir, a.GradlewPath, "assembleDebug")
}

// InstallApp installs the app on the device
func (a *Android) InstallApp(deviceID string) error {
	fmt.Println("Installing Android app on device/emulator...")

	// Try both potential APK locations (old and new AGP paths)
	apkPaths := []string{
		filepath.Join(a.ShellDir, "app", "build", "outputs", "apk", "debug", "app-debug.apk"),
		filepath.Join(a.ShellDir, "app", "build", "intermediates", "apk", "debug", "app-debug.apk"),
	}

	var apkPath string
	var found bool

	for _, path := range apkPaths {
		if _, err := os.Stat(path); err == nil {
			apkPath = path
			found = true
			break
		}
	}

	if !found {
		// Let's look for any APK file
		searchPath := filepath.Join(a.ShellDir, "app", "build")
		fmt.Println("Searching for APK in:", searchPath)

		err := filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && filepath.Ext(path) == ".apk" {
				apkPath = path
				found = true
				fmt.Println("Found APK at:", path)
				return filepath.SkipAll
			}
			return nil
		})

		if err != nil {
			fmt.Printf("Error searching for APK: %v\n", err)
		}
	}

	if !found {
		fmt.Println("No APK file found. Skipping installation for now.")
		return nil // Return nil to avoid failing the entire process
	}

	var args []string
	if deviceID != "" {
		// For specific device, the deviceID must come before the install command
		args = []string{"-s", deviceID, "install", "-r", apkPath}
	} else {
		args = []string{"install", "-r", apkPath}
	}

	return utils.RunCmd("adb", args...)
}

// LaunchApp launches the app on the device
func (a *Android) LaunchApp(deviceID string) error {
	fmt.Println("Launching Android app...")

	var launchArgs []string
	if deviceID != "" {
		// For specific device, the deviceID must come before the command
		launchArgs = []string{"-s", deviceID, "shell", "am", "start", "-n", "com.example.golangmobile/.MainActivity"}
	} else {
		launchArgs = []string{"shell", "am", "start", "-n", "com.example.golangmobile/.MainActivity"}
	}

	err := utils.RunCmd("adb", launchArgs...)
	if err != nil {
		fmt.Printf("Warning: Failed to launch app: %v\n", err)
		fmt.Println("Development server is still running at http://localhost:3000")
		return nil // Return nil to avoid failing the entire process
	}

	return nil
}

// SetupPortForwarding sets up port forwarding for development
func (a *Android) SetupPortForwarding(deviceID, port string) error {
	fmt.Println("Setting up port forwarding to device/emulator...")

	var args []string
	if deviceID != "" {
		// For specific device, the deviceID must come before the reverse command
		args = []string{"-s", deviceID, "reverse", fmt.Sprintf("tcp:%s", port), fmt.Sprintf("tcp:%s", port)}
	} else {
		args = []string{"reverse", fmt.Sprintf("tcp:%s", port), fmt.Sprintf("tcp:%s", port)}
	}

	return utils.RunCmd("adb", args...)
}
