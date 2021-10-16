package server

import (
	"net/http"

	gqlHandler "github.com/99designs/gqlgen/graphql/handler"
)

type Handler struct {
	gqlHandler *gqlHandler.Server
	apiHandler http.Server
	//	logger
}
