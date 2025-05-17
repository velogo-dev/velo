# GTK4 WebView with Go

This project demonstrates how to create a WebView application using Go with GTK4 and WebKitGTK 6.0. The implementation includes:

- A WebView component with GTK4 support
- Modern GTK4 menu implementation with a header bar
- Context menu support
- Clean separation between Go and C code

## Requirements

To build and run this project, you need:

- Go 1.16 or later
- GTK4 development libraries
- WebKitGTK 6.0 development libraries

On Debian/Ubuntu systems, you can install the dependencies with:

```bash
sudo apt-get install libgtk-4-dev libwebkitgtk-6.0-dev
```

On Fedora:

```bash
sudo dnf install gtk4-devel webkitgtk6.0-devel
```

On Arch Linux:

```bash
sudo pacman -S gtk4 webkitgtk-6.0
```

## Building the Application

Use the provided build script:

```bash
./build.sh
```

Or build manually:

```bash
# Compile the C part
gcc -c webview.c $(pkg-config --cflags --libs gtk4 webkitgtk-6.0) -o webview.o

# Build the Go application
go build -o gtk4webview *.go
```

## Running the Application

After building, run:

```bash
./gtk4webview
```

## Project Structure

- `webview.h` and `webview.c` - C implementation of the WebView component
- `webview.go` - Go bindings for the WebView component
- `app.go` - GTK4 application wrapper
- `callback.go` - Callback functions for GTK4 and Go integration
- `main.go` - Main application entry point

## Features

- **Modern GTK4 UI**: Uses a header bar with drop-down menu instead of traditional menu bar
- **Context Menu**: Supports custom context menu items
- **Web Navigation**: Full-featured web browser based on WebKitGTK
- **Clean API**: Simple Go API to interact with the WebView

## License

This project is released under the MIT License.
