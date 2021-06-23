package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/hayashiki/audiy-api/domain/entity"
	"github.com/hayashiki/audiy-api/interfaces/api/graph/generated"
)

func (r *audioResolver) LikeCount(ctx context.Context, obj *entity.Audio) (int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *audioResolver) PlayCount(ctx context.Context, obj *entity.Audio) (int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *audioResolver) Minetype(ctx context.Context, obj *entity.Audio) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

// Audio returns generated.AudioResolver implementation.
func (r *Resolver) Audio() generated.AudioResolver { return &audioResolver{r} }

type audioResolver struct{ *Resolver }
