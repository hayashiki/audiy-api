package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/hayashiki/audiy-api/domain/entity"
	"github.com/hayashiki/audiy-api/interfaces/api/graph/generated"
)

func (r *playResolver) User(ctx context.Context, obj *entity.Play) (*entity.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *playResolver) Audio(ctx context.Context, obj *entity.Play) (*entity.Audio, error) {
	panic(fmt.Errorf("not implemented"))
}

// Play returns generated.PlayResolver implementation.
func (r *Resolver) Play() generated.PlayResolver { return &playResolver{r} }

type playResolver struct{ *Resolver }
