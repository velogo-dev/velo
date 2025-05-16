package desktop

import (
	"fmt"
	"runtime"
	"sync"
	"unsafe"

	webview "github.com/webview/webview_go"
)

// NativeMenu represents a native OS menu
type NativeMenu struct {
	Items        []*CommonMenuItem
	webView      webview.WebView
	registry     *MenuRegistry
	jsManager    *JavaScriptMenuManager
	platformMenu interface{} // Platform-specific menu handle
	mu           sync.Mutex
	hwnd         unsafe.Pointer // untuk Windows
	platform     string         // platform OS
}

// NativeMenuType represents the type of native menu
type NativeMenuType int

const (
	// MenuBarType is the main menu bar at the top of the window
	MenuBarType NativeMenuType = iota
	// ContextMenuType is a right-click context menu
	ContextMenuType
	// PopupMenuType is a popup menu
	PopupMenuType
)

// NativeMenuItem represents a native menu item
type NativeMenuItem struct {
	ID          int              // ID menu item
	Label       string           // Teks yang ditampilkan
	Shortcut    string           // Shortcut keyboard (seperti "Ctrl+S")
	Action      func()           // Callback saat menu item diklik
	SubItems    []NativeMenuItem // Sub menu items
	IsSeparator bool             // Apakah item ini adalah separator
}

// NewNativeMenu creates a new native menu
func NewNativeMenu() *NativeMenu {
	menu := &NativeMenu{
		Items:    make([]*CommonMenuItem, 0),
		registry: NewMenuRegistry(),
	}
	return menu
}

// SetWebView associates the menu with a webview
func (m *NativeMenu) SetWebView(w webview.WebView) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.webView = w
	m.jsManager = NewJavaScriptMenuManager(w, m.registry)
	m.jsManager.BindToWebView()

	// Simpan handle jendela webview untuk Windows
	if runtime.GOOS == "windows" {
		if w != nil {
			m.hwnd = w.Window()
		}
	}
}

// AddItem adds a menu item to the native menu
func (m *NativeMenu) AddItem(item *CommonMenuItem) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Items = append(m.Items, item)
	m.registry.Register(item)
}

// AddSeparator adds a separator to the menu
func (m *NativeMenu) AddSeparator() {
	m.AddItem(MenuSeparator())
}

// AddSubmenu adds a submenu to the native menu
func (m *NativeMenu) AddSubmenu(label string, items ...*CommonMenuItem) *CommonMenuItem {
	submenu := &CommonMenuItem{
		ID:       GenerateUniqueID("submenu"),
		Label:    label,
		SubItems: items,
	}

	// Register all the subitems
	for _, item := range items {
		m.registry.Register(item)
	}

	m.AddItem(submenu)
	return submenu
}

// AddFileMenu adds a standard File menu
func (m *NativeMenu) AddFileMenu(onNew, onOpen, onSave, onExit func()) {
	fileMenu := &CommonMenuItem{
		ID:       "file_menu",
		Label:    "File",
		SubItems: make([]*CommonMenuItem, 0),
	}

	if onNew != nil {
		fileMenu.SubItems = append(fileMenu.SubItems, &CommonMenuItem{
			ID:       "file_new",
			Label:    "New",
			Shortcut: "Ctrl+N",
			Action:   onNew,
		})
	}

	if onOpen != nil {
		fileMenu.SubItems = append(fileMenu.SubItems, &CommonMenuItem{
			ID:       "file_open",
			Label:    "Open",
			Shortcut: "Ctrl+O",
			Action:   onOpen,
		})
	}

	if onSave != nil {
		fileMenu.SubItems = append(fileMenu.SubItems, &CommonMenuItem{
			ID:       "file_save",
			Label:    "Save",
			Shortcut: "Ctrl+S",
			Action:   onSave,
		})
	}

	if onExit != nil {
		fileMenu.SubItems = append(fileMenu.SubItems, &CommonMenuItem{
			ID:       "file_exit",
			Label:    "Exit",
			Shortcut: "Alt+F4",
			Action:   onExit,
		})
	}

	m.AddItem(fileMenu)
}

