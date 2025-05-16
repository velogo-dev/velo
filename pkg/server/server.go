package server

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
)

// PreviewServer represents a local HTTP server for previewing assets
type PreviewServer struct {
	RootDir   string
	AssetsDir string
	Port      string
}

// NewPreviewServer creates a new preview server
func NewPreviewServer(rootDir, port string) *PreviewServer {
	return &PreviewServer{
		RootDir:   rootDir,
		AssetsDir: filepath.Join(rootDir, "mobile-shell", "assets"),
		Port:      port,
	}
}

// Start starts the preview server
func (s *PreviewServer) Start() error {
	fs := http.FileServer(http.Dir(s.AssetsDir))
	http.Handle("/", fs)

	addr := fmt.Sprintf(":%s", s.Port)
	fmt.Printf("Preview server running at http://localhost:%s\n", s.Port)

	return http.ListenAndServe(addr, nil)
}

// StartBackground starts the preview server in a background goroutine
func (s *PreviewServer) StartBackground() {
	go func() {
		if err := s.Start(); err != nil {
			log.Printf("Preview server error: %v", err)
		}
	}()
}
