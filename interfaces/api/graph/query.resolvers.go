package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/hayashiki/audiy-api/interfaces/api/graph/generated"
	"github.com/hayashiki/audiy-api/interfaces/api/graph/model"
)

func (r *queryResolver) Audio(ctx context.Context, id string) (*model.Audio, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Audios(ctx context.Context, cursor string, orderBy []*model.AudioOrder) (*model.AudioConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
