package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"log"

	"github.com/hayashiki/audiy-api/src/domain/model"
	"github.com/hayashiki/audiy-api/src/graph/auth"
	"github.com/hayashiki/audiy-api/src/graph/generated"
	"github.com/hayashiki/audiy-api/src/version"
)

func (r *queryResolver) Version(ctx context.Context) (*model.Version, error) {
	v := &model.Version{
		Version: version.Version,
	}
	return v, nil
}

func (r *queryResolver) Comments(ctx context.Context, audioID string, cursor *string, limit *int, order []string) (*model.CommentConnection, error) {
	_, err := auth.ForContext(ctx)
	if err != nil {
		return nil, err
	}

	if *cursor == "" {
		*cursor = ""
	}
	if *limit == 0 {
		*limit = 100
	}

	log.Printf("%s", order)

	//if len() == 0 {
	orderBy := "CreatedAt"
	//}

	return r.commentUsecase.GetConnection(ctx, audioID, *cursor, *limit, orderBy)
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
