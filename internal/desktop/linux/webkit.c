#include <gtk/gtk.h>
#include <webkitgtk-4.0/webkit2/webkit2.h>
#include "webkit.h"

void destroy_window(GtkWidget* widget, GtkWidget* window) {
    gtk_main_quit();
}

void go_back_cb(GtkWidget* widget, WebKitWebView* web_view) {
    webkit_web_view_go_back(web_view);
}

void go_forward_cb(GtkWidget* widget, WebKitWebView* web_view) {
    webkit_web_view_go_forward(web_view);
}

void refresh_cb(GtkWidget* widget, WebKitWebView* web_view) {
    webkit_web_view_reload(web_view);
}

void zoom_in_cb(GtkWidget* widget, WebKitWebView* web_view) {
    webkit_web_view_set_zoom_level(web_view, webkit_web_view_get_zoom_level(web_view) + 0.1);
}

void zoom_out_cb(GtkWidget* widget, WebKitWebView* web_view) {
    webkit_web_view_set_zoom_level(web_view, webkit_web_view_get_zoom_level(web_view) - 0.1);
}

void zoom_reset_cb(GtkWidget* widget, WebKitWebView* web_view) {
    webkit_web_view_set_zoom_level(web_view, 1.0);
}

void init(int* argc, char** argv[]) {
    gtk_init(argc, argv);
}

GtkWidget* new_window(GtkWidget* widget, GtkWidget* parent_window) {
    GtkWidget* window = gtk_window_new(GTK_WINDOW_TOPLEVEL);
    return window;
}

GtkWidget* set_window_size(GtkWidget* window, int width, int height) {
    gtk_window_set_default_size(GTK_WINDOW(window), width, height);
    return window;
}

GtkWidget* set_window_title(GtkWidget* window, char* title) {
    gtk_window_set_title(GTK_WINDOW(window), title);
    return window;
}

GtkWidget* set_window_icon(GtkWidget* window, char* icon) {
    GError* error = NULL;
    gtk_window_set_icon_from_file(GTK_WINDOW(window), icon, &error);
    if (error) {
        g_error_free(error);
    }
    return window;
}

void run_main_loop(GtkWidget* window) {
    gtk_widget_show_all(window);
    gtk_main();
}

GtkWidget* initialize_gtk_window(int argc, char* argv[], int width, int height, const char* title, const char* icon_path) {
    // Initialize GTK
    gtk_init(&argc, &argv);

    // Create the main window
    GtkWidget* window = gtk_window_new(GTK_WINDOW_TOPLEVEL);
    gtk_window_set_default_size(GTK_WINDOW(window), width, height);
    // Set Title
    gtk_window_set_title(GTK_WINDOW(window), title);
    g_signal_connect(window, "destroy", G_CALLBACK(destroy_window), window);
    // Set Icon
    gtk_window_set_icon_name(GTK_WINDOW(window), "Pineapple");
    // Set Icon from file
    if (icon_path != NULL) {
        GError* error = NULL;
        gtk_window_set_icon_from_file(GTK_WINDOW(window), icon_path, &error);
        if (error) {
            g_error_free(error);
        }
    }
    
    return window;
}

WebKitWebView* create_web_view() {
    // Create a WebKit web view
    WebKitWebView* web_view = WEBKIT_WEB_VIEW(webkit_web_view_new());
    return web_view;
}

GtkWidget* create_zoom_menu(WebKitWebView* web_view) {
    // Create Zoom submenu
    GtkWidget* zoom_menu = gtk_menu_new();
    
    // Add zoom items to zoom submenu
    GtkWidget* zoom_in_item = gtk_menu_item_new_with_label("Zoom In");
    g_signal_connect(zoom_in_item, "activate", G_CALLBACK(zoom_in_cb), web_view);
    gtk_menu_shell_append(GTK_MENU_SHELL(zoom_menu), zoom_in_item);
    
    GtkWidget* zoom_out_item = gtk_menu_item_new_with_label("Zoom Out");
    g_signal_connect(zoom_out_item, "activate", G_CALLBACK(zoom_out_cb), web_view);
    gtk_menu_shell_append(GTK_MENU_SHELL(zoom_menu), zoom_out_item);
    
    GtkWidget* zoom_reset_item = gtk_menu_item_new_with_label("Reset Zoom");
    g_signal_connect(zoom_reset_item, "activate", G_CALLBACK(zoom_reset_cb), web_view);
    gtk_menu_shell_append(GTK_MENU_SHELL(zoom_menu), zoom_reset_item);
    
    return zoom_menu;
}

