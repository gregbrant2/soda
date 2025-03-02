package mapping

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gregbrant2/soda/internal/api/dtos"
	"github.com/gregbrant2/soda/internal/domain/entities"
	"github.com/stretchr/testify/assert"
)

func TestMapServerSuccess(t *testing.T) {
	server := makeServerEntity()
	mapped := MapServer(server)
	assertServer(t, server, mapped)
}

func TestMapServersSuccess(t *testing.T) {
	servers := makeServerEntitySlice()
	mapped := MapServers(servers)
	assert.Equal(t, len(servers), len(mapped))
	assertServer(t, servers[0], mapped[0])
	assertServer(t, servers[1], mapped[1])
	assertServer(t, servers[2], mapped[2])
}

func TestMapNewServer(t *testing.T) {
	dto := makeNewServerDto()
	mapped := MapNewServer(dto)
	assert.Equal(t, dto.Name, mapped.Name)
	assert.Equal(t, dto.Type, mapped.Type)
	assert.Equal(t, dto.IpAddress, mapped.IpAddress)
	assert.Equal(t, dto.Port, mapped.Port)
	assert.Equal(t, dto.Username, mapped.Username)
	assert.Equal(t, dto.Password, mapped.Password)
}

func TestMapDatabaseSuccess(t *testing.T) {
	database := makeDatabaseEntity()
	mapped := MapDatabase(database)
	assertDatabase(t, database, mapped)
}

func TestMapDatabasesSuccess(t *testing.T) {
	databases := makeDatabaseEntitySlice()
	mapped := MapDatabases(databases)
	assert.Equal(t, len(databases), len(mapped))
	assertDatabase(t, databases[0], mapped[0])
	assertDatabase(t, databases[1], mapped[1])
	assertDatabase(t, databases[2], mapped[2])
}

func TestMapNewDatabase(t *testing.T) {
	dto := makeNewDatabaseDto()
	mapped := MapNewDatabase(dto)
	assert.Equal(t, dto.Name, mapped.Name)
	assert.Equal(t, dto.Server, mapped.Server)
}

func assertDatabase(t *testing.T, server entities.Database, mapped dtos.Database) {
	assert.Equal(t, server.Id, mapped.Id)
	assert.Equal(t, server.Name, mapped.Name)
	assert.Equal(t, server.Server, mapped.Server)
}

func assertServer(t *testing.T, server entities.Server, mapped dtos.Server) {
	assert.Equal(t, server.Id, mapped.Id)
	assert.Equal(t, server.Name, mapped.Name)
	assert.Equal(t, server.Type, mapped.Type)
	assert.Equal(t, server.Databases, mapped.Databases)
	assert.Equal(t, server.IpAddress, mapped.IpAddress)
	assert.Equal(t, server.Port, mapped.Port)
	assert.Equal(t, server.Status, mapped.Status)
}

func makeNewDatabaseDto() dtos.NewDatabase {
	s := dtos.NewDatabase{}
	gofakeit.Struct(&s)
	return s
}

func makeDatabaseEntity() entities.Database {
	s := entities.Database{}
	gofakeit.Struct(&s)
	return s
}

func makeDatabaseEntitySlice() []entities.Database {
	s := make([]entities.Database, 3)
	gofakeit.Struct(&s[0])
	gofakeit.Struct(&s[1])
	gofakeit.Struct(&s[2])
	return s
}

func makeNewServerDto() dtos.NewServer {
	s := dtos.NewServer{}
	gofakeit.Struct(&s)
	return s
}

func makeServerEntity() entities.Server {
	s := entities.Server{}
	gofakeit.Struct(&s)
	return s
}

func makeServerEntitySlice() []entities.Server {
	s := make([]entities.Server, 3)
	gofakeit.Struct(&s[0])
	gofakeit.Struct(&s[1])
	gofakeit.Struct(&s[2])
	return s
}
