package app

import (
	"fmt"
	"log"
	"sync"

	"github.com/velogo-dev/velo/pkg/builder"
	"github.com/velogo-dev/velo/pkg/config"
	"github.com/velogo-dev/velo/pkg/server"
)

// App represents the main application
type App struct {
	Config   *config.Config
	Frontend *builder.Frontend
	Android  *builder.Android
	IOS      *builder.IOS
	Server   *server.PreviewServer
}

// New creates a new application instance
func New(cfg *config.Config) *App {
	return &App{
		Config:   cfg,
		Frontend: builder.NewFrontend(cfg.RootDir),
		Android:  builder.NewAndroid(cfg.RootDir),
		IOS:      builder.NewIOS(cfg.RootDir),
		Server:   server.NewPreviewServer(cfg.RootDir, cfg.PreviewPort),
	}
}

// Run executes the application based on the configuration
func (a *App) Run() error {

	// Run in development mode
	if a.Config.DevMode {
		return a.runDevMode()
	}

	// Run in production/build mode
	return a.runBuildMode()
}

// runDevMode runs the application in development mode
func (a *App) runDevMode() error {
	fmt.Println("Running in development mode...")

	// Development with WebView
	var wg sync.WaitGroup
	wg.Add(2) // Increased to 2 to account for the dev server goroutine

	// Start frontend dev server asynchronously
	go func() {
		defer wg.Done()
		if err := a.Frontend.StartDevServer(); err != nil {
			log.Printf("Failed to start dev server: %v", err)
		}
	}()

	// If previewing on device/emulator
	if a.Config.Preview {
		if err := a.setupDevicePreview(); err != nil {
			log.Printf("Warning: Preview setup issue: %v", err)
			// Even if preview setup fails, we'll continue running the dev server
			fmt.Println("Development server is still running at http://localhost:3000")
		}
	} else {
		fmt.Println("Starting preview server...")
		// Just serve the assets locally for quick preview
		a.Server.StartBackground()
	}

	// Keep the program running
	wg.Wait()
	return nil
}

// setupDevicePreview sets up preview on a device
func (a *App) setupDevicePreview() error {
	fmt.Println("Setting up device preview...")

	// Set up port forwarding for Android
	if a.Config.DeviceID != "" && a.Config.BuildAndroid {
		if err := a.Android.SetupPortForwarding(a.Config.DeviceID, a.Config.DevPort); err != nil {
			fmt.Printf("Warning: Port forwarding failed: %v\n", err)
			// Continue anyway, as it might still work
		}
	}

	// Build the apps first when in dev mode
	if a.Config.BuildAndroid {
		// Build Android app before installation
		if err := a.Android.Build(); err != nil {
			fmt.Printf("Warning: Android build failed: %v\n", err)
			fmt.Println("You can access the development server directly at http://localhost:3000")
			return nil // Return nil to continue with the app
		}
		if err := a.launchAndroidPreview(); err != nil {
			fmt.Printf("Warning: Android preview launch failed: %v\n", err)
			fmt.Println("You can access the development server directly at http://localhost:3000")
			return nil // Return nil to continue with the app
		}
	} else if a.Config.BuildIOS {
		// Build iOS app before installation
		if err := a.IOS.Build(); err != nil {
			fmt.Printf("Warning: iOS build failed: %v\n", err)
			fmt.Println("You can access the development server directly at http://localhost:3000")
			return nil // Return nil to continue with the app
		}
		if err := a.launchIOSPreview(); err != nil {
			fmt.Printf("Warning: iOS preview launch failed: %v\n", err)
			fmt.Println("You can access the development server directly at http://localhost:3000")
			return nil // Return nil to continue with the app
		}
	}

	return nil
}

// launchAndroidPreview installs and launches the app on an Android device
func (a *App) launchAndroidPreview() error {
	if err := a.Android.InstallApp(a.Config.DeviceID); err != nil {
		return fmt.Errorf("android app installation failed: %w", err)
	}

	if err := a.Android.LaunchApp(a.Config.DeviceID); err != nil {
		return fmt.Errorf("android app launch failed: %w", err)
	}

	return nil
}

// launchIOSPreview installs and launches the app on an iOS device
func (a *App) launchIOSPreview() error {
	if err := a.IOS.InstallApp(a.Config.DeviceID); err != nil {
		return fmt.Errorf("iOS app installation failed: %w", err)
	}

	if err := a.IOS.LaunchApp(a.Config.DeviceID); err != nil {
		return fmt.Errorf("iOS app launch failed: %w", err)
	}

	return nil
}

// runBuildMode builds the application for production
func (a *App) runBuildMode() error {
	return a.buildWithWebView()
}

// buildWithWebView builds the app using WebView approach
func (a *App) buildWithWebView() error {
	// 1. Build frontend
	if err := a.Frontend.Build(); err != nil {
		return fmt.Errorf("frontend build failed: %w", err)
	}

	// 2. Copy build output to mobile shell
	if err := a.Frontend.CopyBuildToMobile(); err != nil {
		return fmt.Errorf("copying build output failed: %w", err)
	}

	// 3. Build mobile apps if requested
	if a.Config.BuildAndroid {
		if err := a.Android.Build(); err != nil {
			return fmt.Errorf("android build failed: %w", err)
		}
	}

	if a.Config.BuildIOS {
		if err := a.IOS.Build(); err != nil {
			return fmt.Errorf("iOS build failed: %w", err)
		}
	}

	// 4. Preview on device/simulator if requested
	if a.Config.Preview {
		if a.Config.BuildAndroid {
			if err := a.launchAndroidPreview(); err != nil {
				return fmt.Errorf("android preview failed: %w", err)
			}
		} else if a.Config.BuildIOS {
			if err := a.launchIOSPreview(); err != nil {
				return fmt.Errorf("iOS preview failed: %w", err)
			}
		}
	}

	fmt.Println("Build process completed successfully!")
	return nil
}
