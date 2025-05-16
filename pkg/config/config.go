package config

import (
	"flag"
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	// Command-line flags
	DevMode      bool
	BuildAndroid bool
	BuildIOS     bool
	Preview      bool

	// Device settings
	DeviceID string

	// Directory settings
	RootDir string

	// Server settings
	PreviewPort string
	DevPort     string
}

// New creates a new configuration from command-line flags
func New() *Config {
	config := &Config{}

	// Define command-line flags
	flag.BoolVar(&config.DevMode, "dev", false, "Run in development mode with hot reload")
	flag.BoolVar(&config.BuildAndroid, "android", false, "Build Android APK")
	flag.BoolVar(&config.BuildIOS, "ios", false, "Build iOS app")
	flag.BoolVar(&config.Preview, "preview", false, "Preview on simulator/device")
	flag.StringVar(&config.DeviceID, "device", "emulator-5554", "Specify device ID for preview (optional)")
	flag.StringVar(&config.PreviewPort, "preview-port", "8080", "Port for preview server")
	flag.StringVar(&config.DevPort, "dev-port", "3001", "Port for development server")

	// Parse the flags
	flag.Parse()

	// Set root directory to current working directory
	rootDir, err := os.Getwd()
	if err != nil {
		rootDir = "."
	}
	config.RootDir = rootDir

	return config
}

// GetPlatform returns the target platform based on config
func (c *Config) GetPlatform() string {
	if c.BuildAndroid {
		return "android"
	}
	if c.BuildIOS {
		return "ios"
	}
	return "web"
}

// GetAssetsDir returns the path to the assets directory
func (c *Config) GetAssetsDir() string {
	return filepath.Join(c.RootDir, "mobile-shell", "assets")
}

// GetFrontendDir returns the path to the frontend directory
func (c *Config) GetFrontendDir() string {
	return filepath.Join(c.RootDir, "frontend")
}
