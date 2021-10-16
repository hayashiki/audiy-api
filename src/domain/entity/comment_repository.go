package entity

import (
	"context"
)

type Query struct {
	Kind    string
	Filters []*Filter
	Offset  int
	Cursor  string
	Limit   int
	Order   []string
}

type Filter struct {
	key   string
	value interface{}
}

// CommentRepository interface
type CommentRepository interface {
	GetAll(ctx context.Context, userID string, audioID string, cursor string, limit int, sort ...string) ([]*Comment, string, error)
	Save(ctx context.Context, comment *Comment) error
}
