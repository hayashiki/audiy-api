package router

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Router interface {
	CreateHandler() http.Handler
}

func NewRouter(r http.Handler, q http.Handler, l http.Handler) Router {
	return &router{
		r,
		q,
		l,
	}
}

type router struct {
	RootHandler     http.Handler
	QueryHandler    http.Handler
	LivenessHandler http.Handler
}

func (r router) CreateHandler() http.Handler {
	mux := mux.NewRouter()
	mux.Handle("/", r.RootHandler)
	mux.Handle("/query", r.QueryHandler)
	mux.Handle("/liveness", r.LivenessHandler)
	return mux
}
