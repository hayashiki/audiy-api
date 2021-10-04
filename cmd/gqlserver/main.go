package main

import (
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/profiler"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"github.com/hayashiki/audiy-api/interfaces/registry"
	"go.opencensus.io/trace"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	//Profiler initialization, best done as early as possible.
	if err := profiler.Start(profiler.Config{
		ProjectID: os.Getenv("GCP_PROJECT"),
	}); err != nil {
		log.Fatal(err)
	}

	// Create and register a OpenCensus Stackdriver Trace exporter.
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: os.Getenv("GCP_PROJECT"),
	})
	if err != nil {
		log.Fatal(err)
	}
	trace.RegisterExporter(exporter)
	trace.AlwaysSample()

	registry := registry.NewRegistry()

	h := registry.NewHandler()
	http.Handle("/", h)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
