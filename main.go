package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type ConfigHandler struct {
	configDir string
}

func NewConfigHandler(configDir string) *ConfigHandler {
	return &ConfigHandler{
		configDir: configDir,
	}
}

func (h *ConfigHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get config file path from request path
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		http.Error(w, "Path is required", http.StatusBadRequest)
		return
	}

	// Build full path to config directory
	configPath := filepath.Join(h.configDir, path)

	// Check if directory exists
	dirInfo, err := os.Stat(configPath)
	if os.IsNotExist(err) || !dirInfo.IsDir() {
		http.Error(w, "Config directory not found", http.StatusNotFound)
		return
	}

	// Look for JSON files in the directory
	files, err := os.ReadDir(configPath)
	if err != nil {
		http.Error(w, "Error reading config directory", http.StatusInternalServerError)
		return
	}

	// Use the first JSON file found
	var jsonFile string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			jsonFile = file.Name()
			break
		}
	}

	if jsonFile == "" {
		http.Error(w, "No JSON file found in config directory", http.StatusNotFound)
		return
	}

	// Read JSON file
	data, err := os.ReadFile(filepath.Join(configPath, jsonFile))
	if err != nil {
		http.Error(w, "Error reading config file", http.StatusInternalServerError)
		return
	}

	// Validate JSON format
	var jsonData any
	if err := json.Unmarshal(data, &jsonData); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func main() {
	// Set command line arguments
	configDir := flag.String("config", ".dirserve", "Directory to serve. Default is .dirserve.")
	port := flag.Int("port", 8080, "(Optional) Change listen port. Default is 8080.")
	flag.Parse()

	// Check if config directory exists
	if _, err := os.Stat(*configDir); os.IsNotExist(err) {
		log.Fatalf("Config directory '%s' not found", *configDir)
	}

	handler := NewConfigHandler(*configDir)

	// Start server
	fmt.Printf("Config directory: %s\n", *configDir)
	fmt.Printf("Starting server on port: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), handler))
}
