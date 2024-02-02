package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
}

func main() {
	router := mux.NewRouter()

	// Define route for fetching URL content
	router.HandleFunc("/fetch-url-content", FetchURLContentHandler).Methods("GET")

	// Set up HTTP server
	http.Handle("/", router)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// FetchURLContentHandler handles requests to fetch content from a URL after payment verification
func FetchURLContentHandler(w http.ResponseWriter, r *http.Request) {
	// Placeholder for payment verification logic
	// Verify payment here...

	// Read URL from query parameter
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "URL parameter is required", http.StatusBadRequest)
		return
	}

	// Log the URL to the console
	log.Printf("Fetching content from URL: %s", url)

	// Assuming payment is verified and URL parameter is present
	// Placeholder for executing the Extism plugin with the URL
	output, err := executeExtismPluginWithURL(url) // You need to implement this function
	if err != nil {
		http.Error(w, "Failed to fetch URL content: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the fetched content
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(output))
}

// Placeholder for the function to execute the Extism plugin with a given URL
// You need to replace this with actual implementation
func executeExtismPluginWithURL(url string) (string, error) {
	// Implement the logic to execute the Extism plugin and return the content or an error
	return "", nil // Placeholder return
}

