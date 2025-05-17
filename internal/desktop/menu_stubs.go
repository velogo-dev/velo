//go:build !windows
// +build !windows

package desktop

import "unsafe"

// WindowsCreateAndInstallMenu is a stub for non-Windows platforms
func WindowsCreateAndInstallMenu(hwnd unsafe.Pointer, menu *NativeMenu) error {
	// Should never be called on non-Windows platforms
	return nil
}
