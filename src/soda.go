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

	err := http.ListenAndServe(":3030", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("views/soda.gohtml", "views/dashboard.gohtml"))
	log.Println(tpl)

	err := tpl.ExecuteTemplate(w, "soda.gohtml", nil)
	log.Println(err)

	if err != nil {
		log.Println(err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
