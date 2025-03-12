package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gregbrant2/soda/internal/app/viewmodels"
	"github.com/gregbrant2/soda/internal/domain/dataaccess"
	"github.com/gregbrant2/soda/internal/plumbing/utils"
)

func HandleDashboard(dbr dataaccess.DatabaseRepository, sr dataaccess.ServerRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Handle Dashboard")
		dbs, err := dbr.GetDatabases()
		if err != nil {
			utils.Fatal("Error getting databases", err)
		}

		servers, err := sr.GetServers()
		if err != nil {
			utils.Fatal("Error getting servers", err)
		}

		renderTemplate(w, "dashboard", viewmodels.Dashboard{
			Databases: dbs,
			Servers:   servers,
		})
	}
}
