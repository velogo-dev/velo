#!/bin/bash

# Ensure necessary dependencies are installed
if ! pkg-config --exists gtk4 webkitgtk-6.0; then
    echo "Error: Required packages not found. Please install GTK4 and WebKitGTK 6.0."
    echo "On Ubuntu/Debian: sudo apt-get install libgtk-4-dev libwebkitgtk-6.0-dev"
    echo "On Fedora: sudo dnf install gtk4-devel webkitgtk6.0-devel"
    echo "On Arch: sudo pacman -S gtk4 webkitgtk-6.0"
    exit 1
fi

# Build the simplified Go application
go build -o gtk4webview_simple gtk4_webview.go

# Check if build was successful
if [ $? -eq 0 ]; then
    echo "Build successful! Run ./gtk4webview_simple to start the application."
else
    echo "Build failed."
    exit 1
fi 