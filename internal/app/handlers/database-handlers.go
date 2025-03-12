package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/elliotchance/pie/v2"

	"github.com/gregbrant2/soda/internal/app/viewmodels"
	"github.com/gregbrant2/soda/internal/domain/dataaccess"
	"github.com/gregbrant2/soda/internal/domain/entities"
	"github.com/gregbrant2/soda/internal/domain/validation"
	"github.com/gregbrant2/soda/internal/plumbing/clients"
	"github.com/gregbrant2/soda/internal/plumbing/utils"
)

func HandleDatabaseDetails(dbr dataaccess.DatabaseRepository, sr dataaccess.ServerRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Getting database details")
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
		if err != nil {
			slog.Error("Error parsing Id from query", utils.ErrAttr(err))
			errorHandler(w, r, http.StatusBadRequest)
			return
		}

		slog.Info("Database details", "id", id)
		db, err := dbr.GetDatabaseById(id)
		if err != nil {
			utils.Fatal("Error getting db by id", err, "id", id)
		}

		server, err := sr.GetServerByName(db.Server)
		if err != nil {
			utils.Fatal("Error getting server for db", err, "db", db)
		}

		renderTemplate(w, "database-details", viewmodels.DatabaseDetails{
			Database: *db,
			Server:   *server,
		})
	}
}

func HandleDatabaseNew(dbr dataaccess.DatabaseRepository, sr dataaccess.ServerRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("New database")
		servers, err := sr.GetServers()
		if err != nil {
			utils.Fatal("Error getting servers", err)
		}

		var selectedServer entities.Server
		var selectedServerId int64 = -1
		selectedServerQuery := r.URL.Query().Get("serverId")
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

		if r.Method == http.MethodPost {
			database := entities.Database{
				Name:   r.PostFormValue("name"),
				Server: r.PostFormValue("server"),
			}

			slog.Debug("Adding database")
			valid, errors := validation.ValidateDatabaseNew(dbr, sr, database)
			if !valid {
				vm.Errors = errors
				vm.Database = database
				slog.Debug("Returning validation errors", "errors", errors)
				handleDatabaseNewView(w, r, servers, vm)
				return
			}

			server, err := sr.GetServerByName(database.Server)
			if err != nil {
				utils.Fatal("Error getting target server", err)
			}

			id, err := dbr.AddDatabase(database)
			if err != nil {
				utils.Fatal("Error adding database", err)
			}

			c, err := clients.CreateServer(*server)
			if err != nil {
				utils.Fatal("", err)
			}

			err = c.CreateDatabase(*server, database.Name)
			if err != nil {
				utils.Fatal("Error creating database on target server", err)
			}

			// TODO: Implement password
			err = c.CreateUser(*server, database.Name, database.Name, database.Name)
			if err != nil {
				utils.Fatal("Error creating user on target server", err)
			}

			http.Redirect(w, r, fmt.Sprintf("/databases/%d", id), http.StatusSeeOther)
		}

		handleDatabaseNewView(w, r, servers, vm)
	}
}

func handleDatabaseNewView(w http.ResponseWriter, r *http.Request, servers []entities.Server, vm viewmodels.NewDatabase) {
	serverNames := pie.Map(
		servers,
		func(e entities.Server) string {
			return e.Name
		},
	)

	vm.ServerNames = serverNames
	renderTemplate(w, "database-new", vm)
}
