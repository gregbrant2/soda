package utils

type MapFunc[T any, TMapped any] func(T) TMapped

func Map[T any, TMapped any](items []T, f MapFunc[T, TMapped]) []TMapped {
	mapped := make([]TMapped, len(items))
	for i, item := range items {
		mapped[i] = f(item)
	}

	return mapped
}
