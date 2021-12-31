package repository

import (
	"context"
	"github.com/hayashiki/audiy-api/src/domain/model"
	"go.mercari.io/datastore/boom"
)

// UserRepository interface
type UserRepository interface {
	PutTx(tx *boom.Transaction, item *model.User) error
	Get(ctx context.Context, id string) (*model.User, error)
	GetAll(ctx context.Context) ([]*model.User, error)
	Exists(ctx context.Context, id string) (bool, error)
}
