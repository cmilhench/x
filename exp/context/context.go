package context

import (
	"context"
	"reflect"
)

type User struct {
	Name string
}

// WithValueOf is a helper function to set a value of a specific type.
func WithValueOf[T any](parent context.Context, val *T) context.Context {
	return context.WithValue(parent, reflect.TypeOf((*T)(nil)), val)
}

// ValueOf is a helper function to get a value of a specific type.
func ValueOf[T any](ctx context.Context) (*T, bool) {
	v, ok := ctx.Value(reflect.TypeOf((*T)(nil))).(*T)
	return v, ok
}
