#ifndef WEBKIT_H
#define WEBKIT_H

#include <gtk/gtk.h>
#include <webkitgtk-4.0/webkit2/webkit2.h>

// Event handlers
void destroy_window(GtkWidget* widget, GtkWidget* window);
void go_back_cb(GtkWidget* widget, WebKitWebView* web_view);
void go_forward_cb(GtkWidget* widget, WebKitWebView* web_view);
void refresh_cb(GtkWidget* widget, WebKitWebView* web_view);
void zoom_in_cb(GtkWidget* widget, WebKitWebView* web_view);
void zoom_out_cb(GtkWidget* widget, WebKitWebView* web_view);
void zoom_reset_cb(GtkWidget* widget, WebKitWebView* web_view);

// Window management
void init(int* argc, char** argv[]);
GtkWidget* new_window(GtkWidget* widget, GtkWidget* parent_window);
GtkWidget* set_window_size(GtkWidget* window, int width, int height);
GtkWidget* set_window_title(GtkWidget* window, char* title);
GtkWidget* set_window_icon(GtkWidget* window, char* icon);
void run_main_loop(GtkWidget* window);

// Webview components
GtkWidget* create_menu_bar();
GtkWidget* create_file_menu(GtkWidget* window);
GtkWidget* create_navigation_menu(WebKitWebView* web_view);
GtkWidget* create_zoom_menu(WebKitWebView* web_view);
WebKitWebView* create_web_view();

// Main functions
GtkWidget* initialize_gtk_window(int argc, char* argv[], int width, int height, const char* title, const char* icon_path);
GtkWidget* setup_window_layout(GtkWidget* window, WebKitWebView** web_view);
void load_html_content(WebKitWebView* web_view, const char* html_path);
int create_webview(int argc, char* argv[]);

#endif // WEBKIT_H 