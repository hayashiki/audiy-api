package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/hayashiki/audiy-api/src/domain/model"
	"github.com/hayashiki/audiy-api/src/graph/auth"
	"github.com/hayashiki/audiy-api/src/graph/generated"
)

func (r *feedResolver) ID(ctx context.Context, obj *model.Feed) (string, error) {
	return string(obj.ID()), nil
}

func (r *feedResolver) Audio(ctx context.Context, obj *model.Feed) (*model.Audio, error) {
	return r.dataLoaders.AudioGetByID(ctx, obj.AudioID)
	//return r.audioUsecase.Get(ctx, obj.AudioKey.Name)
}

func (r *feedResolver) User(ctx context.Context, obj *model.Feed) (*model.User, error) {
	// TODO: delete
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateFeed(ctx context.Context, input model.CreateFeedInput) (*model.Feed, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateFeed(ctx context.Context, input model.UpdateFeedInput) (*model.Feed, error) {
	auth, err := auth.ForContext(ctx)
	if err != nil {
		return nil, err
	}
	return r.feedUseCase.Put(ctx, auth.ID, input.ID, input.Event)
}

func (r *mutationResolver) DeleteFeed(ctx context.Context, id string) (*model.DeleteFeedResult, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Feeds(ctx context.Context, cursor *string, filter *model.FeedEvent, limit *int, order *model.AudioOrder) (*model.FeedConnection, error) {
	auth, err := auth.ForContext(ctx)
	if err != nil {
		return nil, err
	}

	if *cursor == "" {
		*cursor = ""
	}

	var orderBy string
	if order == nil {
		orderBy = "-PublishedAt"
	} else if order.String() == model.AudioOrderPublishedAtAsc.String() {
		orderBy = "PublishedAt"
	} else if order.String() == model.AudioOrderPublishedAtDesc.String() {
		orderBy = "-PublishedAt"
	}
	return r.feedUseCase.GetConnection(ctx, auth.ID, *cursor, *limit, filter, orderBy)
}

// Feed returns generated.FeedResolver implementation.
func (r *Resolver) Feed() generated.FeedResolver { return &feedResolver{r} }

type feedResolver struct{ *Resolver }
