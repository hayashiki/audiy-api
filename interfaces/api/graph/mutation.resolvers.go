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

func (r *mutationResolver) CreatePlay(ctx context.Context, input entity.UpdateAudioInput) (*entity.Audio, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateComment(ctx context.Context, input entity.UpdateAudioInput) (*entity.Comment, error) {
	auth, err := ForContext(ctx)
	if err != nil {
		return nil, err
	}
	return r.commentUsecase.Save(ctx, auth.ID, input.AudioID)
}

func (r *mutationResolver) ToggleStar(ctx context.Context, input entity.UpdateAudioInput) (*entity.ToggleStarResult, error) {
	auth, err := ForContext(ctx)
	if err != nil {
		return nil, err
	}

	exists, err := r.starUsecase.Exists(ctx, auth.ID, input.AudioID)

	if exists {
		star, err := r.starUsecase.Delete(ctx, auth.ID, input.AudioID)
		if err != nil {
			return nil, err
		}
		return &entity.ToggleStarResult{
			Star:    star,
			Action:  "deleted",
			Success: false,
		}, nil

	}
	star, err := r.starUsecase.Save(ctx, auth.ID, input.AudioID)
	if err != nil {
		return nil, err
	}
	return &entity.ToggleStarResult{
		Star:    star,
		Action:  "created",
		Success: false,
	}, nil
}

func (r *mutationResolver) ToggleLike(ctx context.Context, input entity.UpdateAudioInput) (*entity.ToggleLikeResult, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *mutationResolver) CreateUserAudio(ctx context.Context, input entity.UserAudioInput) (*entity.Audio, error) {
	auth, err := ForContext(ctx)
	if err != nil {
		return nil, err
	}
	r.audioUserUsecase.Save(ctx, auth.ID, input.AudioID)

	// TODO: set create success response
	return nil, nil
}
