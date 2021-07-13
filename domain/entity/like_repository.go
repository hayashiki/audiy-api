package entity

import (
	"context"
)

// LikeRepository interface
type LikeRepository interface {
	Exists(ctx context.Context, userID int64, audioID string) (bool, error)
	Find(ctx context.Context, userID int64, audioID string) (*Like, error)
	FindByRel(ctx context.Context, userID int64, audioID string) (*Star, error)
	Save(context.Context, *Like) error
}
