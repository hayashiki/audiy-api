package main

import (
	"log"
	"net/http"
	"os"

	config2 "github.com/hayashiki/audiy-api/src/config"

	"github.com/go-chi/chi"

	"github.com/hayashiki/audiy-api/src/app"

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
		ProjectID: config2.GetProject(),
	})
	if err != nil {
		log.Fatal(err)
	}
	trace.RegisterExporter(exporter)
	trace.AlwaysSample()

	d := &app.Dependency{}
	d.Inject()
	r := chi.NewRouter()
	app.Routing(r, d)
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
