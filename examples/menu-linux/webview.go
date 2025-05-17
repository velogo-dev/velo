package main

// #cgo pkg-config: gtk4 webkitgtk-6.0
// #include "webview.h"
// #include <stdlib.h>
import "C"
import (
	"runtime"
	"unsafe"
)

// WebView is a Go wrapper for the C WebViewWindow
type WebView struct {
	handle *C.WebViewWindow
}

// NewWebView creates a new WebView instance
func NewWebView(title string, width, height int, url string) *WebView {
	// Ensure we use only one OS thread to run UI functions
	runtime.LockOSThread()

	cTitle := C.CString(title)
	cURL := C.CString(url)
	defer C.free(unsafe.Pointer(cTitle))
	defer C.free(unsafe.Pointer(cURL))

	handle := C.create_webview_window(cTitle, C.int(width), C.int(height), cURL)

	webview := &WebView{
		handle: handle,
	}

	// Set finalizer to clean up resources
	runtime.SetFinalizer(webview, (*WebView).Destroy)

	return webview
}

// SetTitle sets the title of the WebView window
func (w *WebView) SetTitle(title string) {
	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))
	C.webview_set_title(w.handle, cTitle)
}

// Navigate navigates to the specified URL
func (w *WebView) Navigate(url string) {
	cURL := C.CString(url)
	defer C.free(unsafe.Pointer(cURL))
	C.webview_navigate(w.handle, cURL)
}

// Show displays the WebView window
func (w *WebView) Show() {
	C.webview_show(w.handle)
}

// Destroy cleans up resources
func (w *WebView) Destroy() {
	if w.handle != nil {
		C.webview_destroy(w.handle)
		w.handle = nil
	}
}

// AddContextMenuItem adds a context menu item
func (w *WebView) AddContextMenuItem(label string) {
	cLabel := C.CString(label)
	defer C.free(unsafe.Pointer(cLabel))

	C.webview_add_context_menu_item(w.handle, cLabel)
}

// Run starts the main GTK loop
func (w *WebView) Run() {
	C.webview_run()
}
