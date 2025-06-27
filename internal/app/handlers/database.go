package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/elliotchance/pie/v2"
	"github.com/labstack/echo/v4"

	"github.com/gregbrant2/soda/internal/app/viewmodels"
	"github.com/gregbrant2/soda/internal/domain/dataaccess"
	"github.com/gregbrant2/soda/internal/domain/entities"
	"github.com/gregbrant2/soda/internal/domain/validation"
	"github.com/gregbrant2/soda/internal/plumbing/clients"
	"github.com/gregbrant2/soda/internal/plumbing/utils"
)

func HandleDatabaseDetails(uow dataaccess.UnitOfWork) echo.HandlerFunc {
	return func(c echo.Context) error {
		slog.Debug("Getting database details")
		id, err := strconv.ParseInt(c.Param("id"), 10, 32)
		if err != nil {
			slog.Error("Error parsing Id from query", utils.ErrAttr(err))
			return errorHandler(c, http.StatusBadRequest)
		}

		slog.Info("Database details", "id", id)
		db, err := uow.DBs.GetDatabaseById(id)
		if err != nil {
			utils.Fatal("Error getting db by id", err, "id", id)
		}

		server, err := uow.Servers.GetServerByName(db.Server)
		if err != nil {
			utils.Fatal("Error getting server for db", err, "db", db)
		}

		return renderTemplate(c, "database-details", viewmodels.DatabaseDetails{
			Database: *db,
			Server:   *server,
		})
	}
}

func HandleDatabaseNew(uow dataaccess.UnitOfWork) echo.HandlerFunc {
	return func(c echo.Context) error {
		slog.Debug("New database")
		servers, err := uow.Servers.GetServers()
		if err != nil {
			utils.Fatal("Error getting servers", err)
		}

		var selectedServer entities.Server
		var selectedServerId int64 = -1
		selectedServerQuery := c.QueryParam("serverId")
		if len(selectedServerQuery) > 0 {
			selectedServerId, err = strconv.ParseInt(selectedServerQuery, 10, 64)
			if err != nil {
				utils.Fatal("Error parsing selected server", err, "query", selectedServerQuery)
			}

			selectedServer = servers[pie.FindFirstUsing(servers, func(s entities.Server) bool { return s.Id == selectedServerId })]
		}

		vm := viewmodels.NewDatabase{
			Database: entities.Database{
				Server: selectedServer.Name,
			},
			Errors: nil,
		}

		if c.Request().Method == http.MethodPost {
			database := entities.Database{
				Name:   c.FormValue("name"),
				Server: c.FormValue("server"),
			}

			slog.Debug("Adding database")
			valid, errors := validation.ValidateDatabaseNew(uow, database)
			if !valid {
				vm.Errors = errors
				vm.Database = database
				slog.Debug("Returning validation errors", "errors", errors)
				return handleDatabaseNewView(c, servers, vm)
			}

			server, err := uow.Servers.GetServerByName(database.Server)
			if err != nil {
				utils.Fatal("Error getting target server", err)
			}

			id, err := uow.DBs.AddDatabase(database)
			if err != nil {
				utils.Fatal("Error adding database", err)
			}

			s, err := clients.CreateServer(*server)
			if err != nil {
				utils.Fatal("", err)
			}

			err = s.CreateDatabase(*server, database.Name)
			if err != nil {
				utils.Fatal("Error creating database on target server", err)
			}

			// TODO: Implement password
			err = s.CreateUser(*server, database.Name, database.Name, database.Name)
			if err != nil {
				utils.Fatal("Error creating user on target server", err)
			}

			c.Redirect(http.StatusSeeOther, fmt.Sprintf("/databases/%d", id))
		}

		return handleDatabaseNewView(c, servers, vm)
	}
}

func handleDatabaseNewView(c echo.Context, servers []entities.Server, vm viewmodels.NewDatabase) error {
	serverNames := pie.Map(
		servers,
		func(e entities.Server) string {
			return e.Name
		},
	)

	vm.ServerNames = serverNames
	return renderTemplate(c, "database-new", vm)
}
