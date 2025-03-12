package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"text/template"

	"github.com/gregbrant2/soda/internal/plumbing/utils"
)

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	tmpls := template.Must(template.ParseFiles("./web/template/soda.tmpl", "./web/template/"+name+".tmpl"))
	err := tmpls.ExecuteTemplate(w, "soda.tmpl", data)
	if err != nil {
		slog.Error("Error executing template", utils.ErrAttr(err), "name", name)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "custom 404")
	}
}
