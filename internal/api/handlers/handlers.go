package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gregbrant2/soda/internal/api/dtos"
	"github.com/gregbrant2/soda/internal/api/mapping"
	"github.com/gregbrant2/soda/internal/domain/dataaccess"
	"github.com/gregbrant2/soda/internal/domain/validation"
)

func HandleServers(dbr dataaccess.DatabaseRepository, sr dataaccess.ServerRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		servers, err := sr.GetServers()
		if err != nil {
			log.Println(err)
			handleError(w, http.StatusInternalServerError, err.Error(), nil)
		}

		if servers == nil {
			notFound(w)
		}

		returnJson(w, mapping.MapServers(servers))
	}
}

func HandleServerDetails(sr dataaccess.ServerRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
		if err != nil {
			log.Println(err)
			handleError(w, http.StatusBadRequest, err.Error(), nil)
			return
		}

		log.Println(id, "details")

		server, err := sr.GetServerById(id)
		if err != nil {
			log.Println(err)
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

func HandleServerNew(sr dataaccess.ServerRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Attempting to create new server")
		var payload dtos.NewServer
		err := json.NewDecoder(r.Body).Decode(&payload)

		// id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
		if err != nil {
			log.Println("Error decoding JSON request")
			log.Println(err)
			handleError(w, http.StatusBadRequest, err.Error(), nil)
			return
		}

		server := mapping.MapNewServer(payload)
		valid, errors := validation.ValidateServerNew(sr, server)
		if !valid {
			log.Println("Validation errors were encountered", errors)
			handleError(w, http.StatusBadRequest, "Invalid payload", errors)
			return
		}

		serverId, err := sr.AddServer(server)
		if err != nil {
			handleError(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		result := mapping.MapServer(server)
		result.Id = serverId

		returnJson(w, result)
	}
}

func HandleDatabases(dbr dataaccess.DatabaseRepository, sr dataaccess.ServerRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		databases, err := dbr.GetDatabases()
		if err != nil {
			log.Println(err)
			handleError(w, http.StatusInternalServerError, err.Error(), nil)
		}

		if databases == nil {
			notFound(w)
		}

		returnJson(w, mapping.MapDatabases(databases))
	}
}

func HandleDatabaseDetails(sr dataaccess.DatabaseRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
		log.Printf("Attempting to get database %d\n", id)
		if err != nil {
			log.Println(err)
			handleError(w, http.StatusBadRequest, err.Error(), nil)
			return
		}

		database, err := sr.GetDatabaseById(id)
		if err != nil {
			log.Println(err)
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

func HandleDatabaseNew(dbr dataaccess.DatabaseRepository, sr dataaccess.ServerRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Attempting to create new database")
		var payload dtos.NewDatabase
		err := json.NewDecoder(r.Body).Decode(&payload)

		if err != nil {
			log.Println("Error decoding JSON request")
			log.Println(err)
			handleError(w, http.StatusBadRequest, err.Error(), nil)
			return
		}

		database := mapping.MapNewDatabase(payload)
		valid, errors := validation.ValidateDatabaseNew(dbr, sr, database)
		if !valid {
			log.Println("Validation errors were encountered", errors)
			handleError(w, http.StatusBadRequest, "Invalid payload", errors)
			return
		}

		dbId, err := dbr.AddDatabase(database)
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
