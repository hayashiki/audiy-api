package entity

import (
	"context"
)

// AudioUserRepository interface
type AudioUserRepository interface {
	Exists(ctx context.Context, userID int64, audioID string) (bool, error)
	Find(ctx context.Context, userID int64, audioID string) (*AudioUser, error)
	Save(context.Context, *AudioUser) error
}
