package entity

import (
	"cloud.google.com/go/datastore"
	"context"
)

// AudioRepository interface
type AudioRepository interface {
	Exists(ctx context.Context, id string) bool
	FindAll(ctx context.Context, cursor string, limit int, sort ...string) ([]*Audio, string, error)
	Find(ctx context.Context, id string) (*Audio, error)
	Save(ctx context.Context, audio *Audio) error
	Delete(ctx context.Context, audioKey *datastore.Key) error
}
