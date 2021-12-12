package repository

import (
	"context"
	"github.com/hayashiki/audiy-api/src/domain/model"
	"go.mercari.io/datastore/boom"
)

// FCMRepository interface
type FCMRepository interface {
	GetAll(ctx context.Context, cursor string, limit int, orderBy string) ([]*model.Fcm, string, bool, error)
	Put(tx *boom.Transaction, item *model.Fcm) error
	Delete(tx *boom.Transaction, id string) error
}
