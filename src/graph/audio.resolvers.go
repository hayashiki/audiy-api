package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hayashiki/audiy-api/src/domain/model"
	"github.com/hayashiki/audiy-api/src/graph/auth"
	"github.com/hayashiki/audiy-api/src/graph/generated"
	"github.com/hayashiki/audiy-api/src/infrastructure/gcs"
)

func (r *audioResolver) URL(ctx context.Context, obj *model.Audio) (string, error) {
	filePath := gcs.StorageObjectFilePath(obj.ID, "m4a")
	// TODO: read from config
	bucketName := os.Getenv("GCS_INPUT_AUDIO_BUCKET")

	return fmt.Sprintf("https://storage.cloud.google.com/%s/%s?authuser=", bucketName, filePath), nil
	// SignedURLをやめる
	//return gcs.GetGCSSignedURL(context.Background(), bucketName, filePath, "GET", "")
}

func (r *mutationResolver) CreateAudio(ctx context.Context, input *model.CreateAudioInput) (*model.Audio, error) {
	log.Println("CreateAudio")
	log.Println(input)
	return r.audioUsecase.CreateAudio(ctx, input)
}

func (r *mutationResolver) UploadAudio(ctx context.Context, input *model.UploadAudioInput) (*model.Audio, error) {
	return r.audioUsecase.UploadAudio(ctx, input)
}

func (r *queryResolver) Audio(ctx context.Context, id string) (*model.Audio, error) {
	return r.audioUsecase.Get(ctx, id)
}

func (r *queryResolver) Audios(ctx context.Context, cursor *string, filter *model.AudioFilter, limit *int, order *model.AudioOrder) (*model.AudioConnection, error) {
	_, err := auth.ForContext(ctx)
	if err != nil {
		return nil, err
	}

	if *cursor == "" {
		*cursor = ""
	}
	var orderBy string
	if order == nil {
		orderBy = "-PublishedAt"
	} else if order.String() == model.AudioOrderPublishedAtAsc.String() {
		orderBy = "PublishedAt"
	} else if order.String() == model.AudioOrderPublishedAtDesc.String() {
		orderBy = "-PublishedAt"
	}

	return r.audioUsecase.GetConnection(ctx, *cursor, *limit, orderBy)
}

// Audio returns generated.AudioResolver implementation.
func (r *Resolver) Audio() generated.AudioResolver { return &audioResolver{r} }

type audioResolver struct{ *Resolver }
