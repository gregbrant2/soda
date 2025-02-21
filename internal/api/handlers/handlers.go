package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gregbrant2/soda/internal/api/mapping"
	"github.com/gregbrant2/soda/internal/domain/dataaccess"
)

func HandleServers(dbr dataaccess.DatabaseRepository, sr dataaccess.ServerRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		servers, err := sr.GetServers()
		if err != nil {
			log.Println(err)
			handleError(w, http.StatusInternalServerError, err.Error())
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
			handleError(w, http.StatusBadRequest, err.Error())
			return
		}

		log.Println(id, "details")

		server, err := sr.GetServerById(id)
		if err != nil {
			log.Println(err)
			handleError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if server == nil {
			notFound(w)
			return
		}

		returnJson(w, mapping.MapServer(*server))
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

func handleError(w http.ResponseWriter, status int, content string) {
	w.WriteHeader(status)
	fmt.Fprint(w, content)
}
