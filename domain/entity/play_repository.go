package entity

import (
	"context"
)

// PlayRepository interface
type PlayRepository interface {
	Exists(ctx context.Context, userID int64, audioID string) (bool, error)
	Find(ctx context.Context, userID int64, audioID string) (*Play, error)
	Save(context.Context, *Play) error
}
