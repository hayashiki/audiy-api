package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"log"

	"github.com/99designs/gqlgen/graphql"

	"github.com/hayashiki/audiy-api/src/validator"

	"github.com/hayashiki/audiy-api/cmd/errs"

	"github.com/hayashiki/audiy-api/src/domain/entity"
	"github.com/hayashiki/audiy-api/src/graph/auth"
	"github.com/hayashiki/audiy-api/src/graph/generated"
)

func (r *feedResolver) Audio(ctx context.Context, obj *entity.Feed) (*entity.Audio, error) {
	return r.audioUsecase.Get(ctx, obj.AudioKey.Name)
}

func (r *feedResolver) User(ctx context.Context, obj *entity.Feed) (*entity.User, error) {
	// TODO: delete
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateFeed(ctx context.Context, input entity.CreateFeedInput) (*entity.Feed, error) {
	graphql.AddErrorf(ctx, "custom error")
	v := entity.User{
		Email: "hjrke",
		Name:  "klklklklklklkl",
	}

	if err := validator.Validate(v); err != nil {
		err = validator.ManageValidationsErrors(ctx, err)
		return nil, err
	}

	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateFeed(ctx context.Context, input entity.UpdateFeedInput) (*entity.Feed, error) {
	auth, err := auth.ForContext(ctx)
	if err != nil {
		return nil, err
	}
	return r.feedUseCase.Put(ctx, auth.ID, input.ID, input.Event)
}

func (r *mutationResolver) DeleteFeed(ctx context.Context, id string) (*entity.DeleteFeedResult, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Feeds(ctx context.Context, cursor *string, filter *entity.FeedEvent, limit *int, order *entity.AudioOrder) (*entity.FeedConnection, error) {
	var err error

	v := entity.User{
		Email: "hjrke",
		Name:  "klklklklklklkl",
	}

	if err := validator.Validate(v); err != nil {
		log.Println("before")
		log.Println(err)

		err = validator.ManageValidationsErrors(ctx, err)
		log.Println("after")
		log.Println(err)
		return nil, err
	}

	//if err = validate.Struct(user); err != nil {
	//	return nil, util.ValidationError(ctx, err.(validator.ValidationErrors))
	//}

	err = errs.BadCredencials(ctx)
	if err != nil {
		return nil, err
	}

	auth, err := auth.ForContext(ctx)
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
	return r.feedUseCase.GetConnection(ctx, auth.ID, *cursor, *limit, filter, orderStr)
}

// Feed returns generated.FeedResolver implementation.
func (r *Resolver) Feed() generated.FeedResolver { return &feedResolver{r} }

type feedResolver struct{ *Resolver }
