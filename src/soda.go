package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", handleDashboard)
	mux.HandleFunc("/database/cola", handleDatabaseDetails)
	mux.HandleFunc("/database/new", handleDatabaseDetails)
	mux.HandleFunc("/server/db-01", handleServerDetails)

	err := http.ListenAndServe(":3030", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "dashboard")
}

func handleDatabaseDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
	}
	renderTemplate(w, "database-details")
}

func handleServerDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
	}
	renderTemplate(w, "server-details")
}

func renderTemplate(w http.ResponseWriter, name string) {
	tmpls := template.Must(template.ParseFiles("views/soda.gohtml", "views/"+name+".gohtml"))
	err := tmpls.ExecuteTemplate(w, "soda.gohtml", nil)

	if err != nil {
		log.Println(err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
