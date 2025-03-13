package api

import (
	"net/http"

	"github.com/gregbrant2/soda/internal/api/handlers"
	"github.com/gregbrant2/soda/internal/domain/dataaccess"
	"github.com/gregbrant2/soda/internal/plumbing/routing"
)

func RegisterRoutes(
	uow dataaccess.UnitOfWork,
	mux *http.ServeMux) {

	routing.BindRoute(mux, "/api/servers", handlers.HandleServers(uow))
	routing.BindRoute(mux, "/api/server/{id}", handlers.HandleServerDetails(uow))
	routing.BindRoute(mux, "/api/server", handlers.HandleServerNew(uow))
	routing.BindRoute(mux, "/api/databases", handlers.HandleDatabases(uow))
	routing.BindRoute(mux, "/api/database/{id}", handlers.HandleDatabaseDetails(uow))
	routing.BindRoute(mux, "/api/database", handlers.HandleDatabaseNew(uow))
}

func emptyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}
