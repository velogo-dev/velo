// GTK4 WebView with Go
package main

// #cgo pkg-config: gtk4 webkitgtk-6.0
// #include <gtk/gtk.h>
// #include <webkitgtk-6.0/webkit/webkit.h>
// #include <stdlib.h>
//
// // Forward declarations of callback functions
// extern gboolean handle_context_menu(WebKitWebView *web_view, WebKitContextMenu *context_menu,
//                                   GdkEvent *event, WebKitHitTestResult *hit_test_result,
//                                   gpointer user_data);
// extern void handle_custom_action(GSimpleAction *action, GVariant *parameter, gpointer user_data);
//
// // Helper for Cgo
// static GAction* to_gaction(GSimpleAction *action) {
//     return G_ACTION(action);
// }
//
// // Helper for type casting
// static GApplication* to_application(GtkApplication *app) {
//     return G_APPLICATION(app);
// }
//
// static void activate_app(GtkApplication *app, gpointer user_data) {
//     // Create a window
//     GtkWidget *window = gtk_application_window_new(app);
//     gtk_window_set_title(GTK_WINDOW(window), "GTK4 WebView Example");
//     gtk_window_set_default_size(GTK_WINDOW(window), 1024, 768);
//
//     // Create a header bar with menu
//     GtkWidget *header = gtk_header_bar_new();
//     gtk_header_bar_set_show_title_buttons(GTK_HEADER_BAR(header), TRUE);
//     gtk_window_set_titlebar(GTK_WINDOW(window), header);
//
//     // Create menu model
//     GMenu *menu = g_menu_new();
//     GMenu *file_menu = g_menu_new();
//     g_menu_append(file_menu, "Open", "app.open");
//     g_menu_append(file_menu, "Exit", "app.exit");
//     g_menu_append_submenu(menu, "File", G_MENU_MODEL(file_menu));
//
//     // Create menu button
//     GtkWidget *menu_button = gtk_menu_button_new();
//     gtk_menu_button_set_icon_name(GTK_MENU_BUTTON(menu_button), "open-menu-symbolic");
//     gtk_menu_button_set_menu_model(GTK_MENU_BUTTON(menu_button), G_MENU_MODEL(menu));
//     gtk_header_bar_pack_end(GTK_HEADER_BAR(header), menu_button);
//
//     // Create scrolled window
//     GtkWidget *scrolled = gtk_scrolled_window_new();
//     gtk_window_set_child(GTK_WINDOW(window), scrolled);
//
//     // Create WebKit WebView
//     GtkWidget *webview = webkit_web_view_new();
//     webkit_web_view_load_uri(WEBKIT_WEB_VIEW(webview), "https://golang.org");
//     gtk_scrolled_window_set_child(GTK_SCROLLED_WINDOW(scrolled), webview);
//
//     // Connect context menu signal for custom context menu
//     g_signal_connect(webview, "context-menu", G_CALLBACK(handle_context_menu), NULL);
//
//     // Store the app in user data field for later use
//     g_object_set_data(G_OBJECT(webview), "app", app);
//
//     // Show the window
//     gtk_widget_set_visible(window, TRUE);
// }
//
// static void open_callback(GSimpleAction *action, GVariant *parameter, gpointer user_data) {
//     // In a real app, this would open a file dialog
//     g_print("Open action activated\n");
// }
//
// static void quit_app(GtkApplication *app) {
//     g_application_quit(G_APPLICATION(app));
// }
//
// static GtkApplication* create_app(const char *app_id) {
//     // Use DEFAULT_FLAGS instead of deprecated FLAGS_NONE
//     GtkApplication *app = gtk_application_new(app_id, G_APPLICATION_DEFAULT_FLAGS);
//     g_signal_connect(app, "activate", G_CALLBACK(activate_app), NULL);
//
//     // Set up actions
//     GSimpleAction *open_action = g_simple_action_new("open", NULL);
//     g_signal_connect(open_action, "activate", G_CALLBACK(open_callback), NULL);
//     g_action_map_add_action(G_ACTION_MAP(app), G_ACTION(open_action));
//
//     GSimpleAction *exit_action = g_simple_action_new("exit", NULL);
//     g_signal_connect_swapped(exit_action, "activate", G_CALLBACK(quit_app), app);
//     g_action_map_add_action(G_ACTION_MAP(app), G_ACTION(exit_action));
//
//     // Add custom action for context menu
//     GSimpleAction *custom_action = g_simple_action_new("custom", NULL);
//     g_signal_connect(custom_action, "activate", G_CALLBACK(handle_custom_action), NULL);
//     g_action_map_add_action(G_ACTION_MAP(app), G_ACTION(custom_action));
//
//     return app;
// }
import "C"
import (
	"fmt"
	"os"
	"unsafe"
)

//export handle_context_menu
func handle_context_menu(webView *C.WebKitWebView, contextMenu *C.WebKitContextMenu,
	event *C.GdkEvent, hitTestResult *C.WebKitHitTestResult,
	userData C.gpointer) C.gboolean {
	// Add a separator to the context menu
	separator := C.webkit_context_menu_item_new_separator()
	C.webkit_context_menu_append(contextMenu, separator)

	// Create a custom context menu item
	action := C.g_simple_action_new(C.CString("custom"), nil)
	defer C.g_object_unref(C.gpointer(action))

	// Use our helper to convert GSimpleAction to GAction
	gaction := C.to_gaction(action)

	item := C.webkit_context_menu_item_new_from_gaction(
		gaction,
		C.CString("Custom Action"),
		nil)

	// Add the item to the context menu
	C.webkit_context_menu_append(contextMenu, item)

	// Return FALSE to show the context menu with our addition
	return C.FALSE
}

//export handle_custom_action
func handle_custom_action(action *C.GSimpleAction, parameter *C.GVariant, userData C.gpointer) {
	fmt.Println("Custom menu action activated!")
}

func main() {
	// Create the application ID
	appID := C.CString("org.example.gtk4webview")
	defer C.free(unsafe.Pointer(appID))

	// Create and run the GTK4 application
	app := C.create_app(appID)
	defer C.g_object_unref(C.gpointer(app))

	// Run the application using our type casting helper
	gapp := C.to_application(app)
	status := C.g_application_run(gapp, 0, nil)

	os.Exit(int(status))
}
