package repository

import (
	"context"
	"github.com/hayashiki/audiy-api/src/domain/model"
)

// FeedRepository interface
type FeedRepository interface {
	Exists(ctx context.Context, userID string,id int64) (bool, error)
	GetAll(
		ctx context.Context,
		userID string,
		filters map[string]interface{},
		cursor string,
		limit int,
		orderBy string) ([]*model.Feed, string, bool, error)
	Get(ctx context.Context, userID string, id int64) (*model.Feed, error)
	GetByAudio(ctx context.Context, userID string, audioID string) (*model.Feed, error)
	Put(ctx context.Context, userID string, item *model.Feed) error
	PutMulti(ctx context.Context, feeds []*model.Feed) error
	//Delete(ctx context.Context, FeedKey *datastore.Key) error
}
