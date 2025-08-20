package linux

// #cgo pkg-config: gtk+-3.0 webkit2gtk-4.0
// #include "webkit.h"
import "C"
import (
	"os"
	"unsafe"
)

// CreateWebView initializes and runs the GTK WebKit view
func CreateWebView() int {
	argc := C.int(len(os.Args))
	argv := make([]*C.char, len(os.Args))
	for i, arg := range os.Args {
		argv[i] = C.CString(arg)
	}
	defer func() {
		for _, arg := range argv {
			C.free(unsafe.Pointer(arg))
		}
	}()

	return int(C.create_webview(argc, &argv[0]))
}
