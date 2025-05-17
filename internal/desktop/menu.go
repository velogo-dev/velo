package desktop

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"

	webview "github.com/webview/webview_go"
)

// MenuItemType represents the type of a menu item
type MenuItemType string

const (
	// MenuItemTypeNormal is a regular menu item
	MenuItemTypeNormal MenuItemType = "normal"
	// MenuItemTypeSeparator is a separator line
	MenuItemTypeSeparator MenuItemType = "separator"
	// MenuItemTypeSubmenu is a submenu item
	MenuItemTypeSubmenu MenuItemType = "submenu"
	// MenuItemTypeCheckbox is a checkbox menu item
	MenuItemTypeCheckbox MenuItemType = "checkbox"
	// MenuItemTypeRadio is a radio menu item
	MenuItemTypeRadio MenuItemType = "radio"
)

// MenuRole represents predefined menu item roles
type MenuRole string

const (
	// Roles for all platforms
	RoleUndo             MenuRole = "undo"
	RoleRedo             MenuRole = "redo"
	RoleCut              MenuRole = "cut"
	RoleCopy             MenuRole = "copy"
	RolePaste            MenuRole = "paste"
	RoleDelete           MenuRole = "delete"
	RoleSelectAll        MenuRole = "selectAll"
	RoleReload           MenuRole = "reload"
	RoleForceReload      MenuRole = "forceReload"
	RoleToggleDevTools   MenuRole = "toggleDevTools"
	RoleResetZoom        MenuRole = "resetZoom"
	RoleZoomIn           MenuRole = "zoomIn"
	RoleZoomOut          MenuRole = "zoomOut"
	RoleToggleFullScreen MenuRole = "togglefullscreen"
	RoleMinimize         MenuRole = "minimize"
	RoleClose            MenuRole = "close"
	RoleQuit             MenuRole = "quit"

	// macOS specific roles
	RoleAbout      MenuRole = "about"
	RoleServices   MenuRole = "services"
	RoleHide       MenuRole = "hide"
	RoleHideOthers MenuRole = "hideOthers"
	RoleUnhide     MenuRole = "unhide"
	RoleFront      MenuRole = "front"
	RoleWindow     MenuRole = "window"
	RoleHelp       MenuRole = "help"
	RoleZoom       MenuRole = "zoom"

	// Custom role for application menu
	RoleAppMenu    MenuRole = "appMenu"
	RoleFileMenu   MenuRole = "fileMenu"
	RoleEditMenu   MenuRole = "editMenu"
	RoleViewMenu   MenuRole = "viewMenu"
	RoleWindowMenu MenuRole = "windowMenu"
)

// MenuItem represents a menu item in a menu
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
	Before                string // ID of the item to insert this item before
	After                 string // ID of the item to insert this item after
	BeforeGroupContaining string // Placement before the group containing this ID
	AfterGroupContaining  string // Placement after the group containing this ID
}

// MenuItemConstructorOptions is used to create a MenuItem
type MenuItemConstructorOptions struct {
	ID                    string
	Label                 string
	Type                  MenuItemType
	Role                  MenuRole
	Accelerator           string
	Icon                  string
	Enabled               bool
	Checked               bool
	Click                 func()
	Submenu               *Menu
	Position              string
	Before                string
	After                 string
	BeforeGroupContaining string
	AfterGroupContaining  string
}

// Menu represents a collection of menu items
type Menu struct {
	Items        []*MenuItem
	parentWindow webview.WebView
	mu           sync.Mutex // For thread safety
}

// Create a new menu
func NewMenu() *Menu {
	return &Menu{
		Items: make([]*MenuItem, 0),
	}
}

// MenuBar represents the top menu bar of the application
type MenuBar struct {
	Menu         *Menu
	parentWindow webview.WebView
}

// NewMenuBar creates a new menu bar
func NewMenuBar() *MenuBar {
	return &MenuBar{
		Menu: NewMenu(),
	}
}

// SetWebView sets the parent window for the menu bar
func (mb *MenuBar) SetWebView(window webview.WebView) {
	mb.parentWindow = window
}

// SetupMenuHandlers sets up the menu event handlers
func (mb *MenuBar) SetupMenuHandlers() error {
	// This would be implemented differently based on platform
	return nil
}

