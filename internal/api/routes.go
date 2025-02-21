package api

import (
	"net/http"

	"github.com/gregbrant2/soda/internal/api/handlers"
	"github.com/gregbrant2/soda/internal/domain/dataaccess"
	"github.com/gregbrant2/soda/internal/plumbing/routing"
)

func RegisterRoutes(
	dbr dataaccess.DatabaseRepository,
	sr dataaccess.ServerRepository,
	mux *http.ServeMux) {

	routing.BindRoute(mux, "/api/servers", handlers.HandleServers(dbr, sr))
	routing.BindRoute(mux, "/api/server/{id}", handlers.HandleServerDetails(sr))
	routing.BindRoute(mux, "/api/server/new", emptyHandler())
	routing.BindRoute(mux, "/api/databases", emptyHandler())
	routing.BindRoute(mux, "/api/database/{id}", emptyHandler())
	routing.BindRoute(mux, "/api/database/new", emptyHandler())
}

func emptyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}
