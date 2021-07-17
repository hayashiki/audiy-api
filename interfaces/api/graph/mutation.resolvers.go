package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/hayashiki/audiy-api/domain/entity"
	auth2 "github.com/hayashiki/audiy-api/interfaces/api/graph/auth"
	"github.com/hayashiki/audiy-api/interfaces/api/graph/generated"
)

func (r *mutationResolver) CreateAudio(ctx context.Context, input entity.AudiosInput) (*entity.Audio, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreatePlay(ctx context.Context, input entity.UpdateAudioInput) (*entity.Play, error) {
	auth, err := auth2.ForContext(ctx)
	if err != nil {
		return nil, err
	}
	return r.playUsecase.Save(ctx, auth.ID, input.AudioID)
}

func (r *mutationResolver) CreateComment(ctx context.Context, input entity.CreateCommentInput) (*entity.Comment, error) {
	auth, err := auth2.ForContext(ctx)
	if err != nil {
		return nil, err
	}
	return r.commentUsecase.Save(ctx, auth.ID, input)
}

func (r *mutationResolver) UpdateComment(ctx context.Context, input entity.UpdateCommentInput) (*entity.Comment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteComment(ctx context.Context, id string) (*entity.DeleteCommentResult, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
