package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/hayashiki/audiy-api/domain/entity"
	auth2 "github.com/hayashiki/audiy-api/interfaces/api/graph/auth"
	"github.com/hayashiki/audiy-api/interfaces/api/graph/generated"
)

func (r *mutationResolver) CreateStar(ctx context.Context, input entity.UpdateAudioInput) (*entity.Star, error) {
	auth, err := auth2.ForContext(ctx)
	if err != nil {
		return nil, err
	}
	return r.starUsecase.Save(ctx, auth.ID, input.AudioID)
}

func (r *mutationResolver) DeleteStar(ctx context.Context, input entity.UpdateAudioInput) (*entity.Star, error) {
	auth, err := auth2.ForContext(ctx)
	if err != nil {
		return nil, err
	}
	return r.starUsecase.Delete(ctx, auth.ID, input.AudioID)
}

func (r *starResolver) User(ctx context.Context, obj *entity.Star) (*entity.User, error) {
	return r.userUsecase.Get(ctx, obj.UserKey.ID)
}

func (r *starResolver) Audio(ctx context.Context, obj *entity.Star) (*entity.Audio, error) {
	return r.audioUsecase.Get(ctx, obj.AudioKey.Name)
}

// Star returns generated.StarResolver implementation.
func (r *Resolver) Star() generated.StarResolver { return &starResolver{r} }

type starResolver struct{ *Resolver }
