package repository

import (
	"context"
	"github.com/hayashiki/audiy-api/src/domain/model"
	"go.mercari.io/datastore/boom"
)

// AudioRepository interface
type AudioRepository interface {
	Exists(ctx context.Context, id string) (bool, error)
	GetAll(
		ctx context.Context,
		cursor string,
		limit int,
		orderBy string) ([]*model.Audio, string, bool, error)
	GetMulti(ctx context.Context, ids []string) ([]*model.Audio, error)
	Get(ctx context.Context, id string) (*model.Audio, error)
	Put(ctx context.Context, audio *model.Audio) error
	PutTx(tx *boom.Transaction, item *model.Audio) error
	DeleteTx(tx *boom.Transaction, id string) error
}
