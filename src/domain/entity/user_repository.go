package entity

import "context"

// UserRepository interface
type UserRepository interface {
	Save(context.Context, *User) error
	Get(ctx context.Context, id string) (*User, error)
	GetAll(ctx context.Context) ([]*User, error)
}
