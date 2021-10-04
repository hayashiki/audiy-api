package router

import (
	"net/http"

	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"

	"github.com/hayashiki/audiy-api/interfaces/middleware"

	"github.com/gorilla/mux"
)

type Router interface {
	CreateHandler() http.Handler
}

func NewRouter(r http.Handler, q http.Handler, l http.Handler, a http.Handler) Router {
	return &router{
		r,
		q,
		l,
		a,
	}
}

type router struct {
	RootHandler   http.Handler
	QueryHandler  http.Handler
	HealthHandler http.Handler
	APIHandler    http.Handler
}

func (r router) CreateHandler() http.Handler {
	mux := mux.NewRouter()
	httptrace.WrapHandler(r.RootHandler, "audiy-api", "/query")
	mux.Handle("/", r.RootHandler)
	mux.Handle("/query", middleware.Cors(r.QueryHandler))
	mux.Handle("/health", r.HealthHandler)
	mux.Handle("/enqueue/create_audio", r.APIHandler)
	return mux
}
