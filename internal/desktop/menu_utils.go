package desktop

import (
	"fmt"
	"sync"

	webview "github.com/webview/webview_go"
)

// MenuItemType represents different types of menu items
type MenuItemType int

const (
	MenuItemNormal MenuItemType = iota
	MenuItemSeparator
	MenuItemCheckbox
	MenuItemRadio
)

// MenuActionFunc is a function that gets called when a menu item is activated
type MenuActionFunc func()

// CommonMenuItem represents the common structure of menu items
// across different menu implementations
type CommonMenuItem struct {
	ID       string            // Unique identifier
	Label    string            // Display label
	Shortcut string            // Keyboard shortcut (platform dependent format)
	Type     MenuItemType      // Type of menu item
	Checked  bool              // For checkbox/radio items
	Action   MenuActionFunc    // Function to call when activated
	SubItems []*CommonMenuItem // Submenu items (if this is a submenu)
}

// MenuSeparator creates a menu separator item
func MenuSeparator() *CommonMenuItem {
	return &CommonMenuItem{
		Type: MenuItemSeparator,
	}
}

// MenuRegistry keeps track of menu items by ID
type MenuRegistry struct {
	items map[string]*CommonMenuItem
	mu    sync.RWMutex
}

// NewMenuRegistry creates a new menu registry
func NewMenuRegistry() *MenuRegistry {
	return &MenuRegistry{
		items: make(map[string]*CommonMenuItem),
	}
}

// Register adds a menu item to the registry
func (r *MenuRegistry) Register(item *CommonMenuItem) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Register this item
	if item.ID != "" {
		r.items[item.ID] = item
	}

	// Recursively register sub-items
	for _, subItem := range item.SubItems {
		r.Register(subItem)
	}
}

// Get retrieves a menu item by ID
func (r *MenuRegistry) Get(id string) (*CommonMenuItem, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.items[id]
	return item, ok
}

// ExecuteAction executes the action for a menu item with the given ID
func (r *MenuRegistry) ExecuteAction(id string) error {
	item, ok := r.Get(id)
	if !ok {
		return fmt.Errorf("menu item with ID %s not found", id)
	}

	if item.Action != nil {
		item.Action()
	}

	return nil
}

// GetAllItems returns a copy of all registered menu items
func (r *MenuRegistry) GetAllItems() map[string]*CommonMenuItem {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Create a copy to avoid concurrent access issues
	result := make(map[string]*CommonMenuItem, len(r.items))
	for id, item := range r.items {
		result[id] = item
	}

	return result
}

// JavaScriptMenuManager manages the interaction between JavaScript and native code for menus
type JavaScriptMenuManager struct {
	webView  webview.WebView
	registry *MenuRegistry
}

// NewJavaScriptMenuManager creates a new JavaScript menu manager
func NewJavaScriptMenuManager(w webview.WebView, registry *MenuRegistry) *JavaScriptMenuManager {
	return &JavaScriptMenuManager{
		webView:  w,
		registry: registry,
	}
}

// BindToWebView binds menu-related functions to the webview
func (m *JavaScriptMenuManager) BindToWebView() error {
	// Bind the menu action handler
	err := m.webView.Bind("invokeMenuAction", func(id string) {
		m.registry.ExecuteAction(id)
	})

	if err != nil {
		return fmt.Errorf("failed to bind menu action handler: %w", err)
	}

	return nil
}

// SendMenuUpdate sends a menu update to JavaScript
func (m *JavaScriptMenuManager) SendMenuUpdate(menuID string, enable bool, checked bool) {
	script := fmt.Sprintf(`
	if (window.veloMenuUpdate) {
		window.veloMenuUpdate("%s", %t, %t);
	}
	`, menuID, enable, checked)

	m.webView.Eval(script)
}

// Common platform-specific code utilities
// These are helper functions that may be useful across different menu implementations

// FormatShortcut formats a shortcut string based on the current platform
func FormatShortcut(shortcut string) string {
	// This could be expanded to convert platform-independent shortcut descriptions
	// to platform-specific formats (e.g., "Ctrl+C" to "⌘C" on macOS)
	return shortcut
}

// IsWebViewInitialized checks if a webview is valid and initialized
func IsWebViewInitialized(w webview.WebView) bool {
	return w != nil
}

// MenuAccelerator represents a keyboard accelerator for a menu item
type MenuAccelerator struct {
	Key   int  // Virtual key code
	Alt   bool // Alt key modifier
	Ctrl  bool // Ctrl key modifier
	Shift bool // Shift key modifier
	Super bool // Super/Win/Cmd key modifier
}

// ParseShortcut parses a shortcut string into a MenuAccelerator structure
// Format: [Ctrl+][Alt+][Shift+][Super+]Key
func ParseShortcut(shortcut string) *MenuAccelerator {
	// This is a placeholder and would need a proper implementation
	// based on the desired shortcut format
	return nil
}

// GenerateUniqueID generates a unique ID for a menu item
func GenerateUniqueID(prefix string) string {
	// In a real implementation, this would generate a truly unique ID
	// This is just a placeholder
	return fmt.Sprintf("%s_%d", prefix, 0)
}
