package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"

	"github.com/hayashiki/audiy-api/application/server"

	"github.com/hayashiki/audiy-api/etc/config"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/trace"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	//Profiler initialization, best done as early as possible.
	//if err := profiler.Start(profiler.Config{
	//	ProjectID: os.Getenv("GCP_PROJECT"),
	//}); err != nil {
	//	log.Fatal(err)
	//}

	// Create and register a OpenCensus Stackdriver Trace exporter.
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: config.GetProject(),
	})
	if err != nil {
		log.Fatal(err)
	}
	trace.RegisterExporter(exporter)
	trace.AlwaysSample()

	d := &server.Dependency{}
	d.Inject()
	r := chi.NewRouter()
	server.Routing(r, d)
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
