package main

import (
	"log"
	"net/http"

	"github.com/gregbrant2/soda/internal/dataaccess"
	"github.com/gregbrant2/soda/internal/handlers"
	"github.com/gregbrant2/soda/internal/middlewares"
)

func main() {
	log.Println(`
   _____           _       
  / ____|         | |      
 | (___   ___   __| | __ _ 
  \___ \ / _ \ / _` + "` |/ _`" + ` |
  ____) | (_) | (_| | (_| |
 |_____/ \___/ \__,_|\__,_|                           
                           `)

	db := dataaccess.Initialize()
	defer db.Close()

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("web/static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	bindRoute(mux, "/", handlers.HandleDashboard)
	bindRoute(mux, "/database/new", handlers.HandleDatabaseNew)
	bindRoute(mux, "/databases/{id}", handlers.HandleDatabaseDetails)
	bindRoute(mux, "/servers/new", handlers.HandleServerNew)
	bindRoute(mux, "/servers/{id}", handlers.HandleServerDetails)

	err := http.ListenAndServe(":3030", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func bindRoute(mux *http.ServeMux, path string, handler http.HandlerFunc) {
	mux.HandleFunc(path, middlewares.Chain(handler, middlewares.Logging()))
}
