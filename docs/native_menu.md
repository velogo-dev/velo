# Velo Native Menu API

The Velo framework provides a cross-platform native menu API similar to Electron's Menu API. This allows you to create and manage native application menus with consistent behavior across Windows, macOS, and Linux.

## Overview

The menu API is designed to be simple to use while providing a rich set of features:

- Create application menu bars
- Context menus
- Standard OS-specific menu roles
- Keyboard shortcuts
- Menu organization and positioning

## Core Types

### MenuItem

`MenuItem` represents a single item in a menu and has the following properties:

```go
type MenuItem struct {
    ID          string       // Unique identifier for the menu item
    Label       string       // The text to display
    Type        MenuItemType // The type of the menu item
    Role        MenuRole     // Predefined role for the menu item
    Accelerator string       // Keyboard shortcut (e.g., "Ctrl+C")
    Icon        string       // Icon path (may not be supported on all platforms)
    Enabled     bool         // Whether the menu item is enabled
    Checked     bool         // For checkbox and radio items
    Click       func()       // Callback when the menu item is clicked
    Submenu     *Menu        // Submenu for submenu type
    Position    string       // For menu positioning

    // Additional properties for positioning
    Before      string       // ID of the item to insert this item before
    After       string       // ID of the item to insert this item after
    BeforeGroupContaining string // Placement before the group containing this ID
    AfterGroupContaining  string // Placement after the group containing this ID
}
```

### Menu

`Menu` represents a collection of menu items:

```go
type Menu struct {
    Items []*MenuItem
}
```

### NativeMenu

`NativeMenu` represents a native OS menu for the application:

```go
type NativeMenu struct {
    Menu *Menu
}
```

## Creating Menus

### Using Templates

The simplest way to create a menu is using the `BuildFromTemplate` function:

```go
appMenu := desktop.BuildFromTemplate([]*desktop.MenuItemConstructorOptions{
    {
        ID:    "file",
        Label: "File",
        Submenu: desktop.BuildFromTemplate([]*desktop.MenuItemConstructorOptions{
            {ID: "new", Label: "New", Accelerator: "Ctrl+N"},
            {ID: "open", Label: "Open", Accelerator: "Ctrl+O"},
            {Type: desktop.MenuItemTypeSeparator},
            {ID: "quit", Label: "Quit", Role: desktop.RoleQuit},
        }),
    },
    {
        ID:    "edit",
        Label: "Edit",
        Submenu: desktop.BuildFromTemplate([]*desktop.MenuItemConstructorOptions{
            {ID: "undo", Label: "Undo", Role: desktop.RoleUndo},
            {ID: "redo", Label: "Redo", Role: desktop.RoleRedo},
            {Type: desktop.MenuItemTypeSeparator},
            {ID: "cut", Label: "Cut", Role: desktop.RoleCut},
            {ID: "copy", Label: "Copy", Role: desktop.RoleCopy},
            {ID: "paste", Label: "Paste", Role: desktop.RolePaste},
        }),
    },
})
```

### Creating a Default Application Menu

The framework provides a convenience function to create a standard application menu:

```go
// Create a standard application menu with platform-specific customizations
appMenu := desktop.CreateDefaultAppMenu("My App")
```

### Creating Menu Items Individually

You can also create menu items individually:

```go
menu := desktop.NewMenu()

// Add a simple menu item
menu.Append(desktop.NewMenuItem(desktop.MenuItemConstructorOptions{
    ID:    "open",
    Label: "Open...",
    Accelerator: "Ctrl+O",
    Click: func() {
        fmt.Println("Open clicked!")
    },
}))

// Add a separator
menu.Append(desktop.NewMenuItem(desktop.MenuItemConstructorOptions{
    Type: desktop.MenuItemTypeSeparator,
}))

// Add a checkbox item
menu.Append(desktop.NewMenuItem(desktop.MenuItemConstructorOptions{
    ID:      "toggle",
    Label:   "Toggle Feature",
    Type:    desktop.MenuItemTypeCheckbox,
    Checked: true,
    Click: func() {
        fmt.Println("Toggle clicked!")
    },
}))
```

## Menu Types

### MenuItemType

The `MenuItemType` enum defines the types of menu items:

```go
const (
    MenuItemTypeNormal     MenuItemType = "normal"
    MenuItemTypeSeparator  MenuItemType = "separator"
    MenuItemTypeSubmenu    MenuItemType = "submenu"
    MenuItemTypeCheckbox   MenuItemType = "checkbox"
    MenuItemTypeRadio      MenuItemType = "radio"
)
```

### MenuRole

The `MenuRole` enum defines standard OS roles for menu items:

```go
const (
    // Roles for all platforms
    RoleUndo              MenuRole = "undo"
    RoleRedo              MenuRole = "redo"
    RoleCut               MenuRole = "cut"
    RoleCopy              MenuRole = "copy"
    RolePaste             MenuRole = "paste"
    RoleDelete            MenuRole = "delete"
    RoleSelectAll         MenuRole = "selectAll"
    // ... many more roles available
)
```

## Using Native Menus in Your App

```go
// Create a new native menu
nativeMenu := desktop.NewNativeMenu()

// Set a menu
nativeMenu.Menu = appMenu

// Use with the app builder
app := desktop.NewAppBuilder().
    WithTitle("My App").
    WithSize(800, 600).
    WithNativeMenu(nativeMenu).
    Build()
```

## Platform-Specific Considerations

### macOS

- The application menu (first menu) is automatically updated to show the app name
- Standard roles like `about`, `hide`, and `services` work as expected

### Windows & Linux

- Menu bar appears in the application window
- Some roles may be implemented via JavaScript in the webview

## Context Menus

To create a context menu:

```go
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

// Show the context menu on right-click
// Note: Actual implementation depends on your application structure
app.WithBinding("showContextMenu", func() {
    // Logic to show the context menu
})

// Add JavaScript to handle context menu
app.WithInitScript(`
    document.addEventListener('contextmenu', function(e) {
        e.preventDefault();
        window.showContextMenu();
    });
`)
```

## Complete Example

See the `examples/native_menu/main.go` file for a complete example of using the native menu API.
