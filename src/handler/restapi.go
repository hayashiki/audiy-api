package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	entity2 "github.com/hayashiki/audiy-api/src/domain/entity"
	gcs2 "github.com/hayashiki/audiy-api/src/infrastructure/gcs"
	slack2 "github.com/hayashiki/audiy-api/src/infrastructure/slack"
	usecase2 "github.com/hayashiki/audiy-api/src/usecase"

	importer "github.com/hayashiki/audiy-importer"
)

type APIHandler struct {
	slackSvc  slack2.Slack
	gcsSvc    gcs2.Client
	audioRepo entity2.AudioRepository
	feedRepo  entity2.FeedRepository
	userRepo  entity2.UserRepository
}

type PubSubMessage struct {
	Message struct {
		Data []byte `json:"data,omitempty"`
		ID   string `json:"id"`
	} `json:"message"`
}

// NewAPIHandler returns rest api
func NewAPIHandler(
	slackSvc slack2.Slack,
	gcsSvc gcs2.Client,
	audioRepo entity2.AudioRepository,
	feedRepo entity2.FeedRepository,
	userRepo entity2.UserRepository,
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

		input := &usecase2.AudioInput{
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

		auc := usecase2.NewAudio(h.slackSvc, h.gcsSvc, h.audioRepo, h.feedRepo, h.userRepo)
		if err := auc.Do(context.Background(), input); err != nil {
			log.Print(err)
			return
		}
	}
}
