package auth

import (
	"context"
	"errors"
)

type key struct {
	name string
}

type Auth struct {
	ID    string
	Email string
	Name  string
}

var (
	errNoUserInContext = errors.New("no user in context")
)

var KeyAuth = &key{"auth"}

func SetAuth(ctx context.Context, auth *Auth) context.Context {
	return context.WithValue(ctx, KeyAuth, auth)
}

func ForContext(ctx context.Context) (*Auth, error) {
	auth, ok := ctx.Value(KeyAuth).(*Auth)
	if !ok {
		return nil, errNoUserInContext
	}
	return auth, nil
}

// ErrExpiredToken is the error returned if the token has expired
type ErrExpiredToken struct{}

// Error returns the error message for ErrExpiredToken
func (r *ErrExpiredToken) Error() string {
	return "idtoken: token expired"
}
