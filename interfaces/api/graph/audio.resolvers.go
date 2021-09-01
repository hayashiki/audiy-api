package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hayashiki/audiy-api/domain/entity"
	"github.com/hayashiki/audiy-api/infrastructure/gcs"
	auth2 "github.com/hayashiki/audiy-api/interfaces/api/graph/auth"
	"github.com/hayashiki/audiy-api/interfaces/api/graph/generated"
)

func (r *audioResolver) URL(ctx context.Context, obj *entity.Audio) (string, error) {
	filePath := gcs.StorageObjectFilePath(obj.ID, "m4a")
	// TODO: read from config
	bucketName := os.Getenv("GCS_INPUT_AUDIO_BUCKET")

	return fmt.Sprintf("https://storage.cloud.google.com/%s/%s?authuser=", bucketName, filePath), nil
	// SignedURLをやめる
	//return gcs.GetGCSSignedURL(context.Background(), bucketName, filePath, "GET", "")
}

func (r *audioResolver) Played(ctx context.Context, obj *entity.Audio) (bool, error) {
	auth, err := auth2.ForContext(ctx)
	if err != nil {
		return false, err
	}
	return r.playUsecase.Exists(ctx, auth.ID, obj.ID)
}

func (r *audioResolver) Liked(ctx context.Context, obj *entity.Audio) (bool, error) {
	auth, err := auth2.ForContext(ctx)
	if err != nil {
		return false, err
	}
	exists, err := r.likeUsecase.Exists(ctx, auth.ID, obj.ID)
	log.Printf("exists %v %v %v", exists, auth.ID, obj.ID)
	return exists, err
}

func (r *audioResolver) Stared(ctx context.Context, obj *entity.Audio) (bool, error) {
	auth, err := auth2.ForContext(ctx)
	if err != nil {
		return false, err
	}
	return r.starUsecase.Exists(ctx, auth.ID, obj.ID)
}

func (r *queryResolver) Audio(ctx context.Context, id string) (*entity.Audio, error) {
	return r.audioUsecase.Get(ctx, id)
}

func (r *queryResolver) Audios(ctx context.Context, cursor *string, filter *entity.AudioFilter, limit *int, order *entity.AudioOrder) (*entity.AudioConnection, error) {
	_, err := auth2.ForContext(ctx)
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