// AddEditMenu adds a standard Edit menu
func (m *NativeMenu) AddEditMenu() {
	editMenu := &CommonMenuItem{
		ID:       "edit_menu",
		Label:    "Edit",
		SubItems: make([]*CommonMenuItem, 0),
	}

	// Add Cut
	editMenu.SubItems = append(editMenu.SubItems, &CommonMenuItem{
		ID:       "edit_cut",
		Label:    "Cut",
		Shortcut: "Ctrl+X",
		Action: func() {
			if m.webView != nil {
				m.webView.Eval("document.execCommand('cut')")
			}
		},
	})

	// Add Copy
	editMenu.SubItems = append(editMenu.SubItems, &CommonMenuItem{
		ID:       "edit_copy",
		Label:    "Copy",
		Shortcut: "Ctrl+C",
		Action: func() {
			if m.webView != nil {
				m.webView.Eval("document.execCommand('copy')")
			}
		},
	})

	// Add Paste
	editMenu.SubItems = append(editMenu.SubItems, &CommonMenuItem{
		ID:       "edit_paste",
		Label:    "Paste",
		Shortcut: "Ctrl+V",
		Action: func() {
			if m.webView != nil {
				m.webView.Eval("document.execCommand('paste')")
			}
		},
	})

	// Add separator
	editMenu.SubItems = append(editMenu.SubItems, MenuSeparator())

	// Add Select All
	editMenu.SubItems = append(editMenu.SubItems, &CommonMenuItem{
		ID:       "edit_select_all",
		Label:    "Select All",
		Shortcut: "Ctrl+A",
		Action: func() {
			if m.webView != nil {
				m.webView.Eval("document.execCommand('selectAll')")
			}
		},
	})

	m.AddItem(editMenu)
}

