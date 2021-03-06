package app

import (
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/hayashiki/audiy-api/src/handler"
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
		r.Use(middlewareLogger(d.log))
		r.With(d.authenticator.AuthMiddleware).With(d.dataloaderSvc.Middleware).Handle("/", playground.Handler("GraphQL playground", "/query"))
		r.With(d.authenticator.AuthMiddleware).With(d.dataloaderSvc.Middleware).Handle("/query", d.graphQLHandler)
		//})
	})
}
