package entity

import (
	"context"
)

// AudioRepository interface
type AudioRepository interface {
	Exists(context.Context, string) bool
	Find(context.Context, string) (*Audio, error)
	FindAll(ctx context.Context, cursor string, limit int, sort ...string) ([]*Audio, string, error)
	Save(context.Context, *Audio) error
	Remove(ctx context.Context, item *Audio) error
}
