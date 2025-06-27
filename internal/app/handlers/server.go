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
	"github.com/labstack/echo/v4"
)

func HandleServerDetails(uow dataaccess.UnitOfWork) echo.HandlerFunc {
	return func(c echo.Context) error {
		slog.Debug("Getting server details")
		id, err := strconv.ParseInt(c.Param("id"), 10, 32)
		if err != nil {
			slog.Error("Error getting Id from query", utils.ErrAttr(err))
			return errorHandler(c, http.StatusBadRequest)
		}

		slog.Info("Server details for", "id", id)
		server, err := uow.Servers.GetServerById(id)
		if err != nil {
			utils.Fatal("Error getting server", err, "id", id)
		}

		return renderTemplate(c, "server-details", entities.Server{
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

func HandleServerNew(uow dataaccess.UnitOfWork) echo.HandlerFunc {
	return func(c echo.Context) error {
		slog.Debug("Creating new server")

		vm := viewmodels.NewServer{}

		if c.Request().Method == http.MethodPost {
			slog.Info("Adding server")
			server := entities.Server{
				Name:      c.FormValue("name"),
				IpAddress: c.FormValue("ipAddress"),
				Port:      c.FormValue("port"),
				Type:      c.FormValue("type"),
				Username:  c.FormValue("username"),
				Password:  c.FormValue("password"),
			}

			slog.Debug("Saving", "server", server.Name)
			valid, errors := validation.ValidateServerNew(uow, server)
			if !valid {
				vm.Errors = errors
				vm.Server = &server
				return handleServerNewView(c, vm)
			}

			id, err := uow.Servers.AddServer(server)
			if err != nil {
				utils.Fatal("Adding server:", err)
			}

			slog.Debug("Done adding server")
			c.Redirect(http.StatusSeeOther, "/servers/"+strconv.FormatInt(int64(id), 10))
		}

		return handleServerNewView(c, vm)
	}
}

func handleServerNewView(c echo.Context, vm viewmodels.NewServer) error {
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

	return renderTemplate(c, "server-new", vm)
}
