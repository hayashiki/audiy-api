package app

import (
	"context"
	"log"
	"net/http"

	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"github.com/hayashiki/audiy-api/src/config"
	"github.com/hayashiki/audiy-api/src/domain/entity"
	"github.com/hayashiki/audiy-api/src/graph"
	"github.com/hayashiki/audiy-api/src/handler"
	"github.com/hayashiki/audiy-api/src/infrastructure/ds"
	"github.com/hayashiki/audiy-api/src/infrastructure/gcs"
	"github.com/hayashiki/audiy-api/src/infrastructure/slack"
	"github.com/hayashiki/audiy-api/src/logging"
	middleware2 "github.com/hayashiki/audiy-api/src/middleware"
	"github.com/hayashiki/audiy-api/src/usecase"
	"go.opencensus.io/plugin/ochttp"
	"go.uber.org/zap"
)

type Dependency struct {
	log            *zap.SugaredLogger
	audioUsecase   usecase.AudioUsecase
	userUsecase    usecase.UserUsecase
	feedUsecase    usecase.FeedUsecase
	commentUsecase usecase.CommentUsecase
	gcsSvc         gcs.Client
	slackSvc       slack.Slack
	commentRepo    entity.CommentRepository
	audioRepo      entity.AudioRepository
	userRepo       entity.UserRepository
	feedRepo       entity.FeedRepository
	resolver       *graph.Resolver
	authenticator  middleware2.Authenticator
	// TODO: handler struct
	apiHandler     http.Handler
	graphQLHandler http.Handler
}

func (d *Dependency) Inject() {
	// TODO: confを外にだす
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to read gcs client")
	}

	// infrastructure
	dsCli, _ := ds.NewClient(context.Background(), config.GetProject())
	// inject
	gcsClient, err := gcs.NewGCSClient(context.Background(), conf.GCSInputAudioBucket)
	if err != nil {
		log.Fatalf("failed to read gcs client")
	}
	slackSvc := slack.NewClient(conf.SlackBotToken)

	// repository
	commentRepo := ds.NewCommentRepository(dsCli)
	userRepo := ds.NewUserRepository(dsCli)
	audioRepo := ds.NewAudioRepository(dsCli)
	feedRepo := ds.NewFeedRepository(dsCli)

	// middleware
	authenticator := middleware2.NewAuthenticator()

	// usecase
	audioUsecase := usecase.NewAudioUsecase(gcsClient, audioRepo, feedRepo, userRepo)
	commentUsecase := usecase.NewCommentUsecase(commentRepo, audioRepo)
	userUsecase := usecase.NewUserUsecase(userRepo, audioRepo, feedRepo)
	feedUsecase := usecase.NewFeedUsecase(feedRepo, audioRepo)

	logger := logging.NewLogger(conf.IsDev)

	d.log = logger
	d.audioUsecase = audioUsecase
	d.userUsecase = userUsecase
	d.feedUsecase = feedUsecase
	d.commentUsecase = commentUsecase
	d.gcsSvc = gcsClient
	d.slackSvc = slackSvc
	d.commentRepo = commentRepo
	d.audioRepo = audioRepo
	d.userRepo = userRepo
	d.feedRepo = feedRepo
	d.authenticator = authenticator

	resolver := graph.NewResolver(userUsecase, audioUsecase, commentUsecase, feedUsecase)
	d.resolver = resolver
	graphHandler := NewGraphQLHandler(d.resolver)
	d.graphQLHandler = graphHandler
	graphHandler = &ochttp.Handler{
		Handler:     graphHandler,
		Propagation: &propagation.HTTPFormat{},
	}

	rootHandler := NewRootHandler()
	rootHandler = &ochttp.Handler{
		Handler:     rootHandler,
		Propagation: &propagation.HTTPFormat{},
	}

	apiHandler := handler.NewAPIHandler(slackSvc, gcsClient, audioRepo, feedRepo, userRepo)
	d.apiHandler = apiHandler
}
