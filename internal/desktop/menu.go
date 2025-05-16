package desktop

import (
	"fmt"
	"sync"

	webview "github.com/webview/webview_go"
)

// MenuItem represents a menu item in the application
type MenuItem struct {
	ID       string     // Unique identifier for this menu item
	Label    string     // Displayed text
	Shortcut string     // Keyboard shortcut (like "Ctrl+S")
	Action   func()     // Callback function when menu item is clicked
	SubItems []MenuItem // Submenu items (if any)
}

// MenuBar represents the application menu bar
type MenuBar struct {
	Items   []MenuItem
	webview webview.WebView
	mu      sync.Mutex
}

// NewMenuBar creates a new menu bar for the application
func NewMenuBar() *MenuBar {
	return &MenuBar{
		Items: []MenuItem{},
	}
}

// SetWebView associates the menu bar with a webview instance
func (m *MenuBar) SetWebView(w webview.WebView) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.webview = w
}

// AddMenuItem adds a menu item to the menu bar
func (m *MenuBar) AddMenuItem(item MenuItem) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Items = append(m.Items, item)
}

// AddStandardEditMenu adds standard Edit menu items (Cut, Copy, Paste, etc.)
func (m *MenuBar) AddStandardEditMenu() {
	editMenu := MenuItem{
		ID:    "edit",
		Label: "Edit",
		SubItems: []MenuItem{
			{ID: "cut", Label: "Cut", Shortcut: "Ctrl+X", Action: func() {
				m.ExecuteJS("document.execCommand('cut')")
			}},
			{ID: "copy", Label: "Copy", Shortcut: "Ctrl+C", Action: func() {
				m.ExecuteJS("document.execCommand('copy')")
			}},
			{ID: "paste", Label: "Paste", Shortcut: "Ctrl+V", Action: func() {
				m.ExecuteJS("document.execCommand('paste')")
			}},
			{ID: "selectall", Label: "Select All", Shortcut: "Ctrl+A", Action: func() {
				m.ExecuteJS("document.execCommand('selectAll')")
			}},
		},
	}
	m.AddMenuItem(editMenu)
}

// AddStandardFileMenu adds standard File menu items (New, Open, Save, etc.)
func (m *MenuBar) AddStandardFileMenu(onNew, onOpen, onSave, onExit func()) {
	fileMenu := MenuItem{
		ID:    "file",
		Label: "File",
		SubItems: []MenuItem{
			{ID: "new", Label: "New", Shortcut: "Ctrl+N", Action: onNew},
			{ID: "open", Label: "Open", Shortcut: "Ctrl+O", Action: onOpen},
			{ID: "save", Label: "Save", Shortcut: "Ctrl+S", Action: onSave},
			{ID: "exit", Label: "Exit", Shortcut: "Alt+F4", Action: onExit},
		},
	}
	m.AddMenuItem(fileMenu)
}

// SetupMenuHandlers initializes the menu handling in both JavaScript and Go
func (m *MenuBar) SetupMenuHandlers() error {
	if m.webview == nil {
		return fmt.Errorf("webview not set for menu bar")
	}

	// Create menu binding functions that can be called from JavaScript
	err := m.webview.Bind("__menuItemClicked", func(menuID string) {
		m.handleMenuItemClick(menuID)
	})
	if err != nil {
		return fmt.Errorf("failed to bind menu handler: %w", err)
	}

	// Generate JavaScript to create the menu representation
	js := m.generateMenuJS()
	m.webview.Init(js)

	return nil
}

// handleMenuItemClick handles menu item click events
func (m *MenuBar) handleMenuItemClick(menuID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Find the menu item with the given ID and execute its action
	for _, item := range m.Items {
		if item.ID == menuID {
			if item.Action != nil {
				item.Action()
			}
			return
		}
		// Check in submenu items
		for _, subItem := range item.SubItems {
			if subItem.ID == menuID && subItem.Action != nil {
				subItem.Action()
				return
			}
		}
	}
}

// ExecuteJS executes JavaScript code in the webview
func (m *MenuBar) ExecuteJS(js string) {
	if m.webview != nil {
		m.webview.Eval(js)
	}
}

