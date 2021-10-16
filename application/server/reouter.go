package server

import (
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/hayashiki/audiy-api/application/handler"
)

func Routing(r *chi.Mux, d *Dependency) {
	r.Use(middleware.Recoverer)
	r.Get("/health", handler.HealthCheck)
	r.Group(func(r chi.Router) {
		r.Handle("/enqueue/create_audio", d.apiHandler)
	})

	r.Group(func(r chi.Router) {
		//r.Route("query", func(r chi.Router) {
		r.Use(middleware.RequestID)
		r.Use(middlewareCors)
		//r.Use(middlewareLogger(logger))
		r.With(d.authenticator.AuthMiddleware).Handle("/", playground.Handler("GraphQL playground", "/query"))
		r.With(d.authenticator.AuthMiddleware).Handle("/query", d.graphQLHandler)
		//})
	})
}