// Install installs the native menu to the webview
func (m *NativeMenu) Install() error {
	if m.webView == nil {
		return fmt.Errorf("webview not set for native menu")
	}

	// Implementasikan menu berdasarkan platform
	switch runtime.GOOS {
	case "windows":
		return m.installWindowsMenu()
	case "darwin":
		return m.installMacOSMenu()
	case "linux":
		return m.installLinuxMenu()
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// installWindowsMenu installs native menu on Windows
func (m *NativeMenu) installWindowsMenu() error {
	if m.hwnd == nil {
		return fmt.Errorf("window handle not available")
	}

	// Untuk Windows, kita perlu memanggil Win32 API
	// Di aplikasi sebenarnya, ini akan menggunakan syscall untuk membuat menu
	return fmt.Errorf("Windows native menu implementation not yet available")
}

// installMacOSMenu installs native menu on macOS
func (m *NativeMenu) installMacOSMenu() error {
	// Untuk macOS, kita perlu menggunakan Objective-C runtime
	// Di aplikasi sebenarnya, ini akan menggunakan cgo untuk memanggil Cocoa API
	return fmt.Errorf("MacOS native menu implementation not yet available")
}

// installLinuxMenu installs native menu on Linux
func (m *NativeMenu) installLinuxMenu() error {
	// For Linux, implement a JavaScript-based menu bar as a fallback
	// This creates a fixed menu bar at the top of the window that looks native-like
	// without requiring GTK bindings

	if m.webView == nil {
		return fmt.Errorf("webview not set for Linux menu")
	}

	// Create a JavaScript-based menu
	menuJS := `
	(function() {
		// Create menu container
		const menuBar = document.createElement('div');
		menuBar.id = 'velo-native-menu-bar';
		menuBar.style.cssText = 'position: fixed; top: 0; left: 0; right: 0; height: 30px; background: #f0f0f0; border-bottom: 1px solid #ccc; display: flex; z-index: 9999;';
		document.body.style.marginTop = '30px'; // Add margin to body to account for fixed menu
		
		// Menu item styles
		const menuItemStyle = 'padding: 0 15px; line-height: 30px; cursor: pointer; position: relative;';
		const menuItemHoverStyle = 'background-color: #ddd;';
		
		// Submenu styles
		const submenuStyle = 'position: absolute; top: 30px; left: 0; background: white; border: 1px solid #ccc; box-shadow: 2px 2px 5px rgba(0,0,0,0.2); display: none; min-width: 200px; z-index: 10000;';
		const submenuItemStyle = 'padding: 5px 15px; display: flex; justify-content: space-between;';
		
		// Add menu items
		let menuData = [];
	`

	// Add each menu item
	for _, item := range m.Items {
		menuJS += fmt.Sprintf(`
		// Add main menu: %s
		menuData.push({
			id: '%s',
			label: '%s',
			items: [
		`, item.Label, item.ID, item.Label)

		// Add subitems
		for _, subItem := range item.SubItems {
			if subItem.Type == MenuItemSeparator {
				menuJS += fmt.Sprintf(`
				{ isSeparator: true },`)
			} else {
				shortcut := ""
				if subItem.Shortcut != "" {
					shortcut = subItem.Shortcut
				}

				menuJS += fmt.Sprintf(`
				{ id: '%s', label: '%s', shortcut: '%s' },`,
					subItem.ID, subItem.Label, shortcut)
			}
		}

		menuJS += `
			]
		});
		`
	}

	// Add rendering and event handling code
	menuJS += `
		// Render the menu
		for (let menu of menuData) {
			const menuItem = document.createElement('div');
			menuItem.id = menu.id;
			menuItem.textContent = menu.label;
			menuItem.style.cssText = menuItemStyle;
			menuItem.onmouseover = function() { this.style.backgroundColor = '#ddd'; };
			menuItem.onmouseout = function() { this.style.backgroundColor = ''; };
			
			// Create submenu container
			if (menu.items && menu.items.length > 0) {
				const submenu = document.createElement('div');
				submenu.className = 'velo-submenu';
				submenu.style.cssText = submenuStyle;
				
				// Add click handler to show/hide submenu
				menuItem.onclick = function(e) {
					const allSubmenus = document.querySelectorAll('.velo-submenu');
					allSubmenus.forEach(sm => {
						if (sm !== submenu) sm.style.display = 'none';
					});
					
					if (submenu.style.display === 'block') {
						submenu.style.display = 'none';
					} else {
						submenu.style.display = 'block';
					}
					e.stopPropagation();
				};
				
				// Add submenu items
				for (let subItem of menu.items) {
					if (subItem.isSeparator) {
						const separator = document.createElement('div');
						separator.style.cssText = 'height: 1px; background: #ddd; margin: 5px 0;';
						submenu.appendChild(separator);
					} else {
						const subItemEl = document.createElement('div');
						subItemEl.style.cssText = submenuItemStyle;
						subItemEl.onmouseover = function() { this.style.backgroundColor = '#f0f0f0'; };
						subItemEl.onmouseout = function() { this.style.backgroundColor = ''; };
						
						const labelSpan = document.createElement('span');
						labelSpan.textContent = subItem.label;
						
						const shortcutSpan = document.createElement('span');
						shortcutSpan.textContent = subItem.shortcut || '';
						shortcutSpan.style.color = '#999';
						shortcutSpan.style.fontSize = '0.9em';
						
						subItemEl.appendChild(labelSpan);
						subItemEl.appendChild(shortcutSpan);
						
						// Add click handler
						subItemEl.onclick = function(e) {
							window.invokeMenuAction(subItem.id);
							submenu.style.display = 'none';
							e.stopPropagation();
						};
						
						submenu.appendChild(subItemEl);
					}
				}
				
				menuItem.appendChild(submenu);
			}
			
			menuBar.appendChild(menuItem);
		}
		
		// Add event to close all menus when clicking elsewhere
		document.addEventListener('click', function() {
			const allSubmenus = document.querySelectorAll('.velo-submenu');
			allSubmenus.forEach(sm => sm.style.display = 'none');
		});
		
		// Add the menu bar to the document
		document.body.appendChild(menuBar);
	})();
	`

	// Inject the JavaScript to create the menu
	m.webView.Eval(menuJS)

	return nil
}

// handleMenuAction handles menu action events
func (m *NativeMenu) handleMenuAction(id string) {
	if m.registry != nil {
		_ = m.registry.ExecuteAction(id)
	}
}

// Contoh Windows API types dan constants
type HMENU uintptr
type HWND uintptr

const (
	MF_POPUP     = 0x00000010
	MF_STRING    = 0x00000000
	MF_SEPARATOR = 0x00000800
)

// GetWebviewHWND returns the HWND for a webview on Windows
// This is platform-specific and will only work on Windows
func GetWebviewHWND(w webview.WebView) uintptr {
	// This is just a placeholder - the actual implementation would need to
	// access the internal window handle from the webview
	return 0
}
