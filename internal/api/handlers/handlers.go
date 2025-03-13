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
	"golang.org/x/exp/slog"
)

func HandleServers(uow dataaccess.UnitOfWork) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		servers, err := uow.Servers.GetServers()
		if err != nil {
			slog.Error("Error getting servers", utils.ErrAttr(err))
			handleError(w, http.StatusInternalServerError, err.Error(), nil)
		}

		if servers == nil {
			notFound(w)
		}

		returnJson(w, mapping.MapServers(servers))
	}
}

func HandleServerDetails(uow dataaccess.UnitOfWork) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Getting server details")
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
		if err != nil {
			slog.Error("Error getting Id", "id", id, utils.ErrAttr(err))
			handleError(w, http.StatusBadRequest, err.Error(), nil)
			return
		}

		slog.Info("Getting server details for", "id", id)

		server, err := uow.Servers.GetServerById(id)
		if err != nil {
			slog.Error("Error getting server from repository", "id", id, err)
			handleError(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		if server == nil {
			notFound(w)
			return
		}

		returnJson(w, mapping.MapServer(*server))
	}
}

func HandleServerNew(uow dataaccess.UnitOfWork) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Attempting to create new server")
		var payload dtos.NewServer
		err := json.NewDecoder(r.Body).Decode(&payload)

		if err != nil {
			slog.Error("Error decoding JSON request", utils.ErrAttr(err))
			handleError(w, http.StatusBadRequest, err.Error(), nil)
			return
		}

		server := mapping.MapNewServer(payload)
		valid, errors := validation.ValidateServerNew(uow, server)
		if !valid {
			slog.Debug("Validation errors were encountered", "errors", errors)
			handleError(w, http.StatusBadRequest, "Invalid payload", errors)
			return
		}

		serverId, err := uow.Servers.AddServer(server)
		if err != nil {
			handleError(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		result := mapping.MapServer(server)
		result.Id = serverId

		returnJson(w, result)
	}
}

func HandleDatabases(uow dataaccess.UnitOfWork) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Getting databases")
		databases, err := uow.DBs.GetDatabases()
		if err != nil {
			slog.Error("Error getting databases", utils.ErrAttr(err))
			handleError(w, http.StatusInternalServerError, err.Error(), nil)
		}

		if databases == nil {
			slog.Debug("Datanases was nil")
			notFound(w)
		}

		returnJson(w, mapping.MapDatabases(databases))
	}
}

func HandleDatabaseDetails(uow dataaccess.UnitOfWork) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Getting database details")
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
		slog.Info("Attempting to get database", "id", id)
		if err != nil {
			slog.Error("Error getting database", utils.ErrAttr(err))
			handleError(w, http.StatusBadRequest, err.Error(), nil)
			return
		}

		database, err := uow.DBs.GetDatabaseById(id)
		if err != nil {
			slog.Error("Error getting database", "id", id, utils.ErrAttr(err))
			handleError(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		if database == nil {
			notFound(w)
			return
		}

		returnJson(w, mapping.MapDatabase(*database))
	}
}

func HandleDatabaseNew(uow dataaccess.UnitOfWork) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Attempting to create new database")
		var payload dtos.NewDatabase
		err := json.NewDecoder(r.Body).Decode(&payload)

		if err != nil {
			slog.Error("Error decoding JSON request", utils.ErrAttr(err))
			handleError(w, http.StatusBadRequest, err.Error(), nil)
			return
		}

		database := mapping.MapNewDatabase(payload)
		valid, errors := validation.ValidateDatabaseNew(uow, database)
		if !valid {
			slog.Debug("Validation errors were encountered", "errors", errors)
			handleError(w, http.StatusBadRequest, "Invalid payload", errors)
			return
		}

		dbId, err := uow.DBs.AddDatabase(database)
		if err != nil {
			handleError(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		result := mapping.MapDatabase(database)
		result.Id = dbId

		returnJson(w, result)
	}
}

func notFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func returnJson(w http.ResponseWriter, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload)
}

func handleError(w http.ResponseWriter, status int, message string, fields map[string]string) {
	err := dtos.ApiError{Message: message, Fields: fields}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(err)
}
