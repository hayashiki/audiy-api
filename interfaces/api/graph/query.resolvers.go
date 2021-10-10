package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"log"

	"github.com/hayashiki/audiy-api/domain/entity"
	"github.com/hayashiki/audiy-api/etc/version"
	auth2 "github.com/hayashiki/audiy-api/interfaces/api/graph/auth"
	"github.com/hayashiki/audiy-api/interfaces/api/graph/generated"
)

func (r *queryResolver) Version(ctx context.Context) (*entity.Version, error) {
	v := &entity.Version{
		Version: version.Version,
	}
	return v, nil
}

func (r *queryResolver) Comments(ctx context.Context, audioID string, cursor *string, limit *int, order []string) (*entity.CommentConnection, error) {
	auth, err := auth2.ForContext(ctx)
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
	order = []string{"created_at"}
	//}

	return r.commentUsecase.GetConnection(ctx, auth.ID, audioID, *cursor, *limit, order)
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
