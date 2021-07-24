package iam

import (
	"context"
)

type contextKey struct{}

var activeContextKey = contextKey{}

// FromContext returns the Account instance associated with `ctx`
func FromContext(ctx context.Context) (*Account, bool) {
	val := ctx.Value(activeContextKey)
	if o, ok := val.(*Account); ok {
		return o, true
	}
	return nil, false
}

// WithContext returns a copy of parent in which the Account is stored
func WithContext(ctx context.Context, a *Account) context.Context {
	return context.WithValue(ctx, activeContextKey, a)
}
