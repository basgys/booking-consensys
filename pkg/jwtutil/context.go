package jwtutil

import (
	"context"

	"github.com/form3tech-oss/jwt-go"
)

type contextKey struct{}

var activeContextKey = contextKey{}

// FromContext returns the JWT instance associated with `ctx`
func FromContext(ctx context.Context) (*jwt.Token, bool) {
	val := ctx.Value(activeContextKey)
	if o, ok := val.(*jwt.Token); ok {
		return o, true
	}
	return nil, false
}

// WithContext returns a copy of parent in which the JWT is stored
func WithContext(ctx context.Context, t *jwt.Token) context.Context {
	return context.WithValue(ctx, activeContextKey, t)
}
