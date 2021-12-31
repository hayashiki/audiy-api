package repository

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/hayashiki/audiy-api/src/domain/model"
)

// FCMRepository interface
type FCMRepository interface {
	GetAll(ctx context.Context, cursor string, limit int, orderBy string) ([]*model.Fcm, string, bool, error)
	PutTx(tx *datastore.Transaction, item *model.Fcm) error
	Delete(tx *datastore.Transaction, id string) error
}
