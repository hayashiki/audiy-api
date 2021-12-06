package app

import (
	"context"
	"github.com/hayashiki/audiy-api/src/infrastructure/ffmpeg"
	"github.com/hayashiki/audiy-api/src/infrastructure/transcript"
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
	"github.com/hayashiki/audiy-api/src/middleware"
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
	gcsSvc         gcs.Service
	slackSvc       slack.Service
	commentRepo    entity.CommentRepository
	audioRepo      entity.AudioRepository
	userRepo       entity.UserRepository
	feedRepo       entity.FeedRepository
	resolver       *graph.Resolver
	authenticator  middleware.Authenticator
	// TODO: handler struct
	apiHandler     http.Handler
	graphQLHandler http.Handler
}

func (d *Dependency) Inject(conf config.Config) {
	// infrastructure
	dsCli, err := ds.NewClient(context.Background(), config.GetProject())
	if err != nil {
		log.Fatalf("failed to read datastore client")
	}
	// inject
	gcsSvc := gcs.NewService(conf.GCSInputAudioBucket)
	slackSvc := slack.NewClient(conf.SlackBotToken)

	proveSvc := ffmpeg.Service{}
	transcoder := ffmpeg.Transcoder{}
	transcriptSvc  := transcript.NewSpeechRecogniser()

	// repository
	commentRepo := ds.NewCommentRepository(dsCli)
	userRepo := ds.NewUserRepository(dsCli)
	audioRepo := ds.NewAudioRepository(dsCli)
	feedRepo := ds.NewFeedRepository(dsCli)
	transcriptRepo := ds.NewTranscriptRepository(dsCli)

	// middleware
	authenticator := middleware.NewAuthenticator()

	// usecase
	audioUsecase := usecase.NewAudioUsecase(gcsSvc, audioRepo, feedRepo, userRepo)
	commentUsecase := usecase.NewCommentUsecase(commentRepo, audioRepo)
	userUsecase := usecase.NewUserUsecase(userRepo, audioRepo, feedRepo)
	feedUsecase := usecase.NewFeedUsecase(feedRepo, audioRepo)
	transcriptUsecase := usecase.NewTranscriptAudioUsecase(gcsSvc, audioRepo, transcriptRepo, proveSvc, transcoder, transcriptSvc)

	logger := logging.NewLogger(conf.IsDev)

	d.log = logger
	d.audioUsecase = audioUsecase
	d.userUsecase = userUsecase
	d.feedUsecase = feedUsecase
	d.commentUsecase = commentUsecase
	d.gcsSvc = gcsSvc
	d.slackSvc = slackSvc
	d.commentRepo = commentRepo
	d.audioRepo = audioRepo
	d.userRepo = userRepo
	d.feedRepo = feedRepo
	d.authenticator = authenticator

	resolver := graph.NewResolver(userUsecase, audioUsecase, commentUsecase, feedUsecase, transcriptUsecase)
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

	apiHandler := handler.NewAPIHandler(slackSvc, gcsSvc, proveSvc, audioRepo, feedRepo, userRepo)
	d.apiHandler = apiHandler
}
