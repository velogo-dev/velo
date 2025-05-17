package main

import (
	"log"
	"os"
	"runtime"

	"github.com/webview/webview"
)

func simple() {
	// Ensure we're on the main thread for GTK
	runtime.LockOSThread()

	// Create a WebView
	w := webview.New(true)
	defer w.Destroy()

	// Set size and title
	w.SetSize(800, 600, webview.HintNone)
	w.SetTitle("GTK4 WebView Example")

	// Navigate to a URL
	w.Navigate("https://golang.org")

	// Add native callback for context menu items
	w.Bind("contextMenuAction", func() {
		log.Println("Context menu action triggered")
	})

	// Inject JavaScript to create a context menu
	w.Init(`
		window.addEventListener('contextmenu', function(e) {
			// Prevent default context menu
			e.preventDefault();
			
			// Call the Go callback
			contextMenuAction();
			
			// Create a custom menu
			const menu = document.createElement('div');
			menu.style.position = 'absolute';
			menu.style.left = e.pageX + 'px';
			menu.style.top = e.pageY + 'px';
			menu.style.backgroundColor = 'white';
			menu.style.border = '1px solid #ccc';
			menu.style.padding = '5px';
			menu.style.boxShadow = '2px 2px 5px rgba(0,0,0,0.2)';
			menu.style.zIndex = '1000';
			menu.innerHTML = '<div style="cursor:pointer;padding:5px;">Custom Action</div>';
			
			// Handle click
			menu.querySelector('div').addEventListener('click', function() {
				contextMenuAction();
				document.body.removeChild(menu);
			});
			
			// Close menu when clicked outside
			document.addEventListener('click', function removeMenu() {
				if (document.body.contains(menu)) {
					document.body.removeChild(menu);
				}
				document.removeEventListener('click', removeMenu);
			});
			
			document.body.appendChild(menu);
		});
	`)

	// Run the WebView
	w.Run()

	os.Exit(0)
}
