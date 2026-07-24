package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const version = "0.0.1"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Hello Social Media Hack 2026/07!")
	})

	log.Printf("AgencyPulse v%s starting on http://localhost:%s...", version, port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