GtkWidget* create_navigation_menu(WebKitWebView* web_view) {
    // Create "Navigation" menu
    GtkWidget* nav_menu = gtk_menu_new();
    
    // Add navigation items
    GtkWidget* back_item = gtk_menu_item_new_with_label("Back");
    g_signal_connect(back_item, "activate", G_CALLBACK(go_back_cb), web_view);
    gtk_menu_shell_append(GTK_MENU_SHELL(nav_menu), back_item);
    
    GtkWidget* forward_item = gtk_menu_item_new_with_label("Forward");
    g_signal_connect(forward_item, "activate", G_CALLBACK(go_forward_cb), web_view);
    gtk_menu_shell_append(GTK_MENU_SHELL(nav_menu), forward_item);
    
    GtkWidget* refresh_item = gtk_menu_item_new_with_label("Reload");
    g_signal_connect(refresh_item, "activate", G_CALLBACK(refresh_cb), web_view);
    gtk_menu_shell_append(GTK_MENU_SHELL(nav_menu), refresh_item);
    
    // Add a separator
    GtkWidget* separator = gtk_separator_menu_item_new();
    gtk_menu_shell_append(GTK_MENU_SHELL(nav_menu), separator);
    
    // Create Zoom submenu
    GtkWidget* zoom_menu_item = gtk_menu_item_new_with_label("Zoom");
    GtkWidget* zoom_menu = create_zoom_menu(web_view);
    gtk_menu_item_set_submenu(GTK_MENU_ITEM(zoom_menu_item), zoom_menu);
    
    // Add zoom submenu to navigation menu
    gtk_menu_shell_append(GTK_MENU_SHELL(nav_menu), zoom_menu_item);
    
    return nav_menu;
}

GtkWidget* create_file_menu(GtkWidget* window) {
    // Create "File" menu
    GtkWidget* file_menu = gtk_menu_new();
    
    // Add "Quit" to File menu
    GtkWidget* quit_item = gtk_menu_item_new_with_label("Quit");
    g_signal_connect(quit_item, "activate", G_CALLBACK(destroy_window), window);
    gtk_menu_shell_append(GTK_MENU_SHELL(file_menu), quit_item);
    
    return file_menu;
}

GtkWidget* create_menu_bar() {
    // Create the menu bar
    GtkWidget* menu_bar = gtk_menu_bar_new();
    return menu_bar;
}

GtkWidget* setup_window_layout(GtkWidget* window, WebKitWebView** web_view) {
    // Create a vertical box to hold the menu bar and web view
    GtkWidget* vbox = gtk_box_new(GTK_ORIENTATION_VERTICAL, 0);
    gtk_container_add(GTK_CONTAINER(window), vbox);
    
    // Create the menu bar
    GtkWidget* menu_bar = create_menu_bar();
    gtk_box_pack_start(GTK_BOX(vbox), menu_bar, FALSE, FALSE, 0);
    
    // Create "File" menu
    GtkWidget* file_menu_item = gtk_menu_item_new_with_label("File");
    GtkWidget* file_menu = create_file_menu(window);
    gtk_menu_item_set_submenu(GTK_MENU_ITEM(file_menu_item), file_menu);
    
    // Add File menu to menu bar
    gtk_menu_shell_append(GTK_MENU_SHELL(menu_bar), file_menu_item);
    
    // Create "Navigation" menu
    GtkWidget* nav_menu_item = gtk_menu_item_new_with_label("Navigation");
    
    // Create a WebKit web view
    *web_view = create_web_view();
    
    GtkWidget* nav_menu = create_navigation_menu(*web_view);
    gtk_menu_item_set_submenu(GTK_MENU_ITEM(nav_menu_item), nav_menu);
    
    // Add Navigation menu to menu bar
    gtk_menu_shell_append(GTK_MENU_SHELL(menu_bar), nav_menu_item);
    
    // Add web view to the vertical box
    gtk_box_pack_start(GTK_BOX(vbox), GTK_WIDGET(*web_view), TRUE, TRUE, 0);
    
    return vbox;
}

void load_html_content(WebKitWebView* web_view, const char* html_path) {
    // Load from local HTML file or provided path
    gchar* current_dir = g_get_current_dir();
    gchar* file_path;
    
    if (html_path != NULL) {
        file_path = g_strdup_printf("file:///%s", html_path);
    } else {
        file_path = g_strdup_printf("file:///%s/index.html", current_dir);
    }
    
    webkit_web_view_load_uri(web_view, file_path);
    g_free(file_path);
    g_free(current_dir);
}

int create_webview(int argc, char* argv[]) {
    // Initialize window
    GtkWidget* window = initialize_gtk_window(argc, argv, 800, 600, "My Web Browser", "./icon.png");
    
    // Setup WebView
    WebKitWebView* web_view;
    setup_window_layout(window, &web_view);
    
    // Load HTML content
    load_html_content(web_view, NULL);
    
    // Show the window and start the main loop
    gtk_widget_show_all(window);
    gtk_main();

    return 0;
}
