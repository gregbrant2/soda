package routing

import (
	"log/slog"
	"net/http"

	"github.com/gregbrant2/soda/internal/plumbing/middlewares"
)

func BindRoute(
	mux *http.ServeMux,
	path string,
	handler http.HandlerFunc) {
	slog.Debug("Binding route", "path", path, "handler", handler)
	mux.HandleFunc(path, middlewares.Chain(handler, middlewares.Logging()))
}
