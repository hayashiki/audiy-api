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

func (r *feedResolver) Audio(ctx context.Context, obj *entity.Feed) (*entity.Audio, error) {
	return r.audioUsecase.Get(ctx, obj.AudioKey.Name)
}

func (r *feedResolver) User(ctx context.Context, obj *entity.Feed) (*entity.User, error) {
	// TODO: delete
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateFeed(ctx context.Context, input entity.CreateFeedInput) (*entity.Feed, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateFeed(ctx context.Context, input entity.UpdateFeedInput) (*entity.Feed, error) {
	auth, err := auth2.ForContext(ctx)
	if err != nil {
		return nil, err
	}
	return r.feedUseCase.Put(ctx, auth.ID, input.ID, input.Event)
}

func (r *mutationResolver) DeleteFeed(ctx context.Context, id string) (*entity.DeleteFeedResult, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Feeds(ctx context.Context, cursor *string, filter *entity.AudioFilter, limit *int, order *entity.AudioOrder) (*entity.FeedConnection, error) {
	auth, err := auth2.ForContext(ctx)
	if err != nil {
		return nil, err
	}

	if *cursor == "" {
		*cursor = ""
	}

	var orderStr []string
	if order == nil {
		orderStr = []string{"-published_at"}
	} else if order.String() == entity.AudioOrderPublishedAtAsc.String() {
		orderStr = []string{"published_at"}
	} else if order.String() == entity.AudioOrderPublishedAtDesc.String() {
		orderStr = []string{"-published_at"}
	}
	return r.feedUseCase.GetConnection(ctx, auth.ID, *cursor, *limit, orderStr)
}

// Feed returns generated.FeedResolver implementation.
func (r *Resolver) Feed() generated.FeedResolver { return &feedResolver{r} }

type feedResolver struct{ *Resolver }
