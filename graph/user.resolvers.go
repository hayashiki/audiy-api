package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	generated2 "github.com/hayashiki/audiy-api/graph/generated"

	"github.com/hayashiki/audiy-api/domain/entity"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input entity.CreateUserInput) (*entity.User, error) {
	return r.userUsecase.Save(ctx, input)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated2.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