// GenerateCustomMenu generates HTML for a custom menu
func (mb *MenuBar) GenerateCustomMenu() string {
	// Generate HTML for the menu (used for HTML menu, not native menu)
	// This is a placeholder for the function that already exists in your code
	return "<div class='menu-placeholder'>Menu would be generated here</div>"
}

// NativeMenu represents a native OS menu for the application
type NativeMenu struct {
	Menu         *Menu
	parentWindow webview.WebView
	initialized  bool
}

// NewNativeMenu creates a new native menu
func NewNativeMenu() *NativeMenu {
	return &NativeMenu{
		Menu: NewMenu(),
	}
}

// SetWebView sets the parent window for the native menu
func (nm *NativeMenu) SetWebView(window webview.WebView) {
	nm.parentWindow = window
}

// buildFromTemplate builds a menu from a template
func buildFromTemplate(template []*MenuItemConstructorOptions) *Menu {
	menu := NewMenu()
	for _, opt := range template {
		menuItem := &MenuItem{
			ID:                    opt.ID,
			Label:                 opt.Label,
			Type:                  opt.Type,
			Role:                  opt.Role,
			Accelerator:           opt.Accelerator,
			Icon:                  opt.Icon,
			Enabled:               opt.Enabled,
			Checked:               opt.Checked,
			Click:                 opt.Click,
			Submenu:               opt.Submenu,
			Position:              opt.Position,
			Before:                opt.Before,
			After:                 opt.After,
			BeforeGroupContaining: opt.BeforeGroupContaining,
			AfterGroupContaining:  opt.AfterGroupContaining,
		}

		// Set defaults
		if menuItem.Type == "" {
			menuItem.Type = MenuItemTypeNormal
		}
		if menuItem.Enabled == false {
			// Default to enabled if not specified
			menuItem.Enabled = true
		}

		// Add to menu
		menu.append(menuItem)
	}
	return menu
}

// BuildFromTemplate creates a menu from a template
func BuildFromTemplate(template []*MenuItemConstructorOptions) *Menu {
	return buildFromTemplate(template)
}

// append adds a menu item to a menu
func (m *Menu) append(menuItem *MenuItem) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Items = append(m.Items, menuItem)
}

// Append adds a menu item to a menu
func (m *Menu) Append(menuItem *MenuItem) {
	m.append(menuItem)
}

// Insert adds a menu item at a specific position
func (m *Menu) Insert(pos int, menuItem *MenuItem) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if pos < 0 || pos > len(m.Items) {
		// If position is invalid, just append
		m.Items = append(m.Items, menuItem)
		return
	}

	// Insert at the specified position
	m.Items = append(m.Items[:pos], append([]*MenuItem{menuItem}, m.Items[pos:]...)...)
}

// GetMenuItemById finds a menu item by its ID
func (m *Menu) GetMenuItemById(id string) *MenuItem {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, item := range m.Items {
		if item.ID == id {
			return item
		}

		// If item has a submenu, search there too
		if item.Submenu != nil {
			if subItem := item.Submenu.GetMenuItemById(id); subItem != nil {
				return subItem
			}
		}
	}

	return nil
}

// Popup shows the menu as a context menu
func (m *Menu) Popup(options map[string]interface{}) {
	// This would be implemented differently based on platform
	// Not applicable for webview-based apps generally unless using specific plugins
	log.Println("Popup menu functionality is not implemented for this platform")
}

// ClosePopup closes a popup menu
func (m *Menu) ClosePopup() {
	// This would be implemented differently based on platform
	// Not applicable for webview-based apps generally unless using specific plugins
	log.Println("ClosePopup functionality is not implemented for this platform")
}

