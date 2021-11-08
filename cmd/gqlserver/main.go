package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hayashiki/audiy-api/src/config"
	config2 "github.com/hayashiki/audiy-api/src/config"

	"github.com/go-chi/chi"

	"github.com/hayashiki/audiy-api/src/app"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/trace"
)

const defaultPort = "8080"
const shutdownTimeout = 30 * time.Second

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

	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to read config")
	}

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
	d.Inject(conf)
	r := chi.NewRouter()
	app.Routing(r, d)
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}
	go func() {
		log.Printf("Listening on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, os.Interrupt)
	<-sigCh

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("graceful shutdown failure: %s", err)
	}
	log.Printf("graceful shutdown successfully")
}
