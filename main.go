package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Printf("Error loading .env file: %v", err)
	} else {
		log.Println("Environment variables have been successfully loaded or are already set.")
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request received: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "CDN Service is up and running!")
}

func handleCDNContent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	content, ok := vars["content"]
	if !ok {
		http.Error(w, "Content identifier missing", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Serving CDN content for: %s", content)
}

func main() {
	router := mux.NewRouter()
	router.Use(loggingMiddleware)

	router.HandleFunc("/", handleMain).Methods("GET")
	router.HandleFunc("/cdn/{content}", handleCDNContent).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		log.Println("PORT environment variable not set, defaulting to 8080")
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}