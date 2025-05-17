//go:build linux
// +build linux

package desktop

/*
#cgo pkg-config: gtk+-3.0
#include <gtk/gtk.h>
#include <stdlib.h>

// Function to initialize GTK
int initGTK() {
    if (!gtk_init_check(NULL, NULL)) {
        return 0;
    }
    return 1;
}

// Function to create a menubar
void* createMenuBar() {
    return (void*)gtk_menu_bar_new();
}

// Function to create a menu
void* createMenu() {
    return (void*)gtk_menu_new();
}

// Function to create a menu item
void* createMenuItem(const char* label) {
    return (void*)gtk_menu_item_new_with_label(label);
}

// Function to create a separator
void* createSeparator() {
    return (void*)gtk_separator_menu_item_new();
}

// Function to create a check menu item
void* createCheckMenuItem(const char* label, gboolean checked) {
    GtkWidget* item = gtk_check_menu_item_new_with_label(label);
    gtk_check_menu_item_set_active(GTK_CHECK_MENU_ITEM(item), checked);
    return (void*)item;
}

// Function to set submenu
void setSubmenu(void* item, void* submenu) {
    gtk_menu_item_set_submenu(GTK_MENU_ITEM(item), GTK_WIDGET(submenu));
}

// Function to add a menu item to a menu
void appendMenuItem(void* menu, void* item) {
    gtk_menu_shell_append(GTK_MENU_SHELL(menu), GTK_WIDGET(item));
}

// Function to set menu to a window
void setMenuToWindow(void* window, void* menu) {
    GtkWidget* vbox = gtk_box_new(GTK_ORIENTATION_VERTICAL, 0);
    gtk_box_pack_start(GTK_BOX(vbox), GTK_WIDGET(menu), FALSE, FALSE, 0);
    gtk_container_add(GTK_CONTAINER(window), vbox);
    gtk_widget_show_all(GTK_WIDGET(window));
}

// Function to handle menu item clicks
void menuItemClicked(GtkWidget* widget, gpointer data) {
    long menuID = (long)data;
    callGoMenuCallback((int)menuID);
}

// Function to connect signals
void connectMenuItemSignal(void* item, long menuID) {
    g_signal_connect(G_OBJECT(item), "activate", G_CALLBACK(menuItemClicked), (gpointer)menuID);
}

// Function to enable/disable a menu item
void setMenuItemSensitive(void* item, gboolean sensitive) {
    gtk_widget_set_sensitive(GTK_WIDGET(item), sensitive);
}

// Function to call back to Go
extern void callGoMenuCallback(int menuID);

// Utility function to get GTK window from a native window handle
void* getGtkWindowFromHandle(void* handle) {
    // This is a placeholder - implementation depends on how webview provides the handle
    return handle;
}
*/
import "C"
import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"unsafe"
)

// Ensure we're on Linux
func init() {
	if runtime.GOOS != "linux" {
		panic("menu_linux.go is only meant to be compiled on Linux")
	}

	// Initialize GTK
	if C.initGTK() == 0 {
		panic("Failed to initialize GTK")
	}
}

// Map of menu item IDs to callback functions
var (
	menuCallbacks = make(map[int]func())
	menuIDCounter = 0
	menuMutex     sync.Mutex
	menuItems     = make(map[int]unsafe.Pointer)
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

// platformSpecificInstallMenu installs a native menu on Linux using GTK
func (nm *NativeMenu) platformSpecificInstallMenu() error {
	if nm.Menu == nil {
		return fmt.Errorf("no menu to install")
	}

	if nm.parentWindow == nil {
		return fmt.Errorf("parent window not set")
	}

	// Get the GTK window from the webview
	hwnd := nm.getWindowHandle()
	if hwnd == nil {
		// Fallback to the JavaScript-based implementation if we can't get the native handle
		log.Println("Warning: Could not get native window handle, falling back to JavaScript implementation")
		return nm.installViaJS()
	}

	gtkWindow := C.getGtkWindowFromHandle(hwnd)
	if gtkWindow == nil {
		log.Println("Warning: Could not convert handle to GTK window, falling back to JavaScript implementation")
		return nm.installViaJS()
	}

	// Create a menubar
	menuBar := C.createMenuBar()
	if menuBar == nil {
		return fmt.Errorf("failed to create menu bar")
	}

	// Create and add menu items
	for _, item := range nm.Menu.Items {
		menuItem, subMenu := createGtkMenuItem(item)
		if menuItem != nil {
			if subMenu != nil {
				C.setSubmenu(unsafe.Pointer(menuItem), unsafe.Pointer(subMenu))
			}
			C.appendMenuItem(unsafe.Pointer(menuBar), unsafe.Pointer(menuItem))
		}
	}

	// Set the menu bar to the window
	C.setMenuToWindow(unsafe.Pointer(gtkWindow), unsafe.Pointer(menuBar))

	return nil
}

// createGtkMenuItem creates a GTK menu item from a Go MenuItem
func createGtkMenuItem(item *MenuItem) (unsafe.Pointer, unsafe.Pointer) {
	if item.Type == MenuItemTypeSeparator {
		return C.createSeparator(), nil
	}

	// Create the C string for the label
	cLabel := C.CString(item.Label)
	defer C.free(unsafe.Pointer(cLabel))

	var menuItem unsafe.Pointer

	// Handle different menu item types
	if item.Type == MenuItemTypeCheckbox {
		menuItem = C.createCheckMenuItem(cLabel, C.gboolean(btoi(item.Checked)))
	} else {
		menuItem = C.createMenuItem(cLabel)
	}

	// Set sensitivity (enabled/disabled)
	C.setMenuItemSensitive(menuItem, C.gboolean(btoi(item.Enabled)))

	// Register callback
	if item.Click != nil {
		callbackID := registerMenuCallback(item.Click)
		C.connectMenuItemSignal(menuItem, C.long(callbackID))

		// Store the menu item for future reference
		menuMutex.Lock()
		menuItems[callbackID] = menuItem
		menuMutex.Unlock()
	}

	// Handle submenu
	var subMenu unsafe.Pointer
	if item.Submenu != nil && len(item.Submenu.Items) > 0 {
		subMenu = C.createMenu()

		for _, subItem := range item.Submenu.Items {
			subMenuItem, subSubMenu := createGtkMenuItem(subItem)
			if subMenuItem != nil {
				if subSubMenu != nil {
					C.setSubmenu(subMenuItem, subSubMenu)
				}
				C.appendMenuItem(subMenu, subMenuItem)
			}
		}
	}

	return menuItem, subMenu
}

// getWindowHandle gets the native window handle from the webview
func (nm *NativeMenu) getWindowHandle() unsafe.Pointer {
	// This is a placeholder implementation
	// In a real implementation, you would use reflection or the webview's API
	// to get the native window handle
	return nil
}

// btoi converts a boolean to int (1 for true, 0 for false)
func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

// The installLinux method is defined in menu.go and calls this platformSpecificInstallMenu method
