package entity

import (
	"context"
)

// StarRepository interface
type StarRepository interface {
	Exists(ctx context.Context, userID string, audioID string) (bool, error)
	Find(ctx context.Context, id int64) (*Star, error)
	FindByRel(ctx context.Context, userID string, audioID string) (*Star, error)
	Save(context.Context, *Star) error
	Delete(ctx context.Context, id int64) error
}
