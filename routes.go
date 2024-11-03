package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type ContentCache map[string]string

var contentCache ContentCache

func init() {
	contentCache = make(ContentCache)
}

func handleGetContent(w http.ResponseWriter, r *http.Request) {
	contentID := r.URL.Query().Get("id")
	if content, exists := contentCache[contentID]; exists {
		if _, err := w.Write([]byte(content)); err != nil {
			log.Printf("Error sending content for ID %s: %v\n", contentID, err)
			http.Error(w, "Error serving content", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Content not found", http.StatusNotFound)
	}
}

func handleAddContent(w http.ResponseWriter, r *http.Request) {
	contentID := r.URL.Query().Get("id")
	contentValue := r.URL.Query().Get("value")
	contentCache[contentID] = contentValue
	w.WriteHeader(http.StatusCreated)
}

func handleRemoveContent(w http.ResponseWriter, r *http.Request) {
	contentID := r.URL.Query().Get("id")
	delete(contentCache, contentID)
	w.WriteHeader(http.StatusOK)
}

func handleStatusCheck(w http.ResponseWriter, r *http.Request) {
	status := struct {
		Status string `json:"status"`
	}{
		Status: "CDN is operational",
	}
	response, err := json.Marshal(status)
	if err != nil {
		log.Printf("Error marshalling status check response: %v\n", err)
		http.Error(w, "Error generating status", http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(response); err != nil {
		log.Printf("Error sending status check response: %v\n", err)
		http.Error(w, "Error writing status response", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/getContent", handleGetContent)
	http.HandleFunc("/addContent", handleAddContent)
	http.HandleFunc("/removeContent", handleRemoveContent)
	http.HandleFunc("/status", handleStatusCheck)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("CDN Server starting on port:", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}