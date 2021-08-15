package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/hayashiki/audiy-api/domain/entity"
	"github.com/hayashiki/audiy-api/interfaces/api/graph/generated"
)

func (r *playResolver) User(ctx context.Context, obj *entity.Play) (*entity.User, error) {
	return r.userUsecase.Get(ctx, obj.UserKey.Name)
}

func (r *playResolver) Audio(ctx context.Context, obj *entity.Play) (*entity.Audio, error) {
	return r.audioUsecase.Get(ctx, obj.AudioKey.Name)
}

// Play returns generated.PlayResolver implementation.
func (r *Resolver) Play() generated.PlayResolver { return &playResolver{r} }

type playResolver struct{ *Resolver }
