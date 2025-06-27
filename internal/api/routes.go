package api

import (
	"net/http"

	"github.com/gregbrant2/soda/internal/api/handlers"
	"github.com/gregbrant2/soda/internal/domain/dataaccess"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(
	uow dataaccess.UnitOfWork,
	group *echo.Group) {

	group.GET("/api/servers", handlers.HandleServers(uow))
	group.GET("/api/server/:id", handlers.HandleServerDetails(uow))
	group.POST("/api/server", handlers.HandleServerNew(uow))
	group.GET("/api/databases", handlers.HandleDatabases(uow))
	group.POST("/api/database/:id", handlers.HandleDatabaseDetails(uow))
	group.GET("/api/database", handlers.HandleDatabaseNew(uow))
}

func emptyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}
