package handlers

import (
	"fmt"
	"log/slog"

	"github.com/gregbrant2/soda/internal/app/viewmodels"
	"github.com/gregbrant2/soda/internal/domain/dataaccess"
	"github.com/gregbrant2/soda/internal/plumbing/utils"
	"github.com/labstack/echo/v4"
)

func HandleDashboard(uow dataaccess.UnitOfWork) echo.HandlerFunc {
	return func(c echo.Context) error {
		slog.Debug("Handle Dashboard")
		dbs, err := uow.DBs.GetDatabases()
		if err != nil {
			utils.Fatal("Error getting databases", err)
		}

		servers, err := uow.Servers.GetServers()
		if err != nil {
			utils.Fatal("Error getting servers", err)
		}
		d := viewmodels.Dashboard{
			Databases: dbs,
			Servers:   servers,
		}

		fmt.Printf("Dashboard Databases: %+v\n", d.Databases)
		return renderTemplate(c, "dashboard", d)
	}
}
