package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hayashiki/audiy-api/src/domain/entity"
	"github.com/hayashiki/audiy-api/src/graph/auth"
	"github.com/hayashiki/audiy-api/src/graph/generated"
	"github.com/hayashiki/audiy-api/src/infrastructure/gcs"
)

func (r *audioResolver) URL(ctx context.Context, obj *entity.Audio) (string, error) {
	filePath := gcs.StorageObjectFilePath(obj.ID, "m4a")
	// TODO: read from config
	bucketName := os.Getenv("GCS_INPUT_AUDIO_BUCKET")

	return fmt.Sprintf("https://storage.cloud.google.com/%s/%s?authuser=", bucketName, filePath), nil
	// SignedURLをやめる
	//return gcs.GetGCSSignedURL(context.Background(), bucketName, filePath, "GET", "")
}

func (r *mutationResolver) CreateAudio(ctx context.Context, input *entity.CreateAudioInput) (*entity.Audio, error) {
	log.Println("CreateAudio")
	log.Println(input)
	return r.audioUsecase.CreateAudio(ctx, input)
}

func (r *mutationResolver) UploadAudio(ctx context.Context, input *entity.UploadAudioInput) (*entity.Audio, error) {
	return r.audioUsecase.UploadAudio(ctx, input)
}

func (r *queryResolver) Audio(ctx context.Context, id string) (*entity.Audio, error) {
	return r.audioUsecase.Get(ctx, id)
}

func (r *queryResolver) Audios(ctx context.Context, cursor *string, filter *entity.AudioFilter, limit *int, order *entity.AudioOrder) (*entity.AudioConnection, error) {
	_, err := auth.ForContext(ctx)
	if err != nil {
		return nil, err
	}

	if *cursor == "" {
		*cursor = ""
	}
	var orderStr []string
	log.Println("order", order)
	if order == nil {
		orderStr = []string{"-published_at"}
	} else if order.String() == entity.AudioOrderPublishedAtAsc.String() {
		orderStr = []string{"published_at"}
	} else if order.String() == entity.AudioOrderPublishedAtDesc.String() {
		orderStr = []string{"-published_at"}
	}
	//if *filter.Played {
	//	*filter.Played = false
	//}
	//if *filter.Liked {
	//	*filter.Liked = false
	//}
	//if *filter.Stared {
	//	*filter.Stared = false
	//}

	//log.Printf("filter.Liked %v", filter.Liked)
	//log.Printf("filter.Played %v", *filter.Played)
	//log.Printf("filter.Stared %v", *filter.Stared)

	return r.audioUsecase.GetConnection(ctx, *cursor, *limit, orderStr)
}

// Audio returns generated.AudioResolver implementation.
func (r *Resolver) Audio() generated.AudioResolver { return &audioResolver{r} }

type audioResolver struct{ *Resolver }
