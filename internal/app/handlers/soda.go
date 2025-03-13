package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gregbrant2/soda/internal/app/viewmodels"
	"github.com/gregbrant2/soda/internal/domain/dataaccess"
	"github.com/gregbrant2/soda/internal/plumbing/utils"
)

func HandleDashboard(uow dataaccess.UnitOfWork) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Handle Dashboard")
		dbs, err := uow.DBs.GetDatabases()
		if err != nil {
			utils.Fatal("Error getting databases", err)
		}

		servers, err := uow.Servers.GetServers()
		if err != nil {
			utils.Fatal("Error getting servers", err)
		}

		renderTemplate(w, "dashboard", viewmodels.Dashboard{
			Databases: dbs,
			Servers:   servers,
		})
	}
}
