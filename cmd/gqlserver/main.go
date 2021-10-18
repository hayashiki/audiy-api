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

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}
	go func() {
		err := server.ListenAndServe()
		log.Println(err)
	}()

	quitChan := make(chan os.Signal, 1)
	signal.Notify(quitChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	sig := <-quitChan
	log.Printf("received signal %q; shutdown gracefully in %s ...", sig, shutdownTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	errChan := make(chan error)
	go func() { errChan <- server.Shutdown(ctx) }()

	select {
	case sig := <-quitChan:
		log.Printf("received 2nd signal %q; shutdown now", sig)
		cancel()
		server.Close()

	case err := <-errChan:
		if err != nil {
			log.Fatalf("while shutdown: %s", err)
		}
	}
}
