package entity

import (
	"context"

	"cloud.google.com/go/datastore"
)

// FeedRepository interface
type FeedRepository interface {
	Exists(ctx context.Context, id int64, userID string) bool
	FindAll(ctx context.Context, userID string, filters map[string]interface{}, cursor string, limit int, sort ...string) ([]*Feed, string, error)
	Find(ctx context.Context, id int64, userID string) (*Feed, error)
	FindByAudio(ctx context.Context, userID string, audioID string) (*Feed, error)
	Save(ctx context.Context, userID string, feed *Feed) error
	SaveAll(ctx context.Context, userIDs []string, feeds []*Feed) error
	Delete(ctx context.Context, FeedKey *datastore.Key) error
}
