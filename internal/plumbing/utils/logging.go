package utils

import (
	"log/slog"
	"os"

	slogmulti "github.com/samber/slog-multi"
)

func InitLogging() {
	stderr := os.Stderr
	fileLog := OpenCreate("soda.log")

	logger := slog.New(
		slogmulti.Fanout(
			slog.NewJSONHandler(fileLog, &slog.HandlerOptions{}),
			slog.NewTextHandler(stderr, &slog.HandlerOptions{}),
			// ...
		),
	)
	slog.SetDefault(logger)
}

func ErrAttr(err error) slog.Attr {
	return slog.Any("error", err)
}
