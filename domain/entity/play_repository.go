package entity

import (
	"context"
)

// AudioUserRepository interface
type AudioUserRepository interface {
	Exists(ctx context.Context, userID int64, audioID string) (bool, error)
	Find(ctx context.Context, userID int64, audioID string) (*Play, error)
	Save(context.Context, *Play) error
}
