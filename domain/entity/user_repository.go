package entity

import "context"

// UserRepository interface
type UserRepository interface {
	Save(context.Context, *User) error
	Get(ctx context.Context, id int64) (*User, error)
}
