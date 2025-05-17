#include "webview.h"
#include <stdlib.h>
#include <stdio.h>

static void on_activate_exit(GtkWidget *widget, gpointer data) {
    g_application_quit(G_APPLICATION(data));
}

static void on_window_destroy(GtkWidget *widget, gpointer data) {
    // Do nothing - we'll use application lifecycle
}

WebViewWindow *create_webview_window(const char *title, int width, int height, const char *url) {
    WebViewWindow *wv = (WebViewWindow *)malloc(sizeof(WebViewWindow));
    if (!wv) {
        fprintf(stderr, "Memory allocation failed\n");
        return NULL;
    }
    
    // Create window
    wv->window = gtk_window_new();
    gtk_window_set_title(GTK_WINDOW(wv->window), title);
    gtk_window_set_default_size(GTK_WINDOW(wv->window), width, height);
    g_signal_connect(wv->window, "destroy", G_CALLBACK(on_window_destroy), NULL);
    
    // Create main vertical box
    GtkWidget *vbox = gtk_box_new(GTK_ORIENTATION_VERTICAL, 0);
    gtk_window_set_child(GTK_WINDOW(wv->window), vbox);
    
    // ---- Create header bar ----
    GtkWidget *header_bar = gtk_header_bar_new();
    gtk_header_bar_set_show_title_buttons(GTK_HEADER_BAR(header_bar), TRUE);
    gtk_window_set_titlebar(GTK_WINDOW(wv->window), header_bar);
    
    // Create a menu model
    GMenu *menu_model = g_menu_new();
    GMenu *file_menu = g_menu_new();
    
    // Add menu items to the File menu
    g_menu_append(file_menu, "Open", "app.open");
    g_menu_append(file_menu, "Exit", "app.exit");
    
    // Add File menu to the main menu
    g_menu_append_submenu(menu_model, "File", G_MENU_MODEL(file_menu));
    
    // Create a menu button in the header bar
    GtkWidget *menu_button = gtk_menu_button_new();
    GtkWidget *hamburger = gtk_image_new_from_icon_name("open-menu-symbolic");
    // Use set_child instead of set_icon_widget for GTK4
    gtk_menu_button_set_child(GTK_MENU_BUTTON(menu_button), hamburger);
    gtk_menu_button_set_menu_model(GTK_MENU_BUTTON(menu_button), G_MENU_MODEL(menu_model));
    gtk_header_bar_pack_end(GTK_HEADER_BAR(header_bar), menu_button);
    
    // Store the menu bar (not needed but keeping for API compatibility)
    wv->menubar = header_bar;
    
    // Set up context (popup) menu - using GMenu for context menu
    wv->popup_menu = G_MENU(g_menu_new());
    
    // Create scrolled window
    GtkWidget *scrolled_window = gtk_scrolled_window_new();
    gtk_box_append(GTK_BOX(vbox), scrolled_window);
    gtk_widget_set_vexpand(scrolled_window, TRUE);
    
    // Create WebKit WebView
    wv->webview = webkit_web_view_new();
    webkit_web_view_load_uri(WEBKIT_WEB_VIEW(wv->webview), url);
    gtk_scrolled_window_set_child(GTK_SCROLLED_WINDOW(scrolled_window), wv->webview);

    return wv;
}

void webview_set_title(WebViewWindow *wv, const char *title) {
    gtk_window_set_title(GTK_WINDOW(wv->window), title);
}

void webview_navigate(WebViewWindow *wv, const char *url) {
    webkit_web_view_load_uri(WEBKIT_WEB_VIEW(wv->webview), url);
}

void webview_show(WebViewWindow *wv) {
    gtk_widget_set_visible(wv->window, TRUE);
}

void webview_destroy(WebViewWindow *wv) {
    gtk_window_destroy(GTK_WINDOW(wv->window));
    free(wv);
}

void webview_run() {
    // GTK4 requires an application instance to run 
    // This is expected to be handled by the Go side
}

static gboolean on_context_menu(WebKitWebView *web_view, 
                                WebKitContextMenu *context_menu, 
                                GdkEvent *event, 
                                WebKitHitTestResult *hit_test_result, 
                                gpointer user_data) {
    // Add custom context menu items here
    WebViewWindow *wv = (WebViewWindow *)user_data;
    
    if (wv->popup_menu) {
        // Add a separator
        WebKitContextMenuItem *separator = webkit_context_menu_item_new_separator();
        webkit_context_menu_append(context_menu, separator);
        
        // Add our custom menu items from the GMenu
        int n_items = g_menu_model_get_n_items(G_MENU_MODEL(wv->popup_menu));
        for (int i = 0; i < n_items; i++) {
            const char *label = NULL;
            g_menu_model_get_item_attribute(G_MENU_MODEL(wv->popup_menu), i, 
                                           G_MENU_ATTRIBUTE_LABEL, "s", &label);
            if (label) {
                // Create a simple action for the context menu with WebKit API
                GSimpleAction *action = g_simple_action_new(label, NULL);
                WebKitContextMenuItem *item = 
                    webkit_context_menu_item_new_from_gaction(G_ACTION(action), label, NULL);
                webkit_context_menu_append(context_menu, item);
                g_object_unref(action);
            }
        }
    }
    
    return FALSE; // Return FALSE to show the context menu with our additions
}

void webview_add_context_menu_item(WebViewWindow *wv, const char *label) {
    // Initialize popup menu if needed
    if (!wv->popup_menu) {
        wv->popup_menu = G_MENU(g_menu_new());
        g_signal_connect(wv->webview, "context-menu", 
                         G_CALLBACK(on_context_menu), wv);
    }
    
    // Add the item to our GMenu
    g_menu_append(wv->popup_menu, label, "context.custom");
} 