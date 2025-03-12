package main

import (
	"log"
	"net/http"

	"github.com/gregbrant2/soda/internal/api"
	"github.com/gregbrant2/soda/internal/app"
	"github.com/gregbrant2/soda/internal/domain/dataaccess"
	"github.com/gregbrant2/soda/internal/plumbing/utils"
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

	utils.InitLogging()

	db := dataaccess.Initialize()
	defer db.Close()

	dbr := dataaccess.NewMySqlDatabaseRepository(db)
	sr := dataaccess.NewMySqlServerRepository(db)

	mux := http.NewServeMux()
	app.RegisterRoutes(dbr, sr, mux)
	api.RegisterRoutes(dbr, sr, mux)

	err := http.ListenAndServe(":3030", mux)
	if err != nil {
		log.Fatal(err)
	}
}
