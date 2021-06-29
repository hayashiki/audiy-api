package router

import (
	"github.com/hayashiki/audiy-api/interfaces/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

type Router interface {
	CreateHandler(authenticator middleware.Authenticator) http.Handler
}

func NewRouter(r http.Handler, q http.Handler, l http.Handler) Router {
	return &router{
		r,
		q,
		l,
	}
}

type router struct {
	RootHandler   http.Handler
	QueryHandler  http.Handler
	HealthHandler http.Handler
}

func (r router) CreateHandler(authenticator middleware.Authenticator) http.Handler {
	mux := mux.NewRouter()
	mux.Handle("/", r.RootHandler)
	mux.Handle("/query", middleware.Cors(authenticator.AuthMiddleware(r.QueryHandler)))
	mux.Handle("/health", r.HealthHandler)
	return mux
}