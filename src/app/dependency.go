package app

import (
	"context"
	"log"
	"net/http"

	config2 "github.com/hayashiki/audiy-api/src/config"

	entity2 "github.com/hayashiki/audiy-api/src/domain/entity"
	graph2 "github.com/hayashiki/audiy-api/src/graph"
	ds2 "github.com/hayashiki/audiy-api/src/infrastructure/ds"
	gcs2 "github.com/hayashiki/audiy-api/src/infrastructure/gcs"
	slack2 "github.com/hayashiki/audiy-api/src/infrastructure/slack"
	usecase2 "github.com/hayashiki/audiy-api/src/usecase"

	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"github.com/hayashiki/audiy-api/src/handler"
	middleware2 "github.com/hayashiki/audiy-api/src/middleware"
	"go.opencensus.io/plugin/ochttp"
)

type Dependency struct {
	audioUsecase   usecase2.AudioUsecase
	userUsecase    usecase2.UserUsecase
	feedUsecase    usecase2.FeedUsecase
	commentUsecase usecase2.CommentUsecase
	gcsSvc         gcs2.Client
	slackSvc       slack2.Slack
	commentRepo    entity2.CommentRepository
	audioRepo      entity2.AudioRepository
	userRepo       entity2.UserRepository
	feedRepo       entity2.FeedRepository
	resolver       *graph2.Resolver
	authenticator  middleware2.Authenticator
	apiHandler     http.Handler
	graphQLHandler http.Handler
}

func (d *Dependency) Inject() {
	// TODO: confを外にだす
	conf, err := config2.NewConf()
	if err != nil {
		log.Fatalf("failed to read gcs client")
	}

	// infrastructure
	dsCli, _ := ds2.NewClient(context.Background(), config2.GetProject())
	// inject
	gcsClient, err := gcs2.NewGCSClient(context.Background(), conf.GCSInputAudioBucket)
	if err != nil {
		log.Fatalf("failed to read gcs client")
	}
	slackSvc := slack2.NewClient(conf.SlackBotToken)

	// repository
	commentRepo := ds2.NewCommentRepository(dsCli)
	userRepo := ds2.NewUserRepository(dsCli)
	audioRepo := ds2.NewAudioRepository(dsCli)
	feedRepo := ds2.NewFeedRepository(dsCli)

	// middleware
	authenticator := middleware2.NewAuthenticator()

	// usecase
	audioUsecase := usecase2.NewAudioUsecase(gcsClient, audioRepo, feedRepo, userRepo)
	commentUsecase := usecase2.NewCommentUsecase(commentRepo, audioRepo)
	userUsecase := usecase2.NewUserUsecase(userRepo, audioRepo, feedRepo)
	feedUsecase := usecase2.NewFeedUsecase(feedRepo, audioRepo)

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

	resolver := graph2.NewResolver(userUsecase, audioUsecase, commentUsecase, feedUsecase)
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
