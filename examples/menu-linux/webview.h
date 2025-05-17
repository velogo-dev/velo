#ifndef WEBVIEW_H
#define WEBVIEW_H

#include <gtk/gtk.h>
#include <webkitgtk-6.0/webkit/webkit.h>

typedef struct {
    GtkWidget *window;
    GtkWidget *webview;
    GtkWidget *menubar;
    GMenu *popup_menu;
} WebViewWindow;

WebViewWindow *create_webview_window(const char *title, int width, int height, const char *url);
void webview_set_title(WebViewWindow *wv, const char *title);
void webview_navigate(WebViewWindow *wv, const char *url);
void webview_show(WebViewWindow *wv);
void webview_destroy(WebViewWindow *wv);
void webview_run();
void webview_add_context_menu_item(WebViewWindow *wv, const char *label);

#endif /* WEBVIEW_H */ 