package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gregbrant2/soda/internal/api/dtos"
	"github.com/gregbrant2/soda/internal/api/mapping"
	"github.com/gregbrant2/soda/internal/domain/dataaccess"
	"github.com/gregbrant2/soda/internal/domain/validation"
	"github.com/gregbrant2/soda/internal/plumbing/utils"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slog"
)

func HandleServers(uow dataaccess.UnitOfWork) echo.HandlerFunc {
	return func(c echo.Context) error {
		servers, err := uow.Servers.GetServers()
		if err != nil {
			slog.Error("Error getting servers", utils.ErrAttr(err))
			return handleError(c, http.StatusInternalServerError, err.Error(), nil)
		}

		if servers == nil {
			return notFound(c)
		}

		return returnJson(c, mapping.MapServers(servers))
	}
}

func HandleServerDetails(uow dataaccess.UnitOfWork) echo.HandlerFunc {
	return func(c echo.Context) error {
		slog.Debug("Getting server details")
		id, err := strconv.ParseInt(c.Param("id"), 10, 32)
		if err != nil {
			slog.Error("Error getting Id", "id", id, utils.ErrAttr(err))
			return handleError(c, http.StatusBadRequest, err.Error(), nil)
		}

		slog.Info("Getting server details for", "id", id)

		server, err := uow.Servers.GetServerById(id)
		if err != nil {
			slog.Error("Error getting server from repository", "id", id, err)
			return handleError(c, http.StatusInternalServerError, err.Error(), nil)
		}

		if server == nil {
			return notFound(c)
		}

		return returnJson(c, mapping.MapServer(*server))
	}
}

func HandleServerNew(uow dataaccess.UnitOfWork) echo.HandlerFunc {
	return func(c echo.Context) error {
		slog.Debug("Attempting to create new server")
		var payload dtos.NewServer
		err := json.NewDecoder(c.Request().Body).Decode(&payload)

		if err != nil {
			slog.Error("Error decoding JSON request", utils.ErrAttr(err))
			return handleError(c, http.StatusBadRequest, err.Error(), nil)
		}

		server := mapping.MapNewServer(payload)
		valid, errors := validation.ValidateServerNew(uow, server)
		if !valid {
			slog.Debug("Validation errors were encountered", "errors", errors)
			return handleError(c, http.StatusBadRequest, "Invalid payload", errors)
		}

		serverId, err := uow.Servers.AddServer(server)
		if err != nil {
			return handleError(c, http.StatusInternalServerError, err.Error(), nil)
		}

		result := mapping.MapServer(server)
		result.Id = serverId

		return returnJson(c, result)
	}
}

func HandleDatabases(uow dataaccess.UnitOfWork) echo.HandlerFunc {
	return func(c echo.Context) error {
		slog.Debug("Getting databases")
		databases, err := uow.DBs.GetDatabases()
		if err != nil {
			slog.Error("Error getting databases", utils.ErrAttr(err))
			return handleError(c, http.StatusInternalServerError, err.Error(), nil)
		}

		if databases == nil {
			slog.Debug("Datanases was nil")
			return notFound(c)
		}

		return returnJson(c, mapping.MapDatabases(databases))
	}
}

func HandleDatabaseDetails(uow dataaccess.UnitOfWork) echo.HandlerFunc {
	return func(c echo.Context) error {
		slog.Debug("Getting database details")
		id, err := strconv.ParseInt(c.Param("id"), 10, 32)
		slog.Info("Attempting to get database", "id", id)
		if err != nil {
			slog.Error("Error getting database", utils.ErrAttr(err))
			return handleError(c, http.StatusBadRequest, err.Error(), nil)
		}

		database, err := uow.DBs.GetDatabaseById(id)
		if err != nil {
			slog.Error("Error getting database", "id", id, utils.ErrAttr(err))
			return handleError(c, http.StatusInternalServerError, err.Error(), nil)
		}

		if database == nil {
			return notFound(c)
		}

		return returnJson(c, mapping.MapDatabase(*database))
	}
}

func HandleDatabaseNew(uow dataaccess.UnitOfWork) echo.HandlerFunc {
	return func(c echo.Context) error {
		slog.Debug("Attempting to create new database")
		var payload dtos.NewDatabase
		err := json.NewDecoder(c.Request().Body).Decode(&payload)

		if err != nil {
			slog.Error("Error decoding JSON request", utils.ErrAttr(err))
			return handleError(c, http.StatusBadRequest, err.Error(), nil)
		}

		database := mapping.MapNewDatabase(payload)
		valid, errors := validation.ValidateDatabaseNew(uow, database)
		if !valid {
			slog.Debug("Validation errors were encountered", "errors", errors)
			return handleError(c, http.StatusBadRequest, "Invalid payload", errors)
		}

		dbId, err := uow.DBs.AddDatabase(database)
		if err != nil {
			return handleError(c, http.StatusInternalServerError, err.Error(), nil)
		}

		result := mapping.MapDatabase(database)
		result.Id = dbId

		return returnJson(c, result)
	}
}

func notFound(c echo.Context) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusNotFound, "Not Found")
}

func returnJson(c echo.Context, payload interface{}) error {
	return c.JSON(http.StatusOK, payload)
}

func handleError(c echo.Context, status int, message string, fields map[string]string) error {
	err := dtos.ApiError{Message: message, Fields: fields}
	return c.JSON(status, err)
}
