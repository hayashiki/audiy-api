package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/hayashiki/audiy-api/src/domain/entity"
	"github.com/hayashiki/audiy-api/src/graph/generated"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input entity.CreateUserInput) (*entity.User, error) {
	return r.userUsecase.Save(ctx, input)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
