package entity

import "context"

// AudioUserRepository interface
type UserRepository interface {
	Exists(ctx context.Context, userID int64) (bool, error)
	Save(context.Context, *User) error
}
