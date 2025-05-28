#include <gtk/gtk.h>
#include <webkitgtk-4.0/webkit2/webkit2.h>

static void destroy_window_cb(GtkWidget* widget, GtkWidget* window) {
    gtk_main_quit();
}

static void go_back_cb(GtkWidget* widget, WebKitWebView* web_view) {
    webkit_web_view_go_back(web_view);
}

static void go_forward_cb(GtkWidget* widget, WebKitWebView* web_view) {
    webkit_web_view_go_forward(web_view);
}

static void refresh_cb(GtkWidget* widget, WebKitWebView* web_view) {
    webkit_web_view_reload(web_view);
}

static void zoom_in_cb(GtkWidget* widget, WebKitWebView* web_view) {
    webkit_web_view_set_zoom_level(web_view, webkit_web_view_get_zoom_level(web_view) + 0.1);
}

static void zoom_out_cb(GtkWidget* widget, WebKitWebView* web_view) {
    webkit_web_view_set_zoom_level(web_view, webkit_web_view_get_zoom_level(web_view) - 0.1);
}

static void zoom_reset_cb(GtkWidget* widget, WebKitWebView* web_view) {
    webkit_web_view_set_zoom_level(web_view, 1.0);
}

int main(int argc, char* argv[]) {
    // Initialize GTK
    gtk_init(&argc, &argv);

    // Create the main window
    GtkWidget* window = gtk_window_new(GTK_WINDOW_TOPLEVEL);
    gtk_window_set_default_size(GTK_WINDOW(window), 800, 600);
    // Set Title
    gtk_window_set_title(GTK_WINDOW(window), "My Web Browser");
    g_signal_connect(window, "destroy", G_CALLBACK(destroy_window_cb), window);
    // Set Icon
    gtk_window_set_icon_name(GTK_WINDOW(window), "Pineapple");
    // Set Icon from file
    GError* error = NULL;
    gtk_window_set_icon_from_file(GTK_WINDOW(window), "./icon.png", &error);
    if (error) {
        g_error_free(error);
    }
   
    // Create a vertical box to hold the menu bar and web view
    GtkWidget* vbox = gtk_box_new(GTK_ORIENTATION_VERTICAL, 0);
    gtk_container_add(GTK_CONTAINER(window), vbox);
    
    // Create the menu bar
    GtkWidget* menu_bar = gtk_menu_bar_new();
    gtk_box_pack_start(GTK_BOX(vbox), menu_bar, FALSE, FALSE, 0);
    
    // Create "File" menu
    GtkWidget* file_menu_item = gtk_menu_item_new_with_label("File");
    GtkWidget* file_menu = gtk_menu_new();
    gtk_menu_item_set_submenu(GTK_MENU_ITEM(file_menu_item), file_menu);
    
    // Add "Quit" to File menu
    GtkWidget* quit_item = gtk_menu_item_new_with_label("Quit");
    g_signal_connect(quit_item, "activate", G_CALLBACK(destroy_window_cb), window);
    gtk_menu_shell_append(GTK_MENU_SHELL(file_menu), quit_item);
    
    // Add File menu to menu bar
    gtk_menu_shell_append(GTK_MENU_SHELL(menu_bar), file_menu_item);
    
    // Create "Navigation" menu
    GtkWidget* nav_menu_item = gtk_menu_item_new_with_label("Navigation");
    GtkWidget* nav_menu = gtk_menu_new();
    gtk_menu_item_set_submenu(GTK_MENU_ITEM(nav_menu_item), nav_menu);

    // Create a WebKit web view
    WebKitWebView* web_view = WEBKIT_WEB_VIEW(webkit_web_view_new());
    
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
    GtkWidget* zoom_menu = gtk_menu_new();
    gtk_menu_item_set_submenu(GTK_MENU_ITEM(zoom_menu_item), zoom_menu);
    
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
    
    // Add zoom submenu to navigation menu
    gtk_menu_shell_append(GTK_MENU_SHELL(nav_menu), zoom_menu_item);
    
    // Add Navigation menu to menu bar
    gtk_menu_shell_append(GTK_MENU_SHELL(menu_bar), nav_menu_item);
    
    // Add web view to the vertical box
    gtk_box_pack_start(GTK_BOX(vbox), GTK_WIDGET(web_view), TRUE, TRUE, 0);

    // Load from local HTML file
    gchar* current_dir = g_get_current_dir();
    gchar* file_path = g_strdup_printf("file:///%s/index.html", current_dir);
    webkit_web_view_load_uri(web_view, file_path);
    g_free(file_path);
    g_free(current_dir);

    // Show the window and start the main loop
    gtk_widget_show_all(window);
    gtk_main();

    return 0;
}
