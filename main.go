package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "strconv"

    "github.com/extism/go-sdk"
    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    "github.com/kodylow/matador/pkg/auth"
    "github.com/kodylow/matador/pkg/models"
    "github.com/kodylow/matador/pkg/database"
    "github.com/rs/cors"
)

var globalPrice uint64

func init() {
    // Load .env variables
    if err := godotenv.Load(".env"); err != nil {
        log.Fatal("Error loading .env file: ", err)
    }

    // Initialize global settings
    var err error
    globalPrice, err = strconv.ParseUint(os.Getenv("PRICE_SATS"), 10, 64)
    if err != nil {
        log.Fatal("Error converting global price: ", err)
    }

    if err := database.InitDatabase(); err != nil {
        log.Fatal("Error initializing database: ", err)
    }

    if err := auth.InitSecret(); err != nil {
        log.Fatal("Error initializing secret for server side tokens/runes: ", err)
    }
}

func main() {
    router := mux.NewRouter()

    // Define route for fetching URL content with L402 support
    router.HandleFunc("/fetch-url-content", FetchURLContentHandler).Methods("GET")

    // CORS setup
    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"*"},
        AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
        AllowedHeaders:   []string{"*"},
        ExposedHeaders:   []string{"*"},
        AllowCredentials: true,
    })

    // Start the server
    log.Println("Server starting on port 8080")
    log.Fatal(http.ListenAndServe(":8080", c.Handler(router)))
}

func FetchURLContentHandler(w http.ResponseWriter, r *http.Request) {
    url := r.URL.Query().Get("url")
    if url == "" {
        http.Error(w, "URL parameter is required", http.StatusBadRequest)
        return
    }

    // L402 Payment Check
    reqInfo := request.RequestInfo{Path: r.URL.Path}
    if err := auth.CheckAuthorizationHeader(reqInfo); err != nil {
        log.Println("Unauthorized, payment required")
        l402, err := auth.GetL402(globalPrice, reqInfo)
        if err != nil {
            log.Println("Error getting L402:", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }
        w.Header().Set("WWW-Authenticate", l402)
        http.Error(w, "Payment Required", http.StatusPaymentRequired)
        return
    }

    // Proceed with Extism plugin call
    ctx := context.Background()
    manifest := extism.Manifest{
        Wasm: []extism.Wasm{
            extism.WasmUrl{
                Url: "https://github.com/OpenAgentsInc/plugin-url-scraper-go/raw/main/host-functions.wasm",
            },
        },
        AllowedHosts: []string{"*"},
    }
    plugin, err := extism.NewPlugin(ctx, manifest, extism.PluginConfig{EnableWasi: true}, nil)
    if err != nil {
        log.Printf("Failed to initialize plugin: %v\n", err)
        http.Error(w, "Failed to initialize plugin", http.StatusInternalServerError)
        return
    }

    // Call the plugin function
    exit, out, err := plugin.Call("fetch_url_content", []byte(url))
    if err != nil || exit != 0 {
        log.Printf("Plugin call failed: %v\n", err)
        http.Error(w, "Plugin call failed", http.StatusInternalServerError)
        return
    }

    // Return the fetched content
    w.WriteHeader(http.StatusOK)
    w.Write(out)
}

