package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Cache map[string]string

var cache Cache

func init() {
	cache = make(Cache)
}

func getContentHandler(w http.ResponseWriter, r *http.Request) {
	contentID := r.URL.Query().Get("id")
	if content, ok := cache[contentID]; ok {
		if _, err := w.Write([]byte(content)); err != nil {
			log.Printf("Error sending response for getContent: %v\n", err)
			http.Error(w, "Error writing content response", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Content Not Found", http.StatusNotFound)
	}
}

func addContentHandler(w http.ResponseWriter, r *http.Request) {
	contentID := r.URL.Query().Get("id")
	contentValue := r.URL.Query().Get("value")
	cache[contentID] = contentValue
	w.WriteHeader(http.StatusCreated)
}

func removeContentHandler(w http.ResponseWriter, r *http.Request) {
	contentID := r.URL.Query().Get("id")
	delete(cache, contentID)
	w.WriteHeader(http.StatusOK)
}

func statusCheckHandler(w http.ResponseWriter, r *http.Request) {
	status := struct {
		Status string `json:"status"`
	}{
		Status: "CDN is up and running",
	}
	resp, err := json.Marshal(status)
	if err != nil {
		log.Printf("Error marshalling status response: %v\n", err)
		http.Error(w, "Error generating status", http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(resp); err != nil {
		log.Printf("Error sending status response: %v\n", err)
		http.Error(w, "Error writing status response", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/getContent", getContentHandler)
	http.HandleFunc("/addContent", addContentHandler)
	http.HandleFunc("/removeContent", removeContentHandler)
	http.HandleFunc("/status", statusCheckHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" 
	}

	log.Println("CDN Server starting on port:", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}