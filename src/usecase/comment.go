package usecase

import (
	"context"
	"log"

	entity2 "github.com/hayashiki/audiy-api/src/domain/entity"
)

type CommentUsecase interface {
	Save(ctx context.Context, userID string, input entity2.CreateCommentInput) (*entity2.Comment, error)
	GetConnection(ctx context.Context, userID string, audioID string, cursor string, limit int, order []string) (*entity2.CommentConnection, error)
}

func NewCommentUsecase(commentRepo entity2.CommentRepository, audioRepo entity2.AudioRepository) CommentUsecase {
	return &commentUsecase{commentRepo: commentRepo, audioRepo: audioRepo}
}

type commentUsecase struct {
	commentRepo entity2.CommentRepository
	audioRepo   entity2.AudioRepository
}

func (c *commentUsecase) Save(ctx context.Context, userID string, input entity2.CreateCommentInput) (*entity2.Comment, error) {
	newComment := entity2.NewComment(userID, input.AudioID, input.Body)
	err := c.commentRepo.Save(ctx, newComment)
	if err != nil {
		return nil, err
	}

	log.Println(input.AudioID)
	log.Println(c.audioRepo)

	audio, err := c.audioRepo.Find(ctx, input.AudioID)
	if err != nil {
		return nil, err
	}
	audio.CommentCount += 1
	if err := c.audioRepo.Save(ctx, audio); err != nil {
		return nil, err
	}

	return newComment, nil
}

func (c *commentUsecase) GetConnection(ctx context.Context, userID string, audioID string, cursor string, limit int, order []string) (*entity2.CommentConnection, error) {
	comments, nextCursor, err := c.commentRepo.GetAll(ctx, userID, audioID, cursor, limit, order...)
	if err != nil {
		return nil, err
	}
	commentEdges := make([]*entity2.CommentEdge, len(comments))
	for i, a := range comments {
		commentEdges[i] = &entity2.CommentEdge{
			Cursor: nextCursor,
			Node:   a,
		}
	}
	return &entity2.CommentConnection{
		PageInfo: &entity2.PageInfo{
			Cursor:  nextCursor,
			HasMore: len(comments) != 0,
		},
		Edges: commentEdges,
	}, nil
}
