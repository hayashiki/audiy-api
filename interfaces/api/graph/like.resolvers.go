package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/hayashiki/audiy-api/domain/entity"
	auth2 "github.com/hayashiki/audiy-api/interfaces/api/graph/auth"
	"github.com/hayashiki/audiy-api/interfaces/api/graph/generated"
)

func (r *likeResolver) User(ctx context.Context, obj *entity.Like) (*entity.User, error) {
	return r.userUsecase.Get(ctx, obj.UserKey.Name)
}

func (r *likeResolver) Audio(ctx context.Context, obj *entity.Like) (*entity.Audio, error) {
	return r.audioUsecase.Get(ctx, obj.AudioKey.Name)
}

func (r *mutationResolver) CreateLike(ctx context.Context, input entity.UpdateAudioInput) (*entity.Like, error) {
	auth, err := auth2.ForContext(ctx)
	if err != nil {
		return nil, err
	}
	return r.likeUsecase.Save(ctx, auth.ID, input.AudioID)
}

func (r *mutationResolver) DeleteLike(ctx context.Context, input entity.UpdateAudioInput) (*entity.Like, error) {
	auth, err := auth2.ForContext(ctx)
	if err != nil {
		return nil, err
	}
	return r.likeUsecase.Delete(ctx, auth.ID, input.AudioID)
}

// Like returns generated.LikeResolver implementation.
func (r *Resolver) Like() generated.LikeResolver { return &likeResolver{r} }

type likeResolver struct{ *Resolver }
