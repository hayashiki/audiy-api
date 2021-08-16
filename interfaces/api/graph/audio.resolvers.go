package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"log"
	"os"

	"github.com/hayashiki/audiy-api/infrastructure/gcs"

	"github.com/hayashiki/audiy-api/domain/entity"
	auth2 "github.com/hayashiki/audiy-api/interfaces/api/graph/auth"
	"github.com/hayashiki/audiy-api/interfaces/api/graph/generated"
)

func (r *audioResolver) LikeCount(ctx context.Context, obj *entity.Audio) (int, error) {
	auth, err := auth2.ForContext(ctx)
	if err != nil {
		return 0, err
	}
	r.playUsecase.Exists(ctx, auth.ID, obj.ID)

	return 0, nil
}

func (r *audioResolver) PlayCount(ctx context.Context, obj *entity.Audio) (int, error) {
	return 0, nil
}

func (r *audioResolver) URL(ctx context.Context, obj *entity.Audio) (string, error) {
	filePath := gcs.StorageObjectFilePath(obj.ID, "m4a")
	// TODO: read from config
	bucketName := os.Getenv("GCS_INPUT_AUDIO_BUCKET")
	return gcs.GetGCSSignedURL(context.Background(), bucketName, filePath, "GET", "")
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

// Audio returns generated.AudioResolver implementation.
func (r *Resolver) Audio() generated.AudioResolver { return &audioResolver{r} }

type audioResolver struct{ *Resolver }
