# Velo Native Menu CGO Implementation

This document provides details about the platform-specific implementations of native menus using CGO in the Velo framework.

## Overview

The Velo framework implements native menus by using CGO to bind with each platform's native UI toolkit:

- **Windows**: Win32 API
- **macOS**: Cocoa/AppKit
- **Linux**: GTK+ 3.0

Each implementation follows a similar pattern:

1. Access the native window handle from the webview
2. Create native menu structures
3. Attach the menu to the window
4. Handle menu item clicks using callbacks

## Platform-Specific Details

### macOS (Cocoa/AppKit)

The macOS implementation uses Objective-C to interact with AppKit's NSMenu and NSMenuItem classes.

**Key Components:**

- **NSMenu**: Represents a menu or submenu
- **NSMenuItem**: Represents a menu item
- **NSApplication**: The application to which the menu is attached

**CGO Setup:**

```go
/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#include <Cocoa/Cocoa.h>
*/
import "C"
```

**Implementation Notes:**

- The main menu is set using `[NSApp setMainMenu:]`
- Callbacks are implemented using a custom Objective-C class (`VeloMenuHandler`)
- Keyboard shortcuts are handled using key equivalents

### Windows (Win32 API)

The Windows implementation uses the Win32 API to create and manage menus.

**Key Components:**

- **HMENU**: Handle to a Windows menu
- **CreateMenu()**: Creates a menu bar
- **CreatePopupMenu()**: Creates a popup menu (submenu)
- **AppendMenu()**: Adds items to a menu

**CGO Setup:**

```go
/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -luser32 -lgdi32
#include <windows.h>
*/
import "C"
```

**Implementation Notes:**

- Menus are attached to windows using `SetMenu()`
- Menu events are processed through the window procedure (WndProc)
- Menu item clicks are detected via the WM_COMMAND message

### Linux (GTK+)

The Linux implementation uses GTK+ 3.0 to create and manage menus.

**Key Components:**

- **GtkMenuBar**: The top-level menu bar
- **GtkMenu**: A menu container
- **GtkMenuItem**: A menu item

**CGO Setup:**

```go
/*
#cgo pkg-config: gtk+-3.0
#include <gtk/gtk.h>
*/
import "C"
```

**Implementation Notes:**

- GTK must be initialized before creating menus
- Menu items use signal connections (g_signal_connect) for callbacks
- Menus are attached to GtkWindow containers

## Callback Mechanism

All platforms use a similar callback mechanism:

1. Go functions are registered in a map with unique IDs
2. The ID is passed to the native code when creating menu items
3. When a menu item is clicked, the native code calls back to Go with the ID
4. The Go code looks up the callback function and executes it

```go
var (
    menuCallbacks = make(map[int]func())
    menuIDCounter = 0
    menuMutex     sync.Mutex
)

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
```

## Fallback Mechanism

If the native implementation fails (e.g., can't access the window handle), the system falls back to the JavaScript-based implementation:

```go
if hwnd == nil {
    log.Println("Warning: Could not get native window handle, falling back to JavaScript implementation")
    return nm.installViaJS()
}
```

## Window Handle Access

The framework uses reflection to access the native window handle from the webview:

```go
func GetNativeWindowHandle(w webview.WebView) (unsafe.Pointer, error) {
    // Access window handle via reflection...
}
```

## Building and Dependencies

### Build Tags

Each platform-specific file uses a build tag to ensure it's only compiled on the appropriate platform:

```go
// +build darwin
```

### Dependencies

- **macOS**: Requires Xcode and the Cocoa framework
- **Windows**: Requires user32.dll and gdi32.dll
- **Linux**: Requires GTK+ 3.0 development libraries

To install GTK+ dependencies on Ubuntu:

```bash
sudo apt-get install libgtk-3-dev
```

## Future Improvements

Potential enhancements to the CGO menu implementation:

1. Support for custom drawing/styling of menu items
2. More advanced menu features (icons, custom views, etc.)
3. Context menu support with position tracking
4. Better keyboard shortcut handling
5. Improved error handling and diagnostics