// generateMenuJS generates JavaScript code to create the menu
func (m *MenuBar) generateMenuJS() string {
	// This creates a JS object with menu data that can be used by the frontend
	js := `
	window.__menuData = {
		items: [
	`

	for i, item := range m.Items {
		if i > 0 {
			js += ","
		}
		js += fmt.Sprintf(`
		{
			id: '%s',
			label: '%s',
			shortcut: '%s',
			subItems: [
		`, item.ID, item.Label, item.Shortcut)

		for j, subItem := range item.SubItems {
			if j > 0 {
				js += ","
			}
			js += fmt.Sprintf(`
			{
				id: '%s',
				label: '%s',
				shortcut: '%s'
			}`, subItem.ID, subItem.Label, subItem.Shortcut)
		}

		js += `
			]
		}`
	}

	js += `
		]
	};
	
	// Expose function to show native-like menu in web UI
	window.showMenu = function(menuId) {
		// Implementation can be customized by web UI
		if (window.onShowMenu) {
			window.onShowMenu(menuId);
		}
	};
	
	// Trigger menu action from web UI
	window.triggerMenuItem = function(itemId) {
		window.__menuItemClicked(itemId);
	};
	
	// Setup keyboard shortcuts
	document.addEventListener('keydown', function(e) {
		let shortcut = '';
		if (e.ctrlKey) shortcut += 'Ctrl+';
		if (e.altKey) shortcut += 'Alt+';
		if (e.shiftKey) shortcut += 'Shift+';
		
		shortcut += e.key.toUpperCase();
		
		// Check each menu item for matching shortcut
		window.__menuData.items.forEach(function(menuItem) {
			menuItem.subItems.forEach(function(subItem) {
				if (subItem.shortcut === shortcut) {
					window.__menuItemClicked(subItem.id);
					e.preventDefault();
				}
			});
		});
	});
	
	console.log('Menu system initialized');
	`

	return js
}

// GenerateCustomMenu generates HTML for menu that can be embedded in UI
func (m *MenuBar) GenerateCustomMenu() string {
	html := `
	<style>
		.velo-menu-bar {
			display: flex;
			background-color: #f5f5f5;
			border-bottom: 1px solid #ddd;
			padding: 0;
			margin: 0;
			font-family: sans-serif;
			user-select: none;
		}
		.velo-menu-item {
			position: relative;
			padding: 8px 15px;
			cursor: pointer;
		}
		.velo-menu-item:hover {
			background-color: #e8e8e8;
		}
		.velo-submenu {
			display: none;
			position: absolute;
			top: 100%;
			left: 0;
			background-color: white;
			border: 1px solid #ddd;
			box-shadow: 0 2px 5px rgba(0,0,0,0.2);
			z-index: 1000;
			min-width: 180px;
		}
		.velo-menu-item:hover .velo-submenu {
			display: block;
		}
		.velo-submenu-item {
			padding: 8px 15px;
			cursor: pointer;
			display: flex;
			justify-content: space-between;
		}
		.velo-submenu-item:hover {
			background-color: #f0f0f0;
		}
		.velo-shortcut {
			color: #666;
			font-size: 0.9em;
			margin-left: 15px;
		}
	</style>
	<div class="velo-menu-bar">
	`

	for _, item := range m.Items {
		html += fmt.Sprintf(`
		<div class="velo-menu-item" data-menu-id="%s">
			%s
			<div class="velo-submenu">
		`, item.ID, item.Label)

		for _, subItem := range item.SubItems {
			html += fmt.Sprintf(`
			<div class="velo-submenu-item" data-menu-id="%s" onclick="window.triggerMenuItem('%s')">
				<span>%s</span>
				<span class="velo-shortcut">%s</span>
			</div>
			`, subItem.ID, subItem.ID, subItem.Label, subItem.Shortcut)
		}

		html += `
			</div>
		</div>
		`
	}

	html += `</div>
	<script>
		// Initialize menu interactions
		document.addEventListener('DOMContentLoaded', function() {
			console.log('Menu UI initialized');
		});
	</script>
	`

	return html
}
