package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/hayashiki/audiy-api/domain/entity"
	"github.com/hayashiki/audiy-api/interfaces/api/graph/generated"
)

func (r *starResolver) User(ctx context.Context, obj *entity.Star) (*entity.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *starResolver) Audio(ctx context.Context, obj *entity.Star) (*entity.Audio, error) {
	panic(fmt.Errorf("not implemented"))
}

// Star returns generated.StarResolver implementation.
func (r *Resolver) Star() generated.StarResolver { return &starResolver{r} }

type starResolver struct{ *Resolver }
