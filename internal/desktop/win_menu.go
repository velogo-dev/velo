//go:build windowssss
// +build windowssss

package desktop

// #include <windows.h>
// #include <stdlib.h>
import "C"

// Windows constants for menu commands
const (
	WM_COMMAND   = 0x0111
	MF_POPUP     = 0x00000010
	MF_STRING    = 0x00000000
	MF_ENABLED   = 0x00000000
	MF_GRAYED    = 0x00000001
	MF_CHECKED   = 0x00000008
	MF_SEPARATOR = 0x00000800
)

// Windows types
type (
	HMENU uintptr
	HWND  uintptr
)

// CreateWindowsMenu creates a Windows menu from a NativeMenu
func CreateWindowsMenu(menu *NativeMenu) (HMENU, error) {
	// Create main menu bar
	menuBar := createMenu()
	if menuBar == 0 {
		return 0, fmt.Errorf("failed to create menu bar")
	}

	// Process each top-level menu item
	for _, item := range menu.Items {
		if item.Type == MenuItemSeparator {
			continue // Skip separators at the top level
		}

		// If the item has subitems, create a popup menu
		if len(item.SubItems) > 0 {
			popup := createMenu()
			if popup == 0 {
				destroyMenu(menuBar)
				return 0, fmt.Errorf("failed to create popup menu for %s", item.Label)
			}

			// Add subitems to the popup
			for i, subItem := range item.SubItems {
				// Add submenu item or separator
				if subItem.Type == MenuItemSeparator {
					if !appendMenuSeparator(popup) {
						destroyMenu(popup)
						destroyMenu(menuBar)
						return 0, fmt.Errorf("failed to add separator")
					}
				} else {
					// Use the item ID or generate a unique one
					id := subItem.ID
					if id == "" {
						id = fmt.Sprintf("menu_item_%d_%d", i, item)
					}

					if !appendMenuString(popup, id, subItem.Label) {
						destroyMenu(popup)
						destroyMenu(menuBar)
						return 0, fmt.Errorf("failed to add menu item %s", subItem.Label)
					}
				}
			}

			// Add the popup to the menu bar
			if !appendMenuPopup(menuBar, item.ID, item.Label, popup) {
				destroyMenu(popup)
				destroyMenu(menuBar)
				return 0, fmt.Errorf("failed to add popup menu %s to menu bar", item.Label)
			}
		} else {
			// Add a simple menu item to the menu bar (usually not done, but supported)
			if !appendMenuString(menuBar, item.ID, item.Label) {
				destroyMenu(menuBar)
				return 0, fmt.Errorf("failed to add menu item %s to menu bar", item.Label)
			}
		}
	}

	return menuBar, nil
}

// InstallWindowsMenu sets a Windows menu to a window
func InstallWindowsMenu(hwnd HWND, hmenu HMENU) error {
	if !setMenu(hwnd, hmenu) {
		return fmt.Errorf("failed to set menu")
	}

	if !drawMenuBar(hwnd) {
		return fmt.Errorf("failed to draw menu bar")
	}

	return nil
}

// Windows API functions

func createMenu() HMENU {
	ret, _, _ := syscall.NewLazyDLL("user32.dll").NewProc("CreateMenu").Call()
	return HMENU(ret)
}

func destroyMenu(hmenu HMENU) bool {
	ret, _, _ := syscall.NewLazyDLL("user32.dll").NewProc("DestroyMenu").Call(uintptr(hmenu))
	return ret != 0
}

func appendMenuString(hmenu HMENU, id string, text string) bool {
	textPtr, _ := syscall.UTF16PtrFromString(text)
	idInt := getMenuItemID(id)
	ret, _, _ := syscall.NewLazyDLL("user32.dll").NewProc("AppendMenuW").Call(
		uintptr(hmenu),
		MF_STRING,
		uintptr(idInt),
		uintptr(unsafe.Pointer(textPtr)),
	)
	return ret != 0
}

func appendMenuSeparator(hmenu HMENU) bool {
	ret, _, _ := syscall.NewLazyDLL("user32.dll").NewProc("AppendMenuW").Call(
		uintptr(hmenu),
		MF_SEPARATOR,
		0,
		0,
	)
	return ret != 0
}

func appendMenuPopup(hmenu HMENU, id string, text string, hSubMenu HMENU) bool {
	textPtr, _ := syscall.UTF16PtrFromString(text)
	ret, _, _ := syscall.NewLazyDLL("user32.dll").NewProc("AppendMenuW").Call(
		uintptr(hmenu),
		MF_POPUP,
		uintptr(hSubMenu),
		uintptr(unsafe.Pointer(textPtr)),
	)
	return ret != 0
}

func setMenu(hwnd HWND, hmenu HMENU) bool {
	ret, _, _ := syscall.NewLazyDLL("user32.dll").NewProc("SetMenu").Call(
		uintptr(hwnd),
		uintptr(hmenu),
	)
	return ret != 0
}

func drawMenuBar(hwnd HWND) bool {
	ret, _, _ := syscall.NewLazyDLL("user32.dll").NewProc("DrawMenuBar").Call(uintptr(hwnd))
	return ret != 0
}

// Helper functions

// Maps menu item IDs to integer identifiers for Windows
var menuItemIDMap = make(map[string]int)
var nextMenuID = 1000 // Starting ID for menu items

// getMenuItemID returns an integer ID for a menu item
func getMenuItemID(id string) int {
	if id == "" {
		return 0
	}

	if val, ok := menuItemIDMap[id]; ok {
		return val
	}

	// Assign a new ID
	menuItemIDMap[id] = nextMenuID
	nextMenuID++

	return menuItemIDMap[id]
}

// WindowsMenuHandler is a Windows-specific menu handler
func WindowsMenuHandler(hwnd HWND, msg uint32, wparam, lparam uintptr) bool {
	if msg == WM_COMMAND {
		// Extract command ID from WPARAM
		cmdID := int(wparam & 0xFFFF)

		// Find the menu item ID for this command
		for id, idInt := range menuItemIDMap {
			if idInt == cmdID {
				// TODO: Find the menu and trigger the action
				// This would require maintaining a global registry of menus
				return true
			}
		}
	}

	return false
}
