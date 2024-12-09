package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

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

	mux.HandleFunc("/", handleDashboard)
	mux.HandleFunc("/database/new", handleDatabaseNew)
	mux.HandleFunc("/databases/{name}", handleDatabaseDetails)
	mux.HandleFunc("/server/new", handleServerNew)
	mux.HandleFunc("/servers/{name}", handleServerDetails)

	err = http.ListenAndServe(":3030", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	dbs, err := getDatabases()
	if err != nil {
		log.Fatal(err)
	}

	renderTemplate(w, "dashboard", Dashboard{
		Databases: dbs,
	})
}
func handleDatabaseDetails(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
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
	server, err := getServerByName(name)
	if err != nil {
		log.Fatal(server)
	}

	renderTemplate(w, "server-details", Server{
		Name:      server.Name,
		IpAddress: server.IpAddress,
		Status:    "OK",
		Databases: 2,
	})
}

func handleServerNew(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		server := Server{
			Name:      r.PostFormValue("name"),
			IpAddress: r.PostFormValue("ipAddress"),
		}

		_, err := addServer(server)
		if err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, r, "/server/"+server.Name, http.StatusSeeOther)
	}

	renderTemplate(w, "server-new", Server{
		Name:      "",
		IpAddress: "",
	})
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	tmpls := template.Must(template.ParseFiles("views/soda.gohtml", "views/"+name+".gohtml"))
	err := tmpls.ExecuteTemplate(w, "soda.gohtml", data)

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
	res, err := db.Exec("INSERT INTO soda_servers (name, ip_address) VALUES (?, ?)", server.Name, server.IpAddress)
	if err != nil {
		return 0, fmt.Errorf("addServer: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addServer: %v", err)
	}

	return id, nil
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
}

type Server struct {
	Name      string
	Databases int
	IpAddress string
	Status    string
}

type Database struct {
	Name   string
	Server string
}
