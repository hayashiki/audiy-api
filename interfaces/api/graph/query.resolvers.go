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

func (r *queryResolver) Audio(ctx context.Context, id string) (*entity.Audio, error) {
	return r.audioUsecase.Get(ctx, id)
}

func (r *queryResolver) Audios(ctx context.Context, cursor *string, limit *int, order []string) (*entity.AudioConnection, error) {
	_, err := auth2.ForContext(ctx)
	if err != nil {
		return nil, err
	}

	if *cursor == "" {
		*cursor = ""
	}
	if *limit == 0 {
		*limit = 100
	}
	if len(order) == 0 {
		order = []string{"-published_at"}
	}

	return r.audioUsecase.GetConnection(ctx, *cursor, *limit, order)
}

func (r *queryResolver) Version(ctx context.Context) (*entity.Version, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Comments(ctx context.Context, cursor *string, limit *int, order []string) (*entity.AudioConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
