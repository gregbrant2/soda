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
		log.Fatal("Error getting db by id", id, err)
	}

	server, err := dataaccess.GetServerByName(db.Server)
	if err != nil {
		log.Fatal("Error getting server for db", db, err)
	}

	renderTemplate(w, "database-details", viewmodels.DatabaseDetails{
		Database: db,
		Server:   server,
	})
}

func HandleDatabaseNew(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		database := entities.Database{
			Name:   r.PostFormValue("name"),
			Server: r.PostFormValue("server"),
		}

		log.Println("Adding database", database)

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

	var selectedServer entities.Server
	var selectedServerId int64 = -1
	selectedServerQuery := r.URL.Query().Get("serverId")
	if len(selectedServerQuery) > 0 {
		selectedServerId, err = strconv.ParseInt(selectedServerQuery, 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		selectedServer = servers[pie.FindFirstUsing(servers, func(s entities.Server) bool { return s.Id == selectedServerId })]
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
			Server: selectedServer.Name,
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
