package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/hayashiki/audiy-api/src/domain/entity"
	"github.com/hayashiki/audiy-api/src/graph/auth"
	"github.com/hayashiki/audiy-api/src/graph/generated"
)

func (r *commentResolver) User(ctx context.Context, obj *entity.Comment) (*entity.User, error) {
	return r.userUsecase.Get(ctx, obj.UserKey.Name)
}

func (r *commentResolver) Audio(ctx context.Context, obj *entity.Comment) (*entity.Audio, error) {
	return r.audioUsecase.Get(ctx, obj.AudioKey.Name)
}

func (r *mutationResolver) CreateComment(ctx context.Context, input entity.CreateCommentInput) (*entity.Comment, error) {
	auth, err := auth.ForContext(ctx)
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

// Comment returns generated.CommentResolver implementation.
func (r *Resolver) Comment() generated.CommentResolver { return &commentResolver{r} }

type commentResolver struct{ *Resolver }