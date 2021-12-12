package app

import (
	"github.com/hayashiki/audiy-api/src/domain/repository"
	"github.com/hayashiki/audiy-api/src/domain/service"
	"github.com/hayashiki/audiy-api/src/graph/dataloaders"
	"github.com/hayashiki/audiy-api/src/infrastructure/datastore/audio_entity"
	"github.com/hayashiki/audiy-api/src/infrastructure/datastore/comment_entity"
	"github.com/hayashiki/audiy-api/src/infrastructure/datastore/fcm_entity"
	"github.com/hayashiki/audiy-api/src/infrastructure/datastore/feed_entity"
	"github.com/hayashiki/audiy-api/src/infrastructure/datastore/transcript_entity"
	"github.com/hayashiki/audiy-api/src/infrastructure/datastore/user_entity"
	"github.com/hayashiki/audiy-api/src/infrastructure/ffmpeg"
	"github.com/hayashiki/audiy-api/src/infrastructure/transcript"
	"log"
	"net/http"

	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"github.com/hayashiki/audiy-api/src/config"
	"github.com/hayashiki/audiy-api/src/graph"
	"github.com/hayashiki/audiy-api/src/handler"
	"github.com/hayashiki/audiy-api/src/infrastructure/datastore"
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
	commentRepo    repository.CommentRepository
	audioRepo      repository.AudioRepository
	userRepo       repository.UserRepository
	feedRepo       repository.FeedRepository
	resolver       *graph.Resolver
	authenticator  middleware.Authenticator
	dataloaderSvc  dataloaders.DataLoaderService
	// TODO: handler struct
	apiHandler     http.Handler
	graphQLHandler http.Handler
}

func (d *Dependency) Inject(conf config.Config) {
	// infrastructure
	log.Println(config.GetProject())
	mDsClid := datastore.New()
	dsCli := datastore.NewDS()
	//if err != nil {
	//	log.Fatalf("failed to read datastore client")
	//}

	transaction := datastore.NewDatastoreTransactor()

	// inject
	gcsSvc := gcs.NewService(conf.GCSInputAudioBucket)
	slackSvc := slack.NewClient(conf.SlackBotToken)

	proveSvc := ffmpeg.Service{}
	transcoder := ffmpeg.Transcoder{}
	transcriptSvc  := transcript.NewSpeechRecogniser()

	// repository
	commentRepo := comment_entity.NewCommentRepository(mDsClid)
	userRepo := user_entity.NewUserRepository(mDsClid)
	audioRepo := audio_entity.NewAudioRepository(mDsClid)
	feedRepo := feed_entity.NewFeedRepository(dsCli)
	transcriptRepo := transcript_entity.NewTranscriptRepository(mDsClid)

	// middleware
	authenticator := middleware.NewAuthenticator()

	// domain/service
	ds := datastore.New()
	fcmRepo := fcm_entity.NewRepository(ds)
	fcmSvc := service.NewFcmService(fcmRepo)

	// usecase
	audioUsecase := usecase.NewAudioUsecase(gcsSvc, audioRepo, feedRepo, userRepo)
	commentUsecase := usecase.NewCommentUsecase(transaction, commentRepo, audioRepo)
	userUsecase := usecase.NewUserUsecase(userRepo, audioRepo, feedRepo)
	feedUsecase := usecase.NewFeedUsecase(feedRepo, audioRepo)
	transcriptUsecase := usecase.NewTranscriptAudioUsecase(gcsSvc, audioRepo, transcriptRepo, proveSvc, transcoder, transcriptSvc)
	fcmUsecase := usecase.NewFcmUsecase(fcmSvc)

	logger := logging.NewLogger(conf.IsDev)

	d.dataloaderSvc = dataloaders.NewDataLoaderService(audioRepo)

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

	resolver := graph.NewResolver(d.dataloaderSvc, userUsecase, audioUsecase, commentUsecase, feedUsecase, transcriptUsecase, fcmUsecase)
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
