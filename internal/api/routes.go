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
	routing.BindRoute(mux, "/api/server", handlers.HandleServerNew(sr))
	routing.BindRoute(mux, "/api/databases", handlers.HandleDatabases(dbr, sr))
	routing.BindRoute(mux, "/api/database/{id}", handlers.HandleDatabaseDetails(dbr))
	routing.BindRoute(mux, "/api/database", emptyHandler())
}

func emptyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}
