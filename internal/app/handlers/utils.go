package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	p, err := filepath.Abs(".")
	log.Printf("PWD: %s\n", p)
	p, err = filepath.Abs("./web/template")
	log.Printf("web/templates: %s", p)
	info, err := os.Stat("./web/template")
	log.Println(info.Name(), info.IsDir(), err)
	tmpls := template.Must(template.ParseFiles("./web/template/soda.tmpl", "./web/template/"+name+".tmpl"))
	err = tmpls.ExecuteTemplate(w, "soda.tmpl", data)

	if err != nil {
		log.Println(err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "custom 404")
	}
}
