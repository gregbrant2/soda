package app

import (
	"github.com/gregbrant2/soda/internal/app/handlers"
	"github.com/gregbrant2/soda/internal/domain/dataaccess"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(uow dataaccess.UnitOfWork, e *echo.Echo, group *echo.Group) {

	InitRendering(e)

	group.Static("/static", "web/static")

	group.GET("/", handlers.HandleDashboard(uow))

	group.GET("/database/new", handlers.HandleDatabaseNew(uow))
	group.POST("/database/new", handlers.HandleDatabaseNew(uow))
	group.GET("/databases/:id", handlers.HandleDatabaseDetails(uow))

	group.GET("/servers/new", handlers.HandleServerNew(uow))
	group.POST("/servers/new", handlers.HandleServerNew(uow))
	group.GET("/servers/:id", handlers.HandleServerDetails(uow))
}
