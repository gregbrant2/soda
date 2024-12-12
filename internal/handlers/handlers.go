package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/elliotchance/pie/v2"

	"github.com/gregbrant2/soda/internal/clients"
	"github.com/gregbrant2/soda/internal/dataaccess"
	"github.com/gregbrant2/soda/internal/entities"
	"github.com/gregbrant2/soda/internal/viewmodels"
)

func HandleDashboard(w http.ResponseWriter, r *http.Request) {
	dbs, err := dataaccess.GetDatabases()
	if err != nil {
		log.Fatal(err)
	}

	servers, err := dataaccess.GetServers()
	if err != nil {
		log.Fatal(err)
	}

	renderTemplate(w, "dashboard", viewmodels.Dashboard{
		Databases: dbs,
		Servers:   servers,
	})
}

func HandleDatabaseDetails(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		log.Println(err)
		errorHandler(w, r, http.StatusBadRequest)
		return
	}

	log.Println(id, "details")

	db, err := dataaccess.GetDatabaseById(id)
	if err != nil {
		log.Fatal(err)
	}

	renderTemplate(w, "database-details", db)
}

func HandleDatabaseNew(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		database := entities.Database{
			Name:   r.PostFormValue("name"),
			Server: r.PostFormValue("server"),
		}

		server, err := dataaccess.GetServerByName(database.Server)
		if err != nil {
			log.Fatal(err)
		}

		id, err := dataaccess.AddDatabase(database)
		if err != nil {
			log.Fatal(err)
		}

		c := clients.MySqlClient{}
		err = c.CreateDatabase(server, database.Name)
		if err != nil {
			log.Fatal(err)
		}

		err = c.CreateUser(server, database.Name, database.Name, database.Name)
		if err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, r, "/databases/"+strconv.FormatInt(int64(id), 10), http.StatusSeeOther)
	}

	servers, err := dataaccess.GetServers()
	if err != nil {
		log.Fatal(err)
	}

	serverNames := pie.Map(
		servers,
		func(e entities.Server) string {
			return e.Name
		},
	)

	renderTemplate(w, "database-new", viewmodels.NewDatabase{
		Database: entities.Database{
			Name:   "",
			Server: "",
		},
		ServerNames: serverNames,
	})
}

func HandleServerDetails(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		log.Println(err)
		errorHandler(w, r, http.StatusBadRequest)
		return
	}

	log.Println(id, "details")

	server, err := dataaccess.GetServerById(id)
	if err != nil {
		log.Fatal(err)
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

func HandleServerNew(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		log.Println("Adding server")
		server := entities.Server{
			Name:      r.PostFormValue("name"),
			IpAddress: r.PostFormValue("ipAddress"),
			Port:      r.PostFormValue("port"),
			Type:      r.PostFormValue("type"),
			Username:  r.PostFormValue("username"),
			Password:  r.PostFormValue("password"),
		}

		log.Println("Saving", server)

		id, err := dataaccess.AddServer(server)
		if err != nil {
			log.Fatal("Adding server:", err)
		}

		log.Println("Done adding server")
		http.Redirect(w, r, "/servers/"+strconv.FormatInt(int64(id), 10), http.StatusSeeOther)
	}

	renderTemplate(w, "server-new", entities.Server{
		Name:      "",
		IpAddress: "",
		Username:  "",
		Password:  "",
		Type:      "mysql",
		Port:      "3306",
	})
}
