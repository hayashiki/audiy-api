package usecase

import (
	"context"
	"github.com/hayashiki/audiy-api/src/domain/model"
	"github.com/hayashiki/audiy-api/src/domain/service"
)

type FcmUsecase interface {
	GetConnection(ctx context.Context, cursor string, limit int, order string) (*model.FCMConnection, error)
}

func NewFcmUsecase(
	fcmSvc service.FCMService,
) FcmUsecase {
	return &fcmUsecase{
		fcmSvc: fcmSvc,
	}
}

type fcmUsecase struct {
	fcmSvc service.FCMService
}

func (u *fcmUsecase) GetConnection(ctx context.Context, cursor string, limit int, order string) (*model.FCMConnection, error) {
	edges, nextCursor, hasMore, err := u.fcmSvc.GetEdges(ctx, cursor, limit, order)
	if err != nil {
		return nil, err
	}
	return &model.FCMConnection{
		PageInfo: &model.PageInfo{
			Cursor:  nextCursor,
			HasMore: hasMore,
		},
		Edges: edges,
	}, nil
}
