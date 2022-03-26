package util

import (
	"context"
	"fmt"
	"reflect"
)

type ContextKey[T any] struct {
	Name string
}

func (key ContextKey[T]) Value(ctx context.Context) (T, bool) {
	o, ok := ctx.Value(key).(T)
	return o, ok
}

func (key ContextKey[T]) WithValue(ctx context.Context, val T) context.Context {
	return context.WithValue(ctx, key, val)
}

func (key ContextKey[T]) String() string {
	name := key.Name
	if name != "" {
		name = "@" + name
	}
	return fmt.Sprintf("ContextKey(%s%s)", reflect.TypeOf(new(T)).Elem().String(), name)
}

func (key ContextKey[T]) Get(ctx context.Context) T {
	o, _ := ctx.Value(key).(T)
	return o
}

func (key ContextKey[T]) Exists(ctx context.Context) bool {
	_, ok := ctx.Value(key).(T)
	return ok
}

func (key ContextKey[T]) Must(ctx context.Context) T {
	o, ok := ctx.Value(key).(T)
	if !ok {
		panic(fmt.Errorf("%s not found in context", key.String()))
	}
	return o
}
