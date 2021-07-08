package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/hayashiki/audiy-api/domain/entity"
	"github.com/hayashiki/audiy-api/interfaces/api/graph/generated"
)

func (r *likeResolver) User(ctx context.Context, obj *entity.Like) (*entity.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *likeResolver) Audio(ctx context.Context, obj *entity.Like) (*entity.Audio, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *starResolver) User(ctx context.Context, obj *entity.Star) (*entity.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *starResolver) Audio(ctx context.Context, obj *entity.Star) (*entity.Audio, error) {
	panic(fmt.Errorf("not implemented"))
}

// Like returns generated.LikeResolver implementation.
func (r *Resolver) Like() generated.LikeResolver { return &likeResolver{r} }

// Star returns generated.StarResolver implementation.
func (r *Resolver) Star() generated.StarResolver { return &starResolver{r} }

type likeResolver struct{ *Resolver }
type starResolver struct{ *Resolver }
