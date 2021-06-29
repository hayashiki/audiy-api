package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/hayashiki/audiy-api/domain/entity"
	"github.com/hayashiki/audiy-api/interfaces/api/graph/generated"
)

func (r *audioResolver) LikeCount(ctx context.Context, obj *entity.Audio) (int, error) {
	return 0, nil
}

func (r *audioResolver) PlayCount(ctx context.Context, obj *entity.Audio) (int, error) {
	return 0, nil
}

func (r *audioResolver) Played(ctx context.Context, obj *entity.Audio) (bool, error) {
	r.audioUsecase.Get(ctx, obj.ID)

	return true, nil
}

// Audio returns generated.AudioResolver implementation.
func (r *Resolver) Audio() generated.AudioResolver { return &audioResolver{r} }

type audioResolver struct{ *Resolver }
