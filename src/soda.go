package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "password",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "soda",
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	bindRoute(mux, "/", handleDashboard)
	bindRoute(mux, "/database/new", handleDatabaseNew)
	bindRoute(mux, "/databases/{name}", handleDatabaseDetails)
	bindRoute(mux, "/servers/new", handleServerNew)
	bindRoute(mux, "/servers/{name}", handleServerDetails)

	err = http.ListenAndServe(":3030", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func bindRoute(mux *http.ServeMux, path string, handler http.HandlerFunc) {
	mux.HandleFunc(path, Chain(handler, Logging()))
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	dbs, err := getDatabases()
	if err != nil {
		log.Fatal(err)
	}

	servers, err := getServers()
	if err != nil {
		log.Fatal(err)
	}

	renderTemplate(w, "dashboard", Dashboard{
		Databases: dbs,
		Servers:   servers,
	})
}

func handleDatabaseDetails(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	log.Println(name, "details")

	db, err := getDatabaseByName(name)
	if err != nil {
		log.Fatal(err)
	}

	renderTemplate(w, "database-details", Database{
		Name:   db.Name,
		Server: db.Server,
	})
}

func handleDatabaseNew(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		database := Database{
			Name:   r.PostFormValue("name"),
			Server: r.PostFormValue("server"),
		}

		_, err := addDatabase(database)
		if err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, r, "/database/"+database.Name, http.StatusSeeOther)
	}

	renderTemplate(w, "database-new", Database{
		Name:   "",
		Server: "",
	})
}

func handleServerDetails(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	log.Println(name, "details")

	server, err := getServerByName(name)
	if err != nil {
		log.Fatal(server)
	}

	renderTemplate(w, "server-details", Server{
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

func handleServerNew(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		log.Println("Adding server")
		server := Server{
			Name:      r.PostFormValue("name"),
			IpAddress: r.PostFormValue("ipAddress"),
			Port:      r.PostFormValue("port"),
			Type:      r.PostFormValue("type"),
			Username:  r.PostFormValue("username"),
			Password:  r.PostFormValue("password"),
		}

		log.Println("Saving", server)

		_, err := addServer(server)
		if err != nil {
			log.Fatal("Adding server:", err)
		}

		log.Println("Done adding server")
		http.Redirect(w, r, "/servers/"+server.Name, http.StatusSeeOther)
	}

	renderTemplate(w, "server-new", Server{
		Name:      "",
		IpAddress: "",
		Username:  "",
		Password:  "",
		Type:      "mysql",
		Port:      "3306",
	})
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	tmpls := template.Must(template.ParseFiles("views/soda.tmpl", "views/"+name+".tmpl"))
	err := tmpls.ExecuteTemplate(w, "soda.tmpl", data)

	if err != nil {
		log.Println(err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func addDatabase(database Database) (int64, error) {
	res, err := db.Exec("INSERT INTO soda_databases (name, server_name) VALUES (?, ?)", database.Name, database.Server)
	if err != nil {
		return 0, fmt.Errorf("addDatabase: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addDatabase: %v", err)
	}

	return id, nil
}

func getDatabaseByName(name string) (Database, error) {
	row := db.QueryRow("SELECT name, server_name FROM soda_databases WHERE name=?", name)

	var d Database
	err := row.Scan(&d.Name, &d.Server)
	if err != nil {
		return d, err
	}

	return d, nil
}

func getDatabases() ([]Database, error) {
	rows, err := db.Query("SELECT name, server_name FROM soda_databases")
	if err != nil {
		return nil, err
	}

	var databases []Database
	for rows.Next() {
		var d Database
		if err := rows.Scan(&d.Name, &d.Server); err != nil {
			return databases, err
		}
		databases = append(databases, d)
	}

	if err = rows.Err(); err != nil {
		return databases, err
	}

	return databases, nil
}

func addServer(server Server) (int64, error) {
	res, err := db.Exec("INSERT INTO soda_servers (name, ip_address, port, type, username, password) VALUES (?, ?, ?, ?, ?, ?)", server.Name, server.IpAddress, server.Port, server.Type, server.Username, server.Password)
	if err != nil {
		return 0, fmt.Errorf("addServer: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addServer: %v", err)
	}

	return id, nil
}

func getServers() ([]Server, error) {
	rows, err := db.Query("SELECT name, type, ip_address, port FROM soda_servers")
	if err != nil {
		return nil, err
	}

	var servers []Server
	for rows.Next() {
		var s Server
		if err := rows.Scan(&s.Name, &s.Type, &s.IpAddress, &s.Port); err != nil {
			return servers, err
		}
		servers = append(servers, s)
	}

	if err = rows.Err(); err != nil {
		return servers, err
	}

	return servers, nil
}

func getServerByName(name string) (Server, error) {
	row := db.QueryRow("SELECT name, ip_address FROM soda_servers WHERE name=?", name)

	var server Server
	err := row.Scan(&server.Name, &server.IpAddress)
	if err != nil {
		return server, err
	}

	return server, nil
}

type Dashboard struct {
	Databases []Database
	Servers   []Server
}

type Server struct {
	Name      string
	Type      string
	Databases int
	IpAddress string
	Port      string
	Status    string
	Username  string
	Password  string
}

type Database struct {
	Name   string
	Server string
}

type Middleware func(http.HandlerFunc) http.HandlerFunc

// Logging logs all requests with its path and the time it took to process
func Logging() Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			log.Println(r.URL.Path, time.Since(start))
			defer func() { log.Println(r.URL.Path, time.Since(start)) }()

			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}

// Chain applies middlewares to a http.HandlerFunc
func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}
