package main

import (
	"context"
	"fmt"
	"github.com/extism/go-sdk"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func init() {
	// Load environment variables from .env file, if present
	_ = os.Getenv(".env")
}

func main() {
	router := mux.NewRouter()

	// Define route for fetching URL content
	router.HandleFunc("/fetch-url-content", FetchURLContentHandler).Methods("GET")

	// Set up HTTP server
	http.Handle("/", router)
	port := "8080" // Default port
	if customPort := os.Getenv("PORT"); customPort != "" {
		port = customPort
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// FetchURLContentHandler handles requests to fetch content from a URL
func FetchURLContentHandler(w http.ResponseWriter, r *http.Request) {
	// Read URL from query parameter
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "URL parameter is required", http.StatusBadRequest)
		return
	}

	// Initialize the Extism plugin
	ctx := context.Background()
	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmUrl{
				Url: "https://github.com/OpenAgentsInc/plugin-url-scraper-go/raw/main/host-functions.wasm",
			},
		},
	}
	plugin, err := extism.NewPlugin(ctx, manifest, extism.PluginConfig{}, []extism.HostFunction{})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to initialize plugin: %v", err), http.StatusInternalServerError)
		return
	}

	// Call the "fetch_url_content" function on the plugin
	exit, out, err := plugin.Call("fetch_url_content", []byte(url))
	if err != nil || exit != 0 {
		http.Error(w, fmt.Sprintf("Plugin call failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Return the fetched content
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}


