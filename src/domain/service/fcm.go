package service

import (
	"context"
	"github.com/hayashiki/audiy-api/src/domain/model"
	"github.com/hayashiki/audiy-api/src/domain/repository"
)

type FCMService interface {
	GetEdges(ctx context.Context, cursor string, limit int, order string) ([]*model.FCMEdge, string, bool, error)
	GetPageInfo(hasMore bool) *model.PageInfo
}

type fcmService struct {
	repo repository.FCMRepository
}

func NewFcmService(repo repository.FCMRepository) FCMService {
	return &fcmService{
		repo: repo,
	}
}

func (s *fcmService) GetEdges(ctx context.Context, cursor string, limit int, order string) ([]*model.FCMEdge, string, bool, error) {
	fcmTokens, nextCursor, hasMore, err := s.repo.GetAll(ctx, cursor, limit, order)
	if err != nil {
		return nil, "", false, err
	}
	edges := make([]*model.FCMEdge, len(fcmTokens))
	for i, fcm := range fcmTokens {
		edges[i] = &model.FCMEdge{
			Cursor: nextCursor,
			Node:   fcm,
		}
	}
	return edges, nextCursor, hasMore, nil
}

func (s *fcmService) GetPageInfo(hasMore bool) *model.PageInfo {
	return &model.PageInfo{
		Cursor:    "",
		TotalPage: 0,
		HasMore:   hasMore,
	}
}
