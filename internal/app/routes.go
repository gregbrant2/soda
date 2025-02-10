package app

import (
	"net/http"

	"github.com/gregbrant2/soda/internal/app/handlers"
	"github.com/gregbrant2/soda/internal/domain/dataaccess"
	"github.com/gregbrant2/soda/internal/plumbing/routing"
)

func RegisterRoutes(dbr dataaccess.DatabaseRepository, sr dataaccess.ServerRepository, mux *http.ServeMux) {

	fs := http.FileServer(http.Dir("../../web/static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	routing.BindRoute(mux, "/", handlers.HandleDashboard(dbr, sr))
	routing.BindRoute(mux, "/database/new", handlers.HandleDatabaseNew(dbr, sr))
	routing.BindRoute(mux, "/databases/{id}", handlers.HandleDatabaseDetails(dbr, sr))
	routing.BindRoute(mux, "/servers/new", handlers.HandleServerNew(sr))
	routing.BindRoute(mux, "/servers/{id}", handlers.HandleServerDetails(sr))
}
