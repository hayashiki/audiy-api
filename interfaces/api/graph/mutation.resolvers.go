package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/hayashiki/audiy-api/domain/entity"
	"github.com/hayashiki/audiy-api/interfaces/api/graph/generated"
)

func (r *mutationResolver) CreateAudio(ctx context.Context, input entity.AudiosInput) (*entity.Audio, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateUserAudio(ctx context.Context, input entity.UserAudioInput) (*entity.Audio, error) {
	auth, err := ForContext(ctx)
	if err != nil {
		return nil, err
	}
	r.audioUserUsecase.Save(ctx, auth.ID, input.AudioID)

	// TODO: set create success response
	return nil, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
