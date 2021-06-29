package middleware

import (
	"context"
	"errors"
)

type key struct {
	name string
}

var (
	errNoUserInContext = errors.New("no user in context")
)

var keyAuth = &key{"auth"}

func SetAuth(ctx context.Context, auth *Auth) context.Context {
	return context.WithValue(ctx, keyAuth, auth)
}

func ForContext(ctx context.Context) (*Auth, error) {
	auth, ok := ctx.Value(keyAuth).(*Auth)
	if !ok {
		return nil, errNoUserInContext
	}
	return auth, nil
}
