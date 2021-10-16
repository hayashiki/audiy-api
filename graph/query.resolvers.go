package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"log"

	auth3 "github.com/hayashiki/audiy-api/graph/auth"
	generated2 "github.com/hayashiki/audiy-api/graph/generated"

	"github.com/hayashiki/audiy-api/domain/entity"
	"github.com/hayashiki/audiy-api/etc/version"
)

func (r *queryResolver) Version(ctx context.Context) (*entity.Version, error) {
	v := &entity.Version{
		Version: version.Version,
	}
	return v, nil
}

func (r *queryResolver) Comments(ctx context.Context, audioID string, cursor *string, limit *int, order []string) (*entity.CommentConnection, error) {
	auth, err := auth3.ForContext(ctx)
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
func (r *Resolver) Query() generated2.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
