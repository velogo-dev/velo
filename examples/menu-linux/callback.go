package main

// #cgo pkg-config: gtk4 webkitgtk-6.0
// #include <gtk/gtk.h>
// #include "webview.h"
// #include <stdlib.h>
// extern void goActivateCallback(GtkApplication *app, gpointer user_data);
//
// static void setupCallbacks(GtkApplication *app, gpointer user_data) {
//     goActivateCallback(app, user_data);
// }
import "C"
import (
	"fmt"
	"sync"
	"unsafe"
)

var (
	appStateMutex sync.Mutex
	appStateMap   = make(map[unsafe.Pointer]*AppState)
)

//export goActivateCallback
func goActivateCallback(app *C.GtkApplication, userData C.gpointer) {
	// Convert the pointer back to the Go AppState
	appStateMutex.Lock()
	appState, found := appStateMap[unsafe.Pointer(app)]
	appStateMutex.Unlock()

	if !found || appState == nil {
		return
	}

	// The WebView is already created in AppState.Run()
	webview := appState.webview
	if webview == nil {
		return
	}

	// Set parent window for the webview
	win := webview.handle.window
	// Associate the window with the application
	C.gtk_window_set_application((*C.GtkWindow)(unsafe.Pointer(win)), app)

	// Show the window
	C.gtk_widget_set_visible(win, C.TRUE)
}

//export goOpenCallback
func goOpenCallback(action *C.GSimpleAction, parameter C.GVariant, userData C.gpointer) {
	// Convert the pointer back to the Go AppState
	appStateMutex.Lock()
	appState, found := appStateMap[unsafe.Pointer(userData)]
	appStateMutex.Unlock()

	if !found || appState == nil || appState.webview == nil {
		return
	}

	// Handle Open action - for this example, just navigate to a new URL
	fmt.Println("Open action triggered")
	appState.webview.Navigate("https://go.dev/doc/")
}

//export goExitCallback
func goExitCallback(action *C.GSimpleAction, parameter C.GVariant, userData C.gpointer) {
	// Convert the pointer back to the Go AppState
	appStateMutex.Lock()
	appState, found := appStateMap[unsafe.Pointer(userData)]
	appStateMutex.Unlock()

	if !found || appState == nil {
		return
	}

	// Exit the application
	C.g_application_quit(C.GApplication(appState.app))
}

// RegisterApp registers an AppState with a GtkApplication
func RegisterApp(appState *AppState) {
	appStateMutex.Lock()
	defer appStateMutex.Unlock()

	// Register the app state with the application pointer
	appStateMap[unsafe.Pointer(appState.app)] = appState
}

// UnregisterApp removes an AppState registration
func UnregisterApp(appState *AppState) {
	appStateMutex.Lock()
	defer appStateMutex.Unlock()

	delete(appStateMap, unsafe.Pointer(appState.app))
}
