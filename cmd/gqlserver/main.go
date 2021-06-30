package main

import (
	"github.com/hayashiki/audiy-api/interfaces/registry"
	"log"
	"net/http"
	"os"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	registry := registry.NewRegistry()

	h := registry.NewHandler()
	http.Handle("/", h)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
