package registry

import (
	"github.com/hayashiki/audiy-api/interfaces/middleware"
	"net/http"

	"github.com/hayashiki/audiy-api/application/usecase"
	"github.com/hayashiki/audiy-api/domain/entity"
	"github.com/hayashiki/audiy-api/infrastructure/ds"
	"github.com/hayashiki/audiy-api/interfaces/api/graph/handler"
	"github.com/hayashiki/audiy-api/interfaces/api/graph/router"
)

type Registry interface {
	NewHandler() http.Handler
}

type registry struct{}

func NewRegistry() Registry {
	return &registry{}
}

func (s *registry) NewHandler() http.Handler {
	//dsCli := ds.Client()
	//dsStore := &ds.DataStore{Client: dsCli}
	// infrastructure
	dsStore := ds.Connect()

	// middleware
	authenticator := middleware.NewAuthenticator()

	// repository
	repo := entity.NewAudioRepository(dsStore)

	// usecase
	audioUsecase := usecase.NewAudioUsecase(repo)

	// handler
	queryHandler := handler.NewQueryHandler(audioUsecase)
	rootHandler := handler.NewRootHandler()

	// router
	router := router.NewRouter(rootHandler, queryHandler, queryHandler)

	return router.CreateHandler(authenticator)
}
