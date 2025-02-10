package routing

import (
	"net/http"

	"github.com/gregbrant2/soda/internal/plumbing/middlewares"
)

func BindRoute(mux *http.ServeMux, path string, handler http.HandlerFunc) {
	mux.HandleFunc(path, middlewares.Chain(handler, middlewares.Logging()))
}
