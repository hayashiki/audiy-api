package handler

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/hayashiki/audiy-api/interfaces/api/graph"
	"github.com/hayashiki/audiy-api/interfaces/api/graph/generated"
)

// NewQueryHandler returns GraphQL Server.
func NewQueryHandler(
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
	srv.SetQueryCache(lru.New(1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})

	return srv
}

// NewRootHandler returns GraphQL Playground.
func NewRootHandler() http.Handler {
	return playground.Handler("GraphQL playground", "/query")
}
