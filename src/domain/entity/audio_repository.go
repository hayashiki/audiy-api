package entity

import (
	"context"

	"cloud.google.com/go/datastore"
)

// AudioRepository interface
type AudioRepository interface {
	Exists(ctx context.Context, id string) bool
	FindAll(ctx context.Context, filters map[string]interface{}, cursor string, limit int, sort ...string) ([]*Audio, string, error)
	GetMulti(ctx context.Context, IDs []string) ([]*Audio, error)
	Find(ctx context.Context, id string) (*Audio, error)
	Save(ctx context.Context, audio *Audio) error
	Delete(ctx context.Context, audioKey *datastore.Key) error
}
