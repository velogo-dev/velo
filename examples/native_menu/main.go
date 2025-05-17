package main

import (
	"fmt"
	"log"

	"github.com/velogo-dev/velo/internal/desktop"
)

func main() {
	// Create a custom handler for About menu
	aboutHandler := func() {
		fmt.Println("About clicked!")
	}

	// Create a new native menu
	nativeMenu := desktop.NewNativeMenu()

	// Create a menu from template
	appMenu := desktop.CreateDefaultAppMenu("Velo Demo")

	// Replace the default About handler with our custom one
	if aboutItem := appMenu.GetMenuItemById("about"); aboutItem != nil {
		aboutItem.Click = aboutHandler
	}

	// Add a custom menu item to the File menu
	if fileMenu := appMenu.GetMenuItemById("file"); fileMenu != nil && fileMenu.Submenu != nil {
		// Add a separator
		fileMenu.Submenu.Append(desktop.NewMenuItem(desktop.MenuItemConstructorOptions{
			Type: desktop.MenuItemTypeSeparator,
		}))

		// Add Open menu item
		fileMenu.Submenu.Append(desktop.NewMenuItem(desktop.MenuItemConstructorOptions{
			ID:          "open",
			Label:       "Open...",
			Accelerator: "Ctrl+O",
			Click: func() {
				fmt.Println("Open clicked!")
			},
		}))

		// Add Save menu item
		fileMenu.Submenu.Append(desktop.NewMenuItem(desktop.MenuItemConstructorOptions{
			ID:          "save",
			Label:       "Save",
			Accelerator: "Ctrl+S",
			Click: func() {
				fmt.Println("Save clicked!")
			},
		}))
	}

	// Set the menu to our native menu
	nativeMenu.Menu = appMenu

	// Create a custom menu for a context menu example
	contextMenu := desktop.NewMenu()
	contextMenu.Append(desktop.NewMenuItem(desktop.MenuItemConstructorOptions{
		ID:    "cut",
		Label: "Cut",
		Role:  desktop.RoleCut,
	}))
	contextMenu.Append(desktop.NewMenuItem(desktop.MenuItemConstructorOptions{
		ID:    "copy",
		Label: "Copy",
		Role:  desktop.RoleCopy,
	}))
	contextMenu.Append(desktop.NewMenuItem(desktop.MenuItemConstructorOptions{
		ID:    "paste",
		Label: "Paste",
		Role:  desktop.RolePaste,
	}))

	// Create a simple binding to handle right-click and show the context menu
	showContextMenu := func() {
		// This would be implemented in a real app to show the context menu
		// Note: For the webview implementation, this would typically inject
		// the menu via JavaScript into the page
		fmt.Println("Context menu requested")
	}

	// Build and run the app with our custom menu
	err := desktop.NewAppBuilder().
		WithTitle("Velo Native Menu Demo").
		WithSize(800, 600).
		WithDebug(true).
		WithStaticDir("frontend/dist").
		WithNativeMenu(nativeMenu).
		WithBinding("showContextMenu", showContextMenu).
		// Add JavaScript to handle context menu
		WithInitScript(`
			document.addEventListener('contextmenu', function(e) {
				e.preventDefault();
				window.showContextMenu();
			});
		`).
		Run()

	if err != nil {
		log.Fatalf("Error running app: %v", err)
	}
}
