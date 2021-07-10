package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"log"

	"github.com/hayashiki/audiy-api/domain/entity"
	"github.com/hayashiki/audiy-api/interfaces/api/graph/generated"
)

func (r *commentResolver) User(ctx context.Context, obj *entity.Comment) (*entity.User, error) {
	log.Println(obj)
	return r.userUsecase.Get(ctx, obj.UserKey.ID)
}

// Comment returns generated.CommentResolver implementation.
func (r *Resolver) Comment() generated.CommentResolver { return &commentResolver{r} }

type commentResolver struct{ *Resolver }
