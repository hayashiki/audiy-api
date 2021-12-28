package service

import (
	"context"
	"github.com/hayashiki/audiy-api/src/domain/model"
	"github.com/hayashiki/audiy-api/src/domain/repository"
)

type CommentService interface {
	GetEdges(ctx context.Context, audioID string, cursor string, limit int, order string) ([]*model.CommentEdge, string, bool, error)
}

type commentService struct {
	repo repository.CommentRepository
}

func NewCommentService(repo repository.CommentRepository) CommentService {
	return &commentService{
		repo: repo,
	}
}

func (s *commentService) GetEdges(ctx context.Context, audioID string, cursor string, limit int, order string) ([]*model.CommentEdge, string, bool, error) {
	commentTokens, nextCursor, hasMore, err := s.repo.GetAllByAudio(ctx, audioID, cursor, limit, order)
	if err != nil {
		return nil, "", false, err
	}
	edges := make([]*model.CommentEdge, len(commentTokens))
	for i, comment := range commentTokens {
		edges[i] = &model.CommentEdge{
			Cursor: nextCursor,
			Node:   comment,
		}
	}
	return edges, nextCursor, hasMore, nil
}

func (s *commentService) GetPageInfo(hasMore bool) *model.PageInfo {
	return &model.PageInfo{
		Cursor:    "",
		TotalPage: 0,
		HasMore:   hasMore,
	}
}

