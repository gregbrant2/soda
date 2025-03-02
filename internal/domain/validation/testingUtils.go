package validation

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gregbrant2/soda/internal/domain/entities"
)

func ValidServer() *entities.Server {
	server := &entities.Server{}
	gofakeit.Struct(&server)
	server.IpAddress = gofakeit.IPv4Address()
	server.Type = gofakeit.RandomString([]string{"mysql", "mssql", "postgres"})
	return server
}
