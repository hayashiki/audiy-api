package usecase

import (
	"context"
	"log"

	"github.com/hayashiki/audiy-api/domain/entity"
)

type CommentUsecase interface {
	Save(ctx context.Context, userID string, input entity.CreateCommentInput) (*entity.Comment, error)
	GetConnection(ctx context.Context, cursor string, limit int, order []string) (*entity.CommentConnection, error)
}

func NewCommentUsecase(commentRepo entity.CommentRepository) CommentUsecase {
	return &commentUsecase{commentRepo: commentRepo}
}

type commentUsecase struct {
	commentRepo entity.CommentRepository
}

func (c *commentUsecase) Save(ctx context.Context, userID string, input entity.CreateCommentInput) (*entity.Comment, error) {
	newComment := entity.NewComment(userID, input.AudioID, input.Body)
	err := c.commentRepo.Save(ctx, newComment)
	if err != nil {
		return nil, err
	}
	log.Printf("comment +%v", newComment)
	return newComment, nil
}

func (c *commentUsecase) GetConnection(ctx context.Context, cursor string, limit int, order []string) (*entity.CommentConnection, error) {
	comments, nextCursor, err := c.commentRepo.GetAll(ctx, cursor, limit, order...)
	if err != nil {
		return nil, err
	}
	commentEdges := make([]*entity.CommentEdge, len(comments))
	for i, a := range comments {
		commentEdges[i] = &entity.CommentEdge{
			Cursor: nextCursor,
			Node:   a,
		}
	}
	return &entity.CommentConnection{
		PageInfo: &entity.PageInfo{
			Cursor:  nextCursor,
			HasMore: len(comments) != 0,
		},
		Edges: commentEdges,
	}, nil
}
