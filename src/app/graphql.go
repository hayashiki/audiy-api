package app

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/gorilla/websocket"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"net/http"
	"time"

	"github.com/hayashiki/audiy-api/src/graph"
	"github.com/hayashiki/audiy-api/src/graph/generated"
	"github.com/hayashiki/audiy-api/src/graph/trace"

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
	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		KeepAlivePingInterval: 10 * time.Second,
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})
	srv.SetQueryCache(lru.New(1000))
	srv.Use(extension.Introspection{})
	srv.Use(trace.Tracer{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})
	srv.SetErrorPresenter(errorPresenter)
	srv.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		return &gqlerror.Error{
			Message: "server error",
			Extensions: map[string]interface{}{
				"type": "Internal",
				"code": "Unknown",
			},
		}
	})

	return srv
}

func errorPresenter(ctx context.Context, e error) *gqlerror.Error {
	err := graphql.DefaultErrorPresenter(ctx, e)
	return err
}

// NewRootHandler returns GraphQL Playground.
func NewRootHandler() http.Handler {
	return playground.Handler("GraphQL playground", "/query")
}
