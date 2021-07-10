package entity

import (
	"cloud.google.com/go/datastore"
	"context"
)

// StarRepository interface
type StarRepository interface {
	Exists(ctx context.Context, userID *datastore.Key, audioID *datastore.Key) (bool, error)
	Find(ctx context.Context, userID int64, audioID string) (*Star, error)
	Save(context.Context, *Star) error
}
