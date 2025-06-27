package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/gregbrant2/soda/internal/api"
	"github.com/gregbrant2/soda/internal/app"
	"github.com/gregbrant2/soda/internal/domain/dataaccess"
	"github.com/gregbrant2/soda/internal/plumbing/utils"
)

func main() {
	banner()

	utils.InitLogging()
	uow, close := dataaccess.NewUow()
	defer close()

	e := echo.New()

	e.Debug = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	uiGroup := e.Group("")
	app.RegisterRoutes(uow, e, uiGroup)

	apiGroup := e.Group("/api")
	// apiGroup.Use(middleware.JWTWithConfig())
	api.RegisterRoutes(uow, apiGroup)

	e.Logger.Fatal(e.Start(":3030"))
}

func banner() {
	log.Println(`
   _____           _       
  / ____|         | |      
 | (___   ___   __| | __ _ 
  \___ \ / _ \ / _` + "` |/ _`" + ` |
  ____) | (_) | (_| | (_| |
 |_____/ \___/ \__,_|\__,_|                           
                           `)
}
