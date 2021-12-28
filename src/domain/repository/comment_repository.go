package repository

import (
	"context"
	"github.com/hayashiki/audiy-api/src/domain/model"
	"go.mercari.io/datastore/boom"
)

// CommentRepository interface
type CommentRepository interface {
	GetAllByAudio(
		ctx context.Context,
		audioID string,
		cursor string,
		limit int,
		orderBy string) ([]*model.Comment, string, bool, error)
	PutTx(tx *boom.Transaction, comment *model.Comment) error
	DeleteTx(tx *boom.Transaction, id int64) error
}
