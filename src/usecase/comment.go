package usecase

import (
	"context"
	"github.com/hayashiki/audiy-api/src/domain/repository"
	"github.com/hayashiki/audiy-api/src/infrastructure/datastore"
	"go.mercari.io/datastore/boom"

	"github.com/hayashiki/audiy-api/src/domain/model"
)

type CommentUsecase interface {
	Save(ctx context.Context, userID string, input model.CreateCommentInput) (*model.Comment, error)
	GetConnection(ctx context.Context, audioID string, cursor string, limit int, orderBy string) (*model.CommentConnection, error)
}

func NewCommentUsecase(
	transactor datastore.Transactor,
	commentRepo repository.CommentRepository,
	audioRepo repository.AudioRepository) CommentUsecase {
	return &commentUsecase{
		transactor: transactor,
		commentRepo: commentRepo,
		audioRepo: audioRepo,
	}
}

type commentUsecase struct {
	transactor datastore.Transactor
	commentRepo repository.CommentRepository
	audioRepo   repository.AudioRepository
}

func (c *commentUsecase) Save(ctx context.Context, userID string, input model.CreateCommentInput) (*model.Comment, error) {
	comment := &model.Comment{}
	if err := c.transactor.RunInTransaction(ctx, func(tx *boom.Transaction) error {
		audio, err := c.audioRepo.Get(ctx, input.AudioID)
		if err != nil {
			return err
		}
		comment = model.NewComment(userID, audio.ID, input.Body)
		if err := c.commentRepo.PutTx(tx, comment); err != nil {
			return err
		}
		// TODO: create increment method
		audio.CommentCount += 1
		if err := c.audioRepo.Put(ctx, audio); err != nil {
			return err
		}
		return nil
	}); err != nil {
		// TODO: rollback
		return nil, err
	}
	return comment, nil
}

func (c *commentUsecase) GetConnection(ctx context.Context, audioID string, cursor string, limit int, orderBy string) (*model.CommentConnection, error) {
	comments, nextCursor, hasMore, err := c.commentRepo.GetAllByAudio(ctx, audioID, cursor, limit, orderBy)
	if err != nil {
		return nil, err
	}
	commentEdges := make([]*model.CommentEdge, len(comments))
	for i, c := range comments {
		commentEdges[i] = &model.CommentEdge{
			// TODO: カーソル消したい
			Cursor: nextCursor,
			Node:   c,
		}
	}
	return &model.CommentConnection{
		PageInfo: &model.PageInfo{
			Cursor:  nextCursor,
			HasMore: hasMore,
		},
		Edges: commentEdges,
	}, nil
}
