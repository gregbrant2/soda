package app

import (
	"fmt"
	"io"

	"github.com/CloudyKit/jet/v6"
	"github.com/labstack/echo/v4"
)

type JetRenderer struct {
	views *jet.Set
}

func (r *JetRenderer) Render(w io.Writer, name string, model interface{}, c echo.Context) error {
	tmpl, err := r.views.GetTemplate(name)
	if err != nil {
		return err
	}

	var vars jet.VarMap
	// If caller passed a VarMap, use it directly
	// if model != nil {
	// 	if vm, ok := model.(jet.VarMap); ok {
	// 		vars = vm
	// 	} else {
	// 		// Fallback: inject entire data as "data" variable
	// 		vars = make(jet.VarMap)
	// 		vars.Set("Model", model)
	// 	}
	// } else {
	// 	vars = make(jet.VarMap)
	// }

	fmt.Printf("Rendering with Vars: %+v\n", vars)
	return tmpl.Execute(w, vars, model)
}

func InitRendering(e *echo.Echo) {

	views := jet.NewSet(
		jet.NewOSFileSystemLoader("web/template"),
		jet.InDevelopmentMode(),
	)

	renderer := &JetRenderer{
		views: views,
	}

	e.Renderer = renderer
}
