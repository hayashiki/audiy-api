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
	playRepo := ds.NewPlayRepository(dsCli)
	commentRepo := ds.NewCommentRepository(dsCli)
	userRepo := ds.NewUserRepository(dsCli)

	// middleware
	authenticator := middleware.NewAuthenticator()

	// repository
	audioRepo := ds.NewAudioRepository(dsStore)

	// usecase
	audioUsecase := usecase.NewAudioUsecase(audioRepo)
	playUsecase := usecase.NewPlayUsecase(playRepo)
	commentUsecase := usecase.NewCommentUsecase(commentRepo)
	userUsecase := usecase.NewUserUsecase(userRepo)

	// handler
	queryHandler := handler.NewQueryHandler(userUsecase, audioUsecase, playUsecase, commentUsecase)
	authenticator.AuthMiddleware(queryHandler)
	rootHandler := handler.NewRootHandler()

	// router
	router := router.NewRouter(rootHandler, queryHandler, queryHandler)

	return router.CreateHandler()
}
