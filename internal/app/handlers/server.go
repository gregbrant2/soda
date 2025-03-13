package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gregbrant2/soda/internal/app/viewmodels"
	"github.com/gregbrant2/soda/internal/domain/dataaccess"
	"github.com/gregbrant2/soda/internal/domain/entities"
	"github.com/gregbrant2/soda/internal/domain/validation"
	"github.com/gregbrant2/soda/internal/plumbing/utils"
)

func HandleServerDetails(uow dataaccess.UnitOfWork) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Getting server details")
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
		if err != nil {
			slog.Error("Error getting Id from query", utils.ErrAttr(err))
			errorHandler(w, r, http.StatusBadRequest)
			return
		}

		slog.Info("Server details for", "id", id)
		server, err := uow.Servers.GetServerById(id)
		if err != nil {
			utils.Fatal("Error getting server", err, "id", id)
		}

		renderTemplate(w, "server-details", entities.Server{
			Id:        server.Id,
			Name:      server.Name,
			IpAddress: server.IpAddress,
			Type:      server.Type,
			Port:      server.Port,
			Username:  server.Username,
			Password:  server.Password,
			Status:    "OK",
			Databases: 2,
		})
	}
}

func HandleServerNew(uow dataaccess.UnitOfWork) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Creating new server")

		vm := viewmodels.NewServer{}

		if r.Method == http.MethodPost {
			slog.Info("Adding server")
			server := entities.Server{
				Name:      r.PostFormValue("name"),
				IpAddress: r.PostFormValue("ipAddress"),
				Port:      r.PostFormValue("port"),
				Type:      r.PostFormValue("type"),
				Username:  r.PostFormValue("username"),
				Password:  r.PostFormValue("password"),
			}

			slog.Debug("Saving", "server", server.Name)
			valid, errors := validation.ValidateServerNew(uow, server)
			if !valid {
				vm.Errors = errors
				vm.Server = &server
				handleServerNewView(w, vm)
				return
			}

			id, err := uow.Servers.AddServer(server)
			if err != nil {
				utils.Fatal("Adding server:", err)
			}

			slog.Debug("Done adding server")
			http.Redirect(w, r, "/servers/"+strconv.FormatInt(int64(id), 10), http.StatusSeeOther)
		}

		handleServerNewView(w, vm)
	}
}

func handleServerNewView(w http.ResponseWriter, vm viewmodels.NewServer) {
	if vm.Server == nil {
		vm.Server = &entities.Server{
			Name:      "",
			IpAddress: "",
			Username:  "",
			Password:  "",
			Type:      "mysql",
			Port:      "3306",
		}
	}

	renderTemplate(w, "server-new", vm)
}
