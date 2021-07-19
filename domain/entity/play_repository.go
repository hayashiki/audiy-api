package entity

import (
	"context"
)

// PlayRepository interface
type PlayRepository interface {
	Exists(ctx context.Context, userID string, audioID string) (bool, error)
	Find(ctx context.Context, userID string, audioID string) (*Play, error)
	Save(context.Context, *Play) error
}
