package main

// #cgo pkg-config: gtk4 webkitgtk-6.0
// #include <gtk/gtk.h>
// #include "webview.h"
// #include <stdlib.h>
// extern void goActivateCallback(GtkApplication *app, gpointer user_data);
// extern void goOpenCallback(GSimpleAction *action, GVariant *parameter, gpointer user_data);
// extern void goExitCallback(GSimpleAction *action, GVariant *parameter, gpointer user_data);
//
// static void setupAppActions(GtkApplication *app, gpointer user_data) {
//     // Set up app actions
//     GSimpleAction *open_action = g_simple_action_new("open", NULL);
//     g_signal_connect(open_action, "activate", G_CALLBACK(goOpenCallback), user_data);
//     g_action_map_add_action(G_ACTION_MAP(app), G_ACTION(open_action));
//
//     GSimpleAction *exit_action = g_simple_action_new("exit", NULL);
//     g_signal_connect(exit_action, "activate", G_CALLBACK(goExitCallback), user_data);
//     g_action_map_add_action(G_ACTION_MAP(app), G_ACTION(exit_action));
// }
import "C"
import (
	"log"
	"runtime"
	"unsafe"
)

// AppState holds the application state
type AppState struct {
	app     *C.GtkApplication
	webview *WebView
}

// NewApp creates a new GTK4 application
func NewApp(appID string) *AppState {
	// Ensure we use only one OS thread for UI
	runtime.LockOSThread()

	cAppID := C.CString(appID)
	defer C.free(unsafe.Pointer(cAppID))

	// Create a new GTK application
	app := C.gtk_application_new(cAppID, C.G_APPLICATION_FLAGS_NONE)
	if app == nil {
		return nil
	}

	appState := &AppState{
		app: app,
	}

	// Register the app state
	RegisterApp(appState)

	// Connect the activate signal
	C.g_signal_connect(C.gpointer(app), C.CString("activate"),
		C.GCallback(C.goActivateCallback), C.gpointer(unsafe.Pointer(app)))

	// Set up app actions
	C.setupAppActions(app, C.gpointer(unsafe.Pointer(app)))

	return appState
}

// Run runs the GTK application
func (a *AppState) Run() int {
	// Create our webview when the app is activated
	webview := NewWebView("GTK4 WebView Example", 1024, 768, "https://golang.org")
	if webview == nil {
		log.Fatal("Failed to create webview")
	}
	a.webview = webview

	// Add context menu item
	webview.AddContextMenuItem("Custom Action")

	// Show the webview
	webview.Show()

	// Run the GTK application
	status := C.g_application_run(C.GApplication(a.app), 0, nil)
	return int(status)
}

// Cleanup cleans up resources
func (a *AppState) Cleanup() {
	// Unregister the app
	UnregisterApp(a)

	if a.webview != nil {
		a.webview.Destroy()
	}

	if a.app != nil {
		C.g_object_unref(C.gpointer(a.app))
	}
}
