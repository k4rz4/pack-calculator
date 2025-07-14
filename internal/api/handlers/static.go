package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// StaticHandler serves static files
type StaticHandler struct {
	webDir string
}

// NewStaticHandler creates a new static handler
func NewStaticHandler() *StaticHandler {
	webDir := "./web"

	// Debug: Check if web directory exists
	if _, err := os.Stat(webDir); os.IsNotExist(err) {
		log.Printf("WARNING: Web directory %s does not exist", webDir)
	} else {
		log.Printf("Web directory found: %s", webDir)
	}

	return &StaticHandler{
		webDir: webDir,
	}
}

// ServeUI serves the main UI
func (s *StaticHandler) ServeUI(w http.ResponseWriter, r *http.Request) {
	indexPath := filepath.Join(s.webDir, "index.html")

	// Debug: Check if index.html exists
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		log.Printf("ERROR: index.html not found at %s", indexPath)
		http.Error(w, "UI not found", http.StatusNotFound)
		return
	}

	log.Printf("Serving UI: %s", indexPath)
	http.ServeFile(w, r, indexPath)
}

// ServeStatic serves static assets (CSS, JS, etc.)
func (s *StaticHandler) ServeStatic(w http.ResponseWriter, r *http.Request) {
	// Remove /static prefix from path
	path := strings.TrimPrefix(r.URL.Path, "/static")
	if path == "" || path == "/" {
		http.NotFound(w, r)
		return
	}

	// Serve the file
	filePath := filepath.Join(s.webDir, path)

	// Debug: Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("ERROR: Static file not found: %s", filePath)
		http.NotFound(w, r)
		return
	}

	log.Printf("Serving static file: %s", filePath)
	http.ServeFile(w, r, filePath)
}
