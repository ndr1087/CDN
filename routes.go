package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

// Cache structure to simulate caching mechanism for CDN
type Cache map[string]string

var cache Cache

func init() {
	cache = make(Cache)
}

// Handler for content delivery
func getContentHandler(w http.ResponseWriter, r *http.Request) {
	contentID := r.URL.Query().Get("id")
	if content, ok := cache[contentID]; ok {
		w.Write([]byte(content))
	} else {
		http.Error(w, "Content Not Found", http.StatusNotFound)
	}
}

// Handler for adding content to the cache
func addContentHandler(w http.ResponseWriter, r *http.Request) {
	contentID := r.URL.Query().Get("id")
	contentValue := r.URL.Query().Get("value")
	cache[contentID] = contentValue
	w.WriteHeader(http.StatusCreated)
}

// Handler for removing content from the cache
func removeContentHandler(w http.ResponseWriter, r *http.Request) {
	contentID := r.URL.Query().Get("id")
	delete(cache, contentID)
	w.WriteHeader(http.StatusOK)
}

// Handler for checking the CDN status
func statusCheckHandler(w http.ResponseWriter, r *http.Request) {
	status := struct {
		Status string `json:"status"`
	}{
		Status: "CDN is up and running",
	}
	resp, err := json.Marshal(status)
	if err != nil {
		http.Error(w, "Error generating status", http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}

func main() {
	http.HandleFunc("/getContent", getContentHandler)
	http.HandleFunc("/addContent", addContentHandler)
	http.HandleFunc("/removeContent", removeContentHandler)
	http.HandleFunc("/status", statusCheckHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	log.Println("CDN Server starting on port:", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}