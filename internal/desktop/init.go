package desktop

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/webview/webview"
)

// Config holds configuration for the desktop webview
type Config struct {
	Title       string
	Width       int
	Height      int
	Resizable   bool
	Debug       bool
	StaticDir   string
	ServerPort  int
	LocalServer bool
	URL         string
	Bindings    map[string]interface{} // JavaScript bindings
	InitScript  string                 // JavaScript to be injected on page load
	OnNavigate  func(url string)       // Callback when navigation occurs
	OnClose     func()                 // Callback when window is closed
}

// DefaultConfig returns a default configuration
func DefaultConfig() Config {
	return Config{
		Title:       "Velo App",
		Width:       800,
		Height:      600,
		Resizable:   true,
		Debug:       false,
		StaticDir:   "frontend/dist",
		ServerPort:  0, // 0 means random available port
		LocalServer: true,
		Bindings:    make(map[string]interface{}),
		InitScript:  "",
		OnNavigate:  nil,
		OnClose:     nil,
	}
}

// Run starts a webview with the given configuration
func Run(config Config) error {
	var url string

	// Start local server if requested
	if config.LocalServer {
		port, err := startLocalServer(config.StaticDir, config.ServerPort)
		if err != nil {
			return fmt.Errorf("failed to start local server: %w", err)
		}
		url = fmt.Sprintf("http://localhost:%d", port)
		config.URL = url
		log.Printf("Local server started at %s", url)
	} else if config.URL == "" {
		return fmt.Errorf("either LocalServer must be true or URL must be provided")
	} else {
		url = config.URL
	}

	// Create webview
	w := webview.New(config.Debug)
	defer w.Destroy()

	// Configure webview
	w.SetTitle(config.Title)
	if config.Resizable {
		w.SetSize(config.Width, config.Height, webview.HintNone)
	} else {
		w.SetSize(config.Width, config.Height, webview.HintFixed)
	}

	// Set initial JavaScript if provided
	if config.InitScript != "" {
		w.Init(config.InitScript)
	}

	// Setup navigation monitoring
	if config.OnNavigate != nil {
		// We'll inject a script to notify us when navigation occurs
		navigationMonitorScript := `
		(function() {
			let currentUrl = window.location.href;
			function checkUrl() {
				if (window.location.href !== currentUrl) {
					currentUrl = window.location.href;
					window.notifyNavigationChange(currentUrl);
				}
				setTimeout(checkUrl, 100);
			}
			checkUrl();
		})();
		`
		w.Bind("notifyNavigationChange", func(url string) {
			if config.OnNavigate != nil {
				config.OnNavigate(url)
			}
		})
		w.Init(navigationMonitorScript)
	}

	// Register JavaScript bindings
	for name, fn := range config.Bindings {
		err := w.Bind(name, fn)
		if err != nil {
			return fmt.Errorf("failed to bind JavaScript function %s: %w", name, err)
		}
	}

	// Navigate to the URL
	w.Navigate(url)

	// Run the webview
	w.Run()

	// Call OnClose callback if provided
	if config.OnClose != nil {
		config.OnClose()
	}

	return nil
}

// RunWithDefaultConfig runs a webview with default configuration
func RunWithDefaultConfig() error {
	return Run(DefaultConfig())
}

// startLocalServer starts a local HTTP server serving static files
func startLocalServer(staticDir string, port int) (int, error) {
	// Check if the static directory exists
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		return 0, fmt.Errorf("static directory %s does not exist", staticDir)
	}

	// Create a file server handler
	fs := http.FileServer(http.Dir(staticDir))

	// Create a custom handler that adds necessary headers
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers if needed
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Add security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Let the file server serve the request
		fs.ServeHTTP(w, r)
	})

	// Register the handler
	http.Handle("/", handler)

	// Create listener
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return 0, fmt.Errorf("failed to create listener: %w", err)
	}

	// Get the actual port
	actualPort := listener.Addr().(*net.TCPAddr).Port

	// Start the server in a goroutine
	go func() {
		err := http.Serve(listener, nil)
		if err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	return actualPort, nil
}

// ExecuteJS executes JavaScript code in a webview
func ExecuteJS(w webview.WebView, js string) error {
	if w == nil {
		return fmt.Errorf("webview is nil")
	}

	// Sanitize the JavaScript code to prevent injection
	js = strings.ReplaceAll(js, "\\", "\\\\")
	js = strings.ReplaceAll(js, "'", "\\'")
	js = strings.ReplaceAll(js, "\n", " ")

	w.Eval(js)
	return nil
}

// RunApp runs a webview with a specific app configuration
func RunApp(appName string, staticDir string, debug bool) error {
	config := DefaultConfig()
	config.Title = appName
	config.StaticDir = staticDir
	config.Debug = debug

	return Run(config)
}

// GetAbsolutePath returns the absolute path to a directory
func GetAbsolutePath(relativePath string) (string, error) {
	return filepath.Abs(relativePath)
}

// AppBuilder allows for fluent configuration of a desktop app
type AppBuilder struct {
	config Config
}

// NewAppBuilder creates a new AppBuilder with default configuration
func NewAppBuilder() *AppBuilder {
	return &AppBuilder{
		config: DefaultConfig(),
	}
}

// WithTitle sets the app title
func (b *AppBuilder) WithTitle(title string) *AppBuilder {
	b.config.Title = title
	return b
}

// WithSize sets the window size
func (b *AppBuilder) WithSize(width, height int) *AppBuilder {
	b.config.Width = width
	b.config.Height = height
	return b
}

// WithResizable sets whether the window is resizable
func (b *AppBuilder) WithResizable(resizable bool) *AppBuilder {
	b.config.Resizable = resizable
	return b
}

// WithDebug sets debug mode
func (b *AppBuilder) WithDebug(debug bool) *AppBuilder {
	b.config.Debug = debug
	return b
}

// WithStaticDir sets the static directory
func (b *AppBuilder) WithStaticDir(dir string) *AppBuilder {
	b.config.StaticDir = dir
	return b
}

// WithURL sets a custom URL to navigate to
func (b *AppBuilder) WithURL(url string) *AppBuilder {
	b.config.URL = url
	b.config.LocalServer = false
	return b
}

// WithLocalServer sets whether to use a local server
func (b *AppBuilder) WithLocalServer(use bool) *AppBuilder {
	b.config.LocalServer = use
	return b
}

// WithPort sets the server port
func (b *AppBuilder) WithPort(port int) *AppBuilder {
	b.config.ServerPort = port
	return b
}

// WithBinding adds a JavaScript binding
func (b *AppBuilder) WithBinding(name string, fn interface{}) *AppBuilder {
	b.config.Bindings[name] = fn
	return b
}

// WithInitScript sets the initial JavaScript to run
func (b *AppBuilder) WithInitScript(script string) *AppBuilder {
	b.config.InitScript = script
	return b
}

// WithOnNavigate sets the navigation callback
func (b *AppBuilder) WithOnNavigate(callback func(url string)) *AppBuilder {
	b.config.OnNavigate = callback
	return b
}

// WithOnClose sets the close callback
func (b *AppBuilder) WithOnClose(callback func()) *AppBuilder {
	b.config.OnClose = callback
	return b
}

// Build returns the final config
func (b *AppBuilder) Build() Config {
	return b.config
}

// Run runs the app with the built configuration
func (b *AppBuilder) Run() error {
	return Run(b.config)
}
