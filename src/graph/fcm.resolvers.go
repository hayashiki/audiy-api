package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/hayashiki/audiy-api/src/domain/model"
)

func (r *mutationResolver) CreateFcm(ctx context.Context, input model.CreateFCMInput) (*model.Fcm, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateFcm(ctx context.Context, input model.UpdateFCMInput) (*model.Fcm, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteFcm(ctx context.Context, id string) (*model.DeleteFCMResult, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Fcm(ctx context.Context, id string) (*model.Fcm, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Fcms(ctx context.Context, cursor *string, filter *model.AudioFilter, limit *int, order *model.AudioOrder) (*model.FCMConnection, error) {
	if *cursor == "" {
		*cursor = ""
	}
	if *limit == 0 {
		*limit = 100
	}
	return r.fcmUsecase.GetConnection(ctx, *cursor, *limit, "CreatedAt")
}