// Install installs the native menu
func (nm *NativeMenu) Install() error {
	if nm.parentWindow == nil {
		return fmt.Errorf("parent window not set")
	}

	if nm.initialized {
		return nil // Already initialized
	}

	// Install platform-specific menu
	var err error
	switch runtime.GOOS {
	case "darwin":
		err = nm.installDarwin()
	case "windows":
		err = nm.installWindows()
	case "linux":
		err = nm.installLinux()
	default:
		err = fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	if err == nil {
		nm.initialized = true
	}

	return err
}

// installDarwin installs a menu on macOS
func (nm *NativeMenu) installDarwin() error {
	// This would use macOS-specific APIs
	// For webview, this might involve Objective-C bindings or CGO

	// For now, this is just a placeholder - actual implementation would involve
	// platform-specific code using CGO to interact with macOS menus
	return nm.installViaJS()
}

// installWindows installs a menu on Windows
func (nm *NativeMenu) installWindows() error {
	// This would use Windows-specific APIs
	// For webview, this might involve Win32 API calls

	// For now, this is just a placeholder - actual implementation would involve
	// platform-specific code using CGO to interact with Windows menus
	return nm.installViaJS()
}

// installLinux installs a menu on Linux
func (nm *NativeMenu) installLinux() error {
	// This would use GTK or other Linux-specific APIs
	// For webview, this might involve GTK bindings

	// For now, this is just a placeholder - actual implementation would involve
	// platform-specific code using CGO to interact with Linux menus
	return nm.installViaJS()
}

// installViaJS is a fallback that uses JavaScript to create a menu-like experience
func (nm *NativeMenu) installViaJS() error {
	if nm.parentWindow == nil {
		return fmt.Errorf("parent window not set")
	}

	// Generate JavaScript to create a basic menu
	js := nm.generateMenuJS()

	// Execute the JavaScript
	return ExecuteJS(nm.parentWindow, js)
}

// generateMenuJS generates JavaScript code to create a menu
func (nm *NativeMenu) generateMenuJS() string {
	// Basic example that creates a menu-like interface at the top of the page
	js := `
	(function() {
		// Remove existing menu if any
		const existingMenu = document.getElementById('velo-native-menu');
		if (existingMenu) {
			existingMenu.remove();
		}
		
		// Create menu container
		const menuContainer = document.createElement('div');
		menuContainer.id = 'velo-native-menu';
		menuContainer.style.position = 'fixed';
		menuContainer.style.top = '0';
		menuContainer.style.left = '0';
		menuContainer.style.width = '100%';
		menuContainer.style.backgroundColor = '#f5f5f5';
		menuContainer.style.borderBottom = '1px solid #ddd';
		menuContainer.style.zIndex = '1000';
		menuContainer.style.display = 'flex';
		menuContainer.style.fontFamily = 'Arial, sans-serif';
		menuContainer.style.fontSize = '14px';
		
		// Function to create a menu item
		function createMenuItem(item) {
			const menuItem = document.createElement('div');
			menuItem.style.padding = '8px 16px';
			menuItem.style.cursor = 'pointer';
			menuItem.style.position = 'relative';
			menuItem.textContent = item.Label;
			
			if (!item.Enabled) {
				menuItem.style.color = '#999';
				menuItem.style.cursor = 'default';
			}
			
			// Handle click
			if (item.Enabled) {
				menuItem.addEventListener('click', function(e) {
					e.stopPropagation();
					
					// If submenu, toggle it
					if (item.Submenu) {
						const submenu = menuItem.querySelector('.velo-submenu');
						if (submenu) {
							submenu.style.display = submenu.style.display === 'block' ? 'none' : 'block';
						}
					} else {
						// Send message to Go
						window.veloMenuItemClicked(item.ID);
					}
				});
			}
			
			// Handle submenu
			if (item.Submenu && item.Submenu.Items && item.Submenu.Items.length > 0) {
				const submenu = document.createElement('div');
				submenu.className = 'velo-submenu';
				submenu.style.position = 'absolute';
				submenu.style.top = '100%';
				submenu.style.left = '0';
				submenu.style.backgroundColor = '#f5f5f5';
				submenu.style.border = '1px solid #ddd';
				submenu.style.boxShadow = '0 2px 5px rgba(0,0,0,0.1)';
				submenu.style.zIndex = '1001';
				submenu.style.display = 'none';
				submenu.style.minWidth = '150px';
				
				item.Submenu.Items.forEach(function(subItem) {
					if (subItem.Type === 'separator') {
						const separator = document.createElement('div');
						separator.style.height = '1px';
						separator.style.backgroundColor = '#ddd';
						separator.style.margin = '5px 0';
						submenu.appendChild(separator);
					} else {
						const subMenuItem = document.createElement('div');
						subMenuItem.style.padding = '8px 16px';
						subMenuItem.style.cursor = 'pointer';
						subMenuItem.textContent = subItem.Label;
						
						if (!subItem.Enabled) {
							subMenuItem.style.color = '#999';
							subMenuItem.style.cursor = 'default';
						}
						
						if (subItem.Enabled) {
							subMenuItem.addEventListener('click', function(e) {
								e.stopPropagation();
								window.veloMenuItemClicked(subItem.ID);
							});
						}
						
						// Add checkbox/radio styles if needed
						if (subItem.Type === 'checkbox' || subItem.Type === 'radio') {
							if (subItem.Checked) {
								subMenuItem.style.fontWeight = 'bold';
								subMenuItem.prepend('✓ ');
							} else {
								subMenuItem.prepend('　');
							}
						}
						
						submenu.appendChild(subMenuItem);
					}
				});
				
				menuItem.appendChild(submenu);
				
				// Handle hover
				menuItem.addEventListener('mouseenter', function() {
					const allSubmenus = document.querySelectorAll('.velo-submenu');
					allSubmenus.forEach(function(sm) {
						sm.style.display = 'none';
					});
					submenu.style.display = 'block';
				});
				
				menuItem.addEventListener('mouseleave', function(e) {
					if (!submenu.contains(e.relatedTarget)) {
						submenu.style.display = 'none';
					}
				});
			}
			
			return menuItem;
		}
		
		// Add items to the menu
	`

	// Add menu items
	js += "const menuItems = " + nm.menuItemsToJSON() + ";\n"
	js += `
		menuItems.forEach(function(item) {
			menuContainer.appendChild(createMenuItem(item));
		});
		
		// Add the menu to the document
		document.body.appendChild(menuContainer);
		
		// Add padding to body to account for menu
		const menuHeight = menuContainer.offsetHeight;
		document.body.style.paddingTop = menuHeight + 'px';
		
		// Handle click outside to close submenus
		document.addEventListener('click', function() {
			const allSubmenus = document.querySelectorAll('.velo-submenu');
			allSubmenus.forEach(function(sm) {
				sm.style.display = 'none';
			});
		});
		
		// Function to handle menu item clicks
		window.veloMenuItemClicked = function(id) {
			// This function will be bound from Go
			console.log('Menu item clicked:', id);
		};
	})();
	`

	return js
}

// menuItemsToJSON converts menu items to a JSON string for JS
func (nm *NativeMenu) menuItemsToJSON() string {
	if nm.Menu == nil || len(nm.Menu.Items) == 0 {
		return "[]"
	}

	var jsonItems []string
	for _, item := range nm.Menu.Items {
		jsonItems = append(jsonItems, menuItemToJSON(item))
	}

	return "[\n" + strings.Join(jsonItems, ",\n") + "\n]"
}

// menuItemToJSON converts a single menu item to a JSON string
func menuItemToJSON(item *MenuItem) string {
	submenuJSON := "null"
	if item.Submenu != nil && len(item.Submenu.Items) > 0 {
		var subItems []string
		for _, subItem := range item.Submenu.Items {
			subItems = append(subItems, menuItemToJSON(subItem))
		}
		submenuJSON = fmt.Sprintf(`{ "Items": [%s] }`, strings.Join(subItems, ","))
	}

	return fmt.Sprintf(`{
		"ID": %q,
		"Label": %q,
		"Type": %q,
		"Role": %q,
		"Accelerator": %q,
		"Enabled": %t,
		"Checked": %t,
		"Submenu": %s
	}`, item.ID, item.Label, item.Type, item.Role, item.Accelerator, item.Enabled, item.Checked, submenuJSON)
}

// SetApplicationMenu sets the application menu
func SetApplicationMenu(menu *Menu) {
	// This would be implemented differently based on platform
	// Not applicable for webview-based apps generally unless using specific plugins
	log.Println("SetApplicationMenu functionality requires platform-specific implementation")
}

// GetApplicationMenu gets the application menu
func GetApplicationMenu() *Menu {
	// This would be implemented differently based on platform
	// Not applicable for webview-based apps generally unless using specific plugins
	log.Println("GetApplicationMenu functionality requires platform-specific implementation")
	return nil
}

// NewMenuItem creates a new menu item
func NewMenuItem(options MenuItemConstructorOptions) *MenuItem {
	item := &MenuItem{
		ID:                    options.ID,
		Label:                 options.Label,
		Type:                  options.Type,
		Role:                  options.Role,
		Accelerator:           options.Accelerator,
		Icon:                  options.Icon,
		Enabled:               options.Enabled,
		Checked:               options.Checked,
		Click:                 options.Click,
		Submenu:               options.Submenu,
		Position:              options.Position,
		Before:                options.Before,
		After:                 options.After,
		BeforeGroupContaining: options.BeforeGroupContaining,
		AfterGroupContaining:  options.AfterGroupContaining,
	}

	// Set defaults
	if item.Type == "" {
		item.Type = MenuItemTypeNormal
	}
	if !item.Enabled {
		// Default to enabled if not specified
		item.Enabled = true
	}

	return item
}

// CreateDefaultAppMenu creates a platform-appropriate default application menu
func CreateDefaultAppMenu(appName string) *Menu {
	menu := NewMenu()

	isMac := runtime.GOOS == "darwin"

	// macOS application menu
	if isMac {
		appMenuItem := NewMenuItem(MenuItemConstructorOptions{
			ID:    "app",
			Label: appName,
			Role:  RoleAppMenu,
			Submenu: buildFromTemplate([]*MenuItemConstructorOptions{
				{ID: "about", Label: "About " + appName, Role: RoleAbout},
				{Type: MenuItemTypeSeparator},
				{ID: "services", Label: "Services", Role: RoleServices},
				{Type: MenuItemTypeSeparator},
				{ID: "hide", Label: "Hide " + appName, Role: RoleHide, Accelerator: "Command+H"},
				{ID: "hideothers", Label: "Hide Others", Role: RoleHideOthers, Accelerator: "Command+Alt+H"},
				{ID: "unhide", Label: "Show All", Role: RoleUnhide},
				{Type: MenuItemTypeSeparator},
				{ID: "quit", Label: "Quit " + appName, Role: RoleQuit, Accelerator: "Command+Q"},
			}),
		})
		menu.Append(appMenuItem)
	}

	// File menu
	fileMenuItem := NewMenuItem(MenuItemConstructorOptions{
		ID:    "file",
		Label: "File",
		Role:  RoleFileMenu,
		Submenu: buildFromTemplate([]*MenuItemConstructorOptions{
			{ID: "close", Label: getLabelForPlatform(isMac, "Close", "Exit"), Role: getRoleForPlatform(isMac, RoleClose, RoleQuit), Accelerator: getAcceleratorForPlatform(isMac, "Command+W", "Alt+F4")},
		}),
	})
	menu.Append(fileMenuItem)

	// Edit menu
	editMenuItem := NewMenuItem(MenuItemConstructorOptions{
		ID:    "edit",
		Label: "Edit",
		Role:  RoleEditMenu,
		Submenu: buildFromTemplate([]*MenuItemConstructorOptions{
			{ID: "undo", Label: "Undo", Role: RoleUndo, Accelerator: getAcceleratorForPlatform(isMac, "Command+Z", "Ctrl+Z")},
			{ID: "redo", Label: "Redo", Role: RoleRedo, Accelerator: getAcceleratorForPlatform(isMac, "Shift+Command+Z", "Ctrl+Y")},
			{Type: MenuItemTypeSeparator},
			{ID: "cut", Label: "Cut", Role: RoleCut, Accelerator: getAcceleratorForPlatform(isMac, "Command+X", "Ctrl+X")},
			{ID: "copy", Label: "Copy", Role: RoleCopy, Accelerator: getAcceleratorForPlatform(isMac, "Command+C", "Ctrl+C")},
			{ID: "paste", Label: "Paste", Role: RolePaste, Accelerator: getAcceleratorForPlatform(isMac, "Command+V", "Ctrl+V")},
			{Type: MenuItemTypeSeparator},
			{ID: "selectall", Label: "Select All", Role: RoleSelectAll, Accelerator: getAcceleratorForPlatform(isMac, "Command+A", "Ctrl+A")},
		}),
	})
	menu.Append(editMenuItem)

	// View menu
	viewMenuItem := NewMenuItem(MenuItemConstructorOptions{
		ID:    "view",
		Label: "View",
		Role:  RoleViewMenu,
		Submenu: buildFromTemplate([]*MenuItemConstructorOptions{
			{ID: "reload", Label: "Reload", Role: RoleReload, Accelerator: getAcceleratorForPlatform(isMac, "Command+R", "Ctrl+R")},
			{ID: "forcereload", Label: "Force Reload", Role: RoleForceReload, Accelerator: getAcceleratorForPlatform(isMac, "Shift+Command+R", "Ctrl+Shift+R")},
			{ID: "toggledevtools", Label: "Toggle Developer Tools", Role: RoleToggleDevTools, Accelerator: getAcceleratorForPlatform(isMac, "Alt+Command+I", "Ctrl+Shift+I")},
			{Type: MenuItemTypeSeparator},
			{ID: "resetzoom", Label: "Reset Zoom", Role: RoleResetZoom, Accelerator: getAcceleratorForPlatform(isMac, "Command+0", "Ctrl+0")},
			{ID: "zoomin", Label: "Zoom In", Role: RoleZoomIn, Accelerator: getAcceleratorForPlatform(isMac, "Command+=", "Ctrl+=")},
			{ID: "zoomout", Label: "Zoom Out", Role: RoleZoomOut, Accelerator: getAcceleratorForPlatform(isMac, "Command+-", "Ctrl+-")},
			{Type: MenuItemTypeSeparator},
			{ID: "togglefullscreen", Label: "Toggle Full Screen", Role: RoleToggleFullScreen, Accelerator: getAcceleratorForPlatform(isMac, "Ctrl+Command+F", "F11")},
		}),
	})
	menu.Append(viewMenuItem)

	// Window menu
	windowMenuItem := NewMenuItem(MenuItemConstructorOptions{
		ID:    "window",
		Label: "Window",
		Role:  RoleWindowMenu,
		Submenu: buildFromTemplate([]*MenuItemConstructorOptions{
			{ID: "minimize", Label: "Minimize", Role: RoleMinimize, Accelerator: getAcceleratorForPlatform(isMac, "Command+M", "Ctrl+M")},
			{ID: "zoom", Label: "Zoom", Role: RoleZoom},
		}),
	})

	// Add mac-specific window menu items
	if isMac {
		windowSubmenu := windowMenuItem.Submenu
		windowSubmenu.Append(NewMenuItem(MenuItemConstructorOptions{Type: MenuItemTypeSeparator}))
		windowSubmenu.Append(NewMenuItem(MenuItemConstructorOptions{ID: "front", Label: "Bring All to Front", Role: RoleFront}))
	} else {
		// Windows/Linux close window
		windowMenuItem.Submenu.Append(NewMenuItem(MenuItemConstructorOptions{ID: "close", Label: "Close", Role: RoleClose, Accelerator: "Ctrl+W"}))
	}
	menu.Append(windowMenuItem)

	// Help menu
	helpMenuItem := NewMenuItem(MenuItemConstructorOptions{
		ID:    "help",
		Label: "Help",
		Role:  RoleHelp,
		Submenu: buildFromTemplate([]*MenuItemConstructorOptions{
			{ID: "about", Label: "About " + appName, Role: RoleAbout, Click: func() {
				// This would open an About dialog
				log.Printf("Show About dialog for %s", appName)
			}},
		}),
	})
	menu.Append(helpMenuItem)

	return menu
}

// Utility functions for platform-specific options
func getLabelForPlatform(isMac bool, macLabel, otherLabel string) string {
	if isMac {
		return macLabel
	}
	return otherLabel
}

func getRoleForPlatform(isMac bool, macRole, otherRole MenuRole) MenuRole {
	if isMac {
		return macRole
	}
	return otherRole
}

func getAcceleratorForPlatform(isMac bool, macAccelerator, otherAccelerator string) string {
	if isMac {
		return macAccelerator
	}
	return otherAccelerator
}
