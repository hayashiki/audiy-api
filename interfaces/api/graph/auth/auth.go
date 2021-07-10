package auth

import (
	"context"
	"errors"
)

type key struct {
	name string
}

type Auth struct {
	ID   int64
	Email  string
	Name string
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
