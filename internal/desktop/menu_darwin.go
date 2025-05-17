//go:build darwin
// +build darwin

package desktop

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#include <stdlib.h>
#include <stdbool.h>

#import <Cocoa/Cocoa.h>

// Helper functions for working with Objective-C types
void* createNSMenu(const char* title) {
    NSString *menuTitle = [NSString stringWithUTF8String:title];
    NSMenu *menu = [[NSMenu alloc] initWithTitle:menuTitle];
    return (void*)menu;
}

void* createNSMenuItem(const char* title, const char* keyEquivalent, bool enabled, bool checked, int tag) {
    NSString *itemTitle = [NSString stringWithUTF8String:title];
    NSString *keyEq = [NSString stringWithUTF8String:keyEquivalent];

    NSMenuItem *menuItem = [[NSMenuItem alloc] initWithTitle:itemTitle
                                                action:@selector(handleMenuClick:)
                                                keyEquivalent:keyEq];
    [menuItem setEnabled:enabled];
    [menuItem setState:(checked ? NSControlStateValueOn : NSControlStateValueOff)];
    [menuItem setTag:tag];

    return (void*)menuItem;
}

void* createNSSeparatorItem() {
    return (void*)[NSMenuItem separatorItem];
}

void addSubmenu(void* parentMenu, void* childMenu, const char* title) {
    NSMenu *parent = (NSMenu*)parentMenu;
    NSMenu *child = (NSMenu*)childMenu;

    NSString *itemTitle = [NSString stringWithUTF8String:title];
    NSMenuItem *menuItem = [[NSMenuItem alloc] initWithTitle:itemTitle
                                                action:nil
                                                keyEquivalent:@""];
    [menuItem setSubmenu:child];
    [parent addItem:menuItem];
}

void addItemToMenu(void* menu, void* item) {
    NSMenu *nsMenu = (NSMenu*)menu;
    NSMenuItem *nsItem = (NSMenuItem*)item;
    [nsMenu addItem:nsItem];
}

void setMainMenu(void* menu) {
    NSMenu *nsMenu = (NSMenu*)menu;
    [NSApp setMainMenu:nsMenu];
}

// Define the menu handler class
@interface VeloMenuHandler : NSObject
@end

@implementation VeloMenuHandler

// This will be called when a menu item is clicked
- (void)handleMenuClick:(id)sender {
    NSMenuItem *menuItem = (NSMenuItem*)sender;
    int menuID = [menuItem tag];
    // Call back to Go here
    callGoMenuCallback(menuID);
}

@end

// Global handler instance
static VeloMenuHandler *menuHandler = nil;

void initializeMenuHandler() {
    if (menuHandler == nil) {
        menuHandler = [[VeloMenuHandler alloc] init];
    }
}

// Function to call back to Go (implemented in Go code)
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

// Ensure we're on macOS
func init() {
	if runtime.GOOS != "darwin" {
		panic("menu_darwin.go is only meant to be compiled on macOS")
	}

	// Initialize the menu handler
	C.initializeMenuHandler()
}

// Map of menu item IDs to callback functions
var (
	menuCallbacks = make(map[int]func())
	menuIDCounter = 0
	menuMutex     sync.Mutex
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

// platformSpecificInstallMenu installs a native menu on macOS using Cocoa
func (nm *NativeMenu) platformSpecificInstallMenu() error {
	if nm.Menu == nil {
		return fmt.Errorf("no menu to install")
	}

	if nm.parentWindow == nil {
		return fmt.Errorf("parent window not set")
	}

	// Get the native window handle
	hwnd := getWindowHandleFromWebView(nm.parentWindow)
	if hwnd == nil {
		// Fallback to the JavaScript-based implementation if we can't get the native handle
		log.Println("Warning: Could not get native window handle, falling back to JavaScript implementation")
		return nm.installViaJS()
	}

	mainMenu := createNSMenu("")

	// Create menu items and add them to the menu
	for _, item := range nm.Menu.Items {
		nsMenuItem := createMenuItem(item)
		C.addItemToMenu(mainMenu, nsMenuItem)
	}

	// Set as the main menu
	C.setMainMenu(mainMenu)

	return nil
}

// createNSMenu creates an NSMenu from a Go Menu
func createNSMenu(title string) unsafe.Pointer {
	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))

	return C.createNSMenu(cTitle)
}

// createMenuItem creates an NSMenuItem from a Go MenuItem
func createMenuItem(item *MenuItem) unsafe.Pointer {
	if item.Type == MenuItemTypeSeparator {
		return C.createNSSeparatorItem()
	}

	// Create the menu item
	cTitle := C.CString(item.Label)
	defer C.free(unsafe.Pointer(cTitle))

	// Convert accelerator to the macOS key equivalent format
	keyEquivalent := ""
	if item.Accelerator != "" {
		keyEquivalent = convertAcceleratorToKeyEquivalent(item.Accelerator)
	}
	cKeyEquivalent := C.CString(keyEquivalent)
	defer C.free(unsafe.Pointer(cKeyEquivalent))

	// Register the callback if present
	callbackID := registerMenuCallback(item.Click)

	nsMenuItem := C.createNSMenuItem(
		cTitle,
		cKeyEquivalent,
		C.bool(item.Enabled),
		C.bool(item.Checked),
		C.int(callbackID),
	)

	// If it has a submenu, create and attach it
	if item.Submenu != nil && len(item.Submenu.Items) > 0 {
		submenu := createNSMenu(item.Label)

		for _, subItem := range item.Submenu.Items {
			subMenuItem := createMenuItem(subItem)
			C.addItemToMenu(submenu, subMenuItem)
		}

		C.addSubmenu(nsMenuItem, submenu, cTitle)
	}

	return nsMenuItem
}

// convertAcceleratorToKeyEquivalent converts a platform-agnostic accelerator to a macOS key equivalent
func convertAcceleratorToKeyEquivalent(accelerator string) string {
	// TODO: Implement a proper conversion from "Ctrl+C" to the macOS equivalent key
	// This is a simplified example
	if accelerator == "Command+C" || accelerator == "Ctrl+C" {
		return "c"
	} else if accelerator == "Command+X" || accelerator == "Ctrl+X" {
		return "x"
	} else if accelerator == "Command+V" || accelerator == "Ctrl+V" {
		return "v"
	}

	// Default case
	return ""
}

// Override the generic installDarwin implementation
func (nm *NativeMenu) installDarwin() error {
	// Use the macOS-specific implementation with Cocoa
	return nm.platformSpecificInstallMenu()
}
