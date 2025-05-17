//go:build windows
// +build windows

package desktop

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -luser32 -lgdi32

#include <windows.h>
#include <stdlib.h>

// Windows menu item callback
extern LRESULT CALLBACK VeloWndProc(HWND hWnd, UINT message, WPARAM wParam, LPARAM lParam);

// Function to create a main menu
HMENU CreateMainMenu() {
    return CreateMenu();
}

// Function to create a popup menu (for submenus)
HMENU CreatePopupMenu() {
    return CreatePopupMenu();
}

// Function to add a menu item
BOOL AppendMenuItemA(HMENU hMenu, UINT uFlags, UINT_PTR uIDNewItem, const char* lpNewItem) {
    return AppendMenuA(hMenu, uFlags, uIDNewItem, lpNewItem);
}

// Function to add a submenu
BOOL AppendSubMenu(HMENU hMenu, UINT uFlags, HMENU hSubMenu, const char* lpNewItem) {
    return AppendMenuA(hMenu, uFlags | MF_POPUP, (UINT_PTR)hSubMenu, lpNewItem);
}

// Function to check/uncheck a menu item
BOOL CheckMenuItem(HMENU hMenu, UINT uIDCheckItem, UINT uCheck) {
    return CheckMenuItem(hMenu, uIDCheckItem, uCheck);
}

// Function to enable/disable a menu item
BOOL EnableMenuItem(HMENU hMenu, UINT uIDEnableItem, UINT uEnable) {
    return EnableMenuItem(hMenu, uIDEnableItem, uEnable);
}

// Function to set the menu to a window
BOOL SetWindowMenu(HWND hWnd, HMENU hMenu) {
    return SetMenu(hWnd, hMenu);
}

// Menu item constants
#define MF_SEPARATOR 0x00000800L
#define MF_CHECKED 0x00000008L
#define MF_UNCHECKED 0x00000000L
#define MF_ENABLED 0x00000000L
#define MF_DISABLED 0x00000002L
#define MF_STRING 0x00000000L
#define MF_POPUP 0x00000010L

// Function to handle menu events
void HandleMenuEvent(int id) {
    callGoMenuCallback(id);
}

// Function to call back to Go code
extern void callGoMenuCallback(int menuID);
*/
import "C"
import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"unsafe"
)

// Ensure we're on Windows
func init() {
	if runtime.GOOS != "windows" {
		panic("menu_windows.go is only meant to be compiled on Windows")
	}
}

// Map of menu item IDs to callback functions
var (
	menuCallbacks = make(map[int]func())
	menuIDCounter = 1000 // Start from 1000 to avoid conflicts with system IDs
	menuMutex     sync.Mutex
	windowHandle  C.HWND
)

// Register a callback function and get a unique ID
func registerMenuCallback(callback func()) int {
	if callback == nil {
		return 0 // No callback
	}

	menuMutex.Lock()
	defer menuMutex.Unlock()

	menuIDCounter++
	id := menuIDCounter
	menuCallbacks[id] = callback

	return id
}

// This function is called from C when a menu item is clicked
//
//export callGoMenuCallback
func callGoMenuCallback(menuID C.int) {
	id := int(menuID)

	menuMutex.Lock()
	callback, exists := menuCallbacks[id]
	menuMutex.Unlock()

	if exists && callback != nil {
		callback()
	}
}

// Handle Windows messages including WM_COMMAND for menu items
//
//export VeloWndProc
func VeloWndProc(hWnd C.HWND, message C.UINT, wParam C.WPARAM, lParam C.LPARAM) C.LRESULT {
	switch message {
	case C.WM_COMMAND:
		// Extract menu ID from wParam
		menuID := int(wParam) & 0xFFFF
		if menuID >= 1000 { // Our menu items start from 1000
			C.HandleMenuEvent(C.int(menuID))
			return 0
		}
	}

	// Call the default window procedure for other messages
	return C.DefWindowProcA(hWnd, message, wParam, lParam)
}

// Set the window handle for the menu
func setWindowHandle(hwnd unsafe.Pointer) {
	windowHandle = C.HWND(hwnd)
}

// platformSpecificInstallMenu installs a native menu on Windows using Win32 API
func (nm *NativeMenu) platformSpecificInstallMenu() error {
	if nm.Menu == nil {
		return fmt.Errorf("no menu to install")
	}

	if nm.parentWindow == nil {
		return fmt.Errorf("parent window not set")
	}

	// Get the window handle from the webview
	hwnd := getWindowHandleFromWebView(nm.parentWindow)
	if hwnd == nil {
		// Fallback to the JavaScript-based implementation if we can't get the native handle
		log.Println("Warning: Could not get native window handle, falling back to JavaScript implementation")
		return nm.installViaJS()
	}

	setWindowHandle(hwnd)

	// Create the main menu
	mainMenu := C.CreateMainMenu()
	if mainMenu == nil {
		return fmt.Errorf("failed to create main menu")
	}

	// Add menu items to the main menu
	for _, item := range nm.Menu.Items {
		addMenuItemToWin32Menu(mainMenu, item)
	}

	// Set the menu to the window
	if C.SetWindowMenu(windowHandle, mainMenu) == 0 {
		return fmt.Errorf("failed to set window menu")
	}

	// Force redraw of the menu
	C.DrawMenuBar(windowHandle)

	return nil
}

// addMenuItemToWin32Menu adds a menu item to a Win32 menu
func addMenuItemToWin32Menu(hMenu C.HMENU, item *MenuItem) {
	if item.Type == MenuItemTypeSeparator {
		C.AppendMenuItemA(hMenu, C.MF_SEPARATOR, 0, nil)
		return
	}

	// Create a C string for the label
	cLabel := C.CString(item.Label)
	defer C.free(unsafe.Pointer(cLabel))

	// Handle submenu
	if item.Submenu != nil && len(item.Submenu.Items) > 0 {
		// Create a popup menu for the submenu
		subMenu := C.CreatePopupMenu()

		// Add items to the submenu
		for _, subItem := range item.Submenu.Items {
			addMenuItemToWin32Menu(subMenu, subItem)
		}

		// Add the submenu to the parent menu
		C.AppendSubMenu(hMenu, C.MF_STRING, subMenu, cLabel)
	} else {
		// Regular menu item

		// Flags for the menu item
		var flags C.UINT = C.MF_STRING

		// Handle enabled/disabled state
		if !item.Enabled {
			flags |= C.MF_DISABLED
		} else {
			flags |= C.MF_ENABLED
		}

		// Handle checked state for checkable items
		if item.Type == MenuItemTypeCheckbox || item.Type == MenuItemTypeRadio {
			if item.Checked {
				flags |= C.MF_CHECKED
			} else {
				flags |= C.MF_UNCHECKED
			}
		}

		// Register the callback
		callbackID := registerMenuCallback(item.Click)

		// Add the menu item
		C.AppendMenuItemA(hMenu, flags, C.UINT_PTR(callbackID), cLabel)
	}
}

// getWindowHandleFromWebView extracts the window handle from the webview
func getWindowHandleFromWebView(w interface{}) unsafe.Pointer {
	// This is a placeholder - you need to implement this based on your webview library
	// For example, with webview_go, you might need to access a native handle field

	// This is just an example and won't work without the actual implementation:
	// return w.(webview.WebView).Window()

	// For testing, return a dummy handle
	return unsafe.Pointer(uintptr(0))
}

// Override the generic installWindows implementation
func (nm *NativeMenu) installWindows() error {
	// Use the Windows-specific implementation with Win32 API
	return nm.platformSpecificInstallMenu()
}
