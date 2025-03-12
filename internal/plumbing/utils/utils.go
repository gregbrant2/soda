package utils

import (
	"log"
	"log/slog"
)

type MapFunc[T any, TMapped any] func(T) TMapped

func Map[T any, TMapped any](items []T, f MapFunc[T, TMapped]) []TMapped {
	mapped := make([]TMapped, len(items))
	for i, item := range items {
		mapped[i] = f(item)
	}

	return mapped
}

func Fatal(msg string, err error, args ...any) {
	newArgs := append([]any{ErrAttr(err)}, args...)
	slog.Error(msg, newArgs...)
	log.Fatal(msg)
}
