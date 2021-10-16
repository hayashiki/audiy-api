package server

import (
	"net/http"

	"github.com/hayashiki/audiy-api/graph/trace"

	"github.com/hayashiki/audiy-api/graph"
	"github.com/hayashiki/audiy-api/graph/generated"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
)

// NewGraphQLHandler returns GraphQL Server.
func NewGraphQLHandler(
	resolver *graph.Resolver,
) http.Handler {
	srv := handler.New(generated.NewExecutableSchema(
		generated.Config{
			Resolvers:  resolver,
			Directives: generated.DirectiveRoot{},
			Complexity: generated.ComplexityRoot{},
		},
	))
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})
	srv.SetQueryCache(lru.New(1000))
	srv.Use(extension.Introspection{})
	srv.Use(trace.Tracer{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})

	return srv
}

// NewRootHandler returns GraphQL Playground.
func NewRootHandler() http.Handler {
	return playground.Handler("GraphQL playground", "/query")
}
