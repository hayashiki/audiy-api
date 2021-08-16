package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"log"

	"github.com/hayashiki/audiy-api/domain/entity"
	auth2 "github.com/hayashiki/audiy-api/interfaces/api/graph/auth"
	"github.com/hayashiki/audiy-api/interfaces/api/graph/generated"
)

func (r *queryResolver) Audio(ctx context.Context, id string) (*entity.Audio, error) {
	return r.audioUsecase.Get(ctx, id)
}

func (r *queryResolver) Audios(ctx context.Context, cursor *string, filter *entity.AudioFilter, limit *int, order *entity.AudioOrder) (*entity.AudioConnection, error) {
	_, err := auth2.ForContext(ctx)
	if err != nil {
		return nil, err
	}

	if *cursor == "" {
		*cursor = ""
	}
	var orderStr []string
	log.Println("order", order)
	if order == nil {
		orderStr = []string{"-published_at"}
	} else if order.String() == entity.AudioOrderPublishedAtAsc.String() {
		orderStr = []string{"published_at"}
	} else if order.String() == entity.AudioOrderPublishedAtDesc.String() {
		orderStr = []string{"-published_at"}
	}
	//if *filter.Played {
	//	*filter.Played = false
	//}
	//if *filter.Liked {
	//	*filter.Liked = false
	//}
	//if *filter.Stared {
	//	*filter.Stared = false
	//}

	//log.Printf("filter.Liked %v", filter.Liked)
	//log.Printf("filter.Played %v", *filter.Played)
	//log.Printf("filter.Stared %v", *filter.Stared)

	return r.audioUsecase.GetConnection(ctx, *cursor, *limit, orderStr)
}

func (r *queryResolver) Version(ctx context.Context) (*entity.Version, error) {
	panic(fmt.Errorf("not implemented"))
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
