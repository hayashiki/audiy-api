package registry

import (
	"context"
	"github.com/hayashiki/audiy-api/interfaces/middleware"
	"net/http"
	"os"

	"github.com/hayashiki/audiy-api/application/usecase"
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

	dsCli, _ := ds.NewClient(context.Background(), os.Getenv("GCP_PROJECT"))
	audioUserRepo := ds.NewAudioUserRepository(dsCli)

	// middleware
	authenticator := middleware.NewAuthenticator()

	// repository
	repo := ds.NewAudioRepository(dsStore)

	// usecase
	audioUsecase := usecase.NewAudioUsecase(repo)
	audioUserUsecase := usecase.NewAudioUserUsecase(audioUserRepo)

	// handler
	queryHandler := handler.NewQueryHandler(audioUsecase, audioUserUsecase)
	rootHandler := handler.NewRootHandler()

	// router
	router := router.NewRouter(rootHandler, queryHandler, queryHandler)

	return router.CreateHandler(authenticator)
}
