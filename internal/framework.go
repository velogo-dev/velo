package internal

import (
	"fmt"

	"github.com/velogo-dev/velo/constants"
	"github.com/velogo-dev/velo/pkg/utils"
)

type FrameworkInstaller struct {
	Library   constants.Library
	Framework constants.Framework
	AppName   string
}

// InstallFramework installs a framework and its sub-framework
func NewFrameworkInstaller(framework constants.Framework, appName string) *FrameworkInstaller {
	return &FrameworkInstaller{
		Library:   framework.Parent,
		Framework: framework,
		AppName:   appName,
	}
}

func (f *FrameworkInstaller) Install() error {
	switch f.Library {
	case constants.React:
		switch f.Framework {
		case constants.CreateReactApp:
			return f.installCreateReactApp()
		case constants.NextJS:
			return f.installNextJS()
		}
	case constants.Vue:
		switch f.Framework {
		case constants.Nuxt:
			return f.installNuxt()
		case constants.Quasar:
			return f.installQuasar()
		}
	case constants.Svelte:
		switch f.Framework {
		case constants.SvelteKit:
			return f.installSvelteKit()
		case constants.SvelteVite:
			return f.installSvelteVite()
		}
	case constants.Angular:
		switch f.Framework {
		case constants.AngularUniversal:
			return f.installAngularUniversal()
		case constants.Nest:
			return f.installNest()
		}
	case constants.Solid:
		switch f.Framework {
		case constants.SolidStart:
			return f.installSolidStart()
		case constants.SolidVite:
			return f.installSolidVite()
		}
	default:
		return fmt.Errorf("framework not supported")
	}
	return nil
}

// InstallCreateReactApp installs Create React App
func (f *FrameworkInstaller) installCreateReactApp() error {
	return utils.RunCmd("npx", "create-react-app", f.AppName)
}

// InstallNextJS installs Next.js
func (f *FrameworkInstaller) installNextJS() error {
	fmt.Println("⚙️ Installing Next.js with interactive prompts...")
	fmt.Println("✅ Follow the prompts to configure your Next.js application")
	fmt.Println("   - You can customize TypeScript, ESLint, and other options")

	// Run create-next-app with the app name and allow interactive prompts
	return utils.RunCmdWait(".", "npx", "create-next-app@latest", f.AppName)
}

// InstallNuxt installs Nuxt
func (f *FrameworkInstaller) installNuxt() error {
	fmt.Println("⚙️ Installing Nuxt with interactive prompts...")
	fmt.Println("✅ Follow the prompts to configure your Nuxt application")
	return utils.RunCmdWait(".", "npx", "nuxi@latest", "init", f.AppName)
}

// InstallQuasar installs Quasar
func (f *FrameworkInstaller) installQuasar() error {
	fmt.Println("⚙️ Installing Quasar with interactive prompts...")
	fmt.Println("✅ Follow the prompts to configure your Quasar application")
	return utils.RunCmdWait(".", "npm", "init", "quasar@latest", f.AppName)
}

// InstallSvelteKit installs SvelteKit
func (f *FrameworkInstaller) installSvelteKit() error {
	fmt.Println("⚙️ Installing SvelteKit with interactive prompts...")
	fmt.Println("✅ Follow the prompts to configure your SvelteKit application")
	return utils.RunCmdWait(".", "npm", "create", "svelte@latest", f.AppName)
}

// InstallSvelteVite installs Svelte with Vite
func (f *FrameworkInstaller) installSvelteVite() error {
	return utils.RunCmdWait(".", "npm", "create", "vite@latest", f.AppName, "--", "--template", "svelte")
}

// InstallAngularUniversal installs Angular Universal
func (f *FrameworkInstaller) installAngularUniversal() error {
	// First, we need to install Angular CLI
	if err := utils.RunCmd("npm", "install", "-g", "@angular/cli"); err != nil {
		return err
	}
	// Create a new Angular app
	if err := utils.RunCmd("ng", "new", f.AppName); err != nil {
		return err
	}
	// Add Angular Universal
	return utils.RunCmdWithDir(f.AppName, "ng", "add", "@nguniversal/express-engine")
}

// InstallNest installs Nest.js
func (f *FrameworkInstaller) installNest() error {
	fmt.Println("⚙️ Installing Nest.js CLI and setting up a new project...")
	// Install Nest CLI
	if err := utils.RunCmd("npm", "install", "-g", "@nestjs/cli"); err != nil {
		return err
	}
	fmt.Println("✅ Follow the prompts to configure your Nest.js application")
	// Create a new Nest.js project - use RunCmdWait for interactive prompts
	return utils.RunCmdWait(".", "nest", "new", f.AppName)
}

// InstallSolidStart installs SolidStart
func (f *FrameworkInstaller) installSolidStart() error {
	return utils.RunCmdWait(".", "npx", "create-solid@latest", f.AppName, "--template", "start")
}

// InstallSolidVite installs Solid with Vite
func (f *FrameworkInstaller) installSolidVite() error {
	return utils.RunCmdWait(".", "npm", "create", "vite@latest", f.AppName, "--", "--template", "solid")
}
