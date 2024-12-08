package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", handleDashboard)
	mux.HandleFunc("/database/cola", handleDatabaseDetails)
	mux.HandleFunc("/database/new", handleDatabaseNew)
	mux.HandleFunc("/server/db-01", handleServerDetails)

	err := http.ListenAndServe(":3030", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "dashboard", Dashboard{
		Databases: []Database{
			{Name: "cola", Server: "db-01"},
			{Name: "lemonade", Server: "db-02"},
		},
	})
}

func handleDatabaseDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
	}
	renderTemplate(w, "database-details", Database{
		Name:   "cola",
		Server: "db-01",
	})
}

func handleDatabaseNew(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		database := Database{
			Name:   r.PostFormValue("name"),
			Server: r.PostFormValue("server"),
		}

	}

	renderTemplate(w, "database-new", Database{
		Name:   "",
		Server: "",
	})
}

func handleServerDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
	}
	renderTemplate(w, "server-details", Server{
		Name:      "db-01",
		IpAddress: "10.0.0.36",
		Status:    "OK",
		Databases: 2,
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
