package app

import (
	"net/http"

	"github.com/gregbrant2/soda/internal/app/handlers"
	"github.com/gregbrant2/soda/internal/domain/dataaccess"
	"github.com/gregbrant2/soda/internal/plumbing/routing"
)

func RegisterRoutes(
	uow dataaccess.UnitOfWork,
	mux *http.ServeMux) {

	fs := http.FileServer(http.Dir("./web/static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	routing.BindRoute(mux, "/", handlers.HandleDashboard(uow))
	routing.BindRoute(mux, "/database/new", handlers.HandleDatabaseNew(uow))
	routing.BindRoute(mux, "/databases/{id}", handlers.HandleDatabaseDetails(uow))
	routing.BindRoute(mux, "/servers/new", handlers.HandleServerNew(uow))
	routing.BindRoute(mux, "/servers/{id}", handlers.HandleServerDetails(uow))
}
