package registry

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"go.opencensus.io/plugin/ochttp"

	"github.com/hayashiki/audiy-api/domain/entity"

	importer "github.com/hayashiki/audiy-importer"

	"github.com/hayashiki/audiy-api/etc/config"

	"github.com/hayashiki/audiy-api/infrastructure/slack"

	"github.com/hayashiki/audiy-api/infrastructure/gcs"

	"github.com/hayashiki/audiy-api/interfaces/api/graph"

	"github.com/hayashiki/audiy-api/interfaces/middleware"

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
	conf, err := config.NewConf()
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
	playRepo := ds.NewPlayRepository(dsCli)
	commentRepo := ds.NewCommentRepository(dsCli)
	userRepo := ds.NewUserRepository(dsCli)
	audioRepo := ds.NewAudioRepository(dsCli)
	likeRepo := ds.NewLikeRepository(dsCli)
	starRepo := ds.NewStarRepository(dsCli)
	feedRepo := ds.NewFeedRepository(dsCli)

	// middleware
	authenticator := middleware.NewAuthenticator()

	// usecase
	audioUsecase := usecase.NewAudioUsecase(gcsClient, audioRepo, feedRepo, userRepo)
	playUsecase := usecase.NewPlayUsecase(playRepo)
	commentUsecase := usecase.NewCommentUsecase(commentRepo, audioRepo)
	userUsecase := usecase.NewUserUsecase(userRepo, audioRepo, feedRepo)
	likeUsecase := usecase.NewLikeUsecase(likeRepo)
	starUsecase := usecase.NewStarUsecase(starRepo)
	feedUsecase := usecase.NewFeedUsecase(feedRepo, audioRepo)

	// handler
	resolver := graph.NewResolver(userUsecase, audioUsecase, playUsecase, starUsecase, likeUsecase, commentUsecase, feedUsecase)

	queryHandler := handler.NewQueryHandler(resolver)

	queryHandler = &ochttp.Handler{
		Handler:     queryHandler,
		Propagation: &propagation.HTTPFormat{},
	}

	rootHandler := handler.NewRootHandler()
	rootHandler = &ochttp.Handler{
		Handler:     rootHandler,
		Propagation: &propagation.HTTPFormat{},
	}

	apiHandler := NewAPIHandler(slackSvc, gcsClient, audioRepo, feedRepo, userRepo)
	// router
	router := router.NewRouter(rootHandler, authenticator.AuthMiddleware(queryHandler), authenticator.AuthMiddleware(queryHandler), apiHandler)

	return router.CreateHandler()
}

type APIHandler struct {
	slackSvc  slack.Slack
	gcsSvc    gcs.Client
	audioRepo entity.AudioRepository
	feedRepo  entity.FeedRepository
	userRepo  entity.UserRepository
}

type PubSubMessage struct {
	Message struct {
		Data []byte `json:"data,omitempty"`
		ID   string `json:"id"`
	} `json:"message"`
}

// NewAPIHandler returns rest api
func NewAPIHandler(
	slackSvc slack.Slack,
	gcsSvc gcs.Client,
	audioRepo entity.AudioRepository,
	feedRepo entity.FeedRepository,
	userRepo entity.UserRepository,
) http.Handler {
	h := APIHandler{slackSvc: slackSvc, gcsSvc: gcsSvc, audioRepo: audioRepo, feedRepo: feedRepo, userRepo: userRepo}
	return h.Handler()
}

func (h *APIHandler) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("audio called")

		var m PubSubMessage
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			log.Fatalf("fail to parse HTTP body: %v", err)
			http.Error(w, "fail to parse HTTP body", http.StatusBadRequest)
		}
		var e importer.AudioEnqueueMessage
		if err := json.Unmarshal(m.Message.Data, &e); err != nil {
			log.Fatalf("json.Unmarshal: %v", err)
			http.Error(w, "fail to unmarshal data", http.StatusBadRequest)
			return
		}
		log.Printf("e is %+v", e)

		//http.Error(w, "hoge", http.StatusInternalServerError)

		input := &usecase.AudioInput{
			ID:                 e.ID,
			Name:               e.Name,
			Title:              e.Title,
			URLPrivateDownload: e.URLPrivateDownload,
			Created:            e.Created,
			Mimetype:           e.Mimetype,
		}

		if err := input.Validate(); err != nil {
			log.Printf("err %v", err)
			return
		}

		auc := usecase.NewAudio(h.slackSvc, h.gcsSvc, h.audioRepo, h.feedRepo, h.userRepo)
		if err := auc.Do(context.Background(), input); err != nil {
			log.Print(err)
			return
		}
	}
}
