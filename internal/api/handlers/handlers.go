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
