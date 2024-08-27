package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func initEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Assuming environment variables are set.")
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request received: %s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "CDN Service is up and running!")
}

func handleCDNContent(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Serving CDN content...")
}

func main() {
	initEnv()

	router := mux.NewRouter()

	router.Use(loggingMiddleware)

	router.HandleFunc("/", handleMain).Methods("GET")
	router.HandleFunc("/cdn/{content}", handleCDNContent).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}